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
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/wangyysde/sysadmServer"
	"github.com/wangyysde/yaml"
)

// Deinfining default value of configuration file
var DefaultConfigFile = "conf/config.yaml"
var SupportVers = [...]string{"0.1", "0.2","21.0.0"}
var Version = ""
var DefaultIP = "0.0.0.0"
var DefaultPort = 8080
var DefaultSocket = "/var/run/sysadm.sock"
var DefaultAccessLog = "logs/sysadmAccess.log"
var DefaultErrorLog = "logs/sysadmError.log"
var DefaultLogKind = "text"
var DefaultLogLevel = "debug"
var DefaultUser = "admin"
var DefaultPasswd = "Sysadm12345"


var ConfigDefined Config = Config{}
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

	err = yaml.Unmarshal(yamlContent, &ConfigDefined)
	if err != nil {
		return nil, err
	}

	return &ConfigDefined, nil
}

// check the supporting of the version of configuration file specified
// return nil if the version be supported
// Or return err 
func CheckVerIsValid(ver string) error {
	 found := false
	for _,value := range SupportVers {
		if strings.ToLower(ver) == strings.ToLower(value) {
			found = true
			break
		}
	}

	if found {
		return nil
	}

	return fmt.Errorf("The version(%s) of the configuration file specified was not be supported by this release.\n",ver)
}

// check the validity of IP address 
// return IP(net.IP) if the ip address is valid
// Or return nil with error
func CheckIpAddress(address string) (net.IP, error) {
	if len(address) < 1 {
		return nil, fmt.Errorf("The address(%s) is empty or the length of it is less 1",address)
	}

	if ip := net.ParseIP(address); ip != nil {
		if address == "0.0.0.0" || address == "::" {
			return ip, nil
		}

		adds,err := net.InterfaceAddrs()
		if err != nil {
			return nil, fmt.Errorf("Get interface address error: %s",err)
		}

		for _,v := range adds {
			ipnet,ok := v.(*net.IPNet)
			if !ok {
				continue
			}
			if ip.Equal(ipnet.IP) {
				return ip, nil
			}
		}

		return nil, fmt.Errorf("The address(%s) is not any of the addresses of host interfaces.",address)
	}

	ips,err := net.LookupIP(address)
	if err != nil {
		return nil , fmt.Errorf("Lookup the IP of address(%s) error.",err)
	}

	adds,err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("Get interface address error: %s",err)
	}

	for _,ip := range ips {
		for _,v := range adds {
			ipnet,ok := v.(*net.IPNet)
			if !ok {
				continue
			}
			if ip.Equal(ipnet.IP) {
				return ip, nil
			}
		}
	}
	
	return nil, fmt.Errorf("The IP(%v) to the address(hostname:%s) is not any the IP address of host interfaces.",ips,address)
}

// check the validity of port 
// return port with nil if the port is valid 
// Or return 0 with error
func checkPort(port int)(int, error){
	if port > 1024 && port <= 65535 {
		return port,nil
	}

	return 0, fmt.Errorf("The port should be great than 1024 and less than 65535. Now is :%d\n",port)
}

// Converting relative path to absolute path of  file(such as socket, accesslog, errorlog) and return the  file path
// return "" and error if  file can not opened .
// Or return string and nil.
func getFile(f string,cmdRunPath string)(string,error){
	dir ,error := filepath.Abs(filepath.Dir(cmdRunPath))
	if error != nil {
		return "",error
	}

	if ! filepath.IsAbs(f) {
		tmpDir := filepath.Join(dir,"../")
		f = filepath.Join(tmpDir,f)
	}

	fp, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR,7)
	if err != nil {
		return "",err
	}
	fp.Close()

	return f,nil
}

// check the validity of log Level.
// The default level will be return if the value of level is empty or invalid.
// Or the value of level and nil will be returned
func checkLogLevel(level string) (string, error) {
	if len(level) < 1 {
		return defaultConfig.Log.Level, fmt.Errorf("Level is empty,default level has be set")
	}

	for _,l := range sysadmServer.Log.Level {
		if strings.ToLower(level) == strings.ToLower(l) {
			return strings.ToLower(level),nil
		}
	}

	return defaultConfig.Log.Level,fmt.Errorf("Level(%s) was not found,default level has be set.",level)
}

// check the validity of log format.
// The default format will be return if the value of format is empty or invalid.
// Or the value of format and nil will be returned
func checkLogTimeFormat(format string)(string, error){
	if len(format) < 1 {
		return sysadmServer.TimestampFormat["DateTime"], fmt.Errorf("format is empty,default format will be set")
	}

	for _,v := range sysadmServer.TimestampFormat {
		if strings.ToLower(format) == strings.ToLower(v) {
			return format, nil
		}
	}

	return sysadmServer.TimestampFormat["DateTime"], fmt.Errorf("format(%s) is not valid ,default format will be set.",format)
}