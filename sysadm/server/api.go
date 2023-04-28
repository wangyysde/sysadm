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
	//	"encoding/json"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"sysadm/httpclient"
	"sysadm/sysadmerror"
	"sysadm/utils"
	"github.com/wangyysde/sysadmServer"
)

// addFormHandler set delims for template and load template files
// return nil if not error otherwise return error.
func addApiHandler(r *sysadmServer.Engine,cmdRunPath string) {
	// Simple group: v1
	v1 := r.Group("/api/v1.0")
	{
		v1.POST("/:module/*action", apiHandlers)
//		v1.POST("/submit", submitEndpoint)
//		v1.POST("/read", readEndpoint)
	}

}

func apiHandlers(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror
	module := strings.TrimSuffix(strings.TrimPrefix(c.Param("module"),"/"),"/")
	action := strings.TrimSuffix(strings.TrimPrefix(c.Param("action"),"/"),"/")
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1030001,"debug","now handling the request for module %s with action %s.",module,action))
	if !foundModule(module) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1030002,"error","parameters error. module %s was not found.",module))
		logErrors(errs)
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1030002,
			Message: fmt.Sprintf("parameters error.module %s not found", module),
		}
		//respBody,_ := json.Marshal(ret)
		c.JSON(http.StatusOK, ret)
		return 
	}

	if !foundAction(module,action) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1030003,"error"," parameters error. action %s was not found in module %s.",action,module))
		logErrors(errs)
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1030003,
			Message: fmt.Sprintf("parameters error.action %s not found", action),
		}
		//respBody,_ := json.Marshal(ret)
		c.JSON(http.StatusOK, ret)
		return 
	}

	action =strings.ToLower(action)
	mI := Modules[module].Instance
	mI.ActionHanderCaller(action,c)

}

func foundModule(module string) bool {

	found := false
	for k := range Modules {
		if strings.EqualFold(k,module) {
			found = true
			break
		}
	}

	return found
}

func foundAction(m string, action string)bool{
	actions := Modules[m].Actions
	found := false
	for _,value := range actions {
		if strings.EqualFold(value,action) {
			found = true
			break
		}
	}
	 
	return found
}

func buildApiRequestParameters(module string,action string, data map[string] string,headers map[string]string,basicAuthData map[string]string)(*httpclient.RequestParams){
	var reqUrl string = ""
	m := Modules
	
	if RuntimeData.RuningParas.DefinedConfig.ApiServer.Tls {
		if  RuntimeData.RuningParas.DefinedConfig.ApiServer.Port  == 443 {
			reqUrl = "https://" + RuntimeData.RuningParas.DefinedConfig.ApiServer.Address  + "/api/" + apiVersion + "/" + m[module].Path +"/" + action 
		} else {
			reqUrl = "https://" + RuntimeData.RuningParas.DefinedConfig.ApiServer.Address + ":" + strconv.Itoa(RuntimeData.RuningParas.DefinedConfig.ApiServer.Port) + "/api/" + apiVersion + "/" + m[module].Path +"/" + action
		}
	}else {
		if RuntimeData.RuningParas.DefinedConfig.ApiServer.Port == 80 {
			reqUrl = "http://" + RuntimeData.RuningParas.DefinedConfig.ApiServer.Address + "/api/" + apiVersion  + "/" + m[module].Path + "/" + action		
		} else {
			reqUrl = "http://" + RuntimeData.RuningParas.DefinedConfig.ApiServer.Address + ":" + strconv.Itoa(RuntimeData.RuningParas.DefinedConfig.ApiServer.Port) + "/api/" + apiVersion +  "/" + m[module].Path + "/" + action
		}
	}
	var requestParams httpclient.RequestParams = httpclient.RequestParams{}
	requestParams.Url = reqUrl
	requestParams.Method = "POST"
	for k,v := range data {
		requestParams.QueryData = append(requestParams.QueryData,&httpclient.RequestData{Key: k, Value: v})
	}

	if headers != nil  {
		for k,v := range headers {
			requestParams.Headers  = append(requestParams.Headers,httpclient.RequestData{Key: k, Value: v})
		}
	} 

	if basicAuthData != nil {
		requestParams.BasicAuthData = basicAuthData
	}else {
		requestParams.BasicAuthData = nil
	}
	
	return &requestParams
}

func ParseResponseBody(body []byte)([]map[string]string, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1030004,"debug","try to parsing body %s",body))
	if len(body) < 1 {
		
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1030005,"error","the response from  the server is empty"))
		return nil, errs
	}

	res := &ApiResponseStatus{}
	e := json.Unmarshal(body,res)

	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1030006,"error","can not parsing reponse body to json. error: %s",e))
		return nil , errs
	}

	if !res.Status {
		retMsgArray := res.Message.([]interface {})
		if len(retMsgArray) > 0 {
			errMsgMap := retMsgArray[0].(map[string]interface {})
			errMsg := errMsgMap["errorMsg"].(string)
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(res.Errorcode,"error","we got an error: %s",errMsg))
		}else{
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(res.Errorcode,"error","we got an unknow error"))
		}

		return nil , errs
	}
	
	var rets []map[string]string
	iData := res.Message.([]interface {})
	for _,iLine := range iData {
		d := iLine.(map[string]interface{})
		ret := make(map[string]string)
		for k,v := range d {
			value,ok := v.(string)
			if !ok {
				ret[k] = ""
			} else {
				vDecode, err := base64.StdEncoding.DecodeString(value)
				if err != nil { 
					errs = append(errs, sysadmerror.NewErrorWithStringLevel(1030007,"error","decode field(%s)'s content error: %s",k,err))
					return nil,errs
				}
				ret[k] = utils.Bytes2str(vDecode)
			}
		}
		rets = append(rets, ret)
	}

	return rets,errs

}

func buildApiRequestUrl(module string,action string )(string){
	var reqUrl string = ""
	m := Modules
	
	if RuntimeData.RuningParas.DefinedConfig.ApiServer.Tls {
		if  RuntimeData.RuningParas.DefinedConfig.ApiServer.Port  == 443 {
			reqUrl = "https://" + RuntimeData.RuningParas.DefinedConfig.ApiServer.Address + "/api/" + apiVersion + "/" + m[module].Path +"/" + action 
		} else {
			reqUrl = "https://" + RuntimeData.RuningParas.DefinedConfig.ApiServer.Address + ":" + strconv.Itoa(RuntimeData.RuningParas.DefinedConfig.ApiServer.Port) + "/api/" + apiVersion + "/" + m[module].Path +"/" + action
		}
	}else {
		if RuntimeData.RuningParas.DefinedConfig.ApiServer.Port == 80 {
			reqUrl = "http://" + RuntimeData.RuningParas.DefinedConfig.ApiServer.Address + "/api/" + apiVersion  + "/" + m[module].Path + "/" + action		
		} else {
			reqUrl = "http://" + RuntimeData.RuningParas.DefinedConfig.ApiServer.Address + ":" + strconv.Itoa(RuntimeData.RuningParas.DefinedConfig.ApiServer.Port) + "/api/" + apiVersion +  "/" + m[module].Path + "/" + action
		}
	}

	return reqUrl
}

/* 
	status is false if this is a error response, otherwise Status is true
	errorCode is zero if this is a successful response, otherwise Errorcode is nonzero
	// Message is the result sets if this is a successful ,otherwise Message is []map[string]string
	// which has one rows only:["errorMsg"] = errorMsg
*/
func buildResponse(errorCode int, status bool,errMsg string)(ApiResponseStatus){
	var errMapArray []map[string]string
	msgMap := make(map[string]string)
	msgMap["errorMsg"] = errMsg
	errMapArray = append(errMapArray,msgMap)
	ret := ApiResponseStatus {
		Status: status,  
		Errorcode: errorCode,
		Message: errMapArray,
	}

	return ret
}