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
