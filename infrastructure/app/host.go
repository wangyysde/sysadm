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

package app

import(
	"strings"
	"fmt"

	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
)

func addHost(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror

	var requestData ApiHost
	if e := c.ShouldBind(&requestData); e != nil {
		msg := fmt.Sprintf("get host data err %s", e) 
		_ = apiutils.SendResponseForErrorMessage(c,3030001,msg)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3030001,"error",msg))
		logErrors(errs)
		return 
	}
	
	if requestData.Port == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3030002,"warn","can not get SSH Service Port. Default port will be used"))
		requestData.Port = 22
	}

	if strings.TrimSpace(requestData.User) == "" || strings.TrimSpace(requestData.Password) == "" {
		msg := "user account and password must be set" 
		_ = apiutils.SendResponseForErrorMessage(c,3030003,msg)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020005,"error",msg))
		return
	}

	msg := "host has be added successful"
	err := apiutils.SendResponseForSuccessMessage(c,msg)	
	errs=append(errs,err...)
	logErrors(errs)
}

