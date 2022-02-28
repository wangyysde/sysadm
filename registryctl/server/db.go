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
	"strings"

	sysadmDB "github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/registryctl/config"
	"github.com/wangyysde/sysadm/sysadmerror"
)

// initDBEntity initating a dbConfig accroding configurations
// return sysadmDB.DbConfig
 func initDBConfig(definedConfig *config.Config )(sysadmDB.DbConfig) {
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

	return dbConf
 }

 // initDBEntity initating a DB interface entity accroding to the DB configurations.
 // return sysadmDB.DbEntity
 func initDBEntity(dbConfig *sysadmDB.DbConfig)(sysadmDB.DbEntity,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var entity sysadmDB.DbEntity
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202008,"debug","try to initating DB entity..."))
	switch strings.ToLower(dbConfig.Type) {
	case "postgre":
		entity = sysadmDB.Postgre{
			Config: dbConfig,
		}
	case "mysql":
		//TODO 
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202009,"fatal","we can support Postgre DB only now."))
		//config.Entity = Postgre{
		//	Config: config,
		//}
	default:
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202010,"fatal","we can support the following DB serer:%v .",sysadmDB.SupportDBs))
	}

	return entity,errs
}

// openDBConnection open a connection to DB server accroding to the  DBConfig
// closeDBConnection calling should be followed with this function calling.
func openDBConnection(dbConfig *sysadmDB.DbConfig)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202011,"debug","now try to connect to %s db server(%s)",dbConfig.Host,dbConfig.Type))

	entity := dbConfig.Entity
	e := entity.OpenDbConnect()
	if len(e) > 0 {
		errs = appendErrs(errs,e)
	}

	return errs
}

// closeDBConnection closing the connection to the DB.
// the calling of this function should be following with the calling of openDBConnection
func closeDBConnection(dbConfig *sysadmDB.DbConfig){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202016,"debug","now try to closing the connection to %s db server(%s)",dbConfig.Host,dbConfig.Type))
	entity := dbConfig.Entity
	errs =  entity.CloseDB()
	logErrors(errs)
}