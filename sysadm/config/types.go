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

import(
	"github.com/wangyysde/sysadmServer"
)

//Defining server configuration
type Server struct {
	Address string `json:"address"`
	Port int `json:"port"`
	Socket string `json:"socket"`
	Tls bool `json:"tls"`
	Ca string `json:"ca"`
	Cert string `json:"cert"`
	Key string `json:"key"`
}

//Defining log configuration 
type Log struct {
	AccessLog string `json:"accessLog"`
	ErrorLog string `json:"errorLog"`
	Kind string `json:"kind"`
	Level string `json:"level"`
	SplitAccessAndError bool `json:"splitAccessAndError"`
	TimeStampFormat string `json:"timeStampFormat"`
}

type User struct {
	DefaultUser string `json:"defaultUser"`
	DefaultPassword string `json:"defaultPassword"`
}

type DB struct {
	// DB type One of postgre,mysql.If type is not set, then postgre will be set to it. But postgre is available only now TODO MySQL ....
	Type string `json:"type"`
	// DB host address. Normalll,it is either IP adress or hostname of DB server.
	Host string `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Dbname string `json:"dbname"`
	DbMaxConnect int `json:"dbMaxConnect"`
	DbIdleConnect int `json:"dbIdleConnect"`
	Sslmode string `json:"sslmode"`
	Sslrootcert string `json:"sslrootcert"`
	Sslkey string `json:"sslkey"`
	Sslcert string `json:"sslcert"`
}

type Config struct {
	Version string `json:"version"`
	Server Server `json:"server"`
	Log Log `json:"log"`
	User User `json:"user"`
	DB DB `json:"db"`
}

var DefinedConfig Config = Config{}

var defaultConfig Config = Config{
	Version:  Version,
	Server: Server {
		Address: DefaultIP,
		Port: DefaultPort,
		Socket: DefaultSocket,
		Tls: DefaultTls,
		Ca: DefaultCa,
		Cert: DefaultCert,
		Key: DefaultKey,
	},
	Log: Log{
		AccessLog: DefaultAccessLog,
		ErrorLog: DefaultErrorLog,
		Kind: DefaultLogKind,
		Level: DefaultLogLevel,
		SplitAccessAndError: true,
		TimeStampFormat: sysadmServer.TimestampFormat["DateTime"],
	},
	User: User{
		DefaultUser: DefaultUser,
		DefaultPassword: DefaultPasswd,
	},
	DB: DB{
		Type: DefaultDBType,
		Host: DefaultDbHost,
		Port: DefaultDbPort,
		User: DefaultDbUser,
		Password: DefaultDbPassword,
		Dbname: DefaultDbDbName,
		DbMaxConnect: DefaultDbMaxConnect,
		DbIdleConnect: DefaultDbIdleConnect,
		Sslmode: DefaultDbSslmode,
		Sslrootcert: DefaltDbSslrootcert,
		Sslkey: DefaultDbSslkey,
		Sslcert: DefaultDbSslcert,
	},
}