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

package main

import (
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
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
	Use:   "generator",
	Short: "this is a tool for generating go code for apiserver",
	Long: dedent.Dedent(`

	======================================================
	Generator
		this is a tool for generating go code for apiserver

		Please give us feedback at:
		https://sysadm.cn/issues/

	======================================================
	`),
	Run: func(cmd *cobra.Command, args []string) {
		doCommand(cmd, args)
	},
	Args:              cobra.NoArgs,
	CompletionOptions: disableCompletion,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&generationKind, "kind", "k", "", "whick kind code will be generated. one of conversion...... default value is conversion")
	rootCmd.PersistentFlags().StringVarP(&dirs, "dirs", "d", "", "the directories where will be scanned.this value is a string with comma separated")
	rootCmd.PersistentFlags().StringVarP(&generatedFileName, "file", "f", "", "generated file name")
	rootCmd.PersistentFlags().StringVarP(&copyrightFile, "copyright", "c", "", "copyright file path")
}
