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

import (
	"fmt"
	"strings"
)

func (p Project) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("can not get project information with empty ID")
	}

	dbData, e := getObjectInfoByID(p.TableName, p.PkName, id)
	if e != nil {
		return nil, e
	}

	projectData := ProjectSchema{}
	e = Unmarshal(dbData, &projectData)

	return projectData, e
}

func (p Project) setDefaultForObject() {
	p.Name = "project"
	p.TableName = "project"
	p.PkName = "projectid"
}
