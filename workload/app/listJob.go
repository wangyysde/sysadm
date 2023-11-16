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
	"strconv"
	"strings"
	datacenter "sysadm/datacenter/app"
	sysadmK8sClient "sysadm/k8sclient"
	sysadmK8sCluster "sysadm/k8scluster/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/user"
	"time"
)

func listJobHandler(c *sysadmServer.Context) {
	// order fields data of cluster list page
	var allOrderFields = map[string]string{"TD1": "name", "TD7": ""}

	// which field will be order if user has not selected
	var defaultOrderField = "TD1"

	// 1 for DESC 0 for ASC
	var defaultOrderDirection = "1"

	// all popmenu items defined Format:
	// item name, action name, action method
	var allPopMenuItems = []string{"编辑,edit,GET,page", "删除,del,POST,tip"}

	// define all list items(cols) name
	var allListItems = map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "持续时间", "TD4": "标签", "TD5": "Pods(活/成/失/总）", "TD6": "完成时间", "TD7": "创建时间"}

	var additionalJs = []string{"js/sysadmfunctions.js", "/js/workloadList.js"}

	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000500001, "debug", "now handling Job list"))
	listTemplateFile := "workloadlist.html"

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutPutMsg(c, "", "您未登录或超时", runData.logEntity, 8000500002, errs, e)
		return
	}

	// get request data
	requestKeys := []string{"dcID", "clusterID", "namespace", "start", "orderfield", "direction", "searchContent"}
	requestData, e := getRequestData(c, requestKeys)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "参数错误,请确认是从正确地方连接过来的", runData.logEntity, 8000500003, errs, e)
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
	tplData, e := objectsUI.InitTemplateDataForWorkload("/"+defaultObjectName+"/", "工作负载", "Job列表", "", "no",
		allPopMenuItems, additionalJs, []string{}, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000500004, errs, e)
		return
	}
	tplData["objName"] = "job"

	// preparing datacenter data
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	conditions := make(map[string]string, 0)
	conditions["isDeleted"] = "=0"
	order := make(map[string]string, 0)
	var emptyString []string
	dcList, e := dcEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000500005, errs, e)
		return
	}

	e = buildSelectDataWithNs(tplData, dcList, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000500006, errs, e)
		return
	}

	startPos := objectsUI.GetStartPosFromRequest(requestData)
	count, objListData, e := prepareJobData(selectedCluster, selectedNamespace, defaultOrderField, defaultOrderDirection, startPos, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000500008, errs, e)
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

func prepareJobData(selectedCluster, selectedNS, defaultOrderField, defaultOrderDirection string,
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

	jobList, e := clientSet.BatchV1().Jobs(nsStr).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return 0, dataList, e
	}

	totalNum := len(jobList.Items)
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

	var jobItems []interface{}
	for _, item := range jobList.Items {
		jobItems = append(jobItems, item)
	}

	if direction == "1" {
		if orderfield == "TD1" {
			sort.Sort(objectsUI.SortData{Data: jobItems, By: sortJobByName})
		} else {
			sort.Sort(objectsUI.SortData{Data: jobItems, By: sortJobByCreatetime})
		}
	} else {
		if orderfield == "TD1" {
			sort.Sort(sort.Reverse(objectsUI.SortData{Data: jobItems, By: sortJobByName}))
		} else {
			sort.Sort(sort.Reverse(objectsUI.SortData{Data: jobItems, By: sortJobByCreatetime}))
		}
	}

	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}

	for i := startPos; i < endCount; i++ {
		interfaceData := jobItems[i]
		jobData, ok := interfaceData.(batchv1.Job)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not Job schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = jobData.Name
		lineMap["TD1"] = jobData.Name
		lineMap["TD2"] = jobData.Namespace
		v1StartTime := jobData.Status.StartTime
		v1CompletionTime := jobData.Status.CompletionTime
		duration := "未开始"
		if v1StartTime != nil {
			if v1CompletionTime != nil {
				startTime := v1StartTime.Time
				completeTime := v1CompletionTime.Time
				diffTime := completeTime.Sub(startTime)
				duration = diffTime.String()
			} else {
				nowTime := time.Now()
				startTime := v1StartTime.Time
				diffTime := nowTime.Sub(startTime)
				duration = diffTime.String()
			}
		}

		lineMap["TD3"] = duration
		lineMap["TD4"] = objectsUI.ConvertMap2HTML(jobData.Labels)
		activeNum := jobData.Status.Active
		successNum := jobData.Status.Succeeded
		failedNum := jobData.Status.Failed
		totalNum := activeNum + successNum + failedNum
		lineMap["TD5"] = strconv.Itoa(int(activeNum)) + "/" + strconv.Itoa(int(successNum)) + "/" + strconv.Itoa(int(failedNum)) + "/" + strconv.Itoa(int(totalNum))

		completionTimeStr := "未完成"
		if jobData.Status.CompletionTime != nil {
			completionTimeStr = jobData.Status.CompletionTime.Time.Format(objectsUI.DefaultTimeStampFormat)
		}
		lineMap["TD6"] = completionTimeStr
		lineMap["TD7"] = jobData.CreationTimestamp.Time.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := "0,1"

		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortJobByName(p, q interface{}) bool {
	pData, ok := p.(batchv1.Job)
	if !ok {
		return false
	}
	qData, ok := q.(batchv1.Job)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortJobByCreatetime(p, q interface{}) bool {
	pData, ok := p.(batchv1.Job)
	if !ok {
		return false
	}
	qData, ok := q.(batchv1.Job)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}
