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

import(
	"fmt"
	"net/http"

	"github.com/wangyysde/sysadmServer"
)

type formDataStruct struct{
	htmlTitle string
	formTemplateName string
	actionHandler sysadmServer.HandlerFunc
}

var formsData = map[string] formDataStruct {
	"login": formDataStruct {
		htmlTitle: "请你登录",
		formTemplateName: "login.html",
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

	if StartData.sysadmRootPath == "" {
		if _,err := getSysadmRootPath(cmdRunPath); err != nil {
			return err
		}
	}

	r.Delims(templateDelimLeft,templateDelimRight)

	formTmplPath := StartData.sysadmRootPath + "/" + formTemplateDir +"*.html" 
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
	if !ok {
		fmt.Print("We can not got name field value")
	}else{
		fmt.Printf("We got the value of name is: %s\n",value)
		c.String(http.StatusOK, fmt.Sprintf("We got the value of name is: %s\n",value))
	}

	value,ok = c.GetPostForm("password")
	if !ok {
		fmt.Print("We can not got name field value")
	}else{
		fmt.Printf("We got the value of name is: %s\n",value)
		c.String(http.StatusOK, fmt.Sprintf("We got the valvalueue of name is: %s\n",value))
	}
	
}