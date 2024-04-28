/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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
	"net/http"
	"sysadm/config"
)

type RunConf struct {
	// workingDir is X/../ ,X is the path of the directory which binary package of agent locate in it.
	workingDir string
	version    config.Version
	apiServer
	Debug   bool
	LogFile string
}

type runTimeConf struct {
	// keep http or https client for reuse. we should recreate http client if the value of this field is nil
	httpClient *http.Client
}

type RuntimeData struct {
	RunConf
	runTimeConf
}

type apiServer struct {
	// IP address or hostname of apiServer where agent will be connected to
	Address string `form:"address" json:"address" yaml:"address" xml:"address"`

	// service port of apiServer where agent will be connected to
	Port int `form:"port" json:"port" yaml:"port" xml:"port"`

	// tls parameters which agent will use to connect to apiServer
	config.Tls `form:"tls" json:"tls" yaml:"tls" xml:"tls"`
}

var RunData = RuntimeData{}
