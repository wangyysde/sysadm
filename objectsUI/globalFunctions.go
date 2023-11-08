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

package objectsUI

import (
	"fmt"
	"html/template"
	"math"
	"sort"
	"strconv"
	"strings"
)

func InitTemplateData(baseUri, mainCategory, subCategory, addButtonTitle, isSearchForm string,
	allPopmenuItems, addtionalJs []string, requestData map[string]string) (map[string]interface{}, error) {
	tplData := make(map[string]interface{}, 0)

	baseUri = strings.TrimSpace(baseUri)
	mainCategory = strings.TrimSpace(mainCategory)
	subCategory = strings.TrimSpace(subCategory)
	addButtonTitle = strings.TrimSpace(addButtonTitle)
	isSearchForm = strings.TrimSpace(isSearchForm)
	if baseUri == "" || mainCategory == "" || subCategory == "" {
		return nil, fmt.Errorf("data is not valid")
	}
	tplData["baseUri"] = baseUri
	tplData["mainCategory"] = mainCategory
	tplData["subCategory"] = subCategory
	tplData["addButtonTitle"] = addButtonTitle
	tplData["isSearchForm"] = isSearchForm
	tplData["allPopMenuItems"] = allPopmenuItems
	tplData["addtionalJs"] = addtionalJs
	tplData["groupSelectID"] = requestData["groupSelectID"]

	return tplData, nil
}

func InitTemplateDataForWorkload(baseUri, mainCategory, subCategory, addButtonTitle, isSearchForm string,
	allPopmenuItems, addtionalJs, addtionalCss []string, requestData map[string]string) (map[string]interface{}, error) {
	tplData := make(map[string]interface{}, 0)

	baseUri = strings.TrimSpace(baseUri)
	mainCategory = strings.TrimSpace(mainCategory)
	subCategory = strings.TrimSpace(subCategory)
	addButtonTitle = strings.TrimSpace(addButtonTitle)
	isSearchForm = strings.TrimSpace(isSearchForm)
	if baseUri == "" || mainCategory == "" || subCategory == "" {
		return nil, fmt.Errorf("data is not valid")
	}
	tplData["baseUri"] = baseUri
	tplData["mainCategory"] = mainCategory
	tplData["subCategory"] = subCategory
	tplData["addButtonTitle"] = addButtonTitle
	isSearchForm = strings.ToLower(isSearchForm)
	isSearchBool := false
	if isSearchForm == "y" || isSearchForm == "yes" || isSearchForm == "true" || isSearchForm == "1" {
		isSearchBool = true
	}
	tplData["isSearchForm"] = isSearchBool
	tplData["allPopMenuItems"] = allPopmenuItems
	tplData["addtionalJs"] = addtionalJs
	tplData["addtionalCss"] = addtionalCss
	tplData["groupSelectID"] = requestData["groupSelectID"]

	return tplData, nil
}

func GetSearchContentFromRequest(requestData map[string]string) string {
	if searchContent, ok := requestData["searchContent"]; ok {
		return strings.TrimSpace(searchContent)
	}
	return ""
}

func GetObjectIdsFromRequest(requestData map[string]string) []string {
	var ids []string
	if strings.TrimSpace(requestData["objectIds"]) != "" {
		ids = strings.Split(strings.TrimSpace(requestData["objectIds"]), ",")
	}

	return ids
}

func GetStartPosFromRequest(requestData map[string]string) int {
	if startStr, ok := requestData["start"]; ok {
		if start, e := strconv.Atoi(startStr); e == nil {
			return start
		}
	}

	return 0
}

func BuildCondition(requestData map[string]string, isDeleted, groupFieldName string) map[string]string {
	ret := make(map[string]string, 0)
	if strings.TrimSpace(isDeleted) != "" {
		ret["isDeleted"] = strings.TrimSpace(isDeleted)
	}

	groupFieldName = strings.TrimSpace(groupFieldName)
	if requestData["groupSelectID"] != "" && groupFieldName != "" && requestData["groupSelectID"] != "0" {
		ret[groupFieldName] = "=" + requestData["groupSelectID"]

	}

	return ret
}

func BuildOrderDataForQuery(requestData, allOrderFields map[string]string, defaultOrderField, defaultOrderDirction string) map[string]string {
	ret := make(map[string]string, 0)

	requestedOrderField := strings.TrimSpace(strings.ToUpper(requestData["orderfield"]))
	objOrderFieldName := ""
	if tmpField, ok := allOrderFields[requestedOrderField]; ok && tmpField != "" {
		objOrderFieldName = tmpField
	} else {
		objOrderFieldName = allOrderFields[defaultOrderField]
	}

	requestOrderDirection := strings.TrimSpace(requestData["direction"])
	objOrderDirection := ""
	if requestOrderDirection == "" || (requestOrderDirection != "0" && requestOrderDirection != "1") {
		objOrderDirection = defaultOrderDirction
	} else {
		objOrderDirection = requestOrderDirection
	}

	ret[objOrderFieldName] = objOrderDirection

	return ret
}

func BuildThData(requestData, allOrderFields, allListItems map[string]string, tplData map[string]interface{}, defaultOrderField, defaultOrderDirection string) {
	searchContent := GetSearchContentFromRequest(requestData)
	groupSelectID := tplData["groupSelectID"].(string)
	noSelectOrderUri := ""
	if searchContent != "" {
		noSelectOrderUri = "searchContent=" + searchContent
	}

	if groupSelectID != "" {
		if noSelectOrderUri == "" {
			noSelectOrderUri = "groupSelectID=" + groupSelectID
		} else {
			noSelectOrderUri = noSelectOrderUri + "&groupSelectID=" + groupSelectID
		}
	}

	selectedOrderField := strings.TrimSpace(strings.ToUpper(requestData["orderfield"]))
	if selectedOrderField == "" {
		selectedOrderField = defaultOrderField
	}
	selectOrderDirection := strings.TrimSpace(requestData["direction"])
	lastSelectDirction := selectOrderDirection
	if selectOrderDirection == "" {
		selectOrderDirection = defaultOrderDirection
		lastSelectDirction = defaultOrderDirection
	} else {
		if selectOrderDirection == "1" {
			selectOrderDirection = "0"
		} else {
			selectOrderDirection = "1"
		}
	}

	thData := make([]ObjectTitle, 0)
	var thKeys []string
	for k, _ := range allListItems {
		thKeys = append(thKeys, k)
	}
	sort.Strings(thKeys)

	for _, v := range thKeys {
		lineData := ObjectTitle{}
		lineData.ID = v
		lineData.Title = allListItems[v]
		orderUri := ""
		orderDirection := ""
		if isOrder(v, allOrderFields) {
			if v == selectedOrderField {
				if noSelectOrderUri == "" {
					orderUri = "orderfield=" + v + "&direction=" + selectOrderDirection
				} else {
					orderUri = noSelectOrderUri + "&orderfield=" + v + "&direction=" + selectOrderDirection
				}
				orderDirection = lastSelectDirction
			} else {
				if noSelectOrderUri == "" {
					orderUri = "orderfield=" + v + "&direction=" + defaultOrderDirection
				} else {
					orderUri = noSelectOrderUri + "&orderfield=" + v + "&direction=" + defaultOrderDirection
				}
				orderDirection = defaultOrderDirection
			}
		}
		lineData.OrderUri = orderUri
		lineData.OrderDirection = orderDirection
		if v == selectedOrderField {
			lineData.OrderSelected = "yes"
		} else {
			lineData.OrderSelected = ""
		}
		thData = append(thData, lineData)
	}

	tplData["thData"] = thData
}

func isOrder(key string, allOrderFields map[string]string) bool {
	for k, _ := range allOrderFields {
		if strings.TrimSpace(strings.ToUpper(key)) == strings.TrimSpace(strings.ToUpper(k)) {
			return true
		}
	}

	return false
}

func BuildPageNumInfo(tplData map[string]interface{}, requestData map[string]string, totalNum, startPos, numPerPage int,
	defaultOrderField, defaultOrderDirection string) {

	totalPages := int(math.Ceil(float64(totalNum) / float64(numPerPage)))
	currentPage := int(math.Ceil(float64(startPos+1) / float64(numPerPage)))
	preStart := 0
	if startPos >= numPerPage {
		preStart = startPos - numPerPage
	}
	nextStart := 0
	if totalNum > (startPos + numPerPage) {
		nextStart = startPos + numPerPage
	}

	searchContent := GetSearchContentFromRequest(requestData)
	groupSelectID := tplData["groupSelectID"].(string)
	uri := ""
	if searchContent != "" {
		uri = "searchContent=" + searchContent
	}

	if groupSelectID != "" {
		if uri == "" {
			uri = "groupSelectID=" + groupSelectID
		} else {
			uri = uri + "&groupSelectID=" + groupSelectID
		}
	}

	selectedOrderField := strings.TrimSpace(strings.ToUpper(requestData["orderfield"]))
	if selectedOrderField == "" {
		selectedOrderField = defaultOrderField
	}
	selectOrderDirection := strings.TrimSpace(requestData["direction"])
	if selectOrderDirection == "" {
		selectOrderDirection = defaultOrderDirection
	}

	if uri == "" {
		uri = "orderfield=" + selectedOrderField + "&direction=" + selectOrderDirection
	} else {
		uri = uri + "&orderfield=" + selectedOrderField + "&direction=" + selectOrderDirection
	}

	currentPageStr := strconv.Itoa(currentPage)
	prePageUri := ""
	if preStart > 0 || startPos > 0 {
		prePageUri = uri + "&start=" + strconv.Itoa(preStart)
	}

	nextPageUri := ""
	if nextStart != 0 {
		nextPageUri = uri + "&start=" + strconv.Itoa(nextStart)
	}

	totalPageStr := strconv.Itoa(totalPages)

	tplData["currentPage"] = currentPageStr
	tplData["prePageUri"] = prePageUri
	tplData["nextPageUri"] = nextPageUri
	tplData["totalPage"] = totalPageStr
}

func BuildPageNumInfoForWorkloadList(tplData map[string]interface{}, requestData map[string]string, totalNum, startPos, numPerPage int,
	defaultOrderField, defaultOrderDirection string) {

	totalPages := int(math.Ceil(float64(totalNum) / float64(numPerPage)))
	currentPage := int(math.Ceil(float64(startPos+1) / float64(numPerPage)))
	preStart := 0
	if startPos >= numPerPage {
		preStart = startPos - numPerPage
	}
	nextStart := 0
	if totalNum > (startPos + numPerPage) {
		nextStart = startPos + numPerPage
	}

	searchContent := GetSearchContentFromRequest(requestData)
	tplData["searchContent"] = searchContent

	selectedOrderField := strings.TrimSpace(strings.ToUpper(requestData["orderfield"]))
	if selectedOrderField == "" {
		selectedOrderField = defaultOrderField
	}
	selectOrderDirection := strings.TrimSpace(requestData["direction"])
	if selectOrderDirection == "" {
		selectOrderDirection = defaultOrderDirection
	}
	tplData["orderfield"] = selectedOrderField
	tplData["direction"] = selectOrderDirection

	currentPageStr := strconv.Itoa(currentPage)
	prePageStr := ""
	if preStart > 0 || startPos > 0 {
		prePageStr = strconv.Itoa(preStart)
	}

	nextPageStr := ""
	if nextStart != 0 {
		nextPageStr = strconv.Itoa(nextStart)

	}

	totalPageStr := strconv.Itoa(totalPages)

	tplData["currentPage"] = currentPageStr
	tplData["prePage"] = prePageStr
	tplData["nextPage"] = nextPageStr
	tplData["totalPage"] = totalPageStr
}

func InitAddObjectFormTemplateData(baseUri, mainCategory, subCategory, enctype, postUri, submitRedirect, cancelRedirect string, addtionalJs, addtionalCss []string) (map[string]interface{}, error) {
	tplData := make(map[string]interface{}, 0)

	baseUri = strings.TrimSpace(baseUri)
	mainCategory = strings.TrimSpace(mainCategory)
	subCategory = strings.TrimSpace(subCategory)
	enctype = strings.TrimSpace(enctype)
	postUri = strings.TrimSpace(postUri)
	submitRedirect = strings.TrimSpace(submitRedirect)
	cancelRedirect = strings.TrimSpace(cancelRedirect)

	tplData["baseUri"] = baseUri
	tplData["mainCategory"] = mainCategory
	tplData["subCategory"] = subCategory

	tplData["addtionalJs"] = addtionalJs
	tplData["addtionalCss"] = addtionalCss
	tplData["enctype"] = enctype
	tplData["postUri"] = postUri
	tplData["submitRedirect"] = submitRedirect
	tplData["cancelRedirect"] = cancelRedirect

	return tplData, nil
}

func InitTemplateForShowObjectDetails(mainCategory, subCategory, redirectUrl, baseUri string) (map[string]interface{}, error) {
	tplData := make(map[string]interface{}, 0)

	redirectUrl = strings.TrimSpace(redirectUrl)
	mainCategory = strings.TrimSpace(mainCategory)
	subCategory = strings.TrimSpace(subCategory)

	tplData["redirectUrl"] = redirectUrl
	tplData["baseUri"] = baseUri
	tplData["mainCategory"] = mainCategory
	tplData["subCategory"] = subCategory

	return tplData, nil
}

// BuildMultiSelectData 根据传递过来的第一级、第二级和第三级菜单数据firstData,secondData,thirdData设置对象列表页面上下拉菜单数据。
// 该下拉菜单最多可支持三级。如果出错返回错误信息，否则返回nil
func BuildMultiSelectData(firstData, secondData, thirdData SelectData, tplData map[string]interface{}) error {

	// set first select data
	firstOptions := firstData.Options
	tplData["firstGroupTitle"] = strings.TrimSpace(firstData.Title)
	if len(firstOptions) > 0 && (strings.TrimSpace(firstData.SelectedId) == "" || strings.TrimSpace(firstData.SelectedId) == "0") {
		tplData["firstGroupSelectedID"] = "0"
		option := SelectOption{Id: "0", Text: "---选择选项---"}
		firstOptions = append(firstOptions, option)
	} else {
		tplData["firstGroupSelectedID"] = strings.TrimSpace(firstData.SelectedId)
	}
	tplData["firstGroupSelect"] = firstOptions

	// set second select data
	if (strings.TrimSpace(secondData.SelectedId) != "" && strings.TrimSpace(secondData.SelectedId) != "0") && len(secondData.SelectedOptions) < 1 {
		return fmt.Errorf("the second select has be selected but the selected options is empty")
	}
	tplData["secondGroupSelectedID"] = "0"
	if tplData["firstGroupSelectedID"] != "0" && len(secondData.SelectedOptions) > 0 && secondData.SelectedId != "" {
		tplData["secondGroupSelectedID"] = strings.TrimSpace(secondData.SelectedId)
	}
	tplData["secondGroupTitle"] = strings.TrimSpace(secondData.Title)
	tplData["secondSelectedOptions"] = secondData.SelectedOptions
	var secondSelectOptionsJSStr []template.JS
	for _, o := range secondData.Options {
		secondSelectOptionsJSStr = append(secondSelectOptionsJSStr, template.JS(fmt.Sprintf("secondGroupOptions['%s'] = [%s]", o.ParentID, o.OptionsList)))
	}
	tplData["secondGroupSelect"] = secondSelectOptionsJSStr

	// set third select data
	if (strings.TrimSpace(thirdData.SelectedId) != "" && strings.TrimSpace(thirdData.SelectedId) != "0") && len(thirdData.SelectedOptions) < 1 {
		return fmt.Errorf("the third select has be selected but the selected options is empty")
	}
	tplData["thirdGroupSelectedID"] = "0"
	if tplData["secondGroupSelectedID"] != "0" && len(thirdData.SelectedOptions) > 0 && thirdData.SelectedId != "" {
		tplData["thirdGroupSelectedID"] = strings.TrimSpace(thirdData.SelectedId)
	}
	tplData["thirdGroupTitle"] = strings.TrimSpace(thirdData.Title)
	tplData["thirdSelectedOptions"] = thirdData.SelectedOptions
	var thirdSelectOptionsJSStr []template.JS
	for _, o := range thirdData.Options {
		thirdSelectOptionsJSStr = append(thirdSelectOptionsJSStr, template.JS(fmt.Sprintf("thirdGroupOptions['%s'] = [%s]", o.ParentID, o.OptionsList)))
	}
	tplData["thirdGroupSelect"] = thirdSelectOptionsJSStr

	return nil
}

func BuildThDataForWorkloadList(requestData, allOrderFields, allListItems map[string]string, tplData map[string]interface{}, defaultOrderField, defaultOrderDirection string) {

	selectedOrderField := strings.TrimSpace(strings.ToUpper(requestData["orderfield"]))
	if selectedOrderField == "" {
		selectedOrderField = defaultOrderField
	}
	selectOrderDirection := strings.TrimSpace(requestData["direction"])
	if selectOrderDirection == "" {
		selectOrderDirection = defaultOrderDirection
	}

	newOrderDirection := "0"
	if selectOrderDirection == "0" {
		newOrderDirection = "1"
	}

	thData := make([]ObjectTitle, 0)
	var thKeys []string
	for k, _ := range allListItems {
		thKeys = append(thKeys, k)
	}
	sort.Strings(thKeys)

	for _, v := range thKeys {
		lineData := ObjectTitle{}
		lineData.ID = v
		lineData.Title = allListItems[v]
		if isOrder(v, allOrderFields) {
			lineData.IsOrder = true
			if v == selectedOrderField {
				lineData.OrderSelected = "yes"
				lineData.OrderDirection = newOrderDirection
			} else {
				lineData.OrderSelected = ""
				lineData.OrderDirection = defaultOrderDirection
			}
		}

		thData = append(thData, lineData)
	}

	tplData["thData"] = thData
}

func ConvertMap2HTML(data map[string]string) template.HTML {
	ret := ""
	for k, v := range data {
		ret = ret + "<div>" + k + ":" + v + "</div>"
	}

	return template.HTML(ret)
}
