/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
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
*
* NOTE:
* defined some functions are related to handle configurations. 
*/

package app

import(
	"os"
	"strings"
	"path/filepath"

	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/config"
)

func handlerConfig()(bool,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020007, "debug", "try to handle configurations for apiserver"))
	ok, err := handleNotInConfFile()
	errs = append(errs, err...)
	if !ok {
		return false, errs
	}


}

// HandleNotInConfFile handler the configuration items which can not define in define file,such as working dir, configuration file path.
func handleNotInConfFile() (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	confFile := runData.runConf.ConfFile
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020001, "debug", "try to get working dir"))
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020002, "fatal", "get working dir error %s", err))
		return false, errs
	}
	workingDir := filepath.Join(binPath, "../")
	runData.workingRoot = filepath.Join(binPath, "../")

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020003, "debug", "checking configuration file path"))
	var cfgFile string = ""
	if confFile != "" {
		if filepath.IsAbs(confFile) {
			fp, err := os.Open(confFile)
			if err != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020004, "fatal", "can not open configuration file %s error %s", confFile, err))
				return false, errs
			}
			fp.Close()
			cfgFile = confFile
		} else {
			configPath := filepath.Join(workingDir,confFile)
			fp, err := os.Open(configPath)
			if err != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020005, "fatal", "can not open configuration file %s error %s", configPath, err))
				return false, errs
			}
			fp.Close()
			cfgFile = configPath
		}
	} else {
		// try to get configuration file from default path
		configPath := filepath.Join(workingDir, confFilePath)
		fp, err := os.Open(configPath)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020006, "fatal", "can not open configuration file %s error %s", configPath, err))
			return false, errs
		}
		fp.Close()
		cfgFile = configPath
	}

	runData.runConf.ConfFile = cfgFile

	return true, errs

}

// set version data to runData instance
func SetVersion(version *config.Version) {
	if version == nil {
		return
	}

	version.Version = appVer
	version.Author = appAuthor

	runData.runConf.Version = *version
}

// get version data from runData instance
func GetVersion() *config.Version {
	if runData.runConf.Version.Version != "" {
		return &runData.runConf.Version 
	}

	return nil
}

// return the configuration file path of the application from runData
func GetCfgFile() string {
	return strings.TrimSpace(runData.runConf.ConfFile)
}

// set configuration file path what has got from CLI flag to runData
func SetCfgFile(cfgFile string){
	cfgFile = strings.TrimSpace(cfgFile)
	if cfgFile == "" {
		cfgFile = confFilePath
	}

	runData.runConf.ConfFile = cfgFile
}

