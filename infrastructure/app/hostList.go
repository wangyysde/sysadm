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
	"sysadm/db"
	"sysadm/httpclient"
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
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030003, "error", "user should login %s", e))
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
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030006, "error", "get error yum list error %s", e))
		logErrors(errs)
		return
	}

	// get projects list
	tplData, e = getProjectList(c, apiServer, "showmessage.html", tplData)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030007, "error", "get error project list error %s", e))
		logErrors(errs)
		return
	}

	requestData, err, ok := getRequestData(c)
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		errorTplData := map[string]interface{}{
			"errormessage": "get request data error",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, errorTplData)

		return
	}

	// get total number of hosts
	hostCount, err, ok := getHostCountFromDB(requestData, false)
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
	// caclusiong page information for HTML

	// get host list data
	hostData, err, ok := getHostInfoListFromDB(requestData, false)
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

	// get project name
	tplData, err, ok = buildHostProjectInfo(tplData)
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		errorTplData := map[string]interface{}{
			"errormessage": "get project name error ",
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

	selectedprojectid := requestData["projectid"]
	if selectedprojectid == "" {
		selectedprojectid = "0"
	}
	tplData["selectedprojectid"] = selectedprojectid

	templateFile := "infrastructurelist.html"
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

func getProjectList(c *sysadmServer.Context, apiServer *ApiServer, msgTpl string, tplData map[string]interface{}) (map[string]interface{}, error) {
	var errs []sysadmerror.Sysadmerror

	// get project for select menu
	moduleName := "project"
	actionName := "list"

	apiServerData := apiutils.BuildApiServerData(moduleName, actionName, apiServer.ApiVersion, apiServer.Tls.IsTls,
		apiServer.Server.Address, apiServer.Server.Port, apiServer.Tls.Ca, apiServer.Tls.Cert, apiServer.Tls.Key)
	if apiServerData == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101005, "error", "api server parameters error"))
		logErrors(errs)
		tplData := map[string]interface{}{
			"errormessage": "api server parameters error",
		}
		c.HTML(http.StatusOK, msgTpl, tplData)
		return tplData, fmt.Errorf("api server parameters error")
	}

	urlRaw, err := apiutils.BuildApiUrl(apiServerData)
	errs = append(errs, err...)
	if urlRaw == "" {
		err := apiutils.SendResponseForErrorMessage(c, 1101007, "api server parameters error")
		errs = append(errs, err...)
		logErrors(errs)
		tplData := map[string]interface{}{
			"errormessage": "api server parameters error",
		}
		c.HTML(http.StatusOK, msgTpl, tplData)
		return tplData, fmt.Errorf("api server parameters error")
	}

	requestParas := httpclient.RequestParams{}
	requestParas.Url = urlRaw
	requestParas.Method = http.MethodPost
	body, err := httpclient.SendRequest(&requestParas)
	errs = append(errs, err...)
	ret, err := apiutils.ParseResponseBody(body)
	errs = append(errs, err...)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101004, "debug", "now handling project list %#v", ret))
	logErrors(errs)
	if !ret.Status {
		errCode := ret.ErrorCode
		msgLines := ret.Message
		msgLine := msgLines[0]
		errMsg := msgLine["msg"].(string)
		tplData := map[string]interface{}{
			"errormessage": fmt.Sprintf("errorCode: %d Msg: %s", errCode, errMsg),
		}
		c.HTML(http.StatusOK, msgTpl, tplData)
		return tplData, fmt.Errorf("errorCode: %d Msg: %s", errCode, errMsg)
	}

	// preparing project select data
	var projectInfo []map[string]string
	projectInfo = append(projectInfo, map[string]string{"0": "全部项目"})
	res := ret.Message
	for _, line := range res {
		lineMap := make(map[string]string, 0)
		id := utils.Interface2String(line["projectid"])
		name := utils.Interface2String(line["name"])
		lineMap[id] = name
		projectInfo = append(projectInfo, lineMap)
	}

	tplData["projectinfo"] = projectInfo

	return tplData, nil
}

func getRequestData(c *sysadmServer.Context) (map[string]string, []sysadmerror.Sysadmerror, bool) {
	var errs []sysadmerror.Sysadmerror
	ret := make(map[string]string, 0)

	dataMap, err := utils.GetRequestData(c, []string{"projectid", "userid", "hostid", "searchKey", "start", "numperpage", "orderfield", "direction"})
	errs = append(errs, err...)

	projectid := ""
	projectid, ok := dataMap["projectid"]
	if ok {
		projectid = strings.TrimSpace(projectid)
	}
	ret["projectid"] = projectid

	// try to get userid in request data
	userid := ""
	userid, ok = dataMap["userid"]
	if ok {
		userid = strings.TrimSpace(userid)
	}

	// try to get get userid from session if there is not userid in request data
	if userid == "" {
		// get userid
		useridInterface, e := user.GetSessionValue(c, "userid", WorkingData.sessionOption.sessionName)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030008, "error", "can not got user", e))
			return ret, errs, false
		}

		userid = utils.Interface2String(useridInterface)

	}
	ret["userid"] = userid

	hostids := ""
	hostIDMap, _ := utils.GetRequestDataArray(c, []string{"hostid[]"})
	if hostIDMap != nil {
		hostIDSlice, ok := hostIDMap["hostid[]"]
		if ok {
			hostids = strings.Join(hostIDSlice, ",")
		}
	}
	ret["hostids"] = hostids

	searchKey := ""
	searchKey, ok = dataMap["searchKey"]
	if ok {
		searchKey = strings.TrimSpace(searchKey)
	}
	ret["searchKey"] = searchKey

	start := ""
	start, ok = dataMap["start"]
	if ok {
		start = strings.TrimSpace(start)
	}
	ret["start"] = start

	numperpage := ""
	numperpage, ok = dataMap["numperpage"]
	if ok {
		numperpage = strings.TrimSpace(numperpage)
	}
	ret["numperpage"] = numperpage

	orderfield := ""
	orderfield, ok = dataMap["orderfield"]
	if ok {
		orderfield = strings.TrimSpace(orderfield)
	}
	ret["orderfield"] = orderfield

	direction := ""
	direction, ok = dataMap["direction"]
	if ok {
		direction = strings.TrimSpace(direction)
	}
	ret["direction"] = direction

	return ret, errs, true
}

func getHostInfoListFromDB(data map[string]string, isDeleted bool) ([]map[string]string, []sysadmerror.Sysadmerror, bool) {
	var errs []sysadmerror.Sysadmerror
	var ret []map[string]string

	// Qeurying data from DB
	whereMap := make(map[string]string, 0)
	if data["projectid"] != "" {
		whereMap["projectid"] = "=\"" + data["projectid"] + "\""
	}

	if data["userid"] != "" {
		whereMap["userid"] = "=\"" + data["userid"] + "\""
	}

	if data["hostids"] != "" {
		whereMap["hostid"] = "in (" + data["hostids"] + ")"
	}

	if data["searchKey"] != "" {
		searchSql := "=1 and (hostname like \"%" + data["searchKey"] + "%\" or ip like \"%" + data["searchKey"] + "%\")"
		whereMap["1"] = searchSql
	}

	if !isDeleted {
		whereMap["status"] = "!=\"deleted\""
	}

	num := 0
	if data["numperpage"] != "" {
		tmpNum, e := strconv.Atoi(data["numperpage"])
		if e != nil {
			num = tmpNum
		}
	}

	var limit []int
	if num > 0 {
		startInt := 0
		if data["start"] != "" {
			if tmpStartInt, e := strconv.Atoi(data["start"]); e != nil {
				startInt = tmpStartInt
			}
		}
		limit = append(limit, startInt)
		limit = append(limit, num)
	}

	var order []db.OrderData
	direction := 0
	if strings.ToUpper(data["direction"]) == "DESC" {
		direction = 1
	}
	if data["orderfield"] != "" {
		fieldSlice := strings.Split(data["orderfield"], ",")
		for _, f := range fieldSlice {
			line := db.OrderData{Key: f, Order: direction}
			order = append(order, line)
		}
	}
	selectData := db.SelectData{
		Tb:        []string{"host"},
		OutFeilds: []string{"*"},
		Where:     whereMap,
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

func getHostCountFromDB(data map[string]string, isDeleted bool) (int, []sysadmerror.Sysadmerror, bool) {
	var errs []sysadmerror.Sysadmerror
	ret := 0

	// Qeurying data from DB
	whereMap := make(map[string]string, 0)
	if data["projectid"] != "" {
		whereMap["projectid"] = "=\"" + data["projectid"] + "\""
	}

	if data["userid"] != "" {
		whereMap["userid"] = "=\"" + data["userid"] + "\""
	}

	if data["hostids"] != "" {
		whereMap["hostid"] = "in (" + data["hostids"] + ")"
	}

	if data["searchKey"] != "" {
		searchSql := "=1 and (hostname like \"%" + data["searchKey"] + "%\" or ip like \"%" + data["searchKey"] + "%\")"
		whereMap["1"] = searchSql
	}

	if !isDeleted {
		whereMap["status"] = "!=\"deleted\""
	}

	selectData := db.SelectData{
		Tb:        []string{"host"},
		OutFeilds: []string{"count(hostid) as totalNum"},
		Where:     whereMap,
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

func buildHostProjectInfo(tplData map[string]interface{}) (map[string]interface{}, []sysadmerror.Sysadmerror, bool) {
	var errs []sysadmerror.Sysadmerror

	projectinfo := tplData["projectinfo"].([]map[string]string)
	hostData := tplData["HostData"].([]map[string]string)

	var newHostData []map[string]string
	for _, line := range hostData {
		projectid := line["projectid"]
		projectName := ""
		for _, projectLine := range projectinfo {
			if name, ok := projectLine[projectid]; ok {
				projectName = name
				break
			}

		}

		if projectName == "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030013, "error", "get project name error"))
			return tplData, errs, false
		}

		line["projectName"] = projectName
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
