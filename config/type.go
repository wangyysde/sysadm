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
	// version of application. 
	Version string  `form:"version" json:"version" yaml:"version" xml:"version"` 

	// author of application s
	Author string  `form:"author" json:"author" yaml:"author" xml:"author"` 

	// commit ID of the application when the application built based on
	GitCommitId string  `form:"gitCommitId" json:"gitCommitId" yaml:"gitCommitId" xml:"gitCommitId"`  

	// branch name of the application which the application build based on
	Branch	string `form:"branch" json:"branch" yaml:"branch" xml:"branch"`  

	// git status of the branch when the application build based on
	GitTreeStatus string `form:"gitTreeStatus" json:"gitTreeStatus" yaml:"gitTreeStatus" xml:"gitTreeStatus"`  

	// the build time of the application 
	BuildDateTime string `form:"buildDateTime" json:"buildDateTime" yaml:"buildDateTime" xml:"buildDateTime"` 

	// go version which used to build the application 
	GoVersion string `form:"goVersion" json:"goVersion" yaml:"goVersion" xml:"goVersion"`  

	// compiler which used to build the application 
	Compiler string  `form:"compiler" json:"compiler" yaml:"compiler" xml:"compiler"` 

	// architecture the application was build based on 
	Arch string `form:"arch" json:"arch" yaml:"arch" xml:"arch"` 

	// OS name the application was build based on 
	Os string `form:"os" json:"os" yaml:"os" xml:"os"`  

}

//Defining server configuration
type Server struct {
	Address string `form:"address" json:"address" yaml:"address" xml:"address"` 
	Port int  `form:"port" json:"port" yaml:"port" xml:"port"` 
	Socket string `form:"socket" json:"socket" yaml:"socket" xml:"socket"` 
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
	AccessLog string `form:"accessLog" json:"accessLog" yaml:"accessLog" xml:"accessLog"` 
	// descriptor of access log file which will be used to close logger when system exit
	AccessLogFp *os.File
	ErrorLog string `form:"errorLog" json:"errorLog" yaml:"errorLog" xml:"errorLog"`  
	// descriptor of error log file which will be used to close logger when system exit
	ErrorLogFp *os.File
	Kind string `form:"kind" json:"kind" yaml:"kind" xml:"kind"`  
	Level string `form:"level" json:"level" yaml:"level" xml:"level"` 
	SplitAccessAndError bool `form:"splitAccessAndError" json:"splitAccessAndError" yaml:"splitAccessAndError" xml:"splitAccessAndError"`  
	TimeStampFormat string `form:"timeStampFormat" json:"timeStampFormat" yaml:"timeStampFormat" xml:"timeStampFormat"` 
}

type User struct {
	UserName string `form:"userName" json:"userName" yaml:"userName" xml:"userName"`   
	Password string `form:"password" json:"password" yaml:"password" xml:"password"`  
}