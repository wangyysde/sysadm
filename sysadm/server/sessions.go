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
	"strings"
	"net/http"
	"time"

	sessions "github.com/wangyysde/sysadmSessions"
  	"github.com/wangyysde/sysadmSessions/cookie"
	"github.com/wangyysde/sysadmServer"

 )

 var sessionOptions = sessions.Options{
	  Path: sessionPath,
	  Domain: sessionDomain,
	  MaxAge: sessionAge,
	  Secure: false,
	  HttpOnly: true,
	  SameSite: http.SameSiteDefaultMode,
  }


 func initSession(r *sysadmServer.Engine) error{
	 if r == nil {
		return fmt.Errorf("router is nil.")
	}

	store := cookie.NewStore([]byte("secret"))
  	r.Use(sessions.Sessions(sessionName, store))
	return nil
 }

func refreshSession(c *sysadmServer.Context){
	cc,_ :=c.Request.Cookie(sessionName)
	if cc != nil{
   		cookie :=http.Cookie{
    		Name:     sessionName,
      		Value:    cc.Value,
      		Expires: time.Now().Add(time.Duration(sessionAge) * time.Second),
      		Path:    sessionPath,
      		Domain:   "",
      		Secure:   false,
      		HttpOnly: true,
   		}
		http.SetCookie(c.Writer,&cookie)
	}
}

func setSessionValue(c *sysadmServer.Context,key string, value interface{}) error{
	if c == nil {
		return fmt.Errorf("Context is nil.")
	}

	key = strings.TrimSpace(key)
	if len(key) < 1 {
		return fmt.Errorf("The length of session key must bigger 1.")
	}

	session := sessions.Default(c)
	session.Set(key, value)
    err := session.Save()
	refreshSession(c)
	return err
}

func getSessionValue(c *sysadmServer.Context,key string) (interface{},error) {
	if c == nil {
		return nil,fmt.Errorf("Context is nil.")
	}

	key = strings.TrimSpace(key)
	if len(key) < 1 {
		return nil,fmt.Errorf("The length of session key must bigger 1.")
	}

	session := sessions.Default(c)

	value := session.Get(key)
	if value == nil {
		return nil,fmt.Errorf("The value of session key:%s has not found",key)
	}

	refreshSession(c)
	return value, nil
}

func getSessionID(c *sysadmServer.Context)(string, error){
	if c == nil {
		return "",fmt.Errorf("Context is nil.")
	}

	session := sessions.Default(c)

	return session.ID(),nil
}

func clearSession(c *sysadmServer.Context) error{
	if c == nil {
		return fmt.Errorf("Context is nil.")
	}

	session := sessions.Default(c)
	session.Clear()

	return nil
}