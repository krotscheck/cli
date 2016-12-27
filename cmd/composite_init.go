package main

import (
	"io"

	"fmt"
	"os"
	"path"

	"github.com/CloudCoreo/cli/cmd/content"
	"github.com/CloudCoreo/cli/cmd/util"
	"github.com/spf13/cobra"
)

type compositeInitCmd struct {
	out       io.Writer
	directory string
	serverDir bool
}

func newCompositeInitCmd(out io.Writer) *cobra.Command {
	compositeInit := &compositeInitCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   content.CmdInitUse,
		Short: content.CmdCompositeInitShort,
		Long:  content.CmdCompositeInitLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return compositeInit.run()
		},
	}

	f := cmd.Flags()

	f.StringVarP(&compositeInit.directory, content.CmdFlagDirectoryLong, content.CmdFlagDirectoryShort, "", content.CmdFlagDirectoryDescription)
	f.BoolVarP(&compositeInit.serverDir, content.CmdFlagServerLong, content.CmdFlagServerShort, false, content.CmdFlagServerDescription)

	return cmd
}

func (t *compositeInitCmd) run() error {

	if t.directory == "" {
		t.directory, _ = os.Getwd()
	}

	genContent(t.directory)

	if t.serverDir {
		genServerContent(t.directory)
	}

	return nil
}

func genContent(directory string) {
	if directory == "" {
		directory, _ = os.Getwd()
	}

	// config.yml file
	fmt.Println()
	util.CreateFile(content.DefaultFilesConfigYAMLName, directory, content.DefaultFilesConfigYAMLContent, false)

	// override folder
	util.CreateFolder(content.DefaultFilesOverrideFolderName, directory)

	overrideTree := fmt.Sprintf(content.DefaultFilesOverridesReadMeTree, content.DefaultFilesReadMeCodeTicks, content.DefaultFilesReadMeCodeTicks)

	overrideReadmeContent := fmt.Sprintf("%s%s%s", content.DefaultFilesOverridesReadMeHeader, overrideTree, content.DefaultFilesOverridesReadMeFooter)

	err := util.CreateFile(content.DefaultFilesReadMEName, path.Join(directory, content.DefaultFilesOverrideFolderName), overrideReadmeContent, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)

	}

	// services folder
	util.CreateFolder(content.DefaultFilesServicesFolder, directory)

	err = util.CreateFile(content.DefaultFilesConfigRBName, path.Join(directory, content.DefaultFilesServicesFolder), content.DefaultFilesConfigRBContent, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}

	servicesReadMeCode := fmt.Sprintf(content.DefaultFilesServicesReadMeCode, content.DefaultFilesReadMeCodeTicks, content.DefaultFilesReadMeCodeTicks)

	servicesReadMeContent := fmt.Sprintf("%s%s", content.DefaultFilesServicesReadMeHeader, servicesReadMeCode)

	err = util.CreateFile(content.DefaultFilesReadMEName, path.Join(directory+content.DefaultFilesServicesFolder), servicesReadMeContent, false)

	if err != nil {
		fmt.Println(err.Error())
	}

	if err == nil {
		fmt.Println(content.CmdCompositeInitSuccess)
	}
}

func genServerContent(directory string) {
	//operational scripts dir
	util.CreateFolder(content.DefaultFilesOperationalScriptsFolder, directory)

	// generate operational readme file
	err := util.CreateFile(content.DefaultFilesReadMEName, path.Join(directory, content.DefaultFilesOperationalScriptsFolder), content.DefaultFilesOperationalReadMeContent, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}

	//boot scripts dir
	util.CreateFolder(content.DefaultFilesBootScriptsFolder, directory)

	//README.md
	err = util.CreateFile(content.DefaultFilesReadMEName, path.Join(directory, content.DefaultFilesBootScriptsFolder), content.DefaultFilesBootReadMeContent, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}

	//order.yaml
	err = util.CreateFile(content.DefaultFilesOrderYAMLName, path.Join(directory, content.DefaultFilesBootScriptsFolder), content.DefaultFilesBootOrderYAMLContent, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}

	//shutdown scripts dir
	util.CreateFolder(content.DefaultFilesShutdownScriptsFolder, directory)

	//README.md
	err = util.CreateFile(content.DefaultFilesReadMEName, path.Join(directory, content.DefaultFilesShutdownScriptsFolder), content.DefaultFilesShutDownReadMeContent, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}

	//order.yaml
	err = util.CreateFile(content.DefaultFilesOrderYAMLName, path.Join(directory, content.DefaultFilesShutdownScriptsFolder), content.DefaultFilesShutDownOrderYAMLContent, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}
}
