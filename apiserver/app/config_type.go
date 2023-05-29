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
)

// for global block
type ConfGlobal struct {
	// apiserver running mode. apiserver response command data to client and receive command status and command logs from client
	// when this item is true. otherwise apiserver send command data to client actively and get command statuses and command logs from clients.
	Passive bool `form:"passive" json:"passive" yaml:"passive" xml:"passive"`

	// set whether apiserver running in debug mode
	Debug bool `form:"debug" json:"debug" yaml:"debug" xml:"debug"`

	// set whether apiserver running as daemon
	// TODO
	Daemon bool `form:"daemon" json:"daemon" yaml:"daemon" xml:"daemon"`

	// specifies the uri where client get commands to run when apiserver runing as daemon in passive mode. default value is "/getCommand"
	// in other word, apiserver is listening this path for client getting command data when apiserver running in passive mode.
	// when apiserver is run in active mode, apiserver send command data to client on the path of this item specified.default value is "/receiveCommand"
	CommandUri string `form:"commandUri" json:"commandUri" yaml:"commandUri" xml:"commandUri"`

	// specifies the uri where client send command status to when apiserver running as daemon in passive mode.default value is "/receiveCommandStatus"
	// in other word, apiserver is listening this path for client send command status to  when apiserver running in passive mode.
	// when apiserver is run in active mode, apiserver send command status to client on the path of this item specified.default value is "/getCommandStatus"
	CommandStatusUri string `form:"commandStatusUri" json:"commandStatusUri" yaml:"commandStatusUri" xml:"commandStatusUri"`

	// specifies the uri where client send command logs to when apiserver running as daemon in passive mode.default value is "/receiveLogs"
	// in other word, apiserver is listening this path for client send command logs to  when apiserver running in passive mode.
	// when apiserver is run in active mode, apiserver send command status to client on the path of this item specified.default value is "/getLogs"
	CommandLogsUri string `form:"commandLogsUri" json:"commandLogsUri" yaml:"commandLogsUri" xml:"commandLogsUri"`

	// interval of checking new command for client by apiserver when apiserver is running actively. default is 5 second.
	CheckCommandInterval int `form:"checkCommandInterval" json:"checkCommandInterval" yaml:"checkCommandInterval" xml:"checkCommandInterval"`

	// interval of try to get command status from client by apiserver when apiserver is running actively. default is 5 second
	GetStatusInterval int `form:"getStatusInterval" json:"getStatusInterval" yaml:"getStatusInterval" xml:"getStatusInterval"`

	// interval of try to get command log from client by apiserver when apiserver is running actively. default is 5 second
	GetLogInterval int `form:"getLogInterval" json:"getLogInterval" yaml:"getLogInterval" xml:"getLogInterval"`
}

// for server block
type ConfServer struct {
	config.Server

	// insecret specifies whether apiserver listen on a insecret port when it is runing as daemon
	Insecret bool `form:"insecret" json:"insecret" yaml:"insecret" xml:"insecret"`

	// insecret listen port of apiserver listening when it is running ad daemon
	InsecretPort int `form:"insecretPort" json:"insecretPort" yaml:"insecretPort" xml:"insecretPort"`

	config.Tls
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

	// hold the running data related to command. this items can be configurable in the feature
	command runDataForCommand
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
	command: runDataForCommand{},
}
