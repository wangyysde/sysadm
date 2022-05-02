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

package apiutils

import (
	"net/http"

	"github.com/wangyysde/sysadm/sysadmerror"
)

type TlsFile struct {
	Ca string 
	Cert string
	Key string
}
type Server struct {
	Tls bool
	Address string
	Port int
	TlsFile TlsFile
}

type ApiServerData struct {
	ModuleName string
	ActionName string 
	ApiVersion string
	Server Server
}

type ApiResponseData struct {
	// Status is false if this is a error response, otherwise Status is true
	Status bool `json:"status"`
	// Errorcode is zero if this is a successful response, otherwise Errorcode is nonzero
	ErrorCode int `json:"errorCode"`
	// Message is the result sets if this is a successful ,otherwise Message is []map[string]interface
	// which has one rows only:["msg"] = message (encoded by base64)
	Message []map[string]interface{} `json:"message"`
}

type ApiInterface interface {
	GetModuleName() string
	GetActionList() []string
 }


 type Module struct {
	Name string
	Path string
	Entity ApiInterface
}

type HttpAuthData struct {
	IsAuth bool
	AuthType string
	UserName string
	Password string
}

type ProxyRewriteData struct {
	HeaderModifyFunc  func (r *http.Request) 
	Method string
	AuthData *HttpAuthData
	UrlModifyFunc func (r *http.Request,data *ApiServerData) (string, []sysadmerror.Sysadmerror) 
	HostModifyFunc func (r *http.Request,data *ApiServerData)(string, []sysadmerror.Sysadmerror)
	ApiServerData *ApiServerData
}