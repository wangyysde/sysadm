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
	"fmt"
	"io"
	"k8s.io/client-go/kubernetes"
	"os"
	"testing"
)

func TestAppyYamlByClientSet(t *testing.T) {
	clusterID, apiserver, clusterUser, ca, cert, key, e := getClusterData()
	if e != nil {
		t.Fatal(e)
	}

	restConf, e := BuildConfigFromParametes([]byte(ca), []byte(cert), []byte(key), apiserver, clusterID, clusterUser, "default")
	if e != nil {
		t.Fatal(e)
	}

	k8sClient, e := kubernetes.NewForConfig(restConf)
	if e != nil {
		t.Fatal(e)
	}

	fp, e := os.Open("./testbyapp.yaml")
	if e != nil {
		t.Fatal(e)
	}

	yamlContent, e := io.ReadAll(fp)
	if e != nil {
		t.Fatal(e)
	}
	fp.Close()
	
	yamlContentStr := string(yamlContent)
	e = ApplyFromYamlByClientSet(yamlContentStr, k8sClient)
	if e != nil {
		t.Fatal(e)
	}

	fmt.Printf("namespace created\n")

}
