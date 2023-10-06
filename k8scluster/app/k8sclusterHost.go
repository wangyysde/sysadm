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
	"k8s.io/client-go/rest"

	"strings"
	"sysadm/k8sclient"

	hostApp "sysadm/host/app"
	osApp "sysadm/os/app"
	versionApp "sysadm/version/app"
)

func tryReconizeHostsInCluster(clusterData K8sclusterSchema, restConf *rest.Config, userid int) error {
	if restConf == nil {
		clusterUser := strings.TrimSpace(clusterData.ClusterUser)
		apiserver := strings.TrimSpace(clusterData.Apiserver)
		ca := strings.TrimSpace(clusterData.Ca)
		cert := strings.TrimSpace(clusterData.Cert)
		key := strings.TrimSpace(clusterData.Key)
		clusterID := clusterData.Id
		tmpRestConf, e := k8sclient.BuildConfigFromParametes([]byte(ca), []byte(cert), []byte(key), apiserver, clusterID, clusterUser, "default")
		if e != nil {
			return e
		}
		restConf = tmpRestConf
	}

	nodes, e := k8sclient.GetNodeList(restConf)

	if e != nil {
		return e
	}

	host, e := hostApp.New(runData.dbConf, runData.workingRoot)
	if e != nil {
		return e
	}
	var emptyString []string
	systemidCondition := make(map[string]string, 0)
	for _, node := range nodes.Items {
		systemID := node.Status.NodeInfo.SystemUUID
		if id := strings.TrimSpace(systemID); id != "" {
			systemidCondition["systemID"] = `="` + id + `"`
			hostList, e := host.GetObjectList("", emptyString, emptyString, systemidCondition, 0, 0, nil)
			if e != nil {
				continue
			}

			if len(hostList) < 1 {
				e = tryAddHost(clusterData, node, userid)
				if e != nil {
					return e
				}
				continue
			}

			if len(hostList) > 1 {
				//TODO
				// 说明系统中的数据不正确，发出告警，请系统管理员进行处理
				continue
			}

			// 检查并处理主机的信息
		}

	}

	return nil
}

func tryAddHost(clusterData K8sclusterSchema, node corev1.Node, userid int) error {
	hostname := ""
	ips := make([]string, 0)
	for _, addr := range node.Status.Addresses {
		switch addr.Type {
		case corev1.NodeHostName:
			hostname = addr.Address
		case corev1.NodeInternalIP, corev1.NodeExternalIP:
			ips = append(ips, addr.Address)
		}
	}
	osImage := node.Status.NodeInfo.OSImage
	osImages := strings.Split(osImage, " ")
	osVersion, osName := "", ""
	if len(osImages) > 2 {
		osName = osImages[0]
		osVersion = osImages[2]
	}
	osid := 0
	osversionid := 0
	if osName != "" {
		osInst := osApp.New()
		osInfo, e := osInst.GetObjectInfoByName(osName)
		if e != nil {
			return e
		}
		osData := osInfo.(osApp.OSSchema)
		osid = osData.OSID
	}

	if osVersion != "" {
		verInst := versionApp.New()
		verInfo, e := verInst.GetObjectInfoByName(osVersion)
		if e != nil {
			return e
		}
		verData := verInfo.(versionApp.VersionSchema)
		osversionid = verData.VersionID
	}

	status := hostApp.HostStatusUnkown
	k8sclusterid := clusterData.Id
	dcid := clusterData.Dcid
	azid := clusterData.Azid
	machineID := node.Status.NodeInfo.MachineID
	systemID := node.Status.NodeInfo.SystemUUID
	architecture := node.Status.NodeInfo.Architecture
	kernelVersion := node.Status.NodeInfo.KernelVersion

	return hostApp.AddHostFromCluster(userid, osid, osversionid, int(dcid), int(azid), hostname, string(status),
		k8sclusterid, machineID, systemID, architecture, kernelVersion, ips)

}
