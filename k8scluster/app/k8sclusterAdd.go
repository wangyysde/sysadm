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
	"strconv"
	az "sysadm/availablezone/app"
	datacenter "sysadm/datacenter/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/sysadmapi/apiutils"
	"sysadm/user"
	"sysadm/utils"
)

func addformHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(700060001, "debug", "try to display add object form page"))

	messageTemplateFile := "showmessage.html"
	addTemplateFile := "addObjectForm.html"
	var emptyString []string
	baseUri := "/" + DefaultModuleName + "/"
	enctype := `multipart/form-data`
	postUri := "/api/" + DefaultApiVersion + "/" + DefaultModuleName + "/add"
	var tplDataLines []objectsUI.ObjLineData

	messageTplData := make(map[string]interface{}, 0)
	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700060002, "error", "user should login %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "您没有登录或者您没有权限添加集群信息"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	// preparing datacenter data
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	conditions := make(map[string]string, 0)
	conditions["isDeleted"] = "=0"
	order := make(map[string]string, 0)
	dcList, e := dcEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700060003, "error", "get datacenter data error %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	// prepare datacenter select data
	dcOptions := buildDCOptions(dcList)
	dcSelect := objectsUI.ObjItemInfo{Title: "所属数据中心", ID: "dcid", Name: "dcid", Kind: "SELECT",
		ActionUri: "getazbydcidforselect", ItemData: dcOptions, SubObjID: "azid"}
	var lineItems []objectsUI.ObjItemInfo
	lineItems = append(lineItems, dcSelect)

	// prepare empty az select data
	var azOptions []objectsUI.SubItems
	azOption := objectsUI.SubItems{Value: "0", Text: "===选择集群所属可用区===", Checked: true}
	azOptions = append(azOptions, azOption)
	azSelect := objectsUI.ObjItemInfo{Title: "所属可用区", ID: "azid", Name: "azid", Kind: "SELECT", ActionUri: "", ItemData: azOptions}
	lineItems = append(lineItems, azSelect)

	lineData := objectsUI.ObjLineData{Items: lineItems}
	tplDataLines = append(tplDataLines, lineData)

	cnName := objectsUI.ObjItemInfo{Title: "集群中文名称", ID: "cnName", Name: "cnName", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validCNName", Note: "集群的中文名称，长度不大于255个字符"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{cnName}})

	enName := objectsUI.ObjItemInfo{Title: "集群英文名称", ID: "enName", Name: "enName", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validENName", Note: "集群的英文名称，长度不大于255个字符"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{enName}})

	apiserver := objectsUI.ObjItemInfo{Title: "kube-apiserver地址和端口", ID: "apiserver", Name: "apiserver", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validApiserverAddress", Note: "连接集群的kube-apiserver的地址和端口，如x.x.x.x:6443"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{apiserver}})

	clusterUser := objectsUI.ObjItemInfo{Title: "连接集群的用户名", ID: "clusterUser", Name: "clusterUser", Kind: "TEXT", Size: 30, DefaultValue: "admin", ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validClusterUser", Note: "连接集群的用户名，默认是admin"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{clusterUser}})

	ca := objectsUI.ObjItemInfo{Title: "CA证书", ID: "ca", Name: "ca", Kind: "FILE", Size: 30, Note: "连接集群的CA证书"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{ca}})

	cert := objectsUI.ObjItemInfo{Title: "证书", ID: "cert", Name: "cert", Kind: "FILE", Size: 30, Note: "连接集群的证书"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{cert}})

	key := objectsUI.ObjItemInfo{Title: "密钥", ID: "key", Name: "key", Kind: "FILE", Size: 30, Note: "连接集群的密钥"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{key}})

	dutyTel := objectsUI.ObjItemInfo{Title: "值班电话", ID: "dutyTel", Name: "dutyTel", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validDutyTel", Note: "值班电话"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{dutyTel}})

	remark := objectsUI.ObjItemInfo{Title: "备注", ID: "remark", Name: "remark", Kind: "TEXTAREA", Size: 40, Rows: 5}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{remark}})

	tplData, _ := objectsUI.InitAddObjectFormTemplateData(baseUri, "集群管理", "添加集群",
		enctype, postUri, "list", "list", emptyString, emptyString)
	tplData["data"] = tplDataLines

	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, addTemplateFile, tplData)
}

func buildDCOptions(dcList []interface{}) []objectsUI.SubItems {
	var ret []objectsUI.SubItems
	lineData := objectsUI.SubItems{Value: "0", Text: "===选择集群所属数据中心===", Checked: true}
	ret = append(ret, lineData)

	for _, line := range dcList {
		lineData := objectsUI.SubItems{}
		dcData := line.(datacenter.DatacenterSchema)
		lineData.Value = strconv.Itoa(int(dcData.Id))
		lineData.Text = dcData.CnName
		lineData.Checked = false
		ret = append(ret, lineData)
	}

	return ret
}

func getazbydcidforselectHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror

	requestData, e := utils.NewGetRequestData(c, []string{"objID"})
	if e != nil || requestData["objID"] == "" {
		response := apiutils.BuildResponseDataForError(700060004, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	// preparing datacenter data
	azEntity := az.New()
	conditions := make(map[string]string, 0)
	conditions["isDeleted"] = "=0"
	azList, e := azEntity.GetObjectListByDCID(requestData["objID"], conditions)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700060004, "error", "get datacenter data error %s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(700060004, "获取数据失败")
		c.JSON(http.StatusOK, response)
		return
	}

	msg := ""
	for _, line := range azList {
		linAZData := line.(az.AvailablezoneSchema)
		lineStr := strconv.Itoa(int(linAZData.Id)) + ":" + linAZData.CnName
		if msg == "" {
			msg = lineStr
		} else {
			msg = msg + "," + lineStr
		}
	}

	response := apiutils.BuildResponseDataForSuccess(msg)
	c.JSON(http.StatusOK, response)

}
