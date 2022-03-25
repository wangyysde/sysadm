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
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/wangyysde/sysadm/httpclient"
	"github.com/wangyysde/sysadm/sysadmerror"
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
	
	// handling data
	conditionKey, _ := c.GetQuery("conditionKey")
	conditionValue, _ := c.GetQuery("conditionValue")
	conditionKey = strings.TrimSpace(conditionKey)
	conditionValue = strings.TrimSpace(conditionValue)
	data := make(map[string]string)
	if conditionKey != "" && conditionValue != "" {
		data["conditionKey"] = conditionKey
		data["conditionValue"] = conditionValue
	}

	deleted, _ := c.GetQuery("deleted")
	deleted = strings.ToLower(strings.TrimSpace(deleted))
	if deleted != "" {
		data["deleted"] = deleted
	}

	start, _ := c.GetQuery("start")
	num, _ := c.GetQuery("num")
	start = strings.ToLower(strings.TrimSpace(start))
	num = strings.ToLower(strings.TrimSpace(num))
	if start != "" {
		data["start"] = start
	}
	if num != "" {
		data["num"] = num
	}

	orderField, _ := c.GetQuery("orderfield")
	order, _ := c.GetQuery("order")
	orderField = strings.ToLower(strings.TrimSpace(orderField))
	order = strings.ToLower(strings.TrimSpace(order))
	if orderField != "" {
		data["orderField"] = orderField
	}
	if order != "" {
		data["order"] = order
	}

	requestParams := buildApiRequestParameters("project","list",data,nil,nil)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(110003,"debug","try to execute the request with:%s",requestParams.Url))
	body,err := httpclient.SendRequest(requestParams)
	errs = append(errs, err...)
	logErrors(errs)
	ret, errs := ParseResponseBody(body)
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error"){
		var respHMTL string = "<div class=\"table-top\"> No Project or have an error </div>\n"
		tplData := map[string] interface{}{
			"noprojectinfo": template.HTML(respHMTL),
		}
		c.HTML(http.StatusOK,projectTemplates["list"], tplData)
		return 
	}

	var userIdListStr = ""
	first := true 
	for _,v := range ret {
		if first {
			userIdListStr += v["ownerid"]
		}else {
			userIdListStr = userIdListStr + "," + v["ownerid"]
		}
	}
	if userIdListStr == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(110004,"error","ownerid of project no found"))
		logErrors(errs)
		tplData := map[string] interface{}{
			"noprojectinfo": template.HTML("<div class=\"table-top\">No Project or have an error</div>"),
		}
		c.HTML(http.StatusOK,projectTemplates["list"], tplData)
		return 
	}

	userQuerydata := make(map[string]string)
	userQuerydata["userid"] = userIdListStr
	requestParams = buildApiRequestParameters("user","getinfo",userQuerydata,nil,nil)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(110005,"debug","try to execute the request with:%s",requestParams.Url))
	userbody,err := httpclient.SendRequest(requestParams)
	errs = append(errs, err...)
	logErrors(errs)
	userRet, errs := ParseResponseBody(userbody)
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error"){
		tplData := map[string] interface{}{
			"noprojectinfo": template.HTML("<div class=\"table-top\"> No Project or have an error</div>"),
		}
		c.HTML(http.StatusOK,projectTemplates["list"], tplData)
		return 
	}

	usernames := make(map[string]string)
	for _,v := range userRet {
		usernames[v["userid"]] = v["username"]
	}
	
	var htmlData string = ""
	htmlData = htmlData + "<table class=\"list-table\">\n"
	htmlData = htmlData + "<tr>\n"
	htmlData = htmlData + "<th width=\"5%\" align=\"left\">	<input type=\"checkbox\" id=\"projectid[]\" name=\"projectid[]\"></th>\n"
	htmlData = htmlData + "<th width=\"20%\"> 项目名称</th>\n"
	htmlData = htmlData + "<th width=\"10%\">所有者</th>\n"
	htmlData = htmlData + "<th width=\"10%\">删除状态</th>\n"
	htmlData = htmlData + "<th width=\"10%\">镜像数</th>\n"
	htmlData = htmlData + "<th width=\"10%\">创建时间</th>\n"
	htmlData = htmlData + "<th>描述</th>\n"
	htmlData = htmlData + "</tr>\n"
	for _,v := range ret {
		htmlData = htmlData + "<tr>\n"
		htmlData += "<td width=\"5%\">	<input type=\"checkbox\" id=\"projectid[]\" name=\"projectid[]\"></td>\n"
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
	htmlData += "<tr>\n"
	htmlData += "<td colspan=8 ><div class=\"td-foot\">当前 第1页 上一页 下一页 共10页</div></td>\n"
	htmlData += "</tr>\n"
	htmlData += "</table>\n"
	tplData := map[string] interface{}{
		"projectinfo": template.HTML(htmlData),
	}
	c.HTML(http.StatusOK,projectTemplates["list"], tplData)
	
}
