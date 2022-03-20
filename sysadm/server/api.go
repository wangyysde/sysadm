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
 */

package server

import (
	//	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

// addFormHandler set delims for template and load template files
// return nil if not error otherwise return error.
func addApiHandler(r *sysadmServer.Engine,cmdRunPath string) {
	// Simple group: v1
	v1 := r.Group("/api/v1.0")
	{
		v1.POST("/:module/*action", apiHandlers)
//		v1.POST("/submit", submitEndpoint)
//		v1.POST("/read", readEndpoint)
	}

}

func apiHandlers(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror
	module := strings.TrimSuffix(strings.TrimPrefix(c.Param("module"),"/"),"/")
	action := strings.TrimSuffix(strings.TrimPrefix(c.Param("action"),"/"),"/")
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1030001,"debug","now handling the request for module %s with action %s.",module,action))
	if !foundModule(module) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1030002,"error","parameters error. module %s was not found.",module))
		logErrors(errs)
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1030002,
			Message: fmt.Sprintf("parameters error.module %s not found", module),
		}
		//respBody,_ := json.Marshal(ret)
		c.JSON(http.StatusOK, ret)
		return 
	}

	if !foundAction(module,action) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1030003,"error"," parameters error. action %s was not found in module %s.",action,module))
		logErrors(errs)
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1030003,
			Message: fmt.Sprintf("parameters error.action %s not found", action),
		}
		//respBody,_ := json.Marshal(ret)
		c.JSON(http.StatusOK, ret)
		return 
	}

	action =strings.ToLower(action)
	mI := Modules[module].Instance
	mI.ActionHanderCaller(action,c)

	return 
}

func foundModule(module string) bool {

	found := false
	for k := range Modules {
		if strings.EqualFold(k,module) {
			found = true
			break
		}
	}

	return found
}

func foundAction(m string, action string)bool{
	actions := Modules[m].Actions
	found := false
	for _,value := range actions {
		if strings.EqualFold(value,action) {
			found = true
			break
		}
	}
	 
	return found
}

