/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
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
* Note: this file include the function related to user session
 */

package user

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/wangyysde/sysadmServer"
	sessions "github.com/wangyysde/sysadmSessions"
	"sysadm/utils"
)

// SetSessionOptions set sessons options accroding parameters
func SetSessionOptions(path, domain string, maxAge int, secure, httpOnly bool) (*sessions.Options, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		path = DefaultSessionPath
	}

	domain = strings.TrimSpace(domain)
	if domain == "" {
		return nil, fmt.Errorf("session domain must be specified")
	}

	if maxAge <= 0 || maxAge > 86400 {
		maxAge = DefaultMaxAge
	}

	return &sessions.Options{
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: http.SameSiteDefaultMode,
	}, nil
}

func RefreshSession(c *sysadmServer.Context, sessionName string) error {
	sessionName = strings.TrimSpace(sessionName)
	if sessionName == "" {
		return fmt.Errorf("session name and session path must be specified")
	}

	cc, _ := c.Request.Cookie(sessionName)
	if cc == nil {
		return fmt.Errorf("as if cookie has not be enabled")
	}

	cookie := http.Cookie{
		Name:     sessionName,
		Value:    cc.Value,
		Expires:  time.Now().Add(time.Duration(cc.MaxAge) * time.Second),
		Path:     cc.Path,
		Domain:   cc.Domain,
		Secure:   cc.Secure,
		HttpOnly: cc.HttpOnly,
	}
	http.SetCookie(c.Writer, &cookie)

	return nil
}

func SetSessionValue(c *sysadmServer.Context, key, sessionName string, value interface{}) error {
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
	RefreshSession(c, sessionName)
	return err
}

func GetSessionValue(c *sysadmServer.Context, key, sessionName string) (interface{}, error) {
	if c == nil {
		return nil, fmt.Errorf("Context is nil.")
	}

	key = strings.TrimSpace(key)
	if len(key) < 1 {
		return nil, fmt.Errorf("The length of session key must bigger 1.")
	}

	session := sessions.Default(c)

	value := session.Get(key)
	if value == nil {
		return nil, fmt.Errorf("The value of session key:%s has not found", key)
	}

	RefreshSession(c, sessionName)
	return value, nil
}

func GetSessionID(c *sysadmServer.Context) (string, error) {
	if c == nil {
		return "", fmt.Errorf("Context is nil.")
	}

	session := sessions.Default(c)

	return session.ID(), nil
}

func ClearSession(c *sysadmServer.Context) error {
	if c == nil {
		return fmt.Errorf("Context is nil.")
	}

	session := sessions.Default(c)
	session.Clear()

	return nil
}

func IsLogin(c *sysadmServer.Context, sessionName string) (bool, int, error) {
	if c == nil {
		return false, 0, fmt.Errorf("can not check whether login on an empty connection")
	}

	sessionName = strings.TrimSpace(sessionName)
	if sessionName == "" {
		return false, 0, fmt.Errorf("session name should not be empty")
	}

	useridStr, e := GetSessionValue(c, "userid", sessionName)
	if e != nil {
		return false, 0, e
	}
	useridInt, e := utils.Interface2Int(useridStr)
	if e != nil {
		return false, 0, e
	}

	return true, useridInt, nil
}
