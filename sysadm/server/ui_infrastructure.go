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
	"net/http"
	"strconv"

	"github.com/wangyysde/sysadm/httpclient"
	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
	"github.com/wangyysde/sysadmServer"
)

func addInfrastructureUIHandlers(r *sysadmServer.Engine)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030001,"fatal","can not add handlers to nil" ))
		return errs
	}

	r.GET("/infrastructure/list",infrastructureListHandler)

	return errs
}

/*
	handler for handling list of the infrastructure
	Query parameters of request are below: 
	conditionKey: key name for DB query ,such as projectid, ownerid,name....
	conditionValue: the value of the conditionKey.for projectid, ownereid using =, for name, comment using like.
	deleted: 0 :normarl 1: deleted
	start: start number of the result will be returned.
	num: lines of the result will be returned.
*/
func infrastructureListHandler(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030002,"debug","now handling infrastructure list"))

	// get userid
	userid,e := getSessionValue(c, "userid")
	if e != nil || userid == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030003,"error","user should login %s",e))
		logErrors(errs)
		tplData := map[string] interface{}{
			"errormessage": "user should login",
		}
		templateFile := "showmessage.html"
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

	// get yum information list
	moduleName := "yum"
	actionName := "yumlist"
	apiServerData := apiutils.BuildApiServerData(moduleName,actionName,apiVersion,tls,address,port,ca,cert,key)
	urlRaw, _ := apiutils.BuildApiUrl(apiServerData)

	requestParas :=  httpclient.RequestParams{}
	requestParas.Url = urlRaw
	requestParas.Method = http.MethodPost
	
	body,err := httpclient.SendRequest(&requestParas)
	errs = append(errs,err...)
	ret,err := apiutils.ParseResponseBody(body)
	errs = append(errs,err...)
	if !ret.Status {
		message := ret.Message
		messageLine := message[0]
		msg := utils.Interface2String(messageLine["msg"])
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030003,"error","get yum information list error %s",msg))
		logErrors(errs)
		tplData := map[string] interface{}{
			"errormessage": msg,
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK,templateFile, tplData)
		return 
	}

	osVerList := make(map[int]string,0)
	osVerMap := make(map[int]map[int]string,0)
	osNameList := make(map[int]string,0)
	yumInfoList := make(map[int]map[int]string,0)
	for _,line := range ret.Message {
		yumID, errYumID :=  utils.Interface2Int(line["yumid"])
		yumName := utils.Interface2String(line["name"])
		osID, errOSID := utils.Interface2Int(line["osid"])
		osName := utils.Interface2String(line["osName"])
		yumTypeName := utils.Interface2String(line["typeName"])
		yumCatalog := utils.Interface2String(line["catalog"])
		versionID, errVersionID := utils.Interface2Int(line["versionid"])
		versionName :=  utils.Interface2String(line["versionName"])
		if errYumID != nil || yumName == "" || errOSID != nil || errVersionID != nil {
			errMsg := fmt.Sprintf("get yum information error: %s %s %s", errYumID,errOSID,errVersionID)
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030004,"error",errMsg))
			logErrors(errs)
			tplData := map[string] interface{}{
				"errormessage": errMsg,
			}
			templateFile := "showmessage.html"
			c.HTML(http.StatusOK,templateFile, tplData)
			return 
		}

		var yumInfoListStr string = ""
		if osMap,ok := yumInfoList[osID]; ok {
			if versionMap,ok := osMap[versionID]; ok {
				yumInfo := versionMap[:(len(versionMap)-1)]
				yumInfoListStr = yumInfo + ",[" + strconv.Itoa(yumID) +",'" + yumName +"','" + yumTypeName + "','" + yumCatalog +"']]"
			} else {
				yumInfoListStr = "[[" + strconv.Itoa(yumID) +",'" + yumName +"','" + yumTypeName + "','" + yumCatalog +"']]"
			} 
		} else {
			yumInfoListStr = "[[" + strconv.Itoa(yumID) +",'" + yumName +"','" + yumTypeName + "','" + yumCatalog +"']]"
		}

		if verMap,ok := yumInfoList[osID]; ok {
			verMap[versionID] = yumInfoListStr
			yumInfoList[osID] = verMap
		} else {
			verMap := make(map[int]string,0)
			verMap[versionID] = yumInfoListStr
			yumInfoList[osID] = verMap
		}
		
		if verMap,ok := osVerMap[osID]; ok {
			verMap[versionID] = versionName
			osVerMap[osID] = verMap
		} else {
			verMap := make(map[int]string,0)
			verMap[versionID] = versionName
			osVerMap[osID] = verMap
		}
		osNameList[osID] = osName
	}

	for osID,verMap := range osVerMap{
		var osVer string = "["
		i := 0
		for verID,verName := range verMap {
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
	fmt.Printf("yumInfoList:%+v\n",yumInfoList)
	for osid,verMap := range yumInfoList {
		javascriptYumStr = append(javascriptYumStr, template.JS(fmt.Sprintf("yumList[%d] = new Array();", osid)))
		for verID,yumInfo :=  range verMap {
			javascriptYumStr = append(javascriptYumStr, template.JS(fmt.Sprintf("yumList[%d][%d] = %s;", osid,verID, yumInfo))) 
		}
		osLine := make(map[string]string,0)
		osLine["osid"] = strconv.Itoa(osid) 
		osLine["osname"] = osNameList[osid]
		osList = append(osList, osLine)
	}

	tplData := make(map[string]interface{},0)
	tplData["userid"] = userid
	tplData["yumList"] = javascriptYumStr

	var javascriptOsVersionStr []template.JS
	for osid, verStr := range osVerList {
		javascriptOsVersionStr = append(javascriptOsVersionStr, template.JS(fmt.Sprintf("osVerList[%d] = %s;", osid, verStr)))
	}
	tplData["osVerList"] = javascriptOsVersionStr
	tplData["osList"] = osList

	templateFile := "infrastructurelist.html"
	c.HTML(http.StatusOK,templateFile, tplData)
}	
	