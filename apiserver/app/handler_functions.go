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
	"github.com/wangyysde/sysadmServer"
	"net/http"

	runtime "sysadm/apimachinery/runtime/v1beta1"
)

func addResourceHanders(r *sysadmServer.Engine) error {
	if r == nil {
		return fmt.Errorf("router is nil")
	}

	observeredKinds := scheme.GetObservedVersionKinds()
	if len(observeredKinds) < 1 {
		return fmt.Errorf("no version kind has be registered")
	}

	for _, ob := range observeredKinds {
		verbs := ob.Verbs
		gvk := ob.Gvk
		version := gvk.Version
		kind := gvk.Kind
		if (verbs & runtime.Create) == runtime.Create {
			r.POST("/api/"+version+"/"+kind+"/create", createResourceHandler)
		}

		if (verbs & runtime.Delete) == runtime.Delete {
			r.DELETE("/api/"+version+"/"+kind+"/delete", deleteResourceHandler)
		}

		if (verbs & runtime.DeleteCollection) == runtime.DeleteCollection {
			r.DELETE("/api/"+version+"/"+kind+"/deleteclollection", deletecollectionResourceHandler)
		}

		if (verbs & runtime.Get) == runtime.Get {
			r.GET("/api/"+version+"/"+kind+"/get", getResourceHandler)
		}

		if (verbs & runtime.List) == runtime.List {
			r.GET("/api/"+version+"/"+kind+"/list", listResourceHandler)
		}

		if (verbs & runtime.Patch) == runtime.Patch {
			r.PATCH("/api/"+version+"/"+kind+"/patch", patchResourceHandler)
		}

		if (verbs & runtime.Update) == runtime.Update {
			r.PUT("/api/"+version+"/"+kind+"/update", updateResourceHandler)
		}

		if (verbs & runtime.Watch) == runtime.Watch {
			r.GET("/api/"+version+"/"+kind+"/watch", watchResourceHandler)
		}

		r.Any("/api/"+version+"/"+kind+"/", noActionForResourceHandler)
	}

	return nil
}

// addRootHandler adding handler for root path
func addRootHandler(r *sysadmServer.Engine) error {
	if r == nil {
		return fmt.Errorf("router is nil")
	}

	r.Any("/", func(c *sysadmServer.Context) {
		c.JSON(http.StatusOK, sysadmServer.H{
			"status": "ok",
		})
	})

	return nil
}
