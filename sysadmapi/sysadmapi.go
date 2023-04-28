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

Ref: https://docs.docker.com/registry/spec/api/
	https://datatracker.ietf.org/doc/rfc7235/

	errorCode: 7000xxxx
*/

package sysadmapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"sysadm/httpclient"
	sysadm "sysadm/sysadm/server"
	"sysadm/sysadmerror"
	"sysadm/utils"
)

type ApiRequestData struct {
	ModuleName string
	ActionName string 
	Tls bool
	ApiServerAddress string
	ApiServerPort int
	ApiVersion string
	RequestParas httpclient.RequestParams
}

type ApiResponseStatus struct {
	// Status is false if this is a error response, otherwise Status is true
	Status bool `json:"status"`
	// Errorcode is zero if this is a successful response, otherwise Errorcode is nonzero
	Errorcode int `json:"errorcode"`
	// Message is the result sets if this is a successful ,otherwise Message is []map[string]string
	// which has one rows only:["errorMsg"] = errorMsg
	Message interface{} `json:"message"`
}

type ApiResponseData struct {
	// Status is false if this is a error response, otherwise Status is true
	Status bool 
	// Errorcode is zero if this is a successful response, otherwise Errorcode is nonzero
	Errorcode int 
	// DataSet is the result sets like  []map[string]string 
	// otherwise is DataSet["errorMsg"]
	DataSet []map[string]string 
}


/*
	BuildApiUrl: building the url for a API request.
	Parameters: moduleName module name; actionName: action name; tls: ApiServer where supports tls; apiServerAddress:apiServer address.it is the
		sysadm server address now;   apiServerPort: apiServer port.it is the sysadm server port now
	Return: apiUrl if built successful, otherwrise return "" and []sysadmerror.Sysadmerror
*/
func BuildApiUrl(moduleName string,actionName string, tls bool, apiServerAddress string, apiServerPort int,apiVersion string) (string, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	
	if moduleName == "" || actionName == "" || apiServerAddress == "" || apiVersion == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(70000001,"error","can not build apiUrl for moduleName, actionName, apiServerAddress or apiVersion is empty"))
		return "" ,errs
	} 

	found, err := CheckAction(moduleName,actionName)
	errs = append(errs, err...)
	if !found {
		return "" ,errs
	}

	m := sysadm.Modules
	module,_ := m[moduleName]
	
	var apiUrl string = ""
	if tls {
		if  apiServerPort  == 443 {
			apiUrl = "https://" + apiServerAddress  + "/api/" + apiVersion + "/" + module.Path +"/" + actionName 
		} else {
			apiUrl = "https://" + apiServerAddress + ":" + strconv.Itoa(apiServerPort) + "/api/" + apiVersion + "/" + module.Path +"/" + actionName
		}
	}else {
		if apiServerPort == 80 {
			apiUrl = "http://" + apiServerAddress + "/api/" + apiVersion  + "/" + module.Path + "/" + actionName		
		} else {
			apiUrl = "http://" + apiServerAddress + ":" + strconv.Itoa(apiServerPort) + "/api/" + apiVersion +  "/" + module.Path + "/" + actionName
		}
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(70000002,"debug","got apiURL for module(%s) and action(%s) is: %s",moduleName,actionName,apiUrl))

	return apiUrl ,errs

}


/*
    CheckAction: checks the existence of an action in a module. Returns:
	false, []sysadmerror.Sysadmerror if module or action is not exist 
	false, []sysadmerror.Sysadmerror if moduleName or actionName is empty
	true ,[]sysadmerror.Sysadmerror if the action is exist in  the module
*/
func CheckAction(moduleName string, actionName string)(bool,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if moduleName == "" || actionName == ""{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(70001001,"debug","can not check action for moduleName(%s) or actionName(%s) is empty",moduleName,actionName))
		return false ,errs
	}

	m := sysadm.Modules
	module,okModules := m[moduleName]
	if !okModules {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(70001002,"debug","moduleName(%s) is not exist",moduleName))
		return false ,errs
	}

	actions := module.Actions  
	found := false
	for _,value := range actions {
		if strings.EqualFold(value,actionName) {
			found = true
			break
		}
	}
	 
	return found, errs
}

/*
	SetHeaders set data to headers which will be set to the API server
*/
func SetHeaders(apiData *ApiRequestData, headersData map[string]string)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	headers := apiData.RequestParas.Headers
	for k,v := range headersData{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(70001003,"debug","add key %s value %s to headers",k,v))
		data := httpclient.RequestData{Key: k, Value: v}
		headers = append(headers, data)
	}

	apiData.RequestParas.Headers = headers

	return errs
}

/*
	SetQueryData set data to query which will be set to the API server
*/
func SetQueryData(apiData *ApiRequestData, data map[string]string)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	queryData := apiData.RequestParas.QueryData
	for k,v := range data {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(70001004,"debug","add key %s value %s to query",k,v))
		d := &httpclient.RequestData{Key: k, Value: v}
		queryData = append(queryData, d)
	}

	apiData.RequestParas.QueryData = queryData

	return errs
}

/*
	SetBasicAuthData set basic auth data to query which will be set to the API server
	basic auth data including:
	isBasicAuth: "true" or "false"  whether set basic auth data to api server
	username: user account which will be authorization
	password: user password which will be authorization
*/
func SetBasicAuthData(apiData *ApiRequestData, username string, password string, isBasicAuth bool)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	authData := apiData.RequestParas.BasicAuthData
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" || !isBasicAuth {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(70001005,"debug","basic Auth has be set to false."))
		authData["isBasicAuth"] = "false"
		apiData.RequestParas.BasicAuthData = authData
		return errs
	}

	authData["isBasicAuth"] = "true"
	authData["username"] = username
	authData["password"] = password
	apiData.RequestParas.BasicAuthData = authData
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(70001005,"debug","basic Auth has be set to true."))

	return errs
}

/*
	BuildApiRequestData: building the request data  for a API request.
	Parameters: moduleName module name; actionName: action name; tls: ApiServer where supports tls; apiServerAddress:apiServer address.it is the
		sysadm server address now;   apiServerPort: apiServer port.it is the sysadm server port now
	Return: *ApiRequestData if built successful, otherwrise return nil and []sysadmerror.Sysadmerror
*/
func BuildApiRequestData(moduleName string,actionName string, tls bool, apiServerAddress string, apiServerPort int,apiVersion string,method string)(*ApiRequestData,[]sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	var ret *ApiRequestData = nil

	reqUrl, err := BuildApiUrl(moduleName,actionName,tls,apiServerAddress,apiServerPort,apiVersion)
	errs = append(errs,err...)
	if reqUrl == "" {
		return ret,errs
	}

	retData := ApiRequestData{
		ModuleName: moduleName,
		ActionName: actionName,
		Tls: tls, 
		ApiServerAddress: apiServerAddress,
		ApiServerPort: apiServerPort,
		ApiVersion: apiVersion,
		RequestParas: httpclient.RequestParams{},
	}

	if method == "" {
		method = "GET"
	}

	retData.RequestParas.Url = reqUrl
	retData.RequestParas.Method = method
	ret = &retData

	return ret, errs
}

/*
	ParseResponseBody: parses the body content response from api server
	return ApiResponseData and []sysadmerror.Sysadmerror
	ApiResponseData.Status is ture and ApiResponseData.Errorcode is zero for successful response
	ApiResponseData.DataSet is the set of response like map[string]string if 

*/
func ParseResponseBody(body []byte)(ApiResponseData, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	var ret ApiResponseData = ApiResponseData{}
	if len(body) < 1 {
		dataSet := make(map[string]string,0)
		ret.Status = false
		ret.Errorcode = 70001005
		dataSet["errorMsg"] = "the response body from  the server is empty"
		ret.DataSet = append(ret.DataSet,dataSet)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(70001005,"error","the response body from  the server is empty"))
		return ret, errs
	}

	res := &ApiResponseStatus{}
	e := json.Unmarshal(body,res)
	if e != nil {
		dataSet := make(map[string]string,0)
		ret.Status = false
		ret.Errorcode = 70001006
		dataSet["errorMsg"] = fmt.Sprintf("can not parsing reponse body to json. error: %s",e) 
		ret.DataSet = append(ret.DataSet,dataSet)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(70001006,"error","can not parsing reponse body(%s) to json. error: %s",body,e))
		return ret , errs
	}

	if !res.Status {
		dataSet := make(map[string]string,0)
		ret.Status = false
		ret.Errorcode = res.Errorcode
		typeofMsg :=  reflect.TypeOf(res.Message)
		if typeofMsg.Kind() == reflect.String {
			dataSet["errorMsg"] = res.Message.(string)
		} else {
			retMsgArray := res.Message.([]interface {})
			if len(retMsgArray) > 0 {
				errMsgMap := retMsgArray[0].(map[string]interface {})
				errMsg := errMsgMap["errorMsg"].(string)
				dataSet["errorMsg"] = fmt.Sprintf("we got an error: %s",errMsg) 
			}else{
				dataSet["errorMsg"] = "we got an unknow error"
			}
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(res.Errorcode,"error","we got an error: %s",dataSet["errorMsg"]))
		ret.DataSet = append(ret.DataSet,dataSet)

		return ret, errs
	}
	
	var rets []map[string]string
	iData := res.Message.([]interface {})
	for _,iLine := range iData {
		d := iLine.(map[string]interface{})
		dataSet := make(map[string]string,0)
		for k,v := range d {
			value,ok := v.(string)
			if !ok {
				dataSet[k] = ""
			} else {
				vDecode, err := base64.StdEncoding.DecodeString(value)
				if err != nil { 
					dataSet[k] = value
				}else {
					dataSet[k] = utils.Bytes2str(vDecode) 
				}
			}
		}
		rets = append(rets, dataSet)
	}

	ret.Status = true
	ret.Errorcode = 0
	ret.DataSet = rets

	return ret,errs

}
