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

// ErrorCode: 111xxxx
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

	"github.com/wangyysde/sysadm/httpclient"
	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
	"github.com/wangyysde/sysadmServer"
)

var registryctlUri = "/registryctl/"

var registryctlActionsHandlers []actionHandler

// addFormHandler set delims for template and load template files,add handlers according registryctlActionsHandlers
// return nil if not error otherwise return error.
func addRegistryctlHandler(r *sysadmServer.Engine,cmdRunPath string) ([]sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	
	registryctlActionsHandlers = []actionHandler{
		{name: "imagelist", templateFile: "imagelist.html", handler: imageListHandler,method: []string{"GET", "POST"}},
		{name: "taglist", templateFile: "taglist.html", handler: tagListHandler,method: []string{"GET"}},
	}
	
	if RuntimeData.StartParas.SysadmRootPath  == "" {
		if _,err := getSysadmRootPath(cmdRunPath); err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(1110001,"fatal","get the root path of the program error: %s",err))
			return errs
		}
	}

	for _,v := range registryctlActionsHandlers {

		handlerUrl := registryctlUri + v.name
		for _,m := range v.method  {
			switch m{
				case "GET":
					r.GET(handlerUrl,v.handler )
				case "POST":
					r.POST(handlerUrl,v.handler)
				case "HEAD":
					r.HEAD(handlerUrl,v.handler)
				case "PUT":
					r.PUT(handlerUrl,v.handler)
				case "DELETE":
					r.DELETE(handlerUrl,v.handler)
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
func imageListHandler(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101004,"debug","now handling project list"))

	// get template file name 
	templateFile := ""
	for _,v := range registryctlActionsHandlers {
		if v.name == "imagelist" {
			templateFile = v.templateFile
		}
	}	
	
	// get userid
	userid,e := getSessionValue(c, "userid")
	if e != nil || userid == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101010,"error","user should login %s",e))
		logErrors(errs)
		tplData := map[string] interface{}{
			"errormessage": "user should login",
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}

	// get project for select menu
	moduleName := "project"
	actionName := "list"
	definedConfig := RuntimeData.RuningParas.DefinedConfig
	apiVersion := definedConfig.ApiServer.ApiVersion
	tls := definedConfig.ApiServer.Tls
	address := definedConfig.ApiServer.Address
	port := definedConfig.ApiServer.Port
	ca := definedConfig.ApiServer.Ca
	cert := definedConfig.ApiServer.Cert
	key := definedConfig.ApiServer.Key
	apiServerData := apiutils.BuildApiServerData(moduleName,actionName,apiVersion,tls,address,port,ca,cert,key)
	if apiServerData == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101005,"error","api server parameters error"))
		logErrors(errs)
		tplData := map[string] interface{}{
			"errormessage": "api server parameters error",
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}
	
	urlRaw, err := apiutils.BuildApiUrl(apiServerData)
	errs = append(errs,err...)
	if urlRaw == "" {
		err := apiutils.SendResponseForErrorMessage(c,1101007, "api server parameters error")
		errs = append(errs, err...)
		logErrors(errs)
		tplData := map[string] interface{}{
			"errormessage": "api server parameters error",
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}

	requestParas :=  httpclient.RequestParams{}
	requestParas.Url = urlRaw
	requestParas.Method = http.MethodPost
	body,err := httpclient.SendRequest(&requestParas)
	errs = append(errs,err...)
	ret,err := apiutils.ParseResponseBody(body)
	errs = append(errs,err...)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101004,"debug","now handling project list %#v",ret))
	logErrors(errs)
	if !ret.Status {
		errCode := ret.ErrorCode
		msgLines := ret.Message
		msgLine := msgLines[0]
		errMsg := msgLine["msg"].(string)
		tplData := map[string] interface{}{
			"errormessage": fmt.Sprintf("errorCode: %d Msg: %s",errCode,errMsg),
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}
	
	// preparing project select data
	var projectInfo []map[string]string
	projectInfo = append(projectInfo,map[string]string{"0":"全部项目"})
	res := ret.Message
	for _,line := range res{
		lineMap := make(map[string]string,0)
		id := utils.Interface2String(line["projectid"])
		name := utils.Interface2String(line["name"])
		lineMap[id] = name
		projectInfo = append(projectInfo,lineMap)
	}

	// get parameters on connection 
	queryData, _ := utils.GetRequestData(c,[]string{"projectid","searchKey","start","numPerPage"})
	startStr,ok := queryData["start"]
	if !ok {
		startStr = "0"
	}
	start,_ := strconv.Atoi(startStr)

	// get total rows according parametes
	moduleName = "registryctl"
	actionName = "getcount"
	apiVersion = definedConfig.Registryctl.ApiVersion
	tls = definedConfig.Registryctl.Tls
	address = definedConfig.Registryctl.Address
	port = definedConfig.Registryctl.Port
	ca = definedConfig.Registryctl.Ca
	cert = definedConfig.Registryctl.Cert
	key = definedConfig.Registryctl.Key
	apiServerData = apiutils.BuildApiServerData(moduleName,actionName,apiVersion,tls,address,port,ca,cert,key)
	if apiServerData == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101011,"error","api server parameters error"))
		err := apiutils.SendResponseForErrorMessage(c,1101011, "api server parameters error")
		errs = append(errs,err...)
		logErrors(errs)
		tplData := map[string] interface{}{
			"errormessage": "api server parameters error",
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}
	urlRaw, err = apiutils.BuildApiUrl(apiServerData)
	errs = append(errs,err...)
	if urlRaw == "" {
		err := apiutils.SendResponseForErrorMessage(c,1101012, "api server parameters error")
		errs = append(errs, err...)
		logErrors(errs)
		tplData := map[string] interface{}{
			"errormessage": "api server parameters error",
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}

	requestParas =  httpclient.RequestParams{}
	requestParas.Url = urlRaw
	requestParas.Method = http.MethodPost
	
	projectid,ok := queryData["projectid"]
	pageInfoParas := ""
	if ok && projectid != "0" {
		requestParasPr,err := httpclient.AddQueryData(&requestParas,"projectid",projectid)
		requestParas = *requestParasPr
		pageInfoParas = "projectid=" + projectid
		errs=append(errs,err...)
	}
	name,ok := queryData["searchKey"]
	if ok {
		requestParasPr,err := httpclient.AddQueryData(&requestParas,"name",name)
		requestParas = *requestParasPr
		if pageInfoParas == ""{
			pageInfoParas = "searchKey=" + name
		}else {
			pageInfoParas = pageInfoParas + "&searchKey=" + name
		}
		errs=append(errs,err...)
	}
	
	body,err = httpclient.SendRequest(&requestParas)
	errs = append(errs,err...)
	ret,err = apiutils.ParseResponseBody(body)
	errs = append(errs,err...)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101004,"debug","now handling image list %#v",ret))
	logErrors(errs)
	if !ret.Status {
		errCode := ret.ErrorCode
		msgLines := ret.Message
		msgLine := msgLines[0]
		errMsg := msgLine["msg"].(string)
		tplData := map[string] interface{}{
			"errormessage": fmt.Sprintf("errorCode: %d Msg: %s",errCode,errMsg),
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}
	msg := ret.Message
	line := msg[0]
	numInterface,ok := line["num"]
	if !ok {
		tplData := map[string] interface{}{
			"errormessage": "can not got total number of images",
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}

	// get total number and calculate total paget number
	numStr := utils.Interface2String(numInterface)
	num,_ := strconv.Atoi(numStr)
	if num < 1 {
		num = 0
	}
	totalPages := int(math.Ceil(float64(num) / float64(numPerPage)))
	currentPage := 1
	currentPage = int(math.Ceil(float64(start + 1) / float64(numPerPage)))
	currentPageHTML := strconv.Itoa(currentPage)
	totalPageHTML := strconv.Itoa(totalPages)
	prePageHTML := ""
	if currentPage <= 1{
		prePageHTML = " 上一页 "
	}else {
		preNum := start - numPerPage
		prePage := fmt.Sprintf("?start=%d&numPerPage=%d&%s",preNum,numPerPage,pageInfoParas)
		prePageHTML = "<a href=\"javascript:void(0)\" onclick='changePage(\"" + prePage + "\")'>上一页</a>"
	}

	nextPageHTML := ""
	if currentPage >= totalPages{
		nextPageHTML = "下一页 "
	}else{
		nextNum := start + numPerPage
		nextPage := fmt.Sprintf("?start=%d&num=%d&%s",nextNum,numPerPage,pageInfoParas)
		nextPageHTML = "<a href=\"javascript:void(0)\" onclick='changePage(\"" + nextPage + "\")'>下一页</a>"
	}
	

	moduleName = "registryctl"
	actionName = "imagelist"
	apiServerData = apiutils.BuildApiServerData(moduleName,actionName,apiVersion,tls,address,port,ca,cert,key)

	if apiServerData == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101010,"error","api server parameters error"))
		err := apiutils.SendResponseForErrorMessage(c,1101011, "api server parameters error")
		errs = append(errs,err...)
		logErrors(errs)
		tplData := map[string] interface{}{
			"errormessage": "api server parameters error",
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}
	
	urlRaw, err = apiutils.BuildApiUrl(apiServerData)
	errs = append(errs,err...)
	if urlRaw == "" {
		err := apiutils.SendResponseForErrorMessage(c,1101007, "api server parameters error")
		errs = append(errs, err...)
		logErrors(errs)
		tplData := map[string] interface{}{
			"errormessage": "api server parameters error",
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}

	requestParas.Url = urlRaw
	requestParas.Method = http.MethodPost
	requestParasPr,err := httpclient.AddQueryData(&requestParas,"start",startStr)
	requestParas = *requestParasPr
	errs=append(errs,err...)
	requestParasPr,err = httpclient.AddQueryData(&requestParas,"numperpage",strconv.Itoa(numPerPage))
	requestParas = *requestParasPr
	errs=append(errs,err...)
	
	body,err = httpclient.SendRequest(&requestParas)
	errs = append(errs,err...)
	ret,err = apiutils.ParseResponseBody(body)
	errs = append(errs,err...)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101004,"debug","now handling project list %#v",ret))
	logErrors(errs)
	if !ret.Status {
		errCode := ret.ErrorCode
		msgLines := ret.Message
		msgLine := msgLines[0]
		errMsg := msgLine["msg"].(string)
		tplData := map[string] interface{}{
			"errormessage": fmt.Sprintf("errorCode: %d Msg: %s",errCode,errMsg),
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}
	
	type imageData struct {
		Id string
		ProjectName string
		ImageName string
		LastTag string
		PullTimes string
		UpdateTime string 
		Size string
	}

	msg = ret.Message
	var imageList []imageData
	for  _,lineMap := range msg {
		id := utils.Interface2String(lineMap["imageid"])
		imageName := utils.Interface2String(lineMap["name"])
		imageNameArray := strings.Split(imageName,"/")
		projectName := imageNameArray[0]
		lastTag := utils.Interface2String(lineMap["lasttag"])
		pullTimes := utils.Interface2String(lineMap["pulltimes"])
		timeInt,_ := strconv.Atoi(utils.Interface2String(lineMap["update_time"]))
		timeInt64 := int64(timeInt)
		createTimeStamp := time.Unix(timeInt64,0)
		updateTime := createTimeStamp.Format("2006-01-02 15:04:05")
		size, _ := strconv.Atoi(utils.Interface2String(lineMap["size"]))
		s := int((float64(size))/1024/1024)
		sizeStr := fmt.Sprintf("%dMiB",s)
		lineData := imageData{
			Id: id,
			ProjectName: projectName,
			ImageName: imageName,
			LastTag: lastTag,
			PullTimes: pullTimes,
			UpdateTime: updateTime,
			Size: sizeStr,
		}
		imageList = append(imageList, lineData)
	}
	tplData := map[string] interface{}{
		"projectinfo": projectInfo,
		"imagelist": imageList,
		"currentpage": currentPageHTML,
		"totalpage": totalPageHTML,
		"prepage":  template.HTML(prePageHTML),
		"nextpage":  template.HTML(nextPageHTML),
		"userid": template.HTML(userid.(string)),
		"selectedprojectid": projectid,
	}

	c.HTML(http.StatusOK,templateFile, tplData)
	
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
func tagListHandler(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror

	// get template file name 
	templateFile := ""
	for _,v := range registryctlActionsHandlers {
		if v.name == "taglist" {
			templateFile = v.templateFile
		}
	}
	
		// get userid
	userid,e := getSessionValue(c, "userid")
	if e != nil || userid == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101015,"error","user should login %s",e))
		logErrors(errs)
		tplData := map[string] interface{}{
			"errormessage": "user should login",
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}

	// get parameters on connection 
	queryData, _ := utils.GetRequestData(c,[]string{"imageid"})
	imageid,okImageid := queryData["imageid"]
	if !okImageid {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101020,"error","can not get image information"))
		logErrors(errs)
		tplData := map[string] interface{}{
			"errormessage": "can not get image information",
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}

	definedConfig := RuntimeData.RuningParas.DefinedConfig
	apiVersion := definedConfig.Registryctl.ApiVersion 
	tls := definedConfig.Registryctl.Tls
	address := definedConfig.Registryctl.Address
	port := definedConfig.Registryctl.Port
	ca := definedConfig.Registryctl.Ca
	cert := definedConfig.Registryctl.Cert
	key := definedConfig.Registryctl.Key
	moduleName := "registryctl"
	actionName := "imagelist"
	apiServerData := apiutils.BuildApiServerData(moduleName,actionName,apiVersion,tls,address,port,ca,cert,key)
	urlRaw, _ := apiutils.BuildApiUrl(apiServerData)

	requestParas :=  httpclient.RequestParams{}
	requestParas.Url = urlRaw
	requestParas.Method = http.MethodPost
	requestParaPr,err := httpclient.AddQueryData(&requestParas,"imageid",imageid)
	requestParas = *requestParaPr
	errs=append(errs, err...)

	body,err := httpclient.SendRequest(&requestParas)
	errs = append(errs,err...)
	ret,err := apiutils.ParseResponseBody(body)
	errs = append(errs,err...)
	if !ret.Status {
		message := ret.Message
		messageLine := message[0]
		msg := messageLine["msg"]
		tplData := map[string] interface{}{
			"errormessage": msg,
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}

	message := ret.Message
	imageLine := message[0]	
	imageName := utils.Interface2String(imageLine["name"])
	
	moduleName = "registryctl"
	actionName = "taglist"
	apiServerData = apiutils.BuildApiServerData(moduleName,actionName,apiVersion,tls,address,port,ca,cert,key)
	
	urlRaw, err = apiutils.BuildApiUrl(apiServerData)
	errs = append(errs,err...)
	
	requestParas =  httpclient.RequestParams{}
	requestParas.Url = urlRaw
	requestParas.Method = http.MethodPost
	requestParaPr,err = httpclient.AddQueryData(&requestParas,"imageid",imageid)
	requestParas = *requestParaPr

	body,err = httpclient.SendRequest(&requestParas)
	errs = append(errs,err...)
	ret,err = apiutils.ParseResponseBody(body)
	errs = append(errs,err...)
	logErrors(errs)
	if !ret.Status {
		errCode := ret.ErrorCode
		msgLines := ret.Message
		msgLine := msgLines[0]
		errMsg := msgLine["msg"].(string)
		tplData := map[string] interface{}{
			"errormessage": fmt.Sprintf("errorCode: %d Msg: %s",errCode,errMsg),
		}
		templateFile = "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}
	
	type tagData struct {
		Id string
		Name string
		Description string
		Pulltimes string
		CreateTime string
		UpdateTime string 
		Size string
		Digest string
	}

	msg := ret.Message
	var tagList []tagData
	for  _,lineMap := range msg {
		id := utils.Interface2String(lineMap["tagid"])
		tagName := utils.Interface2String(lineMap["name"])
		description := utils.Interface2String(lineMap["description"])
		pulltimes := utils.Interface2String(lineMap["pulltimes"])
		timeInt,_ := strconv.Atoi(utils.Interface2String(lineMap["creation_time"]))
		timeInt64 := int64(timeInt)
		createTimeStamp := time.Unix(timeInt64,0)
		creationTime := createTimeStamp.Format("2006-01-02 15:04:05")
		timeInt,_ = strconv.Atoi(utils.Interface2String(lineMap["update_time"]))
		timeInt64 = int64(timeInt)
		createTimeStamp = time.Unix(timeInt64,0)
		updateTime := createTimeStamp.Format("2006-01-02 15:04:05")
		size, _ := strconv.Atoi(utils.Interface2String(lineMap["size"]))
		s := int((float64(size))/1024/1024)
		sizeStr := fmt.Sprintf("%dMiB",s)
		digest := utils.Interface2String(lineMap["digest"])
		lineData := tagData{
			Id: id,
			Name: tagName,
			Description: description,
			Pulltimes: pulltimes,
			CreateTime: creationTime,
			UpdateTime: updateTime,
			Size: sizeStr,
			Digest: digest,
		}
		tagList = append(tagList, lineData)
	}
	tplData := map[string] interface{}{
		"taglist": tagList,
		"imageid": imageid,
		"userid": template.HTML(userid.(string)),
		"imagename": template.HTML(imageName),
	}

	c.HTML(http.StatusOK,templateFile, tplData)
	
}