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
	"fmt"
	"github.com/wangyysde/sysadmServer"
	"net/http"
	"strconv"
	"strings"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	sysadmOS "sysadm/os/app"
	"sysadm/sysadmLog"
	"sysadm/user"
	"sysadm/utils"
	sysadmVersion "sysadm/version/app"
)

// order fields data of cluster list page
var allOrderFields = map[string]string{"TD1": "name"}

// which field will be order if user has not selected
var defaultOrderField = "TD1"

// 1 for DESC 0 for ASC
var defaultOrderDirection = "1"

// all popmenu items defined Format:
// item name, action name, action method
var allPopMenuItems = []string{"查看详情,detail,GET,page", "编辑命令,edit,GET,page", "删除命令,del,POST,tip", "编辑参数,editparas,GET,page", "设置依赖,setdependent,GET,poppage", "置为主机级别命令,sethostscope,POST,tip", "置为集群级别命令,setclusterscope,POST,tip"}

// define all list items(cols) name
var allListItems = map[string]string{"TD1": "名称", "TD2": "操作系统/版本", "TD3": "执行类型", "TD4": "参数类型", "TD5": "同步类型", "TD6": "命令类型", "TD7": "命令"}

func listHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(7000120001, "debug", "now handling command list"))
	listTemplateFile := "objectlistNew.html"

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000120002, errs, e)
		return
	}

	// get request data
	requestData, e := getRequestData(c)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000120003, errs, e)
		return
	}

	// 初始化模板数据
	tplData, e := objectsUI.InitTemplateData("/"+defaultObjectName+"/", "系统设置", "命令列表", "添加命令定义", "yes",
		allPopMenuItems, []string{}, requestData)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000120004, errs, e)
		return
	}

	// preparing select data
	var osEntity sysadmObjects.ObjectEntity
	osEntity = sysadmOS.New()
	conditions := make(map[string]string, 0)
	order := make(map[string]string, 0)
	var emptyString []string
	osList, e := osEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000120005, errs, e)
		return
	}

	var versionEntity sysadmObjects.ObjectEntity
	versionEntity = sysadmVersion.New()
	conditions["typeID"] = "='" + strconv.Itoa(int(sysadmVersion.VersionTypeOS)) + "'"
	versionList, e := versionEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000120006, errs, e)
		return
	}

	if e := buildSelectData(tplData, osList, versionList, requestData); e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000120007, errs, e)
		return
	}

	// build table header for list objects
	objectsUI.BuildThData(requestData, allOrderFields, allListItems, tplData, defaultOrderField, defaultOrderDirection)

	searchContent := objectsUI.GetSearchContentFromRequest(requestData)
	ids := objectsUI.GetObjectIdsFromRequest(requestData)
	searchKeys := []string{"command", "name"}
	startPos := objectsUI.GetStartPosFromRequest(requestData)
	commandConditions := make(map[string]string, 0)
	commandConditions["deprecated"] = "='0'"
	if requestData["groupSelectID"] != "" && requestData["groupSelectID"] != "0" {
		commandConditions["osversionid"] = "='" + requestData["groupSelectID"] + "'"
	}

	// get total number of list objects
	var commandEntity sysadmObjects.ObjectEntity
	commandEntity, e = New(runData.dbConf, runData.workingRoot)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000120008, errs, e)
		return
	}
	commandCount, e := commandEntity.GetObjectCount(searchContent, ids, searchKeys, commandConditions)
	if e != nil || commandCount < 1 {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000120009, errs, e)
		return
	}

	// get list data
	requestOrder := objectsUI.BuildOrderDataForQuery(requestData, allOrderFields, defaultOrderField, defaultOrderDirection)
	commandList, e := commandEntity.GetObjectList(searchContent, ids, searchKeys, commandConditions, startPos, runData.pageInfo.NumPerPage, requestOrder)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000120010, errs, e)
		return
	}
	// prepare cluster list data
	objListData, e := prepareObjectData(osList, versionList, commandList)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000120011, errs, e)
		return
	}
	tplData["objListData"] = objListData

	// prepare page number information
	objectsUI.BuildPageNumInfo(tplData, requestData, commandCount, startPos, runData.pageInfo.NumPerPage, defaultOrderField, defaultOrderDirection)

	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, listTemplateFile, tplData)
}

func getRequestData(c *sysadmServer.Context) (map[string]string, error) {
	requestData, e := utils.NewGetRequestData(c, []string{"groupSelectID", "searchContent", "start", "orderfield", "direction"})
	if e != nil {
		return requestData, e
	}

	objectIds := ""
	objectIDMap, _ := utils.GetRequestDataArray(c, []string{"objectid[]"})
	if objectIDMap != nil {
		objectIDSlice, ok := objectIDMap["objectid[]"]
		if ok {
			objectIds = strings.Join(objectIDSlice, ",")
		}
	}
	requestData["objectIds"] = objectIds
	if strings.TrimSpace(requestData["start"]) == "" {
		requestData["start"] = "0"
	}

	return requestData, nil
}

func buildSelectData(tplData map[string]interface{}, osList, versionList []interface{}, requestData map[string]string) error {
	secondOptions := make(map[string]string, 0)
	groupSelectID := requestData["groupSelectID"]
	secondSelectedID := "0"
	firstSelectedID := "0"
	var secondSelectedOptions []objectsUI.SelectOption
	for _, line := range versionList {
		verData, ok := line.(sysadmVersion.VersionSchema)
		if !ok {
			return fmt.Errorf("the data is not version schema")
		}
		versionID := verData.VersionID
		name := strings.TrimSpace(verData.Name)
		osID := verData.OSID
		osIDStr := strings.TrimSpace(strconv.Itoa(osID))
		versionIDStr := strings.TrimSpace(strconv.Itoa(versionID))

		if versionIDStr == groupSelectID {
			secondSelectedID = versionIDStr
			firstSelectedID = osIDStr

		}
		if firstSelectedID == osIDStr {
			selectedOption := objectsUI.SelectOption{
				Id:       strconv.Itoa(versionID),
				Text:     name,
				ParentID: osIDStr,
			}
			secondSelectedOptions = append(secondSelectedOptions, selectedOption)
		}

		subOption := "['" + versionIDStr + "','" + name + "']"
		addOption, ok := secondOptions[osIDStr]
		if ok {
			addOption = addOption + "," + subOption
		} else {
			addOption = subOption
		}
		secondOptions[osIDStr] = addOption
	}
	secondSelect := objectsUI.SelectData{Title: "版本", SelectedId: secondSelectedID, SelectedOptions: secondSelectedOptions}
	var secondOptionList []objectsUI.SelectOption
	for code, value := range secondOptions {
		option := objectsUI.SelectOption{
			ParentID:    code,
			OptionsList: value,
		}
		secondOptionList = append(secondOptionList, option)
	}
	secondSelect.Options = secondOptionList

	var firstOptions []objectsUI.SelectOption
	firstSelect := objectsUI.SelectData{
		Title:      "操作系统",
		SelectedId: firstSelectedID,
	}
	for _, line := range osList {
		osData, ok := line.(sysadmOS.OSSchema)
		if !ok {
			return fmt.Errorf("the data is not OS schema")
		}
		option := objectsUI.SelectOption{
			Id:       strconv.Itoa(osData.OSID),
			Text:     osData.Name,
			ParentID: "0",
		}
		firstOptions = append(firstOptions, option)
	}
	firstSelect.Options = firstOptions

	firstSelect.SelectedId = firstSelectedID
	secondSelect.SelectedId = secondSelectedID
	secondSelect.SelectedOptions = secondSelectedOptions

	thirdSelect := objectsUI.SelectData{}
	if e := objectsUI.BuildMultiSelectData(firstSelect, secondSelect, thirdSelect, tplData); e != nil {
		return e
	}

	return nil
}

func prepareObjectData(osList, versionList, commandList []interface{}) ([]map[string]interface{}, error) {
	var dataList []map[string]interface{}

	for _, line := range commandList {
		commandData, ok := line.(CommandDefinedSchema)
		if !ok {
			return dataList, fmt.Errorf("data is not command defined schema")
		}

		lineMap := make(map[string]interface{}, 0)
		lineMap["TD1"] = commandData.Name
		osVerStr := getOSName(osList, commandData.OSID) + "/" + getVerName(versionList, commandData.OsVersionID)
		lineMap["TD2"] = osVerStr
		exeType := "自动执行"
		if commandData.ExecutionType == int(ExecutionTypeHand) {
			exeType = "手动执行"
		}
		lineMap["TD3"] = exeType

		popmenuitems := "0,1,2,3,4"
		paraKind := "非法类型"
		switch ParaKind(commandData.ParaKind) {
		case ParaKindNo:
			paraKind = "无参数"
		case ParaKindFixed:
			paraKind = "固定值"
		case ParaKindObjFieldValue:
			paraKind = "对象字段值"
		case ParaKindGetByCommand:
			paraKind = "通过另一个命令获取"
		}
		lineMap["TD4"] = paraKind
		synchronized := "同步"
		if commandData.Synchronized != 0 {
			synchronized = "异步"
		}
		lineMap["TD5"] = synchronized

		commandType := "非法值"
		switch CommandType(commandData.Type) {
		case CommandTypeBuiltin:
			commandType = "内嵌命令"
		case CommandTypeSys:
			commandType = "系统命令"
		case CommandTypeScript:
			commandType = "脚本或者批处理程序"
		}
		lineMap["TD6"] = commandType
		lineMap["TD7"] = commandData.Command

		if commandData.TransactionScope == int(TransationScopeHost) {
			popmenuitems = popmenuitems + "," + "6"
		} else {
			popmenuitems = popmenuitems + "," + "5"
		}

		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return dataList, nil
}

func getOSName(osList []interface{}, osid int) string {
	name := "未知"

	for _, line := range osList {
		osData, ok := line.(sysadmOS.OSSchema)
		if !ok {
			return name
		}

		if osData.OSID == osid {
			return osData.Name
		}
	}

	return name
}

func getVerName(verList []interface{}, verid int) string {
	name := "未知"
	for _, line := range verList {
		verData, ok := line.(sysadmVersion.VersionSchema)
		if !ok {
			return name
		}
		if verData.VersionID == verid {
			return verData.Name
		}
	}

	return name
}
