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

package config

import(
	"os"
)

type Version struct {
	Version string `json:"version"`
	Author string `json:"author"`
	GitCommitId string `json:"gitCommitId"`
	Branch	string `json:"branch"`
	GitTreeStatus string `json:"gitTreeStatus"`
	BuildDateTime string `json:"buildDateTime"`
	GoVersion string `json:"goVersion"`
	Compiler string `json:"compiler"`
	Arch string `json:"arch"`
	Os string `json:"os"`
}

//Defining server configuration
type Server struct {
	Address string `json:"address"`
	Port int `json:"port"`
	Socket string `json:"socket"`
}

//Define tls structure
type Tls struct {
	IsTls bool `form:"isTls" json:"isTls" yaml:"isTls" xml:"isTls"` 
	Ca string  `form:"ca" json:"ca" yaml:"ca" xml:"ca"`  
	Cert string `form:"cert" json:"cert" yaml:"cert" xml:"cert"`   
	Key string  `form:"key" json:"key" yaml:"key" xml:"key"`   
	InsecureSkipVerify bool  `form:"insecureSkipVerify" json:"insecureSkipVerify" yaml:"insecureSkipVerify" xml:"insecureSkipVerify"` 
}

//Defining log configuration 
type Log struct {
	AccessLog string `json:"accessLog"`
	// descriptor of access log file which will be used to close logger when system exit
	AccessLogFp *os.File
	ErrorLog string `json:"errorLog"`
	// descriptor of error log file which will be used to close logger when system exit
	ErrorLogFp *os.File
	Kind string `json:"kind"`
	Level string `json:"level"`
	SplitAccessAndError bool `json:"splitAccessAndError"`
	TimeStampFormat string `json:"timeStampFormat"`
}

type User struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}