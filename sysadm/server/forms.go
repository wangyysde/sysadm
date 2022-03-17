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
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/wangyysde/sysadm/httpclient"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

// ErrorCode: 105xxxx

type formDataStruct struct{
	htmlTitle string
	formTemplateName string
	formUri string
	actionHandler sysadmServer.HandlerFunc
}

var formsData = map[string] formDataStruct {
	"login":  {
		htmlTitle: "请你登录",
		formTemplateName: "login.html",
		formUri: "login",
		actionHandler: loginHandler,
	},
	/*
	"logout": formDataStruct {
		htmlTitle: "欢迎您再来",
		formUri: "/logout",
		formTemplateName: "test_logout.html",
		actionHandler: handlerLogout,
	},
	*/
}

// addFormHandler set delims for template and load template files
// return nil if not error otherwise return error.
func addFormHandler(r *sysadmServer.Engine,cmdRunPath string) error {
	if r == nil {
		return fmt.Errorf("router is nil.")
	}

	if RuntimeData.StartParas.SysadmRootPath  == "" {
		if _,err := getSysadmRootPath(cmdRunPath); err != nil {
			return err
		}
	}

	r.Delims(templateDelimLeft,templateDelimRight)

	formTmplPath := RuntimeData.StartParas.SysadmRootPath + "/" + formTemplateDir +"*.html" 
	r.LoadHTMLGlob(formTmplPath)

	addForms(r)

	return nil
}

// add handler to router accroding to formsData
func addForms(r *sysadmServer.Engine){
	for k,f := range formsData {
		formUri := formBaseUri + k
		tplData := map[string] interface{}{
			"htmlTitle": f.htmlTitle,
			"formUri": formUri,
			"formId": "login",
		}

		r.GET(formUri,func(c *sysadmServer.Context) {
			c.HTML(http.StatusOK, f.formTemplateName, tplData)
    	})

		actionHandler := f.actionHandler
		r.POST(formUri,actionHandler)
	}
	
}


// handler for handling login form.
func loginHandler(c *sysadmServer.Context) {
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

	if loginWithDB(username,password) {
		_ = setSessionValue(c,"isLogin",true)
		c.JSON(http.StatusOK, sysadmServer.H{"errCode": 0, "msg": "登录成功！"})
		return
	}

	c.JSON(http.StatusOK, sysadmServer.H{"errCode": 102, "msg": "用户名或密码错误！"})
	return 
}


/*
  encrypt data with md5 
  if the length is zero ,then return ""
  otherwise return encrypted data`

*/
func md5Encrypt(data string, salt string) string{
	if len(data) < 1 {
		return ""
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(data))
	if len(salt) > 0 {
		md5Ctx.Write([]byte(salt))
	}
	cipherStr := md5Ctx.Sum(nil)
	encryptedData := hex.EncodeToString(cipherStr)
	return encryptedData
}

// TODO: 
// we are plan to cut user as an independent module, so user login in sysadm should call API 
func loginWithDB(username string, password string) bool {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1050001,"debug","now checking the user is login"))
	if username == "" && password == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1050002,"error","username and password are empty."))
		logErrors(errs)
		return false
	}

	var reqUrl string = ""
	m := Modules
	
	if RuntimeData.RuningParas.DefinedConfig.Server.Tls {
		if  RuntimeData.RuningParas.DefinedConfig.Server.Port  == 443 {
			reqUrl = "https://" + RuntimeData.RuningParas.DefinedConfig.Server.Address + "/api/" + apiVersion + "/" + m["user"].Path +"/login" 
		} else {
			reqUrl = "https://" + RuntimeData.RuningParas.DefinedConfig.Server.Address + ":" + strconv.Itoa(RuntimeData.RuningParas.DefinedConfig.Server.Port) + "/api/" + apiVersion + "/" + m["user"].Path +"/login" 
		}
	}else {
		if RuntimeData.RuningParas.DefinedConfig.Server.Port == 80 {
			reqUrl = "http://" + RuntimeData.RuningParas.DefinedConfig.Server.Address + "/api/" + apiVersion  + "/" + m["user"].Path +"/login"
		} else {
			reqUrl = "http://" + RuntimeData.RuningParas.DefinedConfig.Server.Address + ":" + strconv.Itoa(RuntimeData.RuningParas.DefinedConfig.Server.Port) + "/api/" + apiVersion +  "/" + m["user"].Path +"/login"
		}
	}

	var requestParams httpclient.RequestParams = httpclient.RequestParams{}
	requestParams.Url = reqUrl
	requestParams.Method = "POST"
	requestParams.QueryData = append(requestParams.QueryData,&httpclient.RequestData{Key: "username", Value: username})
	requestParams.QueryData = append(requestParams.QueryData,&httpclient.RequestData{Key: "password", Value: password})
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1050002,"debug","try to execute the request with:%s",reqUrl))
	body,err := httpclient.SendRequest(&requestParams)
	errs = append(errs, err...)

	if len(body) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1050003,"error","the response from  the server is empty"))
		logErrors(errs)
		return false
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1050004,"debug","got response body is: %s",string(body)))
	ret := &ApiResponseStatus{}
	e := json.Unmarshal(body,ret)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1050005,"error","can not parsing reponse body to json. error: %s",e))
		logErrors(errs)
		return false
	}

	if ret.Errorcode  != 0 || ret.Message  != "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1050006,"debug","can not login with errorcode: %d message: %s",ret.Errorcode,ret.Message))
		logErrors(errs)
		return false
	}
	
	logErrors(errs)
	return ret.Status
}