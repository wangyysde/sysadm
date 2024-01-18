/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 202d4 Bzhy Network. All rights reserved.
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
	"mime/multipart"
	"net/http"
	"strings"
	"sysadm/k8sclient"
	"sysadm/sysadmLog"
	"sysadm/sysadmapi/apiutils"
	"sysadm/user"
	"sysadm/utils"
)

func postForApiHandlers(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		e := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001500002, "json", "您没有登录或者没有权限执行本操作")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500002, "info", "user has not login or not permission"))
		if e != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500003, "error", "%s", e))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	module := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(strings.ToLower(c.Param("module"))), "/"), "/")
	action := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(c.Param("action")), "/"), "/")
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500004, "info", "handler for module  %s name with action %s", module, action))
	var objEntity objectEntity = nil
	for m, o := range modulesDefined {
		if m == module {
			o.setObjectInfo()
			objEntity = o
		}
	}

	if objEntity == nil {
		e := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001500005, "json", "操作错误，请稍后再试或联系系统管理员")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500005, "error", "module %s was not found", module))
		if e != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500006, "error", "%s", e))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	switch action {
	case "add":
		postResourceAddHandler(c, module, action)
	default:
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500009, "error", "action %s was not defined", action))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001500009, "json", "操作错误，请稍后再试或联系系统管理员")
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500010, "error", "%s", e1))
		}
	}

	runData.logEntity.LogErrors(errs)

	return
}

func postResourceAddHandler(c *sysadmServer.Context, module, action string) {
	var errs []sysadmLog.Sysadmerror
	formData, e := utils.GetMultipartData(c, []string{"dcid", "clusterID", "namespace", "addType", "objContent", "objFile"})
	if e != nil {
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001500011, "json", "操作错误，请稍后再试或联系系统管理员")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500011, "error", "%s", e))
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500012, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	addTypeSlice := formData["addType"].([]string)
	addType := strings.TrimSpace(addTypeSlice[0])
	if addType != "0" && addType != "1" && addType != "2" {
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001500013, "json", "操作错误，请稍后再试或联系系统管理员")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500013, "error", "the type %s is invalid for adding %s ", addType, module))
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500014, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	if addType == "2" {
		postResourceAdd(c, module, action)
		return
	}

	yamlContent := ""
	if addType == "0" {
		objContentSlice := formData["objContent"].([]string)
		yamlContent = objContentSlice[0]
	}
	if addType == "1" {
		yamlByte, e := utils.ReadUploadedFile(formData["objFile"].(*multipart.FileHeader))
		if e != nil {
			e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001500015, "json", "上传Yaml文件出错，请稍后再试或联系系统管理员")
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500015, "error", "%s ", e))
			if e1 != nil {
				errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500016, "error", "%s", e1))
			}
			runData.logEntity.LogErrors(errs)
			return
		}
		yamlContent = utils.Interface2String(yamlByte)
	}

	clusterIDSlice := formData["clusterID"].([]string)
	clusterID := clusterIDSlice[0]
	clientSet, e := buildClientSetByClusterID(clusterID)
	if e != nil {
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001500016, "json", "连接到指定集群出错，请稍后再试或联系系统管理员")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500016, "error", "%s ", e))
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500017, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	e = k8sclient.ApplyFromYamlByClientSet(yamlContent, clientSet)
	if e != nil {
		eStr := fmt.Sprintf("新建%s失败:%s", module, e)
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001500018, "json", eStr)
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500018, "error", "%s ", e))
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500019, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 0, "json", "ok")
	if e1 != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500020, "error", "%s", e1))
	}
	runData.logEntity.LogErrors(errs)
	return
}

func postResourceAdd(c *sysadmServer.Context, module, action string) {
	var errs []sysadmLog.Sysadmerror

	objEntity, e := newObjEntity(module)
	if e != nil {
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001500021, "json", "操作错误，请稍后再试或联系系统管理员")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500021, "error", "%s ", e))
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500022, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}
	objEntity.setObjectInfo()

	e = objEntity.addNewResource(c, module)
	if e != nil {
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001500023, "json", "资源添加错误，请稍后再试或联系系统管理员")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500023, "error", "%s ", e))
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500024, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 0, "json", "ok")
	if e1 != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001500025, "error", "%s", e1))
	}
	runData.logEntity.LogErrors(errs)
	return
}
