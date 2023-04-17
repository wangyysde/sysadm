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

package join

import (
	"fmt"
	"github.com/pkg/errors"
	etcdphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/etcd"
	markcontrolplanephase "k8s.io/kubernetes/cmd/kubeadm/app/phases/markcontrolplane"
	etcdutil "k8s.io/kubernetes/cmd/kubeadm/app/util/etcd"
)

func controlPlaneJoinPhase(data *joinData) error {
	if err := runEtcdPhase(data); err != nil {
		return err
	}
	if err := runUpdateStatusPhase(data); err != nil {
		return err
	}
	if err := runMarkControlPlanePhase(data); err != nil {
		return err
	}
	return nil
}

func runEtcdPhase(data *joinData) error {
	if data.Cfg().ControlPlane == nil {
		return nil
	}

	// gets access to the cluster using the identity defined in admin.conf
	client, err := data.Client()
	if err != nil {
		return errors.Wrap(err, "couldn't create Kubernetes client")
	}
	cfg, err := data.InitCfg()
	if err != nil {
		return err
	}
	// in case of local etcd
	if cfg.Etcd.External != nil {
		fmt.Println("[control-plane-join] Using external etcd - no local stacked instance added")
		return nil
	}

	if err := etcdutil.CreateDataDirectory(cfg.Etcd.Local.DataDir); err != nil {
		return err
	}

	if err := etcdphase.CreateStackedEtcdStaticPodManifestFile(client, data.ManifestDir(), data.PatchesDir(), cfg.NodeRegistration.Name, &cfg.ClusterConfiguration, &cfg.LocalAPIEndpoint, data.DryRun(), data.CertificateWriteDir()); err != nil {
		return errors.Wrap(err, "error creating local etcd static pod manifest file")
	}

	return nil
}

func runUpdateStatusPhase(data *joinData) error {
	if data.Cfg().ControlPlane != nil {
		fmt.Println("The 'update-status' phase is deprecated and will be removed in a future release. " +
			"Currently it performs no operation")
	}
	return nil
}

func runMarkControlPlanePhase(data *joinData) error {
	if data.Cfg().ControlPlane == nil {
		return nil
	}

	// gets access to the cluster using the identity defined in admin.conf
	client, err := data.Client()
	if err != nil {
		return errors.Wrap(err, "couldn't create Kubernetes client")
	}
	cfg, err := data.InitCfg()
	if err != nil {
		return err
	}

	if err := markcontrolplanephase.MarkControlPlane(client, cfg.NodeRegistration.Name, cfg.NodeRegistration.Taints); err != nil {
		return errors.Wrap(err, "error applying control-plane label and taints")
	}

	return nil
}
