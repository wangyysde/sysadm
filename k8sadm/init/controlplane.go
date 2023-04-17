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

package init

import (
	"fmt"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/controlplane"
)

func runControlPlanePhase(data *initData) error {
	cfg := data.Cfg()

	if err := controlplane.CreateStaticPodFiles(data.ManifestDir(), data.PatchesDir(), &cfg.ClusterConfiguration, &cfg.LocalAPIEndpoint, data.DryRun(), kubeadmconstants.KubeAPIServer); err != nil {
		return err
	}
	if err := controlplane.CreateStaticPodFiles(data.ManifestDir(), data.PatchesDir(), &cfg.ClusterConfiguration, &cfg.LocalAPIEndpoint, data.DryRun(), kubeadmconstants.KubeControllerManager); err != nil {
		return err
	}
	if err := controlplane.CreateStaticPodFiles(data.ManifestDir(), data.PatchesDir(), &cfg.ClusterConfiguration, &cfg.LocalAPIEndpoint, data.DryRun(), kubeadmconstants.KubeScheduler); err != nil {
		return err
	}
	fmt.Printf("[control-plane] Using manifest folder %q\n", data.ManifestDir())
	return nil
}
