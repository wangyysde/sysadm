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
	"github.com/pkg/errors"
	etcdphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/etcd"
	etcdutil "k8s.io/kubernetes/cmd/kubeadm/app/util/etcd"
)

func runEtcd(data *initData) error {
	cfg := data.Cfg()

	if err := etcdutil.CreateDataDirectory(cfg.Etcd.Local.DataDir); err != nil {
		return err
	}
	fmt.Printf("[etcd] Creating static Pod manifest for local etcd in %q\n", data.ManifestDir())
	if err := etcdphase.CreateLocalEtcdStaticPodManifestFile(data.ManifestDir(), data.PatchesDir(), cfg.NodeRegistration.Name, &cfg.ClusterConfiguration, &cfg.LocalAPIEndpoint, data.DryRun()); err != nil {
		return errors.Wrap(err, "error creating local etcd static pod manifest file")
	}
	return nil
}
