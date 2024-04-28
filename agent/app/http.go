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

package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	sysadmApiServer "sysadm/apiserver/app"
	"sysadm/httpclient"
)

func buildHttpClient() error {
	dailer, e := httpclient.BuildDailer(defaultTcpTimeout, defaultKeepAliveProbeInterval, "")

	if e != nil {
		return e
	}

	var rt http.RoundTripper = nil
	if RunData.IsTls {
		tlsConf, e := httpclient.BuildTlsClientConfig(RunData.Ca, RunData.Cert, RunData.Key, RunData.workingDir, RunData.InsecureSkipVerify)
		if e != nil {
			return e
		}
		rt, e = httpclient.BuildTlsRoundTripper(dailer, tlsConf, defaultTLSHandshakeTimeout, defaultIdleConnTimeout,
			defaultMaxIdleConns, defaultMaxIdleConnsPerHost, defaultMaxConnsPerHost, defaultReadBufferSize,
			defaultWriteBufferSize, defaultDisableKeepAives, defaultDisableCompression, defaultForceAttemptHTTP2)
		if e != nil {
			return e
		}
	} else {
		rt, e = httpclient.NewBuildRoundTripper(dailer, defaultIdleConnTimeout, defaultMaxIdleConns, defaultMaxIdleConnsPerHost,
			defaultMaxConnsPerHost, defaultReadBufferSize, defaultWriteBufferSize, defaultDisableKeepAives, defaultDisableCompression,
			defaultForceAttemptHTTP2)
		if e != nil {
			return e
		}
	}
	client := httpclient.BuildHttpClient(rt, defaultHTTPTimeOut)

	RunData.httpClient = client

	return nil
}

// buildGetCommandUrl build complete url address where agent send request to
func buildGetCommandUrl() string {

	address := RunData.Address
	port := RunData.Port
	uri := sysadmApiServer.GetCommandUri
	apiVersion := sysadmApiServer.ApiVersion
	url := ""

	if RunData.IsTls {
		if port == 443 {
			url = "https://" + address + "/api/" + apiVersion + "/" + uri
		} else {
			portStr := strconv.Itoa(port)
			url = "https://" + address + ":" + portStr + "/api/" + apiVersion + "/" + uri
		}
		return url
	}

	if port == 80 {
		url = "http://" + address + "/api/" + apiVersion + "/" + uri
	} else {
		portStr := strconv.Itoa(port)
		url = "http://" + address + ":" + portStr + "/api/" + apiVersion + "/" + uri
	}

	return url

}

func handleHTTPBody(body []byte) error {
	var gotCommand sysadmApiServer.CommandData = sysadmApiServer.CommandData{}

	err := json.Unmarshal(body, &gotCommand)
	if err != nil {
		return err
	}

	err = doRouteCommand(&gotCommand)

	return err
}
