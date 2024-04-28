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
	"os"
	"path/filepath"
	"strings"
	sysadmApiServer "sysadm/apiserver/app"
	"sysadm/config"
	"sysadm/utils"
)

func SetVersion(version *config.Version) {
	if version == nil {
		return
	}

	version.Version = ver
	version.Author = author

	RunData.version = *version
}

func GetVersion() *config.Version {
	if RunData.version.Version != "" {
		return &RunData.version
	}

	return nil
}

func validateConf() error {
	binPath, e := filepath.Abs(filepath.Dir(os.Args[0]))
	if e != nil {
		return e
	}

	RunData.workingDir = filepath.Join(binPath, "../")
	_, e = utils.ValidateAddress(RunData.Address, false)
	if e != nil {
		return e
	}

	if RunData.IsTls {
		if RunData.Port == 0 {
			RunData.Port = sysadmApiServer.DefaultTlsPort
		}

		caPath, e := utils.CheckFileIsReadable(RunData.Ca, RunData.workingDir)
		if e != nil {
			return e
		}
		RunData.Ca = caPath

		certPath, e := utils.CheckFileIsReadable(RunData.Cert, RunData.workingDir)
		if e != nil {
			return e
		}
		RunData.Cert = certPath

		keyPath, e := utils.CheckFileIsRead(RunData.Key, RunData.workingDir)
		if e != nil {
			return e
		}
		RunData.Key = keyPath
	} else {
		if RunData.Port == 0 {
			RunData.Port = sysadmApiServer.DefaultPort
		}
		RunData.Ca = ""
		RunData.Cert = ""
		RunData.Key = ""
	}

	logFile := strings.TrimSpace(RunData.LogFile)
	if logFile == "" {
		logFile = defaultLogFile
	}
	tmpLogFile, e := utils.CheckFileWritable(logFile, RunData.workingDir, true, true)
	if e != nil {
		return e
	}
	RunData.LogFile = tmpLogFile

	return nil
}
