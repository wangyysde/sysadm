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
		connectType := strings.TrimSpace(clusterData.ConnectType)
		token := strings.TrimSpace(clusterData.Token)
		kubeConfig := strings.TrimSpace(clusterData.KubeConfig)
		clusterID := clusterData.Id
		tmpRestConf, e := k8sclient.BuildConfigFromParasWithConnectType(connectType, apiserver, clusterID, clusterUser, "", ca, cert, key, token, kubeConfig)
		if e != nil {
			return e
		}
		restConf = tmpRestConf
	}

	nodes, e := k8sclient.GetNodeList(restConf)

	if e != nil {
		return e
	}

	for _, node := range nodes.Items {
		e = tryAddHost(clusterData, node, userid)
		if e != nil {
			return e
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

	status := ""
	for _, c := range node.Status.Conditions {
		if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
			status = string(HostStatusReady)
			break
		}
		if c.Status == corev1.ConditionTrue {
			if status == "" {
				status = string(c.Type)
			} else {
				status = status + "," + string(c.Type)
			}
		}
	}
	if status == "" {
		status = string(HostStatusUnkown)
	}
	k8sclusterid := clusterData.Id
	dcid := clusterData.Dcid
	azid := clusterData.Azid
	machineID := node.Status.NodeInfo.MachineID
	systemID := node.Status.NodeInfo.SystemUUID
	architecture := node.Status.NodeInfo.Architecture
	kernelVersion := node.Status.NodeInfo.KernelVersion

	hostInst, e := HostNew(runData.dbConf, runData.workingRoot)
	if e != nil {
		return e
	}
	return hostInst.AddHostFromCluster(userid, osid, osversionid, int(dcid), int(azid), hostname, string(status),
		k8sclusterid, machineID, systemID, architecture, kernelVersion, ips)

}

/*
func tryUpdateHostInfoForClusterAdd(clusterData K8sclusterSchema, node corev1.Node, userid int,
	hostList []interface{}) error {

	hostLine := hostList[0]
	hostData, ok := hostLine.(hostApp.HostSchema)
	if !ok {
		return fmt.Errorf("the data is not hostSchema data")
	}

	// 如果对应的主机存在关联的集群，


	hostID := hostData.HostId
	hostName := node.Name
	status := ""
	for _, c := range node.Status.Conditions {
		if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
			status = string(hostApp.HostStatusReady)
			break
		}
		if c.Status == corev1.ConditionTrue {
			if status == "" {
				status = string(c.Type)
			} else {
				status = status + "," + string(c.Type)
			}
		}
	}
	if status == "" {
		status = string(hostApp.HostStatusUnkown)
	}

	k8sclusterid := clusterData.Id
	dcid := clusterData.Dcid
	azid := clusterData.Azid
	machineID := node.Status.NodeInfo.MachineID
	systemID := node.Status.NodeInfo.SystemUUID
	architecture := node.Status.NodeInfo.Architecture
	kernelVersion := node.Status.NodeInfo.KernelVersion
	ips := make([]string, 0)
	for _, addr := range node.Status.Addresses {
		ips = append(ips, addr.Address)
	}

	return hostApp.UpdateHostInfoForClusterAdd(hostID, hostName, status, k8sclusterid, machineID, systemID, architecture, kernelVersion, dcid, azid, userid)
}
*/
