package k8sclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	sigyaml "sigs.k8s.io/yaml"
	"strings"
)

func ApplyFromYaml(yamlContent string, config *restclient.Config, dyClient *dynamic.DynamicClient, apiRL []*metav1.APIResourceList) error {
	yamlContent = strings.TrimSpace(yamlContent)
	if len(yamlContent) < 20 {
		return fmt.Errorf("yamlContent %s is not valid", yamlContent)
	}

	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(yamlContent), len(yamlContent))
	for {
		var rawObj runtime.RawExtension
		e := decoder.Decode(&rawObj)
		if e == io.EOF {
			break
		}
		if e != nil {
			return fmt.Errorf("decode yaml content %s error %s", yamlContent, e)
		}

		//obj, gvk, e := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		obj, _, e := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if e != nil {
			return fmt.Errorf("yaml content %s can not be serialized or can not decoded error %s", yamlContent, e)
		}

		unstructuredMap, e := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if e != nil {
			return fmt.Errorf("unstructrue object error %s ", e)
		}
		unstructureObj := &unstructured.Unstructured{Object: unstructuredMap}
		gvr, e := GetGVR(config, unstructureObj.GroupVersionKind())
		if e != nil {
			return fmt.Errorf("get gvr of the object error %s", e)
		}
		unstructuredYaml, e := sigyaml.Marshal(unstructureObj)
		if e != nil {
			return fmt.Errorf("unable to marshal resource as yaml: %w", e)
		}

		if len(apiRL) < 1 {
			apiRL, e = GetApiResourcesList(config)
			if e != nil {
				return e
			}
		}
		namespaced, e := IsNamespaced(unstructureObj.GetKind(), apiRL)
		if e != nil {
			return e
		}
		namespace := strings.TrimSpace(unstructureObj.GetNamespace())
		if namespace == "" && namespaced {
			namespace = "default"
		}
		var getErr error = nil
		if namespaced {
			_, getErr = dyClient.Resource(gvr).Namespace(namespace).Get(context.Background(), unstructureObj.GetName(), metav1.GetOptions{})
		} else {
			_, getErr = dyClient.Resource(gvr).Get(context.Background(), unstructureObj.GetName(), metav1.GetOptions{})
		}

		if getErr != nil {
			var createErr error = nil
			if namespaced {
				_, createErr = dyClient.Resource(gvr).Namespace(namespace).Create(context.Background(), unstructureObj, metav1.CreateOptions{})
			} else {
				_, createErr = dyClient.Resource(gvr).Create(context.Background(), unstructureObj, metav1.CreateOptions{})
			}
			if createErr != nil {
				return fmt.Errorf("create resource %s error %s namespaced %v namespace: %s \n", createErr, unstructureObj.GetName(), namespaced, namespace)
			}
		} else {
			var applyErr error = nil
			applyOption := metav1.ApplyOptions{Force: true, FieldManager: unstructureObj.GetName()}
			if namespaced {
				_, applyErr = dyClient.Resource(gvr).Namespace(namespace).Apply(context.Background(), unstructureObj.GetName(), unstructureObj, applyOption)
			} else {
				//_, updateErr = dyClient.Resource(gvr).Update(context.Background(), unstructureObj, metav1.UpdateOptions{})
				_, applyErr = dyClient.Resource(gvr).Apply(context.Background(), unstructureObj.GetName(), unstructureObj, applyOption)
			}
			if applyErr != nil {
				return fmt.Errorf("apply resource error %s", applyErr)
			}
		}

		force := true
		if namespaced {
			_, e = dyClient.Resource(gvr).
				Namespace(namespace).
				Patch(context.Background(),
					unstructureObj.GetName(),
					types.ApplyPatchType,
					unstructuredYaml, metav1.PatchOptions{
						FieldManager: unstructureObj.GetName(),
						Force:        &force,
					})

			if e != nil {
				return fmt.Errorf("unable to patch namespaced resource: %w", e)
			}
		} else {
			_, e = dyClient.Resource(gvr).
				Patch(context.Background(),
					unstructureObj.GetName(),
					types.ApplyPatchType,
					unstructuredYaml, metav1.PatchOptions{
						Force:        &force,
						FieldManager: unstructureObj.GetName(),
					})
			if e != nil {
				return fmt.Errorf("unable to patch unnamespaced resource: %w", e)
			}
		}

		fmt.Printf("%s \"%s\" applyed\n", unstructureObj.GetKind(), unstructureObj.GetName())
	}

	return nil
}

func ApplyFromYamlByClientSet(yamlContent string, clientSet *kubernetes.Clientset) error {

	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(yamlContent), len(yamlContent))
	for {
		var rawObj runtime.RawExtension
		e := decoder.Decode(&rawObj)
		if e == io.EOF {
			break
		}
		if e != nil {
			return fmt.Errorf("decode yaml content %s error %s", yamlContent, e)
		}

		//obj, gvk, e := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		obj, _, e := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if e != nil {
			return fmt.Errorf("yaml content %s can not be serialized or can not decoded error %s", yamlContent, e)
		}

		unstructuredMap, e := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if e != nil {
			return fmt.Errorf("unstructrue object error %s ", e)
		}
		unstructureObj := &unstructured.Unstructured{Object: unstructuredMap}
		gvr, e := GetGVRByClientSet(clientSet, unstructureObj.GroupVersionKind())
		if e != nil {
			return fmt.Errorf("get gvr of the object error %s", e)
		}

		//	unstructuredYaml, e := sigyaml.Marshal(unstructureObj)
		//	if e != nil {
		//		return fmt.Errorf("unable to marshal resource as yaml: %w", e)
		//	}

		apiRL, e := discovery.ServerPreferredResources(clientSet.DiscoveryClient)
		if e != nil {
			return e
		}

		namespaced, e := IsNamespaced(unstructureObj.GetKind(), apiRL)
		if e != nil {
			return e
		}
		namespace := strings.TrimSpace(unstructureObj.GetNamespace())
		if namespace == "" && namespaced {
			namespace = "default"
		}
		var getErr error = nil
		dyClient := dynamic.New(clientSet.RESTClient())
		if namespaced {
			_, getErr = dyClient.Resource(gvr).Namespace(namespace).Get(context.Background(), unstructureObj.GetName(), metav1.GetOptions{})
		} else {
			_, getErr = dyClient.Resource(gvr).Get(context.Background(), unstructureObj.GetName(), metav1.GetOptions{})
		}

		if getErr != nil {
			var createErr error = nil
			if namespaced {
				_, createErr = dyClient.Resource(gvr).Namespace(namespace).Create(context.Background(), unstructureObj, metav1.CreateOptions{})
			} else {
				_, createErr = dyClient.Resource(gvr).Create(context.Background(), unstructureObj, metav1.CreateOptions{})
			}
			if createErr != nil {
				return fmt.Errorf("create resource %s error %s namespaced %v namespace: %s \n", createErr, unstructureObj.GetName(), namespaced, namespace)
			}
		} else {
			var applyErr error = nil
			applyOption := metav1.ApplyOptions{Force: true, FieldManager: unstructureObj.GetName()}
			if namespaced {
				_, applyErr = dyClient.Resource(gvr).Namespace(namespace).Apply(context.Background(), unstructureObj.GetName(), unstructureObj, applyOption)
			} else {
				//_, updateErr = dyClient.Resource(gvr).Update(context.Background(), unstructureObj, metav1.UpdateOptions{})
				_, applyErr = dyClient.Resource(gvr).Apply(context.Background(), unstructureObj.GetName(), unstructureObj, applyOption)
			}
			if applyErr != nil {
				return fmt.Errorf("apply resource error %s", applyErr)
			}

		}
	}

	return nil
}
