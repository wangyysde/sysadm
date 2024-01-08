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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"mime/multipart"
	"strconv"
	"strings"
	datacenter "sysadm/datacenter/app"
	sysadmK8sClient "sysadm/k8sclient"
	sysadmK8sCluster "sysadm/k8scluster/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/user"
	"sysadm/utils"
)

func buildSelectDataWithNs(tplData map[string]interface{}, dcList []interface{}, requestData map[string]string) error {

	selectedDC := strings.TrimSpace(requestData["dcID"])
	if selectedDC == "" {
		selectedDC = "0"
	}
	selectedCluster := strings.TrimSpace(requestData["clusterID"])
	if selectedCluster == "" {
		selectedCluster = "0"
	}

	selectedNamespace := strings.TrimSpace(requestData["namespace"])
	if selectedNamespace == "" {
		selectedNamespace = "0"
	}

	// preparing datacenter options
	var dcOptions []objectsUI.SelectOption
	dcOption := objectsUI.SelectOption{
		Id:       "0",
		Text:     "===选择数据中心===",
		ParentID: "0",
	}
	dcOptions = append(dcOptions, dcOption)
	for _, line := range dcList {
		dcData, ok := line.(datacenter.DatacenterSchema)
		if !ok {
			return fmt.Errorf("the data is not datacenter schema")
		}
		dcOption := objectsUI.SelectOption{
			Id:       strconv.Itoa(int(dcData.Id)),
			Text:     dcData.CnName,
			ParentID: "0",
		}
		dcOptions = append(dcOptions, dcOption)
	}
	dcSelect := objectsUI.SelectData{Title: "数据中心", SelectedId: selectedDC, Options: dcOptions}
	tplData["dcSelect"] = dcSelect

	// preparing cluster options
	var clusterOptions []objectsUI.SelectOption
	clusterOption := objectsUI.SelectOption{
		Id:       "0",
		Text:     "===选择集群===",
		ParentID: "0",
	}
	clusterOptions = append(clusterOptions, clusterOption)
	if selectedDC != "0" {
		var k8sclusterEntity sysadmObjects.ObjectEntity
		k8sclusterEntity = sysadmK8sCluster.New()
		conditions := make(map[string]string, 0)
		var emptyString []string
		conditions["isDeleted"] = "='0'"
		conditions["dcid"] = "='" + selectedDC + "'"
		clusterList, e := k8sclusterEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, make(map[string]string))
		if e != nil {
			return e
		}

		for _, line := range clusterList {
			clusterData, ok := line.(sysadmK8sCluster.K8sclusterSchema)
			if !ok {
				return fmt.Errorf("the data is not cluster schema")
			}

			clusterOption := objectsUI.SelectOption{
				Id:       clusterData.Id,
				Text:     clusterData.CnName,
				ParentID: strconv.Itoa(int(clusterData.Dcid)),
			}
			clusterOptions = append(clusterOptions, clusterOption)
		}
	}
	clusterSelect := objectsUI.SelectData{Title: "集群", SelectedId: selectedCluster, Options: clusterOptions}
	tplData["clusterSelect"] = clusterSelect

	// preparing namespace options
	var namespaceOptions []objectsUI.SelectOption
	if selectedCluster != "0" {
		namespaceOption := objectsUI.SelectOption{
			Id:       "0",
			Text:     "所有命名空间",
			ParentID: "0",
		}
		namespaceOptions = append(namespaceOptions, namespaceOption)
		var k8sclusterEntity sysadmObjects.ObjectEntity
		k8sclusterEntity = sysadmK8sCluster.New()
		clusterInfo, e := k8sclusterEntity.GetObjectInfoByID(selectedCluster)
		if e != nil {
			return e
		}
		clusterData, ok := clusterInfo.(sysadmK8sCluster.K8sclusterSchema)
		if !ok {
			return fmt.Errorf("the data is not K8S data schema")
		}
		ca := []byte(clusterData.Ca)
		cert := []byte(clusterData.Cert)
		key := []byte(clusterData.Key)
		restConf, e := sysadmK8sClient.BuildConfigFromParametes(ca, cert, key, clusterData.Apiserver, clusterData.Id, clusterData.ClusterUser, "default")
		if e != nil {
			return e
		}

		clientSet, e := sysadmK8sClient.BuildClientset(restConf)
		if e != nil {
			return e
		}
		nsList, e := clientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if e != nil {
			return e
		}
		for _, line := range nsList.Items {
			name := line.Name
			nsOption := objectsUI.SelectOption{
				Id:       name,
				Text:     name,
				ParentID: selectedCluster,
			}
			namespaceOptions = append(namespaceOptions, nsOption)
		}
	} else {
		namespaceOption := objectsUI.SelectOption{
			Id:       "0",
			Text:     "===选择命名空间===",
			ParentID: "0",
		}
		namespaceOptions = append(namespaceOptions, namespaceOption)
	}
	nsSelect := objectsUI.SelectData{Title: "命名空间", SelectedId: selectedNamespace, Options: namespaceOptions}
	tplData["nsSelect"] = nsSelect

	return nil
}

func buildSelectData(tplData map[string]interface{}, dcList []interface{}, requestData map[string]string) error {

	selectedDC := strings.TrimSpace(requestData["dcID"])
	if selectedDC == "" {
		selectedDC = "0"
	}
	selectedCluster := strings.TrimSpace(requestData["clusterID"])
	if selectedCluster == "" {
		selectedCluster = "0"
	}

	// preparing datacenter options
	var dcOptions []objectsUI.SelectOption
	dcOption := objectsUI.SelectOption{
		Id:       "0",
		Text:     "===选择数据中心===",
		ParentID: "0",
	}
	dcOptions = append(dcOptions, dcOption)
	for _, line := range dcList {
		dcData, ok := line.(datacenter.DatacenterSchema)
		if !ok {
			return fmt.Errorf("the data is not datacenter schema")
		}
		dcOption := objectsUI.SelectOption{
			Id:       strconv.Itoa(int(dcData.Id)),
			Text:     dcData.CnName,
			ParentID: "0",
		}
		dcOptions = append(dcOptions, dcOption)
	}
	dcSelect := objectsUI.SelectData{Title: "数据中心", SelectedId: selectedDC, Options: dcOptions}
	tplData["dcSelect"] = dcSelect

	// preparing cluster options
	var clusterOptions []objectsUI.SelectOption
	clusterOption := objectsUI.SelectOption{
		Id:       "0",
		Text:     "===选择集群===",
		ParentID: "0",
	}
	clusterOptions = append(clusterOptions, clusterOption)
	if selectedDC != "0" {
		var k8sclusterEntity sysadmObjects.ObjectEntity
		k8sclusterEntity = sysadmK8sCluster.New()
		conditions := make(map[string]string, 0)
		var emptyString []string
		conditions["isDeleted"] = "='0'"
		conditions["dcid"] = "='" + selectedDC + "'"
		clusterList, e := k8sclusterEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, make(map[string]string))
		if e != nil {
			return e
		}

		for _, line := range clusterList {
			clusterData, ok := line.(sysadmK8sCluster.K8sclusterSchema)
			if !ok {
				return fmt.Errorf("the data is not cluster schema")
			}

			clusterOption := objectsUI.SelectOption{
				Id:       clusterData.Id,
				Text:     clusterData.CnName,
				ParentID: strconv.Itoa(int(clusterData.Dcid)),
			}
			clusterOptions = append(clusterOptions, clusterOption)
		}
	}
	clusterSelect := objectsUI.SelectData{Title: "集群", SelectedId: selectedCluster, Options: clusterOptions}
	tplData["clusterSelect"] = clusterSelect

	nsSelect := objectsUI.SelectData{}
	tplData["nsSelect"] = nsSelect

	return nil
}

func getRequestData(c *sysadmServer.Context, keys []string) (map[string]string, error) {
	requestData, e := utils.NewGetRequestData(c, keys)
	if e != nil {
		return requestData, e
	}

	objectIds := ""
	objectIDMap, _ := utils.GetRequestDataArray(c, []string{"objectid[]"})
	if objectIDMap != nil {
		objectIDSlice, ok := objectIDMap["objectid[]"]
		if ok {
			objectIds = strings.Join(objectIDSlice, ",")
		}
	}
	requestData["objectIds"] = objectIds
	if strings.TrimSpace(requestData["start"]) == "" {
		requestData["start"] = "0"
	}

	return requestData, nil
}

func newObjEntity(module string) (objectEntity, error) {
	for n, v := range modulesDefined {
		if n == module {
			return v, nil
		}
	}

	return nil, fmt.Errorf("module %s  has not found", module)
}

func buildClientSetByClusterID(clusterID string) (*kubernetes.Clientset, error) {
	clusterID = strings.TrimSpace(clusterID)
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID is empty")
	}

	var k8sclusterEntity sysadmObjects.ObjectEntity
	k8sclusterEntity = sysadmK8sCluster.New()
	clusterInfo, e := k8sclusterEntity.GetObjectInfoByID(clusterID)
	if e != nil {
		return nil, e
	}
	clusterData, ok := clusterInfo.(sysadmK8sCluster.K8sclusterSchema)
	if !ok {
		return nil, fmt.Errorf("the data is not K8S data schema")
	}
	ca := []byte(clusterData.Ca)
	cert := []byte(clusterData.Cert)
	key := []byte(clusterData.Key)
	restConf, e := sysadmK8sClient.BuildConfigFromParametes(ca, cert, key, clusterData.Apiserver, clusterData.Id, clusterData.ClusterUser, "default")
	if e != nil {
		return nil, e
	}

	clientSet, e := sysadmK8sClient.BuildClientset(restConf)
	if e != nil {
		return nil, e
	}

	return clientSet, nil
}

func checkUserLoginAndGetRequestData(c *sysadmServer.Context, keys []string) (map[string]string, error) {
	requestData, e := utils.NewGetRequestData(c, keys)
	if e != nil {
		return requestData, e
	}

	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		return requestData, fmt.Errorf("user has not login or not permission")
	}

	requestData["userid"] = userid.(string)

	return requestData, nil
}

func getFormDataFromMultipartForm(c *sysadmServer.Context, keys []string) (map[string]string, error) {
	data := make(map[string]string, 0)

	keys = append(keys, "dcid")
	keys = append(keys, "clusterID")
	keys = append(keys, "namespace")
	keys = append(keys, "addType")
	keys = append(keys, "objContent")

	formData, e := utils.GetMultipartData(c, keys)
	if e != nil {
		return data, e
	}

	for _, k := range keys {
		keyData := formData[k].([]string)
		value := ""
		if len(keyData) > 1 {
			for _, v := range keyData {
				if value == "" {
					value = v
				} else {
					value = value + "," + v
				}
			}
		} else {
			if len(keyData) > 0 {
				value = keyData[0]
			} else {
				value = ""
			}
		}
		data[k] = value
	}

	objContent := formData["objContent"].([]string)
	yamlContent := objContent[0]
	if data["addType"] == "1" {
		yamlFile, e := utils.GetMultipartData(c, []string{"objFile"})
		if e != nil {
			return data, e
		}
		yamlByte, e := utils.ReadUploadedFile(yamlFile["objFile"].(*multipart.FileHeader))
		if e != nil {
			return data, e
		}
		yamlContent = utils.Interface2String(yamlByte)
	}

	data["yamlContent"] = yamlContent

	return data, nil
}
