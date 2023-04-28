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
 */

package server

import (
	"strings"

	"sysadm/sysadmapi/apiutils"
	"sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

var yumActions = []string{"yumlist", "add", "del","getosversion","getobject","getcount"}

func (y Yum) ModuleName() string {
	return "yum"
}

func (y Yum) ActionHanderCaller(action string, c *sysadmServer.Context) {
	switch action {
	case "yumlist": // TODO
		y.yumActionHandler(c, "yumlist")
	case "add": //TODO
		y.yumActionHandler(c, "add")
	case "del": //TODO
		y.yumActionHandler(c, "del")
	case "getosversion":
		y.yumActionHandler(c, "getosversion")
	case "getobject":
		y.yumActionHandler(c, "getobject")
	case "getcount":
		y.yumActionHandler(c, "getcount")
	}

}

/*
export
*/
func (y Yum) yumActionHandler(c *sysadmServer.Context, action string) {
	var errs []sysadmerror.Sysadmerror

	moduleName := "yum"
	found := false
	for _, a := range yumActions {
		if strings.EqualFold(action, a) {
			found = true
			break
		}
	}

	if !found {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(7000101101, "error", "api server parameters error"))
		err := apiutils.SendResponseForErrorMessage(c, 7000101101, "api server parameters error")
		errs = append(errs, err...)
		logErrors(errs)
	}

	definedConfig := RuntimeData.RuningParas.DefinedConfig
	apiVersion := definedConfig.Registryctl.ApiVersion
	tls := definedConfig.Registryctl.Tls
	address := definedConfig.Registryctl.Address
	port := definedConfig.Registryctl.Port
	ca := definedConfig.Registryctl.Ca
	cert := definedConfig.Registryctl.Cert
	key := definedConfig.Registryctl.Key

	apiServerData := apiutils.BuildApiServerData(moduleName, strings.TrimSpace(strings.ToLower(action)), apiVersion, tls, address, port, ca, cert, key)
	if apiServerData == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(7000101002, "error", "api server parameters error"))
		err := apiutils.SendResponseForErrorMessage(c, 7000101002, "api server parameters error")
		errs = append(errs, err...)
		logErrors(errs)
	}

	err := apiutils.PassProxy(c, apiServerData)
	errs = append(errs, err...)
	logErrors(errs)
}
