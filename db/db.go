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

package db

import (
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/wangyysde/sysadm/sysadmerror"
)

var SupportDBs = []string{
	"postgre",
	"mysql",
	/* =============================
		TODO:
	    "info",
	    "warning",
	    "error",
	    "fatal",
	    "panic",
	*/
}

/*
  Checking database parametes and initating an instance
  return DbConfig and sysadmerror.Sysadmerror.
  the error level will be set to warning if any of certification files is not exist and  slmode is not disable,
  then set sslmode to disable.
  cmdRunPath: the path of executeable file
*/
func InitDbConfig(config *DbConfig, cmdRunPath string) (*DbConfig, []sysadmerror.Sysadmerror) {
	// Checking the type of DB
	// Initating an entity if type is valid otherwise return fatal error
	// TODO we shoud add other db support 

	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(100001, "debug", "Now checking database configurations."))
	if !IsSupportedDB(config.Type) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100002, "fatal", "DB type %s not be suupported.", config.Type))
		return config, errs
	}
	switch strings.ToLower(config.Type) {
	case "postgre":
		config.Entity = Postgre{
			Config: config,
		}
	case "mysql":
		//TODO 
		config.Entity = MySQL{
			Config: config,
		}
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(100003, "debug", "Database type %s is right.", config.Type))

	// Checking db host is validly. 
	host := config.Host
	if len(host) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100004, "fatal", "The db host(%s) is empty or the length of it is less 1", host))
		return config, errs
	}
	if ip := net.ParseIP(host); ip != nil {
		config.HostIP = ip
	} else {
		ips, err := net.LookupIP(host)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(100005, "fatal", "We can not get the ip address for hostname:%s error: %s", host, err))
			return config, errs
		}
		config.HostIP = ips[0]
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(100006, "debug", "Database server host %s is right.", host))

	// Checking db port is validly. 
	port := config.Port
	if port < 1024 || port > 65535 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100007, "fatal", "DB Port(%d) should be large than 1024 and less 65535", port))
		return config, errs
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(100008, "debug", "Database server port %d is right.", port))

	if config.SslMode == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100009, "debug", "sslmode not be set. It will be set to disable"))
		config.SslMode = "disable"
	}

	// Checking certification files if sslmode is not disable. 
	// the error level will be set to warning if any of certification files is not exist and  slmode is not disable,
	// then set sslmode to disable.
	if strings.ToLower(config.SslMode) != "disable" {
		ca, errCa := getFile(config.SslCa, cmdRunPath)
		cert, errCert := getFile(config.SslCert, cmdRunPath)
		key, errKey := getFile(config.SslKey, cmdRunPath)

		if errCa != nil || errCert != nil || errKey != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(100010, "warning", "DB SslMode has be set to %s ,but ca(%s), cert(%s) or key(%s) not exist. We try to connect to DB server with disable sslmode.P", config.SslMode, config.SslCa, config.SslCert, config.SslKey))
			config.SslMode = "disable"
			config.SslCa = ""
			config.SslCert = ""
			config.SslKey = ""
		} else {
			config.SslCa = ca
			config.SslCert = cert
			config.SslKey = key
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100011, "debug", "the certification files has be checked."))
	}

	if config.MaxOpenConns < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100012, "warning", "the value of MaxOpenConns(Now: %d) must be large 1. It will be set to 10", config.MaxOpenConns))
		config.MaxOpenConns = 10
	}

	if config.MaxOpenConns < 10 || config.MaxOpenConns > 2000 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100013, "warning", "the value of MaxOpenConns(Now: %d) should be set large 10 and less 2000. ", config.MaxOpenConns))
	}

	if config.MaxIdleConns < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100014, "warning", "the value of MaxIdleConns(Now: %d) must be large 1. It will be set to 10", config.MaxIdleConns))
		config.MaxIdleConns = 10
	}

	if config.MaxIdleConns < 10 || config.MaxIdleConns > 2000 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100015, "warning", "the value of MaxIdleConns(Now: %d) should be set large 10 and less 2000. ", config.MaxIdleConns))
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(100016, "debug", "all database configuration parametes have be checked."))

	return config, errs
}

// Converting relative path to absolute path of  file(such as certs) and return the  file path
// return "" and error if  file can not opened .
// Or return string and nil.
func getFile(f string, cmdRunPath string) (string, error) {
	dir, error := filepath.Abs(filepath.Dir(cmdRunPath))
	if error != nil {
		return "", error
	}

	if !filepath.IsAbs(f) {
		tmpDir := filepath.Join(dir, "../")
		f = filepath.Join(tmpDir, f)
	}

	fp, err := os.OpenFile(f, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	fp.Close()

	return f, nil
}

// IsSupportedDB checking the DB of t has be supported now.
// return true if DB is supported
// otherwise return false
func IsSupportedDB(t string) bool {
	for _, v := range SupportDBs {
		if strings.EqualFold(v, t) {
			return true
		}
	}

	return false
}

// check whether dbname is valid.
func CheckIdentifier(dbType string, identifier string) bool {
	switch strings.ToLower(dbType) {
	case "postgre":
		entity := Postgre{}
		return entity.Identifier(identifier)
	case "mysql":
		entity := MySQL{}
		return entity.Identifier(identifier)
	default:
		return false
	}
}

