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
	"github.com/lithammer/dedent"
	"github.com/pkg/errors"
	"io"
	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmscheme "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
	kubeadmapiv1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta3"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/options"
	phases "k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/join"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/discovery"
	configutil "k8s.io/kubernetes/cmd/kubeadm/app/util/config"
	kubeconfigutil "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	joinWorkerNodeDoneMsg = dedent.Dedent(`
		This node has joined the cluster:
		* Certificate signing request was sent to apiserver and a response was received.
		* The Kubelet was informed of the new secure connection details.

		Run 'kubectl get nodes' on the control-plane to see this node join the cluster.

		`)

	joinControPlaneDoneTemp = template.Must(template.New("join").Parse(dedent.Dedent(`
		This node has joined the cluster and a new control plane instance was created:

		* Certificate signing request was sent to apiserver and approval was received.
		* The Kubelet was informed of the new secure connection details.
		* Control plane label and taint were applied to the new node.
		* The Kubernetes control plane instances scaled up.
		{{.etcdMessage}}

		To start administering your cluster from this node, you need to run the following as a regular user:

			mkdir -p $HOME/.kube
			sudo cp -i {{.KubeConfigPath}} $HOME/.kube/config
			sudo chown $(id -u):$(id -g) $HOME/.kube/config

		Run 'kubectl get nodes' to see this node join the cluster.

		`)))
)

// compile-time assert that the local data object satisfies the phases data interface.
var _ phases.JoinData = &joinData{}

// joinData defines all the runtime information used when running the kubeadm join workflow;
// this data is shared across all the phases that are included in the workflow.
type joinData struct {
	cfg                   *kubeadmapi.JoinConfiguration
	initCfg               *kubeadmapi.InitConfiguration
	tlsBootstrapCfg       *clientcmdapi.Config
	client                clientset.Interface
	ignorePreflightErrors sets.String
	outputWriter          io.Writer
	patchesDir            string
	dryRun                bool
	dryRunDir             string
}

type JoinParam struct {
	APIServerAdvertiseAddress string
	APIServerBindPort         int32
	//token
	TokenStr string
	//discovery-token-ca-cert-hash
	TokenDiscoveryCAHash string
	//cri-socket
	NodeCRISocket string
	//control-plane
	ControlPlane bool
	//certificate-key
	CertificateKey           string
	TokenDiscovery           string
	TokenDiscoverySkipCAHash string
	FileDiscovery            string
	TLSBootstrapToken        string
	NodeName                 string
	Out                      io.Writer
}

func StartJoin(joinParam JoinParam) error {
	data, err := newJoinData(joinParam)
	if err != nil {
		fmt.Printf("newJoinData: %v\n", err)
		return err
	}
	fmt.Printf("%#v\n", data.cfg)
	fmt.Printf("%#v\n", data.tlsBootstrapCfg)
	//预检测
	if err = runPreflight(data); err != nil {
		return err
	}
	//控制平面准备
	if err = controlPlanePreparePhase(data); err != nil {
		return err
	}
	//检查etcd
	if err = checkEtcdPhase(data); err != nil {
		return err
	}
	//启动kubelet
	if err = kubeletStartJoinPhase(data); err != nil {
		fmt.Println("===========================", err.Error())
		return err
	}
	//控制平面加入
	if err = controlPlaneJoinPhase(data); err != nil {
		fmt.Println("===========================", err.Error())
		return err
	}
	// if the node is hosting a new control plane instance
	if data.cfg.ControlPlane != nil {
		// outputs the join control plane done message and exit
		etcdMessage := ""
		if data.initCfg.Etcd.External == nil {
			etcdMessage = "* A new etcd member was added to the local/stacked etcd cluster."
		}

		ctx := map[string]string{
			"KubeConfigPath": kubeadmconstants.GetAdminKubeConfigPath(),
			"etcdMessage":    etcdMessage,
		}
		if err := joinControPlaneDoneTemp.Execute(data.outputWriter, ctx); err != nil {
			return err
		}

	} else {
		fmt.Fprint(data.outputWriter, joinWorkerNodeDoneMsg)
	}
	return nil
}

func newJoinData(data JoinParam) (*joinData, error) {
	adminKubeConfigPath := kubeadmconstants.GetAdminKubeConfigPath()
	// initialize the public kubeadm config API by applying defaults
	externalcfg := &kubeadmapiv1.JoinConfiguration{}

	// Add optional config objects to host flags.
	// un-set objects will be cleaned up afterwards (into newJoinData func)
	externalcfg.Discovery.File = &kubeadmapiv1.FileDiscovery{}
	externalcfg.Discovery.BootstrapToken = &kubeadmapiv1.BootstrapTokenDiscovery{}
	externalcfg.ControlPlane = &kubeadmapiv1.JoinControlPlane{}

	// This object is used for storage of control-plane flags.
	joinControlPlane := &kubeadmapiv1.JoinControlPlane{}

	// Apply defaults
	kubeadmscheme.Scheme.Default(externalcfg)
	kubeadmapiv1.SetDefaults_JoinControlPlane(joinControlPlane)

	setJoinConfig(data, externalcfg)

	if len(data.TokenStr) > 0 {
		if len(externalcfg.Discovery.TLSBootstrapToken) == 0 {
			externalcfg.Discovery.TLSBootstrapToken = data.TokenStr
		}
		if len(externalcfg.Discovery.BootstrapToken.Token) == 0 {
			externalcfg.Discovery.BootstrapToken.Token = data.TokenStr
		}
	}

	if len(externalcfg.Discovery.File.KubeConfigPath) == 0 {
		externalcfg.Discovery.File = nil
	}
	externalcfg.Discovery.BootstrapToken.APIServerEndpoint = fmt.Sprintf("%s:%d", data.APIServerAdvertiseAddress, data.APIServerBindPort)

	if !data.ControlPlane {
		// Use a defaulted JoinControlPlane object to detect if the user has passed
		// other control-plane related flags.
		defaultJCP := &kubeadmapiv1.JoinControlPlane{}
		kubeadmapiv1.SetDefaults_JoinControlPlane(defaultJCP)

		// This list must match the JCP flags in addJoinConfigFlags()
		joinControlPlaneFlags := []string{
			options.CertificateKey,
			options.APIServerAdvertiseAddress,
			options.APIServerBindPort,
		}

		if *externalcfg.ControlPlane != *defaultJCP {
			fmt.Printf("[preflight] WARNING: --%s is also required when passing control-plane "+
				"related flags such as [%s]", options.ControlPlane, strings.Join(joinControlPlaneFlags, ", "))
		}
		externalcfg.ControlPlane = nil
	}

	var tlsBootstrapCfg *clientcmdapi.Config
	if _, err := os.Stat(adminKubeConfigPath); err == nil && data.ControlPlane {
		// use the admin.conf as tlsBootstrapCfg, that is the kubeconfig file used for reading the kubeadm-config during discovery
		fmt.Printf("[preflight] found %s. Use it for skipping discovery", adminKubeConfigPath)
		tlsBootstrapCfg, err = clientcmd.LoadFromFile(adminKubeConfigPath)
		if err != nil {
			return nil, errors.Wrapf(err, "Error loading %s", adminKubeConfigPath)
		}
	}

	cfg, err := configutil.LoadOrDefaultJoinConfiguration("", externalcfg)
	if err != nil {
		return nil, err
	}

	if externalcfg.NodeRegistration.CRISocket != "" {
		cfg.NodeRegistration.CRISocket = externalcfg.NodeRegistration.CRISocket
	}

	if cfg.ControlPlane != nil {
		if err := configutil.VerifyAPIServerBindAddress(cfg.ControlPlane.LocalAPIEndpoint.AdvertiseAddress); err != nil {
			return nil, err
		}
	}

	return &joinData{
		cfg:                   cfg,
		tlsBootstrapCfg:       tlsBootstrapCfg,
		ignorePreflightErrors: sets.NewString(),
		outputWriter:          data.Out,
	}, nil
}

func setJoinConfig(data JoinParam, cfg *kubeadmapiv1.JoinConfiguration) {
	//如果是控制平面加入，再设置这两个属性
	if data.ControlPlane {
		cfg.ControlPlane.LocalAPIEndpoint.AdvertiseAddress = data.APIServerAdvertiseAddress
		cfg.ControlPlane.LocalAPIEndpoint.BindPort = data.APIServerBindPort
	}
	cfg.NodeRegistration.CRISocket = data.NodeCRISocket
	cfg.Discovery.BootstrapToken.CACertHashes = []string{data.TokenDiscoveryCAHash}
}

// CertificateKey returns the key used to encrypt the certs.
func (j *joinData) CertificateKey() string {
	if j.cfg.ControlPlane != nil {
		return j.cfg.ControlPlane.CertificateKey
	}
	return ""
}

// Cfg returns the JoinConfiguration.
func (j *joinData) Cfg() *kubeadmapi.JoinConfiguration {
	return j.cfg
}

// DryRun returns the DryRun flag.
func (j *joinData) DryRun() bool {
	return j.dryRun
}

// KubeConfigDir returns the path of the Kubernetes configuration folder or the temporary folder path in case of DryRun.
func (j *joinData) KubeConfigDir() string {
	if j.dryRun {
		return j.dryRunDir
	}
	return kubeadmconstants.KubernetesDir
}

// KubeletDir returns the path of the kubelet configuration folder or the temporary folder in case of DryRun.
func (j *joinData) KubeletDir() string {
	if j.dryRun {
		return j.dryRunDir
	}
	return kubeadmconstants.KubeletRunDirectory
}

// ManifestDir returns the path where manifest should be stored or the temporary folder path in case of DryRun.
func (j *joinData) ManifestDir() string {
	if j.dryRun {
		return j.dryRunDir
	}
	return kubeadmconstants.GetStaticPodDirectory()
}

// CertificateWriteDir returns the path where certs should be stored or the temporary folder path in case of DryRun.
func (j *joinData) CertificateWriteDir() string {
	if j.dryRun {
		return j.dryRunDir
	}
	return j.initCfg.CertificatesDir
}

// TLSBootstrapCfg returns the cluster-info (kubeconfig).
func (j *joinData) TLSBootstrapCfg() (*clientcmdapi.Config, error) {
	if j.tlsBootstrapCfg != nil {
		return j.tlsBootstrapCfg, nil
	}
	fmt.Println("[preflight] Discovering cluster-info")
	tlsBootstrapCfg, err := discovery.For(j.cfg)
	j.tlsBootstrapCfg = tlsBootstrapCfg
	return tlsBootstrapCfg, err
}

// InitCfg returns the InitConfiguration.
func (j *joinData) InitCfg() (*kubeadmapi.InitConfiguration, error) {
	if j.initCfg != nil {
		return j.initCfg, nil
	}
	if _, err := j.TLSBootstrapCfg(); err != nil {
		return nil, err
	}
	fmt.Println("[preflight] Fetching init configuration")
	initCfg, err := fetchInitConfigurationFromJoinConfiguration(j.cfg, j.tlsBootstrapCfg)
	j.initCfg = initCfg
	return initCfg, err
}

// Client returns the Client for accessing the cluster with the identity defined in admin.conf.
func (j *joinData) Client() (clientset.Interface, error) {
	if j.client != nil {
		return j.client, nil
	}
	path := filepath.Join(j.KubeConfigDir(), kubeadmconstants.AdminKubeConfigFileName)

	client, err := kubeconfigutil.ClientSetFromFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "[preflight] couldn't create Kubernetes client")
	}
	j.client = client
	return client, nil
}

// IgnorePreflightErrors returns the list of preflight errors to ignore.
func (j *joinData) IgnorePreflightErrors() sets.String {
	return j.ignorePreflightErrors
}

// OutputWriter returns the io.Writer used to write messages such as the "join done" message.
func (j *joinData) OutputWriter() io.Writer {
	return j.outputWriter
}

// PatchesDir returns the folder where patches for components are stored
func (j *joinData) PatchesDir() string {
	// If provided, make the flag value override the one in config.
	if len(j.patchesDir) > 0 {
		return j.patchesDir
	}
	if j.cfg.Patches != nil {
		return j.cfg.Patches.Directory
	}
	return ""
}

func fetchInitConfigurationFromJoinConfiguration(cfg *kubeadmapi.JoinConfiguration, tlsBootstrapCfg *clientcmdapi.Config) (*kubeadmapi.InitConfiguration, error) {
	// Retrieves the kubeadm configuration
	fmt.Println("[preflight] Retrieving KubeConfig objects")
	initConfiguration, err := fetchInitConfiguration(tlsBootstrapCfg)
	if err != nil {
		return nil, err
	}

	// Create the final KubeConfig file with the cluster name discovered after fetching the cluster configuration
	clusterinfo := kubeconfigutil.GetClusterFromKubeConfig(tlsBootstrapCfg)
	tlsBootstrapCfg.Clusters = map[string]*clientcmdapi.Cluster{
		initConfiguration.ClusterName: clusterinfo,
	}
	tlsBootstrapCfg.Contexts[tlsBootstrapCfg.CurrentContext].Cluster = initConfiguration.ClusterName

	// injects into the kubeadm configuration the information about the joining node
	initConfiguration.NodeRegistration = cfg.NodeRegistration
	if cfg.ControlPlane != nil {
		initConfiguration.LocalAPIEndpoint = cfg.ControlPlane.LocalAPIEndpoint
	}

	return initConfiguration, nil
}

func fetchInitConfiguration(tlsBootstrapCfg *clientcmdapi.Config) (*kubeadmapi.InitConfiguration, error) {
	// creates a client to access the cluster using the bootstrap token identity
	tlsClient, err := kubeconfigutil.ToClientSet(tlsBootstrapCfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to access the cluster")
	}

	// Fetches the init configuration
	initConfiguration, err := configutil.FetchInitConfigurationFromCluster(tlsClient, nil, "preflight", true, false)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch the kubeadm-config ConfigMap")
	}

	return initConfiguration, nil
}
