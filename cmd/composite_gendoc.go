package main

import (
	"io"

	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudCoreo/cli/cmd/content"
	"github.com/CloudCoreo/cli/cmd/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var docHeaders = map[string]string{
	"head.md":        "",
	"description.md": "Description",
	"config.yaml":    "",
	"tags.md":        "Tags",
	"categories.md":  "Categories",
	"diagram.md":     "Diagram",
	"icon.md":        "Icon",
	"hierarchy.md":   "Hierarchy",
}

var docOrder = []string{"head.md", "description.md", "hierarchy.md", "config.yaml", "tags.md", "categories.md", "diagram.md", "icon.md"}

type compositeGendocCmd struct {
	out       io.Writer
	directory string
	serverDir bool
}

func newCompositeGendocCmd(out io.Writer) *cobra.Command {
	compositeGendoc := &compositeGendocCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   content.CmdGendocUse,
		Short: content.CmdCompositeGendocShort,
		Long:  content.CmdCompositeGendocLong,
		Run: func(cmd *cobra.Command, args []string) {
			compositeGendoc.ProcessCmdGendocUse(args)
		},
	}

	f := cmd.Flags()

	f.StringVarP(&compositeGendoc.directory, content.CmdFlagDirectoryLong, content.CmdFlagDirectoryShort, "", content.CmdFlagDirectoryDescription)
	f.BoolVarP(&compositeGendoc.serverDir, content.CmdFlagServerLong, content.CmdFlagServerShort, false, content.CmdFlagServerDescription)

	return cmd
}

func (t *compositeGendocCmd) run() error {

	if t.directory == "" {
		t.directory, _ = os.Getwd()
	}

	genContent(t.directory)

	if t.serverDir {
		genServerContent(t.directory)
	}

	return nil
}

//ProcessCmdGendocUse Process Cmd GenDoc
func (t *compositeGendocCmd) ProcessCmdGendocUse(args []string) {
	util.CheckArgsCount(args)

	if t.directory == "" {
		t.directory, _ = os.Getwd()
	}

	var readmeFileContent bytes.Buffer

	for index := range docOrder {

		fileName := docOrder[index]

		if fileName == content.DefaultFilesConfigYAMLName {
			configFileContent, _ := generateConfigContent(path.Join(t.directory, fileName))
			readmeFileContent.WriteString(configFileContent)
		} else {

			fileContent, err := ioutil.ReadFile(path.Join(t.directory, fileName))
			if err != nil {
				fmt.Println(fmt.Sprintf(content.ErrorMissingFile, fileName))
				err := util.CreateFile(fileName, t.directory, "", false)
				if err != nil {
					fmt.Fprintf(os.Stderr, err.Error())
					os.Exit(-1)
				}

			}

			// create headers when non empty
			if docHeaders[fileName] != "" {
				readmeFileContent.WriteString(fmt.Sprintf("## %s\n", docHeaders[docOrder[index]]))
			}

			readmeFileContent.WriteString(string(fileContent) + "\n\n")
		}
	}

	err := util.CreateFile(content.DefaultFilesReadMEName, t.directory, readmeFileContent.String(), true)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)

	}

	fmt.Println(content.CmdCompositeGendocSuccess)
}

// YamlConfig struct for parsing config.yaml file
type YamlConfig struct {
	Variables yaml.MapSlice
}

type varOption struct {
	key           string
	description   string
	required      bool
	valueType     string      `schema:"type"`
	defaultValues interface{} `schema:"default"`
}

func generateConfigContent(configFilePath string) (string, error) {

	filename, err := filepath.Abs(configFilePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}

	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, "Could not read "+filename)
	}

	var config YamlConfig

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not parse YAML")
		os.Exit(-1)
	}

	missingRequiredContent, err := generateVariablesContent(
		config,
		func(required bool, defaultValue interface{}) bool {
			return required && defaultValue == nil
		},
		fmt.Sprintf("\n## %s\n\n", content.DefautlFilesGenDocReadMeRquiredNoDefaultHeader),
		false)

	requiredContent, err := generateVariablesContent(
		config,
		func(required bool, defaultValue interface{}) bool {
			return required && defaultValue != nil
		},
		fmt.Sprintf("\n## %s\n\n", content.DefautlFilesGenDocReadMeRquiredDefaultHeader),
		true)

	notRequiredDefaultContent, err := generateVariablesContent(
		config,
		func(required bool, defaultValue interface{}) bool {
			return !required && defaultValue != nil
		},
		fmt.Sprintf("\n## %s\n\n", content.DefautlFilesGenDocReadMeNoRquiredDefaultHeader),
		true)

	theRestContent, err := generateVariablesContent(
		config,
		func(required bool, defaultValue interface{}) bool {
			return !required && defaultValue == nil
		},
		fmt.Sprintf("\n## %s\n\n", content.DefautlFilesGenDocReadMeNoRquiredNoDefaultHeader),
		true)

	return missingRequiredContent + requiredContent + notRequiredDefaultContent + theRestContent, err

}

func generateVariablesContent(config YamlConfig, check func(bool, interface{}) bool, header string, printVar bool) (string, error) {

	var contentBytes bytes.Buffer
	counter := 0
	contentBytes.WriteString(header)

	// loop over mapslice and create option object. TODO: Better way would be to use reflection instead.
	for _, variable := range config.Variables {
		option := varOption{}
		option.key = variable.Key.(string)
		for _, o := range variable.Value.(yaml.MapSlice) {
			if o.Value != nil {
				switch o.Key {
				case "description":
					option.description = o.Value.(string)
				case "required":
					option.required = o.Value.(bool)
				case "type":
					option.valueType = o.Value.(string)
				case "default":
					option.defaultValues = o.Value
				}
			}
		}

		// check if
		if check(option.required, option.defaultValues) {
			counter++
			contentBytes.WriteString("### `" + option.key + "`:\n")
			contentBytes.WriteString("  * description: " + option.description)

			if printVar && option.defaultValues != nil {
				switch strings.ToLower(option.valueType) {
				case "array":
					// Convert to string[] and then join items with ,
					defaultValues := convertStringSlice(option.defaultValues)

					contentBytes.WriteString("\n  * default: ")
					contentBytes.WriteString(fmt.Sprint(strings.Join(defaultValues, ", ")))
				case "boolean":
					contentBytes.WriteString("\n  * default: ")
					contentBytes.WriteString(strconv.FormatBool(option.defaultValues.(bool)))
				case "case":
					contentBytes.WriteString("\n  * default: ")
					contentBytes.WriteString(fmt.Sprint(option.defaultValues))
				case "number":
					contentBytes.WriteString("\n  * default: ")
					contentBytes.WriteString(fmt.Sprint(option.defaultValues))
				case "string":
					contentBytes.WriteString("\n  * default: ")
					contentBytes.WriteString(fmt.Sprint(option.defaultValues))

				default:
					// if unknown type is provided try to cast it to string
					if c, ok := option.defaultValues.(string); ok {
						if strings.Contains(c, "\n") {
							contentBytes.WriteString("\n  * default: \n" + content.DefaultFilesReadMeCodeTicks + "\n")
							contentBytes.WriteString(option.defaultValues.(string) + "\n")
							contentBytes.WriteString(content.DefaultFilesReadMeCodeTicks)
						} else {
							contentBytes.WriteString("\n  * default: ")
							contentBytes.WriteString(option.defaultValues.(string) + "\n")
						}
					}
				}
			}

			contentBytes.WriteString("\n\n")
		}
	}

	if counter == 0 {
		contentBytes.WriteString("**None**\n\n")
	}

	return contentBytes.String(), nil
}

func convertStringSlice(slice interface{}) []string {

	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]string, s.Len())

	for i := 0; i < s.Len(); i++ {

		switch s.Index(i).Interface().(type) {
		case string:
			ret[i] = s.Index(i).Interface().(string)
		case int:
			ret[i] = strconv.FormatInt(int64(s.Index(i).Interface().(int)), 16)
		case bool:
			ret[i] = strconv.FormatBool(s.Index(i).Interface().(bool))
		}
	}

	return ret
}
