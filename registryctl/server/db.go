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
	sysadmDB "github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/registryctl/config"
	"github.com/wangyysde/sysadm/sysadmerror"
)

/* 
	initDB checking the configuration of DB and then initating DB connections.
	return an entity and errors when a connection to DB server opened successfully
	otherwise nil and errors when a connection to DB server opened false
*/
func initDB(definedConfig *config.Config,cmdRunPath string)(sysadmDB.DbEntity,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202008,"debug","try to initating DB entity..."))
	dbConf := sysadmDB.DbConfig {
		Type: definedConfig.DB.Type,
		Host: definedConfig.DB.Host, 
		Port: definedConfig.DB.Port, 
		User: definedConfig.DB.User,
		Password: definedConfig.DB.Password,
		DbName: definedConfig.DB.Dbname,
		SslMode: definedConfig.DB.Sslmode,
		SslCa: definedConfig.DB.Sslrootcert,
		SslCert: definedConfig.DB.Sslcert,
		SslKey: definedConfig.DB.Sslkey, 
		MaxOpenConns: definedConfig.DB.DbMaxConnect,
		MaxIdleConns: definedConfig.DB.DbMaxConnect,
		Connect: nil,
		Entity: nil,
	}

	newDBConf,err := sysadmDB.InitDbConfig(&dbConf,cmdRunPath)
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

	RuntimeData.RuningParas.DBConfig = newDBConf

	return entity,errs
}