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
	"github.com/wangyysde/sysadmServer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	sysadmK8sClient "sysadm/k8sclient"
	sysadmK8sCluster "sysadm/k8scluster/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/sysadmLog"
	"sysadm/sysadmapi/apiutils"
	"sysadm/utils"
)

func getclusteroptionfbydcHandler(c *sysadmServer.Context) {
	// order fields data of cluster list page
	var errs []sysadmLog.Sysadmerror
	requestData, e := utils.NewGetRequestData(c, []string{"objID"})
	if e != nil || requestData["objID"] == "" {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000200001, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(8000200001, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	var k8sclusterEntity sysadmObjects.ObjectEntity
	k8sclusterEntity = sysadmK8sCluster.New()
	conditions := make(map[string]string, 0)
	var emptyString []string
	conditions["isDeleted"] = "='0'"
	conditions["dcid"] = "='" + requestData["objID"] + "'"
	clusterList, e := k8sclusterEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, make(map[string]string))
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000200002, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(8000200002, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	msg := "0:===选择集群==="
	for _, line := range clusterList {
		lineClusterData := line.(sysadmK8sCluster.K8sclusterSchema)
		lineStr := lineClusterData.Id + ":" + lineClusterData.CnName
		msg = msg + "," + lineStr
	}

	response := apiutils.BuildResponseDataForSuccess(msg)
	c.JSON(http.StatusOK, response)

}

func getnsoptionbyclusterHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror

	requestData, e := utils.NewGetRequestData(c, []string{"objID"})
	if e != nil || requestData["objID"] == "" {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000200003, "error", "cluster id is empty or %s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(8000200003, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	var k8sclusterEntity sysadmObjects.ObjectEntity
	k8sclusterEntity = sysadmK8sCluster.New()
	clusterInfo, e := k8sclusterEntity.GetObjectInfoByID(requestData["objID"])
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000200004, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(8000200004, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}
	clusterData, ok := clusterInfo.(sysadmK8sCluster.K8sclusterSchema)
	if !ok {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000200005, "error", "the data is not k8scluster schema"))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(8000200005, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}
	ca := []byte(clusterData.Ca)
	cert := []byte(clusterData.Cert)
	key := []byte(clusterData.Key)
	restConf, e := sysadmK8sClient.BuildConfigFromParametes(ca, cert, key, clusterData.Apiserver, clusterData.Id, clusterData.ClusterUser, "default")
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000200006, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(8000200006, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	clientSet, e := sysadmK8sClient.BuildClientset(restConf)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000200007, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(8000200007, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}
	nsList, e := clientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000200008, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(8000200008, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	msg := "0:所有命名空间"
	for _, line := range nsList.Items {
		lineStr := line.Name + ":" + line.Name
		msg = msg + "," + lineStr
	}

	response := apiutils.BuildResponseDataForSuccess(msg)
	c.JSON(http.StatusOK, response)
}
