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

package main

// The possible classes of types.
type TypeKind string

const (
	defaultGeneratedRuntimeCodeFilename string = "zz_generated.runtime.go"
	defaultGoHeadFileName               string = "CopyRight"
	defaultModuleName                   string = "sysadm"

	// Builtin is a primitive, like bool, string, int.
	Builtin TypeKind = "Builtin"
	Struct  TypeKind = "Struct"
	Map     TypeKind = "Map"
	Slice   TypeKind = "Slice"
	Pointer TypeKind = "Pointer"
)

var (
	generatedRuntimeCodeFilename = ""
	goHeadFileName               = ""
	moduleName                   = ""
	sysadmRoot                   = ""
)
