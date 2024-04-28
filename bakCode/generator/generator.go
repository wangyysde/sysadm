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
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	execute()
}

func doCommand(cmd *cobra.Command, args []string) {
	generationKind = strings.ToLower(strings.TrimSpace(generationKind))
	if generationKind == "" {
		generationKind = defaultGenerationKind
	}

	generatedFileName = strings.ToLower(strings.TrimSpace(generatedFileName))
	if generatedFileName == "" {
		generatedFileName = defaultGeneratedFileName
	}
	dirs = strings.TrimSpace(dirs)
	if dirs == "" {
		fmt.Printf("the directories where will be scanned must be specified.")
		os.Exit(1)
	}

	copyrightPath, e := filepath.Abs(filepath.Dir(copyrightFile))
	if e != nil {
		fmt.Printf("get copyright file path error: %s\n", e)
		os.Exit(1)
	}
	copyrightFile = copyrightPath

	switch generationKind {
	case "conversion":
		e := generateConversionFiles()
		if e != nil {
			fmt.Printf("generate conversion files error: %s\n", e)
			os.Exit(1)
		}
	default:
		fmt.Printf("the kind %s is not valid\n", generationKind)
		os.Exit(1)
	}

	fmt.Printf("files have be generated\n")
	os.Exit(0)
}
