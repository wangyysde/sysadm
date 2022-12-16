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

ErrorCode 1001xxx
*/

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
	"github.com/wangyysde/sysadmServer"
	"github.com/wangyysde/yaml"
)

/*
  Get the absolute path of configurationn file and try to open it
  cfgPath: the path of configuration file user specified
  cmdRunPath: args[0]
*/
func GetCfgFilePath(cfgPath string, cmdRunPath string) (string, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	
	dir ,error := filepath.Abs(filepath.Dir(cmdRunPath))
	if error != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001001,"error","Can not get absolute path for promgram. error: %s.",error))
		return "",errs
	}

	if cfgPath == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001002,"error","configuration file path is empty."))
		return "",errs
	}

	if ! filepath.IsAbs(cfgPath) {
		tmpDir := filepath.Join(dir,"../")
		cfgPath = filepath.Join(tmpDir,cfgPath)
	}

	fp, err := os.Open(cfgPath)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001003,"error","Can not open the configuration file:%s, error: %s.",cfgPath,err))
		return "",errs
	}
	fp.Close()

	return cfgPath,errs
}


/*
	GetCfgContent get the content of the configuration file specified by configPath.
	Note, configPath must be an absolute path of configuration file 
	Then yaml.Unmarshal the content into o
	Retrun o and  []sysadmerror.Sysadmerror if no error was occurred.
	otherwise return nil and []sysadmerror.Sysadmerror
*/

func GetCfgContent(configPath string,o interface{}) (interface{},[]sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	if configPath == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001004,"warn","The configration file path is empty."))
		return nil, errs
	}

	yamlContent, err := ioutil.ReadFile(configPath) 
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001005,"error","Can not read configuration file: %s error: %s.",configPath,err))
		return nil, errs
	}

	err = yaml.Unmarshal(yamlContent, o)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001006,"error","Can not Unmarshal configuration contenet error: %s.",err))
		return nil, errs
	}

	return o, errs
}

/* 
	ValidateListenAddress 
	1. get the value of the environment variable named envName if envName is not empty, the check it is valid
	   return the value of the environment variable named envName if it passed check.
	2. check the validity of the value of confValue and return it if it is valid.
	3. otherwise return defaultValue
*/
func ValidateListenAddress(confValue string,defaultValue string, envName string )(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		tmpAddress := os.Getenv(envName)
		if strings.TrimSpace(tmpAddress) != "" {
			ip,err := utils.CheckIpAddress(tmpAddress,true)
			errs = append(errs,err...)
			if ip != nil {
				return tmpAddress,errs
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001007,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid server address(%s)",envName,tmpAddress,envName,err))
		}
	}

	if strings.TrimSpace(confValue) != "" {
		ip,err := utils.CheckIpAddress(confValue,true)
		errs = append(errs,err...)
		if ip != nil {
			return confValue,errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001008,"warning","Address of server has be found in the configuration file, but the value(%s) is not a valid server address(%s)",confValue,err))
	}

	return defaultValue,errs
}

/*
	ValidateListenPort 
	1. get the value of the environment variable named envName if envName is not empty, the check it is valid
	   return the value of the environment variable named envName if it passed check.
	2. check the validity of the value of confValue and return it if it is valid.
	3. otherwise return defaultValue
*/
func ValidateListenPort(confValue int,defaultValue int, envName string )(int,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		tmpPort := os.Getenv(envName)
		if strings.TrimSpace(tmpPort) != "" {
			tmpPortInt,e := strconv.Atoi(tmpPort)
			if e == nil {
				port,err := utils.CheckPort(tmpPortInt)
				errs = append(errs,err...)
				if port > 0 {
					return port,errs
				}
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001009,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid server port",envName,tmpPort,envName))
		}
	}

	if confValue > 1024 && confValue <= 65536 {
		return confValue,errs
	}
	
	return defaultValue,errs
}

/*
	ValidateListenSocket 
	1. get the value of the environment variable named envName if envName is not empty, and check it is valid
	   return the absolute path of socket file if it passed check.
	2. check the validity of the value of confValue and return the absolute path of socket file if it is valid.
	3. check the validity of the value of defaultValue and return the absolute path of socket file if it is valid.
	4. otherwrise return ""
*/
func ValidateListenSocket(confValue string,defaultValue string, envName string,cmdRunPath string)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		socket := os.Getenv(envName)
		if strings.TrimSpace(socket) != "" {
			tmpSocket,_ := utils.CheckFileRW(socket,cmdRunPath,true)
			if tmpSocket != "" {
				return tmpSocket, errs
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001010,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid socket file",envName,socket,envName))
		}

	}

	if strings.TrimSpace(confValue) != "" {
		tmpSocket,_ := utils.CheckFileRW(confValue,cmdRunPath,true)
		if tmpSocket != "" {
			return tmpSocket, errs
		}

		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001011,"warning","socket file(%s) set in the configuration file  is not a valid socket file",confValue))
	}

	tmpSocket,_ := utils.CheckFileRW(defaultValue,cmdRunPath,true)
	if tmpSocket != "" {
		return tmpSocket, errs
	}
	
	return "", errs
}

/*
	ValidateIsTls 
	1. get the value of the environment variable named envName if envName is not empty, and check it is the bool value
	   return the value it passed check. 
	2.otherwise return confValue 
*/
func ValidateIsTls(confValue bool, envName string)(bool,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		isTls := strings.TrimSpace(os.Getenv(envName))
		if strings.TrimSpace(isTls) != "" {
			if strings.ToLower(isTls) == "true" || strings.ToLower(isTls) == "yes" || strings.ToLower(isTls) == "1"{
				return true,errs
			}

			if strings.ToLower(isTls) == "false" || strings.ToLower(isTls) == "no" || strings.ToLower(isTls) == "0"{
				return false,errs
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001012,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid isTls",envName,isTls,envName))
		}
	}

	return confValue,errs

}

/*
	ValidateTlsFile 
	1. get the value of the environment variable named envName if envName is not empty, and check it is valid
	   return the absolute path of PKI file if it passed check.
	2. check the validity of the value of confValue and return the absolute path of PKI file if it is valid.
	3. check the validity of the value of defaultValue and return the absolute path of socket file if it is valid.
	4. otherwrise return ""
*/
func ValidateTlsFile(confValue string,defaultValue string, envName string,cmdRunPath string)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		tlsFile := os.Getenv(envName)
		if strings.TrimSpace(tlsFile) != "" {
			tmpTlsFile,_ := utils.CheckFileIsRead(tlsFile,cmdRunPath)
			if tmpTlsFile != "" {
				return tmpTlsFile, errs
			}
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001013,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid TLS file",envName,tlsFile,envName))
	}

	if strings.TrimSpace(confValue) != "" {
		tmpTlsFile,_ := utils.CheckFileIsRead(confValue,cmdRunPath)
		if tmpTlsFile != "" {
			return tmpTlsFile, errs
		}

		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001014,"warning","file(%s) set in the configuration file  is not a valid PKI file",confValue))
	}

	tmpTlsFile, _:= utils.CheckFileIsRead(defaultValue,cmdRunPath)
	if tmpTlsFile != "" {
		return tmpTlsFile, errs
	}
	
	return "", errs
}


/*
	ValidateLogFile 
	1. get the value of the environment variable named envName if envName is not empty, and check it is valid
	   return the absolute path of log file if it passed check.
	2. check the validity of the value of confValue and return the absolute path of log file if it is valid.
	3. check the validity of the value of defaultValue and return the absolute path of log file if it is valid.
	4. otherwrise return ""
*/
func ValidateLogFile(confValue string,defaultValue string, envName string,cmdRunPath string)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		logFile := os.Getenv(envName)
		if strings.TrimSpace(logFile) != "" {
			tmpLogfile,_ := utils.CheckFileRW(logFile,cmdRunPath,false)
			if tmpLogfile != "" {
				return tmpLogfile, errs
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001015,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid log file",envName,logFile,envName))
		}
		
	}

	if strings.TrimSpace(confValue) != "" {
		tmpLogfile,_ := utils.CheckFileRW(confValue,cmdRunPath,false)
		if tmpLogfile != "" {
			return tmpLogfile, errs
		}

		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001016,"warning","log file(%s) set in the configuration file  is not a valid log file",confValue))
	}

	tmpLogfile,_ := utils.CheckFileRW(defaultValue,cmdRunPath,false)
	if tmpLogfile != "" {
		return tmpLogfile, errs
	}
	
	return "", errs
}

/* 
   Try to get log kind  from one of  envName,configuration file or default value.
   The order for getting log kind is envName,configuration file and default value.
*/
func ValidateLogKind(confValue string,defaultValue string, envName string)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	logKind := strings.TrimSpace(os.Getenv(envName))
	if logKind != ""{
		if strings.ToLower(logKind) == "text" || strings.ToLower(logKind) == "json" {
			return strings.ToLower(logKind), errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001017,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid log kind",envName,logKind,envName))
	}

	if strings.TrimSpace(confValue) != "" {
		if strings.ToLower(confValue) == "text" || strings.ToLower(confValue) == "json" {
			return strings.ToLower(confValue),errs
		}
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001018,"warning","The default log kind %s will be used",defaultValue))
	return defaultValue,errs
}

/*
  Try to get log level from one of envName, configuration file or default value.
  The order for getting log level is envName, configuration file and default value.
*/
func ValidateLogLevel(confValue string,defaultValue string, envName string) (string,[]sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	logLevel := strings.TrimSpace(os.Getenv(envName))
	if logLevel != "" {
		if checkLogLevel(logLevel) {
			return strings.ToLower(logLevel), errs
		}
		
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001019,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid log level",envName,logLevel,envName))
	}

	confValue = strings.TrimSpace(confValue)
	if confValue != "" {
		if checkLogLevel(confValue) {
			return strings.ToLower(confValue), errs
		}
		
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001020,"warning","The Log level(%s) has been set in the configuration file environment is not a valid log level",confValue))
	}

	return defaultValue,errs
}

/* 
   CheckLogLevel check level if is a log level string.
   return true if it is a log level string otherwise return false
*/
func CheckLogLevel(level string) bool {
	return checkLogLevel(level)
}
/* 
   checkLogLevel check level if is a log level string.
   return true if it is a log level string otherwise return false
*/
func checkLogLevel(level string) bool {
	if len(level) < 1 {
		return false
	}

	for _,l := range sysadmServer.Levels {
		if strings.ToLower(level) == strings.ToLower(l) {
			return true
		}
	}

	return false
}

/*
	ValidateIsSplitLog 
	1. get the value of the environment variable named envName if envName is not empty, and check it is the bool value
	   return the value it passed check. 
	2.otherwise return confValue 
*/
func ValidateIsSplitLog(confValue bool, envName string)(bool,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		isSplitLog := strings.TrimSpace(os.Getenv(envName))
		if strings.TrimSpace(isSplitLog) != "" {
			if strings.ToLower(isSplitLog) == "true" || strings.ToLower(isSplitLog) == "yes" || strings.ToLower(isSplitLog) == "1"{
				return true,errs
			}

			if strings.ToLower(isSplitLog) == "false" || strings.ToLower(isSplitLog) == "no" || strings.ToLower(isSplitLog) == "0"{
				return false,errs
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001021,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid isSplitLog",envName,isSplitLog,envName))
		}
	}

	return confValue,errs

}

/*
  Try to get log time format from one of envName, configuration file or default value.
  The order for getting log time format is envName, configuration file and default value.
*/
func ValidateLogTimeFormat(confValue string,defaultValue string, envName string)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		format := strings.TrimSpace(os.Getenv(envName))
		if strings.TrimSpace(format) != "" {
			if checkLogTimeFormat(strings.TrimSpace(format)){
				return strings.TrimSpace(format),errs
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001022,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid log time format",envName,format,envName))
		}
		
	}

	if strings.TrimSpace(confValue) != "" {
		if checkLogTimeFormat(strings.TrimSpace(confValue)){
			return strings.TrimSpace(confValue),errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001023,"warning","The Log time format(%s) has been set in the configuration file,but it is not a valid log time format",confValue))
	}

	return defaultValue,errs
}

/*
   checkLogTimeFormat check valid of log format.
   return true if format is a log time format string otherwise return false
*/
func checkLogTimeFormat(format string) bool{
	if len(format) < 1 {
		return false
	}

	for _,v := range sysadmServer.TimestampFormat {
		if strings.EqualFold(format,v) {
			return true
		}
	}

	return false
}

/*
  Try to get log DB type from one of envName, configuration file or default value.
  The order for getting DB type is envName, configuration file and default value.
*/
func ValidateDBType(confValue string,defaultValue string, envName string)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		dbType := strings.TrimSpace(os.Getenv(envName))
		if strings.TrimSpace(dbType) != "" {
			if db.IsSupportedDB(dbType) {
				return strings.ToLower(strings.TrimSpace(dbType)),errs
			}
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001024,"warning","The environment variable %s has be found, but the value(%s) of %s is not a DB type value",envName,dbType,envName))

	}

	confValue = strings.ToLower(strings.TrimSpace(confValue))
	if confValue != "" {
		if db.IsSupportedDB(confValue) {
			return strings.ToLower(strings.TrimSpace(confValue)),errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001025,"warning","The DB type has been set in the configuration file,but it is not a valid DB type",confValue))
	}

	return defaultValue,errs

}

/*
  Try to get DB name from one of envName, configuration file or default value.
  The order for getting DB name is envName, configuration file and default value.
*/
func ValidateDBName(confValue string,defaultValue string, envName string,dbType string)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		dbName := strings.TrimSpace(os.Getenv(envName))
		if strings.TrimSpace(dbName) != "" {
			if db.CheckIdentifier(dbType,dbName) {
				return strings.TrimSpace(dbType),errs
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001026,"warning","The environment variable %s has be found, but the value(%s) of %s is not a DB name value",envName,dbType,envName))
		}
	}

	if strings.TrimSpace(confValue) != "" {
		if db.CheckIdentifier(dbType,strings.TrimSpace(confValue)) {
			return strings.TrimSpace(confValue),errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001027,"warning","The DB name has been set in the configuration file,but it is not a valid DB name",confValue))
	} 

	return defaultValue,errs
}

/* 
	ValidateServerAddress 
	1. get the value of the environment variable named envName if envName is not empty, the check it is valid
	   return the value of the environment variable named envName if it passed check.
	2. check the validity of the value of confValue and return it if it is valid.
	3. otherwise return defaultValue
*/
func ValidateServerAddress(confValue string,defaultValue string, envName string )(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		tmpAddress := os.Getenv(envName)
		if strings.TrimSpace(tmpAddress) != "" {
			ip,err := utils.CheckIpAddress(tmpAddress,false)
			errs = append(errs,err...)
			if ip != nil {
				return tmpAddress,errs
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001028,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid server address(%s)",envName,tmpAddress,envName,err))
		}
	}

	if strings.TrimSpace(confValue) != "" {
		ip,err := utils.CheckIpAddress(confValue,false)
		errs = append(errs,err...)
		if ip != nil {
			return confValue,errs
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001029,"warning","Address of server has be found in the configuration file, but the value(%s) is not a valid server address(%s)",confValue,err))
	}

	return defaultValue,errs
}

/*
	ValidateServerPort 
	1. get the value of the environment variable named envName if envName is not empty, the check it is valid
	   return the value of the environment variable named envName if it passed check.
	2. check the validity of the value of confValue and return it if it is valid.
	3. otherwise return defaultValue
*/
func ValidateServerPort(confValue int,defaultValue int, envName string )(int,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		tmpPort := os.Getenv(envName)
		if strings.TrimSpace(tmpPort) != "" {
			tmpPortInt,e := strconv.Atoi(tmpPort)
			if e == nil {
				port,err := utils.CheckPort(tmpPortInt)
				errs = append(errs,err...)
				if port > 0 {
					return port,errs
				}
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001030,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid server port",envName,tmpPort,envName))
		}
	}

	if confValue > 1024 && confValue <= 65536 {
		return confValue,errs
	}
	
	return defaultValue,errs
}

/*
	ValidateServerSocket 
	1. get the value of the environment variable named envName if envName is not empty, and check it is valid
	   return the absolute path of socket file if it passed check.
	2. check the validity of the value of confValue and return the absolute path of socket file if it is valid.
	3. check the validity of the value of defaultValue and return the absolute path of socket file if it is valid.
	4. otherwrise return ""
*/
func ValidateServerSocket(confValue string,defaultValue string, envName string,cmdRunPath string)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		socket := os.Getenv(envName)
		if strings.TrimSpace(socket) != "" {
			tmpSocket,_ := utils.CheckFileIsRead(socket,cmdRunPath)
			if tmpSocket != "" {
				return tmpSocket, errs
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001031,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid socket file",envName,socket,envName))
		}
	}

	if strings.TrimSpace(confValue) != "" {
		tmpSocket,_ := utils.CheckFileIsRead(confValue,cmdRunPath)
		if tmpSocket != "" {
			return tmpSocket, errs
		}

		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001032,"warning","socket file(%s) set in the configuration file  is not a valid socket file",confValue))
	}

	tmpSocket,_ := utils.CheckFileIsRead(defaultValue,cmdRunPath)
	if tmpSocket != "" {
		return tmpSocket, errs
	}
	
	return "", errs
}

/*
  Try to get user from one of envName, configuration file or default value.
  The order for getting user is envName, configuration file and default value.
*/
func ValidateUser(confValue string,defaultValue string, envName string)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		dbUser := strings.TrimSpace(os.Getenv(envName))
		matched,err := regexp.MatchString("^[a-zA-Z0-9]{1,64}",dbUser)
		if matched {
			return dbUser, errs
		}

		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001031,"warning","The environment variable %s has be found, but the value(%s) of %s is not a DB user",envName,dbUser,envName))

		}
	}

	if strings.TrimSpace(confValue) != "" {
		matched,_ := regexp.MatchString("^[a-zA-Z0-9]{1,64}",confValue)
		if matched {
			return confValue, errs
		}

		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001032,"warning","The DB User has been set in the configuration file,but it is not a valid DB name",confValue))
	} 

	return defaultValue,errs
}

/*
  Try to get password from one of envName, configuration file or default value.
  The order for getting password is envName, configuration file and default value.
*/
func ValidatePassword(confValue string,defaultValue string, envName string)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		password := strings.TrimSpace(os.Getenv(envName))
		if len(password) > 0 && len(password) < 65 {
			return password,errs
		}
		
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001033,"warning","The environment variable %s has be found, but the value(%s) of %s is not a password",envName,password,envName))

	}

	if strings.TrimSpace(confValue) != "" {
		if len(strings.TrimSpace(confValue)) > 0 && len(strings.TrimSpace(confValue)) < 65 {
			return strings.TrimSpace(confValue),errs
		}
		
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001034,"warning","The password has been set in the configuration file,but it is not a valid password ",confValue))
	} 

	return defaultValue,errs
}

/*
	ValidateConns 
	1. get the value of the environment variable named envName if envName is not empty, the check it is valid
	   return the value of the environment variable named envName if it passed check.
	2. check the validity of the value of confValue and return it if it is valid.
	3. otherwise return defaultValue
*/
func ValidateConns(confValue int,defaultValue int, envName string )(int,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(envName) != "" {
		tmpConns := os.Getenv(envName)
		if strings.TrimSpace(tmpConns) != "" {
			tmpConnsInt,e := strconv.Atoi(tmpConns)
			if e == nil {
				if tmpConnsInt > 0 {
					return tmpConnsInt,errs
				}
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1001035,"warning","The environment variable %s has be found, but the value(%s) of %s is not a valid value",envName,tmpConns,envName))
		}
	}

	if confValue > 0  {
		return confValue,errs
	}
	
	return defaultValue,errs
}
