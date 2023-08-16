/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
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
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"strings"
)

func GetPodInfo(restConf *rest.Config, ns, podName string) (*corev1.Pod, error) {
	var ret *corev1.Pod = nil

	if restConf == nil {
		return ret, fmt.Errorf("can not get pod information on an empty client")
	}

	ns = strings.TrimSpace(ns)
	if ns == "" {
		ns = corev1.NamespaceDefault
	}

	podName = strings.TrimSpace(podName)
	if podName == "" {
		return ret, fmt.Errorf("pod name should be specifioed")
	}

	clientSet, e := BuildClientset(restConf)
	if e != nil {
		return ret, e
	}

	pod, e := clientSet.CoreV1().Pods(ns).Get(context.Background(), podName, metav1.GetOptions{})
	if e != nil {
		return ret, e
	}

	return pod, nil
}

func GetPodInfoWithPrefix(restConf *rest.Config, ns, prefix string) ([]corev1.Pod, error) {
	var ret []corev1.Pod

	if restConf == nil {
		return ret, fmt.Errorf("can not get pod information on an empty client")
	}

	ns = strings.TrimSpace(ns)
	if ns == "" {
		ns = corev1.NamespaceDefault
	}

	prefix = strings.TrimSpace(prefix)
	clientSet, e := BuildClientset(restConf)
	if e != nil {
		return ret, e
	}

	pods, e := clientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return ret, e
	}

	for _, p := range pods.Items {
		podName := p.Name
		if strings.HasPrefix(podName, prefix) {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func GetPodCount(restConf *rest.Config, ns string) (ObjectCount, error) {
	ns = strings.TrimSpace(ns)
	ret := ObjectCount{Namespace: ""}

	if restConf == nil {
		return ret, fmt.Errorf("can not get pod count on an empty client")
	}

	clientSet, e := BuildClientset(restConf)
	if e != nil {
		return ret, e
	}

	pods, e := clientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return ret, e
	}

	total := 0
	ready := 0
	unready := 0
	for _, p := range pods.Items {
		conditions := p.Status.Conditions
		statusReady := false
		for _, c := range conditions {
			if c.Type != corev1.PodReady {
				continue
			}
			if c.Status == corev1.ConditionTrue {
				statusReady = true
			}
		}
		total = total + 1
		if statusReady {
			ready = ready + 1
		} else {
			unready = unready + 1
		}
	}

	ret.Kind = PodKind
	ret.Total = int32(total)
	ret.Ready = int32(ready)
	ret.Unready = int32(unready)

	return ret, nil

}
