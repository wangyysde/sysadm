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
	"github.com/spf13/cobra"
	"sysadm/agent/app"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start agent as daemon",
	Run: func(cmd *cobra.Command, args []string) {
		app.Start(cmd, args)
	},
	Args: cobra.NoArgs,
}

func init() {

	rootCmd.AddCommand(startCmd)

	// IP address or hostname of apiServer where agent will be connected to
	startCmd.PersistentFlags().StringVarP(&app.RunData.Address, "apiserver", "a", "", "IP address or hostname of apiServer where agent will be connected to.")

	// service port of apiServer where agent will be connected to
	port := startCmd.PersistentFlags().IntP("port", "p", 0, "service port of apiServer where agent will be connected to.")
	app.RunData.Port = *port

	// whether using TLS when agent connect to apiServer
	enableTls := startCmd.PersistentFlags().BoolP("enable-tls", "e", false, "whether using TLS when agent connect to apiServer.")
	app.RunData.IsTls = *enableTls

	// Path to a cert file for the certificate authority
	startCmd.PersistentFlags().StringVarP(&app.RunData.Ca, "certificate-authority", "", "", "Path to a cert file for the certificate authority.")

	// Path to a client certificate file for TLS
	startCmd.PersistentFlags().StringVarP(&app.RunData.Cert, "client-certificate", "", "", "Path to a client certificate file for TLS.")

	// Path to a client key file for TLS
	startCmd.PersistentFlags().StringVarP(&app.RunData.Key, "client-key", "", "", "Path to a client key file for TLS.")

	// If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
	insecureSkipVerify := startCmd.PersistentFlags().BoolP("insecure-skip-tls-verify=false", "", false, "If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure.")
	app.RunData.InsecureSkipVerify = *insecureSkipVerify

	// enable debug mode
	debug := startCmd.PersistentFlags().BoolP("debug", "", false, "enable debug mode.")
	app.RunData.Debug = *debug

	// Path to log file. default is /var/log/agent.log
	startCmd.PersistentFlags().StringVarP(&app.RunData.LogFile, "log-file", "", "", "Path to log file. default is /var/log/agent.log.")
}
