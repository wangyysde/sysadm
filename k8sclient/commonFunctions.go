package k8sclient

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

func BuildKubeConf(masterUrl, kubeconfigPath string) (*restclient.Config, error) {
	return clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
}

func BuildDynamicClient(config *restclient.Config) (*dynamic.DynamicClient, error) {
	return dynamic.NewForConfig(config)
}

func BuildClientset(config *restclient.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(config)
}

func GetGVR(config *restclient.Config, gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	clientset, e := BuildClientset(config)
	if e != nil {
		return schema.GroupVersionResource{}, fmt.Errorf("build clientset error %s", e)
	}

	gr, e := restmapper.GetAPIGroupResources(clientset.Discovery())
	if e != nil {
		return schema.GroupVersionResource{}, e
	}

	mapper := restmapper.NewDiscoveryRESTMapper(gr)

	mapping, e := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if e != nil {
		return schema.GroupVersionResource{}, e
	}

	return mapping.Resource, nil
}

func GetApiResourcesList(config *restclient.Config) ([]*v1.APIResourceList, error) {
	var ret []*v1.APIResourceList
	clientset, e := kubernetes.NewForConfig(config)
	if e != nil {
		return ret, fmt.Errorf("build clientset error %s", e)
	}
	return discovery.ServerPreferredResources(clientset.Discovery())
}

func IsNamespaced(objKind string, apiRl []*v1.APIResourceList) (bool, error) {
	objKind = strings.ToLower(strings.TrimSpace(objKind))
	if objKind == "" {
		return false, fmt.Errorf("object kind must be not empty")
	}

	if len(apiRl) < 1 {
		return false, fmt.Errorf("api resource list must be not empty")
	}

	namespaced := false
	found := false
	for _, rl := range apiRl {
		for _, r := range rl.APIResources {
			kind := strings.ToLower(strings.TrimSpace(r.Kind))
			if objKind == kind {
				found = true
				namespaced = r.Namespaced
				break
			}
		}
		if found {
			break
		}
	}

	return namespaced, nil
}
