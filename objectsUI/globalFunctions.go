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
	subCategory = strings.TrimSpace(mainCategory)
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
