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

package server

import (
	sysadmDB "sysadm/db"
	"sysadm/sysadmerror"
)

/*
	initDB checking the configuration of DB and then initating DB connections.
	return an entity and errors when a connection to DB server opened successfully
	otherwise nil and errors when a connection to DB server opened false
*/
func initDB()(sysadmDB.DbEntity,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	definedConf := &CurrentRuningData.Config
	var sslModle string = ""
	if !definedConf.DB.Tls.IsTls {
		sslModle = "disable"
	} else {
		sslModle = "enable"
	}

	dbConf := sysadmDB.DbConfig {
		Type: definedConf.DB.Type,
		Host: definedConf.DB.Server.Address, 
		Port: definedConf.DB.Server.Port, 
		User: definedConf.DB.Credit.UserName,
		Password: definedConf.DB.Credit.Password,
		DbName: definedConf.DB.DBName,
		SslMode: sslModle,
		SslCa: definedConf.DB.Tls.Ca,
		SslCert: definedConf.DB.Tls.Cert, 
		SslKey: definedConf.DB.Tls.Key, 
		MaxOpenConns: definedConf.DB.MaxOpenConns, 
		MaxIdleConns: definedConf.DB.MaxIdleConns,
		Connect: nil,
		Entity: nil,
	}

	newDBConf,err := sysadmDB.InitDbConfig(&dbConf, CurrentRuningData.Options.workingRoot )
	errs = append(errs, err...)
	maxLevel := sysadmerror.GetMaxLevel(errs)
	fatalLevel := sysadmerror.GetLevelNum("fatal")
	if maxLevel >= fatalLevel {
		return nil,errs
	}

	entity := newDBConf.Entity
	err = entity.OpenDbConnect()
	errs = append(errs, err...)
	maxLevel = sysadmerror.GetMaxLevel(errs)
	if maxLevel >= fatalLevel {
		return nil,err
	}

	CurrentRuningData.workingData.dbConf = newDBConf

	return entity,errs
}