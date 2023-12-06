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
	corev1 "k8s.io/api/core/v1"
	"sysadm/objectsUI"
)

var allQuotaFormItems = map[string]QuotaFormItem{
	"cpuRequest":     {ResourceName: corev1.ResourceRequestsCPU, UintStr: ""},
	"cpuLimit":       {ResourceName: corev1.ResourceLimitsCPU, UintStr: ""},
	"memRequest":     {ResourceName: corev1.ResourceRequestsMemory, UintStr: "Mi"},
	"memLimit":       {ResourceName: corev1.ResourceLimitsMemory, UintStr: "Mi"},
	"storageRequest": {ResourceName: corev1.ResourceRequestsStorage, UintStr: "Gi"},
	"pvcNum":         {ResourceName: corev1.ResourcePersistentVolumeClaims, UintStr: ""},
	"podNum":         {ResourceName: corev1.ResourcePods, UintStr: ""},
	"serviceNum":     {ResourceName: corev1.ResourceServices, UintStr: ""},
	"secretNum":      {ResourceName: corev1.ResourceSecrets, UintStr: ""},
	"configMapNum":   {ResourceName: corev1.ResourceConfigMaps, UintStr: ""},
}

var allLimitRangeFormItems = map[string]QuotaFormItem{
	"cpuMin":         {ResourceName: corev1.ResourceCPU, UintStr: ""},
	"cpuMax":         {ResourceName: corev1.ResourceCPU, UintStr: ""},
	"cpuDefault":     {ResourceName: corev1.ResourceCPU, UintStr: ""},
	"memMin":         {ResourceName: corev1.ResourceMemory, UintStr: "Mi"},
	"memMax":         {ResourceName: corev1.ResourceMemory, UintStr: "Mi"},
	"memDefault":     {ResourceName: corev1.ResourceMemory, UintStr: "Mi"},
	"storageMin":     {ResourceName: corev1.ResourceStorage, UintStr: "Gi"},
	"storageMax":     {ResourceName: corev1.ResourceStorage, UintStr: "Gi"},
	"storageDefault": {ResourceName: corev1.ResourceStorage, UintStr: "Gi"},
}

var quotaListPagePopmenu = []string{"详情,quotaDetail,GET,poppage", "删除,delQuota,POST,tip", "编辑,editQuota,GET,page"}
var quotaListAllListItems = map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "CPU请求", "TD4": "CPU上限", "TD5": "内存请求", "TD6": "内存上限", "TD7": "创建时间"}
var quotaListDefaultOrderField = "TD1"
var quotaListDefaultOrderDirection = "1"
var quotaListAllOrderFields = map[string]objectsUI.SortBy{"TD1": sortQuotaByName, "TD7": sortQuotaByCreatetime}
