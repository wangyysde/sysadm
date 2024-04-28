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
	runtime "sysadm/apimachinery/runtime/v1beta1"
)

var scheme = &runtime.Scheme{}
var registeredResources map[runtime.GroupVersion]*runtime.RegistryType

func prepareSchema() error {
	registeredResources = runtime.GetRegisteredResources()

	if len(registeredResources) < 1 {
		return fmt.Errorf("no registered resource")
	}

	for _, r := range registeredResources {
		e := scheme.AddSchemeData(r)
		if e != nil {
			return e
		}
		if r.ConversionRegFn != nil {
			fn := r.ConversionRegFn
			e := fn(scheme)
			if e != nil {
				return e
			}
		}
	}

	return nil
}
