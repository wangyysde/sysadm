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

	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
	"github.com/wangyysde/yaml"
)

var ConfigDefined Config = Config{}
// Getting the path of configurationn file and try to open it
// configPath: the path of configuration file user specified
// cmdRunPath: args[0]
func getConfigPath(configPath string, cmdRunPath string) (string, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	dir ,error := filepath.Abs(filepath.Dir(cmdRunPath))
	if error != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201000,"error","Can not get absolute path for promgram. error: %s.",error))
		return "",errs
	}

	if configPath == "" {
		configPath = filepath.Join(dir,".../")
		configPath = filepath.Join(configPath, DefaultConfigFile)
		fp, err := os.Open(configPath)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201001,"error","Can not open the configuration file:%s, error: %s.",configPath,err))
			return "",errs
		}

		fp.Close()
		return configPath,errs
	}

	if ! filepath.IsAbs(configPath) {
		tmpDir := filepath.Join(dir,"../")
		configPath = filepath.Join(tmpDir,configPath)
	}

	fp, err := os.Open(configPath)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201002,"error","Can not open the configuration file:%s, error: %s.",configPath,err))
		return "",errs
	}
	fp.Close()

	return configPath,errs
}

// Reading the content of configuration from configPath and parsing the content 
// returning a pointer to Config if it is successfully parsed
// Or returning an error and nil
func getConfigContent(configPath string) (*Config, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	if configPath == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201003,"error","The configration file path must not empty."))
		return nil, errs
	}

	yamlContent, err := ioutil.ReadFile(configPath) 
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201004,"error","Can not read configuration file: %s error: %s.",configPath,err))
		return nil, errs
	}

	err = yaml.Unmarshal(yamlContent, &ConfigDefined)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201005,"error","Can not Unmarshal configuration contenet error: %s.",err))
		return nil, errs
	}

	return &ConfigDefined, errs
}

// check the supporting of the version of configuration file specified
// return nil if the version be supported
// Or return err 
func checkVerIsValid(ver string) ([]sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	found := false
	for _,value := range SupportVers {
		if strings.EqualFold(ver,value)  {
			found = true
			break
		}
	}

	if found {
		return errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(201006,"error","The version(%s) of the configuration file specified was not be supported by this release.",ver))

	return errs
}

// check the validity of IP address 
// return IP(net.IP) if the ip address is valid
// Or return nil with error
func checkIpAddress(address string) (net.IP, error) {
	if len(address) < 1 {
		return nil, fmt.Errorf("the address(%s) is empty or the length of it is less 1",address)
	}

	if ip := net.ParseIP(address); ip != nil {
		if address == "0.0.0.0" || address == "::" {
			return ip, nil
		}

		adds,err := net.InterfaceAddrs()
		if err != nil {
			return nil, fmt.Errorf("get interface address error: %s",err)
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

		return nil, fmt.Errorf("the address(%s) is not any of the addresses of host interfaces",address)
	}

	ips,err := net.LookupIP(address)
	if err != nil {
		return nil , fmt.Errorf("lookup the ip of address(%s) error",err)
	}

	adds,err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("get interface address error: %s",err)
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
	
	return nil, fmt.Errorf("the ip(%v) to the address(hostname:%s) is not any the IP address of host interfaces",ips,address)
}

// check the validity of port 
// return port with nil if the port is valid 
// Or return 0 with error
func checkPort(port int)(int, error){
	if port > 1024 && port <= 65535 {
		return port,nil
	}

	return 0, fmt.Errorf("the port should be great than 1024 and less than 65535. Now is :%d",port)
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
		return defaultConfig.Log.Level, fmt.Errorf("level is empty,default level has be set")
	}

	for _,l := range sysadmServer.Levels {
		if  strings.EqualFold(level,l)  {
			return strings.ToLower(level),nil
		}
	}

	return defaultConfig.Log.Level,fmt.Errorf("level(%s) was not found,default level has be set",level)
}

// check the validity of log format.
// The default format will be return if the value of format is empty or invalid.
// Or the value of format and nil will be returned
func checkLogTimeFormat(format string)(string, error){
	if len(format) < 1 {
		return sysadmServer.TimestampFormat["DateTime"], fmt.Errorf("format is empty,default format will be set")
	}

	for _,v := range sysadmServer.TimestampFormat {
		if strings.EqualFold(format,v) {
			return format, nil
		}
	}

	return sysadmServer.TimestampFormat["DateTime"], fmt.Errorf("format(%s) is not valid ,default format will be set",format)
}

// Try to get listenIP from one of  SERVER_IP,configuration file or default value.
// The order for getting listenIP is SERVER_IP,configuration file and default value.
func getServerAddress(confContent *Config)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	address := os.Getenv("SERVER_IP")
	if address != "" {
		_, err := checkIpAddress(address)
		if err == nil {
			return address,errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201007,"warning","We have found environment variable SERVER_IP: %s,but the value of SERVER_IP is not a valid server IP(%s)",address,err))
	}

	if confContent != nil {
		_,err := checkIpAddress(confContent.Server.Address)
		if err == nil {
			return confContent.Server.Address,errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201008,"warning","We have found server address(%s) from configuration file,but the value of server address is not a valid server IP(%s).default value of server address:%s will be used.",confContent.Server.Address,err,defaultConfig.Server.Address))
	}

	return defaultConfig.Server.Address,errs
}

// Try to get listenPort from one of  SERVER_PORT,configuration file or default value.
// The order for getting listenPort is SERVER_PORT,configuration file and default value.
func getServerPort(confContent *Config) (int,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	port := os.Getenv("SERVER_PORT")
	if port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201009,"warning","We have found environment variable SERVER_PORT: %s,but the value of SERVER_PORT is not a valid server port(%s)",port,err))
		}else{
			_, err = checkPort(p)
			if err == nil {
				return p,errs
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201010,"warning","We have found environment variable SERVER_PORT: %d,but the value of SERVER_PORT is not a valid server port(%s)",p,err))
		}
	}

	if confContent != nil {
		_,err := checkPort(confContent.Server.Port)
		if err == nil {
			return confContent.Server.Port,errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201011,"warning","We have found server port(%d) from configuration file,but the value of server port is not a valid server port(%s).default value of server port:%s will be used.",confContent.Server.Port,err,defaultConfig.Server.Port))
	}

	return defaultConfig.Server.Port,errs
}

// Try to get socket file  from one of  SOCKET,configuration file or default value.
// The order for getting socket file is SOCKET,configuration file and default value.
func getSockFile(confContent *Config,  cmdRunPath string) (string, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	sockFile := os.Getenv("SOCKET")
	if sockFile != "" {
		f, err := getFile(sockFile,cmdRunPath,true)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201012,"warning","We have found environment variable SOCKET: %s,but the value of SOCKET is not a valid socket file(%s)",sockFile,err))
		}else{
			return f,errs
		}
	}

	if confContent != nil {
		f,err := getFile(confContent.Server.Socket,cmdRunPath,true)
		if err == nil {
			return f,errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201013,"warning","We have found server socket file (%s) from configuration file,but the value of server socket file is not a valid server socket file(%s).default value of server socket file: %s will be used.",confContent.Server.Socket,err,defaultConfig.Server.Socket))

	}

	f,err := getFile(defaultConfig.Server.Socket,cmdRunPath,true)
	if err == nil {
		return f,errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(201014,"error","we can not open socket file (%s): %s .",defaultConfig.Server.Socket,err))

	return "",errs
}

// Try to get access log file  from one of  ACCESSLOG,configuration file or default value.
// The order for getting access log file is ACCESSLOG,configuration file and default value.
func getAccessLogFile(confContent *Config,  cmdRunPath string) (string,[]sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	accessFile := os.Getenv("ACCESSLOG")
	if accessFile != "" {
		f, err := getFile(accessFile,cmdRunPath,false)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201015,"warning","We have found environment variable ACCESSLOG: %s,but the value of ACCESSLOG is not a valid access log file(%s)",accessFile,err))
		}else{
			return f,errs
		}
	}

	if confContent != nil {
		f,err := getFile(confContent.Log.AccessLog,cmdRunPath,false)
		if err == nil {
			return f,errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201016,"warning","We have found server access log file (%s) from configuration file,but the value of server access log file is not a valid server access log file(%s).default value of server access log file: %s will be used.",confContent.Log.AccessLog,err,defaultConfig.Log.AccessLog))
	}

	f,err := getFile(defaultConfig.Log.AccessLog,cmdRunPath,false)
	if err == nil {
		return f,errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(201017,"error","We can not using default access log file %s, error: %s.",defaultConfig.Log.AccessLog,err))
	return "",errs
}

// Try to get error log file  from one of  ERRORLOG,configuration file or default value.
// The order for getting error log file is ERRORLOG,configuration file and default value.
func getErrorLogFile(confContent *Config,  cmdRunPath string) (string, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	errorFile := os.Getenv("ERRORLOG")
	if errorFile != "" {
		f, err := getFile(errorFile,cmdRunPath,false)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201018,"warning","We have found environment variable ERRORLOG: %s,but the value of ERRORLOG is not a valid error log file(%s)",errorFile,err))
		}else{
			return f,errs
		}
	}

	if confContent != nil {
		f,err := getFile(confContent.Log.ErrorLog,cmdRunPath,false)
		if err == nil {
			return f,errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201019,"warning","We have found server error log file (%s) from configuration file,but the value of server error log file is not a valid server error log file(%s).default value of server error log file: %s will be used.",confContent.Log.ErrorLog,err,defaultConfig.Log.ErrorLog))

	}

	f,err := getFile(defaultConfig.Log.ErrorLog,cmdRunPath,false)
	if err == nil {
		return f,errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(201020,"error","We can not using default error log file %s, error: %s.",defaultConfig.Log.ErrorLog,err))
	return "",errs
}

// Try to get log kind  from one of  LOGKIND,configuration file or default value.
// The order for getting log kind is LOGKIND,configuration file and default value.
func getLogKind(confContent *Config) (string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	logKind := os.Getenv("LOGKIND")
	if logKind != ""{
		if strings.ToLower(logKind) == "text" || strings.ToLower(logKind) == "json" {
			return strings.ToLower(logKind),errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201021,"warning","We have found environment variable LOGKIND: %s,but the value of LOGKIND is not a valid kind of log(should be text or json)",logKind))
	}

	if confContent != nil {
		if strings.ToLower(confContent.Log.Kind) == "text" || strings.ToLower(confContent.Log.Kind) == "json" {
			return strings.ToLower(confContent.Log.Kind),errs
		}

		if len(confContent.Log.Kind) > 0 {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201022,"warning","We have found log kind (%s) from configuration file,but the value of  log kind is not valid .default value of log kind: %s will be used.",confContent.Log.Kind ,defaultConfig.Log.Kind))
		}
	}

	return defaultConfig.Log.Kind,errs
}

// Try to get log level   from one of  LOGLEVEL,configuration file or default value.
// The order for getting log level is LOGLEVEL,configuration file and default value.
func getLogLevel(confContent *Config) (string, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	logLevel := os.Getenv("LOGLEVEL")
	if logLevel != "" {
		level, err := checkLogLevel(logLevel)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201023,"warning","we have found environment variable LOGLEVEL: %s,but the value of LOGLEVEL is not valid(%s)",logLevel,err))
		}else{
			return level,errs
		}
	}

	if confContent != nil {
		level,err := checkLogLevel(confContent.Log.Level)
		if err == nil {
			return level,errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201024,"warning","we have found log level %(%s) from configuration file,but the value of log level  is not a valid(%s).default value of log level: %s will be used.",confContent.Log.Level,err,defaultConfig.Log.Level))
	}

	return defaultConfig.Log.Level,errs
}

// Try to get log Timestampformat from one of  LOGTIMEFORMAT,configuration file or default value.
// The order for getting log level is LOGTIMEFORMAT,configuration file and default value.
func getLogTimeFormat(confContent *Config) (string, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	logTimeFormat:= os.Getenv("LOGTIMEFORMAT")
	if logTimeFormat != "" {
		timeFormat, err := checkLogTimeFormat(logTimeFormat)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201025,"warning","We have found environment variable LOGTIMEFORMAT: %s,but the value of LOGTIMEFORMAT is not valid(%s)",logTimeFormat,err))
		}else{
			return timeFormat,errs
		}
	}

	if confContent != nil {
		timeFormat,err := checkLogTimeFormat(confContent.Log.TimeStampFormat)
		if err == nil {
			return timeFormat,errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201026,"warning","We have found log timestampformat from configuration file,but the value of log timestampformat is not a valid(%s).default value of log timestampformat: %s will be used.",err,defaultConfig.Log.TimeStampFormat))
	}

	return defaultConfig.Log.TimeStampFormat,errs
}

// Try to get issplitlog for log  from one of  SPLITLOG,configuration file or default value.
// The order for getting log level is SPLITLOG,configuration file and default value.
func getIsSplitLog(confContent *Config) (bool,[]sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	isSplitLog:= os.Getenv("SPLITLOG")
	if isSplitLog != "" {
		if strings.ToLower(isSplitLog) == "true" || strings.ToLower(isSplitLog) == "false" {
			if strings.ToLower(isSplitLog) == "true" {
				return true,errs
			}else {
				return false,errs
			}
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201027,"warning","we have found environment variable SPLITLOG ,but the value of SPLITLOG is not bool value."))
	}

	if confContent != nil {
		return confContent.Log.SplitAccessAndError,errs
	}

	return defaultConfig.Log.SplitAccessAndError,errs
}

// Try to get default user from one of  DEFAULTUSER,configuration file or default value.
// The order for getting default user is DEFAULTUSER,configuration file and default value.
func getDefaultUser(confContent *Config) string {
	envUser := os.Getenv("DEFAULTUSER")
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

// Try to get default password from one of  DEFAULTPASSWD,configuration file or default value.
// The order for getting default user is DEFAULTPASSWD,configuration file and default value.
func getDefaultPassword(confContent *Config) string {
	envPasswd := os.Getenv("DEFAULTPASSWD")
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
		return "", fmt.Errorf("the address(%s) is empty or the length of it is less 1",address)
	}

	if ip := net.ParseIP(address); ip != nil {
		return address,nil
	}

	ips,err := net.LookupIP(address)
	if err != nil {
		return "" , fmt.Errorf("lookup the IP of address(%s) error",err)
	}
	
	return ips[0].String(), nil
}

// Getting host address of Postgre  from environment and checking the validity of it
// return the address of it is valid ,otherwise getting host address of Postgre  from 
// configuration file and checking the validity of it. return the address of it is valid.
// otherwise return the default address of Postgre.
func getPostgreHost(confContent *Config)string{
	dbHost := os.Getenv("DBHOST")
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
func getPostgrePort(confContent *Config) int{
	dbPort := os.Getenv("DBPORT")
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
func getPostgreUser(confContent *Config) string{
	dbUser := os.Getenv("DBUSER")
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
func getPostgrePassword(confContent *Config) string{
	dbPassword := os.Getenv("DBPASSWORD")
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
func getPostgreDBName(confContent *Config) string{
	dbDBName := os.Getenv("DBDBNAME")
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
func getPostgreMaxConnect(confContent *Config) int{
	dbMaxConnect := os.Getenv("DBMAXCONNECT")
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
func getPostgreDbIdleConnect(confContent *Config) int{
	dbIdleConnect := os.Getenv("DBIDLECONNECT")
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
func getPostgreSslmode(confContent *Config) string{
	dbSslmode := os.Getenv("DBSSLMODE")
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
func checkFileExists(f string,cmdRunPath string ) (bool,error) {
	dir ,error := filepath.Abs(filepath.Dir(cmdRunPath))
	if error != nil {
		return false,error
	}

	if ! filepath.IsAbs(f) {
		tmpDir := filepath.Join(dir,"../")
		f = filepath.Join(tmpDir,f)
	}

	_, err := os.Stat(f)
	
	return !os.IsNotExist(err),err
}

// Getting Sslrootcert of Postgre  from environment and checking the validity of it
// return the Sslrootcert if it is valid ,otherwise getting Sslrootcert of Postgre  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default Sslrootcert of Postgre.
func getPostgreSslrootcert(confContent *Config) string{
	dbSslCa := os.Getenv("DBSSLCA")
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
	dbSslKey := os.Getenv("DBSSLKEY")
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
	dbSslcert := os.Getenv("DBSSLCERT")
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

// Getting host address of Registry from environment and checking the validity of it
// return the address of it is valid ,otherwise getting host address of Registry  from 
// configuration file and checking the validity of it. return the address of it is valid.
// otherwise return the default address of Registry.
func getRegistryHost(confContent *Config)(string,error){
	registryHost := os.Getenv("REGISTRYHOST")
	if registryHost != ""{
		if host,err := checkHostAddress(registryHost); err == nil{
			return host,nil
		}
	}

	if confContent != nil  {
		if host,err := checkHostAddress(confContent.Registry.Server.Host); err == nil{
			return host,nil
		}
	}

	host, err := checkHostAddress(defaultConfig.Registry.Server.Host)
	if err == nil{
		return host,nil
	}

	return defaultConfig.Registry.Server.Host, err
}

// Getting port of registry  from environment and checking the validity of it
// return the port if it is valid ,otherwise getting port of registry  from 
// configuration file and checking the validity of it. return the port if it is valid.
// otherwise return the default port of Postgre.
func getRegistryPort(confContent *Config) int{
	registryPort := os.Getenv("REGISTRYPORT")
	if registryPort != ""{
		port, err := strconv.Atoi(registryPort)
		if err == nil {
			if port > 1024 && port < 65536 {
				return port
			}
		}
	}

	if confContent != nil  {
		if confContent.Registry.Server.Port > 1024 && confContent.Registry.Server.Port <= 65535 {
			return confContent.Registry.Server.Port 
		}
	}

	return defaultConfig.Registry.Server.Port
}

// Getting Sslkey of Registry  from environment and checking the validity of it
// return the Sslkey if it is valid ,otherwise getting Sslkey of Registry  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default Sslkey of Postgre.
func getRegistrySslkey(confContent *Config) string{
	registrySslKey := os.Getenv("REGISTRYSSLKEY")
	if registrySslKey != ""{
		return registrySslKey
	}

	if confContent != nil  {
		if confContent.Registry.Server.Certs.Key != "" {
			return confContent.Registry.Server.Certs.Key
		}
	}

	return defaultConfig.Registry.Server.Certs.Key
}

// Getting Sslcert of Registry  from environment and checking the validity of it
// return the Sslcert if it is valid ,otherwise getting Sslkey of Registry  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default Sslcert of Postgre.
func getRegistrySslCert(confContent *Config) string{
	registrySslCert := os.Getenv("REGISTRYSSLCERT")
	if registrySslCert != ""{
		return registrySslCert
	}

	if confContent != nil  {
		if confContent.Registry.Server.Certs.Cert != "" {
			return confContent.Registry.Server.Certs.Cert
		}
	}

	return defaultConfig.Registry.Server.Certs.Cert
}

// Getting Sslca of Registry  from environment and checking the validity of it
// return the Sslrootcert if it is valid ,otherwise getting Sslca of Registry  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default Sslca of Postgre.
func getRegistrySslCa(confContent *Config) string{
	registrySslCa := os.Getenv("REGISTRYSSLCA")
	if registrySslCa != ""{
		return registrySslCa
	}

	if confContent != nil  {
		if confContent.Registry.Server.Certs.Ca != "" {
			return confContent.Registry.Server.Certs.Ca
		}
	}

	return defaultConfig.Registry.Server.Certs.Ca
}

// Getting user of Registry  from environment and checking the validity of it
// return the user if it is valid ,otherwise getting user of Postgre  from 
// configuration file and checking the validity of it. return the user if it is valid.
// otherwise return the default user of Postgre.
func getRegistryUser(confContent *Config) string{
	registryUser := os.Getenv("REGISTRYUSER")
	if registryUser != ""{
		return registryUser
	}

	if confContent != nil  {
		if confContent.Registry.Credit.Username != "" {
			return confContent.Registry.Credit.Username
		}
	}

	return defaultConfig.Registry.Credit.Username
}

// Getting password of Registry  from environment and checking the validity of it
// return the password if it is valid ,otherwise getting password of Registry  from 
// configuration file and checking the validity of it. return the password if it is valid.
// otherwise return the default password of Registry.
func getRegistryPassword(confContent *Config) string{
	registryPassword := os.Getenv("REGISTRYPASSWORD")
	if registryPassword != ""{
		return registryPassword
	}

	if confContent != nil  {
		if confContent.Registry.Credit.Password != "" {
			return confContent.Registry.Credit.Password
		}
	}

	return defaultConfig.Registry.Credit.Password
}

func appendErrs(dst []sysadmerror.Sysadmerror,from []sysadmerror.Sysadmerror)([]sysadmerror.Sysadmerror){

	dst = append(dst,from...)
		
	return dst
}

// Try to get the values of items of configuration from OS variables ,configuratio file or default value.
// The value of a item will be come from OS variables first ,then come from configuration file and last come from default value.
// All the values of items should be passed check when set it to ConfigDefined
func HandleConfig(configPath string, cmdRunPath string) (*Config,[]sysadmerror.Sysadmerror) {
	var confContent *Config = nil
	cfgFile,errs := getConfigPath(configPath,cmdRunPath)
	if cfgFile != "" {
		tmpConfContent,err := getConfigContent(cfgFile)
		if len(err) > 0 {
			errs =  appendErrs(errs,err) 
		}else {
			confContent = tmpConfContent
			e := checkVerIsValid(confContent.SysadmVersion)
			if len(e) >0  {
				errs = appendErrs(errs,e)
			}else{
				ConfigDefined.SysadmVersion = confContent.SysadmVersion
			}
		}
	}
	
	address,err := getServerAddress(confContent)
	ConfigDefined.Server.Address = address
	if len(err) > 0 {
		errs = appendErrs(errs,err)
	}

	ConfigDefined.Server.Port,err = getServerPort(confContent)
	if len(err) > 0 {
		errs = appendErrs(errs,err)
	}

	ConfigDefined.Server.Socket,err = getSockFile(confContent,cmdRunPath)
	if len(err) > 0 {
		errs = appendErrs(errs,err)
	}

	ConfigDefined.Log.AccessLog,err = getAccessLogFile(confContent,cmdRunPath)
	if len(err) > 0 {
		errs = appendErrs(errs,err)
	}

	ConfigDefined.Log.ErrorLog,err = getErrorLogFile(confContent,cmdRunPath)
	if len(err) > 0 {
		errs = appendErrs(errs,err)
	}

	ConfigDefined.Log.Kind,err = getLogKind(confContent)
	if len(err) > 0 {
		errs = appendErrs(errs,err)
	}

	ConfigDefined.Log.Level,err = getLogLevel(confContent)
	if len(err) > 0 {
		errs = appendErrs(errs,err)
	}

	ConfigDefined.Log.TimeStampFormat,err = getLogTimeFormat(confContent)
	if len(err) > 0 {
		errs = appendErrs(errs,err)
	}

	ConfigDefined.Log.SplitAccessAndError,err = getIsSplitLog(confContent)
	if len(err) > 0 {
		errs = appendErrs(errs,err)
	}

	ConfigDefined.User.DefaultUser = getDefaultUser(confContent)
	ConfigDefined.User.DefaultPassword = getDefaultPassword(confContent)
	ConfigDefined.DB.Host = getPostgreHost(confContent)
	ConfigDefined.DB.Port = getPostgrePort(confContent)
	ConfigDefined.DB.User = getPostgreUser(confContent)
	ConfigDefined.DB.Password = getPostgrePassword(confContent)
	ConfigDefined.DB.Dbname = getPostgreDBName(confContent)
	ConfigDefined.DB.Sslmode = getPostgreSslmode(confContent)
	if strings.ToLower(ConfigDefined.DB.Sslmode) != "disable" {
		ConfigDefined.DB.Sslrootcert = getPostgreSslrootcert(confContent)
		ret, err := checkFileExists(ConfigDefined.DB.Sslrootcert, cmdRunPath)
		if !ret {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201028,"warning","SslMode of Postgre has be set to %s But sslCA(%s) can not be found. We will try to set SslMode to disable. Error: %s",ConfigDefined.DB.Sslmode, ConfigDefined.DB.Sslrootcert,err))
			ConfigDefined.DB.Sslmode = "disable"
		}
	}

	if strings.ToLower(ConfigDefined.DB.Sslmode) != "disable" {
		ConfigDefined.DB.Sslkey = getPostgreSslkey(confContent)
		ret,err := checkFileExists(ConfigDefined.DB.Sslkey, cmdRunPath)
		if !ret {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201029,"warning","SslMode of Postgre has be set to %s But SslKey(%s) can not be found. We will try to set SslMode to disable error: %s",ConfigDefined.DB.Sslmode, ConfigDefined.DB.Sslkey,err))
			ConfigDefined.DB.Sslmode = "disable"
		}
	}

	if strings.ToLower(ConfigDefined.DB.Sslmode) != "disable" {
		ConfigDefined.DB.Sslcert = getPostgreSslcert(confContent)
		ret, err := checkFileExists(ConfigDefined.DB.Sslcert, cmdRunPath)
		if !ret {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201030,"warning","SslMode of Postgre has be set to %s But Sslcert(%s) can not be found. We will try to set SslMode to disable,error: %s ",ConfigDefined.DB.Sslmode, ConfigDefined.DB.Sslcert,err))
			ConfigDefined.DB.Sslmode = "disable"
		}
	}

	ConfigDefined.DB.DbMaxConnect = getPostgreMaxConnect(confContent)
	ConfigDefined.DB.DbIdleConnect = getPostgreDbIdleConnect(confContent)

	host,e := getRegistryHost(confContent)
	if e == nil {
		ConfigDefined.Registry.Server.Host = host
	}else{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(201031,"fatal","registry host(%s) is not valid.error:%s. ",host, e))
	}

	ConfigDefined.Registry.Server.Port = getRegistryPort(confContent)

	if strings.ToLower(ConfigDefined.Registry.Server.Sslmode) != "disable" {
		ConfigDefined.Registry.Server.Certs.Ca = getRegistrySslCa(confContent)
		ret, err := checkFileExists(ConfigDefined.Registry.Server.Certs.Ca, cmdRunPath)
		if !ret {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201032,"warning","SslMode of Registry has be set to %s But sslCA(%s) can not be found. We will try to set SslMode to disable. Error: %s",ConfigDefined.Registry.Server.Sslmode, ConfigDefined.Registry.Server.Certs.Ca,err))
			ConfigDefined.Registry.Server.Sslmode = "disable"
		}
	}

	if strings.ToLower(ConfigDefined.Registry.Server.Sslmode) != "disable" {
		ConfigDefined.Registry.Server.Certs.Key = getRegistrySslkey(confContent)
		ret,err := checkFileExists(ConfigDefined.Registry.Server.Certs.Key, cmdRunPath)
		if !ret {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201033,"warning","SslMode of Registry has be set to %s But Key(%s) file can not be found. We will try to set SslMode to disable error: %s",ConfigDefined.Registry.Server.Sslmode, ConfigDefined.Registry.Server.Certs.Key,err))
			ConfigDefined.Registry.Server.Sslmode = "disable"
		}
	}

	if strings.ToLower(ConfigDefined.Registry.Server.Sslmode) != "disable" {
		ConfigDefined.Registry.Server.Certs.Cert = getRegistrySslCert(confContent)
		ret,err := checkFileExists(ConfigDefined.Registry.Server.Certs.Cert, cmdRunPath)
		if !ret {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(201034,"warning","SslMode of Registry has be set to %s But Cert(%s) file can not be found. We will try to set SslMode to disable error: %s",ConfigDefined.Registry.Server.Sslmode, ConfigDefined.Registry.Server.Certs.Cert,err))
			ConfigDefined.Registry.Server.Sslmode = "disable"
		}
	}

	ConfigDefined.Registry.Credit.Username = getRegistryUser(confContent)
	ConfigDefined.Registry.Credit.Password = getRegistryPassword(confContent)

	return &ConfigDefined,errs
}

