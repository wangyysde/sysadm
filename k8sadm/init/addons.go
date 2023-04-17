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
*
*/

package init

import (
	"io"
	clientset "k8s.io/client-go/kubernetes"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	dnsaddon "k8s.io/kubernetes/cmd/kubeadm/app/phases/addons/dns"
	proxyaddon "k8s.io/kubernetes/cmd/kubeadm/app/phases/addons/proxy"
)

func runAddon(data *initData) error {
	if err := runCoreDNSAddon(data); err != nil {
		return err
	}
	if err := runKubeProxyAddon(data); err != nil {
		return err
	}
	return nil
}

func getInitData(data *initData) (*kubeadmapi.InitConfiguration, clientset.Interface, io.Writer, error) {
	cfg := data.Cfg()
	var client clientset.Interface
	var err error
	client, err = data.Client()
	if err != nil {
		return nil, nil, nil, err
	}

	out := data.OutputWriter()
	return cfg, client, out, err
}

// runCoreDNSAddon installs CoreDNS addon to a Kubernetes cluster
func runCoreDNSAddon(data *initData) error {
	cfg, client, out, err := getInitData(data)
	if err != nil {
		return err
	}
	return dnsaddon.EnsureDNSAddon(&cfg.ClusterConfiguration, client, out, false)
}

// runKubeProxyAddon installs KubeProxy addon to a Kubernetes cluster
func runKubeProxyAddon(data *initData) error {
	cfg, client, out, err := getInitData(data)
	if err != nil {
		return err
	}
	return proxyaddon.EnsureProxyAddon(&cfg.ClusterConfiguration, &cfg.LocalAPIEndpoint, client, out, false)
}
