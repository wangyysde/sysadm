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

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wangyysde/yaml"
)

// Deinfining default value of configuration file
var DefaultConfigFile = "conf/config.yaml"

// Getting the path of configurationn file and try to open it
// configPath: the path of configuration file user specified
// cmdRunPath: args[0]
func GetConfigPath(configPath string, cmdRunPath string) (string, error) {
	
	dir ,error := filepath.Abs(filepath.Dir(cmdRunPath))
	if error != nil {
		return "",error
	}

	if configPath == "" {
		configPath = filepath.Join(dir,".../")
		configPath = filepath.Join(configPath, DefaultConfigFile)
		fp, err := os.Open(configPath)
		if err != nil {
			return "",err
		}

		fp.Close()
		return configPath,nil
	}

	if ! filepath.IsAbs(configPath) {
		tmpDir := filepath.Join(dir,"../")
		configPath = filepath.Join(tmpDir,configPath)
	}

	fp, err := os.Open(configPath)
	if err != nil {
		return "",err
	}
	fp.Close()

	return configPath,nil
}

// Reading the content of configuration from configPath and parsing the content 
// returning a pointer to Config if it is successfully parsed
// Or returning an error and nil
func GetConfigContent(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("The configration file path must not empty")
	}

	yamlContent, err := ioutil.ReadFile(configPath) 
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, &DefinedConfig)
	if err != nil {
		return nil, err
	}

	return &DefinedConfig, nil
}