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
	"net/http" //

	"github.com/wangyysde/sysadm/httpclient"
	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
	"github.com/wangyysde/sysadmServer"
)

var registryctlActions = []string{"imagelist", "getinfo", "imagedel", "taglist", "tagdel", "yumlist", "yumadd", "yumdel"}

func (r Registryctl) ModuleName() string {
	return "registryctl"
}

func (r Registryctl) ActionHanderCaller(action string, c *sysadmServer.Context) {
	switch action {
	case "imagelist":
		r.imagelistHandler(c)
	case "imagedel":
		r.imagedelHandler(c)
	case "tagdel":
		r.tagdelHandler(c)
	case "getinfo":
		r.getInfoHandler(c)
	}

}

/*
handling user login according to username and password provided by rquest's URL
response the client with Status: false, Erorrcode: int, and Message: string if login is failed
otherwise response the client with Status: true, Erorrcode: 0, and Message: "" if login is successful
*/
func (r Registryctl) imagelistHandler(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror

	moduleName := "registryctl"
	actionName := "imagelist"

	definedConfig := RuntimeData.RuningParas.DefinedConfig
	apiVersion := definedConfig.Registryctl.ApiVersion
	tls := definedConfig.Registryctl.Tls
	address := definedConfig.Registryctl.Address
	port := definedConfig.Registryctl.Port
	ca := definedConfig.Registryctl.Ca
	cert := definedConfig.Registryctl.Cert
	key := definedConfig.Registryctl.Key

	apiServerData := apiutils.BuildApiServerData(moduleName, actionName, apiVersion, tls, address, port, ca, cert, key)
	if apiServerData == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700010006, "error", "api server parameters error"))
		err := apiutils.SendResponseForErrorMessage(c, 700010006, "api server parameters error")
		errs = append(errs, err...)
		logErrors(errs)
	}

	err := apiutils.PassProxy(c, apiServerData)
	errs = append(errs, err...)
	logErrors(errs)

	// 1. user authorization and privelege checks
}

/*
handling user login according to username and password provided by rquest's URL
response the client with Status: false, Erorrcode: int, and Message: string if login is failed
otherwise response the client with Status: true, Erorrcode: 0, and Message: "" if login is successful
*/
func (r Registryctl) imagedelHandler(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror

	moduleName := "registryctl"
	actionName := "imagedel"

	definedConfig := RuntimeData.RuningParas.DefinedConfig
	apiVersion := definedConfig.Registryctl.ApiVersion
	tls := definedConfig.Registryctl.Tls
	address := definedConfig.Registryctl.Address
	port := definedConfig.Registryctl.Port
	ca := definedConfig.Registryctl.Ca
	cert := definedConfig.Registryctl.Cert
	key := definedConfig.Registryctl.Key

	apiServerData := apiutils.BuildApiServerData(moduleName, actionName, apiVersion, tls, address, port, ca, cert, key)
	if apiServerData == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700010007, "error", "api server parameters error"))
		err := apiutils.SendResponseForErrorMessage(c, 700010007, "api server parameters error")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	keys := []string{"imageid[]"}
	datas, err := utils.GetRequestDataArray(c, keys)
	errs = append(errs, err...)
	data, okdata := datas["imageid[]"]
	if !okdata || len(data) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700010008, "error", "no image has be selected"))
		err := apiutils.SendResponseForErrorMessage(c, 700010008, "no image has be selected")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	imageidStr := ""
	for _, id := range data {
		if imageidStr == "" {
			imageidStr = id
		} else {
			imageidStr = imageidStr + "," + id
		}
	}

	rawURL, err := apiutils.BuildApiUrl(apiServerData)
	errs = append(errs, err...)

	var queryData []*httpclient.RequestData
	queryData = append(queryData, &httpclient.RequestData{Key: "imageid", Value: imageidStr})
	requestParams := httpclient.RequestParams{
		Headers:       []httpclient.RequestData{},
		QueryData:     queryData,
		BasicAuthData: map[string]string{},
		Method:        http.MethodDelete,
		Url:           rawURL,
	}

	body, err := httpclient.SendRequest(&requestParams)
	errs = append(errs, err...)
	if len(body) < 1 {
		err := apiutils.SendResponseForErrorMessage(c, 700010009, "unkown error has occurred.")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	parsedRet, err := apiutils.ParseResponseBody(body)
	errs = append(errs, err...)
	c.JSON(http.StatusOK, parsedRet)

	logErrors(errs)
}

/*
 */
func (r Registryctl) tagdelHandler(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror

	moduleName := "registryctl"
	actionName := "tagdel"

	definedConfig := RuntimeData.RuningParas.DefinedConfig
	apiVersion := definedConfig.Registryctl.ApiVersion
	tls := definedConfig.Registryctl.Tls
	address := definedConfig.Registryctl.Address
	port := definedConfig.Registryctl.Port
	ca := definedConfig.Registryctl.Ca
	cert := definedConfig.Registryctl.Cert
	key := definedConfig.Registryctl.Key

	apiServerData := apiutils.BuildApiServerData(moduleName, actionName, apiVersion, tls, address, port, ca, cert, key)

	keys := []string{"tagid[]"}
	datas, err := utils.GetRequestDataArray(c, keys)
	errs = append(errs, err...)
	data, okdata := datas["tagid[]"]
	if !okdata || len(data) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700010010, "error", "no tag has be selected"))
		err := apiutils.SendResponseForErrorMessage(c, 700010010, "no tag has be selected")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	tagidStr := ""
	for _, id := range data {
		if tagidStr == "" {
			tagidStr = id
		} else {
			tagidStr = tagidStr + "," + id
		}
	}

	rawURL, err := apiutils.BuildApiUrl(apiServerData)
	errs = append(errs, err...)

	var queryData []*httpclient.RequestData
	queryData = append(queryData, &httpclient.RequestData{Key: "tagid", Value: tagidStr})
	requestParams := httpclient.RequestParams{
		Headers:       []httpclient.RequestData{},
		QueryData:     queryData,
		BasicAuthData: map[string]string{},
		Method:        http.MethodDelete,
		Url:           rawURL,
	}

	body, err := httpclient.SendRequest(&requestParams)
	errs = append(errs, err...)
	if len(body) < 1 {
		err := apiutils.SendResponseForErrorMessage(c, 700010011, "unkown error has occurred.")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	parsedRet, err := apiutils.ParseResponseBody(body)
	errs = append(errs, err...)
	c.JSON(http.StatusOK, parsedRet)

	logErrors(errs)
}

func (r Registryctl) getInfoHandler(c *sysadmServer.Context) {

}
