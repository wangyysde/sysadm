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

package v1beta1

const (
	ContentTypeJSON string = "application/json"
	ContentTypeYAML string = "application/yaml"
	//	ContentTypeProtobuf string = "application/vnd.kubernetes.protobuf"

	APIVersionInternal string = "__internal"
)

var (
	registeredResources = make(map[GroupVersion]*RegistryType)
)

const (
	// Create creates a new resource HTTP verb is POST
	Create VerbKind = 1 << iota

	// Delete delete a resource from system  HTTP verb is DELETE
	Delete

	// DeleteCollection delete a set of a resource HTTP verb is DELETE
	DeleteCollection

	// Get gets the information of a resource. HTTP verb is GET
	Get

	// List lists the information of a resource list HTTP verb is GET
	List

	// HTTP verb is Patch
	Patch

	// Update update the information of a resource. HTTP verb is PUT
	Update

	// Watch watching an individual resource or collection of resources. HTTP verb is GET
	Watch
)

const (
	ResourcePkFieldName                  = "ID"
	ResourcepKDbFieldName                = "id"
	ResourceRelationParentDBFieldName    = "object_id"
	ResourceRelationChildDBFieldName     = "child_id"
	ReferenceFieldName                   = "ReferenceObject"
	ResourceReferenceDBObjectIdFieldName = "objectId"
)
