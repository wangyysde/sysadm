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

func TestGetKubernetesVersion(t *testing.T) {
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

	version, e := GetKubernetesVersion(restConf)
	if e != nil {
		t.Fatalf("%s", e)
	}

	fmt.Printf("kubernetes version:%s \n", version)
}

func TestGetPlatformInfo(t *testing.T) {
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

	info, e := GetPlatformInfo(restConf)
	if e != nil {
		t.Fatalf("%s", e)
	}

	fmt.Printf("platform is:%s \n", info)
}

func TestGetPodCIDR(t *testing.T) {
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

	podCIDR, e := GetPodCIDR(restConf)
	if e != nil {
		t.Fatalf("%s\n", e)
	}

	fmt.Printf("podCIDR: %s\n", podCIDR)
}

func TestGetSvcCIDR(t *testing.T) {
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

	svcCIDR, e := GetSvcCIDR(restConf)
	if e != nil {
		t.Fatalf("%s\n", e)
	}

	fmt.Printf("svcCIDR: %s\n", svcCIDR)
}
