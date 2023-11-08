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
	sysadmsessions "github.com/wangyysde/sysadmSessions"
	"strings"
	sysadmDB "sysadm/db"
	sysadmObjects "sysadm/objects/app"
	"sysadm/sysadmLog"
	sysadmSetting "sysadm/syssetting/app"
)

// 设置运行期数据
func SetRunData(dbConf *sysadmDB.DbConfig, logEntity *sysadmLog.LoggerConfig, workingRoot string) error {
	workingRoot = strings.TrimSpace(workingRoot)
	if dbConf == nil || logEntity == nil || workingRoot == "" {
		return fmt.Errorf("数据库配置为空，或日志配置为空，或工作根目录为空")
	}

	runData.dbConf = dbConf
	runData.logEntity = logEntity
	runData.workingRoot = workingRoot

	if e := sysadmObjects.SetRunDataForDBConf(dbConf); e != nil {
		return e
	}

	if e := sysadmObjects.SetWorkingRoot(workingRoot); e != nil {
		return e
	}

	return nil
}

// 设置session数据
func SetSessionOptions(options sysadmsessions.Options, sessionName string) {
	runData.sessionOption = options
	runData.sessionName = sessionName
}

// 设置页面信息
func SetPageInfo(pageInfo sysadmSetting.PageInfo) {
	runData.pageInfo = pageInfo
}

// 添加可用区处理的事件侦听器
func AddHandlers(r *sysadmServer.Engine) error {
	if r == nil {
		return fmt.Errorf("can not add handlers on nil ")
	}

	// 为api设置事件处理器
	/*
		groupPath := "/api/" + DefaultApiVersion + "/" + DefaultModuleName
		v1 := r.Group(groupPath)
		{
			v1.POST("/validName", validNameHandler)
			v1.POST("/validCommand", validCommandHandler)
			v1.POST("/validCrontab", validCrontabHandler)
		}
	*/
	// 为前端显示设置事件处理器
	groupPath := "/" + defaultObjectName
	display := r.Group(groupPath)
	{
		// for cluster options on the list page
		display.GET("/getclusteroptionfbydc", getclusteroptionfbydcHandler)
		// for namespace options on the list page
		display.GET("/getnsoptionbycluster", getnsoptionbyclusterHandler)
	}

	deploymentPath := "/" + defaultObjectName + "/deployment"
	displayDeployment := r.Group(deploymentPath)
	{
		//	display.GET("/add", addK8scluster)
		displayDeployment.GET("/list", listDeploymentHandler)
	}

	statefulsetPath := "/" + defaultObjectName + "/statefulset"
	displayStatefulset := r.Group(statefulsetPath)
	{
		displayStatefulset.GET("/list", listStatefulsetHandler)
	}

	daemonsetPath := "/" + defaultObjectName + "/daemonset"
	displayDaemonset := r.Group(daemonsetPath)
	{
		displayDaemonset.GET("/list", listDaemonsetHandler)
	}

	jobPath := "/" + defaultObjectName + "/job"
	displayJob := r.Group(jobPath)
	{
		displayJob.GET("/list", listJobHandler)
	}

	cronjobPath := "/" + defaultObjectName + "/cronjob"
	displayCronjob := r.Group(cronjobPath)
	{
		displayCronjob.GET("/list", listCronjobHandler)
	}

	return nil
}
