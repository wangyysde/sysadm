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
	"sysadm/sysadmapi/apiutils"
	"sysadm/user"
	"sysadm/utils"
)

func validCnNameHandler(c *sysadmServer.Context) {
	var response apiutils.ApiResponseData
	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		response = apiutils.BuildResponseDataForError(7000170001, "您没有登录或者没有权限执行本操作")
		c.JSON(http.StatusOK, response)
		return
	}
	requestData, e := utils.NewGetRequestData(c, []string{"objvalue"})

	if e != nil || !validCnName(requestData["objvalue"]) {
		response = apiutils.BuildResponseDataForError(7000170002, "中文名称不能为空,不能重复，且其长度不得大于255个字符.")
	} else {
		response = apiutils.BuildResponseDataForSuccess("ok")
	}

	c.JSON(http.StatusOK, response)
}

func validEnNameHandler(c *sysadmServer.Context) {
	var response apiutils.ApiResponseData
	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		response = apiutils.BuildResponseDataForError(7000170003, "您没有登录或者没有权限执行本操作")
		c.JSON(http.StatusOK, response)
		return
	}
	requestData, e := utils.NewGetRequestData(c, []string{"objvalue"})

	if e != nil || !validEnName(requestData["objvalue"]) {
		response = apiutils.BuildResponseDataForError(7000170004, "英文名称不能重复，且其长度不得大于255个字符.")
	} else {
		response = apiutils.BuildResponseDataForSuccess("ok")
	}

	c.JSON(http.StatusOK, response)
}
