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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"strings"
)

func GetStatefulSetCount(restConf *rest.Config, ns string) (ObjectCount, error) {
	ns = strings.TrimSpace(ns)
	ret := ObjectCount{Namespace: ""}

	if restConf == nil {
		return ret, fmt.Errorf("can not get deployment count on an empty client")
	}

	clientSet, e := BuildClientset(restConf)
	if e != nil {
		return ret, e
	}

	sts, e := clientSet.AppsV1().StatefulSets(ns).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return ret, e
	}

	total := 0
	ready := 0
	unready := 0
	for _, s := range sts.Items {
		numReadyRs := s.Status.ReadyReplicas
		numCurrRs := s.Status.CurrentReplicas
		if numReadyRs == numCurrRs {
			ready = ready + 1
		} else {
			unready = unready + 1
		}
	}

	ret.Kind = statefulsetKind
	ret.Total = int32(total)
	ret.Ready = int32(ready)
	ret.Unready = int32(unready)

	return ret, nil
}
