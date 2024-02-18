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
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"testing"
)

func TestBuildConfigFromParasWithConnectType(t *testing.T) {
	var apiServer = "9.9.33.36:6443"
	var clusterName = "kubernetes"
	var clusterUser = "kubernetes-admin"
	var namespace = "default"
	var caFile = "./data/ca.crt"
	var certFile = "./data/cert.pem"
	var keyFile = "./data/key.pem"
	var tokenFile = "./data/token"
	var kubeConfigFile = "./data/kubeconfig"

	caByte, e1 := os.ReadFile(caFile)
	if e1 != nil {
		t.Fatal(e1)
		return
	}

	certByte, e := os.ReadFile(certFile)
	if e != nil {
		t.Fatal(e)
		return
	}

	keyByte, e := os.ReadFile(keyFile)
	if e != nil {
		t.Fatal(e)
		return
	}

	tokenByte, e := os.ReadFile(tokenFile)
	if e != nil {
		t.Fatal(e)
		return
	}

	kubeconfigByte, e := os.ReadFile(kubeConfigFile)
	if e != nil {
		t.Fatal(e)
		return
	}

	tokenRestConf, e := BuildConfigFromParasWithConnectType("1", apiServer, "", "", "", string(caByte), "", "", string(tokenByte), "")
	if e != nil {
		t.Fatal(e)
	}
	tokenClient, e := kubernetes.NewForConfig(tokenRestConf)
	if e != nil {
		t.Fatal(e)
	}
	podList, e := tokenClient.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("got total %d pods\n", len(podList.Items))

	certRestConf, e := BuildConfigFromParasWithConnectType("0", apiServer, clusterName, clusterUser, namespace, string(caByte), string(certByte), string(keyByte), "", "")
	if e != nil {
		t.Fatal(e)
		return
	}

	k8sClient, e := kubernetes.NewForConfig(certRestConf)
	if e != nil {
		t.Fatal(e)
	}

	deployList, e := k8sClient.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("got total %d deployments\n", len(deployList.Items))

	kubeconfRestConf, e := BuildConfigFromParasWithConnectType("2", apiServer, clusterName, clusterUser, namespace, "", "", "", "", string(kubeconfigByte))
	if e != nil {
		t.Fatal(e)
	}
	kubeconfClient, e := kubernetes.NewForConfig(kubeconfRestConf)
	if e != nil {
		t.Fatal(e)
	}
	nsList, e := kubeconfClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("got total %d namespaces \n", len(nsList.Items))

	t.Logf("test ok\n")
}
