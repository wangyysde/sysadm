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
* structures for apiserver configurations
 */

package app

import (
	"context"
	sysadmDB "sysadm/db"

	"sysadm/config"
	"sysadm/redis"
	"sysadm/sysadmLog"
)

// for global block
type ConfGlobal struct {

	// set whether apiserver running in debug mode
	Debug bool `form:"debug" json:"debug" yaml:"debug" xml:"debug" `

	// virtualIP is the  virtual IP for apiserver(s) when apiserver(s) are behind a LB or HA
	VirtualIP []string `form:"virtualIP" json:"virtualIP" yaml:"virtualIP" xml:"virtualIP"`

	// Optional extra Subject Alternative Names (SANs) to use for the API Server serving certificate. Can be both IP addresses and DNS names.
	ExtraSans []string `form:"extraSans" json:"extraSans" yaml:"extraSans" xml:"extraSans"`
}

// for server block
type ConfServer struct {
	config.Server

	// insecret specifies whether apiserver listen on a insecret port when it is runing as daemon
	Insecret bool `form:"insecret" json:"insecret" yaml:"insecret" xml:"insecret"`

	// insecret listen port of apiserver listening when it is running ad daemon
	InsecretPort int `form:"insecretPort" json:"insecretPort" yaml:"insecretPort" xml:"insecretPort"`

	// IsTls identify whether apiServer listening TLS Port
	// apiServer will get certificate from DB if IsTls is true
	IsTls bool `form:"isTls" json:"isTls" yaml:"isTls" xml:"isTls"`

	// ca  path of apiServer if apiServer listen on TLS
	Ca string `form:"ca" json:"ca" yaml:"ca" xml:"ca"`

	// certification  path of apiServer if apiServer listen on TLS
	Cert string `form:"cert" json:"cert" yaml:"cert" xml:"cert"`

	// key path of apiServer if apiServer listen on TLS
	Key string `form:"key" json:"key" yaml:"key" xml:"key"`
}

// for DB block
type ConfDB struct {
	// tppe, one of mysql,postgre
	Type string `form:"type" json:"type" yaml:"type" xml:"type"`

	// DB name
	DBName string `form:"dbName" json:"dbName" yaml:"dbName" xml:"dbName"`

	config.Server

	config.Tls

	config.User

	// max number of connections of concurrent openned
	MaxOpenConns int `form:"maxOpenConns" json:"maxOpenConns" yaml:"maxOpenConns" xml:"maxOpenConns"`

	// max number of idle connections
	MaxIdleConns int `form:"maxIdleConns" json:"maxIdleConns" yaml:"maxIdleConns" xml:"maxIdleConns"`
}

// apiserver configuration
type Conf struct {
	// version information for apiserver
	Version config.Version

	// save the path of apiserver's configuration file
	ConfFile string

	// hold global block items
	ConfGlobal ConfGlobal `form:"global" json:"global" yaml:"global" xml:"global"`

	// hold server block items
	ConfServer ConfServer `form:"server" json:"server" yaml:"server" xml:"server"`

	// hold log block items
	ConfLog config.Log `form:"log" json:"log" yaml:"log" xml:"log"`

	// hold redis block items
	ConfRedis redis.ClientConf `form:"redis" json:"redis" yaml:"redis" xml:"redis"`

	// hold db block items
	ConfDB ConfDB `form:"db" json:"db" yaml:"db" xml:"db"`
}

// hold running data
type runningData struct {
	// root path of apiserver working
	workingRoot string

	// redis entity
	redisEntity redis.RedisEntity

	// redisContext
	redisCtx context.Context

	// db entity
	dbEntity sysadmDB.DbEntity

	// configuration data for apiserver running
	runConf Conf

	// logger entity
	logEntity *sysadmLog.LoggerConfig
}

var runData runningData = runningData{
	workingRoot: "",
	redisEntity: nil,
	redisCtx:    nil,
	dbEntity:    nil,
	runConf: Conf{
		Version:    config.Version{},
		ConfFile:   "",
		ConfGlobal: ConfGlobal{},
		ConfServer: ConfServer{},
		ConfLog:    config.Log{},
		ConfRedis:  redis.ClientConf{},
		ConfDB:     ConfDB{},
	},
	logEntity: nil,
}
