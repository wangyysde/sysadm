/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2021 Bzhy Network. All rights reserved.
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
    "fmt"
	"encoding/json"

    "github.com/spf13/cobra"
    "github.com/wangyysde/sysadmServer"
	"github.com/wangyysde/yaml"
)

type Version struct {
        ServerVersion string `json:"serverVersion"`
}

var versionCmd = &cobra.Command{
    Use: "version",
    Short: "Print the version of sysadm",
    Run: func(cmd *cobra.Command, args []string){
		const flag = "output"
		v := Version{
			ServerVersion: "v1.0.0",
		}

        of, err := cmd.Flags().GetString(flag)
        if err != nil {
			sysadmServer.Log(fmt.Sprintf("error accessing flag %s for command %s", v.ServerVersion , cmd.Name()),"error")
			return
		}
		switch of {
		case "":
			sysadmServer.Log(fmt.Sprintf("sysadm Server version %s", v.ServerVersion),"info")
		case "short":
			sysadmServer.Log(fmt.Sprintf("sysadm Server version %s",v.ServerVersion),"info")
		case "yaml":
			y, err := yaml.Marshal(&v)
            if err != nil {
                return 
            }
        	sysadmServer.Log(fmt.Sprintf("%s",string(y)),"info")
		case "json":
			y, err := json.MarshalIndent(&v, "", "  ")
            if err != nil {
                return 
            }
            sysadmServer.Log(fmt.Sprintf(string(y)),"info")
		default:
			sysadmServer.Log(fmt.Sprintf("invalid output format: %s", of),"error")
		}
	},
	Args: cobra.NoArgs,
}
        
func init(){
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
}



