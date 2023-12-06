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
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/user"
)

func detailHandler(c *sysadmServer.Context, module, action string) {
	var errs []sysadmLog.Sysadmerror
	var additionalJs = []string{}
	var additionalCss = []string{}
	var detailTemplateFile = "resourceDetail.html"
	tplData := make(map[string]interface{}, 0)

	errs = append(errs, sysadmLog.NewErrorWithStringLevel(8001300001, "debug", "try to show datails for module %s with action %s", module, action))
	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutputResourceDetail(c, detailTemplateFile, "未登录或超时", runData.logEntity, 8001300002, additionalJs, additionalCss, tplData, errs, fmt.Errorf("user have not login"))
		return
	}

	// get request data
	requestKeys := []string{"dcID", "clusterID", "namespace", "objID"}
	requestData, e := getRequestData(c, requestKeys)
	if e != nil {
		objectsUI.OutputResourceDetail(c, detailTemplateFile, "参数错误,请确认是从正确地方连接过来的", runData.logEntity, 8001300003, additionalJs, additionalCss, tplData, errs, e)
		return
	}

	objEntity, e := newObjEntity(module)
	if e != nil {
		objectsUI.OutputResourceDetail(c, detailTemplateFile, "", runData.logEntity, 8001300004, additionalJs, additionalCss, tplData, errs, e)
		return
	}
	objEntity.setObjectInfo()

	e = objEntity.showResourceDetail(action, tplData, requestData)
	objectsUI.OutputResourceDetail(c, detailTemplateFile, "", runData.logEntity, 8001300005, additionalJs, additionalCss, tplData, errs, e)
	return

}
