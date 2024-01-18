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
	coreV1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
	"sysadm/sysadmLog"
	"sysadm/sysadmapi/apiutils"
	"sysadm/user"
)

func getForApiHandlers(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		e := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400002, "json", "您没有登录或者没有权限执行本操作")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400002, "info", "user has not login or not permission"))
		if e != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400003, "error", "%s", e))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	module := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(strings.ToLower(c.Param("module"))), "/"), "/")
	action := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(c.Param("action")), "/"), "/")
	var objEntity objectEntity = nil
	for m, o := range modulesDefined {
		if m == module {
			o.setObjectInfo()
			objEntity = o
		}
	}

	if objEntity == nil {
		e := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400004, "json", "操作错误，请稍后再试或联系系统管理员")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400004, "error", "module %s was not found", module))
		if e != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400005, "error", "%s", e))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	if objEntity.getNamespaced() {
		getForApiWithNamespacedHandler(c, module, action)
		return
	}

	getForApiNonNamespacedHandler(c, module, action)
	return
}

func getForApiWithNamespacedHandler(c *sysadmServer.Context, module, action string) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400029, "info", "namespaced handler for module  %s name with action %s", module, action))
	requestKeys := []string{"clusterID", "namespace", "objValue"}
	requestData, e := getRequestData(c, requestKeys)
	if e != nil {
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400006, "json", "操作错误，请稍后再试或联系系统管理员")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400006, "error", "%s", e))
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400007, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	switch action {
	case "getNameList":
		getNamespacedResourceNameListHandler(c, module, action, requestData)
		return
	case "validateNewName":
		validateNamespacedResourceNewName(c, module, requestData)
	default:
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400008, "error", "action %s was not defined", action))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400008, "json", "操作错误，请稍后再试或联系系统管理员")
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400009, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)

		return
	}

	runData.logEntity.LogErrors(errs)
	return
}

func getNamespacedResourceNameListHandler(c *sysadmServer.Context, module, action string, requestData map[string]string) {
	var errs []sysadmLog.Sysadmerror

	clientSet, e := buildClientSetByClusterID(requestData["clusterID"])
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400010, "error", "%s", e))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400010, "json", "操作错误，请稍后再试或联系系统管理员")
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400011, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	var objList interface{} = nil
	var er error = nil
	switch module {
	case "pvc":
		objList, er = clientSet.CoreV1().PersistentVolumeClaims(requestData["namespace"]).List(context.Background(), metav1.ListOptions{})
	case "configmap":
		objList, er = clientSet.CoreV1().ConfigMaps(requestData["namespace"]).List(context.Background(), metav1.ListOptions{})
	case "secret":
		objList, er = clientSet.CoreV1().Secrets(requestData["namespace"]).List(context.Background(), metav1.ListOptions{})
	case "service":
		objList, er = clientSet.CoreV1().Services(requestData["namespace"]).List(context.Background(), metav1.ListOptions{})
	case "serviceaccount":
		objList, er = clientSet.CoreV1().ServiceAccounts(requestData["namespace"]).List(context.Background(), metav1.ListOptions{})
	}

	if er != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400012, "error", "%s", er))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400012, "json", "操作错误，请稍后再试或联系系统管理员")
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400013, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	var nameList = []string{}
	errorMsg := ""
	switch module {
	case "pvc":
		pvList, ok := objList.(*coreV1.PersistentVolumeClaimList)
		if !ok {
			errorMsg = "the data is not persistent volume Claim List Schema"
			break
		}
		for _, v := range pvList.Items {
			nameList = append(nameList, v.Name)
		}
	case "configmap":
		cmList, ok := objList.(*coreV1.ConfigMapList)
		if !ok {
			errorMsg = "the data is not configMap List Schema"
			break
		}

		for _, v := range cmList.Items {
			nameList = append(nameList, v.Name)
		}
	case "secret":
		secretList, ok := objList.(*coreV1.SecretList)
		if !ok {
			errorMsg = "the data is not secret List Schema"
			break
		}
		for _, v := range secretList.Items {
			nameList = append(nameList, v.Name)
		}
	case "service":
		serviceList, ok := objList.(*coreV1.ServiceList)
		if !ok {
			errorMsg = "the data is not service List Schema"
			break
		}
		for _, v := range serviceList.Items {
			nameList = append(nameList, v.Name)
		}
	case "serviceaccount":
		saList, ok := objList.(*coreV1.ServiceAccountList)
		if !ok {
			errorMsg = "the data is not service account List Schema"
			break
		}
		for _, v := range saList.Items {
			nameList = append(nameList, v.Name)
		}
	}

	if errorMsg != "" {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400014, "error", "%s", errorMsg))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400015, "json", "操作错误，请稍后再试或联系系统管理员")
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400016, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 0, "json", nameList)
	if e1 != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400017, "error", "%s", e1))
	}
	runData.logEntity.LogErrors(errs)
	return
}

func getForApiNonNamespacedHandler(c *sysadmServer.Context, module, action string) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400028, "info", "non namespaced handler for module  %s name with action %s", module, action))
	requestKeys := []string{"clusterID", "objValue"}
	requestData, e := getRequestData(c, requestKeys)
	if e != nil {
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400029, "json", "操作错误，请稍后再试或联系系统管理员")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400029, "error", "%s", e))
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400030, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	switch action {
	case "getNameList":
		getNonNamespacedResourceNameListHandler(c, module, requestData)
		return
	case "validateNewName":
		validateNonNamespacedResourceNewName(c, module, requestData)
	default:
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400031, "error", "action %s was not defined", action))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400031, "json", "操作错误，请稍后再试或联系系统管理员")
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400032, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)

		return
	}

	runData.logEntity.LogErrors(errs)
	return
}

func validateNamespacedResourceNewName(c *sysadmServer.Context, module string, requestData map[string]string) {
	var errs []sysadmLog.Sysadmerror
	if requestData["objValue"] == "" {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400018, "info", "validate new %s name with empty value", module))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400018, "json", (module + "名字不能为空"))
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400028, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	objValue := requestData["objValue"]
	validateSlice := []string{}
	switch module {
	case "deployment":
		validateSlice = apimachineryvalidation.NameIsDNSLabel(objValue, false)
	}

	if len(validateSlice) > 0 {
		invalidateStr := ""
		for _, v := range validateSlice {
			if invalidateStr == "" {
				invalidateStr = v
			} else {
				invalidateStr = invalidateStr + ";" + v
			}
		}
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400020, "debug", "new name %s  of deployment is not validate in: ", objValue, invalidateStr))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400020, "json", ("应用名称" + objValue + "不合法"))
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400021, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	clientSet, e := buildClientSetByClusterID(requestData["clusterID"])
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400022, "error", "%s", e))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400022, "json", "系统出现未知错误，请稍后再试或联系系统管理员")
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400023, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	existFlag := false
	ns := requestData["namespace"]
	switch module {
	case "deployment":
		objList, e := clientSet.AppsV1().Deployments(ns).List(context.Background(), metav1.ListOptions{})
		if e != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400024, "error", "%s", e))
			e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400023, "json", "系统出现未知错误，请稍后再试或联系系统管理员")
			if e1 != nil {
				errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400024, "error", "%s", e1))
			}
			runData.logEntity.LogErrors(errs)
			return
		}

		for _, o := range objList.Items {
			if o.Name == objValue {
				existFlag = true
				break
			}
		}
	}

	if existFlag {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400025, "debug", "name %s exist in namespace %s", objValue, ns))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400025, "json", ("应用名称" + objValue + "在命名空间" + ns + "已经存在"))
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400026, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 0, "json", "ok")
	if e1 != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400027, "error", "%s", e1))
	}
	runData.logEntity.LogErrors(errs)
	return
}

func getNonNamespacedResourceNameListHandler(c *sysadmServer.Context, module string, requestData map[string]string) {
	var errs []sysadmLog.Sysadmerror

	clientSet, e := buildClientSetByClusterID(requestData["clusterID"])
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400033, "error", "%s", e))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400033, "json", "操作错误，请稍后再试或联系系统管理员")
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400034, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	var objList interface{} = nil
	var er error = nil
	switch module {
	case "ingressclass":
		objList, er = clientSet.NetworkingV1().IngressClasses().List(context.Background(), metav1.ListOptions{})
	default:
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400035, "error", "action of getNameList for module %s was not defined", module))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400035, "json", "操作错误，请稍后再试或联系系统管理员")
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400036, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}
	if er != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400037, "error", "%s", er))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400037, "json", "操作错误，请稍后再试或联系系统管理员")
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400038, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	var nameList = []string{}
	errorMsg := ""
	switch module {
	case "ingressclass":
		ingressClassList, ok := objList.(*networkingV1.IngressClassList)
		if !ok {
			errorMsg = "the data is not ingress Classes List Schema"
			break
		}
		for _, v := range ingressClassList.Items {
			nameList = append(nameList, v.Name)
		}
	}

	if errorMsg != "" {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400039, "error", "%s", errorMsg))
		e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 80001400039, "json", "操作错误，请稍后再试或联系系统管理员")
		if e1 != nil {
			errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400040, "error", "%s", e1))
		}
		runData.logEntity.LogErrors(errs)
		return
	}

	e1 := apiutils.ResponseDataToClient(c, nil, http.StatusOK, 0, "json", nameList)
	if e1 != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(80001400041, "error", "%s", e1))
	}
	runData.logEntity.LogErrors(errs)
	return
}

func validateNonNamespacedResourceNewName(c *sysadmServer.Context, module string, requestData map[string]string) {
	// var errs []sysadmLog.Sysadmerror

	//TODO

	return
}
