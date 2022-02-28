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
*/

package server

import (
	"fmt"
	"strconv"

	"github.com/wangyysde/sysadm/sysadmerror"
)

func getRepositories()([]sysadmerror.Sysadmerror){
	var requestParams requestParams = requestParams{}
	var regUrl string = "`"
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

/*
func addRegistryHandlersV2(startParams *StartParameters)(([]sysadmerror.Sysadmerror)){
	var errs []sysadmerror.Sysadmerror
	

	return errs
}

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

