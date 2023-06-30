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
 */

package app

type Project struct {
	// object name. the value should be "project"
	Name string
	// table name which hold project data in DB
	TableName string
	// field name of primary key in the table
	PkName string
}

type ProjectSchema struct {
	//projectid identified a project
	ProjectID int `form:"id" json:"id" yaml:"id" xml:"id"`
	// the owner of the project. owner is the user who created the project normally
	UserID int `form:"userid" json:"userid" yaml:"userid" xml:"userid"`
	// project name. it must be a string in english. Is is a part of url of image.
	Name string `form:"name" json:"name" yaml:"name" xml:"name"`
	// description of a project
	Comment string `form:"comment" json:"comment" yaml:"comment" xml:"comment"`
	// the value is true if a user has be deleted
	Deleted int `form:"deleted" json:"deleted" yaml:"deleted" xml:"deleted"`
	// the time when the project has be create
	Creation_time int `form:"creation_time" json:"creation_time" yaml:"creation_time" xml:"creation_time"`
	// the time when the project has be update
	Update_time int `form:"update_time" json:"update_time" yaml:"update_time" xml:"update_time"`
}
