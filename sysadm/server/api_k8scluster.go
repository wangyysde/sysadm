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
	"github.com/wangyysde/sysadmServer"
	k8scluster "sysadm/k8scluster/app"
	"sysadm/sysadmerror"
	sysadmSetting "sysadm/syssetting/app"
)

func addK8sclusterHandlers(r *sysadmServer.Engine) []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror

	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030001, "fatal", "can not add handlers to nil"))
		return errs
	}

	dbConf := RuntimeData.RuningParas.DBConfig

	workingRoot := RuntimeData.StartParas.SysadmRootPath

	if e := k8scluster.SetRunData(dbConf, RuntimeData.sysadmLogEntity, workingRoot); e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030002, "fatal", "set run data error error message is：%s", e))
	}

	k8scluster.SetSessionOptions(sessionOptions, sessionName)
	pageInfo := sysadmSetting.PageInfo{
		NumPerPage: numPerPage,
	}

	k8scluster.SetPageInfo(pageInfo)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("fatal") {
		return errs
	}

	if e := k8scluster.AddHandlers(r); e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030003, "fatal", "add handlers error, error message is %s", e))
	}

	return errs
}
