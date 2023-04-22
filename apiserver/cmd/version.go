/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* https://www.sysadm.cn/licenses/apache-2.0.txt
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
*
 */

package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
	apiserverApp "github.com/wangyysde/sysadm/apiserver/app"
	"github.com/wangyysde/sysadmLog"
	"github.com/wangyysde/yaml"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of apiserver",
	Run: func(cmd *cobra.Command, args []string) {
		version := apiserverApp.GetVersion() 
		const flag = "format"
		if version != nil {
			of, err := cmd.Flags().GetString(flag)

			if err != nil {
				sysadmLog.Error("error accessing flag %s for command %s", version.Version, cmd.Name())
				return
			}

			switch of {
			case "":
				sysadmLog.Info("apiserver version %+v", version)
			case "short":
				sysadmLog.Info("apiserver version %+v", version)
			case "yaml":
				y, err := yaml.Marshal(&version)
				if err != nil {
					return
				}
				sysadmLog.Info("%s", string(y))
			case "json":
				y, err := json.MarshalIndent(&version, "", "  ")
				if err != nil {
					return
				}
				sysadmLog.Info(string(y))
			default:
				sysadmLog.Info("invalid output format: %s", of)
			}
		} else {
			sysadmLog.Error("can not get version inforation  %s", cmd.Name())
			return
		}

	},
	Args: cobra.NoArgs,
}

func init() {
	versionCmd.Flags().StringP("format", "f", "", "Output format; available options are 'yaml', 'json' and 'short'")

}
