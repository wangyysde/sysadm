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
	"bytes"
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/pkg/errors"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	cmdutil "k8s.io/kubernetes/cmd/kubeadm/app/cmd/util"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
	"k8s.io/kubernetes/cmd/kubeadm/app/preflight"
	utilsexec "k8s.io/utils/exec"
	"text/template"
)

var (
	preflightExample = cmdutil.Examples(`
		# Run join pre-flight checks using a config file.
		kubeadm join phase preflight --config kubeadm-config.yaml
		`)

	notReadyToJoinControlPlaneTemp = template.Must(template.New("join").Parse(dedent.Dedent(`
		One or more conditions for hosting a new control plane instance is not satisfied.

		{{.Error}}

		Please ensure that:
		* The cluster has a stable controlPlaneEndpoint address.
		* The certificates that must be shared among control plane instances are provided.

		`)))
)

func runPreflight(j *joinData) error {
	fmt.Println("[preflight] Running pre-flight checks")

	// Start with general checks
	fmt.Println("[preflight] Running general checks")
	if err := preflight.RunJoinNodeChecks(utilsexec.New(), j.Cfg(), j.IgnorePreflightErrors()); err != nil {
		return err
	}

	initCfg, err := j.InitCfg()
	if err != nil {
		return err
	}

	// Continue with more specific checks based on the init configuration
	fmt.Println("[preflight] Running configuration dependant checks")
	if j.Cfg().ControlPlane != nil {
		// Checks if the cluster configuration supports
		// joining a new control plane instance and if all the necessary certificates are provided
		hasCertificateKey := len(j.CertificateKey()) > 0
		if err := checkIfReadyForAdditionalControlPlane(&initCfg.ClusterConfiguration, hasCertificateKey); err != nil {
			// outputs the not ready for hosting a new control plane instance message
			ctx := map[string]string{
				"Error": err.Error(),
			}

			var msg bytes.Buffer
			notReadyToJoinControlPlaneTemp.Execute(&msg, ctx)
			return errors.New(msg.String())
		}

		// run kubeadm init preflight checks for checking all the prerequisites
		fmt.Println("[preflight] Running pre-flight checks before initializing the new control plane instance")

		if err := preflight.RunInitNodeChecks(utilsexec.New(), initCfg, j.IgnorePreflightErrors(), true, hasCertificateKey); err != nil {
			return err
		}

		if j.DryRun() {
			fmt.Println("[preflight] Would pull the required images (like 'kubeadm config images pull')")
			return nil
		}

		fmt.Println("[preflight] Pulling images required for setting up a Kubernetes cluster")
		fmt.Println("[preflight] This might take a minute or two, depending on the speed of your internet connection")
		fmt.Println("[preflight] You can also perform this action in beforehand using 'kubeadm config images pull'")
		if err := preflight.RunPullImagesCheck(utilsexec.New(), initCfg, j.IgnorePreflightErrors()); err != nil {
			return err
		}
	}
	return nil
}

func checkIfReadyForAdditionalControlPlane(initConfiguration *kubeadmapi.ClusterConfiguration, hasCertificateKey bool) error {
	// blocks if the cluster was created without a stable control plane endpoint
	if initConfiguration.ControlPlaneEndpoint == "" {
		return errors.New("unable to add a new control plane instance to a cluster that doesn't have a stable controlPlaneEndpoint address")
	}

	if !hasCertificateKey {
		// checks if the certificates are provided and are still valid, not expired yet.
		if ret, err := certs.SharedCertificateExists(initConfiguration); !ret {
			return err
		}
	}

	return nil
}
