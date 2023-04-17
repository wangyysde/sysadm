/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* https://www.sysadm.cn/licenses/apache-2.0.txt
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
*
 */

package main

import (
	"os"

	kubeadmInit "github.com/wangyysde/sysadm/k8sadm/init"
	"github.com/wangyysde/sysadm/k8sadm/join"
	"github.com/wangyysde/sysadm/k8sadm/token"
)

func main() {
	var err error
	args := os.Args
	switch args[1] {
	case "init":
		err = kubeadmInit.StartInit(kubeadmInit.InitParam{
			AdvertiseAddress:  "172.28.2.10",
			BindPort:          6443,
			CRISocket:         "unix:///run/containerd/containerd.sock",
			ServiceSubnet:     "10.1.0.0/16",
			PodSubnet:         "10.2.0.0/16",
			ImageRepository:   "hb.sincerecloud.com/k8s/v1.26.3",
			KubernetesVersion: "1.26.3",
			UploadCerts:       true,
			Out:               os.Stdout,
		})
	case "join":
		err = join.StartJoin(join.JoinParam{
			APIServerAdvertiseAddress: "172.28.2.10",
			APIServerBindPort:         6443,
			TokenStr:                  "hk21p8.lvuu2pg2ztcefyvp",
			TokenDiscoveryCAHash:      "sha256:882b06787a51f6e7403e5b337317fbcebadc895913cfafb087b47f468718294a",
			NodeCRISocket:             "unix:///run/containerd/containerd.sock",
			Out:                       os.Stdout,
		})
	case "createToken":
		err = token.CreateToken("unix://run/containerd/containerd.sock", "")
	case "tokenList":
		err = token.GetToken()
	}

	if err != nil {
		return
	}
}
