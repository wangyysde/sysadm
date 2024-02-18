/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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
 */

package k8sclient

import (
	"fmt"
	restclient "k8s.io/client-go/rest"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"strings"
)

func BuildConfigFromParametes(caData, certData, keyData []byte, apiServer, clusterName, userName, nameSpace string) (*restclient.Config, error) {
	if len(caData) < 1 || len(certData) < 1 || len(keyData) < 1 {
		return nil, fmt.Errorf("ca certification or certification or key is not valid")
	}

	apiServer = strings.ToLower(strings.TrimSpace(apiServer))
	if len(apiServer) < 1 {
		return nil, fmt.Errorf("the address of api server is not valid")
	}
	if !strings.HasPrefix(apiServer, "https://") {
		apiServer = "https://" + apiServer
	}

	conf := clientcmdapi.Config{Kind: "Config", APIVersion: "v1"}
	cluster := clientcmdapi.Cluster{Server: apiServer, CertificateAuthorityData: caData}
	clusterMap := make(map[string]*clientcmdapi.Cluster, 0)
	clusterMap[clusterName] = &cluster
	conf.Clusters = clusterMap

	authInfo := clientcmdapi.AuthInfo{ClientCertificateData: certData, ClientKeyData: keyData, Username: userName}
	authInfoMap := make(map[string]*clientcmdapi.AuthInfo, 0)
	authInfoMap[userName] = &authInfo
	conf.AuthInfos = authInfoMap

	context := clientcmdapi.Context{Cluster: clusterName, AuthInfo: userName, Namespace: nameSpace}
	contextMap := make(map[string]*clientcmdapi.Context, 0)
	contextMap[clusterName] = &context
	conf.Contexts = contextMap
	conf.CurrentContext = clusterName

	confByte, e := clientcmd.Write(conf)
	if e != nil {
		return nil, e
	}

	return clientcmd.RESTConfigFromKubeConfig(confByte)
}

func BuildConfigFromParasWithConnectType(connectType, apiServer, clusterName, clusterUser, namespace, ca, cert, key, token,
	kubeConfig string) (*restclient.Config, error) {

	apiServer = strings.ToLower(strings.TrimSpace(apiServer))

	switch connectType {
	case "0":
		if len(apiServer) < 1 {
			return nil, fmt.Errorf("the address of api server is not valid")
		}
		if !strings.HasPrefix(apiServer, "https://") && !strings.HasPrefix(apiServer, "http://") {
			apiServer = "https://" + apiServer
		}
		return BuildConfigByCert(ca, cert, key, apiServer, clusterName, clusterUser, namespace)
	case "1":
		if len(apiServer) < 1 {
			return nil, fmt.Errorf("the address of api server is not valid")
		}
		if !strings.HasPrefix(apiServer, "https://") && !strings.HasPrefix(apiServer, "http://") {
			apiServer = "https://" + apiServer
		}
		return BuildConfigByToken(apiServer, ca, token)
	case "2":
		return BuildConfigByKubeConf(kubeConfig)
	}

	return nil, fmt.Errorf("connect type %s is not valid", connectType)
}

func BuildConfigByCert(caData, certData, keyData, apiServer, clusterName, userName, nameSpace string) (*restclient.Config, error) {
	conf := clientcmdapi.Config{Kind: "Config", APIVersion: "v1"}
	cluster := clientcmdapi.Cluster{Server: apiServer, CertificateAuthorityData: []byte(caData)}
	clusterMap := make(map[string]*clientcmdapi.Cluster, 0)
	clusterMap[clusterName] = &cluster
	conf.Clusters = clusterMap

	userName = strings.TrimSpace(userName)
	if userName == "" {
		userName = defaultClusterUserName
	}
	authInfo := clientcmdapi.AuthInfo{ClientCertificateData: []byte(certData), ClientKeyData: []byte(keyData), Username: userName}
	authInfoMap := make(map[string]*clientcmdapi.AuthInfo, 0)
	authInfoMap[userName] = &authInfo
	conf.AuthInfos = authInfoMap

	nameSpace = strings.TrimSpace(nameSpace)
	if nameSpace == "" {
		nameSpace = defaultNamespace
	}

	context := clientcmdapi.Context{Cluster: clusterName, AuthInfo: userName, Namespace: nameSpace}
	contextMap := make(map[string]*clientcmdapi.Context, 0)
	contextMap[clusterName] = &context
	conf.Contexts = contextMap
	conf.CurrentContext = clusterName

	confByte, e := clientcmd.Write(conf)
	if e != nil {
		return nil, e
	}

	return clientcmd.RESTConfigFromKubeConfig(confByte)
}

func BuildConfigByKubeConf(kubeConfig string) (*restclient.Config, error) {
	kubeConfig = strings.TrimSpace(kubeConfig)

	return clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
}

// func BuildConfigByToken(apiServer, caData, token string) (*restclient.Config, error) {
func BuildConfigByToken(apisServer, ca, token string) (*restclient.Config, error) {
	tlsClientConfig := restclient.TLSClientConfig{}
	tlsClientConfig.CAData = []byte(ca)
	bearer := strings.TrimSpace(token)
	config := restclient.Config{
		// TODO: switch to using cluster DNS.
		Host:            apisServer,
		TLSClientConfig: tlsClientConfig,
		BearerToken:     bearer,
		BearerTokenFile: "",
	}
	return &config, nil
}

func GetClusterDefaultUser() string {
	return defaultClusterUserName
}
