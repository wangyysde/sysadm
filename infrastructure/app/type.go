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

package app

import (
	"sysadm/config"
	sysadmDB "sysadm/db"
	"sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

// This structure is for configuration of apiServer Block
type ApiServer struct {
	Server     config.Server `form:"server" json:"server" yaml:"server" xml:"server"`
	Tls        config.Tls    `form:"tls" json:"tls" yaml:"tls" xml:"tls"`
	ApiVersion string        `form:"apiVersion" json:"apiVersion" yaml:"apiVersion" xml:"apiVersion"`
}

// Saving running data of an instance
type RuningData struct {
	dbConf      *sysadmDB.DbConfig
	logConf     *config.Log
	apiServer   *ApiServer
	workingRoot string
}

// Initating working data for an instance
var WorkingData RuningData = RuningData{
	dbConf:      nil,
	logConf:     nil,
	apiServer:   nil,
	workingRoot: "",
}

// Saving the status of a host
type HostStatus struct {
	// Identifier of a status come from DB
	StatusID int

	// Status Name. not null
	Name string

	// description of a status. Maybe null
	Description string
}

// Saving host ip
type HostIP struct {
	// Identifier of a IP come from DB
	IpID int

	// interface name of which the ip set
	DevName string

	// IP address for IPV4
	Ipv4 string

	// netmask address for ipv4 address
	Maskv4 string

	// IP address for IPV6
	Ipv6 string

	// netmask address for ipv6 address
	Maskv6 string

	// 0 for offline, 1 for online
	Status int

	// 0 not management ip, 1 management IP
	isManage int
}

// Saving host user
type HostUser struct {
	// userID identified a user on a host come from DB
	UserID int

	// username on a host
	UserName string

	// password encode by base64
	SecurePassword string

	// clear password
	ClearPassword string
}

// OS information
type Os struct {
	// specify which OS distrubition,such as centos,readhat, ubantu. come from DB
	OsID int

	// distribution name.such as centos,redhat. this field must be unique
	Name string

	// distribution description
	Description string
}

// OS Version information
type OsVersion struct {
	// version id identified a version come from DB
	VersionID int

	// version name
	Name string

	// OS Type
	Os *Os

	// description of the version
	Description string
}

// host infromation
type Host struct {
	// hostid identified a host
	Hostid int

	// host name of OS
	Hostname string

	// OsVersion information
	OsVersion *OsVersion

	// host status
	Status *HostStatus

	// IP list with a host
	HostIps []HostIP

	// User list on a host
	HostUsers []HostUser
}

// Infrastructure
type Infrastructure struct {
	ModuleName string
	ApiVersion string
}

type handlerAdder func(*sysadmServer.Engine, string, Infrastructure) []sysadmerror.Sysadmerror

// define host information for API server using
type ApiHost struct {
	// host name of OS
	Hostname string `form:"hostname" json:"hostname" xml:"hostname" binding:"-"`

	// ip address for connecting to
	Ip string `form:"ip" json:"ip" xml:"ip"`

	// ip type 4 for ipv4 6 for ipv6
	Iptype string `form:"iptype" json:"iptype" xml:"iptype"`

	// wheather agent running in passive mode
	PassiveMode int `form:"passiveMode" json:"passiveMode" xml:"passiveMode"`

	// agent listen port number
	AgentPort int `form:"agentPort" json:"agentPort" xml:"agentPort"`

	// the path where agent listen to receiving command when is running in active mode
	ReceiveCommandUri string `form:"receiveCommandUri" json:"receiveCommandUri" xml:"receiveCommandUri"`

	// which os has be installed on a node. the value of osid is reference to table os in DB
	OsID int `form:"osid" json:"osid" xml:"osid"`

	// which version of os has be installed on a node. the value of osversionid is reference to table version in DB
	OsVersionID int `form:"osversionid" json:"osversionid" xml:"osversionid"`

	// which yum information should be deploy on a node. the value of yumid is reference to table yum in DB
	YumID []string `form:"yumid[]" json:"yumid[]" xml:"yumid[]"`

	// operator userid what used to check permissions
	Userid int `form:"userid" json:"userid" xml:"userid"`
}
