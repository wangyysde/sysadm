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
*/

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"

	sysadm "github.com/wangyysde/sysadm/sysadm/server"
)

type BodyError struct {
	Code string `json:"code"` 
	Message string `json:"message"`
	Detail string `json:"detail"`
}

type apiResponseForBool struct {
	status bool `json:"status"`
	errorcode string `json:"error"`
	message string `json:"message"`
}

type ReponseError struct {
	Errors []BodyError `json:"errors"`
}

func getRepositories()([]sysadmerror.Sysadmerror){
	var requestParams requestParams = requestParams{}
	var regUrl string = ""
	if definedConfig.Registry.Server.Tls {
		if definedConfig.Registry.Server.Port == 443 {
			regUrl = "https://" + definedConfig.Registry.Server.Host + "/v2/_catalog"
		} else {
			regUrl = "https://" + definedConfig.Registry.Server.Host + ":" + strconv.Itoa(definedConfig.Registry.Server.Port) + "/v2/_catalog"
		}
	}else {
		if definedConfig.Registry.Server.Port == 80 {
			regUrl = "http://" + definedConfig.Registry.Server.Host + "/v2/_catalog"
		} else {
			regUrl = "http://" + definedConfig.Registry.Server.Host + ":" + strconv.Itoa(definedConfig.Registry.Server.Port) + "/v2/_catalog"
		}
	}

	requestParams.url = regUrl
	requestParams.method = "GET"
	body,err := sendRequest(&requestParams)
	logErrors(err)
	fmt.Println(string(body))
	
	return err
}


func addRegistryV2RootHandler()(([]sysadmerror.Sysadmerror)){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030001,"debug","now adding /v2 handler"))
	r := StartData.router
	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030002,"fatal","we can not add handler to a nil router."))
		return errs
	}
	
	r.GET("/v2/", handerRootV2)
	r.POST("/v2/", handerRootV2)
	r.HEAD("/v2/", handerRootV2)
	r.PUT("/v2/", handerRootV2)
	r.PATCH("/v2/",handerRootV2)
	
	return errs
}

func handerRootV2(c *sysadmServer.Context) {
	r := c.Request
	username,password,_ := r.BasicAuth()
	ok := isLogin(username,password)
	if !ok {

		//Bearer realm="http://harbor.bzhy.com/service/token",service="harbor-registry"
		c.Header("Docker-Distribution-API-Version","registry/2.0")
		//c.Header("Content-Length","83")
		c.Header("WWW-Authenticate","Basic realm=\"basic-realm\"")
		//	c.Header("WWW-Authenticate","Bearer realm=\"http://harbor.bzhy.com/service/token\",service=\"harbor-registry\"")
		be := []BodyError{{
			Code: "UNAUTHORIZED",
			Message: "unauthorized:unauthorized",
			Detail: "",},
		}
		
		
		var re = ReponseError{
			Errors: be,
		}
	
		reponseBody,err := json.Marshal(re)
		sysadmServer.Logf("debug","json is: %s\n",reponseBody)
		if err == nil {
			c.JSON(http.StatusUnauthorized,re)
		}else {
			c.JSON(http.StatusUnauthorized,sysadmServer.H{})
		}
		return 
	}

	
	c.Header("Docker-Distribution-API-Version","registry/2.0")
	c.JSON(http.StatusOK,  sysadmServer.H{})
}

func isLogin(username string, password string) bool {
	var errs []sysadmerror.Sysadmerror
	if username == "" && password == "" {
		return false
	}

	var reqUrl string = ""
	m := sysadm.Modules
	if definedConfig.Sysadm.Server.Tls {
		if definedConfig.Sysadm.Server.Port == 443 {
			reqUrl = "https://" + definedConfig.Sysadm.Server.Host + "/api/" + definedConfig.Sysadm.ApiVerion + m["user"].Path +"/login" 
		} else {
			reqUrl = "https://" + definedConfig.Sysadm.Server.Host + ":" + strconv.Itoa(definedConfig.Sysadm.Server.Port) + "/api/" + definedConfig.Sysadm.ApiVerion+ m["user"].Path +"/login" 
		}
	}else {
		if definedConfig.Sysadm.Server.Port == 80 {
			reqUrl = "http://" + definedConfig.Sysadm.Server.Host + "/api/" + definedConfig.Sysadm.ApiVerion + m["user"].Path +"/login"
		} else {
			reqUrl = "http://" + definedConfig.Sysadm.Server.Host + ":" + strconv.Itoa(definedConfig.Sysadm.Server.Port) + "/api/" + definedConfig.Sysadm.ApiVerion + m["user"].Path +"/login"
		}
	}

	var requestParams requestParams = requestParams{}
	requestParams.url = reqUrl
	requestParams.method = "POST"
	requestParams.data = append(requestParams.data,&requestData{key: "username", value: username})
	requestParams.data = append(requestParams.data,&requestData{key: "password", value: password})
	
	body,err := sendRequest(&requestParams)
	logErrors(err)
	var ret *apiResponseForBool = &apiResponseForBool{} 
	if len(body) > 1{
		e := json.Unmarshal(body,ret)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030003,"error","can not unmarshal reponse body. error %s",e))
		logErrors(errs)
		return false
	}
	
	if ret.errorcode != "" || ret.message != "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030004,"debug","can not login with errorcode: %s message: %s",ret.errorcode,ret.message))
	}
	
	return ret.status
}
/*
func addRegistryHandlersRoot(startParams *StartParameters)(([]sysadmerror.Sysadmerror)){
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(203001,"debug","try to add handler for registry v2 root path"))
	r := startParams.router
	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(203001,"fatal","Occurred a fatal error, router for server is nil."))
		return errs
	}

	r.GET("/v2",func(c *sysadmServer.Context) {
			c.HTML(http.StatusOK, f.formTemplateName, tplData)
    	})

	return errs
}

// handler for handling login form.
func registryHandlerRoot(c *sysadmServer.Context) {
	value,ok := c.GetPostForm("username")
	if !ok || strings.TrimSpace(value) == "" {
		c.JSON(http.StatusOK, sysadmServer.H{"errCode": 100, "msg": "请输入帐号！"})
		return
	}
	username := strings.TrimSpace(value);

	value,ok = c.GetPostForm("password")
	if !ok || strings.TrimSpace(value) == "" {
		c.JSON(http.StatusOK, sysadmServer.H{"errCode": 101, "msg": "请输入密码！"})
		return 
	}
	password := strings.TrimSpace(value)

	if strings.ToLower(username) == strings.ToLower(definedConfig.User.DefaultUser) && 
		md5Encrypt(strings.ToLower(password)) == md5Encrypt(strings.ToLower(definedConfig.User.DefaultPassword)){
			if err := setSessionValue(c,"isLogin",true); err != nil {
				c.JSON(http.StatusOK, sysadmServer.H{"errCode": 102, "msg": err})
				return 	
			} else {
				c.JSON(http.StatusOK, sysadmServer.H{"errCode": 0, "msg": "登录成功！"})
				return 
			}
	} 

	c.JSON(http.StatusOK, sysadmServer.H{"errCode": 103, "msg": "用户名或密码错误！"})
	return 
}

*/

