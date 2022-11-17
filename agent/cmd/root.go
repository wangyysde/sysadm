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
	"os"
	"path/filepath"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wangyysde/sysadm/agent/app"
	"github.com/wangyysde/sysadm/config"
	"github.com/wangyysde/sysadmLog"
)

// Disable completion to agent package
var disableCompletion = cobra.CompletionOptions{
	DisableDefaultCmd: false,
	DisableNoDescFlag: false,
	DisableDescriptions: false,
	HiddenDefaultCmd: false,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "agent",
	Short: "agent is a management tool for managing node of kubernetes clusters",
	Long:  dedent.Dedent(`

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
		Version: "",
		Author: "",
		GitCommitId: gitCommitId,
		Branch: branchName,
		GitTreeStatus: gitTreeStatus,
		BuildDateTime: buildDateTime,
		GoVersion: goVersion,
		Compiler: compiler,
		Arch: arch,
		Os: hostos,
	}

	app.SetVersion(&version)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// specifing configuration file path.
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.CfgFile, "config", "c", app.CfgFile, "specified config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// specifies whether agent using TLS protocol when it is communicate with  a server (agent send the reponses message to the server)
	serverIsTls := rootCmd.PersistentFlags().BoolP("tls", "",app.DefaultServerIsTls, "specifies whether agent using TLS protocol when it is communicate with  a server (agent send the reponses message to the server)")
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
	skipVerifyCert := rootCmd.PersistentFlags().BoolP("skip-verify-cert", "",app.DefaultskipVerifyCert, "skipVerifyCert set whether check the certs which got from a server isvalid")
	app.CliOps.Global.Tls.InsecureSkipVerify = *skipVerifyCert

	// where the results of a command running will be send to. one of server: a server receiving the results; stdout, file
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.Output, "output", "o", "server", "where the results of a command running will be send to. one of server,stdout,file")

	//  the path of output file. this value must not empty if output be set to "file"
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.OutputFile, "outputfile", "of", "", `the path of output file. this value must not empty if output be set to "file"`)	

	//  logfile for agent which is used to log runing log messages of agent to 
	rootCmd.PersistentFlags().StringVarP(&app.CliOps.Global.LogFile, "log", "l", app.LogFile, `logfile for agent which is used to log runing log messages of agent to`)	

}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if app.CliOps.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(app.CliOps.CfgFile)
	} else {

		dir, err :=  filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			sysadmLog.Error("get absolute path error %s",err)
			os.Exit(1)
		}

		configPath := filepath.Join(dir,app.CfgFile)
		app.CliOps.CfgFile = configPath
		viper.SetConfigFile(configPath)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		sysadmLog.Info("Using config file: %s\n", viper.ConfigFileUsed())
	}
}
