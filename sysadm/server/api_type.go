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
	"github.com/wangyysde/sysadmServer"
)

type ApiResponseBodyData map[string] string
type ApiResponseStatus struct {
	// Status is false if this is a error response, otherwise Status is true
	Status bool `json:"status"`
	// Errorcode is zero if this is a successful response, otherwise Errorcode is nonzero
	Errorcode int `json:"errorcode"`
	// Message is the result sets if this is a successful ,otherwise Message is []map[string]string
	// which has one rows only:["errorMsg"] = errorMsg
	Message interface{} `json:"message"`
}

type apiHander func (c *sysadmServer.Context)
type Module struct {
	Name string
	Path string
	Instance ModuleInterface
	Actions []string
 }

 type User struct {}
 type Registry struct {}
 type Sysadm struct {}
 type Project struct {
	Projectid int `json:"projectid"`
	Ownerid int `json:"ownerid"`
	Name string `json:"name"`
	Comment string `json:"comment"`
	Deleted bool `json:"deleted"`
	Creation_time int `json:"creation_time"`
	Update_time int `json:"update_time"`
 }

 type ModuleInterface interface {
	ModuleName() string
	ActionHanderCaller(action string, c *sysadmServer.Context)
 }


var Modules = map[string]Module{
	"user": {
		Name: "user",
		Path: "user",
		Instance: User{},
		Actions: userActions,
	},
	"registry": {
		Name: "registry",
		Path: "registry",
		Instance: Registry{},
	},
	"sysadm": {
		Name: "sysadm",
		Path: "sysadm",
		Instance: Sysadm{},
	},
	"project": {
		Name: "project",
		Path: "project",
		Instance: Project{},
		Actions: projectActions,
	},
 }