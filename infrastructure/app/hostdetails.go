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
	"net/http"
	"strings"
	"sysadm/db"
	"sysadm/sysadmerror"
	"sysadm/utils"
)

func detailHost(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	errMsg := ""
	hostidMap, err := utils.GetRequestData(c, []string{"hostid"})
	errs = append(errs, err...)
	hostid, ok := hostidMap["hostid"]
	if !ok {
		errMsg = "no host has be selected to display"
	}

	var hostData map[string]string
	var e error
	if errMsg == "" {
		hostData, e = getHostDataByHostid(hostid)
		if e != nil {
			errMsg = fmt.Sprintf("%s", e)
		}
	}
	hostTplData := make(map[string]interface{}, 0)
	msgTplData := make(map[string]interface{}, 0)
	if errMsg == "" {
		status := hostData["status"]
		statusMsg, ok := hostStatus[status]
		if !ok || statusMsg == "" {
			statusMsg = hostStatus["unkown"]
		}
		hostData["statusMsg"] = statusMsg

		modeStr := "被动模式"
		agentPortStr := "----"
		if hostData["passiveMode"] == "0" {
			modeStr = "主动模式"
			agentPortStr = hostData["agentPort"]
		}
		hostData["modeStr"] = modeStr
		hostData["agentPortStr"] = agentPortStr

		offlineTime := "----"
		if hostData["status"] == "maintenance" || hostData["status"] == "offline" {
			offlineTime = hostData["offlineStartTime"]
		}
		hostData["offlineTime"] = offlineTime

		deleteTime := "----"
		if hostData["status"] == "deleted" {
			deleteTime = hostData["deletetime"]
		}
		hostData["deleteTime"] = deleteTime

		hostTplData["hostData"] = hostData
		templateFile := "hostdetails.html"
		c.HTML(http.StatusOK, templateFile, hostTplData)
		return
	}
	templateFile := "hostdetailsErrorMsg.html"
	msgTplData["errormessage"] = errMsg
	c.HTML(http.StatusOK, templateFile, msgTplData)
}

func getHostDataByHostid(hostid string) (map[string]string, error) {
	hostid = strings.TrimSpace(hostid)
	if hostid == "" {
		return nil, fmt.Errorf("no host has be selected")
	}

	// Qeurying data from DB
	whereMap := make(map[string]string, 0)
	whereMap["hostID"] = "=\"" + hostid + "\""
	whereMap["status"] = "!= \"deleted\""
	selectData := db.SelectData{
		Tb:        []string{"host"},
		OutFeilds: []string{"*"},
		Where:     whereMap,
	}
	dbEntity := WorkingData.dbConf.Entity
	dbData, _ := dbEntity.QueryData(&selectData)
	if dbData == nil {
		return nil, fmt.Errorf("no host data has be got")
	}

	if len(dbData) < 1 {
		return nil, fmt.Errorf("no host data has be got")
	}

	hostData := make(map[string]string)
	for k, v := range dbData[0] {
		vStr := utils.Interface2String(v)
		hostData[k] = vStr
	}

	return hostData, nil
}
