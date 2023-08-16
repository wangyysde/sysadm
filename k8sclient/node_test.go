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
	"testing"
)

func TestGetNodeCount(t *testing.T) {
	if restConf == nil {
		id, apiserver, clusterUser, ca, cert, key, e := getClusterData()
		if e != nil {
			t.Fatalf("get cluster data error:%s\n", e)
		}

		tmpRestConf, e := BuildConfigFromParametes([]byte(ca), []byte(cert), []byte(key), apiserver, id, clusterUser, "default")
		if e != nil {
			t.Fatalf("build rest config error:%s\n", e)
		}
		restConf = tmpRestConf
	}

	cpCount, e := GetNodeCount(restConf, NodeRoleCP)
	if e != nil {
		t.Fatalf("%s", e)
	}

	wkCount, e := GetNodeCount(restConf, NodeRoleWK)
	if e != nil {
		t.Fatalf("%s", e)
	}

	fmt.Printf("control plane=Healthy: %d Unhealthy: %d  ++++++ work=Healthy:%d  Unhealthy: %d\n",
		cpCount.Ready, cpCount.Unready, wkCount.Ready, wkCount.Unready)
}

func TestGetNodeInfo(t *testing.T) {
	if restConf == nil {
		id, apiserver, clusterUser, ca, cert, key, e := getClusterData()
		if e != nil {
			t.Fatalf("get cluster data error:%s\n", e)
		}

		tmpRestConf, e := BuildConfigFromParametes([]byte(ca), []byte(cert), []byte(key), apiserver, id, clusterUser, "default")
		if e != nil {
			t.Fatalf("build rest config error:%s\n", e)
		}
		restConf = tmpRestConf
	}

	node, e := GetNodeInfo(restConf, "cp1")
	if e != nil {
		t.Fatalf("there is an error occurred when get node infromation: %s\n", e)
	}

	fmt.Printf("node info:%#v\n", node.Status.NodeInfo)
	fmt.Printf("pod cidr:%s\n", node.Spec.PodCIDR)
}
