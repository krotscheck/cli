// Copyright © 2016 Paul Allen <paul@cloudcoreo.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"os"

	"github.com/CloudCoreo/cli/client"
	"github.com/CloudCoreo/cli/cmd/content"
	"github.com/CloudCoreo/cli/cmd/util"
	"github.com/spf13/cobra"
)

// CloudListCmd represents the based command for cloud subcommands
var CloudListCmd = &cobra.Command{
	Use:   content.CmdListUse,
	Short: content.CmdCloudListShort,
	Long:  content.CmdCloudListLong,
	PreRun: func(cmd *cobra.Command, args []string) {
		util.CheckArgsCount(args)
		SetupCoreoCredentials()
		SetupCoreoDefaultTeam()

	},
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.MakeClient(key, secret, content.EndpointAddress)
		if err != nil {
			util.PrintError(err, json)
			return
		}

		t, err := c.GetCloudAccounts(context.Background(), teamID)
		if err != nil {
			util.PrintError(err, json)
			os.Exit(-1)
		}

		b := make([]interface{}, len(t))
		for i := range t {
			b[i] = t[i]
		}

		util.PrintResult(b, []string{"ID", "Name", "TeamID"}, json)
	},
}

func init() {
	CloudCmd.AddCommand(CloudListCmd)
}