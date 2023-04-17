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
	"io"
	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	kubeproxyconfig "k8s.io/kube-proxy/config/v1alpha1"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmscheme "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
	kubeadmapiv1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta3"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/options"
	phases "k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/init"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
	configutil "k8s.io/kubernetes/cmd/kubeadm/app/util/config"
	kubeconfigutil "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"
	"os"
	"path/filepath"
)

type InitParam struct {
	AdvertiseAddress  string
	BindPort          int32
	CRISocket         string
	ServiceSubnet     string
	PodSubnet         string
	ImageRepository   string
	KubernetesVersion string
	UploadCerts       bool
	Out               io.Writer
}

// compile-time assert that the local data object satisfies the phases data interface.
var _ phases.InitData = &initData{}

// initData defines all the runtime information used when running the kubeadm init workflow;
// this data is shared across all the phases that are included in the workflow.
type initData struct {
	cfg                     *kubeadmapi.InitConfiguration
	skipTokenPrint          bool
	dryRun                  bool
	kubeconfigDir           string
	kubeconfigPath          string
	ignorePreflightErrors   sets.String
	certificatesDir         string
	dryRunDir               string
	externalCA              bool
	client                  clientset.Interface
	outputWriter            io.Writer
	uploadCerts             bool
	skipCertificateKeyPrint bool
	patchesDir              string
}

func StartInit(initParam InitParam) error {
	data, err := newInitData(initParam, initParam.Out)
	if err != nil {
		return err
	}
	//预检测
	if err = runPreflight(data); err != nil {
		fmt.Println("======== runPreflight", err)
		return err
	}
	//创建证书
	if err = runCerts(data); err != nil {
		fmt.Println("======== runCerts", err)
		return err
	}
	//创建配置文件
	if err = runKubeConfig(data); err != nil {
		fmt.Println("======== runKubeConfig", err)
		return err
	}
	//生成kubelet环境变量文件，启动kubelet
	if err = runKubeletStart(data); err != nil {
		fmt.Println("======== runKubeletStart", err)
		return err
	}
	//生成静态pod文件
	if err = runControlPlanePhase(data); err != nil {
		fmt.Println("======== runControlPlanePhase", err)
		return err
	}
	//创建etcd静态pod文件
	if err = runEtcd(data); err != nil {
		fmt.Println("======== runEtcd", err)
		return err
	}
	//初始化控制平面
	if err = runWaitControlPlane(data); err != nil {
		fmt.Println("======== runWaitControlPlane", err)
		return err
	}
	//上传配置文件
	if err = runUploadConfig(data); err != nil {
		fmt.Println("======== runUploadConfig", err)
		return err
	}
	//上传证书
	if err = runUploadCerts(data); err != nil {
		fmt.Println("======== runUploadCerts", err)
		return err
	}
	//执行控制平面检查
	if err = runMarkControlPlane(data); err != nil {
		fmt.Println("======== runMarkControlPlane", err)
		return err
	}
	//生成token
	if err = runBootstrapToken(data); err != nil {
		fmt.Println("======== runBootstrapToken", err)
		return err
	}
	//检测是否启用了kubelet证书轮换，并更新kubelet.conf文件以指向Node用户的可轮换证书和密钥
	if err = runKubeletFinalizeCertRotation(data); err != nil {
		fmt.Println("======== runKubeletFinalizeCertRotation", err)
		return err
	}
	//加载coreDNS和kube-proxy插件
	if err = runAddon(data); err != nil {
		fmt.Println("======== runAddon", err)
		return err
	}
	//打印加入集群的命令
	if err = showJoinCommand(data); err != nil {
		fmt.Println("======== showJoinCommand", err)
		return err
	}
	return nil
}

func SetInitConfig(cfg *kubeadmapiv1.InitConfiguration, initParam InitParam) {
	cfg.LocalAPIEndpoint.AdvertiseAddress = initParam.AdvertiseAddress
	cfg.LocalAPIEndpoint.BindPort = initParam.BindPort
	cfg.NodeRegistration.CRISocket = initParam.CRISocket
}

func SetClusterConfigFlags(cfg *kubeadmapiv1.ClusterConfiguration, initParam InitParam) {
	cfg.Networking.ServiceSubnet = initParam.ServiceSubnet
	cfg.Networking.PodSubnet = initParam.PodSubnet
	cfg.ImageRepository = initParam.ImageRepository
	cfg.KubernetesVersion = initParam.KubernetesVersion
}

func newInitData(initParam InitParam, out io.Writer) (*initData, error) {
	// initialize the public kubeadm config API by applying defaults
	externalInitCfg := &kubeadmapiv1.InitConfiguration{}
	kubeadmscheme.Scheme.Default(externalInitCfg)
	externalClusterCfg := &kubeadmapiv1.ClusterConfiguration{}
	kubeadmscheme.Scheme.Default(externalClusterCfg)
	// Create the options object for the bootstrap token-related flags, and override the default value for .Description
	bto := options.NewBootstrapTokenOptions()
	bto.Description = "The default bootstrap token generated by 'kubeadm init'."
	SetInitConfig(externalInitCfg, initParam)
	SetClusterConfigFlags(externalClusterCfg, initParam)

	if err := bto.ApplyTo(externalInitCfg); err != nil {
		return nil, err
	}

	// Either use the config file if specified, or convert public kubeadm API to the internal InitConfiguration
	// and validates InitConfiguration
	cfg, err := configutil.LoadOrDefaultInitConfiguration("", externalInitCfg, externalClusterCfg)
	if err != nil {
		return nil, err
	}

	//todo 借鉴ip地址的校验规则和版本校验规则
	//if err := configutil.VerifyAPIServerBindAddress(cfg.LocalAPIEndpoint.AdvertiseAddress); err != nil {
	//	return nil, err
	//}
	//if err := features.ValidateVersion(features.InitFeatureGates, cfg.FeatureGates, cfg.KubernetesVersion); err != nil {
	//	return nil, err
	//}

	//启用ipvs
	proxyCfg := cfg.ClusterConfiguration.ComponentConfigs["kubeproxy.config.k8s.io"].Get().(*kubeproxyconfig.KubeProxyConfiguration)
	proxyCfg.Mode = "ipvs"
	cfg.ClusterConfiguration.ComponentConfigs["kubeproxy.config.k8s.io"].Set(proxyCfg)
	return &initData{
		cfg:                   cfg,
		certificatesDir:       cfg.CertificatesDir,
		kubeconfigDir:         kubeadmconstants.KubernetesDir,
		kubeconfigPath:        kubeadmconstants.GetAdminKubeConfigPath(),
		ignorePreflightErrors: sets.NewString(),
		outputWriter:          out,
		uploadCerts:           initParam.UploadCerts,
	}, nil
}

// UploadCerts returns Uploadcerts flag.
func (d *initData) UploadCerts() bool {
	return d.uploadCerts
}

// CertificateKey returns the key used to encrypt the certs.
func (d *initData) CertificateKey() string {
	return d.cfg.CertificateKey
}

// SetCertificateKey set the key used to encrypt the certs.
func (d *initData) SetCertificateKey(key string) {
	d.cfg.CertificateKey = key
}

// SkipCertificateKeyPrint returns the skipCertificateKeyPrint flag.
func (d *initData) SkipCertificateKeyPrint() bool {
	return d.skipCertificateKeyPrint
}

// Cfg returns initConfiguration.
func (d *initData) Cfg() *kubeadmapi.InitConfiguration {
	return d.cfg
}

// DryRun returns the DryRun flag.
func (d *initData) DryRun() bool {
	return d.dryRun
}

// SkipTokenPrint returns the SkipTokenPrint flag.
func (d *initData) SkipTokenPrint() bool {
	return d.skipTokenPrint
}

// IgnorePreflightErrors returns the IgnorePreflightErrors flag.
func (d *initData) IgnorePreflightErrors() sets.String {
	return d.ignorePreflightErrors
}

// CertificateWriteDir returns the path to the certificate folder or the temporary folder path in case of DryRun.
func (d *initData) CertificateWriteDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return d.certificatesDir
}

// CertificateDir returns the CertificateDir as originally specified by the user.
func (d *initData) CertificateDir() string {
	return d.certificatesDir
}

// KubeConfigDir returns the path of the Kubernetes configuration folder or the temporary folder path in case of DryRun.
func (d *initData) KubeConfigDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return d.kubeconfigDir
}

// KubeConfigPath returns the path to the kubeconfig file to use for connecting to Kubernetes
func (d *initData) KubeConfigPath() string {
	if d.dryRun {
		d.kubeconfigPath = filepath.Join(d.dryRunDir, kubeadmconstants.AdminKubeConfigFileName)
	}
	return d.kubeconfigPath
}

// ManifestDir returns the path where manifest should be stored or the temporary folder path in case of DryRun.
func (d *initData) ManifestDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return kubeadmconstants.GetStaticPodDirectory()
}

// KubeletDir returns path of the kubelet configuration folder or the temporary folder in case of DryRun.
func (d *initData) KubeletDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return kubeadmconstants.KubeletRunDirectory
}

// ExternalCA returns true if an external CA is provided by the user.
func (d *initData) ExternalCA() bool {
	return d.externalCA
}

// OutputWriter returns the io.Writer used to write output to by this command.
func (d *initData) OutputWriter() io.Writer {
	return d.outputWriter
}

// Client returns a Kubernetes client to be used by kubeadm.
// This function is implemented as a singleton, thus avoiding to recreate the client when it is used by different phases.
// Important. This function must be called after the admin.conf kubeconfig file is created.
func (d *initData) Client() (clientset.Interface, error) {
	if d.client == nil {
		if d.dryRun {
			svcSubnetCIDR, err := kubeadmconstants.GetKubernetesServiceCIDR(d.cfg.Networking.ServiceSubnet)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to get internal Kubernetes Service IP from the given service CIDR (%s)", d.cfg.Networking.ServiceSubnet)
			}
			// If we're dry-running, we should create a faked client that answers some GETs in order to be able to do the full init flow and just logs the rest of requests
			dryRunGetter := apiclient.NewInitDryRunGetter(d.cfg.NodeRegistration.Name, svcSubnetCIDR.String())
			d.client = apiclient.NewDryRunClient(dryRunGetter, os.Stdout)
		} else {
			// If we're acting for real, we should create a connection to the API server and wait for it to come up
			var err error
			d.client, err = kubeconfigutil.ClientSetFromFile(d.KubeConfigPath())
			if err != nil {
				return nil, err
			}
		}
	}
	return d.client, nil
}

// Tokens returns an array of token strings.
func (d *initData) Tokens() []string {
	tokens := []string{}
	for _, bt := range d.cfg.BootstrapTokens {
		tokens = append(tokens, bt.Token.String())
	}
	return tokens
}

// PatchesDir returns the folder where patches for components are stored
func (d *initData) PatchesDir() string {
	// If provided, make the flag value override the one in config.
	if len(d.patchesDir) > 0 {
		return d.patchesDir
	}
	if d.cfg.Patches != nil {
		return d.cfg.Patches.Directory
	}
	return ""
}
