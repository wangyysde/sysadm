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
	"sysadm/config"
	sysadmDB "sysadm/db"
)

type DBServer struct {
	Type string `json:"type"`
	DBName string `json:"dbName"`
	Server config.Server `json:"server"`
	Tls config.Tls `json:"tls"`
	Credit config.User `json:"credit"`
	MaxOpenConns int `json:"maxOpenConns"`
	MaxIdleConns int `json:"maxIdleConns"`
}

// This structure is for configuration of apiServer Block
type ApiServer struct {
	Server config.Server `form:"server" json:"server" yaml:"server" xml:"server"`
	Tls config.Tls `form:"tls" json:"tls" yaml:"tls" xml:"tls"`
	ApiVersion string `form:"apiVersion" json:"apiVersion" yaml:"apiVersion" xml:"apiVersion"`
}

type Config struct {
	Version config.Version 
	Server config.Server `json:"server,omitempty"`
	ServerTls config.Tls `json:"tls"`
	Log config.Log `json:"log"`
	DB DBServer `json:"db,omitempty"`
	ApiServer ApiServer `form:"apiServer" json:"apiServer" yaml:"apiServer" xml:"apiServer"`
}

type CliOptions struct {
	CfgFile string 
	// absolute path to the working root
	workingRoot string
}

type workingData struct {
	dbConf *sysadmDB.DbConfig
}

type RuningData struct {
	Options CliOptions 
	Config Config
	workingData workingData
}

var CurrentRuningData *RuningData = nil
var LastRuningData *RuningData = nil