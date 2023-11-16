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

package app

import (
	networkingv1 "k8s.io/api/networking/v1"
)

func sortIngressByName(p, q interface{}) bool {
	pData, ok := p.(networkingv1.Ingress)
	if !ok {
		return false
	}
	qData, ok := q.(networkingv1.Ingress)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortIngressByCreatetime(p, q interface{}) bool {
	pData, ok := p.(networkingv1.Ingress)
	if !ok {
		return false
	}
	qData, ok := q.(networkingv1.Ingress)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}
