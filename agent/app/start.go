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

package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
}

func Start(cmd *cobra.Command, args []string) {
	e := validateConf()
	if e != nil {
		log.Fatal("parameter is not valid %s", e)
		os.Exit(-1)
	}

	logFile, e := os.OpenFile(RunData.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if e != nil {
		log.WithFields(log.Fields{"logFile": RunData.LogFile,
			"errorMsg": fmt.Sprintf("%s", e)}).Error("open log file error")
	} else {
		defer logFile.Close()
		log.Info("log file %s has be opened, log message will be log into the file", RunData.LogFile)
		log.SetOutput(logFile)
	}
	if RunData.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	e = startLoop()
	if e != nil {
		log.Errorf("%s", e)
		os.Exit(-1)
	}

	os.Exit(0)
}
