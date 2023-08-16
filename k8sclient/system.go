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

package k8sclient

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"strings"
)

func GetKubernetesVersion(restConf *rest.Config) (string, error) {
	if restConf == nil {
		return "", fmt.Errorf("can not get nodes count on an empty client")
	}

	ds := discovery.NewDiscoveryClientForConfigOrDie(restConf)
	version, e := ds.ServerVersion()
	if e != nil {
		return "", fmt.Errorf("discovering infomation error %s", e)
	}

	return version.GitVersion, nil
}

func GetPlatformInfo(restConf *rest.Config) (string, error) {
	if restConf == nil {
		return "", fmt.Errorf("can not get nodes count on an empty client")
	}

	ds := discovery.NewDiscoveryClientForConfigOrDie(restConf)
	version, e := ds.ServerVersion()
	if e != nil {
		return "", fmt.Errorf("discovering infomation error %s", e)
	}

	return version.Platform, nil
}

func GetPodCIDR(restConf *rest.Config) (string, error) {
	if restConf == nil {
		return "", fmt.Errorf("can not get pod CIDR on an empty client")
	}

	nodeList, e := GetNodeList(restConf)
	if e != nil {
		return "", e
	}

	node := nodeList.Items[0]

	return node.Spec.PodCIDR, nil
}

func GetSvcCIDR(restConf *rest.Config) (string, error) {
	if restConf == nil {
		return "", fmt.Errorf("can not get pod CIDR on an empty client")
	}

	pods, e := GetPodInfoWithPrefix(restConf, "kube-system", "kube-controller-manager-")
	if e != nil {
		return "", e
	}

	pod := pods[0]
	svcCIDR := ""
	podContainers := pod.Spec.Containers
	for _, c := range podContainers {
		if c.Name == "kube-controller-manager" {
			commands := c.Command
			for _, cmd := range commands {
				if strings.Contains(cmd, "service-cluster-ip-range") {
					lineArray := strings.Split(cmd, "=")
					if len(lineArray) > 1 {
						svcCIDR = lineArray[1]
						break
					}
				}
			}
		}
		if svcCIDR != "" {
			break
		}
	}

	return svcCIDR, nil
}

func GetCRIInfo(restConf *rest.Config) (string, error) {
	if restConf == nil {
		return "", fmt.Errorf("we can not get CRI information on an empty client")
	}

	nodeList, e := GetNodeList(restConf)
	if e != nil {
		return "", e
	}

	var node corev1.Node
	found := false
	for _, n := range nodeList.Items {
		if IsNodeStatusReady(n) {
			node = n
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("system can not get cri information for no node is ready")
	}

	criInfo := node.Status.NodeInfo.ContainerRuntimeVersion

	return criInfo, nil
}
