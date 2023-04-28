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

ErrorCode 90200xxx
*/

package server

import (
	"path/filepath"

	"sysadm/config"
	"sysadm/sysadmerror"
	log "github.com/wangyysde/sysadmLog"
	"github.com/wangyysde/sysadmServer"
)

/*
handleConfig: 1. get the content of configuration file and parsed it
2. get the values of configurations from envirement, configuration and default value
*/
func handleConfig(cmdPath string) []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror

	cfgFile := CurrentRuningData.Options.CfgFile
	cfgFile, err := config.GetCfgFilePath(cfgFile, cmdPath)
	errs = append(errs, err...)

	fileConf := &Config{}
	if cfgFile != "" {
		CurrentRuningData.Options.CfgFile = cfgFile
		tmpConf, err := config.GetCfgContent(cfgFile, fileConf)
		if tmpConf == nil {
			fileConf = nil
		}
		errs = append(errs, err...)
	} else {
		fileConf = nil
	}

	workingPath, e := getBinRootPath(cmdPath)
	CurrentRuningData.Options.workingRoot = workingPath
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(90200001, "error", "Can not get working path , error: %s.", e))
	}

	defaultConf := getDefaultConf()
	if fileConf == nil {
		fileConf = defaultConf
	}

	definedConf := &CurrentRuningData.Config
	err = validateServerConf(definedConf, fileConf, defaultConf, cmdPath)
	errs = append(errs, err...)

	err = validateServerTlsConf(definedConf, fileConf, defaultConf, cmdPath)
	errs = append(errs, err...)

	err = validateServerLog(definedConf, fileConf, defaultConf, cmdPath)
	errs = append(errs, err...)

	err = validateDBServerConf(definedConf, fileConf, defaultConf, cmdPath)
	errs = append(errs, err...)

	err = validateApiServer(definedConf, fileConf, defaultConf, cmdPath)
	errs = append(errs, err...)

	return errs
}

/*
validateServerConf check the validation of server configuration.
*/
func validateServerConf(defindConf *Config, fileConf *Config, defaultConf *Config, cmdPath string) []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror
	var err []sysadmerror.Sysadmerror

	defindConf.Server.Address, errs = config.ValidateListenAddress(fileConf.Server.Address, defaultConf.Server.Address, "SERVERADDRESS")
	defindConf.Server.Port, err = config.ValidateListenPort(fileConf.Server.Port, defaultConf.Server.Port, "SERVERPORT")
	errs = append(errs, err...)
	defindConf.Server.Socket, err = config.ValidateListenSocket(fileConf.Server.Socket, defaultConf.Server.Socket, "SERVERSOCKET", cmdPath)
	errs = append(errs, err...)

	return errs
}

/*
validateServerTlsConf validate tls file for server Listen.
set the pasths of tls files to absolute path if istls is true.
otherwise set the pasths of tls files to  "".
*/
func validateServerTlsConf(definedConf *Config, fileConf *Config, defaultConf *Config, cmdPath string) []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror
	var err []sysadmerror.Sysadmerror

	definedConf.ServerTls.IsTls, err = config.ValidateIsTls(fileConf.ServerTls.IsTls, "SERVERISTLS")
	errs = append(errs, err...)

	if definedConf.ServerTls.IsTls {
		definedConf.ServerTls.Ca, err = config.ValidateTlsFile(fileConf.ServerTls.Ca, defaultConf.ServerTls.Ca, "SERVERISTLSCA", cmdPath)
		errs = append(errs, err...)
		definedConf.ServerTls.Cert, err = config.ValidateTlsFile(fileConf.ServerTls.Cert, defaultConf.ServerTls.Cert, "SERVERISTLSCERT", cmdPath)
		errs = append(errs, err...)
		definedConf.ServerTls.Key, err = config.ValidateTlsFile(fileConf.ServerTls.Key, defaultConf.ServerTls.Key, "SERVERISTLSKEY", cmdPath)
		errs = append(errs, err...)
		definedConf.ServerTls.InsecureSkipVerify, err = config.ValidateIsTls(fileConf.ServerTls.InsecureSkipVerify, "SERVERISTLSINSECURESKIPVERIFY")
		errs = append(errs, err...)
	}

	return errs
}

func validateServerLog(definedConf *Config, fileConf *Config, defaultConf *Config, cmdPath string) []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror
	var err []sysadmerror.Sysadmerror

	definedConf.Log.AccessLog, err = config.ValidateLogFile(fileConf.Log.AccessLog, defaultConf.Log.AccessLog, "ACCESSLOG", cmdPath)
	errs = append(errs, err...)

	definedConf.Log.ErrorLog, err = config.ValidateLogFile(fileConf.Log.ErrorLog, defaultConf.Log.ErrorLog, "ERRORLOG", cmdPath)
	errs = append(errs, err...)

	definedConf.Log.Kind, err = config.ValidateLogKind(fileConf.Log.Kind, defaultConf.Log.Kind, "ERRORLOGKIND")
	errs = append(errs, err...)

	definedConf.Log.Level, err = config.ValidateLogLevel(fileConf.Log.Level, defaultConf.Log.Level, "LOGLEVEL")
	errs = append(errs, err...)

	definedConf.Log.SplitAccessAndError, err = config.ValidateIsSplitLog(fileConf.Log.SplitAccessAndError, "ISSPLITLOG")
	errs = append(errs, err...)

	definedConf.Log.TimeStampFormat, err = config.ValidateLogTimeFormat(fileConf.Log.TimeStampFormat, defaultConf.Log.TimeStampFormat, "LOGTIMEFORMAT")
	errs = append(errs, err...)

	return errs
}

/*
getDefaultConf return a pointor point to &Config with the default value
*/
func getDefaultConf() *Config {
	return &Config{
		Version: CurrentRuningData.Config.Version,
		Server: config.Server{
			Address: serverAddress,
			Port:    serverPort,
			Socket:  serverSocket,
		},
		ServerTls: config.Tls{
			IsTls:              serverIsTls,
			Ca:                 serverCa,
			Cert:               serverCert,
			Key:                serverKey,
			InsecureSkipVerify: serverInsecureSkipVerify,
		},
		Log: config.Log{
			AccessLog:           serverAccessLog,
			ErrorLog:            serverErrorLog,
			Kind:                serverLogKind,
			Level:               serverLogLevel,
			SplitAccessAndError: serverLogSplitAccessAndError,
			TimeStampFormat:     serverLogTimeStampFormat,
		},
		DB: DBServer{
			Type:   dbType,
			DBName: dbName,
			Server: config.Server{
				Address: dbServerAddress,
				Port:    dbServerPort,
				Socket:  dbServerSocket,
			},
			MaxOpenConns: dbMaxOpenConns,
			MaxIdleConns: dbMaxIdeleConns,
			Tls: config.Tls{
				IsTls:              dbServerIsTls,
				Ca:                 dbServerCa,
				Cert:               dbServerCert,
				Key:                dbServerKey,
				InsecureSkipVerify: dbServerInsecureSkipVerify,
			},
		},
		ApiServer: ApiServer{
			Server: config.Server{
				Address: apiServerAddress,
				Port:    apiServerPort,
			},
			Tls: config.Tls{
				IsTls:              apiServerIsTls,
				Ca:                 apiServerCa,
				Cert:               apiServerCert,
				Key:                apiServerKey,
				InsecureSkipVerify: apiServerInsecureSkipVerify,
			},
			ApiVersion: apiVersion,
		},
	}
}

/*
validateDBServerConf check the validation of DB server configuration.
*/
func validateDBServerConf(definedConf *Config, fileConf *Config, defaultConf *Config, cmdPath string) []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror
	var err []sysadmerror.Sysadmerror

	definedConf.DB.Type, err = config.ValidateDBType(fileConf.DB.Type, defaultConf.DB.Type, "DBTYPE")
	errs = append(errs, err...)

	definedConf.DB.DBName, err = config.ValidateDBName(fileConf.DB.DBName, defaultConf.DB.DBName, "DBNAME", defaultConf.DB.Type)
	errs = append(errs, err...)

	definedConf.DB.Server.Address, errs = config.ValidateServerAddress(fileConf.DB.Server.Address, defaultConf.DB.Server.Address, "DBADDRESS")
	errs = append(errs, err...)

	definedConf.DB.Server.Port, errs = config.ValidateServerPort(fileConf.DB.Server.Port, defaultConf.DB.Server.Port, "DBPORT")
	errs = append(errs, err...)

	definedConf.DB.Server.Socket, errs = config.ValidateServerSocket(fileConf.DB.Server.Socket, defaultConf.DB.Server.Socket, "DBSOCKET", cmdPath)
	errs = append(errs, err...)

	definedConf.DB.Tls.IsTls, err = config.ValidateIsTls(fileConf.DB.Tls.IsTls, "DBISTLS")
	errs = append(errs, err...)

	if definedConf.DB.Tls.IsTls {
		definedConf.DB.Tls.Ca, err = config.ValidateTlsFile(fileConf.DB.Tls.Ca, defaultConf.DB.Tls.Ca, "DBTLSCA", cmdPath)
		errs = append(errs, err...)
		definedConf.DB.Tls.Cert, err = config.ValidateTlsFile(fileConf.DB.Tls.Cert, defaultConf.DB.Tls.Cert, "DBTLSCERT", cmdPath)
		errs = append(errs, err...)
		definedConf.DB.Tls.Key, err = config.ValidateTlsFile(fileConf.DB.Tls.Key, defaultConf.DB.Tls.Key, "DBTLSKEY", cmdPath)
		errs = append(errs, err...)
		definedConf.DB.Tls.InsecureSkipVerify, err = config.ValidateIsTls(fileConf.DB.Tls.InsecureSkipVerify, "DBTLSINSECURESKIPVERIFY")
		errs = append(errs, err...)
	}

	definedConf.DB.Credit.UserName, err = config.ValidateUser(fileConf.DB.Credit.UserName, defaultConf.DB.Credit.UserName, "DBUSERNAME")
	errs = append(errs, err...)

	definedConf.DB.Credit.Password, err = config.ValidateUser(fileConf.DB.Credit.Password, defaultConf.DB.Credit.Password, "DBPASSWORD")
	errs = append(errs, err...)

	definedConf.DB.MaxOpenConns, err = config.ValidateConns(fileConf.DB.MaxOpenConns, defaultConf.DB.MaxOpenConns, "DBMAXOPENCONNS")
	errs = append(errs, err...)

	definedConf.DB.MaxIdleConns, err = config.ValidateConns(fileConf.DB.MaxIdleConns, defaultConf.DB.MaxIdleConns, "DBMAXIDLECONNS")
	errs = append(errs, err...)

	return errs

}

// Get the absolute path to the working root
func getBinRootPath(cmdPath string) (string, error) {
	dir, error := filepath.Abs(filepath.Dir(cmdPath))
	if error != nil {
		return "", error
	}

	dir = filepath.Join(dir, "../")

	return dir, nil
}

// set parameters to accessLogger and errorLooger
func setLogger() []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror

	definedConf := &CurrentRuningData.Config

	sysadmServer.SetLoggerKind(definedConf.Log.Kind)
	sysadmServer.SetLogLevel(definedConf.Log.Level)
	sysadmServer.SetTimestampFormat(definedConf.Log.TimeStampFormat)
	if definedConf.Log.AccessLog != "" {
		_, fp, err := sysadmServer.SetAccessLogFile(definedConf.Log.AccessLog)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(102000, "error", "can not set access log file(%s) error: %s", definedConf.Log.AccessLog, err))
		} else {
			definedConf.Log.AccessLogFp = fp
		}

	}

	if definedConf.Log.SplitAccessAndError && definedConf.Log.ErrorLog != "" {
		_, fp, err := sysadmServer.SetErrorLogFile(definedConf.Log.ErrorLog)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(102001, "error", "can not set error log file(%s) error: %s", definedConf.Log.ErrorLog, err))
		} else {
			definedConf.Log.ErrorLogFp = fp
		}
	}
	sysadmServer.SetIsSplitLog(definedConf.Log.SplitAccessAndError)

	level, e := log.ParseLevel(definedConf.Log.Level)
	if e != nil {
		sysadmServer.SetMode(sysadmServer.DebugMode)
	} else {
		if level >= log.DebugLevel {
			sysadmServer.SetMode(sysadmServer.DebugMode)
		} else {
			sysadmServer.SetMode(sysadmServer.ReleaseMode)
		}
	}

	return errs
}

// close access log file descriptor and error log file descriptor
// set AccessLogger  and ErrorLogger to nil
func closeLogger() {
	definedConf := &CurrentRuningData.Config
	if definedConf.Log.AccessLogFp != nil {
		fp := definedConf.Log.AccessLogFp
		fp.Close()
		definedConf.Log.AccessLogFp = nil
		definedConf.Log.AccessLog = ""
	}

	if definedConf.Log.ErrorLogFp != nil {
		fp := definedConf.Log.ErrorLogFp
		fp.Close()
		definedConf.Log.ErrorLogFp = nil
		definedConf.Log.ErrorLog = ""
	}
}

/*
validateApiServer check the validation of apiServer configuration.
*/
func validateApiServer(defindConf *Config, fileConf *Config, defaultConf *Config, cmdPath string) []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror
	var err []sysadmerror.Sysadmerror

	defindConf.ApiServer.Server.Address, errs = config.ValidateServerAddress(fileConf.ApiServer.Server.Address, defaultConf.ApiServer.Server.Address, "APISERVERADDRESS")

	defindConf.ApiServer.Server.Port, err = config.ValidateServerPort(fileConf.ApiServer.Server.Port, defaultConf.ApiServer.Server.Port, "APISERVERPORT")
	errs = append(errs, err...)

	defindConf.ApiServer.Tls.IsTls, err = config.ValidateIsTls(fileConf.ApiServer.Tls.IsTls, "APISERVERISTLS")
	errs = append(errs, err...)

	if defindConf.ApiServer.Tls.IsTls {
		defindConf.ApiServer.Tls.Ca, err = config.ValidateTlsFile(fileConf.ApiServer.Tls.Ca, defaultConf.ApiServer.Tls.Ca, "APISERVERTLSCA", cmdPath)
		errs = append(errs, err...)
		defindConf.ApiServer.Tls.Cert, err = config.ValidateTlsFile(fileConf.ApiServer.Tls.Cert, defaultConf.ApiServer.Tls.Cert, "APISERVERTLSCERT", cmdPath)
		errs = append(errs, err...)
		defindConf.ApiServer.Tls.Key, err = config.ValidateTlsFile(fileConf.ApiServer.Tls.Key, defaultConf.ApiServer.Tls.Key, "APISERVERTLSKEY", cmdPath)
		errs = append(errs, err...)
		defindConf.ApiServer.Tls.InsecureSkipVerify, err = config.ValidateIsTls(fileConf.ApiServer.Tls.InsecureSkipVerify, "APISERVERTLSINSECURESKIPVERIFY")
		errs = append(errs, err...)
	}

	defindConf.ApiServer.ApiVersion, err = config.ValidateApiVersion(fileConf.ApiServer.ApiVersion, defaultConf.ApiServer.ApiVersion, "APIVERSION")
	errs = append(errs, err...)

	return errs
}
