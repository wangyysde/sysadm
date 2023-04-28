/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
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
	"sysadm/sysadmerror"
	"sysadm/registryctl/config"
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
	Use:   "registryctl",
	Short: "registryctl： registry contoller of sysadm registry",
	Long:  dedent.Dedent(`

	======================================================
	SYSADM
		registry contoller of sysadm registry

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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c",config.DefaultConfigFile, "specified config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var errs []sysadmerror.Sysadmerror
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		dir, err :=  filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(200000,"fatal","occurred an error:%s",err))
			logErrors(errs)
			os.Exit(200000)
		}


		viper.AddConfigPath(dir+"/conf")
		viper.SetConfigType("yaml")
		viper.SetConfigName("registryctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(200001,"info","Using config filer:%s",viper.ConfigFileUsed()))
		sysadmServer.Logf("info","Using config file: %s\n", viper.ConfigFileUsed())
	}
}


func logErrors(errs []sysadmerror.Sysadmerror){

	for _,e := range errs {
		l := sysadmerror.GetErrorLevelString(e)
		no := e.ErrorNo
		sysadmServer.Logf(l,"erroCode: %d Msg: %s",no,e.ErrorMsg)
	}
	
	return
}