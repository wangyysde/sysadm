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

// +sysadm:
package v1beta1

import (
	runtimeBetaV1 "sysadm/apimachinery/runtime/v1beta1"
)

// GroupName is the group name use in this package
const GroupName = "sys.sysadm.cn"

// SchemeGroupVersion is group version used to register these objects
var SchemaGroupVersion = runtimeBetaV1.GroupVersion{GroupName, "v1beta1"}

// TypeRegistryFunc used to register this resource type when sysadm-apiserver start
// the name of this variable MUST NOT be changed
var TypeRegistryFunc runtimeBetaV1.FuncRegistry = addNewType

func addNewType(schema *runtimeBetaV1.Scheme) error {
	return schema.AddKnowTypes(SchemaGroupVersion,
		&Command{})
}
