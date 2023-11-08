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
	"context"
	"fmt"
	"github.com/wangyysde/sysadmServer"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"sort"
	"strings"
	datacenter "sysadm/datacenter/app"
	sysadmK8sClient "sysadm/k8sclient"
	sysadmK8sCluster "sysadm/k8scluster/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/user"
)

func listCronjobHandler(c *sysadmServer.Context) {
	// order fields data of cluster list page
	var allOrderFields = map[string]string{"TD1": "name", "TD6": ""}

	// which field will be order if user has not selected
	var defaultOrderField = "TD1"

	// 1 for DESC 0 for ASC
	var defaultOrderDirection = "1"

	// all popmenu items defined Format:
	// item name, action name, action method
	var allPopMenuItems = []string{"编辑,edit,GET,page", "删除,del,POST,tip"}

	// define all list items(cols) name
	var allListItems = map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "标签", "TD4": "上次调度时间", "TD5": "上次成功执行时间", "TD6": "创建时间"}

	var additionalJs = []string{"js/sysadmfunctions.js", "/js/workloadList.js"}

	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000400001, "debug", "now handling CronJob list"))
	listTemplateFile := "workloadlist.html"

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutPutMsg(c, "", "您未登录或超时", runData.logEntity, 8000400002, errs, e)
		return
	}

	// get request data
	requestKeys := []string{"dcID", "clusterID", "namespace", "start", "orderfield", "direction", "searchContent"}
	requestData, e := getRequestData(c, requestKeys)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "参数错误,请确认是从正确地方连接过来的", runData.logEntity, 8000400003, errs, e)
		return
	}

	// preparing object list data
	selectedCluster := strings.TrimSpace(requestData["clusterID"])
	if selectedCluster == "" {
		selectedCluster = "0"
	}

	selectedNamespace := strings.TrimSpace(requestData["namespace"])
	if selectedNamespace == "" {
		selectedNamespace = "0"
	}

	// 初始化模板数据
	tplData, e := objectsUI.InitTemplateDataForWorkload("/"+defaultObjectName+"/", "工作负载", "CronJobt列表", "", "no",
		allPopMenuItems, additionalJs, []string{}, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000400004, errs, e)
		return
	}
	tplData["objName"] = "cronjob"

	// preparing datacenter data
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	conditions := make(map[string]string, 0)
	conditions["isDeleted"] = "=0"
	order := make(map[string]string, 0)
	var emptyString []string
	dcList, e := dcEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000400005, errs, e)
		return
	}

	e = buildSelectDataWithNs(tplData, dcList, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000400006, errs, e)
		return
	}

	startPos := objectsUI.GetStartPosFromRequest(requestData)
	count, objListData, e := prepareCronJobData(selectedCluster, selectedNamespace, defaultOrderField, defaultOrderDirection, startPos, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000400008, errs, e)
		return
	}

	if count == 0 {
		tplData["noData"] = "当前系统中没有数据"
	} else {
		tplData["noData"] = ""
		tplData["objListData"] = objListData

		// build table header for list objects
		objectsUI.BuildThDataForWorkloadList(requestData, allOrderFields, allListItems, tplData, defaultOrderField, defaultOrderDirection)

		// prepare page number information
		objectsUI.BuildPageNumInfoForWorkloadList(tplData, requestData, count, startPos, runData.pageInfo.NumPerPage, defaultOrderField, defaultOrderDirection)
	}
	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, listTemplateFile, tplData)

}

func prepareCronJobData(selectedCluster, selectedNS, defaultOrderField, defaultOrderDirection string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	if selectedCluster == "" || selectedCluster == "0" {
		return 0, dataList, nil
	}

	nsStr := ""
	if selectedNS != "0" {
		nsStr = selectedNS
	}

	var k8sclusterEntity sysadmObjects.ObjectEntity
	k8sclusterEntity = sysadmK8sCluster.New()
	clusterInfo, e := k8sclusterEntity.GetObjectInfoByID(selectedCluster)
	if e != nil {
		return 0, dataList, e
	}
	clusterData, ok := clusterInfo.(sysadmK8sCluster.K8sclusterSchema)
	if !ok {
		return 0, dataList, fmt.Errorf("the data is not K8S data schema")
	}
	ca := []byte(clusterData.Ca)
	cert := []byte(clusterData.Cert)
	key := []byte(clusterData.Key)
	restConf, e := sysadmK8sClient.BuildConfigFromParametes(ca, cert, key, clusterData.Apiserver, clusterData.Id, clusterData.ClusterUser, "default")
	if e != nil {
		return 0, dataList, e
	}

	clientSet, e := sysadmK8sClient.BuildClientset(restConf)
	if e != nil {
		return 0, dataList, e
	}

	cronJobList, e := clientSet.BatchV1().CronJobs(nsStr).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return 0, dataList, e
	}

	totalNum := len(cronJobList.Items)
	if totalNum < 1 {
		return 0, dataList, nil
	}

	orderfield := requestData["orderfield"]
	direction := requestData["direction"]
	if orderfield == "" {
		orderfield = defaultOrderField
	}
	if direction == "" || (direction != "0" && direction != "1") {
		direction = defaultOrderDirection
	}

	var cronJobItems []interface{}
	for _, item := range cronJobList.Items {
		cronJobItems = append(cronJobItems, item)
	}

	if direction == "1" {
		if orderfield == "TD1" {
			sort.Sort(sysadmObjects.SortData{Data: cronJobItems, By: sortCronJobByName})
		} else {
			sort.Sort(sysadmObjects.SortData{Data: cronJobItems, By: sortCronJobByCreatetime})
		}
	} else {
		if orderfield == "TD1" {
			sort.Sort(sort.Reverse(sysadmObjects.SortData{Data: cronJobItems, By: sortCronJobByName}))
		} else {
			sort.Sort(sort.Reverse(sysadmObjects.SortData{Data: cronJobItems, By: sortCronJobByCreatetime}))
		}
	}

	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}

	for i := startPos; i < endCount; i++ {
		interfaceData := cronJobItems[i]
		cronJobData, ok := interfaceData.(batchv1.CronJob)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not CronJob schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = cronJobData.Name
		lineMap["TD1"] = cronJobData.Name
		lineMap["TD2"] = cronJobData.Namespace
		lineMap["TD3"] = objectsUI.ConvertMap2HTML(cronJobData.Labels)
		lineMap["TD4"] = cronJobData.Status.LastScheduleTime.Format(objectsUI.DefaultTimeStampFormat)
		lineMap["TD5"] = cronJobData.Status.LastSuccessfulTime.Format(objectsUI.DefaultTimeStampFormat)
		lineMap["TD6"] = cronJobData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := "0,1"
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortCronJobByName(p, q interface{}) bool {
	pData, ok := p.(batchv1.CronJob)
	if !ok {
		return false
	}
	qData, ok := q.(batchv1.CronJob)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortCronJobByCreatetime(p, q interface{}) bool {
	pData, ok := p.(batchv1.CronJob)
	if !ok {
		return false
	}
	qData, ok := q.(batchv1.CronJob)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}
