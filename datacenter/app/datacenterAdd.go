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
	"sysadm/objectsUI"
	sysadmRegion "sysadm/region/app"
	"sysadm/sysadmLog"
	"sysadm/sysadmapi/apiutils"
	"sysadm/user"
	"sysadm/utils"
)

func addformHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(7000160001, "debug", "try to display add object form page"))

	addTemplateFile := "addObjectFormNew.html"
	var emptyString []string
	baseUri := "/" + DefaultModuleName + "/"
	enctype := ""
	postUri := "/api/" + DefaultApiVersion + "/" + DefaultModuleName + "/add"
	var tplDataLines []objectsUI.ObjLineData

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000160002, errs, e)
		return
	}

	// preparing select data
	conditions := make(map[string]string, 0)
	order := make(map[string]string, 0)
	regionEntiy := sysadmRegion.New()
	conditions["display"] = "='" + sysadmRegion.CountryDisplay + "'"
	countryList, e := regionEntiy.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000160003, errs, e)
		return
	}
	countryOptions := buildCountryOptions(countryList)
	var regionItems []objectsUI.ObjItemInfo
	countrySelect := objectsUI.ObjItemInfo{Title: "国家", ID: "country", Name: "country", Kind: "SELECT",
		ActionUri: "getprovincebycountrycodeforselect", ItemData: countryOptions, SubObjID: "province", JsActionKind: objectsUI.JsActionKind_Select_Change_SelectOptions}
	regionItems = append(regionItems, countrySelect)

	var provinceOptions []objectsUI.SubItems
	provinceOptions = append(provinceOptions, objectsUI.SubItems{Value: "0", Text: "===选择省份===", Checked: true})
	provinceSelect := objectsUI.ObjItemInfo{Title: "省份", ID: "province", Name: "province", Kind: "SELECT",
		ActionUri: "getcitybyprovincecodeforselect", ItemData: provinceOptions, SubObjID: "city", JsActionKind: objectsUI.JsActionKind_Select_Change_SelectOptions}
	regionItems = append(regionItems, provinceSelect)

	var cityOptions []objectsUI.SubItems
	cityOptions = append(cityOptions, objectsUI.SubItems{Value: "0", Text: "===选择城市===", Checked: true})
	citySelect := objectsUI.ObjItemInfo{Title: "城市", ID: "city", Name: "city", Kind: "SELECT", ItemData: cityOptions}
	regionItems = append(regionItems, citySelect)
	regionLineData := objectsUI.ObjLineData{Items: regionItems}
	tplDataLines = append(tplDataLines, regionLineData)

	cnName := objectsUI.ObjItemInfo{Title: "中文名称", ID: "cnName", Name: "cnName", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validCnName", Note: "数据中心中文名称"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{cnName}})

	enName := objectsUI.ObjItemInfo{Title: "English Name", ID: "enName", Name: "enName", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validEnName", Note: "English Name"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{enName}})

	address := objectsUI.ObjItemInfo{Title: "地址", ID: "address", Name: "address", Kind: "TEXT", Size: 30, Note: "数据中心详细地址"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{address}})

	dutyTel := objectsUI.ObjItemInfo{Title: "值班电话", ID: "dutyTel", Name: "dutyTel", Kind: "TEXT", Size: 30, Note: "值班电话"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{dutyTel}})

	typeOptions := buildTypeOptions()
	lineType := objectsUI.ObjItemInfo{Title: "线路类型", ID: `type`, Name: `type`, Kind: `RADIO`, Note: "", ItemData: typeOptions}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{lineType}})

	remark := objectsUI.ObjItemInfo{Title: "备注", ID: "remark", Name: "remark", Kind: "TEXTAREA", Size: 40, Rows: 5}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{remark}})

	tplData, _ := objectsUI.InitAddObjectFormTemplateData(baseUri, "基础设施", "数据中心",
		enctype, postUri, "list", "list", emptyString, emptyString)
	tplData["data"] = tplDataLines

	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, addTemplateFile, tplData)
}

func buildCountryOptions(countryList []interface{}) []objectsUI.SubItems {
	var ret []objectsUI.SubItems
	lineData := objectsUI.SubItems{Value: "0", Text: "===选择国家===", Checked: true}
	ret = append(ret, lineData)

	for _, line := range countryList {
		lineData := objectsUI.SubItems{}
		countryData := line.(sysadmRegion.CountrySchema)
		lineData.Value = countryData.Code
		lineData.Text = countryData.ChineseName
		lineData.Checked = false
		ret = append(ret, lineData)
	}

	return ret
}

func buildTypeOptions() []objectsUI.SubItems {
	var ret []objectsUI.SubItems
	subItem := objectsUI.SubItems{Value: strconv.Itoa(LineTypeCT), Text: "电信", Checked: true}
	ret = append(ret, subItem)

	subItem = objectsUI.SubItems{Value: strconv.Itoa(LineTypeCUCC), Text: "联通", Checked: false}
	ret = append(ret, subItem)

	subItem = objectsUI.SubItems{Value: strconv.Itoa(LinetypeCMCC), Text: "移动", Checked: false}
	ret = append(ret, subItem)

	subItem = objectsUI.SubItems{Value: strconv.Itoa(LineTypeCBN), Text: "广电", Checked: false}
	ret = append(ret, subItem)

	subItem = objectsUI.SubItems{Value: strconv.Itoa(LineTypeBGP2), Text: "双线BGP", Checked: false}
	ret = append(ret, subItem)

	subItem = objectsUI.SubItems{Value: strconv.Itoa(LineTypeBGP3), Text: "三线BGP", Checked: false}
	ret = append(ret, subItem)

	subItem = objectsUI.SubItems{Value: strconv.Itoa(LineTypeBGP4), Text: "四线BGP", Checked: false}
	ret = append(ret, subItem)

	subItem = objectsUI.SubItems{Value: strconv.Itoa(LineTypeOverseas), Text: "海外", Checked: false}
	ret = append(ret, subItem)

	return ret
}

func getprovincebycountrycodeforselectHandler(c *sysadmServer.Context) {
	requestData, e := utils.NewGetRequestData(c, []string{"objID"})
	if e != nil || requestData["objID"] == "" {
		response := apiutils.BuildResponseDataForError(7000160004, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	regionEntity := sysadmRegion.New()
	var emptyString []string
	conditions := make(map[string]string, 0)
	order := make(map[string]string, 0)
	conditions["countryCode"] = "='" + requestData["objID"] + "'"

	provinceList, e := regionEntity.GetProvinceList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		response := apiutils.BuildResponseDataForError(7000160005, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	msg := ""
	for _, line := range provinceList {
		lineData, ok := line.(sysadmRegion.ProvinceSchema)
		if !ok {
			response := apiutils.BuildResponseDataForError(7000160006, "数据处理错误")
			c.JSON(http.StatusOK, response)
			return
		}
		lineStr := lineData.Code + ":" + lineData.Name
		if msg == "" {
			msg = lineStr
		} else {
			msg = msg + "," + lineStr
		}
	}

	response := apiutils.BuildResponseDataForSuccess(msg)
	c.JSON(http.StatusOK, response)
}

func getcitybyprovincecodeforselectHandler(c *sysadmServer.Context) {
	requestData, e := utils.NewGetRequestData(c, []string{"objID"})
	if e != nil || requestData["objID"] == "" {
		response := apiutils.BuildResponseDataForError(7000160007, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	regionEntity := sysadmRegion.New()
	var emptyString []string
	conditions := make(map[string]string, 0)
	order := make(map[string]string, 0)
	conditions["provinceCode"] = "='" + requestData["objID"] + "'"

	cityList, e := regionEntity.GetCityList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		response := apiutils.BuildResponseDataForError(7000160008, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	msg := ""
	for _, line := range cityList {
		lineData, ok := line.(sysadmRegion.CitySchema)
		if !ok {
			response := apiutils.BuildResponseDataForError(7000160009, "数据处理错误")
			c.JSON(http.StatusOK, response)
			return
		}
		lineStr := lineData.Code + ":" + lineData.Name
		if msg == "" {
			msg = lineStr
		} else {
			msg = msg + "," + lineStr
		}
	}

	response := apiutils.BuildResponseDataForSuccess(msg)
	c.JSON(http.StatusOK, response)
}