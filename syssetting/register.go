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

package syssetting

import (
	runtime "sysadm/apimachinery/runtime/v1beta1"
)

// GroupName is the group name use in this package
const GroupName = "syssetting.sysadm.cn"

// SchemeGroupVersion is group version used to register these objects
var SchemaGroupVersion = runtime.GroupVersion{GroupName, runtime.APIVersionInternal}

// TypeRegistryFunc used to register this resource type when sysadm-apiserver start
// the name of this variable MUST NOT be changed
var TypeRegistryFunc runtime.FuncRegistry = AddNewType

var allowedVerbs runtime.VerbKind = 0

func AddNewType(schema *runtime.Scheme) error {
	return schema.AddKnowTypes(SchemaGroupVersion, allowedVerbs,
		&Syssetting{})
}

func GetKind() (string, error) {
	return runtime.GetKindByType(&Syssetting{})
}

func init() {
	runtime.Register(SchemaGroupVersion, TypeRegistryFunc, nil)
	return
}
