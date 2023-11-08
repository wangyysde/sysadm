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
	datacenter "sysadm/datacenter/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/user"
)

func addformHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(7000130001, "debug", "try to display add object form page"))

	addTemplateFile := "addObjectForm.html"
	var emptyString []string
	baseUri := "/" + DefaultModuleName + "/"
	enctype := ""
	postUri := "/api/" + DefaultApiVersion + "/" + DefaultModuleName + "/add"
	var tplDataLines []objectsUI.ObjLineData

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000130002, errs, e)
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
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000130003, errs, e)
		return
	}
	// prepare datacenter select data
	dcOptions := buildDCOptions(dcList)
	dcSelect := objectsUI.ObjItemInfo{Title: "所属数据中心", ID: "dcid", Name: "dcid", Kind: "SELECT", ActionUri: "", ItemData: dcOptions}
	var lineItems []objectsUI.ObjItemInfo
	lineItems = append(lineItems, dcSelect)
	lineData := objectsUI.ObjLineData{Items: lineItems}
	tplDataLines = append(tplDataLines, lineData)

	cnName := objectsUI.ObjItemInfo{Title: "可用区名称", ID: "cnName", Name: "cnName", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validCNName", Note: "可用区中文名称，长度不大于255个字符且不能为空"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{cnName}})

	enName := objectsUI.ObjItemInfo{Title: "可用区英文名称", ID: "enName", Name: "enName", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validENName", Note: "集群的英文名称，长度不大于255个字符"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{enName}})

	dutyTel := objectsUI.ObjItemInfo{Title: "值班电话", ID: "dutyTel", Name: "dutyTel", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validDutyTel", Note: "值班电话"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{dutyTel}})

	remark := objectsUI.ObjItemInfo{Title: "备注", ID: "remark", Name: "remark", Kind: "TEXTAREA", Size: 40, Rows: 5}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{remark}})

	tplData, _ := objectsUI.InitAddObjectFormTemplateData(baseUri, "基础设施", "添加可用区",
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
