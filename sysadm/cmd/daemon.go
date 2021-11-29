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
	//    "fmt"
	//"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/wangyysde/sysadm/sysadm/server"
	"github.com/wangyysde/sysadm/sysadm/config"
	//    "github.com/wangyysde/sysadmServer"
	//	"github.com/wangyysde/yaml"
)

// define daemon sub-command
var daemonCmd = &cobra.Command{
    Use: "daemon",
    Short: "start restart stop or reload the daemon of sysadm",
  //  Run: func(cmd *cobra.Command, args []string){
//		sysadmServer.Log(fmt.Sprintf("cfg:%s",cfgFile),"info")
//	},
	Args: cobra.NoArgs,
}

// define start sub-command of daemon
var startCmd =  &cobra.Command{
	Use: "start",
	Short: "Start the daemon of sysadm server",
	Run: func(cmd *cobra.Command, args []string){
		server.StartData.OldConfigPath = ""
		server.StartData.ConfigPath = cfgFile
		config.Version = versionStr
		server.DaemonStart(cmd, os.Args[0])
	},
	Args: cobra.NoArgs,

}
// define restart sub-command of daemon
var restartCmd = &cobra.Command{
	Use: "restart",
	Short: "Restart the daemon of sysadm server",
	Run: func(cmd *cobra.Command, args []string){},
	Args: cobra.NoArgs,
}

// define stop sub-command of daemon commands
var stopCmd = &cobra.Command{
	Use: "stop",
	Short: "Stop the daemon of sysadm",
	Run: func(cmd *cobra.Command, args []string){},
	Args: cobra.NoArgs,
}

// define reload sub-command of daemon command
var reloadCmd = &cobra.Command{
	Use: "reload",
	Short: "Reload the daemon of sysadm",
	Run: func(cmd *cobra.Command, args []string){},
	Args: cobra.NoArgs,
}

func init(){
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.AddCommand(startCmd)
	daemonCmd.AddCommand(restartCmd)
	daemonCmd.AddCommand(stopCmd)
	daemonCmd.AddCommand(reloadCmd)
}



