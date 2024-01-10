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
	apiBatchV1 "k8s.io/api/batch/v1"
	batchv1 "k8s.io/api/batch/v1"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyconfigBatchv1 "k8s.io/client-go/applyconfigurations/batch/v1"
	"strconv"
	"strings"
	"sysadm/k8sclient"
	"sysadm/objectsUI"
	"sysadm/utils"
)

func (c *cronjob) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortCronJobByName, "TD6": sortCronJobByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "标签", "TD4": "上次调度时间", "TD5": "上次成功执行时间", "TD6": "创建时间"}
	additionalJs := []string{"js/sysadmfunctions.js", "/js/workloadList.js"}
	additionalCss := []string{}
	templateFile := "addWorkload.html"

	c.mainModuleName = "工作负载"
	c.moduleName = "定时任务"
	c.allPopMenuItems = allPopMenuItems
	c.allListItems = allListItems
	c.addButtonTile = "创建定时任务"
	c.isSearchForm = "no"
	c.allOrderFields = allOrderFields
	c.defaultOrderField = "TD1"
	c.defaultOrderDirection = "1"
	c.namespaced = true
	c.moduleID = "cronjob"
	c.additionalJs = additionalJs
	c.additionalCss = additionalCss
	c.templateFile = templateFile

}

func (c *cronjob) getMainModuleName() string {
	return c.mainModuleName
}

func (c *cronjob) getModuleName() string {
	return c.moduleName
}

func (c *cronjob) getAddButtonTitle() string {
	return c.addButtonTile
}

func (c *cronjob) getIsSearchForm() string {
	return c.isSearchForm
}

func (c *cronjob) getAllPopMenuItems() []string {
	return c.allPopMenuItems
}

func (c *cronjob) getAllListItems() map[string]string {
	return c.allListItems
}

func (c *cronjob) getDefaultOrderField() string {
	return c.defaultOrderField
}

func (c *cronjob) getDefaultOrderDirection() string {
	return c.defaultOrderDirection
}

func (c *cronjob) getAllorderFields() map[string]objectsUI.SortBy {
	return c.allOrderFields
}

func (c *cronjob) getNamespaced() bool {
	return c.namespaced
}

// for daemonSet
func (c *cronjob) listObjectData(selectedCluster, selectedNS string,
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
		orderfield = c.getDefaultOrderField()
	}
	if direction == "" || (direction != "0" && direction != "1") {
		direction = c.getDefaultOrderDirection()
	}

	var cronJobItems []interface{}
	for _, item := range cronJobList.Items {
		cronJobItems = append(cronJobItems, item)
	}

	sortWorkloadData(cronJobItems, direction, orderfield, c.getAllorderFields())

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
		if cronJobData.Status.LastScheduleTime != nil {
			lineMap["TD4"] = cronJobData.Status.LastScheduleTime.Format(objectsUI.DefaultTimeStampFormat)
		} else {
			lineMap["TD4"] = "---"
		}
		if cronJobData.Status.LastSuccessfulTime != nil {
			lineMap["TD5"] = cronJobData.Status.LastSuccessfulTime.Format(objectsUI.DefaultTimeStampFormat)
		} else {
			lineMap["TD5"] = "---"
		}
		lineMap["TD6"] = cronJobData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := "0,1"
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func (c *cronjob) getModuleID() string {
	return c.moduleID
}

func (c *cronjob) buildAddFormData(tplData map[string]interface{}) error {
	tplData["thirdCategory"] = "创建定时任务"
	formData, e := objectsUI.InitFormData("addCrontJob", "addCrontJob", "POST", "_self", "yes", "addWorkload", "")
	if e != nil {
		return e
	}
	tplData["formData"] = formData

	e = buildCrontJobBasiceFormData(tplData)
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

func (c *cronjob) getAdditionalJs() []string {
	return c.additionalJs
}
func (c *cronjob) getAdditionalCss() []string {
	return c.additionalCss
}

func (c *cronjob) addNewResource(cc *sysadmServer.Context, module string) error {
	requestKeys := []string{"dcid", "clusterID", "namespace", "addType", "nsSelected", "name", "restartPolicySelected", "completion", "parallelism", "labelKey[]", "labelValue[]", "annotationKey[]", "annotationValue[]"}
	requestKeys = append(requestKeys, "selectorKey[]")
	requestKeys = append(requestKeys, "selectorValue[]")
	requestKeys = append(requestKeys, "containerData[]")
	requestKeys = append(requestKeys, "storageMountData[]")
	requestKeys = append(requestKeys, "storageMountData[]")
	requestKeys = append(requestKeys, "concurrencyPolicy", "minute", "hour", "day", "month", "week")
	formData, e := utils.GetMultipartData(cc, requestKeys)
	if e != nil {
		return e
	}

	ns := formData["nsSelected"].([]string)
	name := formData["name"].([]string)
	cronjobApplyConfig := applyconfigBatchv1.CronJob(name[0], ns[0])

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
	cronjobApplyConfig = cronjobApplyConfig.WithLabels(labels)

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
	cronjobApplyConfig = cronjobApplyConfig.WithAnnotations(annotations)

	crontjobSpecApplyConfig := applyconfigBatchv1.CronJobSpecApplyConfiguration{}
	concurrencyPolicy := formData["concurrencyPolicy"].([]string)[0]
	concurrencyPolicyV1 := batchv1.ConcurrencyPolicy(concurrencyPolicy)
	crontjobSpecApplyConfig.ConcurrencyPolicy = &concurrencyPolicyV1

	minute := formData["minute"].([]string)[0]
	minute = strings.TrimSpace(minute)
	if minute == "" {
		minute = "*"
	}
	hour := formData["hour"].([]string)[0]
	hour = strings.TrimSpace(hour)
	if hour == "" {
		hour = "*"
	}
	day := formData["day"].([]string)[0]
	day = strings.TrimSpace(day)
	if day == "" {
		day = "*"
	}
	month := formData["month"].([]string)[0]
	month = strings.TrimSpace(month)
	if month == "" {
		month = "*"
	}
	week := formData["week"].([]string)[0]
	week = strings.TrimSpace(week)
	if week == "" {
		week = "*"
	}
	scheduleStr := minute + " " + hour + " " + day + " " + month + " " + week
	crontjobSpecApplyConfig.Schedule = &scheduleStr

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

	jobTemplateSpec := applyconfigBatchv1.JobTemplateSpecApplyConfiguration{Spec: &jobSpecApplyConfig}
	crontjobSpecApplyConfig.JobTemplate = &jobTemplateSpec
	cronjobApplyConfig.WithSpec(&crontjobSpecApplyConfig)

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

	_, e = clientSet.BatchV1().CronJobs(ns[0]).Apply(context.Background(), cronjobApplyConfig, applyOption)

	return e

	return nil

}

func (c *cronjob) delResource(s *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (c *cronjob) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

func (c *cronjob) getTemplateFile(action string) string {
	switch action {
	case "list":
		return crontjobTemplateFiles["list"]
	case "addform":
		return crontjobTemplateFiles["addform"]
	default:
		return ""

	}
	return ""
}

func buildCrontJobBasiceFormData(tplData map[string]interface{}) error {
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
	_ = objectsUI.AddSelectData("nsSelectedID", "nsSelected", defaultSelectedNs, "", "", "选择命名空间", "", 1, false, false, nsOptions, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("nameLine", false, false, false)
	_ = objectsUI.AddTextData("name", "name", "", "任务名称", "validateNewName", "addWorkloadValidateNewName", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("scheduleLine", false, false, false)
	_ = objectsUI.AddTextData("minute", "minute", "*", "运行计划  分钟", "", "", "", 5, false, false, lineData)
	_ = objectsUI.AddTextData("hour", "hour", "*", "小时", "", "", "", 5, false, false, lineData)
	_ = objectsUI.AddTextData("day", "day", "*", "日", "", "", "", 5, false, false, lineData)
	_ = objectsUI.AddTextData("month", "month", "*", "月", "", "", "", 5, false, false, lineData)
	_ = objectsUI.AddTextData("week", "week", "*", "星期", "", "", "计划格式遵守Linux系统下Crontab格式，如0 0 13 * 5", 5, false, false, lineData)
	basicData = append(basicData, lineData)

	concurrencyPolicyOptions := make(map[string]string, 0)
	concurrencyPolicyOptions[string(apiBatchV1.AllowConcurrent)] = "允许并发运行"
	concurrencyPolicyOptions[string(apiBatchV1.ForbidConcurrent)] = "禁止并发运行"
	concurrencyPolicyOptions[string(apiBatchV1.ReplaceConcurrent)] = "替换未完成任务"
	lineData = objectsUI.InitLineData("concurrencyPolicyLine", false, false, false)
	_ = objectsUI.AddSelectData("concurrencyPolicyID", "concurrencyPolicy", string(apiBatchV1.AllowConcurrent), "", "", "并行运行策略", "是指在前后多个计划周期内,当前一个执行周期的任务还没有执行完成时,如何对待已经来临的新的计划周期的任务的策略", 1, false, false, concurrencyPolicyOptions, lineData)
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
