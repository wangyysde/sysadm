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

package objectsUI

import (
	"github.com/wangyysde/sysadmServer"
	"net/http"
	"strings"
	"sysadm/sysadmLog"
)

func OutPutErrorMsg(c *sysadmServer.Context, templateFile string, logEntity *sysadmLog.LoggerConfig, errcode int,
	errs []sysadmLog.Sysadmerror, e error) {
	templateFile = strings.TrimSpace(templateFile)
	if templateFile == "" {
		templateFile = "showmessage.html"
	}
	messageTplData := make(map[string]interface{}, 0)
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(errcode, "error", "%s", e))
	logEntity.LogErrors(errs)
	messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
	c.HTML(http.StatusOK, templateFile, messageTplData)

	return
}

func OutPutMsg(c *sysadmServer.Context, templateFile, msg string, logEntity *sysadmLog.LoggerConfig, errcode int,
	errs []sysadmLog.Sysadmerror, e error) {

	templateFile = strings.TrimSpace(templateFile)
	if templateFile == "" {
		templateFile = "infoBox.html"
	}

	errorBox := false
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(errcode, "error", "%s", e))
		logEntity.LogErrors(errs)
		errorBox = true
	}

	messageTplData := make(map[string]interface{}, 0)
	msg = strings.TrimSpace(msg)
	if msg == "" {
		msg = "系统内部出错，请稍后再试或联系系统管理员"
	}

	messageTplData["errorBox"] = errorBox
	messageTplData["message"] = msg

	c.HTML(http.StatusOK, templateFile, messageTplData)
}

func OutputResourceDetail(c *sysadmServer.Context, templateFile, errorMsg string, logEntity *sysadmLog.LoggerConfig,
	errcode int, additionalJs, additionCss []string, tplData map[string]interface{}, errs []sysadmLog.Sysadmerror, e error) {

	templateFile = strings.TrimSpace(templateFile)
	if templateFile == "" {
		templateFile = defaultResourceDetailTemplateFile
	}

	tplData["additionalJs"] = additionalJs
	tplData["additionCss"] = additionCss
	errorFlag := false
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(errcode, "error", "%s", e))
		logEntity.LogErrors(errs)
		errorFlag = true
		tplData["errorFlag"] = errorFlag
		errorMsg = strings.TrimSpace(errorMsg)
		if errorMsg == "" {
			errorMsg = "系统发生未知错误,请稍后再试或联系系统管理员"
		}
		tplData["errorMsg"] = errorMsg

	}

	tplData["errorFlag"] = errorFlag
	c.HTML(http.StatusOK, templateFile, tplData)
	return
}
