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
	"fmt"
	"strings"
	runtime "sysadm/apimachinery/runtime/v1beta1"
)

func getGvkBasedPath(path string) (*runtime.GroupVersionKind, error) {
	path = strings.TrimSpace(path)
	// resource path should be form /api/<version>/<kind>/<action>
	if len(path) < 9 {
		return nil, fmt.Errorf("request path %s is not a valid resource request path", path)
	}

	pathSlice := strings.Split(path, "/")
	if len(pathSlice) != 5 {
		return nil, fmt.Errorf("request path %s is not a valid resource request path.resource request path must "+
			"be like /api/<version>/<kind>/<action>", path)
	}
	version := strings.TrimSpace(pathSlice[2])
	kind := strings.TrimSpace(pathSlice[3])
	if version == "" || kind == "" {
		return nil, fmt.Errorf("request path %s is not a valid resource request path", path)
	}

	group := scheme.GetGroupByKind(kind)
	if group == "" {
		return nil, fmt.Errorf("we can not get the group of the resource")
	}

	gvk := runtime.GroupVersionKind{Group: group, Version: version, Kind: kind}

	return &gvk, nil
}
