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

	errorCode: 8000xxxx
*/

package apiutils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"sysadm/httpclient"
	"sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

/*
	BuildResponseDataForString: build response data for error message and successful message
	ApiResponseData.Status is false if this is a error response, otherwise Status is true
	ApiResponseData.ErrorCode is zero if this is a successful response, otherwise Errorcode is nonzero
	ApiResponseData.Message is the result sets if this is a successful ,otherwise Message is []map[string]interface
	   For success: Message is the data sets like []map[string]interface{}
*/
func BuildResponseDataForString(errorCode int, msg string)(ApiResponseData){
	
	var status bool = false
	if errorCode == 0 {
		status = true
		if len(strings.TrimSpace(msg)) < 1 {
			msg = "successful"
		}
	} else {
		if len(strings.TrimSpace(msg)) < 1 {
			msg = "unknow error"
		}
	} 

	encodeMsg := base64.StdEncoding.EncodeToString([]byte(msg))
	var dataSets []map[string]interface{}
	data := make(map[string]interface{},0)
	data["msg"] = encodeMsg
	dataSets = append(dataSets,data)

	retData := ApiResponseData{
		Status: status,
		ErrorCode: errorCode,
		Message: dataSets,
	}

	return retData
}

/* 
	BuildResponseDataForMap: build response data for data sets
	ApiResponseData.Status is true
	ApiResponseData.ErrorCode is zero 
	ApiResponseData.Message is the result sets 
*/
func BuildResponseDataForMap(data []map[string]interface{}) (ApiResponseData){
	
	retData := ApiResponseData{
		Status: true,
		ErrorCode: 0,
		Message: data,
	}

	return retData
}

/* 
	NewBuildResponseDataForMap: build response data for data sets
	ApiResponseData.Message is the result sets 
*/
func NewBuildResponseDataForMap(status bool, errorCode int, data []map[string]interface{}) (ApiResponseData){
	if status {
		errorCode = 0
	}

	retData := ApiResponseData{
		Status: status,
		ErrorCode: errorCode,
		Message: data,
	}

	return retData
}



/*
	BuildResponseDataForSuccess build response data for string with successful 
*/
func BuildResponseDataForSuccess(msg string)(ApiResponseData){
	return BuildResponseDataForString(0,msg)
}

/*
	BuildResponseDataForSuccess build response data for string with error
*/
func BuildResponseDataForError(errorCode int, msg string)(ApiResponseData){
	return BuildResponseDataForString(errorCode,msg)
}

/*
	SendResponseForErrorMessage: build response data using errorCode and msg. 
	then send the response data to the client with json format.
*/
func SendResponseForErrorMessage(c *sysadmServer.Context,errorCode int, msg string)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	retData := BuildResponseDataForError(errorCode,msg)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010001,"debug","set errorCode %d message: %s to the client",errorCode,msg))
	
	c.JSON(http.StatusOK, retData)

	return errs
}

/*
	SendResponseForSuccessMessage: build response data using  msg. 
	then send the response data to the client with json format.
*/
func SendResponseForSuccessMessage(c *sysadmServer.Context,msg string)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	retData := BuildResponseDataForSuccess(msg)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010002,"debug","set successful message: %s to the client",msg))
	
	c.JSON(http.StatusOK, retData)

	return errs
}


/*
	NewSendResponseForErrorMessage: build response data using  msg ,value ,errorCode
	then send the response data to the client with specified format. 
	if format is not one of the json, html,xml, yaml, then default format is json
*/
func NewSendResponseForErrorMessage(c *sysadmServer.Context,headers http.Header,status int,format string, errorCode int, msg string,value ...interface{})([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	for key,value := range headers {
		if len(value) < 1 {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010022,"debug","delete  key %s from header",key))
			c.Header(key,"")
		} else {
			for _,v := range value {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010023,"debug","add key %s value %s to http header",key,value))
				c.Header(key,v)
			}
		}
	}

	if http.StatusText(status) == "" {
		status = http.StatusOK
	}

	retMsg := fmt.Sprintf(msg,value...)
	switch strings.ToLower(format){
	case "json":
		retData := BuildResponseDataForError(errorCode,retMsg)
		c.JSON(status, retData)
	case "html":
		c.Header("Content-Type","text/html")
		c.String(status,msg,value...)
	case "xml":
		retData := BuildResponseDataForSuccess(retMsg)
		c.XML(status,retData)
	case "yaml":
		retData := BuildResponseDataForSuccess(retMsg)
		c.YAML(status,retData)
	default:
		retData := BuildResponseDataForSuccess(retMsg)
		c.JSON(status, retData)
	}

	return errs
}

/*
	NewSendResponseForSuccessMessage: build response data using  msg and value 
	then send the response data to the client with specified format. 
	if format is not one of the json, html,xml, yaml, then default format is json
*/
func NewSendResponseForSuccessMessage(c *sysadmServer.Context,headers http.Header,status int,format string,msg string, value ...interface{})([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	for key,value := range headers {
		if len(value) < 1 {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010020,"debug","delete  key %s from header",key))
			c.Header(key,"")
		} else {
			for _,v := range value {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010021,"debug","add key %s value %s to http header",key,value))
				c.Header(key,v)
			}
		}
	}

	if http.StatusText(status) == "" {
		status = http.StatusOK
	}

	retMsg := fmt.Sprintf(msg,value...)
	switch strings.ToLower(format){
	case "json":
		retData := BuildResponseDataForSuccess(retMsg)
		c.JSON(status, retData)
	case "html":
		c.Header("Content-Type","text/html")
		c.String(status,msg,value...)
	case "xml":
		retData := BuildResponseDataForSuccess(retMsg)
		c.XML(status,retData)
	case "yaml":
		retData := BuildResponseDataForSuccess(retMsg)
		c.YAML(status,retData)
	default:
		retData := BuildResponseDataForSuccess(retMsg)
		c.JSON(status, retData)
	}

	return errs
}
/*
	SendResponseForMap: build response data using  map. 
	then send the response data to the client with json format.
*/
func SendResponseForMap(c *sysadmServer.Context,data []map[string]interface{})([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	retData := BuildResponseDataForMap(data)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010003,"debug","send data: %#v to the client",data))
	
	c.JSON(http.StatusOK, retData)

	return errs
}

/*
	BuildApiUrl: building the url for a API request.
	Parameters: moduleName module name; actionName: action name; tls: ApiServer where supports tls; apiServerAddress:apiServer address.it is the
		sysadm server address now;   apiServerPort: apiServer port.it is the sysadm server port now
	Return: apiUrl(like: http://data.Server.Address:data.Server.Port/api/data.ApiVersion/data.ModuleName/data.ActionName) returned, if built successful, otherwrise return "" and []sysadmerror.Sysadmerror
*/
func BuildApiUrl(data *ApiServerData)  (string, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	
	if data == nil  || data.ModuleName == "" || data.ActionName == "" || data.ApiVersion == "" || data.Server.Address == "" || data.Server.Port == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010004,"error","parameters error. module name, action, apiVersion or server is empty."))
		return "",errs
	}

	var apiUrl string = ""
	if data.Server.Tls {
		if  data.Server.Port  == 443 {
			apiUrl = "https://" + data.Server.Address + "/api/" + data.ApiVersion + "/" + data.ModuleName +"/" + data.ActionName 
		} else { 
			apiUrl = "https://" + data.Server.Address + ":" + strconv.Itoa(data.Server.Port) + "/api/" + data.ApiVersion + "/" + data.ModuleName +"/" + data.ActionName
		}
	}else {
		if data.Server.Port == 80 {
			apiUrl = "http://" + data.Server.Address + "/api/" + data.ApiVersion  + "/" + data.ModuleName + "/" + data.ActionName 		
		} else {
			apiUrl = "http://" + data.Server.Address + ":" + strconv.Itoa(data.Server.Port) + "/api/" + data.ApiVersion +  "/" + data.ModuleName + "/" + data.ActionName
		}
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010005,"debug","proxy destination URL is: %s",apiUrl))

	return apiUrl ,errs

}


func BuildReverseProxyDirector(c *sysadmServer.Context,data *ApiServerData)(func(r *http.Request), []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
		
	rawURL,err := BuildApiUrl(data)
	errs = append(errs,err...)
	if rawURL == "" {
		return nil,errs
	}
	
	url := c.Request.URL
	query := url.RawQuery
	rawURL = rawURL + "?" + query
	url,e := url.Parse(rawURL)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010006,"error","parse proxy url(%s) error: %s",rawURL,e))
		return nil, errs
	}
	
	host := data.Server.Address + ":" + strconv.Itoa(data.Server.Port)

	return func(r *http.Request) {
		r.Host = host
		r.URL = url
	},errs
}

/*
	BuildApiUrl: building the url for a API request.
	Parameters: moduleName module name; actionName: action name; tls: ApiServer where supports tls; apiServerAddress:apiServer address.it is the
		sysadm server address now;   apiServerPort: apiServer port.it is the sysadm server port now
	Return: apiUrl(like: http://data.Server.Address:data.Server.Port/api/data.ApiVersion/data.ModuleName/data.ActionName) returned, if built successful, otherwrise return "" and []sysadmerror.Sysadmerror
*/
func NewBuildApiUrl(r *http.Request,data *ApiServerData)  (string, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	
	if data == nil  || data.ModuleName == "" || data.ActionName == "" || data.ApiVersion == "" || data.Server.Address == "" || data.Server.Port == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(800020004,"error","parameters error. parameters is %#v",data))
		return "",errs
	}

	var apiUrl string = ""
	if data.Server.Tls {
		if  data.Server.Port  == 443 {
			apiUrl = "https://" + data.Server.Address + "/api/" + data.ApiVersion + "/" + data.ModuleName +"/" + data.ActionName 
		} else { 
			apiUrl = "https://" + data.Server.Address + ":" + strconv.Itoa(data.Server.Port) + "/api/" + data.ApiVersion + "/" + data.ModuleName +"/" + data.ActionName
		}
	}else {
		if data.Server.Port == 80 {
			apiUrl = "http://" + data.Server.Address + "/api/" + data.ApiVersion  + "/" + data.ModuleName + "/" + data.ActionName 		
		} else {
			apiUrl = "http://" + data.Server.Address + ":" + strconv.Itoa(data.Server.Port) + "/api/" + data.ApiVersion +  "/" + data.ModuleName + "/" + data.ActionName
		}
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(800020005,"debug","proxy destination URL is: %s",apiUrl))

	url,e := url.Parse(apiUrl)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(800020005,"error","parse proxy url(%s) error: %s",apiUrl,e))
		return "", errs
	}
	
	r.URL = url

	return apiUrl ,errs

}

func NewBuildApiHost(r *http.Request,data *ApiServerData)(string, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	address := data.Server.Address
	if address == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(800020006,"error","destination server is empty"))
		return "",errs
	}
	port := data.Server.Port
	if port == 0 {
		if data.Server.Tls {
			port = 443 
		} else {
			port = 80
		}
	}

	r.Host = (address + ":" + strconv.Itoa(port))
	return (address + ":" + strconv.Itoa(port)),errs
}

func NewBuildReverseProxyDirector(c *sysadmServer.Context, data *ProxyRewriteData)(func(r *http.Request), []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	return func(r *http.Request) {
		if data.HeaderModifyFunc != nil {
			mHeader := data.HeaderModifyFunc
			mHeader(r)
		}

		if strings.TrimSpace(data.Method) != "" {
			r.Method = data.Method
		}

		if data.AuthData != nil && data.AuthData.IsAuth {
			if strings.ToLower(data.AuthData.AuthType) == "basic"{
				r.SetBasicAuth(data.AuthData.UserName,data.AuthData.Password)
			}
		}

		if data.UrlModifyFunc != nil {
			mUrlFunc := data.UrlModifyFunc
			_,err := mUrlFunc(r,data.ApiServerData)
			errs = append(errs,err...)
		}

		if data.HostModifyFunc != nil {
			mHostFunc := data.HostModifyFunc
			_,err := mHostFunc(r,data.ApiServerData)
			errs = append(errs,err...)
		}

	},errs

}

func BuildApiServerData(moduleName string,actionName string,apiVersion string,tls bool,address string,port int,ca string,cert string,key string) *ApiServerData {
	var ret *ApiServerData = nil
	if moduleName == "" || actionName == "" || apiVersion == "" || address == "" || port == 0 {
		return ret
	}
	
	tlsFiles := TlsFile{
		Ca: ca,
		Cert: cert,
		Key: key,
	}

	return &ApiServerData{
		ModuleName: moduleName,
		ActionName: actionName,
		ApiVersion: apiVersion,
		Server: Server{
			Tls: tls,
			Address: address,
			Port: port,
			TlsFile: tlsFiles,
		},
	}
}


/*
	passProxy: 
	1. change the host and port of the request to registry server
	2. change the url of the request to registry server
	3. pass the request to registry server
*/
func PassProxy(c *sysadmServer.Context, data *ApiServerData) ([]sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	r := c.Request
	method := r.Method
	uri := r.RequestURI
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010004,"debug","module %s action %s method %s uri %s",data.ModuleName,data.ActionName,method,uri))

	// build proxy director
	reverseProxyDirector,err := BuildReverseProxyDirector(c,data)
	errs=append(errs,err...)
	if reverseProxyDirector == nil {
		err = SendResponseForErrorMessage(c,800010005,"internal error")
		errs=append(errs,err...)
		return errs
	}

	// build roundTripper
	roundTripper := httpclient.BuildRoundTripper(nil)
	
	// set ReverseProxy
	registryProxy := httputil.ReverseProxy{
		Director: reverseProxyDirector,
		Transport: roundTripper,
	}

	registryProxy.ServeHTTP(c.Writer,c.Request)
	
	return errs

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
		dataSet := make(map[string]interface{},0)
		ret.Status = false
		ret.ErrorCode = 800010005
		encodeMsg := "the response body from  the server is empty"
		dataSet["msg"] = encodeMsg
		ret.Message = append(ret.Message,dataSet)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010005,"error","the response body from  the server is empty"))
		return ret, errs
	}

	res := &ApiResponseData{}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010008,"debug","Unmarshal body %s",body))
	e := json.Unmarshal(body,res)
	if e != nil {
		dataSet := make(map[string]interface{},0)
		ret.Status = false
		ret.ErrorCode = 800010006
		dataSet["msg"] = fmt.Sprintf("can not parsing reponse body to json. error: %s",e) 
		ret.Message = append(ret.Message,dataSet)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010006,"error","can not parsing reponse body(%s) to json. error: %s",body,e))
		return ret , errs
	}

	ret.Status = res.Status
	ret.ErrorCode = res.ErrorCode
	var tmpMessage []map[string]interface{}
	message := res.Message
	for _,m := range message {
		tmpMap := make(map[string]interface{},0)
		for k,v := range m {
			vStr,ok := v.(string)
			if ok {
				vDecode, err := base64.StdEncoding.DecodeString(vStr)
				if err != nil { 
					errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010007,"debug","can not decode content(%s) for key %s  error %s",vStr,k,err))
					tmpMap[k] = vStr
				}else {
					errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010009,"debug","decode content(%s) for key %s",vDecode,k))
					tmpMap[k] = vDecode
				}
			}else {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(800010010,"debug","the value of  key %s is %#v",k,v))
				tmpMap[k] = v
			}
		}
		tmpMessage = append(tmpMessage,tmpMap)
	}
	ret.Message = tmpMessage

	return ret,errs
}


func ActionNotFound(c *sysadmServer.Context,module string, action string, method string) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	errs = NewSendResponseForErrorMessage(c,make(map[string][]string,0),http.StatusNotFound,"json",800010023,"request data error: module %s has not action  %s with %s method",module,action,method)

	return errs
}