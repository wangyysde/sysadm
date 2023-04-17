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
	clientset "k8s.io/client-go/kubernetes"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	certsphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/controlplane"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/copycerts"
	kubeconfigphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/kubeconfig"
	kubeconfigutil "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"
)

func controlPlanePreparePhase(data *joinData) error {
	if err := runControlPlanePrepareDownloadCertsPhaseLocal(data); err != nil {
		return err
	}
	if err := runControlPlanePrepareCertsPhaseLocal(data); err != nil {
		return err
	}
	if err := runControlPlanePrepareKubeconfigPhaseLocal(data); err != nil {
		return err
	}
	if err := runControlPlanePrepareControlPlaneSubphase(data); err != nil {
		return err
	}
	return nil
}

func runControlPlanePrepareDownloadCertsPhaseLocal(data *joinData) error {
	if data.Cfg().ControlPlane == nil || len(data.CertificateKey()) == 0 {
		fmt.Println("[download-certs] Skipping certs download")
		return nil
	}

	cfg, err := data.InitCfg()
	if err != nil {
		return err
	}

	// If we're dry-running, download certs to tmp dir, and defer to restore to the path originally specified by the user
	certsDir := cfg.CertificatesDir
	cfg.CertificatesDir = data.CertificateWriteDir()
	defer func() { cfg.CertificatesDir = certsDir }()

	client, err := bootstrapClient(data)
	if err != nil {
		return err
	}

	if err := copycerts.DownloadCerts(client, cfg, data.CertificateKey()); err != nil {
		return errors.Wrap(err, "error downloading certs")
	}
	return nil
}

func runControlPlanePrepareCertsPhaseLocal(data *joinData) error {
	// Skip if this is not a control plane
	if data.Cfg().ControlPlane == nil {
		return nil
	}

	cfg, err := data.InitCfg()
	if err != nil {
		return err
	}

	fmt.Printf("[certs] Using certificateDir folder %q\n", cfg.CertificatesDir)

	certsDir := cfg.CertificatesDir
	cfg.CertificatesDir = data.CertificateWriteDir()
	defer func() { cfg.CertificatesDir = certsDir }()
	return certsphase.CreatePKIAssets(cfg)
}

func runControlPlanePrepareKubeconfigPhaseLocal(data *joinData) error {
	// Skip if this is not a control plane
	if data.Cfg().ControlPlane == nil {
		return nil
	}

	cfg, err := data.InitCfg()
	if err != nil {
		return err
	}

	fmt.Println("[kubeconfig] Generating kubeconfig files")
	fmt.Printf("[kubeconfig] Using kubeconfig folder %q\n", data.KubeConfigDir())

	if err := kubeconfigphase.CreateJoinControlPlaneKubeConfigFiles(data.KubeConfigDir(), cfg); err != nil {
		return errors.Wrap(err, "error generating kubeconfig files")
	}

	return nil
}

func runControlPlanePrepareControlPlaneSubphase(data *joinData) error {
	// Skip if this is not a control plane
	if data.Cfg().ControlPlane == nil {
		return nil
	}

	cfg, err := data.InitCfg()
	if err != nil {
		return err
	}

	fmt.Printf("[control-plane] Using manifest folder %q\n", data.ManifestDir())

	for _, component := range kubeadmconstants.ControlPlaneComponents {
		fmt.Printf("[control-plane] Creating static Pod manifest for %q\n", component)
		err := controlplane.CreateStaticPodFiles(
			data.ManifestDir(),
			data.PatchesDir(),
			&cfg.ClusterConfiguration,
			&cfg.LocalAPIEndpoint,
			data.DryRun(),
			component,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func bootstrapClient(data *joinData) (clientset.Interface, error) {
	tlsBootstrapCfg, err := data.TLSBootstrapCfg()
	if err != nil {
		return nil, errors.Wrap(err, "unable to access the cluster")
	}
	client, err := kubeconfigutil.ToClientSet(tlsBootstrapCfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to access the cluster")
	}
	return client, nil
}
