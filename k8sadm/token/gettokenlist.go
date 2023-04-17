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
	"context"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcertutil "k8s.io/client-go/util/cert"
	bootstrapapi "k8s.io/cluster-bootstrap/token/api"
	bootstraptokenv1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/bootstraptoken/v1"
	outputapiv1alpha2 "k8s.io/kubernetes/cmd/kubeadm/app/apis/output/v1alpha2"
	cmdutil "k8s.io/kubernetes/cmd/kubeadm/app/cmd/util"
	kubeconfigutil "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/pubkeypin"
)

func GetToken() error {
	kubeConfigFile := cmdutil.GetKubeConfigPath("")
	client, err := cmdutil.GetClientSet(kubeConfigFile, false)
	if err != nil {
		return err
	}
	tokens, err := runListTokens(client)
	for _, token := range tokens {
		fmt.Printf("%s\t%v\n", token.BootstrapToken.Token.String(), token.Usages)
	}
	ca, _ := getCertificate(kubeConfigFile)
	fmt.Println(ca)
	return nil
}

func runListTokens(client clientset.Interface) ([]outputapiv1alpha2.BootstrapToken, error) {
	var list []outputapiv1alpha2.BootstrapToken
	fmt.Println("[token] preparing selector for bootstrap token")
	tokenSelector := fields.SelectorFromSet(
		map[string]string{
			"type": string(bootstrapapi.SecretTypeBootstrapToken),
		},
	)
	listOptions := metav1.ListOptions{
		FieldSelector: tokenSelector.String(),
	}

	fmt.Println("[token] retrieving list of bootstrap tokens")
	secrets, err := client.CoreV1().Secrets(metav1.NamespaceSystem).List(context.TODO(), listOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list bootstrap tokens")
	}
	for _, secret := range secrets.Items {
		token, err := bootstraptokenv1.BootstrapTokenFromSecret(&secret)
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}

		outputToken := outputapiv1alpha2.BootstrapToken{
			BootstrapToken: bootstraptokenv1.BootstrapToken{
				Token:       &bootstraptokenv1.BootstrapTokenString{ID: token.Token.ID, Secret: token.Token.Secret},
				Description: token.Description,
				TTL:         token.TTL,
				Expires:     token.Expires,
				Usages:      token.Usages,
				Groups:      token.Groups,
			},
		}
		list = append(list, outputToken)
	}
	return list, nil
}

func getCertificate(kubeConfigFile string) ([]string, error) {
	// load the kubeconfig file to get the CA certificate and endpoint
	config, err := clientcmd.LoadFromFile(kubeConfigFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load kubeconfig")
	}

	// load the default cluster config
	clusterConfig := kubeconfigutil.GetClusterFromKubeConfig(config)
	if clusterConfig == nil {
		return nil, errors.New("failed to get default cluster config")
	}
	// load CA certificates from the kubeconfig (either from PEM data or by file path)
	var caCerts []*x509.Certificate
	if clusterConfig.CertificateAuthorityData != nil {
		caCerts, err = clientcertutil.ParseCertsPEM(clusterConfig.CertificateAuthorityData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse CA certificate from kubeconfig")
		}
	} else if clusterConfig.CertificateAuthority != "" {
		caCerts, err = clientcertutil.CertsFromFile(clusterConfig.CertificateAuthority)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load CA certificate referenced by kubeconfig")
		}
	} else {
		return nil, errors.New("no CA certificates found in kubeconfig")
	}

	// hash all the CA certs and include their public key pins as trusted values
	publicKeyPins := make([]string, 0, len(caCerts))
	for _, caCert := range caCerts {
		publicKeyPins = append(publicKeyPins, pubkeypin.Hash(caCert))
	}
	return publicKeyPins, nil
}
