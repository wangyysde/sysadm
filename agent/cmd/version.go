/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
* @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html
 */

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangyysde/yaml"
	"sysadm/agent/app"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of agent",
	Run: func(cmd *cobra.Command, args []string) {
		version := app.GetVersion()
		const flag = "format"
		if version != nil {
			of, err := cmd.Flags().GetString(flag)

			if err != nil {
				fmt.Printf("error accessing flag %s for command %s\n", version.Version, cmd.Name())
				return
			}

			switch of {
			case "":
				fmt.Printf("agent version %s\n", version.Version)
			case "short":
				fmt.Printf("agent version %s\n", version.Version)
			case "yaml":
				y, err := yaml.Marshal(&version)
				if err != nil {
					return
				}
				fmt.Printf("%s\n", string(y))
			case "json":
				y, err := json.MarshalIndent(&version, "", "  ")
				if err != nil {
					return
				}
				fmt.Printf("%s\n", string(y))
			default:
				fmt.Printf("invalid output format: %s\n", of)
			}
		} else {
			fmt.Printf("can not get version of %s", cmd.Name())
			return
		}

	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringP("format", "f", "", "Output format; available options are 'yaml', 'json' and 'short'")

}
