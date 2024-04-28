/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
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
* Note: define root command
 */

package cmd

import (
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wangyysde/sysadmLog"
	apiserverApp "sysadm/apiserver/app"
	"sysadm/config"
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
	Use:   "apiserver",
	Short: "apiserver send command,collect command status and command logs for sysadm platform",
	Long: dedent.Dedent(`

	======================================================
	Apiserver
		sysadm ApiServer 

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

	apiserverApp.SetVersion(&version)
	cobra.OnInitialize(initConfig)

	// add version subcommand
	rootCmd.AddCommand(versionCmd)

	// add start subcommand
	rootCmd.AddCommand(startCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgFile := apiserverApp.GetCfgFile()
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		sysadmLog.Info("Using config file: %s\n", viper.ConfigFileUsed())
	}
}
