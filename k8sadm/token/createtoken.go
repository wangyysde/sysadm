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

package token

import (
	"fmt"
	clientset "k8s.io/client-go/kubernetes"
	kubeadmscheme "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
	kubeadmapiv1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta3"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/options"
	cmdutil "k8s.io/kubernetes/cmd/kubeadm/app/cmd/util"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	tokenphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/bootstraptoken/node"
	configutil "k8s.io/kubernetes/cmd/kubeadm/app/util/config"
)

func CreateToken(criSocket string, certificateKey string) error {
	kubeConfigFile := constants.GetAdminKubeConfigPath()
	cfg := cmdutil.DefaultInitConfiguration()
	kubeadmscheme.Scheme.Default(cfg)
	//将cri重写为用户指定的
	cfg.NodeRegistration.CRISocket = criSocket

	bto := options.NewBootstrapTokenOptions()
	if err := bto.ApplyTo(cfg); err != nil {
		return err
	}
	fmt.Println("[token] getting Clientsets from kubeconfig file")
	kubeConfigFile = cmdutil.GetKubeConfigPath(kubeConfigFile)
	client, err := cmdutil.GetClientSet(kubeConfigFile, false)
	if err != nil {
		return err
	}
	return runCreateToken(client, cfg, certificateKey, kubeConfigFile)
}

func runCreateToken(client clientset.Interface, initCfg *kubeadmapiv1.InitConfiguration, certificateKey string, kubeConfigFile string) error {
	// ClusterConfiguration is needed just for the call to LoadOrDefaultInitConfiguration
	clusterCfg := &kubeadmapiv1.ClusterConfiguration{
		// KubernetesVersion is not used, but we set this explicitly to avoid
		// the lookup of the version from the internet when executing LoadOrDefaultInitConfiguration
		KubernetesVersion: kubeadmconstants.CurrentKubernetesVersion.String(),
	}
	kubeadmscheme.Scheme.Default(clusterCfg)

	// This call returns the ready-to-use configuration based on the configuration file that might or might not exist and the default cfg populated by flags
	fmt.Println("[token] loading configurations")

	internalcfg, err := configutil.LoadOrDefaultInitConfiguration("", initCfg, clusterCfg)
	if err != nil {
		return err
	}

	fmt.Println("[token] creating token")
	if err := tokenphase.CreateNewTokens(client, internalcfg.BootstrapTokens); err != nil {
		return err
	}

	return nil
}
