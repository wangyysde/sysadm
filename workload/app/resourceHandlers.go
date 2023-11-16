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

package app

import (
	"fmt"
	"github.com/wangyysde/sysadmServer"
	"strings"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
)

func resourceHandler(c *sysadmServer.Context) {
	module := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(strings.ToLower(c.Param("module"))), "/"), "/")
	action := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(strings.ToLower(c.Param("action"))), "/"), "/")
	var errs []sysadmLog.Sysadmerror

	switch action {
	case "list":
		listResourceHandler(c, module)
		return
	default:
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000800000, "debug", "action %s for %s was not found", action, module))
		e := fmt.Errorf("action %s for %s was not found", action, module)
		objectsUI.OutPutMsg(c, "", "您未登录或超时", runData.logEntity, 8000800001, errs, e)
		return
	}
}
