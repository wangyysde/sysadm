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
	"net/http"
	"strings"

	"github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

type apiUserHandler func (u User)(c *sysadmServer.Context)

var  userActions = []string{"login"}


 func (u User) Name()string{
	return "user"
}

func (u User) ActionHanderCaller(action string, c *sysadmServer.Context){
	switch action{
		case "login":
			u.loginHandler(c)
	}
	
	return
}

/* 
	handling user login according to username and password provided by rquest's URL
	response the client with Status: false, Erorrcode: int, and Message: string if login is failed
	otherwise response the client with Status: true, Erorrcode: 0, and Message: "" if login is successful
*/
func (u User) loginHandler(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1040001,"debug","now handling login action handler through api."))
	   
	username,okUsername := c.GetQuery("username")
	password,okPassword := c.GetQuery("password")
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	if username == "" || password == "" || !okUsername || !okPassword {
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1040004,
			Message: "username or password incorrect",
		}
		c.JSON(http.StatusOK, ret)
		return 
	}

	// Qeurying data from DB
	selectData := db.SelectData{
		Tb: []string{"user"},
		OutFeilds: []string{"userid","username","password","deleted","salt",},
		Where: map[string]string{"username": "= '" + username + "'"},
	}
	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
		errs = append(errs,err...)
		logErrors(errs)
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1040005,
			Message: "database query error",
		}
		c.JSON(http.StatusOK, ret)
		return 
	} 

	// if the user is not exist in DB 
	if len(retData) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1040006,"debug","user %s try to login by api ,but username %s is not exist.",username,username))
		logErrors(err)
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1040006,
			Message: "username or password incorrect",
		}
		c.JSON(http.StatusOK, ret)
		return 
	}
	
	// checking password 
	row := retData[0]
	dbPassword := row["password"].(string)
	salt := row["salt"].(string)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1040007,"debug","dbpassword %s.",dbPassword))
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1040008,"debug","password %s.",password))
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1040009,"debug","salt %s.",salt))
	if(md5Encrypt(password,salt) == strings.TrimSpace(dbPassword)) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1040010,"debug","user %s login successful.",username))
		logErrors(err)
		ret := ApiResponseStatus {
			Status: true,
			Errorcode: 0,
			Message: "",
		}
		c.JSON(http.StatusOK, ret)
		return 
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1040011,"debug","user %s try to login by api, but password error.",username))
	logErrors(err)
	ret := ApiResponseStatus {
		Status: false,
		Errorcode: 1040008,
		Message: "username or password incorrect",
	}
	c.JSON(http.StatusOK, ret)
	
	return 

}