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
	"os"
	"path/filepath"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wangyysde/sysadmServer"

	"sysadm/config"
	"sysadm/infrastructure/server"
)

var cfgFile string

// Disable completion to sysadm package
var disableCompletion = cobra.CompletionOptions{
	DisableDefaultCmd: true,
	DisableNoDescFlag: true,
	DisableDescriptions: true,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "infrastructure",
	Short: "infrastructure module of sysadm",
	Long:  dedent.Dedent(`

	======================================================
	Infrastructure
		An easily system administration platform

		Please give us feedback at:
		https://sysadm/issues/

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
	if server.CurrentRuningData != nil {
		server.LastRuningData = server.CurrentRuningData
		server.CurrentRuningData = &server.RuningData{}
	} else {
		server.CurrentRuningData = &server.RuningData{}
	}

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

	server.SetVersion(&version)
	server.CurrentRuningData.Config.Version = version
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "conf/infrastructure.yaml", "specified config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		dir, err :=  filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			sysadmServer.Logf("error","%s",err)
			os.Exit(1)
		}


		viper.AddConfigPath(dir+"/conf")
		viper.SetConfigType("yaml")
		viper.SetConfigName("infrastructure")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		sysadmServer.Logf("info","Using config file: %s\n", viper.ConfigFileUsed())
	}
}
