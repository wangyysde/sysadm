/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyconfigBatchv1 "k8s.io/client-go/applyconfigurations/batch/v1"
	"strconv"
	"strings"
	"sysadm/k8sclient"
	"sysadm/objectsUI"
	"sysadm/utils"
	"time"
)

func (j *job) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortJobByName, "TD7": sortJobByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "持续时间", "TD4": "标签", "TD5": "Pods(活/成/失/总）", "TD6": "完成时间", "TD7": "创建时间"}
	additionalJs := []string{"js/sysadmfunctions.js", "/js/workloadList.js"}
	additionalCss := []string{}
	templateFile := "addWorkload.html"

	j.mainModuleName = "工作负载"
	j.moduleName = "任务"
	j.allPopMenuItems = allPopMenuItems
	j.allListItems = allListItems
	j.addButtonTile = "创建任务"
	j.isSearchForm = "no"
	j.allOrderFields = allOrderFields
	j.defaultOrderField = "TD1"
	j.defaultOrderDirection = "1"
	j.namespaced = true
	j.moduleID = "job"
	j.additionalJs = additionalJs
	j.additionalCss = additionalCss
	j.templateFile = templateFile

}

func (j *job) getMainModuleName() string {
	return j.mainModuleName
}

func (j *job) getModuleName() string {
	return j.moduleName
}

func (j *job) getAddButtonTitle() string {
	return j.addButtonTile
}

func (j *job) getIsSearchForm() string {
	return j.isSearchForm
}

func (j *job) getAllPopMenuItems() []string {
	return j.allPopMenuItems
}

func (j *job) getAllListItems() map[string]string {
	return j.allListItems
}

func (j *job) getDefaultOrderField() string {
	return j.defaultOrderField
}

func (j *job) getDefaultOrderDirection() string {
	return j.defaultOrderDirection
}

func (j *job) getAllorderFields() map[string]objectsUI.SortBy {
	return j.allOrderFields
}

func (j *job) getNamespaced() bool {
	return j.namespaced
}

// for daemonSet
func (j *job) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	if selectedCluster == "" || selectedCluster == "0" {
		return 0, dataList, nil
	}

	nsStr := ""
	if selectedNS != "0" {
		nsStr = selectedNS
	}

	clientSet, e := buildClientSetByClusterID(selectedCluster)
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
		orderfield = j.getDefaultOrderField()
	}
	if direction == "" || (direction != "0" && direction != "1") {
		direction = j.getDefaultOrderDirection()
	}

	var jobItems []interface{}
	for _, item := range jobList.Items {
		jobItems = append(jobItems, item)
	}

	sortWorkloadData(jobItems, direction, orderfield, j.getAllorderFields())

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

func (j *job) getModuleID() string {
	return j.moduleID
}

func (j *job) buildAddFormData(tplData map[string]interface{}) error {
	tplData["thirdCategory"] = "创建任务"
	formData, e := objectsUI.InitFormData("addJob", "addJob", "POST", "_self", "yes", "addWorkload", "")
	if e != nil {
		return e
	}
	tplData["formData"] = formData

	e = buildJobBasiceFormData(tplData)
	if e != nil {
		return e
	}

	e = builtContainerFormData(tplData)
	if e != nil {
		return e
	}

	e = buildStorageFormData(tplData)
	if e != nil {
		return e
	}

	return nil
}

func (j *job) getAdditionalJs() []string {
	return j.additionalJs
}
func (j *job) getAdditionalCss() []string {
	return j.additionalCss
}

func (j *job) addNewResource(c *sysadmServer.Context, module string) error {
	requestKeys := []string{"dcid", "clusterID", "namespace", "addType", "nsSelected", "name", "restartPolicySelected", "completion", "parallelism", "labelKey[]", "labelValue[]", "annotationKey[]", "annotationValue[]"}
	requestKeys = append(requestKeys, "selectorKey[]")
	requestKeys = append(requestKeys, "selectorValue[]")
	requestKeys = append(requestKeys, "containerData[]")
	requestKeys = append(requestKeys, "storageMountData[]")
	formData, e := utils.GetMultipartData(c, requestKeys)
	if e != nil {
		return e
	}

	ns := formData["nsSelected"].([]string)
	name := formData["name"].([]string)
	jobApplyConfig := applyconfigBatchv1.Job(name[0], ns[0])

	// 配置labels
	labelKeys := formData["labelKey[]"].([]string)
	labelValue := formData["labelValue[]"].([]string)
	if len(labelKeys) != len(labelValue) {
		return fmt.Errorf("label's key is not equal to label's value")
	}

	labels := make(map[string]string, 0)
	if len(labelKeys) > 0 {
		for i, k := range labelKeys {
			value := labelValue[i]
			labels[k] = value
		}
	} else {
		labels[defaultLabelKey] = name[0]
	}
	for k, v := range extraLabels {
		labels[k] = v
	}
	jobApplyConfig = jobApplyConfig.WithLabels(labels)

	// 配置注解
	annotationKey := formData["annotationKey[]"].([]string)
	annotationValue := formData["annotationValue[]"].([]string)
	if len(annotationKey) != len(annotationValue) {
		return fmt.Errorf("annotation's key is not equal to annotation's value")
	}
	annotations := make(map[string]string, 0)
	for i, k := range annotationKey {
		value := annotationValue[i]
		annotations[k] = value
	}
	jobApplyConfig = jobApplyConfig.WithAnnotations(annotations)

	jobSpecApplyConfig := applyconfigBatchv1.JobSpecApplyConfiguration{}
	completion := strings.TrimSpace(formData["completion"].([]string)[0])
	parallelism := strings.TrimSpace(formData["parallelism"].([]string)[0])
	if completion != "" && parallelism != "" {
		completionInt, e1 := strconv.Atoi(completion)
		parallelismInt, e2 := strconv.Atoi(parallelism)
		if e1 != nil || e2 != nil {
			return fmt.Errorf("completion or parallelism is error")
		}
		completionInt32 := int32(completionInt)
		parallelismInt32 := int32(parallelismInt)
		jobSpecApplyConfig.Completions = &completionInt32
		jobSpecApplyConfig.Parallelism = &parallelismInt32
	}

	podTemplateSpecApplyConfiguration, e := buildPodTemplateSpecApplyConfig(formData, labels, annotations)
	if e != nil {
		return e
	}
	podRestartPolicy := apicorev1.RestartPolicyOnFailure
	restartPolicy := strings.TrimSpace(formData["restartPolicySelected"].([]string)[0])
	if strings.ToLower(restartPolicy) == strings.ToLower(strings.TrimSpace(string(apicorev1.RestartPolicyNever))) {
		podRestartPolicy = apicorev1.RestartPolicyNever
	}
	podTemplateSpecApplyConfiguration.Spec.RestartPolicy = &podRestartPolicy

	jobSpecApplyConfig.Template = podTemplateSpecApplyConfiguration
	jobApplyConfig = jobApplyConfig.WithSpec(&jobSpecApplyConfig)

	clusterIDSlice := formData["clusterID"].([]string)
	clusterID := clusterIDSlice[0]
	clientSet, e := buildClientSetByClusterID(clusterID)
	if e != nil {
		return e
	}
	applyOption := metav1.ApplyOptions{
		Force:        true,
		FieldManager: k8sclient.FieldManager,
	}

	_, e = clientSet.BatchV1().Jobs(ns[0]).Apply(context.Background(), jobApplyConfig, applyOption)

	return e

}

func (j *job) delResource(s *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (j *job) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

func (j *job) getTemplateFile(action string) string {
	switch action {
	case "list":
		return jobTemplateFiles["list"]
	case "addform":
		return jobTemplateFiles["addform"]
	default:
		return ""

	}
	return ""
}

func buildJobBasiceFormData(tplData map[string]interface{}) error {
	clusterID := tplData["clusterID"].(string)
	if clusterID == "" || clusterID == "0" {
		return fmt.Errorf("cluster must be specified when add a new job")
	}
	clientSet, e := buildClientSetByClusterID(clusterID)
	if e != nil {
		return e
	}
	nsObjectList, e := clientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return e
	}
	nsExistList := []string{}
	for _, item := range nsObjectList.Items {
		found := false
		for _, n := range denyDeployWokloadNSList {
			if strings.TrimSpace(strings.ToLower(item.Name)) == n {
				found = true
			}
		}

		if found {
			continue
		}
		nsExistList = append(nsExistList, item.Name)

	}

	// 准备命名空间下拉菜单option数据
	nsOptions := make(map[string]string, 0)
	defaultSelectedNs := ""
	for _, n := range nsExistList {
		nsOptions[n] = n
		if defaultSelectedNs == "" {
			defaultSelectedNs = n
		}
	}

	// 准备基本信息表单内容
	var basicData []interface{}
	lineData := objectsUI.InitLineData("namespaceSelectLine", false, false, false)
	e = objectsUI.AddSelectData("nsSelectedID", "nsSelected", defaultSelectedNs, "", "", "选择命名空间", "", 1, false, false, nsOptions, lineData)
	if e != nil {
		return e
	}
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("nameLine", false, false, false)
	_ = objectsUI.AddTextData("name", "name", "", "任务名称", "validateNewName", "addWorkloadValidateNewName", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
	basicData = append(basicData, lineData)

	restartPolicyOptions := make(map[string]string, 0)
	restartPolicyOptions["OnFailure"] = "运行失败时重启"
	restartPolicyOptions["Never"] = "从不重启"
	lineData = objectsUI.InitLineData("restartPolicySelectLine", false, false, false)
	_ = objectsUI.AddSelectData("restartPolicySelectedID", "restartPolicySelected", "OnFailure", "", "", "选择容器故障时重启策略", "", 1, false, false, restartPolicyOptions, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("completionsLine", false, false, false)
	_ = objectsUI.AddTextData("completion", "completion", "", "预定成功容器组数", "", "", "值为大于0的正整数。", 30, false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("parallelismLine", false, false, false)
	_ = objectsUI.AddTextData("parallelism", "parallelism", "", "并行容器组数", "", "", "并行运行容器组伯数量，其值为大于0小于等于预定成功容器组数之间的正整数", 30, false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("newNsLabel", true, true, false)
	_ = objectsUI.AddTextData("labelKey", "labelKey[]", "", "标签", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("equal", "equal", "=", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("labelValue", "labelValue[]", "", "值", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("delLabel", "delLabel", "fa-trash", "#", "workloadDelLabel", false, true, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("addlabelanchor", false, false, false)
	_ = objectsUI.AddWordsInputData("addLabel", "addLabel", "增加标签", "#", "workloadAddLabel", false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("newAnnotationLabel", true, true, false)
	_ = objectsUI.AddTextData("annotationKey", "annotationKey[]", "", "注解", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("equal", "equal", "=", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("annotationValue", "annotationValue[]", "", "值", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("annotationLabel", "annotationLabel", "fa-trash", "#", "workloadDelAnnotaion", false, true, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("annotationanchor", false, false, false)
	_ = objectsUI.AddWordsInputData("annotationLabel", "annotationLabel", "增加注解", "#", "workloadAddAnnotation", false, false, lineData)
	basicData = append(basicData, lineData)

	tplData["BasicData"] = basicData

	return nil
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
