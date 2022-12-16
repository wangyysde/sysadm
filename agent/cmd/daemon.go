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

	"github.com/spf13/cobra"
	"github.com/wangyysde/sysadm/agent/app"
)

var daemonCmd = &cobra.Command{
    Use: "daemon",
    Short: "start agent as daemon",
    Run: func(cmd *cobra.Command, args []string){
		app.Daemon(cmd, args)
	},
	Args: cobra.NoArgs,
}
        
func init(){

	rootCmd.AddCommand(daemonCmd)

	// the method of getting commands by agent. agent gets commands from the server periodically and run them if this value is true
	// otherwise the server send a command to a agent when it  want the agent to run the command on a host. 
	
	passive := rootCmd.PersistentFlags().BoolP("passive", "p",app.DefaultPassive, "the method of getting commands by agent. agent gets commands from the server periodically and run them if this value is true. otherwise the server send a command to a agent when it  want the agent to run the command on a host.")
	app.CliOps.Agent.Passive =  *passive

	// the ca file for agent when agent running as daemon. 
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Agent.Tls.Ca, "agent-ca", "", "", "the ca file for agent when agent running as daemon.")

	// the cert file  for agent when agent running as daemon.
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Agent.Tls.Cert, "agent-cert", "", "", "the cert file  for agent when agent running as daemon.")

	// the  key file  for agent when agent running as daemon. 
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Agent.Tls.Key, "agent-key", "", "", "the  key file  for agent when agent running as daemon.")

	// listen address of agent using when agent running as daemon. 
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Agent.Server.Address, "listen", "a", "", "listen address of agent using when agent running as daemon.")

	// listen port of agent using when agent running as daemon.
	port := rootCmd.PersistentFlags().Int("listen-port", app.DefaultListenPort, "listen port of agent using when agent running as daemon.")
	app.CliOps.Agent.Server.Port = *port

	// period(second) for agent gets command from server when agent running as passive 
	period := rootCmd.PersistentFlags().Int("period", app.DefaultPeriod, "period(second) for agent gets command from server when agent running as passive")
	app.CliOps.Agent.Period = *period

	// insecret specifies whether agent listen on a insecret port when it is runing as daemon
	insecret := rootCmd.PersistentFlags().BoolP("insecret", "i",app.DefaultInsecret, "insecret specifies whether agent listen on a insecret port when it is runing as daemon")
	app.CliOps.Agent.Insecret = *insecret

	// insecret listen port of agent listening when it is running ad daemon  
	insecretPort := rootCmd.PersistentFlags().Int("insecret-port", app.DefaultInsecretPort, "insecret listen port of agent listening when it is running ad daemon")
	app.CliOps.Agent.InsecretPort = *insecretPort

}
