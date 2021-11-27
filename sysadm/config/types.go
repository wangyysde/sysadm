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

//Defining server configuration
type Server struct {
	Address string `json:"address"`
	Port int `json:"port"`
	Socket string `json:"socket"`
}

//Defining log configuration 
type Log struct {
	AccessLog string `json:"accessLog"`
	ErrorLog string `json:"errorLog"`
	Kind string `json:"kind"`
	Level string `json:"level"`
	SplitAccessAndError bool `json:"splitAccessAndError"`
	TimeStampFormat string `json:"timeStampFormat"`
}

type User struct {
	DefaultUser string `json:"defaultUser"`
	DefaultPassword string `json:"defaultPassword"`
}

type Config struct {
	Version string `json:"version"`
	Server Server `json:"server"`
	Log Log `json:"log"`
	User User `json:"user"`
}

var DefinedConfig Config = Config{}