/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
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

package server

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sysadm/httpclient"
	"sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

// ErrorCode: 110xxxx

var projectUri = "/project/"

var projectTemplates = map[string] string {
	"list": "projectlist.html",
}

type projectDataStruct struct{
	actionHandler sysadmServer.HandlerFunc
	method []string
}

var projectData = map[string] projectDataStruct {
	"list":  {
		actionHandler: projectListHandler,
		method: []string{"GET"},
	},
	
}

// addFormHandler set delims for template and load template files
// return nil if not error otherwise return error.
func addProjectsHandler(r *sysadmServer.Engine,cmdRunPath string) ([]sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(110001,"debug","now adding handlers for project"))
	
	if RuntimeData.StartParas.SysadmRootPath  == "" {
		if _,err := getSysadmRootPath(cmdRunPath); err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(110002,"fatal","get the root path of the program error: %s",err))
			return errs
		}
	}

	for k,p := range projectData {
		
		listUri := projectUri + k
		for _,m := range p.method {
			switch m{
				case "GET":
					r.GET(listUri,p.actionHandler)
				case "POST":
					r.POST(listUri,p.actionHandler)
			}
		}
			
	}
	
	return errs
}

/*
	handler for handling list of the project
	Query parameters of request are below: 
	conditionKey: key name for DB query ,such as projectid, ownerid,name....
	conditionValue: the value of the conditionKey.for projectid, ownereid using =, for name, comment using like.
	deleted: 0 :normarl 1: deleted
	start: start number of the result will be returned.
	num: lines of the result will be returned.
*/
func projectListHandler(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(110004,"debug","now handling project list"))

	addProjectFormUrl := buildApiRequestUrl("project","add")
	delProjectFormUrl := buildApiRequestUrl("project","del")
	// checking user have login 
	userid,err := getSessionValue(c, "userid")
	if err != nil || userid == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(110005,"error","user should login"))
		var respHMTL string = "<div class=\"table-top\"> user should login </div>\n"
		tplData := map[string] interface{}{
			"noprojectinfo": template.HTML(respHMTL),
			"addProjectFormUrl": addProjectFormUrl,
			"delProjectFormUrl": delProjectFormUrl,
		}
		c.HTML(http.StatusOK,projectTemplates["list"], tplData)
		return 
	}
	
	// handling data
	conditionKey, _ := c.GetQuery("conditionKey")
	conditionValue, _ := c.GetQuery("conditionValue")
	conditionKey = strings.TrimSpace(conditionKey)
	conditionValue = strings.TrimSpace(conditionValue)
	data := make(map[string]string)
	numData := make(map[string]string)
	pageInfoParas := ""
	if conditionKey != "" && conditionValue != "" {
		data["conditionKey"] = conditionKey
		data["conditionValue"] = conditionValue
		numData["conditionKey"] = conditionKey
		numData["conditionValue"] = conditionValue
		pageInfoParas = pageInfoParas + "conditionKey=" + conditionKey
		pageInfoParas =pageInfoParas + "&conditionValue=" + conditionValue
	}

	deleted, _ := c.GetQuery("deleted")
	deleted = strings.ToLower(strings.TrimSpace(deleted))
	if deleted != "" {
		data["deleted"] = deleted
		numData["deleted"] = deleted
		if strings.TrimSpace(pageInfoParas) == ""{
			pageInfoParas = pageInfoParas + "deleted=" + deleted
		}else{
			pageInfoParas = pageInfoParas + "&deleted=" + deleted
		}
	}

	start, _ := c.GetQuery("start")
	start = strings.ToLower(strings.TrimSpace(start))
	if start != "" {
		data["start"] = start
	}else{
		data["start"] = "0"
	}
	data["num"] = strconv.Itoa(numPerPage)

	orderField, _ := c.GetQuery("orderfield")
	order, _ := c.GetQuery("order")
	orderField = strings.ToLower(strings.TrimSpace(orderField))
	order = strings.ToLower(strings.TrimSpace(order))
	if orderField != "" {
		data["orderField"] = orderField
		if strings.TrimSpace(pageInfoParas) == ""{
			pageInfoParas = pageInfoParas + "orderfield=" + orderField
		}else{
			pageInfoParas = pageInfoParas + "&orderfield=" + orderField
		}
	}
	if order != "" {
		if strings.TrimSpace(pageInfoParas) == ""{
			pageInfoParas = pageInfoParas + "order=" + order
		}else{
			pageInfoParas = pageInfoParas + "&order=" + order
		}
		data["order"] = order
	}

	requestParams := buildApiRequestParameters("project","list",data,nil,nil)
	requestNumParams := buildApiRequestParameters("project","getcount",numData,nil,nil)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(110003,"debug","try to execute the request with:%s",requestParams.Url))
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(110004,"debug","try to execute the request with:%s",requestNumParams.Url))
	body,e := httpclient.SendRequest(requestParams)
	numBody,numE := httpclient.SendRequest(requestNumParams)
	errs = append(errs, e...)
	errs = append(errs, numE...)
	logErrors(errs)
	errs = errs[0:0]
	ret, errs := ParseResponseBody(body)
	numRet, numErrs := ParseResponseBody(numBody)
	errs = append(errs, numErrs...)
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error"){
		var respHMTL string = "<div class=\"table-top\"> No Project or have an error </div>\n"
		tplData := map[string] interface{}{
			"noprojectinfo": template.HTML(respHMTL),
			"addProjectFormUrl": addProjectFormUrl,
			"delProjectFormUrl": delProjectFormUrl,
			"userid": template.HTML(userid.(string)),
		}
		c.HTML(http.StatusOK,projectTemplates["list"], tplData)
		return 
	}

	numLine := numRet[0]
	numStr := numLine["num"]
	total := 0
	if strings.TrimSpace(numStr) != "" {
		total,_ = strconv.Atoi(strings.TrimSpace(numStr))
	}
	
	var userIdListStr = ""
	first := true 
	for _,v := range ret {
		if first {
			userIdListStr += v["ownerid"]
			first = false
		}else {
			userIdListStr = userIdListStr + "," + v["ownerid"]
		}
	}
	if userIdListStr == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(110004,"error","ownerid of project no found"))
		logErrors(errs)
		tplData := map[string] interface{}{
			"noprojectinfo": template.HTML("<div class=\"table-top\">No Project or have an error</div>"),
			"addProjectFormUrl": addProjectFormUrl,
			"delProjectFormUrl": delProjectFormUrl,
			"userid": template.HTML(userid.(string)),
		}
		c.HTML(http.StatusOK,projectTemplates["list"], tplData)
		return 
	}

	userQuerydata := make(map[string]string)
	userQuerydata["userid"] = userIdListStr
	requestParams = buildApiRequestParameters("user","getinfo",userQuerydata,nil,nil)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(110005,"debug","try to execute the request with:%s",requestParams.Url))
	userbody,e := httpclient.SendRequest(requestParams)
	errs = append(errs, e...)
	logErrors(errs)
	userRet, errs := ParseResponseBody(userbody)
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error"){
		tplData := map[string] interface{}{
			"noprojectinfo": template.HTML("<div class=\"table-top\"> No Project or have an error</div>"),
			"addProjectFormUrl": addProjectFormUrl,
			"delProjectFormUrl": delProjectFormUrl,
			"userid": template.HTML(userid.(string)),
		}
		c.HTML(http.StatusOK,projectTemplates["list"], tplData)
		return 
	}

	usernames := make(map[string]string)
	for _,v := range userRet {
		usernames[v["userid"]] = v["username"]
	}
	
	var htmlData string = ""
	htmlData += "<form id=\"delProject\" method=\"post\" target=\"_self\" onsubmit=\"return false\">\n"
	htmlData = htmlData + "<table class=\"list-table\">\n"
	htmlData = htmlData + "<tr>\n"
	htmlData = htmlData + "<th width=\"5%\" align=\"left\">	<input type=\"checkbox\" id=\"projectidth\" name=\"projectid[]\" onclick='selectAllCheckbox()'></th>\n"
	htmlData = htmlData + "<th width=\"20%\"> 项目名称</th>\n"
	htmlData = htmlData + "<th width=\"10%\">所有者</th>\n"
	htmlData = htmlData + "<th width=\"10%\">删除状态</th>\n"
	htmlData = htmlData + "<th width=\"10%\">镜像数</th>\n"
	htmlData = htmlData + "<th width=\"10%\">创建时间</th>\n"
	htmlData = htmlData + "<th>描述</th>\n"
	htmlData = htmlData + "</tr>\n"
	for _,v := range ret {
		htmlData = htmlData + "<tr>\n"
		htmlData += "<td width=\"5%\">	<input type=\"checkbox\" id=\"projectid[]\" name=\"projectid[]\" value=\"" + v["projectid"] + "\" ></td>\n"
		htmlData += "<td>" + v["name"]+"</td>\n"
		username,ok := usernames[v["ownerid"]]
		if ok {
			htmlData += "<td>" + username + "</td>\n"
		} else {
			htmlData += "<td>  </td>\n"
		}
		deleted = ""
		if v["deleted"] == "1" {
			deleted = "删除"
		}else {
			deleted = "正常"
		}
		htmlData += "<td>" + deleted + "</td>\n"
		htmlData += "<td>10</td>\n"
		timeInt,_ := strconv.Atoi(v["creation_time"])
		timeInt64 := int64(timeInt)
		createTimeStamp := time.Unix(timeInt64,0)
		createTime := createTimeStamp.Format("2006-01-02 15:04:05")
		htmlData += "<td>" + createTime + "</td>\n"
		htmlData += "<td>" + v["comment"] + "</td>\n"
		htmlData += "</tr>\n"

	}
	htmlData += "</table>\n"
	htmlData += "</form>\n"
	pageStr := "<td ><div class=\"div-foot\">当前第"
	totalPages := int(math.Ceil(float64(total) / float64(numPerPage)))
	currentPage := 1
	startInt,_ := strconv.Atoi(start)
	currentPage = int(math.Ceil(float64(startInt + 1) / float64(numPerPage)))
	pageStr += strconv.Itoa(currentPage) + "页"
	if currentPage <= 1{
		pageStr += " 上一页 "
	}else {
		preNum := startInt - numPerPage
		prePage := fmt.Sprintf("?start=%d&num=%d&%s",preNum,numPerPage,pageInfoParas)
		pageStr = pageStr + "<a href=\"javascript:void(0)\" onclick='changePage(\"" + prePage + "\")'>上一页</a>"
	}

	if currentPage >= totalPages {
		pageStr += "下一页 "
	}else{
		nextNum := startInt + numPerPage
		nextPage := fmt.Sprintf("?start=%d&num=%d&%s",nextNum,numPerPage,pageInfoParas)
		pageStr = pageStr + "<a href=\"javascript:void(0)\" onclick='changePage(\"" + nextPage + "\")'>下一页</a>"
	}
	pageStr = pageStr + " 共" + strconv.Itoa(totalPages) + "页"

	htmlData += "<table class=\"foot-table\"><tr>\n"
	htmlData += "<td ><div class=\"div-foot\">" + pageStr + "</div></td></tr>\n"
	htmlData += "</table>\n"
	
	
	tplData := map[string] interface{}{
		"projectinfo": template.HTML(htmlData),
		"userid": template.HTML(userid.(string)),
		"addProjectFormUrl": addProjectFormUrl,
		"delProjectFormUrl": delProjectFormUrl,
	}
	c.HTML(http.StatusOK,projectTemplates["list"], tplData)
	
}


