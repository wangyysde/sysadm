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

// 注意：
// 1.命令参数获取方法中从对象字段获取和通过另一个命令获取实现比较比较复杂，暂不支持这两种方法
// 2.命令的依赖性暂不实现
// 3. 对应的命令的事务性暂不实现
package app

import (
	"fmt"
	"github.com/wangyysde/sysadmServer"
	"net/http"
	"strconv"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	sysadmOS "sysadm/os/app"
	"sysadm/sysadmLog"
	"sysadm/sysadmapi/apiutils"
	"sysadm/user"
	"sysadm/utils"
	sysadmVersion "sysadm/version/app"
)

func addformHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(7000140001, "debug", "try to display add object form page"))

	addTemplateFile := "addObjectFormNew.html"
	var emptyString []string
	baseUri := "/" + DefaultModuleName + "/"
	enctype := ""
	postUri := "/api/" + DefaultApiVersion + "/" + DefaultModuleName + "/add"
	var tplDataLines []objectsUI.ObjLineData

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000140002, errs, e)
		return
	}

	// preparing select data
	var osEntity sysadmObjects.ObjectEntity
	osEntity = sysadmOS.New()
	conditions := make(map[string]string, 0)
	order := make(map[string]string, 0)
	osList, e := osEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000140003, errs, e)
		return
	}
	// prepare OS select data
	osOptions := buildOSOptions(osList)
	osSelect := objectsUI.ObjItemInfo{Title: "适用的操作系统", ID: "osID", Name: "osID", Kind: "SELECT",
		ActionUri: "getversionbyosidforselect", ItemData: osOptions, SubObjID: "osversionid", JsActionKind: objectsUI.JsActionKind_Select_Change_SelectOptions}
	var osItems []objectsUI.ObjItemInfo
	osItems = append(osItems, osSelect)

	// prepare empty version select data
	var verOptions []objectsUI.SubItems
	verOption := objectsUI.SubItems{Value: "0", Text: "===选择所适用的版本===", Checked: true}
	verOptions = append(verOptions, verOption)
	verSelect := objectsUI.ObjItemInfo{Title: "所适用版本", ID: "osversionid", Name: "osversionid", Kind: "SELECT", ActionUri: "", ItemData: verOptions}
	osItems = append(osItems, verSelect)

	lineData := objectsUI.ObjLineData{Items: osItems}
	tplDataLines = append(tplDataLines, lineData)

	name := objectsUI.ObjItemInfo{Title: "命令名称", ID: "name", Name: "name", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validName", Note: "命令的名称，用于在前端显示时使用的"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{name}})

	command := objectsUI.ObjItemInfo{Title: "命令", ID: "command", Name: "command", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validCommand", Note: "需执行的命令，如果是内嵌命令则是内嵌命令标识，如果是脚本或系统系统则是命令的绝对路径，不能为空，且不能有重复"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{command}})

	var executionTypeOptions []objectsUI.SubItems
	subItem := objectsUI.SubItems{Value: strconv.Itoa(int(ExecutionTypeAuto)), Text: "自动执行", Checked: false, RelatedObjectIsDisplay: true}
	executionTypeOptions = append(executionTypeOptions, subItem)
	subItem = objectsUI.SubItems{Value: strconv.Itoa(int(ExecutionTypeHand)), Text: "手动执行", Checked: true, RelatedObjectIsDisplay: false}
	executionTypeOptions = append(executionTypeOptions, subItem)
	executionType := objectsUI.ObjItemInfo{Title: "执行类型", ID: "executionType", Name: "executionType", Kind: "RADIO", ActionUri: "", Note: "", ItemData: executionTypeOptions, SubObjID: "automationKind", JsActionKind: objectsUI.JsActionKind_Radio_ChangeSubDisplay}
	var executionTypeItems []objectsUI.ObjItemInfo
	executionTypeItems = append(executionTypeItems, executionType)

	autoKindItems := buildAutoKindItems()
	automationKind := objectsUI.ObjItemInfo{Title: "执行时刻", ID: "automationKind", Name: "automationKind", Kind: "RADIO", ActionUri: "", Note: "", ItemData: autoKindItems, JsActionKind: objectsUI.JsActionKind_Radio_CustomizeAction, NoDisplay: true}
	executionTypeItems = append(executionTypeItems, automationKind)
	lineData = objectsUI.ObjLineData{Items: executionTypeItems}
	tplDataLines = append(tplDataLines, lineData)

	objectInfo, e := sysadmObjects.GetCommandRelatedObjectList()
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000140004, errs, e)
		return
	}

	// prepare object information select data
	objOptions, e := buildObjectInfoOptions(objectInfo)
	objSelect := objectsUI.ObjItemInfo{Title: "对象", ID: "objectName", Name: "objectName", Kind: "SELECT",
		ActionUri: "", ItemData: objOptions, NoDisplay: true}
	var objSelectItems []objectsUI.ObjItemInfo
	objSelectItems = append(objSelectItems, objSelect)
	lineData = objectsUI.ObjLineData{Items: objSelectItems}
	tplDataLines = append(tplDataLines, lineData)

	crontab := objectsUI.ObjItemInfo{Title: "Crontab时间", ID: "crontab", Name: "crontab", Kind: "TEXT", ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validCrontab", Note: "如果命令是crontab类别的自动执行命令，本字段指定crontab格式的执行时间和周期", NoDisplay: true}
	var crontItems []objectsUI.ObjItemInfo
	crontItems = append(crontItems, crontab)
	lineData = objectsUI.ObjLineData{Items: crontItems}
	tplDataLines = append(tplDataLines, lineData)

	var synchronizedOptions []objectsUI.SubItems
	subItem = objectsUI.SubItems{Value: strconv.Itoa(CommandKindSynchronized), Text: "同步命令", Checked: true}
	synchronizedOptions = append(synchronizedOptions, subItem)
	subItem = objectsUI.SubItems{Value: strconv.Itoa(CommandKindAsynchronized), Text: "异步命令", Checked: false}
	synchronizedOptions = append(synchronizedOptions, subItem)
	synchronizedKind := objectsUI.ObjItemInfo{Title: "通讯类型", ID: "synchronized", Name: "synchronized", Kind: "RADIO", ActionUri: "", Note: "", ItemData: synchronizedOptions}
	var synchronizedItems []objectsUI.ObjItemInfo
	synchronizedItems = append(synchronizedItems, synchronizedKind)

	var typeOptions []objectsUI.SubItems
	subItem = objectsUI.SubItems{Value: strconv.Itoa(int(CommandTypeSys)), Text: "系统命令", Checked: true}
	typeOptions = append(typeOptions, subItem)
	subItem = objectsUI.SubItems{Value: strconv.Itoa(int(CommandTypeScript)), Text: "脚本命令", Checked: false}
	typeOptions = append(typeOptions, subItem)
	commandType := objectsUI.ObjItemInfo{Title: "命令类型", ID: "type", Name: "type", Kind: "RADIO", ActionUri: "", Note: "", ItemData: typeOptions}
	synchronizedItems = append(synchronizedItems, commandType)
	lineData = objectsUI.ObjLineData{Items: synchronizedItems}
	tplDataLines = append(tplDataLines, lineData)

	descriptions := objectsUI.ObjItemInfo{Title: "描述", ID: "descriptions", Name: "descriptions", Kind: "TEXTAREA", Size: 40, Rows: 5}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{descriptions}})

	addionalJs := []string{"/js/addCommand.js"}
	tplData, _ := objectsUI.InitAddObjectFormTemplateData(baseUri, "系统设置", "添加命令",
		enctype, postUri, "list", "list", addionalJs, emptyString)
	tplData["data"] = tplDataLines

	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, addTemplateFile, tplData)
}

func buildOSOptions(osList []interface{}) []objectsUI.SubItems {
	var ret []objectsUI.SubItems
	lineData := objectsUI.SubItems{Value: "0", Text: "===选择适用的操作系统===", Checked: true}
	ret = append(ret, lineData)

	for _, line := range osList {
		lineData := objectsUI.SubItems{}
		osData := line.(sysadmOS.OSSchema)
		lineData.Value = strconv.Itoa(osData.OSID)
		lineData.Text = osData.Name
		lineData.Checked = false
		ret = append(ret, lineData)
	}

	return ret
}

func buildAutoKindItems() []objectsUI.SubItems {
	var ret []objectsUI.SubItems
	subItem := objectsUI.SubItems{Value: strconv.Itoa(int(AutomationKindObjectCreate)), Text: "对象创建时", Checked: false}
	ret = append(ret, subItem)

	subItem = objectsUI.SubItems{Value: strconv.Itoa(int(AutomationKindObjectConfChange)), Text: "对象配置修改时", Checked: false}
	ret = append(ret, subItem)

	subItem = objectsUI.SubItems{Value: strconv.Itoa(int(AutomationKindObjectStatusChange)), Text: "对象状态改变时", Checked: false}
	ret = append(ret, subItem)

	subItem = objectsUI.SubItems{Value: strconv.Itoa(int(AutomationKindObjectDelete)), Text: "对象删除时", Checked: false}
	ret = append(ret, subItem)

	subItem = objectsUI.SubItems{Value: strconv.Itoa(int(AutomationKindCrontab)), Text: "定时任务", Checked: false}
	ret = append(ret, subItem)

	return ret
}

func buildObjectInfoOptions(objInfos []interface{}) ([]objectsUI.SubItems, error) {
	var ret []objectsUI.SubItems
	lineData := objectsUI.SubItems{Value: "0", Text: "===选择对象===", Checked: true}
	ret = append(ret, lineData)

	for _, line := range objInfos {
		lineData := objectsUI.SubItems{}
		objData, ok := line.(sysadmObjects.ObjectInfoSchema)
		if !ok {
			return ret, fmt.Errorf("data is not object information schema")
		}
		lineData.Value = strconv.Itoa(objData.ID)
		lineData.Text = objData.CnName
		lineData.Checked = false
		ret = append(ret, lineData)
	}

	return ret, nil
}

func getversionbyosidforselectHandler(c *sysadmServer.Context) {

	requestData, e := utils.NewGetRequestData(c, []string{"objID"})
	if e != nil || requestData["objID"] == "" {
		response := apiutils.BuildResponseDataForError(7000140005, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	// preparing version data
	var versionEntity sysadmObjects.ObjectEntity
	versionEntity = sysadmVersion.New()
	conditions := make(map[string]string, 0)
	conditions["typeID"] = "='" + strconv.Itoa(int(sysadmVersion.VersionTypeOS)) + "'"
	conditions["osid"] = "='" + requestData["objID"] + "'"
	var emptyString []string
	order := make(map[string]string, 0)
	versionList, e := versionEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		response := apiutils.BuildResponseDataForError(7000140006, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	msg := ""
	for _, line := range versionList {
		lineData, ok := line.(sysadmVersion.VersionSchema)
		if !ok {
			response := apiutils.BuildResponseDataForError(7000140007, "数据处理错误")
			c.JSON(http.StatusOK, response)
			return
		}
		lineStr := strconv.Itoa(lineData.VersionID) + ":" + lineData.Name
		if msg == "" {
			msg = lineStr
		} else {
			msg = msg + "," + lineStr
		}
	}

	response := apiutils.BuildResponseDataForSuccess(msg)
	c.JSON(http.StatusOK, response)
}
