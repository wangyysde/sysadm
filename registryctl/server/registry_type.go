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

Ref: https://docs.docker.com/registry/spec/api/
	https://datatracker.ietf.org/doc/rfc7235/
*/

package server

type blob struct {
	digest string
	size int64
}

type image struct {
	username string
	name string
	size int64
	tag string
	architecture string
	digest string
	blobs []blob
}

var processImages map[string]image = make(map[string]image,0)

type FsLayer struct {
	BlobSum string `json:"blobSum"`
}

type History struct {
	V1Compatibility string `json:"v1Compatibility"`
}
type Manifest struct {
	SchemaVersion int `json:"schemaVersion"`
	Name string `json:"name"`
	Tag string `json:"tag"`
	Architecture string `json:"architecture"`
	FsLayers []FsLayer `json:"fsLayers"`
	History []History `json:"history"`
}

type RegistryCtl struct {}

var  registryctlActions = []string{"imagelist","getcount","taglist"}