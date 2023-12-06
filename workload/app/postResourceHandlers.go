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
	"github.com/wangyysde/sysadmServer"
	"net/http"
	"sysadm/sysadmLog"
	"sysadm/sysadmapi/apiutils"
	"sysadm/user"
)

func addNewObjectHandler(c *sysadmServer.Context, module string) {
	var errs []sysadmLog.Sysadmerror
	var response apiutils.ApiResponseData

	errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001000001, "debug", "now handling %s adding", module))
	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		response = apiutils.BuildResponseDataForError(80001000002, "您没有登录或者没有权限执行本操作")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001000002, "info", "user has not login or not permission when adding new namespace"))
		runData.logEntity.LogErrors(errs)
		c.JSON(http.StatusOK, response)
		return
	}

	objEntity, e := newObjEntity(module)
	if e != nil {
		response = apiutils.BuildResponseDataForError(80001000003, "请求的地址不正确，请确认是从正确地方连接过来的")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001000003, "error", "module %s was not defined", module))
		runData.logEntity.LogErrors(errs)
		c.JSON(http.StatusOK, response)
		return
	}
	requestKeys := []string{"dcID", "clusterID", "namespace"}
	requestData, e := getRequestData(c, requestKeys)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001000004, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response = apiutils.BuildResponseDataForError(80001000004, "系统出错，请稍后再试或者联系系统管理员")
		c.JSON(http.StatusOK, response)
		return
	}

	objEntity.setObjectInfo()
	if objEntity.getNamespaced() && (requestData["namespace"] == "" || requestData["namespace"] == "0") {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001000005, "info", "resource is namespaced, but namespace is not specified"))
		runData.logEntity.LogErrors(errs)
		response = apiutils.BuildResponseDataForError(80001000005, "系统出错，请稍后再试或者联系系统管理员")
		c.JSON(http.StatusOK, response)
		return

	}

	e = objEntity.addNewResource(c, module)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001000006, "info", "%s", e))
		runData.logEntity.LogErrors(errs)
		response = apiutils.BuildResponseDataForError(80001000006, "系统出错，请稍后再试或者联系系统管理员")
		c.JSON(http.StatusOK, response)
		return
	}

	response = apiutils.BuildResponseDataForSuccess("内容已经添加成功")
	c.JSON(http.StatusOK, response)
	return
}

func actionNotFoundHandler(c *sysadmServer.Context, module, action string) {
	var errs []sysadmLog.Sysadmerror
	var response apiutils.ApiResponseData

	errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001000007, "info", "action %s for module %s was not defined", action, module))
	runData.logEntity.LogErrors(errs)
	response = apiutils.BuildResponseDataForError(80001000007, "请求地址出错，请确认是从正确的地方连接过来的")
	c.JSON(http.StatusOK, response)
}
