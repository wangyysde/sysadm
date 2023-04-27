/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
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
*
* NOTE:
* defined some functions are related to handle configurations.
 */

package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/wangyysde/sysadm/config"
	"github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/redis"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
)

func HandlerConfig() (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020007, "debug", "try to handle configurations for apiserver"))
	ok, err := handleNotInConfFile()
	errs = append(errs, err...)
	if !ok {
		return false, errs
	}

	// read configuration content from configuration file
	var conf *Conf = &Conf{}
	tmpConf, err := config.GetCfgContent(runData.runConf.ConfFile, conf)
	errs = append(errs, err...)
	if tmpConf == nil {
		return false, errs
	}

	// validate configurations defined in global block
	ok, err = validateGlobalBlock(conf)
	errs = append(errs, err...)
	if !ok {
		return false, errs
	}

	// validate configurations defined in server block
	ok, err = validateServerBlock(conf)
	errs = append(errs, err...)
	if !ok {
		return false, errs
	}

	// validate configurations defined in log block
	ok, err = validateLogBlock(conf)
	errs = append(errs, err...)
	if !ok {
		return false, errs
	}

	// validata configurations defined in redis block
	ok, err = validateRedisBlock(conf)
	errs = append(errs, err...)
	if !ok {
		return false, errs
	}

	// validata configurations defined in db block
	ok, err = validateDbBlock(conf)
	errs = append(errs, err...)
	if !ok {
		return false, errs
	}

	return true, errs
}

// HandleNotInConfFile handler the configuration items which can not define in define file,such as working dir, configuration file path.
func handleNotInConfFile() (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	confFile := runData.runConf.ConfFile
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020001, "debug", "try to get working dir"))
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020002, "fatal", "get working dir error %s", err))
		return false, errs
	}
	workingDir := filepath.Join(binPath, "../")
	runData.workingRoot = filepath.Join(binPath, "../")

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020003, "debug", "checking configuration file path"))
	var cfgFile string = ""
	if confFile != "" {
		if filepath.IsAbs(confFile) {
			fp, err := os.Open(confFile)
			if err != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020004, "fatal", "can not open configuration file %s error %s", confFile, err))
				return false, errs
			}
			fp.Close()
			cfgFile = confFile
		} else {
			configPath := filepath.Join(workingDir, confFile)
			fp, err := os.Open(configPath)
			if err != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020005, "fatal", "can not open configuration file %s error %s", configPath, err))
				return false, errs
			}
			fp.Close()
			cfgFile = configPath
		}
	} else {
		// try to get configuration file from default path
		configPath := filepath.Join(workingDir, confFilePath)
		fp, err := os.Open(configPath)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020006, "fatal", "can not open configuration file %s error %s", configPath, err))
			return false, errs
		}
		fp.Close()
		cfgFile = configPath
	}

	runData.runConf.ConfFile = cfgFile

	return true, errs

}

// set version data to runData instance
func SetVersion(version *config.Version) {
	if version == nil {
		return
	}

	version.Version = appVer
	version.Author = appAuthor

	runData.runConf.Version = *version
}

// get version data from runData instance
func GetVersion() *config.Version {
	if runData.runConf.Version.Version != "" {
		return &runData.runConf.Version
	}

	return nil
}

// return the configuration file path of the application from runData
func GetCfgFile() string {
	return strings.TrimSpace(runData.runConf.ConfFile)
}

// set configuration file path what has got from CLI flag to runData
func SetCfgFile(cfgFile string) {
	cfgFile = strings.TrimSpace(cfgFile)
	if cfgFile == "" {
		cfgFile = confFilePath
	}

	runData.runConf.ConfFile = cfgFile
}

// validate configurations read from configuration file, then pass them to runData if them are valid.
func validateGlobalBlock(conf *Conf) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020008, "debug", "try to handle configuration items in global block"))

	runData.runConf.ConfGlobal.Passive = conf.ConfGlobal.Passive
	runData.runConf.ConfGlobal.Debug = conf.ConfGlobal.Debug

	commandUri := strings.TrimSpace(conf.ConfGlobal.CommandUri)
	if commandUri == "" {
		if conf.ConfGlobal.Passive {
			commandUri = passiveResponseCommandUri
		} else {
			commandUri = activeSendCommandUri
		}
	}
	runData.runConf.ConfGlobal.CommandUri = commandUri

	commandStatusUri := strings.TrimSpace(conf.ConfGlobal.CommandStatusUri)
	if commandStatusUri == "" {
		if conf.ConfGlobal.Passive {
			commandStatusUri = passiveResponseCommandStatusUri
		} else {
			commandStatusUri = activeSendCommandStatusUri
		}
	}
	runData.runConf.ConfGlobal.CommandStatusUri = commandStatusUri

	commandLogsUri := strings.TrimSpace(conf.ConfGlobal.CommandLogsUri)
	if commandLogsUri == "" {
		if conf.ConfGlobal.Passive {
			commandLogsUri = passiveResponseCommandLogsUri
		} else {
			commandLogsUri = activeSendCommandLogsUri
		}
	}
	runData.runConf.ConfGlobal.CommandLogsUri = commandLogsUri

	runData.runConf.ConfGlobal.CheckCommandInterval = conf.ConfGlobal.CheckCommandInterval
	if conf.ConfGlobal.CheckCommandInterval == 0 {
		runData.runConf.ConfGlobal.CheckCommandInterval = checkCommandInterval
	}

	runData.runConf.ConfGlobal.GetStatusInterval = conf.ConfGlobal.GetStatusInterval
	if conf.ConfGlobal.GetStatusInterval == 0 {
		runData.runConf.ConfGlobal.GetStatusInterval = getStatusInterval
	}

	runData.runConf.ConfGlobal.GetLogInterval = conf.ConfGlobal.GetLogInterval
	if conf.ConfGlobal.GetLogInterval == 0 {
		runData.runConf.ConfGlobal.GetLogInterval = getLogInterval
	}

	return true, errs
}

// validate configurations read from configuration file, then pass them to runData if them are valid.
func validateServerBlock(conf *Conf) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020009, "debug", "try to handle configuration items in server block"))

	address := strings.TrimSpace(conf.ConfServer.Address)
	tmpAddress, err := config.ValidateListenAddress(address, apiserverAddress, "")
	errs = append(errs, err...)
	runData.runConf.ConfServer.Address = tmpAddress

	tmpPort, err := config.ValidateListenPort(conf.ConfServer.Port, apiserverPort, "")
	errs = append(errs, err...)
	runData.runConf.ConfServer.Port = tmpPort

	socketFile := strings.TrimSpace(conf.ConfServer.Socket)
	tmpSocket, err := config.ValidateListenSocket(socketFile, "", "", runData.workingRoot)
	errs = append(errs, err...)
	runData.runConf.ConfServer.Socket = tmpSocket

	runData.runConf.ConfServer.Insecret = conf.ConfServer.Insecret
	runData.runConf.ConfServer.IsTls = conf.ConfServer.IsTls
	if runData.runConf.ConfServer.IsTls {
		ca := strings.TrimSpace(conf.ConfServer.Ca)
		cert := strings.TrimSpace(conf.ConfServer.Cert)
		key := strings.TrimSpace(conf.ConfServer.Key)
		tmpCa, err := config.ValidateTlsFile(ca, "", "", runData.workingRoot)
		errs = append(errs, err...)
		tmpCert, err := config.ValidateTlsFile(cert, "", "", runData.workingRoot)
		errs = append(errs, err...)
		tmpKey, err := config.ValidateTlsFile(key, "", "", runData.workingRoot)
		errs = append(errs, err...)
		if tmpCa == "" || tmpCert == "" || tmpKey == "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020010, "warning", "isTls has be set to true, but one(s) of certifications file is not readable. apiserver will be started with insecret"))
			runData.runConf.ConfServer.IsTls = false
			runData.runConf.ConfServer.InsecretPort = 0
			runData.runConf.ConfServer.Ca = ""
			runData.runConf.ConfServer.Cert = ""
			runData.runConf.ConfServer.Key = ""
		} else {
			runData.runConf.ConfServer.Ca = tmpCa
			runData.runConf.ConfServer.Cert = tmpCert
			runData.runConf.ConfServer.Key = tmpKey
		}
	} else {
		runData.runConf.ConfServer.InsecretPort = 0
		runData.runConf.ConfServer.Ca = ""
		runData.runConf.ConfServer.Cert = ""
		runData.runConf.ConfServer.Key = ""
	}

	insecretPort := conf.ConfServer.InsecretPort
	if insecretPort == 0 {
		insecretPort = apiserverInsecretPort
	}

	if runData.runConf.ConfServer.IsTls {
		runData.runConf.ConfServer.InsecretPort = insecretPort
	} else {
		runData.runConf.ConfServer.InsecretPort = 0
	}

	return true, errs
}

// validate configurations read from configuration file in log block, then pass them to runData if them are valid.
func validateLogBlock(conf *Conf) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020011, "debug", "try to handle configuration items in log block"))

	accessLog := strings.TrimSpace(conf.ConfLog.AccessLog)
	errorLog := strings.TrimSpace(conf.ConfLog.ErrorLog)
	tmpAccessLog, err := config.ValidateLogFile(accessLog, accessLogFile, "", runData.workingRoot)
	errs = append(errs, err...)
	if tmpAccessLog == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020012, "warning", "access log file can not be writeable. all access logs will be write to stdout"))
	}
	runData.runConf.ConfLog.AccessLog = tmpAccessLog

	tmpErrorLog, err := config.ValidateLogFile(errorLog, errorLogFile, "", runData.workingRoot)
	errs = append(errs, err...)
	if tmpErrorLog == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020013, "warning", "error log file can not be writeable. all error logs will be write to access log file"))
	}
	runData.runConf.ConfLog.ErrorLog = tmpErrorLog

	kind := strings.TrimSpace(conf.ConfLog.Kind)
	tmpKind, err := config.ValidateLogKind(kind, defaultLogKind, "")
	errs = append(errs, err...)
	runData.runConf.ConfLog.Kind = tmpKind

	level := strings.TrimSpace(conf.ConfLog.Level)
	tmpLevel, err := config.ValidateLogLevel(level, defaultLogLevel, "")
	errs = append(errs, err...)
	runData.runConf.ConfLog.Level = tmpLevel

	if runData.runConf.ConfLog.AccessLog != "" && runData.runConf.ConfLog.ErrorLog != "" {
		runData.runConf.ConfLog.SplitAccessAndError = true
	} else {
		runData.runConf.ConfLog.SplitAccessAndError = false
	}

	timeFormat := strings.TrimSpace(conf.ConfLog.TimeStampFormat)
	tmpTimeFormat, err := config.ValidateLogTimeFormat(timeFormat, defaultTimeStampFormat, "")
	errs = append(errs, err...)
	runData.runConf.ConfLog.TimeStampFormat = tmpTimeFormat

	return true, errs
}

// validate configurations read from configuration file in redis block, then pass them to runData if them are valid.
func validateRedisBlock(conf *Conf) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020014, "debug", "try to handle configuration items in log block"))

	if conf.ConfRedis.Mode < 1 || conf.ConfRedis.Mode > 3 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020015, "warning", "redis mode(%d) is not valid. single mode will be used"))
		runData.runConf.ConfRedis.Mode = redis.RedisModeSingle
	} else {
		runData.runConf.ConfRedis.Mode = conf.ConfRedis.Mode
	}

	redisMaster := strings.TrimSpace(conf.ConfRedis.Master)
	if !redis.IsValidMaster(runData.runConf.ConfRedis.Mode, redisMaster) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020016, "error", "redis master is not valid"))
		return false, errs
	}
	runData.runConf.ConfRedis.Master = redisMaster

	redisAddrs := strings.TrimSpace(conf.ConfRedis.Addrs)
	if !redis.IsValidAddrs(runData.runConf.ConfRedis.Mode, redisAddrs) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020017, "error", "redis server address is not valid"))
		return false, errs
	}
	runData.runConf.ConfRedis.Addrs = redisAddrs

	runData.runConf.ConfRedis.UserName = strings.TrimSpace(conf.ConfRedis.UserName)
	runData.runConf.ConfRedis.Password = strings.TrimSpace(conf.ConfRedis.Password)
	runData.runConf.ConfRedis.SentinelUsername = strings.TrimSpace(conf.ConfRedis.SentinelUsername)
	runData.runConf.ConfRedis.SentinelPassword = strings.TrimSpace(conf.ConfRedis.SentinelPassword)
	runData.runConf.ConfRedis.DB = conf.ConfRedis.DB

	runData.runConf.ConfRedis.IsTls = conf.ConfRedis.IsTls
	if runData.runConf.ConfRedis.IsTls {
		ca := strings.TrimSpace(conf.ConfRedis.Ca)
		cert := strings.TrimSpace(conf.ConfRedis.Cert)
		key := strings.TrimSpace(conf.ConfRedis.Key)
		tmpCa, err := config.ValidateTlsFile(ca, "", "", runData.workingRoot)
		errs = append(errs, err...)
		tmpCert, err := config.ValidateTlsFile(cert, "", "", runData.workingRoot)
		errs = append(errs, err...)
		tmpKey, err := config.ValidateTlsFile(key, "", "", runData.workingRoot)
		errs = append(errs, err...)
		if tmpCa == "" || tmpCert == "" || tmpKey == "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020018, "warning", "isTls has be set to true, but one(s) of certifications file is not readable. apiserver will connect to redis server withoud TLS"))
			runData.runConf.ConfRedis.IsTls = false
			runData.runConf.ConfRedis.Ca = ""
			runData.runConf.ConfRedis.Cert = ""
			runData.runConf.ConfRedis.Key = ""
		} else {
			runData.runConf.ConfRedis.Ca = tmpCa
			runData.runConf.ConfRedis.Cert = tmpCert
			runData.runConf.ConfRedis.Key = tmpKey
		}
	} else {
		runData.runConf.ConfRedis.Ca = ""
		runData.runConf.ConfRedis.Cert = ""
		runData.runConf.ConfRedis.Key = ""
	}

	return true, errs

}

// validate configurations read from configuration file in db block, then pass them to runData if them are valid.
func validateDbBlock(conf *Conf) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020019, "debug", "try to handle configuration items in db block"))

	dbType := strings.TrimSpace(strings.ToLower(conf.ConfDB.Type))
	if !db.IsSupportedDB(dbType) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020020, "error", "db type %s can not be supported", dbType))
		return false, errs
	}
	runData.runConf.ConfDB.Type = dbType

	dbName := strings.TrimSpace(conf.ConfDB.DBName)
	if !db.CheckIdentifier(dbType, dbName) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020021, "error", "DB name %s for %s server is not valid", dbName, dbType))
		return false, errs
	}
	runData.runConf.ConfDB.DBName = dbName

	dbAddress := strings.TrimSpace(conf.ConfDB.Address)
	ip, err := utils.CheckIpAddress(dbAddress, false)
	errs = append(errs, err...)
	if ip == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020022, "error", "address %s is not valid DB server address", dbAddress))
		return false, errs
	}
	runData.runConf.ConfDB.Address = dbAddress

	dbPort := conf.ConfDB.Port
	dbPort, err = utils.CheckPort(dbPort)
	if dbPort == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020023, "error", "DB port %s is not valid", dbPort))
		return false, errs
	}
	runData.runConf.ConfDB.Port = dbPort

	runData.runConf.ConfDB.Socket = ""
	runData.runConf.ConfDB.IsTls = conf.ConfDB.IsTls
	if runData.runConf.ConfDB.IsTls {
		ca := strings.TrimSpace(conf.ConfDB.Ca)
		cert := strings.TrimSpace(conf.ConfDB.Cert)
		key := strings.TrimSpace(conf.ConfDB.Key)
		tmpCa, err := config.ValidateTlsFile(ca, "", "", runData.workingRoot)
		errs = append(errs, err...)
		tmpCert, err := config.ValidateTlsFile(cert, "", "", runData.workingRoot)
		errs = append(errs, err...)
		tmpKey, err := config.ValidateTlsFile(key, "", "", runData.workingRoot)
		errs = append(errs, err...)
		if tmpCa == "" || tmpCert == "" || tmpKey == "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020024, "warning", "isTls has be set to true, but one(s) of certifications file is not readable. apiserver will connect to db server withoud TLS"))
			runData.runConf.ConfDB.IsTls = false
			runData.runConf.ConfDB.Ca = ""
			runData.runConf.ConfDB.Cert = ""
			runData.runConf.ConfDB.Key = ""
		} else {
			runData.runConf.ConfDB.Ca = tmpCa
			runData.runConf.ConfDB.Cert = tmpCert
			runData.runConf.ConfDB.Key = tmpKey
		}
	} else {
		runData.runConf.ConfDB.Ca = ""
		runData.runConf.ConfDB.Cert = ""
		runData.runConf.ConfDB.Key = ""
	}

	dbUser := strings.TrimSpace(conf.ConfDB.UserName)
	dbPasswd := strings.TrimSpace(conf.ConfDB.Password)
	if dbUser == "" || dbPasswd == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020025, "error", "db user and db password must be not empty"))
		return false, errs
	}
	runData.runConf.ConfDB.UserName = dbUser
	runData.runConf.ConfDB.Password = dbPasswd

	maxOpenConns := conf.ConfDB.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = defaultMaxDBOpenConns
	}
	runData.runConf.ConfDB.MaxOpenConns = maxOpenConns

	maxIdleConns := conf.ConfDB.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = defaultMaxDBIdleConns
	}
	runData.runConf.ConfDB.MaxIdleConns = maxIdleConns

	return true, errs
}
