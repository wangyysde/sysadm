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
*
 */

package app

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"sysadm/config"
	"sysadm/sysadmerror"
	sysadmUtils "sysadm/utils"
)

// we define fileConf as global variable as it is used by a few of functions.
var fileConf *FileConf = &FileConf{}

/*
handleNotInConfFile handler the configuration items which can not define in define file,
such as working dir, configuration file path.
*/
func handleNotInConfFile() []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080001, "debug", "try to get working dir"))
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080002, "fatal", "get working dir error %s", err))
		return errs
	}
	workingDir := filepath.Join(binPath, "../")
	RunConf.WorkingDir = workingDir

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080003, "debug", "checking configuration file path"))
	var cfgFile string = ""
	if strings.TrimSpace(CliOps.CfgFile) != "" {
		if filepath.IsAbs(CliOps.CfgFile) {
			fp, err := os.Open(CliOps.CfgFile)
			if err != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080004, "fatal", "can not open configuration file %s error %s", CliOps.CfgFile, err))
				return errs
			}
			fp.Close()
		} else {
			configPath := filepath.Join(workingDir, CliOps.CfgFile)
			fp, err := os.Open(configPath)
			if err != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080005, "fatal", "can not open configuration file %s error %s", configPath, err))
				return errs
			}
			fp.Close()
			CliOps.CfgFile = configPath
		}
		cfgFile = CliOps.CfgFile
	} else {
		// try to get configuration file from default path
		configPath := filepath.Join(workingDir, DefaultConf)
		fp, err := os.Open(configPath)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080006, "fatal", "can not open configuration file %s error %s", configPath, err))
			return errs
		}
		fp.Close()
		CliOps.CfgFile = configPath
		cfgFile = configPath
	}
	RunConf.CfgFile = cfgFile

	return errs

}

/*
handle configuration items set in global block.
*/
func handleGlobalBlock() []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081001, "debug", "try to handle configuration items in global block"))

	if strings.TrimSpace(RunConf.CfgFile) != "" {
		_, tmpErrs := config.GetCfgContent(RunConf.CfgFile, fileConf)
		errs = append(errs, tmpErrs...)
	}

	ret := validateTlsConf(CliOps.Global.Tls, fileConf.Global.Tls, "global")
	RunConf.Global.Tls = *ret

	if RunConf.Global.Tls.IsTls && (RunConf.Global.Tls.Ca == "" || RunConf.Global.Tls.Cert == "" || RunConf.Global.Tls.Key == "") {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081002, "fatal", "IsTls has be set to true, but ca, cert or key file is not found."))
		return errs
	}

	retServer, err := validateServerConf(CliOps.Global.Server, fileConf.Global.Server, "global", false)
	errs = append(errs, err...)
	tmpIp, _ := sysadmUtils.CheckIpAddress(retServer.Address, false)
	tmpPort, _ := sysadmUtils.CheckPort(retServer.Port)
	if tmpIp == nil || tmpPort == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081004, "fatal", "the server address %s or server port %d is not valid.", retServer.Address, retServer.Port))
		return errs
	}
	RunConf.Global.Server = *retServer

	retLog, err := validateLogConf(CliOps.Global.Log, fileConf.Global.Log, "global")
	RunConf.Global.Log = *retLog
	errs = append(errs, err...)

	RunConf.Global.DebugMode = validateRunMode(CliOps.Global.DebugMode, fileConf.Global.DebugMode, "global")
	if RunConf.Global.DebugMode && (sysadmerror.GetLevelNum(RunConf.Global.Log.Level) < sysadmerror.GetLevelNum("debug")) {
		RunConf.Global.Log.Level = "debug"
	}

	RunConf.Global.NodeIdentifer = validateNodeIdentifer(CliOps.Global.NodeIdentifer, fileConf.Global.NodeIdentifer, "global")
	RunConf.Global.Uri = validateUri(CliOps.Global.Uri, fileConf.Global.Uri, "global")
	RunConf.Global.CommandStatusUri = validateListenUri(CliOps.Global.CommandStatusUri, fileConf.Global.CommandStatusUri,"global","commandStatusUri")
	RunConf.Global.CommandLogsUri = validateListenUri(CliOps.Global.CommandLogsUri, fileConf.Global.CommandLogsUri, "global", "commandLogsUri")
	RunConf.Global.SourceIP = validateSourceIP(CliOps.Global.SourceIP, fileConf.Global.SourceIP, "global")

	return errs
}

/*
validateTlsConf validate the tls values in cliConf (set by command line flags), fileConf(set by configuration file) and envBlock (set by environment)
the priority of the valules for secret is higher the priority for insecret.
the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher
the priority of the defination in configuration file.
*/
func validateTlsConf(cliConf config.Tls, fileConf config.Tls, envBlock string) *config.Tls {
	ret := &config.Tls{
		IsTls:              false,
		Ca:                 "",
		Cert:               "",
		Key:                "",
		InsecureSkipVerify: true,
	}

	envIsTls := false
	var envMap map[string]string
	envMapP := getEnvDefineForBlock(envBlock)
	if envMapP == nil {
		envMap = map[string]string{}
	} else {
		envMap = *envMapP
	}

	if envName, ok := envMap["IsTls"]; ok {
		isTlsValue := os.Getenv(envName)
		isTlsValue = strings.ToLower(strings.TrimSpace(isTlsValue))
		if isTlsValue == "yes" || isTlsValue == "y" || isTlsValue == "on" || isTlsValue == "1" {
			envIsTls = true
		}
	}

	if envIsTls {
		envCaName, okCa := envMap["Ca"]
		envCertName, okCert := envMap["Cert"]
		envKeyName, okKey := envMap["Key"]
		envInsecureSkipVerifyName, okInsecureSkipVerify := envMap["InsecureSkipVerify"]
		if okCa && okCert && okKey && okInsecureSkipVerify {
			caValue := strings.TrimSpace(os.Getenv(envCaName))
			certValue := strings.TrimSpace(os.Getenv(envCertName))
			keyValue := strings.TrimSpace(os.Getenv(envKeyName))
			insecureSkipVerifyValue := strings.ToLower(strings.TrimSpace(os.Getenv(envInsecureSkipVerifyName)))

			caExist, _ := sysadmUtils.CheckFileExists(caValue, "")
			certExist, _ := sysadmUtils.CheckFileExists(certValue, "")
			keyExist, _ := sysadmUtils.CheckFileExists(keyValue, "")
			if caExist && certExist && keyExist {
				if insecureSkipVerifyValue == "yes" || insecureSkipVerifyValue == "y" || insecureSkipVerifyValue == "on" || insecureSkipVerifyValue == "1" {
					ret.InsecureSkipVerify = true
				} else {
					ret.InsecureSkipVerify = false
				}
				ret.Ca = caValue
				ret.Cert = certValue
				ret.Key = keyValue
				ret.IsTls = true
			}
		}
	}

	// 1. the priority for secret is higher the priority for insecret.
	// 2.the priority of the defination in configuration file is higher than enverionments.
	if fileConf.IsTls {
		caValue := strings.ToLower(strings.TrimSpace(fileConf.Ca))
		certValue := strings.ToLower(strings.TrimSpace(fileConf.Cert))
		keyValue := strings.ToLower(strings.TrimSpace(fileConf.Key))
		insecureSkipVerifyValue := fileConf.InsecureSkipVerify

		_, _ = sysadmUtils.CheckFileExists(caValue, "")
		certExist, _ := sysadmUtils.CheckFileExists(certValue, "")
		keyExist, _ := sysadmUtils.CheckFileExists(keyValue, "")
		if certExist && keyExist {
			ret.InsecureSkipVerify = insecureSkipVerifyValue
			ret.Ca = caValue
			ret.Cert = certValue
			ret.Key = keyValue
			ret.IsTls = true
		}
	}

	// 1. the priority for secret is higher the priority for insecret.
	// 2.the priority of the defination by command line flags is higher than the priority of the defination in configuration file.
	if cliConf.IsTls {
		caValue := strings.ToLower(strings.TrimSpace(cliConf.Ca))
		certValue := strings.ToLower(strings.TrimSpace(cliConf.Cert))
		keyValue := strings.ToLower(strings.TrimSpace(cliConf.Key))
		insecureSkipVerifyValue := cliConf.InsecureSkipVerify

		_, _ = sysadmUtils.CheckFileExists(caValue, "")
		certExist, _ := sysadmUtils.CheckFileExists(certValue, "")
		keyExist, _ := sysadmUtils.CheckFileExists(keyValue, "")
		if certExist && keyExist {
			ret.InsecureSkipVerify = insecureSkipVerifyValue
			ret.Ca = caValue
			ret.Cert = certValue
			ret.Key = keyValue
			ret.IsTls = true
		}
	}

	return ret
}

// get environment name list by blockName . these list are defined in env_xxxx.go files
func getEnvDefineForBlock(blockName string) *map[string]string {

	if strings.TrimSpace(blockName) == "" {
		return nil
	}

	blockName = strings.ToLower(strings.TrimSpace(blockName))
	switch {
	case strings.Compare(blockName, "global") == 0:
		return &env_global
	case strings.Compare(blockName, "agent") == 0:
		return &env_agent
	case strings.Compare(blockName, "agentredis") == 0:
		return &env_agentRedis
	default:
		return nil
	}

}

/*
validateServerConf validate the server values in cliConf (set by command line flags), fileConf(set by configuration file) and envBlock (set by environment)
the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher
the priority of the defination in configuration file.
*/
func validateServerConf(cliConf config.Server, fileConf config.Server, envBlock string, isCheckSocket bool) (ret *config.Server, errs []sysadmerror.Sysadmerror) {
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081201, "debug", "try to validate server block settings in global block"))
	ret = &config.Server{
		Address: "",
		Port:    0,
		Socket:  "",
	}

	var envMap map[string]string
	envMapP := getEnvDefineForBlock(envBlock)
	if envMapP == nil {
		envMap = map[string]string{}
	} else {
		envMap = *envMapP
	}

	addressName, okAddress := envMap["Address"]
	portName, okPort := envMap["Port"]

	// got the variable name for server address and port
	if okAddress && okPort {
		ip, err := sysadmUtils.CheckIpAddress(strings.TrimSpace(os.Getenv(addressName)), false)
		errs = append(errs, err...)

		portStr := os.Getenv(portName)
		portInt, e := strconv.Atoi(portStr)

		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081202, "warning", "server port %s was found in environment with name %s,but it can not be convert to int.%s", portStr, portName, e))
		} else {
			port, err := sysadmUtils.CheckPort(portInt)
			errs = append(errs, err...)

			// if the values of address and port from environment are valid,then try to set them to running conf
			if ip != nil && port != 0 {
				address := strings.TrimSpace(os.Getenv(addressName))
				ret.Address = address
				ret.Port = port
				if isCheckSocket {
					socketName, okSocket := envMap["Socket"]
					if !okSocket {
						errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081203, "warning", "can not get environment variable name for %s block", envBlock))
					} else {
						socketFile := strings.TrimSpace(os.Getenv(socketName))
						if socketFile != "" {
							newSocketFile, e := sysadmUtils.CheckFileIsRead(socketFile, RunConf.WorkingDir)
							if e != nil {
								errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081204, "warning", "we got socket file path %s with environment name %s from environment variable, but this socket file is not valid %s", socketFile, socketName, e))
								ret.Socket = ""
							} else {
								ret.Socket = newSocketFile
							}
						} else {
							errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081205, "warning", "we can not got socket file path  with environment name %s from environment variables", socketName))
							ret.Socket = ""
						}
					}
				} else {
					ret.Socket = ""
				}
			}

		}
	}

	// 1. the priority of the defination in configuration file is higher than enverionments.
	// 2.try to validate the values set in configuration. the values  of ret will be relaced with this values if they are valid.
	address := fileConf.Address
	port := fileConf.Port
	tmpIp, eI := sysadmUtils.CheckIpAddress(address, false)
	tmpPort, eP := sysadmUtils.CheckPort(port)
	errs = append(errs, eI...)
	errs = append(errs, eP...)

	if tmpIp != nil && tmpPort != 0 {
		ret.Address = address
		ret.Port = port
		if isCheckSocket {
			socketFile := fileConf.Socket
			newSocketFile, e := sysadmUtils.CheckFileIsRead(socketFile, RunConf.WorkingDir)
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081206, "warning", "we got socket file path %s from configuration file, but this socket file is not valid %s", socketFile, e))
				ret.Socket = ""
			} else {
				ret.Socket = newSocketFile
			}
		}
	}

	// 1. the priority of the defination by command line flags is higher than the priority of the defination in configuration file.
	// 2.try to validate the values set by command flags. the values  of ret will be relaced with this values if they are valid.
	address = cliConf.Address
	port = cliConf.Port
	if strings.TrimSpace(address) != "" && port != 0 {
		tmpIp, eI = sysadmUtils.CheckIpAddress(address, false)
		tmpPort, eP = sysadmUtils.CheckPort(port)
		errs = append(errs, eI...)
		errs = append(errs, eP...)

		if tmpIp != nil && tmpPort != 0 {
			ret.Address = address
			ret.Port = port
			if isCheckSocket {
				socketFile := cliConf.Socket
				newSocketFile, e := sysadmUtils.CheckFileIsRead(socketFile, RunConf.WorkingDir)
				if e != nil {
					errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081207, "warning", "we got socket file path %s from command flags, but this socket file is not valid %s", socketFile, e))
					ret.Socket = ""
				} else {
					ret.Socket = newSocketFile
				}
			}
		}
	}

	return ret, errs
}

/*
validateLogConf validate the log values in cliConf (set by command line flags), fileConf(set by configuration file) and envBlock (set by environment)
the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher
the priority of the defination in configuration file.
*/
func validateLogConf(cliConf config.Log, fileConf config.Log, envBlock string) (ret *config.Log, errs []sysadmerror.Sysadmerror) {
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081301, "debug", "try to validate log block settings in global block"))
	ret = &config.Log{
		AccessLog:           "",
		AccessLogFp:         nil,
		ErrorLog:            "",
		ErrorLogFp:          nil,
		Kind:                DefaultLogKind,
		Level:               DefaultLogLevel,
		SplitAccessAndError: false,
		TimeStampFormat:     DefaultTimeStampFormat,
	}

	var envMap map[string]string
	envMapP := getEnvDefineForBlock(envBlock)
	if envMapP == nil {
		envMap = map[string]string{}
	} else {
		envMap = *envMapP
	}

	accessLogName, okAccessLog := envMap["AccessLog"]
	if !okAccessLog {
		accessLogName = ""
	}
	accessLog, err := validateLogFile(cliConf.AccessLog, fileConf.AccessLog, accessLogName)
	errs = append(errs, err...)
	if accessLog != "" {
		ret.AccessLog = accessLog
	} else {
		ret.AccessLog = DefaultLogFile
	}

	errorLogName, okErrorLog := envMap["ErrorLog"]
	if !okErrorLog {
		errorLogName = ""
	}
	errorLog, err := validateLogFile(cliConf.ErrorLog, fileConf.ErrorLog, errorLogName)
	errs = append(errs, err...)
	if errorLog != "" {
		ret.ErrorLog = errorLog
	}

	if strings.Compare(ret.AccessLog, ret.ErrorLog) != 0 && strings.TrimSpace(ret.ErrorLog) != "" {
		ret.SplitAccessAndError = true
	} else {
		ret.ErrorLog = ""
	}

	// check log kind
	kindName, okKind := envMap["Kind"]
	if okKind {
		kindValue := os.Getenv(kindName)
		kindValue = strings.ToLower(strings.TrimSpace(kindValue))
		if kindValue != "" && (kindValue == "text" || kindValue == "json") {
			ret.Kind = kindValue
		}
	}

	kindValue := strings.ToLower(strings.TrimSpace(fileConf.Kind))
	if kindValue != "" && (kindValue == "text" || kindValue == "json") {
		ret.Kind = kindValue
	}

	kindValue = strings.ToLower(strings.TrimSpace(cliConf.Kind))
	if kindValue != "" && (kindValue == "text" || kindValue == "json") {
		ret.Kind = kindValue
	}

	// check log level
	levelName, okLevel := envMap["Level"]
	if okLevel {
		levelValue := os.Getenv(levelName)
		levelValue = strings.ToLower(strings.TrimSpace(levelValue))
		if levelValue != "" && checkLogLevel(levelValue) {
			ret.Level = levelValue
		}
	}

	fileLevel := strings.TrimSpace(strings.ToLower(fileConf.Level))
	if checkLogLevel(fileLevel) {
		ret.Level = fileLevel
	}

	cliLevel := strings.TrimSpace(strings.ToLower(cliConf.Level))
	if checkLogLevel(cliLevel) {
		ret.Level = cliLevel
	}

	// check log time format
	// timeStampFormat can set by configuration file or environment. command flag can not set timeStampFormat
	timeFormatName, okTimeFormat := envMap["TimeFormat"]
	if okTimeFormat {
		timeFormatValue := os.Getenv(timeFormatName)
		timeFormatValue = strings.TrimSpace(timeFormatValue)
		if checkLogTimeFormat(timeFormatValue) {
			ret.TimeStampFormat = timeFormatValue
		}
	}

	fileFormat := strings.TrimSpace(fileConf.TimeStampFormat)
	if checkLogTimeFormat(fileFormat) {
		ret.TimeStampFormat = fileFormat
	}

	return ret, errs
}

/*
validateLogFile validate the logfile in cliConf (set by command line flags), fileConf(set by configuration file) and envBlock (set by environment)
the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher
the priority of the defination in configuration file.
*/
func validateLogFile(cliLogfile string, fileLogfile string, envName string) (ret string, errs []sysadmerror.Sysadmerror) {
	ret = ""

	// try to get logfile from environment and check it
	if strings.TrimSpace(envName) != "" {
		envLogfile := os.Getenv(envName)
		if strings.TrimSpace(envLogfile) != "" {
			logFile, e := sysadmUtils.CheckFileRW(envLogfile, RunConf.WorkingDir, false)
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081301, "warning", "we got log file  %s from environment with name %s, but this logfile can not open %s", logFile, logFile, e))
			}

			if strings.TrimSpace(logFile) != "" {
				ret = logFile
			}
		}
	}

	// try to get logfile from configuration file and validate it. ret will be overrided by this value if it passed the validation.
	if strings.TrimSpace(fileLogfile) != "" {
		fileLogfile := strings.TrimSpace(fileLogfile)
		logFile, e := sysadmUtils.CheckFileRW(fileLogfile, RunConf.WorkingDir, false)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081302, "warning", "we got log file  %s from configuration file, but this logfile can not open %s", fileLogfile, e))
		}

		if strings.TrimSpace(logFile) != "" {
			ret = logFile
		}
	}

	// try to get logfile from command line  and validate it. ret will be overrided by this value if it passed the validation.
	if strings.TrimSpace(cliLogfile) != "" {
		cliLogfile := strings.TrimSpace(cliLogfile)
		logFile, e := sysadmUtils.CheckFileRW(cliLogfile, RunConf.WorkingDir, false)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081303, "warning", "we got log file  %s from command line, but this logfile can not open %s", cliLogfile, e))
		}

		if strings.TrimSpace(logFile) != "" {
			ret = logFile
		}
	}

	return ret, errs
}

/*
validateRunMode validate the run mode  in cliConf (set by command line flags), fileConf(set by configuration file) and envBlock (set by environment)
the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher
the priority of the defination in configuration file.
*/
func validateRunMode(cliConf bool, fileConf bool, envBlock string) (ret bool) {
	ret = false

	var envMap map[string]string
	envMapP := getEnvDefineForBlock(envBlock)
	if envMapP == nil {
		envMap = map[string]string{}
	} else {
		envMap = *envMapP
	}

	debugName, okDebugMode := envMap["DebugMode"]
	if okDebugMode {
		debugValue := os.Getenv(debugName)
		debugValue = strings.ToLower(strings.TrimSpace(debugValue))
		if debugValue == "y" || debugValue == "yes" || debugValue == "on" {
			ret = true
		}
	}

	if fileConf {
		ret = fileConf
	}

	if cliConf {
		ret = cliConf
	}

	return ret
}

/*
validateNodeIdentifer  validate node identifer  in cliConf (set by command line flags), fileConf(set by configuration file) and envBlock (set by environment)
the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher
the priority of the defination in configuration file.
*/
func validateNodeIdentifer(cliConf string, fileConf string, envBlock string) (ret string) {
	ret = ""

	var envMap map[string]string
	envMapP := getEnvDefineForBlock(envBlock)
	if envMapP == nil {
		envMap = map[string]string{}
	} else {
		envMap = *envMapP
	}

	identiferName, okIdentifer := envMap["NodeIdentifer"]
	if okIdentifer {
		identiferValue := strings.TrimSpace(os.Getenv(identiferName))
		if identiferValue != "" && len(identiferValue) <= 63 {
			ret = identiferValue
		}
	}

	fileConf = strings.TrimSpace(fileConf)
	if fileConf != "" && len(fileConf) <= 63 {
		ret = fileConf
	}

	cliConf = strings.TrimSpace(cliConf)
	if cliConf != "" && len(cliConf) <= 63 {
		ret = cliConf
	}

	if ret == "" {
		ret = DefaultNodeIdentifer
	}

	return ret
}

/*
validateUri  validate uri in cliConf (set by command line flags), fileConf(set by configuration file) and envBlock (set by environment)
the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher
the priority of the defination in configuration file.
*/
func validateUri(cliConf string, fileConf string, envBlock string) (ret string) {
	ret = ""

	var envMap map[string]string
	envMapP := getEnvDefineForBlock(envBlock)
	if envMapP == nil {
		envMap = map[string]string{}
	} else {
		envMap = *envMapP
	}

	uriName, okUri := envMap["Uri"]
	if okUri {
		uriValue := strings.TrimSpace(os.Getenv(uriName))
		if uriValue != "" && len(uriValue) <= 63 {
			ret = uriValue
		}
	}

	fileConf = strings.TrimSpace(fileConf)
	if fileConf != "" && len(fileConf) <= 63 {
		ret = fileConf
	}

	cliConf = strings.TrimSpace(cliConf)
	if cliConf != "" && len(cliConf) <= 63 {
		ret = cliConf
	}

	if ret == "" {
		ret = "/"
	}

	return ret
}

/*
validateListenUri  validate fieldName in cliConf (set by command line flags), fileConf(set by configuration file) and envBlock (set by environment)
the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher
the priority of the defination in configuration file.
*/
func validateListenUri(cliConf, fileConf, envBlock,fieldName string) (ret string) {
	ret = ""

	var envMap map[string]string
	envMapP := getEnvDefineForBlock(envBlock)
	if envMapP == nil {
		envMap = map[string]string{}
	} else {
		envMap = *envMapP
	}

	uriName, okUri := envMap[fieldName]
	if okUri {
		uriValue := strings.TrimSpace(os.Getenv(uriName))
		if uriValue != "" && len(uriValue) <= 63 {
			ret = uriValue
		}
	}

	fileConf = strings.TrimSpace(fileConf)
	if fileConf != "" && len(fileConf) <= 63 {
		ret = fileConf
	}

	cliConf = strings.TrimSpace(cliConf)
	if cliConf != "" && len(cliConf) <= 63 {
		ret = cliConf
	}

	if ret == "" {
		ret = "/"
	}

	return ret
}

/*
validateSourceIP  validate source IP in cliConf (set by command line flags), fileConf(set by configuration file) and envBlock (set by environment)
the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher
the priority of the defination in configuration file.
*/
func validateSourceIP(cliConf string, fileConf string, envBlock string) (ret string) {
	ret = ""

	var envMap map[string]string
	envMapP := getEnvDefineForBlock(envBlock)
	if envMapP == nil {
		envMap = map[string]string{}
	} else {
		envMap = *envMapP
	}

	sourceIPName, okSourceIP := envMap["SourceIP"]
	if okSourceIP {
		sourceValue := strings.TrimSpace(os.Getenv(sourceIPName))
		if sourceValue != "" {
			tmpIp, _ := sysadmUtils.CheckIpAddress(sourceValue, true)
			if tmpIp != nil {
				ret = sourceValue
			}
		}
	}

	fileConf = strings.TrimSpace(fileConf)
	if fileConf != "" {
		tmpIp, _ := sysadmUtils.CheckIpAddress(fileConf, true)
		if tmpIp != nil {
			ret = fileConf
		}
	}

	cliConf = strings.TrimSpace(cliConf)
	if cliConf != "" {
		tmpIp, _ := sysadmUtils.CheckIpAddress(cliConf, true)
		if tmpIp != nil {
			ret = cliConf
		}
	}

	return ret
}
