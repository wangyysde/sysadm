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
	"context"
	"github.com/wangyysde/sysadmServer"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"sysadm/sysadmLog"
	"sysadm/sysadmapi/apiutils"
	"sysadm/user"
)

func validateNameHandler(c *sysadmServer.Context, module string) {
	switch module {
	case "namespace":
		// validate object name according RFC1123
		validateNameWith1123Label(c, module)
		return
	}
}

func validateNameWith1123Label(c *sysadmServer.Context, module string) {
	var errs []sysadmLog.Sysadmerror
	var response apiutils.ApiResponseData

	errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001100001, "debug", "now validate  %s name for adding", module))
	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		response = apiutils.BuildResponseDataForError(80001100002, "您没有登录或者没有权限执行本操作")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001100002, "info", "user has not login or not permission when validating namespace name"))
		runData.logEntity.LogErrors(errs)
		c.JSON(http.StatusOK, response)
		return
	}

	requestKeys := []string{"dcID", "clusterID", "namespace", "objValue"}
	requestData, e := getRequestData(c, requestKeys)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001100003, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response = apiutils.BuildResponseDataForError(80001100003, "系统出错，请稍后再试或者联系系统管理员")
		c.JSON(http.StatusOK, response)
		return
	}

	validateSlice := apimachineryvalidation.ValidateNamespaceName(requestData["objValue"], false)
	if len(validateSlice) > 0 {
		validateStr := "所输入的名称在以下方面不符合命名空间名称要求："
		for _, v := range validateSlice {
			validateStr = validateStr + "." + v
		}
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001100004, "info", "name(%s) of namespace is not valid", requestData["objValue"]))
		runData.logEntity.LogErrors(errs)
		response = apiutils.BuildResponseDataForError(80001100004, validateStr)
		c.JSON(http.StatusOK, response)
		return
	}

	clientSet, e := buildClientSetByClusterID(requestData["clusterID"])
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001100005, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response = apiutils.BuildResponseDataForError(80001100005, "系统内部错误，请稍后再试或者联系系统管理员")
		c.JSON(http.StatusOK, response)
		return
	}

	nsList, e := clientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001100006, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response = apiutils.BuildResponseDataForError(80001100006, "系统内部错误，请稍后再试或者联系系统管理员")
		c.JSON(http.StatusOK, response)
		return
	}

	found := false
	for _, item := range nsList.Items {
		if item.Name == requestData["objValue"] {
			found = true
			break
		}
	}

	if found {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001100007, "error", "namespace(%s) is exist in the cluster", requestData["objValue"]))
		runData.logEntity.LogErrors(errs)
		response = apiutils.BuildResponseDataForError(80001100007, "拟创建的命名空间在集群中已经存在")
		c.JSON(http.StatusOK, response)
		return
	}

	response = apiutils.BuildResponseDataForSuccess("ok")
	c.JSON(http.StatusOK, response)
	return
}
