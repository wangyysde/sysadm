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
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"strings"
)

func GetNodeCount(restConf *rest.Config, role NodeRoleType) (ObjectCount, error) {
	ret := ObjectCount{Namespace: ""}

	if restConf == nil {
		return ret, fmt.Errorf("can not get nodes count on an empty client")
	}

	clientSet, e := BuildClientset(restConf)
	if e != nil {
		return ret, fmt.Errorf("build  client set error %s", e)
	}
	nodes, e := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return ret, e
	}

	total := 0
	ready := 0
	unready := 0

	for _, n := range nodes.Items {
		conditions := n.Status.Conditions
		statusReady := false
		for _, c := range conditions {
			if c.Type != corev1.NodeReady {
				continue
			}
			if c.Status == corev1.ConditionTrue {
				statusReady = true
			}
		}
		nodeRole := NodeRoleWK
		for l, _ := range n.Labels {
			if strings.TrimSpace(strings.ToLower(l)) == "node-role.kubernetes.io/control-plane" || strings.TrimSpace(strings.ToLower(l)) == "node-role.kubernetes.io/master" {
				nodeRole = NodeRoleCP
				break
			}
		}
		switch role {
		case NodeRoleALL:
			total = total + 1
			if statusReady {
				ready = ready + 1
			} else {
				unready = unready + 1
			}
		case NodeRoleCP:
			if nodeRole == NodeRoleCP {
				total = total + 1
				if statusReady {
					ready = ready + 1
				} else {
					unready = unready + 1
				}
			}
		case NodeRoleWK:
			if nodeRole == NodeRoleWK {
				total = total + 1
				if statusReady {
					ready = ready + 1
				} else {
					unready = unready + 1
				}
			}

		}
	}

	ret.Kind = NodeKind
	ret.Total = int32(total)
	ret.Ready = int32(ready)
	ret.Unready = int32(unready)

	return ret, nil
}

func GetNodeInfo(restConf *rest.Config, nodeName string) (*corev1.Node, error) {
	var ret *corev1.Node = nil
	if restConf == nil {
		return ret, fmt.Errorf("can not get nodes count on an empty client")
	}

	nodeName = strings.TrimSpace(nodeName)
	if nodeName == "" {
		return ret, fmt.Errorf("node name should not empty")
	}

	clientSet, e := BuildClientset(restConf)
	if e != nil {
		return ret, e
	}

	node, e := clientSet.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if e != nil {
		return ret, e
	}

	return node, nil
}

func GetNodeList(restConf *rest.Config) (*corev1.NodeList, error) {
	var ret *corev1.NodeList = nil

	if restConf == nil {
		return ret, fmt.Errorf("can not get nodes count on an empty client")
	}

	clientSet, e := BuildClientset(restConf)
	if e != nil {
		return ret, e
	}

	nodeList, e := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return ret, e
	}

	return nodeList, nil
}

func IsNodeStatusReady(node corev1.Node) bool {
	conditions := node.Status.Conditions
	for _, c := range conditions {
		if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
			return true
		}
	}

	return false
}
