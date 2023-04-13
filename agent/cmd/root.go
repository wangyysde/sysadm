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
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wangyysde/sysadm/agent/app"
	"github.com/wangyysde/sysadm/config"
	"github.com/wangyysde/sysadmLog"
)

// Disable completion to agent package
var disableCompletion = cobra.CompletionOptions{
	DisableDefaultCmd:   false,
	DisableNoDescFlag:   false,
	DisableDescriptions: false,
	HiddenDefaultCmd:    false,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "agent",
	Short: "agent is a management tool for managing node of kubernetes clusters",
	Long: dedent.Dedent(`

	======================================================
	Agent
		agent is a management tool for managing node of kubernetes clusters

		Please give us feedback at:
		https://github.com/wangyysde/sysadm/issues/

	======================================================
	`),
	CompletionOptions: disableCompletion,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	version := config.Version{
		Version:       "",
		Author:        "",
		GitCommitId:   gitCommitId,
		Branch:        branchName,
		GitTreeStatus: gitTreeStatus,
		BuildDateTime: buildDateTime,
		GoVersion:     goVersion,
		Compiler:      compiler,
		Arch:          arch,
		Os:            hostos,
	}

	app.SetVersion(&version)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// specifing configuration file path.
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.CfgFile, "config", "c", "", "specified config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// specifies whether agent using TLS protocol when it is communicate with  a server (agent send the reponses message to the server)
	serverIsTls := rootCmd.PersistentFlags().BoolP("tls", "", app.DefaultServerIsTls, "specifies whether agent using TLS protocol when it is communicate with  a server (agent send the reponses message to the server)")
	app.CliOps.Global.Tls.IsTls = *serverIsTls

	// ca file which agent will used to connect to a server(agent send the reponses message to the server).
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.Tls.Ca, "ca", "", "", "ca file which agent will used to connect to a server")

	//  cert file which agent will used to connect to a server(agent send the reponses message to the server)
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.Tls.Cert, "cert", "", "", "cert file which agent will used to connect to a servers")

	// key file which agent will used to connect to a server(agent send the reponses message to the server)
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.Tls.Key, "keyfile", "", "", "key file which agent will used to connect to a server")

	// the address of a server which agent send the reponses message to
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.Server.Address, "server", "s", "", "the address of a server which agent send the reponses message to")

	// the port of a server listen which agent send the reponses message to
	serverPort := rootCmd.PersistentFlags().Int("port", app.DefaultServerPort, "the address of a server which agent send the reponses message to")
	app.CliOps.Global.Server.Port = *serverPort

	// skipVerifyCert set whether check the certs which got from a server is valid.  this value will be set to InsecureSkipVerify
	skipVerifyCert := rootCmd.PersistentFlags().BoolP("skip-verify-cert", "", app.DefaultskipVerifyCert, "skipVerifyCert set whether check the certs which got from a server isvalid")
	app.CliOps.Global.Tls.InsecureSkipVerify = *skipVerifyCert

	// the path of access log file
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.Log.AccessLog, "access-logfile", "", "", `the path of access log file`)

	// the path of error log file. both access log messages and error log messages will be log into access log file if error log file not set.
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.Log.ErrorLog, "error-logfile", "", "", `the path of error log file. both access log messages and error log messages will be log into access log file if error log file not set.`)

	// log message with the format(kind) will be output. its value is one of "text" and "json". default value is text
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.Log.Kind, "log-kind", "", app.DefaultLogKind, `log message with the format(kind) will be output. its value is one of "text" and "json". default value is text.`)

	// specifies log level. just the log messages will be output what the level of the log message is higher "logLevel".
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.Log.Level, "log-level", "", app.DefaultLogLevel, `specifies log level. just the log messages will be output what the level of the log message is higher "logLevel".`)

	// specifies whether agent running in Debug mode
	debugMode := rootCmd.PersistentFlags().BoolP("debug", "", app.DefalutDebugMode, "specifies whether agent running in Debug mode")
	app.CliOps.Global.DebugMode = *debugMode

	// specifies a identifer of the node which agent running on it.
	// It is any combination of the IP,HOSTNAME and MAC joined by commas  or a customize string what the leght of the string is less 63
	// agent will get all IPs without not active and reponse these IPs in list to the server by nodeIdentifer.IPs filed if IP is included in NodeIdentifer
	// agent will get hostname and reponse the hostname  to the server by nodeIdentifer.Hostname filed if hostname is included in NodeIdentifer
	// agent will get all MACs without not active and reponse these MACs in list to the server by nodeIdentifer.MACs filed if MAC is included in NodeIdentifer
	// customize string is reponse to the server directly .
	// customize string is conflicted with IP,HOSTNAME and MAC. the nodeIdentifer can be changed by the server during agent communicate with the server
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.NodeIdentifer, "node-identifer", "", app.DefaultNodeIdentifer, `It is any combination of the IP,HOSTNAME and MAC joined by commas  or a customize string what the leght of the string is less 63`)

	// specifies the uri where agent get commands to run when agent runing as daemon in passive mode.
	// agent will send the requests to "/" on the server if GetUri is empty.
	// Uri is the path where agent will send result message to when is running as command.
	// Uri is the listen path where agent receives commands to run when  agent runing as daemon in active mode.
	// the length of this value shoule less 63
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.Uri, "uri", "", "", `specifies uri path for agent get command from (daemon in passive), where listen on (daemon in active) or where agent send result to (run as CLI)`)

	// sourceIP specifies the source IP address which will be use to connect to a server by agent. this ip address must be configurated on one of the
	// interfaces  on the host where agent running on.  agent will get a source IP address from host  automatically if the value of this field is "".
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.SourceIP, "source", "", "", `the source IP address which will be use to connect to a server by agent.`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if app.CliOps.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(app.CliOps.CfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		sysadmLog.Info("Using config file: %s\n", viper.ConfigFileUsed())
	}
}
