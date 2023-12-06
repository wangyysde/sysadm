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
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

func listNamespaceHandler(c *sysadmServer.Context) {
	// title of add sign
	var addButtonTitle = "创建新的命名空间"

	// order fields data of cluster list page
	var allOrderFields = map[string]string{"TD1": "name", "TD2": "", "TD4": ""}

	// which field will be order if user has not selected
	var defaultOrderField = "TD1"

	// 1 for DESC 0 for ASC
	var defaultOrderDirection = "1"

	// all popmenu items defined Format:
	// item name, action name, action method
	var allPopMenuItems = []string{"详情,detail,GET,poppage", "编辑,edit,GET,page", "删除,del,POST,tip", "新增配额,addQuota,GET,page", "配额列表,QuotaList,Get,page", "新增默认配额,addLimitRange,GET,page", "默认资源配额详情,limitRangeDetail,GET,poppage", "编辑默认配额,limitRangeEdit,GET,page", "删除默认配额,limitRangeDel,post,tip"}

	// define all list items(cols) name
	var allListItems = map[string]string{"TD1": "名称", "TD2": "状态", "TD3": "标签", "TD4": "创建时间"}

	var additionalJs = []string{"js/sysadmfunctions.js", "/js/workloadList.js"}

	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000600001, "debug", "now handling namespace list"))
	listTemplateFile := "workloadlist.html"

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutPutMsg(c, "", "您未登录或超时", runData.logEntity, 8000600002, errs, e)
		return
	}

	// get request data
	requestKeys := []string{"dcID", "clusterID", "start", "orderfield", "direction", "searchContent"}
	requestData, e := getRequestData(c, requestKeys)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "参数错误,请确认是从正确地方连接过来的", runData.logEntity, 8000600003, errs, e)
		return
	}

	// preparing object list data
	selectedCluster := strings.TrimSpace(requestData["clusterID"])
	if selectedCluster == "" {
		selectedCluster = "0"
	}

	// 初始化模板数据
	tplData, e := objectsUI.InitTemplateDataForWorkload("/"+defaultObjectName+"/", "集群管理", "命名空间列表", addButtonTitle, "no",
		allPopMenuItems, additionalJs, []string{}, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000600004, errs, e)
		return
	}
	tplData["objName"] = "namespace"

	// preparing datacenter data
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	conditions := make(map[string]string, 0)
	conditions["isDeleted"] = "=0"
	order := make(map[string]string, 0)
	var emptyString []string
	dcList, e := dcEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000600005, errs, e)
		return
	}

	e = buildSelectData(tplData, dcList, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000600006, errs, e)
		return
	}

	startPos := objectsUI.GetStartPosFromRequest(requestData)
	count, objListData, e := prepareNsData(selectedCluster, defaultOrderField, defaultOrderDirection, startPos, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000100008, errs, e)
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

func prepareNsData(selectedCluster, defaultOrderField, defaultOrderDirection string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	if selectedCluster == "" || selectedCluster == "0" {
		return 0, dataList, nil
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

	nsList, e := clientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return 0, dataList, e
	}

	totalNum := len(nsList.Items)
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

	var nsItems []interface{}
	for _, item := range nsList.Items {
		nsItems = append(nsItems, item)
	}

	if direction == "1" {
		switch orderfield {
		case "TD1":
			sort.Sort(objectsUI.SortData{Data: nsItems, By: sortNsByName})
		case "TD2":
			sort.Sort(objectsUI.SortData{Data: nsItems, By: sortNsByStatus})
		case "TD4":
			sort.Sort(objectsUI.SortData{Data: nsItems, By: sortNsByCreatetime})
		}
	} else {
		switch orderfield {
		case "TD1":
			sort.Sort(sort.Reverse(objectsUI.SortData{Data: nsItems, By: sortNsByName}))
		case "TD2":
			sort.Sort(sort.Reverse(objectsUI.SortData{Data: nsItems, By: sortNsByStatus}))
		case "TD4":
			sort.Sort(sort.Reverse(objectsUI.SortData{Data: nsItems, By: sortNsByCreatetime}))
		}
	}

	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}

	for i := startPos; i < endCount; i++ {
		interfaceData := nsItems[i]
		nsData, ok := interfaceData.(coreV1.Namespace)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not Namespace schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = nsData.Name
		lineMap["TD1"] = nsData.Name
		lineMap["TD2"] = nsData.Status.Phase
		lineMap["TD3"] = objectsUI.ConvertMap2HTML(nsData.Labels)
		lineMap["TD4"] = nsData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)

		popmenuitems, e := preparePopMenuItems(clientSet, nsData)
		if e != nil {
			return 0, dataList, e
		}

		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortNsByName(p, q interface{}) bool {
	pData, ok := p.(coreV1.Namespace)
	if !ok {
		return false
	}
	qData, ok := q.(coreV1.Namespace)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortNsByCreatetime(p, q interface{}) bool {
	pData, ok := p.(coreV1.Namespace)
	if !ok {
		return false
	}
	qData, ok := q.(coreV1.Namespace)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

func sortNsByStatus(p, q interface{}) bool {
	pData, ok := p.(coreV1.Namespace)
	if !ok {
		return false
	}
	qData, ok := q.(coreV1.Namespace)
	if !ok {
		return false
	}

	return pData.Status.Phase < qData.Status.Phase
}

func preparePopMenuItems(clientSet *kubernetes.Clientset, nsData coreV1.Namespace) (string, error) {
	popmenuitems := ""

	if nsData.Status.Phase == coreV1.NamespaceActive {
		popmenuitems = "0,1,2,3"
		quotaList, e := clientSet.CoreV1().ResourceQuotas(nsData.Name).List(context.Background(), metav1.ListOptions{})
		if e != nil {
			return "", e
		}

		if len(quotaList.Items) > 0 {
			popmenuitems = popmenuitems + ",4"
		}

		limitRangeList, e := clientSet.CoreV1().LimitRanges(nsData.Name).List(context.Background(), metav1.ListOptions{})
		if e != nil {
			return "", e
		}
		if len(limitRangeList.Items) > 0 {
			popmenuitems = popmenuitems + ",6,7,8"
		} else {
			popmenuitems = popmenuitems + ",5"
		}
	}
	return popmenuitems, nil
}
