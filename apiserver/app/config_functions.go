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
	"context"
	"errors"
	"fmt"
	"github.com/wangyysde/sysadmServer"
	"io/fs"
	certutil "k8s.io/client-go/util/cert"
	netutils "k8s.io/utils/net"
	"net"
	"os"
	"path/filepath"
	"strings"
	sysadmPki "sysadm/apiserver/pki"
	"sysadm/sysadmLog"
	"time"

	"sysadm/config"
	"sysadm/db"
	sysadmDB "sysadm/db"
	"sysadm/redis"
	"sysadm/sysadmerror"
	"sysadm/utils"
)

func handlerConfig() (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020007, "debug", "try to handle configurations for apiServer"))
	ok, err := handleNotInConfFile()
	errs = append(errs, err...)
	if !ok {
		return false, errs
	}

	// read configuration content from configuration file
	var conf = &Conf{}
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
	var cfgFile = ""
	if confFile != "" {
		if filepath.IsAbs(confFile) {
			fp, err := os.Open(confFile)
			if err != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020004, "fatal", "can not open configuration file %s error %s", confFile, err))
				return false, errs
			}
			_ = fp.Close()
			cfgFile = confFile
		} else {
			configPath := filepath.Join(workingDir, confFile)
			fp, err := os.Open(configPath)
			if err != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020005, "fatal", "can not open configuration file %s error %s", configPath, err))
				return false, errs
			}
			_ = fp.Close()
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
		_ = fp.Close()
		cfgFile = configPath
	}

	runData.runConf.ConfFile = cfgFile

	return true, errs

}

// SetVersion set version data to runData instance
func SetVersion(version *config.Version) {
	if version == nil {
		return
	}

	version.Version = appVer
	version.Author = appAuthor

	runData.runConf.Version = *version
}

// GetVersion get version data from runData instance
func GetVersion() *config.Version {
	if runData.runConf.Version.Version != "" {
		return &runData.runConf.Version
	}

	return nil
}

// GetCfgFile return the configuration file path of the application from runData
func GetCfgFile() string {
	return strings.TrimSpace(runData.runConf.ConfFile)
}

// SetCfgFile set configuration file path what has got from CLI flag to runData
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

	runData.runConf.ConfGlobal.Debug = conf.ConfGlobal.Debug
	var newVirtualIP []string
	for _, ip := range conf.ConfGlobal.VirtualIP {
		validatedIp, e := utils.ValidateAddress(ip, false)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020030, "fatal", "virtual IP %s is not valid ", ip))
			return false, errs
		}
		newVirtualIP = append(newVirtualIP, validatedIp)
	}
	runData.runConf.ConfGlobal.VirtualIP = newVirtualIP

	var newExtraSans []string
	for _, san := range conf.ConfGlobal.ExtraSans {
		san = strings.TrimSpace(san)
		_, e := utils.ValidateAddress(san, false)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020031, "fatal", "extra san  %s is not an IP or hostname ", san))
			return false, errs
		}
		newExtraSans = append(newExtraSans, san)
	}
	runData.runConf.ConfGlobal.ExtraSans = newExtraSans

	return true, errs
}

// validate configurations read from configuration file, then pass them to runData if them are valid.
func validateServerBlock(conf *Conf) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	address := strings.TrimSpace(conf.ConfServer.Address)
	tmpAddress, err := config.ValidateListenAddress(address, apiserverAddress, "")
	errs = append(errs, err...)
	runData.runConf.ConfServer.Address = tmpAddress

	tmpPort, err := config.ValidateListenPort(conf.ConfServer.Port, apiserverPort, "")
	errs = append(errs, err...)
	runData.runConf.ConfServer.Port = tmpPort

	insecretPort, err := config.ValidateListenPort(conf.ConfServer.InsecretPort, apiserverInsecretPort, "")
	errs = append(errs, err...)
	runData.runConf.ConfServer.InsecretPort = insecretPort

	runData.runConf.ConfServer.Insecret = conf.ConfServer.Insecret
	runData.runConf.ConfServer.IsTls = conf.ConfServer.IsTls

	if runData.runConf.ConfServer.IsTls {
		e := prepareApiServerCerts()
		if e != nil {
			runData.runConf.ConfServer.IsTls = false
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020040, "warn", "an error(%s) has occurred when prepare certifications, we try start server without TLS", e))
		}
	} else {
		runData.runConf.ConfServer.Ca = ""
		runData.runConf.ConfServer.Cert = ""
		runData.runConf.ConfServer.Key = ""
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
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020014, "debug", "try to handle configuration items in redis block"))

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

	runData.runConf.ConfRedis.Username = strings.TrimSpace(conf.ConfRedis.Username)
	runData.runConf.ConfRedis.Password = strings.TrimSpace(conf.ConfRedis.Password)
	runData.runConf.ConfRedis.SentinelUsername = strings.TrimSpace(conf.ConfRedis.SentinelUsername)
	runData.runConf.ConfRedis.SentinelPassword = strings.TrimSpace(conf.ConfRedis.SentinelPassword)
	runData.runConf.ConfRedis.DB = conf.ConfRedis.DB

	runData.runConf.ConfRedis.Tls.IsTls = conf.ConfRedis.Tls.IsTls
	if runData.runConf.ConfRedis.Tls.IsTls {
		ca := strings.TrimSpace(conf.ConfRedis.Tls.Ca)
		cert := strings.TrimSpace(conf.ConfRedis.Tls.Cert)
		key := strings.TrimSpace(conf.ConfRedis.Tls.Key)
		tmpCa, err := config.ValidateTlsFile(ca, "", "", runData.workingRoot)
		errs = append(errs, err...)
		tmpCert, err := config.ValidateTlsFile(cert, "", "", runData.workingRoot)
		errs = append(errs, err...)
		tmpKey, err := config.ValidateTlsFile(key, "", "", runData.workingRoot)
		errs = append(errs, err...)
		if tmpCa == "" || tmpCert == "" || tmpKey == "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020018, "warning", "isTls has be set to true, but one(s) of certifications file is not readable. apiserver will connect to redis server withoud TLS"))
			runData.runConf.ConfRedis.Tls.IsTls = false
			runData.runConf.ConfRedis.Tls.Ca = ""
			runData.runConf.ConfRedis.Tls.Cert = ""
			runData.runConf.ConfRedis.Tls.Key = ""
		} else {
			runData.runConf.ConfRedis.Tls.Ca = tmpCa
			runData.runConf.ConfRedis.Tls.Cert = tmpCert
			runData.runConf.ConfRedis.Tls.Key = tmpKey
		}
	} else {
		runData.runConf.ConfRedis.Tls.Ca = ""
		runData.runConf.ConfRedis.Tls.Cert = ""
		runData.runConf.ConfRedis.Tls.Key = ""
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

// SetLogger set parameters to accessLogger and errorLoger
func setLogger() (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	logger := sysadmLog.NewSysadmLogger()
	_ = sysadmServer.SetLoggerKind(runData.runConf.ConfLog.Kind)
	logger.SetLoggerKind(runData.runConf.ConfLog.Kind)
	_ = sysadmServer.SetLogLevel(runData.runConf.ConfLog.Level)
	logger.SetLoggerLevel(runData.runConf.ConfLog.Level)
	_ = sysadmServer.SetTimestampFormat(runData.runConf.ConfLog.TimeStampFormat)
	logger.SetTimestampFormat(runData.runConf.ConfLog.TimeStampFormat)

	if runData.runConf.ConfLog.AccessLog != "" {
		_, fp, err := sysadmServer.SetAccessLogFile(runData.runConf.ConfLog.AccessLog)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020026, "error", "access log %s can not be openned. error %. access logs will be output to standard device", runData.runConf.ConfLog.AccessLog, err))
			e := sysadmServer.SetAccessLoggerWithFp(os.Stdout)
			_ = logger.SetAccessLoggerWithFp(os.Stdout)
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020027, "error", "can not set access logger error %s", e))
				return false, errs
			}
		} else {
			runData.runConf.ConfLog.AccessLogFp = fp
			_ = logger.SetAccessLoggerWithFp(fp)
		}
	}

	if runData.runConf.ConfLog.SplitAccessAndError && runData.runConf.ConfLog.ErrorLog != "" {
		_, fp, err := sysadmServer.SetErrorLogFile(runData.runConf.ConfLog.ErrorLog)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020028, "error", "can not open error log file(%s) error: %s", runData.runConf.ConfLog.ErrorLog, err))
		} else {
			runData.runConf.ConfLog.ErrorLogFp = fp
			logger.SetErrorLoggerWithFp(fp)
		}
	}

	sysadmServer.SetIsSplitLog(runData.runConf.ConfLog.SplitAccessAndError)
	logger.SetIsSplitLog(runData.runConf.ConfLog.SplitAccessAndError)
	if runData.runConf.ConfGlobal.Debug {
		sysadmServer.SetMode(sysadmServer.DebugMode)
	}

	runData.logEntity = logger
	return true, errs
}

// closeLogger close access log file descriptor and error log file descriptor
// set AccessLogger  and ErrorLogger to nil
func closeLogger() {
	if runData.runConf.ConfLog.AccessLogFp != nil {
		fp := runData.runConf.ConfLog.AccessLogFp
		_ = fp.Close()
		runData.runConf.ConfLog.AccessLogFp = nil
	}

	if runData.runConf.ConfLog.ErrorLogFp != nil {
		fp := runData.runConf.ConfLog.ErrorLogFp
		_ = fp.Close()
		runData.runConf.ConfLog.ErrorLogFp = nil
	}
}

// InitRedis initate a new redis entity
func initRedis() (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	entity, e := redis.NewClient(runData.runConf.ConfRedis, runData.workingRoot)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20020029, "fatal", "can not open connection to redis server %s", e))
		return false, errs
	}

	runData.redisEntity = entity
	var ctx = context.Background()
	runData.redisCtx = ctx

	return true, errs
}

// close the entity of redis
func closeRedisEntity() {

	if runData.redisEntity != nil {
		entity := runData.redisEntity
		_ = entity.Close()
	}

	runData.redisEntity = nil
	runData.redisCtx = nil
}

// InitDB initate a new DB entity
func initDBEntity() (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	definedConf := runData.runConf.ConfDB
	var sslModle string = ""
	if !definedConf.IsTls {
		sslModle = "disable"
	} else {
		sslModle = "enable"
	}

	dbConf := sysadmDB.DbConfig{
		Type:         definedConf.Type,
		Host:         definedConf.Address,
		Port:         definedConf.Port,
		User:         definedConf.UserName,
		Password:     definedConf.Password,
		DbName:       definedConf.DBName,
		SslMode:      sslModle,
		SslCa:        definedConf.Ca,
		SslCert:      definedConf.Cert,
		SslKey:       definedConf.Key,
		MaxOpenConns: definedConf.MaxOpenConns,
		MaxIdleConns: definedConf.MaxIdleConns,
		Connect:      nil,
		Entity:       nil,
	}

	newDBConf, err := sysadmDB.InitDbConfig(&dbConf, runData.workingRoot)
	errs = append(errs, err...)
	maxLevel := sysadmerror.GetMaxLevel(errs)
	fatalLevel := sysadmerror.GetLevelNum("fatal")
	if maxLevel >= fatalLevel {
		return false, errs
	}

	entity := newDBConf.Entity
	err = entity.OpenDbConnect()
	errs = append(errs, err...)
	maxLevel = sysadmerror.GetMaxLevel(errs)
	if maxLevel >= fatalLevel {
		return false, errs
	}

	runData.dbEntity = entity
	return true, errs
}

// close the entity of DB
func closeDBEntity() {
	dbEntity := runData.dbEntity
	if dbEntity != nil {
		dbEntity.CloseDB()
	}

	dbEntity = nil
}

func prepareApiServerCerts() error {
	var caContent, caKeyContent, certContent, certKeyContent []byte
	createCa := false
	createCert := false

	// try to read ca from disk and check them validity.
	// regenerate the ca and caKey if them is not exist or NOT validate
	workingRoot := runData.workingRoot
	ca := strings.TrimSpace(runData.runConf.ConfServer.Ca)
	if ca == "" {
		ca = filepath.Join(pkiPath, caFile)
	}
	if !filepath.IsAbs(ca) {
		ca = filepath.Join(workingRoot, ca)
	}
	caContent, e := os.ReadFile(ca)
	if e != nil {
		switch {
		case errors.Is(e, fs.ErrInvalid):
			return fmt.Errorf("maybe ca path %s is not valid", ca)
		case errors.Is(e, fs.ErrPermission):
			return fmt.Errorf("no permission to read ca file %s", ca)
		case errors.Is(e, fs.ErrNotExist):
			createCa = true
		case errors.Is(e, os.ErrDeadlineExceeded):
			return fmt.Errorf("read ca file %s timeout", ca)
		}
	}

	if !createCa {
		caCerts, e := certutil.ParseCertsPEM(caContent)
		if e != nil {
			createCa = true
		} else {
			caCert := caCerts[0]
			e := sysadmPki.ValidateCertPeriod(caCert, time.Duration(0))
			if e != nil {
				createCa = true
			}
		}
	}

	if createCa {
		_, _, newCaPem, newCaKeyPem, e := sysadmPki.CreateCertificateAuthority(caCommonName, caOrgnaization, caPeriodDays, publicKeyAlgorithm)
		if e != nil {
			return e
		}
		caContent = newCaPem
		caKeyContent = newCaKeyPem
		createCert = true
	}

	// try to read certification and key from  the disk and check them validity
	// set createCert to true if certification(or key) is not exist or one of them is NOT validate
	cert := strings.TrimSpace(runData.runConf.ConfServer.Cert)
	certKey := strings.TrimSpace(runData.runConf.ConfServer.Key)
	if cert == "" {
		cert = filepath.Join(pkiPath, apiServerCertFile)
	}
	if certKey == "" {
		certKey = filepath.Join(pkiPath, apiServerCertKeyFile)
	}
	if !filepath.IsAbs(cert) {
		cert = filepath.Join(workingRoot, cert)
	}
	if !filepath.IsAbs(certKey) {
		certKey = filepath.Join(workingRoot, certKey)
	}

	certContent, e = os.ReadFile(cert)
	if e != nil {
		switch {
		case errors.Is(e, fs.ErrInvalid):
			return fmt.Errorf("maybe certification path %s is not valid", cert)
		case errors.Is(e, fs.ErrPermission):
			return fmt.Errorf("no permission to read certification file %s", cert)
		case errors.Is(e, fs.ErrNotExist):
			createCert = true
		case errors.Is(e, os.ErrDeadlineExceeded):
			return fmt.Errorf("read certification file %s timeout", cert)
		}
	}
	if !createCert {
		certs, e := certutil.ParseCertsPEM(certContent)
		if e != nil {
			createCert = true
		} else {
			e := sysadmPki.ValidateCertPeriod(certs[0], time.Duration(0))
			if e != nil {
				createCert = true
			}
		}
	}

	if !createCert {
		certKeyContent, e = os.ReadFile(certKey)
		if e != nil {
			switch {
			case errors.Is(e, fs.ErrInvalid):
				return fmt.Errorf("maybe path of certification key %s is not valid", certKey)
			case errors.Is(e, fs.ErrPermission):
				return fmt.Errorf("no permission to read certification key file %s", certKey)
			case errors.Is(e, fs.ErrNotExist):
				createCert = true
			case errors.Is(e, os.ErrDeadlineExceeded):
				return fmt.Errorf("read certification key file %s timeout", certKey)
			}
		}
		_, e = sysadmPki.ParseKeyPEM(string(certKeyContent))
		if e != nil {
			createCert = true
		}
	}

	if createCert {
		if !createCa {
			caKey := filepath.Join(workingRoot, pkiPath, caKeyFile)
			caKeyContent, e = os.ReadFile(caKey)
			if e != nil {
				switch {
				case errors.Is(e, fs.ErrInvalid):
					return fmt.Errorf("maybe ca key path %s is not valid", caKey)
				case errors.Is(e, fs.ErrPermission):
					return fmt.Errorf("no permission to read ca key file %s", caKey)
				case errors.Is(e, fs.ErrNotExist):
					_, _, newCaPem, newCaKeyPem, e := sysadmPki.CreateCertificateAuthority(caCommonName, caOrgnaization, caPeriodDays, publicKeyAlgorithm)
					if e != nil {
						return e
					}
					caContent = newCaPem
					caKeyContent = newCaKeyPem
					createCa = true
				case errors.Is(e, os.ErrDeadlineExceeded):
					return fmt.Errorf("read ca key  %s from disk  timeout", caKey)
				}
			}
		}

		localIPs, e := utils.GetLocalIPs()
		if e != nil {
			return e
		}
		var altNameIPs []net.IP
		for _, ip := range localIPs {
			netIP := netutils.ParseIPSloppy(ip)
			altNameIPs = append(altNameIPs, netIP)
		}
		for _, ip := range runData.runConf.ConfGlobal.VirtualIP {
			netIP := netutils.ParseIPSloppy(ip)
			altNameIPs = append(altNameIPs, netIP)
		}

		hostName, e := os.Hostname()
		if e != nil {
			return fmt.Errorf("an error has occurred when get hostname: %s", e)
		}
		dnsNames := []string{"localhost", hostName}
		for _, v := range runData.runConf.ConfGlobal.ExtraSans {
			dnsNames = append(dnsNames, v)
		}

		_, _, certPem, keyPem, e := sysadmPki.CreateCertAndKey(publicKeyAlgorithm, string(caContent), string(caKeyContent),
			altNameIPs, apiServerCertPeriodDays, apiServerCertCommonName, apiServerCertOrgnaization, dnsNames, sysadmPki.CertUsageTypeServerAuth)
		if e != nil {
			return fmt.Errorf("an error has occurred when generate certification and key for apiServer: %s", e)
		}
		certContent = certPem
		certKeyContent = keyPem
	}

	if createCa {
		e = os.WriteFile(ca, caContent, 0644)
		if e != nil {
			return fmt.Errorf("an error has occurred when write ca certification to disk: %s", e)
		}
		caKey := filepath.Join(workingRoot, pkiPath, caKeyFile)
		e = os.WriteFile(caKey, caKeyContent, 0600)
		if e != nil {
			return fmt.Errorf("an error has occurred when write key of ca to disk: %s", e)
		}
	}

	if createCert {
		e = os.WriteFile(cert, certContent, 0644)
		if e != nil {
			return fmt.Errorf("an error has occurred when write ca certification to disk: %s", e)
		}
		e = os.WriteFile(certKey, certKeyContent, 0600)
		if e != nil {
			return fmt.Errorf("an error has occurred when write key of certification to disk: %s", e)
		}
	}

	fullCertFile := filepath.Join(workingRoot, pkiPath, apiServerFullCertFile)
	if !utils.IsFileReadable(fullCertFile) || createCert {
		fullCertContent := string(certContent) + "\n" + string(caContent)
		e = os.WriteFile(fullCertFile, []byte(fullCertContent), 0644)
		if e != nil {
			return fmt.Errorf("an error has occurred when write full certification to disk: %s", e)
		}
	}

	return nil
}
