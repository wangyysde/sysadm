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
	"strconv"
	"strings"

	"github.com/wangyysde/sysadmServer"
	"github.com/wangyysde/yaml"
)

var ConfigDefined Config = Config{}
// Getting the path of configurationn file and try to open it
// configPath: the path of configuration file user specified
// cmdRunPath: args[0]
func getConfigPath(configPath string, cmdRunPath string) (string, error) {
	
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
func getConfigContent(configPath string) (*Config, error) {
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
func checkVerIsValid(ver string) error {
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
func checkIpAddress(address string) (net.IP, error) {
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

	fp, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR|os.O_APPEND,os.ModeAppend|os.ModePerm)
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

	for _,l := range sysadmServer.Levels {
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

// Try to get listenIP from one of  SYSADMSERVER_IP,configuration file or default value.
// The order for getting listenIP is SYSADMSERVER_IP,configuration file and default value.
func getServerAddress(confContent *Config) string{
	address := os.Getenv("SYSADMSERVER_IP")
	if address != "" {
		_, err := checkIpAddress(address)
		if err == nil {
			return address
		}
		sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_IP: %s,but the value of SYSADMSERVER_IP is not a valid server IP(%s)",address,err)
	}

	if confContent != nil {
		_,err := checkIpAddress(confContent.Server.Address)
		if err == nil {
			return confContent.Server.Address
		}
		sysadmServer.Logf("warning","We have found server address(%s) from configuration file,but the value of server address is not a valid server IP(%s).default value of server address:%s will be used.",confContent.Server.Address,err,defaultConfig.Server.Address)
	}

	return defaultConfig.Server.Address
}

// Try to get listenPort from one of  SYSADMSERVER_PORT,configuration file or default value.
// The order for getting listenPort is SYSADMSERVER_PORT,configuration file and default value.
func getServerPort(confContent *Config) int{
	port := os.Getenv("SYSADMSERVER_PORT")
	if port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_PORT: %s,but the value of SYSADMSERVER_PORT is not a valid server port(%s)",port,err)
		}else{
			_, err = checkPort(p)
			if err == nil {
				return p
			}
			sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_PORT: %d,but the value of SYSADMSERVER_PORT is not a valid server port(%s)",p,err)
		}
	}

	if confContent != nil {
		_,err := checkPort(confContent.Server.Port)
		if err == nil {
			return confContent.Server.Port
		}
		sysadmServer.Logf("warning","We have found server port(%d) from configuration file,but the value of server port is not a valid server port(%s).default value of server port:%s will be used.",confContent.Server.Port,err,defaultConfig.Server.Port)
	}

	return defaultConfig.Server.Port
}

// Try to get socket file  from one of  SYSADMSERVER_SOCK,configuration file or default value.
// The order for getting socket file is SYSADMSERVER_SOCK,configuration file and default value.
func getSockFile(confContent *Config,  cmdRunPath string) (string, error) {
	sockFile := os.Getenv("SYSADMSERVER_SOCK")
	if sockFile != "" {
		f, err := getFile(sockFile,cmdRunPath)
		if err != nil {
			sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_SOCK: %s,but the value of SYSADMSERVER_SOCK is not a valid socket file(%s)",sockFile,err)
		}else{
			return f,nil
		}
	}

	if confContent != nil {
		f,err := getFile(confContent.Server.Socket,cmdRunPath)
		if err == nil {
			return f,nil
		}
		sysadmServer.Logf("warning","We have found server socket file (%s) from configuration file,but the value of server socket file is not a valid server socket file(%s).default value of server socket file: %s will be used.",confContent.Server.Socket,err,defaultConfig.Server.Socket)
	}

	f,err := getFile(defaultConfig.Server.Socket,cmdRunPath)
	if err == nil {
		return f,nil
	}

	return "",fmt.Errorf("we can not open socket file (%s): %s .",defaultConfig.Server.Socket,err)
}

// Try to get access log file  from one of  SYSADMSERVER_ACCESSLOG,configuration file or default value.
// The order for getting access log file is SYSADMSERVER_ACCESSLOG,configuration file and default value.
func getAccessLogFile(confContent *Config,  cmdRunPath string) string {
	accessFile := os.Getenv("SYSADMSERVER_ACCESSLOG")
	if accessFile != "" {
		f, err := getFile(accessFile,cmdRunPath)
		if err != nil {
			sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_ACCESSLOG: %s,but the value of SYSADMSERVER_ACCESSLOG is not a valid access log file(%s)",accessFile,err)
		}else{
			return f
		}
	}

	if confContent != nil {
		f,err := getFile(confContent.Log.AccessLog,cmdRunPath)
		if err == nil {
			return f
		}
		sysadmServer.Logf("warning","We have found server access log file (%s) from configuration file,but the value of server access log file is not a valid server access log file(%s).default value of server access log file: %s will be used.",confContent.Log.AccessLog,err,defaultConfig.Log.AccessLog)
	}

	f,err := getFile(defaultConfig.Log.AccessLog,cmdRunPath)
	if err == nil {
		return f
	}

	return ""
}

// Try to get error log file  from one of  SYSADMSERVER_ERRORLOG,configuration file or default value.
// The order for getting error log file is SYSADMSERVER_ERRORLOG,configuration file and default value.
func getErrorLogFile(confContent *Config,  cmdRunPath string) string {
	errorFile := os.Getenv("SYSADMSERVER_ERRORLOG")
	if errorFile != "" {
		f, err := getFile(errorFile,cmdRunPath)
		if err != nil {
			sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_ERRORLOG: %s,but the value of SYSADMSERVER_ERRORLOG is not a valid error log file(%s)",errorFile,err)
		}else{
			return f
		}
	}

	if confContent != nil {
		f,err := getFile(confContent.Log.ErrorLog,cmdRunPath)
		if err == nil {
			return f
		}
		sysadmServer.Logf("warning","We have found server error log file (%s) from configuration file,but the value of server error log file is not a valid server error log file(%s).default value of server error log file: %s will be used.",confContent.Log.ErrorLog,err,defaultConfig.Log.ErrorLog)
	}

	f,err := getFile(defaultConfig.Log.ErrorLog,cmdRunPath)
	if err == nil {
		return f
	}

	return ""
}

// Try to get log kind  from one of  SYSADMSERVER_LOGKIND,configuration file or default value.
// The order for getting log kind is SYSADMSERVER_LOGKIND,configuration file and default value.
func getLogKind(confContent *Config) string{
	logKind := os.Getenv("SYSADMSERVER_LOGKIND")
	if logKind != ""{
		if strings.ToLower(logKind) == "text" || strings.ToLower(logKind) == "json" {
			return strings.ToLower(logKind)
		}
		sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_LOGKIND: %s,but the value of SYSADMSERVER_LOGKIND is not a valid kind of log(should be text or json)",logKind)
	}

	if confContent != nil {
		if strings.ToLower(confContent.Log.Kind) == "text" || strings.ToLower(confContent.Log.Kind) == "json" {
			return strings.ToLower(confContent.Log.Kind)
		}

		if len(confContent.Log.Kind) > 0 {
			sysadmServer.Logf("warning","We have found log kind (%s) from configuration file,but the value of  log kind is not valid .default value of log kind: %s will be used.",confContent.Log.Kind ,defaultConfig.Log.Kind)
		}
	}

	return defaultConfig.Log.Kind
}

// Try to get log level   from one of  SYSADMSERVER_LOGLEVEL,configuration file or default value.
// The order for getting log level is SYSADMSERVER_LOGLEVEL,configuration file and default value.
func getLogLevel(confContent *Config) string {
	logLevel := os.Getenv("SYSADMSERVER_LOGLEVEL")
	if logLevel != "" {
		level, err := checkLogLevel(logLevel)
		if err != nil {
			sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_LOGLEVEL: %s,but the value of SYSADMSERVER_LOGLEVEL is not valid(%s)",logLevel,err)
		}else{
			return level
		}
	}

	if confContent != nil {
		level,err := checkLogLevel(confContent.Log.Level)
		if err == nil {
			return level
		}
		sysadmServer.Logf("warning","We have found log level %(%s) from configuration file,but the value of log level  is not a valid(%s).default value of log level: %s will be used.",confContent.Log.Level,err,defaultConfig.Log.Level)
	}

	return defaultConfig.Log.Level
}

// Try to get log Timestampformat from one of  SYSADMSERVER_LOGTIMEFORMAT,configuration file or default value.
// The order for getting log level is SYSADMSERVER_LOGTIMEFORMAT,configuration file and default value.
func getLogTimeFormat(confContent *Config) string {
	logTimeFormat:= os.Getenv("SYSADMSERVER_LOGTIMEFORMAT")
	if logTimeFormat != "" {
		timeFormat, err := checkLogTimeFormat(logTimeFormat)
		if err != nil {
			sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_LOGTIMEFORMAT: %s,but the value of SYSADMSERVER_LOGTIMEFORMAT is not valid(%s)",logTimeFormat,err)
		}else{
			return timeFormat
		}
	}

	if confContent != nil {
		timeFormat,err := checkLogTimeFormat(confContent.Log.TimeStampFormat)
		if err == nil {
			return timeFormat
		}
		sysadmServer.Logf("warning","We have found log timestampformat from configuration file,but the value of log timestampformat is not a valid(%s).default value of log timestampformat: %s will be used.",err,defaultConfig.Log.TimeStampFormat)
	}

	return defaultConfig.Log.TimeStampFormat
}

// Try to get issplitlog for log  from one of  SYSADMSERVER_SPLITLOG,configuration file or default value.
// The order for getting log level is SYSADMSERVER_SPLITLOG,configuration file and default value.
func getIsSplitLog(confContent *Config) bool {
	isSplitLog:= os.Getenv("SYSADMSERVER_SPLITLOG")
	if isSplitLog != "" {
		if strings.ToLower(isSplitLog) == "true" || strings.ToLower(isSplitLog) == "false" {
			if strings.ToLower(isSplitLog) == "true" {
				return true
			}else {
				return false
			}
		}
		sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_SPLITLOG ,but the value of SYSADMSERVER_SPLITLOG is not bool value.")
	}

	if confContent != nil {
		return confContent.Log.SplitAccessAndError
	}

	return defaultConfig.Log.SplitAccessAndError
}

// Try to get default user from one of  SYSADMSERVER_DEFAULTUSER,configuration file or default value.
// The order for getting default user is SYSADMSERVER_DEFAULTUSER,configuration file and default value.
func getDefaultUser(confContent *Config) string {
	envUser := os.Getenv("SYSADMSERVER_DEFAULTUSER")
	if envUser != "" {
		return envUser
	}
	
	if confContent != nil {
		if confContent.User.DefaultUser != "" {
			return confContent.User.DefaultUser 
		}
	}

	return defaultConfig.User.DefaultUser
}

// Try to get default password from one of  SYSADMSERVER_DEFAULTPASSWD,configuration file or default value.
// The order for getting default user is SYSADMSERVER_DEFAULTPASSWD,configuration file and default value.
func getDefaultPassword(confContent *Config) string {
	envPasswd := os.Getenv("SYSADMSERVER_DEFAULTPASSWD")
	if envPasswd != "" {
		return envPasswd
	}
	
	if confContent != nil {
		if confContent.User.DefaultPassword != "" {
			return confContent.User.DefaultPassword 
		}
	}

	return defaultConfig.User.DefaultPassword
}

// Try to get the values of items of configuration from OS variables ,configuratio file or default value.
// The value of a item will be come from OS variables first ,then come from configuration file and last come from default value.
// All the values of items should be passed check when set it to ConfigDefined
func HandleConfig(configPath string, cmdRunPath string) (*Config,error) {
	var confContent *Config = nil
	cfgFile,err := getConfigPath(configPath,cmdRunPath)
	if err != nil {
		sysadmServer.Logf("warning","Can not get configuration file: %s",err)
		ConfigDefined.Version = confContent.Version
	}else{
		tmpConfContent,err := getConfigContent(cfgFile)
		if err != nil {
			sysadmServer.Logf("warning","Can not get the content of the configuration file: %s error: %s configuration file: %s",cfgFile, err)
			ConfigDefined.Version = confContent.Version
		}else {
			confContent = tmpConfContent
			e := checkVerIsValid(confContent.Version)
			if e != nil {
				sysadmServer.Logf("warning","%s",e)
				ConfigDefined.Version = defaultConfig.Version
			}else{
				ConfigDefined.Version = confContent.Version
			}
		}
	}

	ConfigDefined.Server.Address = getServerAddress(confContent)
	ConfigDefined.Server.Port = getServerPort(confContent) 
	ConfigDefined.Server.Socket,err = getSockFile(confContent,cmdRunPath)
	if err != nil {
		return nil,err
	}
	ConfigDefined.Log.AccessLog = getAccessLogFile(confContent,cmdRunPath)
	ConfigDefined.Log.ErrorLog = getErrorLogFile(confContent,cmdRunPath)
	ConfigDefined.Log.Kind = getLogKind(confContent)
	ConfigDefined.Log.Level = getLogLevel(confContent)
	ConfigDefined.Log.TimeStampFormat = getLogTimeFormat(confContent)
	ConfigDefined.Log.SplitAccessAndError = getIsSplitLog(confContent)
	ConfigDefined.User.DefaultUser = getDefaultUser(confContent)
	ConfigDefined.User.DefaultPassword = getDefaultPassword(confContent)

	return &ConfigDefined,nil
}