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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/wangyysde/sysadmServer"
	apiCoreV1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyconfigCoreV1 "k8s.io/client-go/applyconfigurations/core/v1"
	"sort"
	"strings"
	"sysadm/k8sclient"
	"sysadm/objectsUI"
	"sysadm/utils"
)

func (s *secret) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortSecretByName, "TD6": sortSecretByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "类型", "TD4": "标签", "TD5": "数据项数", "TD6": "是否可编辑", "TD7": "创建时间"}
	additionalJs := []string{}
	additionalCss := []string{}
	templateFile := ""

	s.mainModuleName = "配置和存储"
	s.moduleName = "密文"
	s.allPopMenuItems = allPopMenuItems
	s.allListItems = allListItems
	s.addButtonTile = "添加密文"
	s.isSearchForm = "no"
	s.allOrderFields = allOrderFields
	s.defaultOrderField = "TD1"
	s.defaultOrderDirection = "1"
	s.moduleID = "secret"
	s.additionalJs = additionalJs
	s.additionalCss = additionalCss
	s.templateFile = templateFile
	s.namespaced = true
}

func (s *secret) getMainModuleName() string {
	return s.mainModuleName
}

func (s *secret) getModuleName() string {
	return s.moduleName
}

func (s *secret) getAddButtonTitle() string {
	return s.addButtonTile
}

func (s *secret) getIsSearchForm() string {
	return s.isSearchForm
}

func (s *secret) getAllPopMenuItems() []string {
	return s.allPopMenuItems
}

func (s *secret) getAllListItems() map[string]string {
	return s.allListItems
}

func (s *secret) getDefaultOrderField() string {
	return s.defaultOrderField
}

func (s *secret) getDefaultOrderDirection() string {
	return s.defaultOrderDirection
}

func (s *secret) getAllorderFields() map[string]objectsUI.SortBy {
	return s.allOrderFields
}

func (s *secret) getNamespaced() bool {
	return s.namespaced
}

// for Ingress
func (s *secret) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	nsStr, orderField, direction := checkRequestData(selectedNS, s.defaultOrderField, s.defaultOrderDirection, requestData)
	secretList, e := clientSet.CoreV1().Secrets(nsStr).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return 0, dataList, e
	}
	totalNum := len(secretList.Items)
	if totalNum < 1 {
		return 0, dataList, nil
	}
	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}

	var secretItems []interface{}
	for _, item := range secretList.Items {
		secretItems = append(secretItems, item)
	}

	moduleAllOrderFields := s.allOrderFields
	for field, fn := range moduleAllOrderFields {
		if field == orderField {
			if direction == "1" {
				sort.Sort(objectsUI.SortData{Data: secretItems, By: fn})
			} else {
				sort.Sort(sort.Reverse(objectsUI.SortData{Data: secretItems, By: fn}))
			}

		}
	}

	for i := startPos; i < endCount; i++ {
		interfaceData := secretItems[i]
		secretData, ok := interfaceData.(corev1.Secret)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not Secret schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = secretData.Name
		lineMap["TD1"] = secretData.Name
		lineMap["TD2"] = secretData.Namespace
		lineMap["TD3"] = secretData.Type
		lineMap["TD4"] = objectsUI.ConvertMap2HTML(secretData.Labels)
		dataCount := len(secretData.Data)
		stringDataCount := len(secretData.StringData)
		totalNum := dataCount + stringDataCount
		lineMap["TD5"] = totalNum
		editable := "是"
		popmenuitems := "0,1"
		if secretData.Immutable != nil && *secretData.Immutable {
			editable = "否"
			popmenuitems = "0,1"
		}
		lineMap["TD6"] = editable
		lineMap["TD7"] = secretData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortSecretByName(p, q interface{}) bool {
	pData, ok := p.(corev1.Secret)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.Secret)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortSecretByCreatetime(p, q interface{}) bool {
	pData, ok := p.(corev1.Secret)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.Secret)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

func (s *secret) getModuleID() string {
	return s.moduleID
}

func (s *secret) buildAddFormData(tplData map[string]interface{}) error {
	tplData["thirdCategory"] = "创建密文"
	formData, e := objectsUI.InitFormData("addSecret", "addSecret", "POST", "_self", "yes", "addWorkload", "")
	if e != nil {
		return e
	}
	tplData["formData"] = formData

	clusterID := tplData["clusterID"].(string)
	if clusterID == "" || clusterID == "0" {
		return fmt.Errorf("cluster must be specified when add a new secret")
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
		for _, n := range denyStatefulSetWokloadNSList {
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
	lineData := objectsUI.InitLineData("basicInfoID", false, false, false)
	_ = objectsUI.AddLabelData("basicInfoID", "mid", "Left", "基本信息", false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("namespaceSelectLine", false, false, false)
	e = objectsUI.AddSelectData("nsSelectedID", "nsSelected", defaultSelectedNs, "", "", "选择命名空间", "", 1, false, false, nsOptions, lineData)
	if e != nil {
		return e
	}
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("nameLine", false, false, false)
	_ = objectsUI.AddTextData("name", "name", "", "字典名称", "validateNewName", "addWorkloadValidateNewName", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
	var immutableOptions []objectsUI.Option
	immutableOptions, _ = objectsUI.AddCheckBoxOption("创建后不可修改", "1", false, false, immutableOptions)
	_ = objectsUI.AddCheckBoxData("immutableID", "immutable", "", "", false, immutableOptions, lineData)
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

	lineData = objectsUI.InitLineData("configMapDataZoneID", false, false, false)
	_ = objectsUI.AddLabelData("configMapDataZoneID", "mid", "Left", "字典数据", false, lineData)
	basicData = append(basicData, lineData)

	var secretTypeOptions []objectsUI.Option
	secretTypeOption := objectsUI.Option{Text: "Opaque", Value: "0", Checked: true, Disabled: false}
	secretTypeOptions = append(secretTypeOptions, secretTypeOption)
	secretTypeOption = objectsUI.Option{Text: "镜像仓库密码", Value: "1", Checked: false, Disabled: false}
	secretTypeOptions = append(secretTypeOptions, secretTypeOption)
	secretTypeOption = objectsUI.Option{Text: "Service Account Token", Value: "2", Checked: false, Disabled: false}
	secretTypeOptions = append(secretTypeOptions, secretTypeOption)
	secretTypeOption = objectsUI.Option{Text: "TLS", Value: "3", Checked: false, Disabled: false}
	secretTypeOptions = append(secretTypeOptions, secretTypeOption)
	lineData = objectsUI.InitLineData("secretTypeLineID", false, false, false)
	_ = objectsUI.AddRadioData("secretTypeID", "secretType", "密文类型", "secretTypeChangedForAddSecret", false, secretTypeOptions, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("OpaquekeyID", true, true, false)
	_ = objectsUI.AddTextData("opaquekeyID[]", "opaquekey[]", "", "Key", "", "", "", 20, false, false, lineData)
	_ = objectsUI.AddTextareaData("opaqueDataID[]", "opaqueData[]", "", "  数据", "", "", "", 60, 5, false, false, lineData)
	_ = objectsUI.AddWordsInputData("selectorLabel[]", "selectorLabel", "fa-trash", "#", "workloadDelSelector", false, true, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("OpaqueAnchorID", false, false, false)
	_ = objectsUI.AddWordsInputData("opaquekeyID", "opaquekey", "添加数据项", "#", "workloadAddSelectorBlock", false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("registryServerLine", false, false, true)
	_ = objectsUI.AddTextData("registryServerID", "registryServer", "", "镜像仓库地址:", "", "", "不带HTTP和HTTPS的镜像仓库访问地址,如hb.sysad.cn", 30, false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("registryUsernameLine", false, false, true)
	_ = objectsUI.AddTextData("registryUsernameID", "registryusername", "", "用户名:", "", "", "登陆镜像仓库的用户名", 30, false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("registryPasswordLine", false, false, true)
	_ = objectsUI.AddTextData("registryPasswordID", "registryPassword", "", "用户名:", "", "", "登陆镜像仓库的密码", 30, false, false, lineData)
	basicData = append(basicData, lineData)

	saOptions := make(map[string]string, 0)
	saOptions["0"] = "===选择SA==="
	lineData = objectsUI.InitLineData("saSelectLine", false, false, true)
	_ = objectsUI.AddSelectData("saSelectID", "saSelected", "0", "", "", "选择SA", "", 1, false, false, saOptions, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("certLine", false, false, true)
	_ = objectsUI.AddTextareaData("certID[]", "cert", "", "证书内容: ", "", "", "", 60, 5, false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("keyLine", false, false, true)
	_ = objectsUI.AddTextareaData("keyID[]", "key", "", "密钥内容: ", "", "", "", 60, 5, false, false, lineData)
	basicData = append(basicData, lineData)

	tplData["BasicData"] = basicData
	return nil
}

func (s *secret) getAdditionalJs() []string {
	return s.additionalJs
}
func (s *secret) getAdditionalCss() []string {
	return s.additionalCss
}

func (s *secret) addNewResource(c *sysadmServer.Context, module string) error {
	requestKeys := []string{"dcid", "clusterID", "namespace", "addType", "nsSelected", "name"}
	requestKeys = append(requestKeys, "immutable", "secretType", "opaquekey[]", "opaqueData[]", "registryServer")
	requestKeys = append(requestKeys, "registryusername", "registryPassword", "saSelected", "cert", "key")
	requestKeys = append(requestKeys, "labelKey[]", "labelValue[]", "annotationKey[]", "annotationValue[]")
	formData, e := utils.GetMultipartData(c, requestKeys)
	if e != nil {
		return e
	}
	ns := formData["nsSelected"].([]string)
	name := formData["name"].([]string)
	secretApplyConfig := applyconfigCoreV1.Secret(name[0], ns[0])

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
		labels[defaultSecretLabelKey] = name[0]
	}
	for k, v := range extraLabels {
		labels[k] = v
	}
	secretApplyConfig = secretApplyConfig.WithLabels(labels)

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
	secretApplyConfig = secretApplyConfig.WithAnnotations(annotations)

	immutable := formData["immutable"].([]string)
	isMmutable := false
	if len(immutable) > 0 {
		immutableStr := strings.TrimSpace(immutable[0])
		if immutableStr == "1" {
			isMmutable = true
		}
	}
	secretApplyConfig.Immutable = &isMmutable

	secretData := make(map[string][]byte, 0)
	secretStringData := make(map[string]string, 0)

	secretTypeStr := strings.TrimSpace(formData["secretType"].([]string)[0])
	var secretType *apiCoreV1.SecretType
	switch secretTypeStr {
	case "0":
		tmpSecretType := apiCoreV1.SecretTypeOpaque
		secretType = &tmpSecretType
		opaquekeys := formData["opaquekey[]"].([]string)
		opaqueDatas := formData["opaqueData[]"].([]string)
		var addedKeys []string
		for i, k := range opaquekeys {
			k = strings.TrimSpace(k)
			if utils.FoundStrInSlice(addedKeys, k, true) {
				return fmt.Errorf("the data key %s is duplicate")
			}
			addedKeys = append(addedKeys, k)
			data := strings.TrimSpace(opaqueDatas[i])
			secretStringData[k] = data
		}
	case "1":
		tmpSecretType := apiCoreV1.SecretTypeDockerConfigJson
		secretType = &tmpSecretType
		registryusername := strings.TrimSpace(formData["registryusername"].([]string)[0])
		registryPassword := strings.TrimSpace(formData["registryPassword"].([]string)[0])
		registryServer := strings.TrimSpace(formData["registryServer"].([]string)[0])
		authStr := registryusername + ":" + registryPassword
		authEncodeStr := base64.StdEncoding.EncodeToString([]byte(authStr))
		auths := RegistryAuths{}
		auth := RegistryHost{registryServer: RegistryAuth{Auth: authEncodeStr}}
		auths.Auths = auth
		jsonStr, e := json.Marshal(auths)
		if e != nil {
			return e
		}
		secretStringData[".dockerconfigjson"] = string(jsonStr)
	case "2":
		saSelected := strings.TrimSpace(formData["saSelected"].([]string)[0])
		if saSelected == "0" {
			return fmt.Errorf("service accout %s is not valid")
		}
		tmpSecretType := apiCoreV1.SecretTypeServiceAccountToken
		secretType = &tmpSecretType
		annotations := secretApplyConfig.Annotations
		annotations["kubernetes.io/service-account.name"] = saSelected
		secretApplyConfig = secretApplyConfig.WithAnnotations(annotations)
	case "3":
		cert := strings.TrimSpace(formData["cert"].([]string)[0])
		key := strings.TrimSpace(formData["key"].([]string)[0])
		tmpSecretType := apiCoreV1.SecretTypeTLS
		secretType = &tmpSecretType
		certEncodeStr := base64.StdEncoding.EncodeToString([]byte(cert))
		keyEncodeStr := base64.StdEncoding.EncodeToString([]byte(key))
		secretData["tls.crt"] = []byte(certEncodeStr)
		secretData["tls.key"] = []byte(keyEncodeStr)
	}
	secretApplyConfig.Type = secretType
	secretApplyConfig.Data = secretData
	secretApplyConfig.StringData = secretStringData

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

	_, e = clientSet.CoreV1().Secrets(ns[0]).Apply(context.Background(), secretApplyConfig, applyOption)

	return e
}

func (s *secret) delResource(c *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (s *secret) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

func (s *secret) getTemplateFile(action string) string {
	switch action {
	case "list":
		return secretTemplateFiles["list"]
	case "addform":
		return secretTemplateFiles["addform"]
	default:
		return ""

	}
	return ""

}
