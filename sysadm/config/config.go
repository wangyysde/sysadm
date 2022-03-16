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
	sysadmDB "github.com/wangyysde/sysadm/db"
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
func getFile(f string,cmdRunPath string, isRmTest bool)(string,error){
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
	if isRmTest {
		_ = os.Remove(f)
	}
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
		f, err := getFile(sockFile,cmdRunPath,true)
		if err != nil {
			sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_SOCK: %s,but the value of SYSADMSERVER_SOCK is not a valid socket file(%s)",sockFile,err)
		}else{
			return f,nil
		}
	}

	if confContent != nil {
		f,err := getFile(confContent.Server.Socket,cmdRunPath,true)
		if err == nil {
			return f,nil
		}
		sysadmServer.Logf("warning","We have found server socket file (%s) from configuration file,but the value of server socket file is not a valid server socket file(%s).default value of server socket file: %s will be used.",confContent.Server.Socket,err,defaultConfig.Server.Socket)
	}

	f,err := getFile(defaultConfig.Server.Socket,cmdRunPath,true)
	if err == nil {
		return f,nil
	}

	return "",fmt.Errorf("we can not open socket file (%s): %s .",defaultConfig.Server.Socket,err)
}

// Getting the tls value for sysadm server from environment and checking the validity of it
// return the tls if it is valid ,otherwise getting the value from 
// configuration file and checking the validity of it. return the it if it is valid.
// otherwise return the default value.
func getSysadmTls(confContent *Config) bool{
	tls := os.Getenv("SYSADMSERVER_TLS")
	if strings.ToLower(tls) == "true" || strings.ToLower(tls) == "false" {
		if strings.ToLower(tls) == "true" {
			return true
		}else {
			return false
		}
	}
	
	if confContent != nil  {
		return confContent.Server.Tls
	}

	return defaultConfig.Server.Tls
}

/* 
  Getting ca of sysadm server  from environment and checking the validity of it
  return the ca if it is valid ,otherwise getting ca of sysadm  from 
  configuration file and checking the validity of it. return the user if it is valid.
  otherwise return the default ca of sysadm server.
*/
  func getSysadmCa(confContent *Config) string{
	ca := os.Getenv("SYSADMSERVER_CA")
	if ca != ""{
		return ca
	}

	if confContent != nil  {
		if confContent.Server.Ca != "" {
			return confContent.Server.Ca
		}
	}

	return defaultConfig.Server.Ca
}

/* 
  Getting cert of sysadm server  from environment and checking the validity of it
  return the cert if it is valid ,otherwise getting cert of sysadm  from 
  configuration file and checking the validity of it. return the user if it is valid.
  otherwise return the default ca of sysadm server.
*/
  func getSysadmCert(confContent *Config) string{
	cert := os.Getenv("SYSADMSERVER_CERT")
	if cert != ""{
		return cert
	}

	if confContent != nil  {
		if confContent.Server.Cert != "" {
			return confContent.Server.Cert
		}
	}

	return defaultConfig.Server.Cert
}

/* 
  Getting key of sysadm server  from environment and checking the validity of it
  return the cert if it is valid ,otherwise getting cert of sysadm  from 
  configuration file and checking the validity of it. return the user if it is valid.
  otherwise return the default ca of sysadm server.
*/
  func getSysadmKey(confContent *Config) string{
	key := os.Getenv("SYSADMSERVER_KEY")
	if key != ""{
		return key
	}

	if confContent != nil  {
		if confContent.Server.Key != "" {
			return confContent.Server.Key
		}
	}

	return defaultConfig.Server.Key
}


// Try to get access log file  from one of  SYSADMSERVER_ACCESSLOG,configuration file or default value.
// The order for getting access log file is SYSADMSERVER_ACCESSLOG,configuration file and default value.
func getAccessLogFile(confContent *Config,  cmdRunPath string) string {
	accessFile := os.Getenv("SYSADMSERVER_ACCESSLOG")
	if accessFile != "" {
		f, err := getFile(accessFile,cmdRunPath,false)
		if err != nil {
			sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_ACCESSLOG: %s,but the value of SYSADMSERVER_ACCESSLOG is not a valid access log file(%s)",accessFile,err)
		}else{
			return f
		}
	}

	if confContent != nil {
		f,err := getFile(confContent.Log.AccessLog,cmdRunPath,false)
		if err == nil {
			return f
		}
		sysadmServer.Logf("warning","We have found server access log file (%s) from configuration file,but the value of server access log file is not a valid server access log file(%s).default value of server access log file: %s will be used.",confContent.Log.AccessLog,err,defaultConfig.Log.AccessLog)
	}

	f,err := getFile(defaultConfig.Log.AccessLog,cmdRunPath,false)
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
		f, err := getFile(errorFile,cmdRunPath,false)
		if err != nil {
			sysadmServer.Logf("warning","We have found environment variable SYSADMSERVER_ERRORLOG: %s,but the value of SYSADMSERVER_ERRORLOG is not a valid error log file(%s)",errorFile,err)
		}else{
			return f
		}
	}

	if confContent != nil {
		f,err := getFile(confContent.Log.ErrorLog,cmdRunPath,false)
		if err == nil {
			return f
		}
		sysadmServer.Logf("warning","We have found server error log file (%s) from configuration file,but the value of server error log file is not a valid server error log file(%s).default value of server error log file: %s will be used.",confContent.Log.ErrorLog,err,defaultConfig.Log.ErrorLog)
	}

	f,err := getFile(defaultConfig.Log.ErrorLog,cmdRunPath,false)
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

// check the validity of IP address 
// return IP(string) if the ip address is valid
// Or return nil with error
func checkHostAddress(address string) (string, error) {
	if len(address) < 1 {
		return "", fmt.Errorf("The address(%s) is empty or the length of it is less 1",address)
	}

	if ip := net.ParseIP(address); ip != nil {
		return address,nil
	}

	ips,err := net.LookupIP(address)
	if err != nil {
		return "" , fmt.Errorf("Lookup the IP of address(%s) error.",err)
	}
	
	return ips[0].String(), nil
}

// Getting type  of DB from environment and checking the validity of it
// returns the type value  if it is valid ,otherwise getting type  of DB  from 
// configuration file and checking the validity of it. return the type value if it is valid.
// otherwise return the default DB type(postgre).
func getDBType(confContent *Config)string{
	dbType := os.Getenv("SYSADMSERVER_DBTYPE")
	if dbType != ""{
		if sysadmDB.IsSupportedDB(dbType) {
			return strings.ToLower(dbType)
		}
	}

	if confContent != nil  {
		dbType := confContent.DB.Type
		if sysadmDB.IsSupportedDB(dbType) {
			return strings.ToLower(dbType)
		}
	}

	return defaultConfig.DB.Type
}

// Getting host address of Postgre  from environment and checking the validity of it
// return the address of it is valid ,otherwise getting host address of Postgre  from 
// configuration file and checking the validity of it. return the address of it is valid.
// otherwise return the default address of Postgre.
func getDBHost(confContent *Config)string{
	dbHost := os.Getenv("SYSADMSERVER_DBHOST")
	if dbHost != ""{
		if host,err := checkHostAddress(dbHost); err == nil{
			return host
		}
	}

	if confContent != nil  {
		if host,err := checkHostAddress(confContent.DB.Host); err == nil{
			return host
		}
	}

	return defaultConfig.DB.Host
}

// Getting port of Postgre  from environment and checking the validity of it
// return the port if it is valid ,otherwise getting port of Postgre  from 
// configuration file and checking the validity of it. return the port if it is valid.
// otherwise return the default port of Postgre.
func getDBPort(confContent *Config) int{
	dbPort := os.Getenv("SYSADMSERVER_DBPORT")
	if dbPort != ""{
		port, err := strconv.Atoi(dbPort)
		if err == nil {
			if port > 1024 && port < 65536 {
				return port
			}
		}
	}

	if confContent != nil  {
		if confContent.DB.Port > 1024 && confContent.DB.Port <= 65535 {
			return confContent.DB.Port 
		}
	}

	return defaultConfig.DB.Port
}

// Getting user of Postgre  from environment and checking the validity of it
// return the user if it is valid ,otherwise getting user of Postgre  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default user of Postgre.
func getDBUser(confContent *Config) string{
	dbUser := os.Getenv("SYSADMSERVER_DBUSER")
	if dbUser != ""{
		return dbUser
	}

	if confContent != nil  {
		if confContent.DB.User != "" {
			return confContent.DB.User
		}
	}

	return defaultConfig.DB.User
}

// Getting Password of Postgre  from environment and checking the validity of it
// return the Password if it is valid ,otherwise getting Password of Postgre  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default user of Postgre.
func getDBPassword(confContent *Config) string{
	dbPassword := os.Getenv("SYSADMSERVER_DBPASSWORD")
	if dbPassword != ""{
		return dbPassword
	}

	if confContent != nil  {
		if confContent.DB.Password != "" {
			return confContent.DB.Password
		}
	}

	return defaultConfig.DB.Password
}

// Getting DBName of Postgre  from environment and checking the validity of it
// return the DBName if it is valid ,otherwise getting DBName of Postgre  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default DBName of Postgre.
func getDBDBName(confContent *Config) string{
	dbDBName := os.Getenv("SYSADMSERVER_DBDBNAME")
	if dbDBName != ""{
		return dbDBName
	}

	if confContent != nil  {
		if confContent.DB.Dbname != "" {
			return confContent.DB.Dbname
		}
	}

	return defaultConfig.DB.Dbname
}

// Getting MaxConnect of Postgre  from environment and checking the validity of it
// return the MaxConnect if it is valid ,otherwise getting MaxConnect of Postgre  from 
// configuration file and checking the validity of it. return the port if it is valid.
// otherwise return the default MaxConnect of Postgre.
func getDBMaxConnect(confContent *Config) int{
	dbMaxConnect := os.Getenv("SYSADMSERVER_DBMAXCONNECT")
	if dbMaxConnect != ""{
		maxConnect,err := strconv.Atoi(dbMaxConnect)
		if err == nil{
			if maxConnect >1 && maxConnect < 20000 {
				return maxConnect
			}
		}
	}

	if confContent != nil  {
		if confContent.DB.DbMaxConnect > 1 && confContent.DB.DbMaxConnect <= 20000 {
			return confContent.DB.DbMaxConnect 
		}
	}

	return defaultConfig.DB.DbMaxConnect
}

// Getting dbIdleConnect of Postgre  from environment and checking the validity of it
// return the dbIdleConnect if it is valid ,otherwise getting dbIdleConnect of Postgre  from 
// configuration file and checking the validity of it. return the port if it is valid.
// otherwise return the default dbIdleConnect of Postgre.
func getDBDbIdleConnect(confContent *Config) int{
	dbIdleConnect := os.Getenv("SYSADMSERVER_DBIDLECONNECT")
	if dbIdleConnect != ""{
		idleConnect, err := strconv.Atoi(dbIdleConnect)
		if err == nil {
			if idleConnect >1 && idleConnect < 20000 {
				return idleConnect
			}
		}
	}

	if confContent != nil  {
		if confContent.DB.DbIdleConnect > 1 && confContent.DB.DbIdleConnect <= 20000 {
			return confContent.DB.DbIdleConnect 
		}
	}

	return defaultConfig.DB.DbIdleConnect
}

// Getting Sslmode of Postgre  from environment and checking the validity of it
// return the Sslmode if it is valid ,otherwise getting Sslmode of Postgre  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default Sslmode of Postgre.
func getDBSslmode(confContent *Config) string{
	dbSslmode := os.Getenv("SYSADMSERVER_DBSSLMODE")
	if dbSslmode != ""{
		return dbSslmode
	}

	if confContent != nil  {
		if confContent.DB.Sslmode != "" {
			return confContent.DB.Sslmode
		}
	}

	return defaultConfig.DB.Sslmode
}

// Checking a file if is exists.
func checkFileExists(f string,cmdRunPath string ) bool {
	dir ,error := filepath.Abs(filepath.Dir(cmdRunPath))
	if error != nil {
		sysadmServer.Logf("error","Getting program root path error: %s",error)
		return false
	}

	if ! filepath.IsAbs(f) {
		tmpDir := filepath.Join(dir,"../")
		f = filepath.Join(tmpDir,f)
	}

	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// Getting Sslrootcert of Postgre  from environment and checking the validity of it
// return the Sslrootcert if it is valid ,otherwise getting Sslrootcert of Postgre  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default Sslrootcert of Postgre.
func getPostgreSslrootcert(confContent *Config) string{
	dbSslCa := os.Getenv("SYSADMSERVER_DBSSLCA")
	if dbSslCa != ""{
		return dbSslCa
	}

	if confContent != nil  {
		if confContent.DB.Sslrootcert != "" {
			return confContent.DB.Sslrootcert
		}
	}

	return defaultConfig.DB.Sslrootcert
}

// Getting Sslkey of Postgre  from environment and checking the validity of it
// return the Sslkey if it is valid ,otherwise getting Sslkey of Postgre  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default Sslkey of Postgre.
func getPostgreSslkey(confContent *Config) string{
	dbSslKey := os.Getenv("SYSADMSERVER_DBSSLKEY")
	if dbSslKey != ""{
		return dbSslKey
	}

	if confContent != nil  {
		if confContent.DB.Sslkey != "" {
			return confContent.DB.Sslkey
		}
	}

	return defaultConfig.DB.Sslkey
}

// Getting Sslcert of Postgre  from environment and checking the validity of it
// return the Sslcert if it is valid ,otherwise getting Sslcert of Postgre  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default Sslcert of Postgre.
func getPostgreSslcert(confContent *Config) string{
	dbSslcert := os.Getenv("SYSADMSERVER_DBSSLCERT")
	if dbSslcert != ""{
		return dbSslcert
	}

	if confContent != nil  {
		if confContent.DB.Sslcert != "" {
			return confContent.DB.Sslcert
		}
	}

	return defaultConfig.DB.Sslcert
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
	ConfigDefined.Server.Tls = getSysadmTls(confContent)
	if ConfigDefined.Server.Tls {
		caFile := getSysadmCa(confContent)
		if caFile == "" || !checkFileExists(ConfigDefined.Server.Ca, cmdRunPath) {
			sysadmServer.Logf("warning","Tls of Sysadm server  has be set to true But CA(%s) can not be found. We will try to set Tls to false",confContent.Server.Ca)
			 ConfigDefined.Server.Tls =false
		}
	}

	if ConfigDefined.Server.Tls {
		certFile := getSysadmCert(confContent)
		if certFile == "" || !checkFileExists(ConfigDefined.Server.Cert, cmdRunPath) {
			sysadmServer.Logf("warning","Tls of Sysadm server  has be set to true But Cert(%s) can not be found. We will try to set Tls to false",confContent.Server.Cert)
			 ConfigDefined.Server.Tls =false
		}
	}

	if ConfigDefined.Server.Tls {
		keyFile := getSysadmCert(confContent)
		if keyFile == "" || !checkFileExists(ConfigDefined.Server.Key, cmdRunPath) {
			sysadmServer.Logf("warning","Tls of Sysadm server  has be set to true But Key(%s) can not be found. We will try to set Tls to false",confContent.Server.Key)
			 ConfigDefined.Server.Tls =false
		}
	}

	ConfigDefined.Log.AccessLog = getAccessLogFile(confContent,cmdRunPath)
	ConfigDefined.Log.ErrorLog = getErrorLogFile(confContent,cmdRunPath)
	ConfigDefined.Log.Kind = getLogKind(confContent)
	ConfigDefined.Log.Level = getLogLevel(confContent)
	ConfigDefined.Log.TimeStampFormat = getLogTimeFormat(confContent)
	ConfigDefined.Log.SplitAccessAndError = getIsSplitLog(confContent)
	ConfigDefined.User.DefaultUser = getDefaultUser(confContent)
	ConfigDefined.User.DefaultPassword = getDefaultPassword(confContent)
	ConfigDefined.DB.Type = getDBType(confContent)
	ConfigDefined.DB.Host = getDBHost(confContent)
	ConfigDefined.DB.Port = getDBPort(confContent)
	ConfigDefined.DB.User = getDBUser(confContent)
	ConfigDefined.DB.Password = getDBPassword(confContent)
	ConfigDefined.DB.Dbname = getDBDBName(confContent)
	ConfigDefined.DB.Sslmode = getDBSslmode(confContent)
	if strings.ToLower(ConfigDefined.DB.Sslmode) != "disable" {
		ConfigDefined.DB.Sslrootcert = getPostgreSslrootcert(confContent)
		if !checkFileExists(ConfigDefined.DB.Sslrootcert, cmdRunPath) {
			sysadmServer.Logf("warning","SslMode of Postgre has be set to %s But sslCA(%s) can not be found. We will try to set SslMode to disable ",ConfigDefined.DB.Sslmode, ConfigDefined.DB.Sslrootcert)
			ConfigDefined.DB.Sslmode = "disable"
		}
	}

	if strings.ToLower(ConfigDefined.DB.Sslmode) != "disable" {
		ConfigDefined.DB.Sslkey = getPostgreSslkey(confContent)
		if !checkFileExists(ConfigDefined.DB.Sslkey, cmdRunPath) {
			sysadmServer.Logf("warning","SslMode of Postgre has be set to %s But SslKey(%s) can not be found. We will try to set SslMode to disable ",ConfigDefined.DB.Sslmode, ConfigDefined.DB.Sslkey)
			ConfigDefined.DB.Sslmode = "disable"
		}
	}

	if strings.ToLower(ConfigDefined.DB.Sslmode) != "disable" {
		ConfigDefined.DB.Sslcert = getPostgreSslcert(confContent)
		if !checkFileExists(ConfigDefined.DB.Sslcert, cmdRunPath) {
			sysadmServer.Logf("warning","SslMode of Postgre has be set to %s But Sslcert(%s) can not be found. We will try to set SslMode to disable ",ConfigDefined.DB.Sslmode, ConfigDefined.DB.Sslcert)
			ConfigDefined.DB.Sslmode = "disable"
		}
	}

	ConfigDefined.DB.DbMaxConnect = getDBMaxConnect(confContent)
	ConfigDefined.DB.DbIdleConnect = getDBDbIdleConnect(confContent)

	return &ConfigDefined,nil
}

