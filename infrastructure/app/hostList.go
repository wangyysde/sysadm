/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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
* Note this file for list host information
 */

package app

import (
	"fmt"
	"github.com/wangyysde/sysadmServer"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	sysadmAZObj "sysadm/availablezone/app"
	datacenter "sysadm/datacenter/app"
	"sysadm/db"
	"sysadm/httpclient"
	sysadmK8sCluster "sysadm/k8scluster/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/sysadmapi/apiutils"
	"sysadm/sysadmerror"
	"sysadm/user"
	"sysadm/utils"
)

/*
handler for handling list of the infrastructure
Query parameters of request are below:
conditionKey: key name for DB query ,such as hostid, userid,hostname....
conditionValue: the value of the conditionKey.for hostid, userid,hostname using =, for name, comment using like.
deleted: 0 :normarl 1: deleted
start: start number of the result will be returned.
num: lines of the result will be returned.
*/
func listHost(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(700040001, "debug", "now handling infrastructure list"))

	// get userid
	userid, e := user.GetSessionValue(c, "userid", WorkingData.sessionOption.sessionName)
	if e != nil || userid == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700040002, "error", "user should login %s", e))
		logErrors(errs)
		tplData := map[string]interface{}{
			"errormessage": "user should login",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, tplData)
		return
	}

	apiServer := WorkingData.apiServer
	tplData := make(map[string]interface{}, 0)
	tplData["userid"] = userid
	tplData, e = getYumList(c, apiServer, "showmessage.html", tplData)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700040003, "error", "get error yum list error %s", e))
		logErrors(errs)
		tplData := map[string]interface{}{
			"errormessage": "系统出错，请稍后再试或联系系统管理员",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, tplData)
		return
	}

	requestKeys := []string{"dcID", "azID", "clusterID", "userid", "start", "orderfield", "direction", "searchContent", "objID"}
	requestData, e := objectsUI.GetRequestData(c, requestKeys)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700040004, "error", "get request data error %s", e))
		logErrors(errs)
		errorTplData := map[string]interface{}{
			"errormessage": "系统出错，请稍后再试或联系系统管理员",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, errorTplData)
		return
	}

	// 为前端下拉菜单的数据中心部分准备数据
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	conditions := make(map[string]string, 0)
	conditions["isDeleted"] = "=0"
	order := make(map[string]string, 0)
	var emptyString []string
	dcList, e := dcEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700040005, "error", "get datacenter data error %s", e))
		logErrors(errs)
		errorTplData := map[string]interface{}{
			"errormessage": "系统出错，请稍后再试或联系系统管理员",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, errorTplData)
		return
	}
	e = buildSelectData(tplData, dcList, requestData)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700040006, "error", "building select data  error %s", e))
		logErrors(errs)
		errorTplData := map[string]interface{}{
			"errormessage": "系统出错，请稍后再试或联系系统管理员",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, errorTplData)
		return
	}

	listConditions, e := listHostCondition(requestData, false)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700040006, "error", "build query conditions error %s", e))
		logErrors(errs)
		errorTplData := map[string]interface{}{
			"errormessage": "系统出错，请稍后再试或联系系统管理员",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, errorTplData)
		return
	}

	// get total number of hosts
	hostCount, err, ok := getHostCountFromDB(listConditions)
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		errorTplData := map[string]interface{}{
			"errormessage": "there is an error occurred when got host total number",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, errorTplData)
		return
	}

	if hostCount < 1 {
		logErrors(errs)
		errorTplData := map[string]interface{}{
			"errormessage": "系统中没有查询到主机信息",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, errorTplData)
		return
	}

	startPos := objectsUI.GetStartPosFromRequest(requestData)
	// prepare page number information
	objectsUI.BuildPageNumInfoForWorkloadList(tplData, requestData, hostCount, startPos, WorkingData.pageInfo.numPerPage, "TD1", "1")

	// get host list data
	hostData, err, ok := getHostInfoListFromDB(requestData, listConditions, startPos)
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		errorTplData := map[string]interface{}{
			"errormessage": "there is an error occurred when got host list",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, errorTplData)

		return
	}

	// get os name and version name of os
	tplData["HostData"] = hostData
	tplData, err, ok = buildOsAndVersionInfo(tplData)
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		errorTplData := map[string]interface{}{
			"errormessage": "get os and its os version information error",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, errorTplData)
		return
	}

	// get host status text
	tplData = buildHostStatusInfo(tplData)

	// get foot page information
	pageInfoPara := buildPageInfoParas(requestData)
	tplData, err, ok = buildFootPageInfo(requestData, tplData, hostCount, pageInfoPara)
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		errorTplData := map[string]interface{}{
			"errormessage": "can not build foot page information",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, errorTplData)

		return
	}

	templateFile := "hostlist.html"
	c.HTML(http.StatusOK, templateFile, tplData)
}

func getYumList(c *sysadmServer.Context, apiServer *ApiServer, msgTpl string, tplData map[string]interface{}) (map[string]interface{}, error) {
	var errs []sysadmerror.Sysadmerror

	// get yum information list
	moduleName := "yum"
	actionName := "yumlist"
	apiServerData := apiutils.BuildApiServerData(moduleName, actionName, apiServer.ApiVersion, apiServer.Tls.IsTls,
		apiServer.Server.Address, apiServer.Server.Port, apiServer.Tls.Ca, apiServer.Tls.Cert, apiServer.Tls.Key)
	urlRaw, _ := apiutils.BuildApiUrl(apiServerData)

	requestParas := httpclient.RequestParams{}
	requestParas.Url = urlRaw
	requestParas.Method = http.MethodPost

	body, err := httpclient.SendRequest(&requestParas)
	errs = append(errs, err...)
	ret, err := apiutils.ParseResponseBody(body)
	errs = append(errs, err...)
	if !ret.Status {
		message := ret.Message
		messageLine := message[0]
		msg := utils.Interface2String(messageLine["msg"])
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030003, "error", "get yum information list error %s", msg))
		logErrors(errs)
		tplData := map[string]interface{}{
			"errormessage": msg,
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, tplData)
		return tplData, fmt.Errorf("%s", msg)
	}

	osVerList := make(map[int]string, 0)
	osVerMap := make(map[int]map[int]string, 0)
	osNameList := make(map[int]string, 0)
	yumInfoList := make(map[int]map[int]string, 0)
	for _, line := range ret.Message {
		yumID, errYumID := utils.Interface2Int(line["yumid"])
		yumName := utils.Interface2String(line["name"])
		osID, errOSID := utils.Interface2Int(line["osid"])
		osName := utils.Interface2String(line["osName"])
		yumTypeName := utils.Interface2String(line["typeName"])
		yumCatalog := utils.Interface2String(line["catalog"])
		versionID, errVersionID := utils.Interface2Int(line["versionid"])
		versionName := utils.Interface2String(line["versionName"])
		if errYumID != nil || yumName == "" || errOSID != nil || errVersionID != nil {
			errMsg := fmt.Sprintf("get yum information error: %s %s %s", errYumID, errOSID, errVersionID)
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030004, "error", errMsg))
			logErrors(errs)
			tplData := map[string]interface{}{
				"errormessage": errMsg,
			}
			templateFile := msgTpl // "showmessage.html"
			c.HTML(http.StatusOK, templateFile, tplData)
			return tplData, fmt.Errorf("%s", errMsg)
		}

		var yumInfoListStr string = ""
		if osMap, ok := yumInfoList[osID]; ok {
			if versionMap, ok := osMap[versionID]; ok {
				yumInfo := versionMap[:(len(versionMap) - 1)]
				yumInfoListStr = yumInfo + ",[" + strconv.Itoa(yumID) + ",'" + yumName + "','" + yumTypeName + "','" + yumCatalog + "']]"
			} else {
				yumInfoListStr = "[[" + strconv.Itoa(yumID) + ",'" + yumName + "','" + yumTypeName + "','" + yumCatalog + "']]"
			}
		} else {
			yumInfoListStr = "[[" + strconv.Itoa(yumID) + ",'" + yumName + "','" + yumTypeName + "','" + yumCatalog + "']]"
		}

		if verMap, ok := yumInfoList[osID]; ok {
			verMap[versionID] = yumInfoListStr
			yumInfoList[osID] = verMap
		} else {
			verMap := make(map[int]string, 0)
			verMap[versionID] = yumInfoListStr
			yumInfoList[osID] = verMap
		}

		if verMap, ok := osVerMap[osID]; ok {
			verMap[versionID] = versionName
			osVerMap[osID] = verMap
		} else {
			verMap := make(map[int]string, 0)
			verMap[versionID] = versionName
			osVerMap[osID] = verMap
		}
		osNameList[osID] = osName
	}
	tplData["osNameList"] = osNameList
	tplData["osVerMap"] = osVerMap

	for osID, verMap := range osVerMap {
		var osVer string = "["
		i := 0
		for verID, verName := range verMap {
			if i == 0 {
				osVer = osVer + "[" + strconv.Itoa(verID) + ",'" + verName + "']"
			} else {
				osVer = osVer + ",[" + strconv.Itoa(verID) + ",'" + verName + "']"
			}
			i = i + 1
		}
		osVer = osVer + "]"
		osVerList[osID] = osVer
	}

	var osList []map[string]string
	var javascriptYumStr []template.JS
	javascriptYumStr = append(javascriptYumStr, template.JS("var yumList = new Array();"))
	for osid, verMap := range yumInfoList {
		javascriptYumStr = append(javascriptYumStr, template.JS(fmt.Sprintf("yumList[%d] = new Array();", osid)))
		for verID, yumInfo := range verMap {
			javascriptYumStr = append(javascriptYumStr, template.JS(fmt.Sprintf("yumList[%d][%d] = %s;", osid, verID, yumInfo)))
		}
		osLine := make(map[string]string, 0)
		osLine["osid"] = strconv.Itoa(osid)
		osLine["osname"] = osNameList[osid]
		osList = append(osList, osLine)
	}

	tplData["yumList"] = javascriptYumStr
	var javascriptOsVersionStr []template.JS
	for osid, verStr := range osVerList {
		javascriptOsVersionStr = append(javascriptOsVersionStr, template.JS(fmt.Sprintf("osVerList[%d] = %s;", osid, verStr)))
	}
	tplData["osVerList"] = javascriptOsVersionStr
	tplData["osList"] = osList

	return tplData, nil
}

func getHostInfoListFromDB(requestData, conditions map[string]string, startPos int) ([]map[string]string, []sysadmerror.Sysadmerror, bool) {
	var errs []sysadmerror.Sysadmerror
	var ret []map[string]string

	var limit []int
	limit = append(limit, startPos)
	limit = append(limit, WorkingData.pageInfo.numPerPage)

	var order []db.OrderData
	direction := 0
	if requestData["direction"] == "1" {
		direction = 1
	}

	if requestData["orderfield"] != "" {
		fieldSlice := strings.Split(requestData["orderfield"], ",")
		for _, f := range fieldSlice {
			line := db.OrderData{Key: f, Order: direction}
			order = append(order, line)
		}
	}
	selectData := db.SelectData{
		Tb:        []string{"host"},
		OutFeilds: []string{"*"},
		Where:     conditions,
		Limit:     limit,
		Order:     order,
	}

	dbEntity := WorkingData.dbConf.Entity
	dbData, err := dbEntity.QueryData(&selectData)
	errs = append(errs, err...)
	if dbData == nil {
		return ret, errs, false
	}

	for _, l := range dbData {
		row := make(map[string]string, 0)
		for k, v := range l {
			vStr := utils.Interface2String(v)
			row[k] = vStr
		}
		ret = append(ret, row)
	}

	return ret, errs, true
}

func getHostCountFromDB(conditons map[string]string) (int, []sysadmerror.Sysadmerror, bool) {
	var errs []sysadmerror.Sysadmerror
	ret := 0

	selectData := db.SelectData{
		Tb:        []string{"host"},
		OutFeilds: []string{"count(hostid) as totalNum"},
		Where:     conditons,
	}

	dbEntity := WorkingData.dbConf.Entity
	dbData, err := dbEntity.QueryData(&selectData)
	errs = append(errs, err...)
	if dbData == nil {
		return ret, errs, false
	}

	row := dbData[0]
	result, ok := row["totalNum"]
	if !ok {
		return ret, errs, false
	}

	ret, e := utils.Interface2Int(result)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030009, "error", "can not get host count %s", e))
		return ret, errs, false
	}

	return ret, errs, true
}

func buildFootPageInfo(data map[string]string, tplData map[string]interface{}, hostCount int, pageInfoParas string) (map[string]interface{}, []sysadmerror.Sysadmerror, bool) {
	var errs []sysadmerror.Sysadmerror

	numPerPage := 0
	numperpageStr := data["numperpage"]
	if numperpageStr == "" {
		numPerPage = WorkingData.pageInfo.numPerPage
	} else {
		tmpNumPerPage, e := strconv.Atoi(numperpageStr)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030010, "error", "build foot page information error %s", e))
			return tplData, errs, false
		}
		numPerPage = tmpNumPerPage
	}

	start := 0
	startStr := data["start"]
	if startStr != "" {
		tmpStart, e := strconv.Atoi(startStr)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030011, "error", "build foot page information error %s", e))
			return tplData, errs, false
		}
		start = tmpStart
	}

	totalPages := int(math.Ceil(float64(hostCount) / float64(numPerPage)))
	currentPage := 1
	currentPage = int(math.Ceil(float64(start+1) / float64(numPerPage)))
	currentPageHTML := strconv.Itoa(currentPage)
	totalPageHTML := strconv.Itoa(totalPages)
	prePageHTML := ""
	if currentPage <= 1 {
		prePageHTML = " 上一页 "
	} else {
		preNum := start - numPerPage
		prePage := fmt.Sprintf("?start=%d&numPerPage=%d&%s", preNum, numPerPage, pageInfoParas)
		prePageHTML = "<a href=\"javascript:void(0)\" onclick='changePage(\"" + prePage + "\")'>上一页</a>"
	}

	nextPageHTML := ""
	if currentPage >= totalPages {
		nextPageHTML = "下一页 "
	} else {
		nextNum := start + numPerPage
		nextPage := fmt.Sprintf("?start=%d&num=%d&%s", nextNum, numPerPage, pageInfoParas)
		nextPageHTML = "<a href=\"javascript:void(0)\" onclick='changePage(\"" + nextPage + "\")'>下一页</a>"
	}

	tplData["currentpage"] = currentPageHTML
	tplData["totalpage"] = totalPageHTML
	tplData["prepage"] = template.HTML(prePageHTML)
	tplData["nextpage"] = template.HTML(nextPageHTML)

	return tplData, errs, true
}

func buildPageInfoParas(data map[string]string) string {
	pageInfoParas := ""
	if data["projectid"] != "" {
		pageInfoParas = pageInfoParas + "projectid=" + data["projectid"] + "&"
	}

	if data["userid"] != "" {
		pageInfoParas = pageInfoParas + "userid=" + data["userid"] + "&"
	}

	if data["hostid"] != "" {
		pageInfoParas = pageInfoParas + "hostid=" + data["hostid"] + "&"
	}

	if data["searchKey"] != "" {
		pageInfoParas = pageInfoParas + "searchKey=" + data["searchKey"] + "&"
	}

	if data["orderfield"] != "" {
		pageInfoParas = pageInfoParas + "orderfield=" + data["orderfield"] + "&"
	}

	if data["direction"] != "" {
		pageInfoParas = pageInfoParas + "direction" + data["direction"]
	}

	return pageInfoParas
}

func buildOsAndVersionInfo(tplData map[string]interface{}) (map[string]interface{}, []sysadmerror.Sysadmerror, bool) {
	var errs []sysadmerror.Sysadmerror

	osNameList := tplData["osNameList"].(map[int]string)
	osVerMap := tplData["osVerMap"].(map[int]map[int]string)
	hostData := tplData["HostData"].([]map[string]string)

	var newHostData []map[string]string
	for _, line := range hostData {
		hostOsIDStr := line["osID"]
		hostOsversionidStr := line["osversionid"]
		osName := ""

		osIDInt := 0
		for osID, name := range osNameList {
			osIDStr := strconv.Itoa(osID)
			if hostOsIDStr == osIDStr {
				osName = name
				osIDInt = osID
				break
			}
		}

		verMap, okVerMap := osVerMap[osIDInt]
		verName := ""
		if okVerMap {
			versionIDInt, e := strconv.Atoi(hostOsversionidStr)
			if e == nil {
				verName = verMap[versionIDInt]
			}
		}

		if osName == "" || verName == "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030012, "error", "get host os and os version error "))
			return tplData, errs, false
		}

		line["OsVerInfo"] = osName + "/" + verName

		newHostData = append(newHostData, line)
	}

	tplData["HostData"] = newHostData

	return tplData, errs, true
}

func buildHostStatusInfo(tplData map[string]interface{}) map[string]interface{} {

	hostData := tplData["HostData"].([]map[string]string)

	var newHostData []map[string]string
	for _, line := range hostData {
		statusCode := line["status"]
		if v, ok := hostStatus[statusCode]; ok {
			line["statusText"] = v
		} else {
			line["statusText"] = hostStatus["unkown"]
		}
		newHostData = append(newHostData, line)
	}

	tplData["HostData"] = newHostData

	return tplData
}

func listHostCondition(data map[string]string, isDeleted bool) (map[string]string, error) {
	whereMap := make(map[string]string, 0)
	if data["objectIds"] != "" {
		whereMap["hostid"] = "in (" + data["objectIds"] + ")"
		if !isDeleted {
			whereMap["status"] = "!=\"deleted\""
		}

		return whereMap, nil
	}

	if data["clusterID"] != "" {
		whereMap["k8sclusterid"] = "=\"" + data["clusterID"] + "\""
		if !isDeleted {
			whereMap["status"] = "!=\"deleted\""
		}

		return whereMap, nil
	}

	if data["azID"] != "" {
		whereMap["azid"] = "=\"" + data["azID"] + "\""
		if !isDeleted {
			whereMap["status"] = "!=\"deleted\""
		}

		return whereMap, nil
	}

	if data["dcID"] != "" {
		whereMap["dcid"] = "=\"" + data["dcID"] + "\""
		if !isDeleted {
			whereMap["status"] = "!=\"deleted\""
		}

		return whereMap, nil
	}

	if data["searchKey"] != "" {
		searchSql := "=1 and (hostname like \"%" + data["searchKey"] + "%\" or ip like \"%" + data["searchKey"] + "%\")"
		whereMap["1"] = searchSql
		if !isDeleted {
			whereMap["status"] = "!=\"deleted\""
		}

		return whereMap, nil
	}

	return whereMap, nil
}

func buildSelectData(tplData map[string]interface{}, dcList []interface{}, requestData map[string]string) error {

	selectedDC := strings.TrimSpace(requestData["dcID"])
	if selectedDC == "" {
		selectedDC = "0"
	}

	selectedAZ := strings.TrimSpace(requestData["azID"])
	if selectedAZ == "" {
		selectedAZ = "0"
	}

	selectedCluster := strings.TrimSpace(requestData["clusterID"])
	if selectedCluster == "" {
		selectedCluster = "0"
	}

	// preparing datacenter options
	var dcOptions []objectsUI.SelectOption
	dcOption := objectsUI.SelectOption{
		Id:       "0",
		Text:     "===选择数据中心===",
		ParentID: "0",
	}
	dcOptions = append(dcOptions, dcOption)
	for _, line := range dcList {
		dcData, ok := line.(datacenter.DatacenterSchema)
		if !ok {
			return fmt.Errorf("the data is not datacenter schema")
		}
		dcOption := objectsUI.SelectOption{
			Id:       strconv.Itoa(int(dcData.Id)),
			Text:     dcData.CnName,
			ParentID: "0",
		}
		dcOptions = append(dcOptions, dcOption)
	}
	dcSelect := objectsUI.SelectData{Title: "数据中心", SelectedId: selectedDC, Options: dcOptions}
	tplData["dcSelect"] = dcSelect

	// preparing AZ options
	var azOptions []objectsUI.SelectOption
	azOption := objectsUI.SelectOption{
		Id:       "0",
		Text:     "===选择可用区===",
		ParentID: "0",
	}
	azOptions = append(azOptions, azOption)
	if selectedDC != "0" {
		var azEntity sysadmObjects.ObjectEntity
		azEntity = sysadmAZObj.New()
		conditions := make(map[string]string, 0)
		var emptyString []string
		conditions["isDeleted"] = "='0'"
		conditions["dcid"] = "='" + selectedDC + "'"
		azList, e := azEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, make(map[string]string))
		if e != nil {
			return e
		}
		for _, line := range azList {
			azData, ok := line.(sysadmAZObj.AvailablezoneSchema)
			if !ok {
				return fmt.Errorf("the data is not AZ schema")
			}

			azOption := objectsUI.SelectOption{
				Id:       strconv.Itoa(int(azData.Id)),
				Text:     azData.CnName,
				ParentID: strconv.Itoa(int(azData.Datacenterid)),
			}
			azOptions = append(azOptions, azOption)
		}
	}
	azSelect := objectsUI.SelectData{Title: "可用区", SelectedId: selectedAZ, Options: azOptions}
	tplData["azSelect"] = azSelect

	// preparing cluster options
	var clusterOptions []objectsUI.SelectOption
	clusterOption := objectsUI.SelectOption{
		Id:       "0",
		Text:     "===选择集群===",
		ParentID: "0",
	}
	clusterOptions = append(clusterOptions, clusterOption)
	if selectedDC != "0" {
		var k8sclusterEntity sysadmObjects.ObjectEntity
		k8sclusterEntity = sysadmK8sCluster.New()
		conditions := make(map[string]string, 0)
		var emptyString []string
		conditions["isDeleted"] = "='0'"
		conditions["azid"] = "='" + selectedAZ + "'"
		clusterList, e := k8sclusterEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, make(map[string]string))
		if e != nil {
			return e
		}

		for _, line := range clusterList {
			clusterData, ok := line.(sysadmK8sCluster.K8sclusterSchema)
			if !ok {
				return fmt.Errorf("the data is not cluster schema")
			}

			clusterOption := objectsUI.SelectOption{
				Id:       clusterData.Id,
				Text:     clusterData.CnName,
				ParentID: strconv.Itoa(int(clusterData.Azid)),
			}
			clusterOptions = append(clusterOptions, clusterOption)
		}
	}
	clusterSelect := objectsUI.SelectData{Title: "集群", SelectedId: selectedCluster, Options: clusterOptions}
	tplData["clusterSelect"] = clusterSelect

	return nil
}
