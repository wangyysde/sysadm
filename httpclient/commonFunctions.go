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

package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// BuildUrl build completion url address which will be send to a server
func BuildUrl(server, uri string, port int, isTLS bool) (string, error) {
	var url string = ""

	server = strings.TrimSpace(server)
	uri = strings.TrimSpace(uri)

	if server == "" {
		return "", fmt.Errorf("server address must be not empty")
	}

	if strings.TrimSpace(uri) == "" {
		uri = "/"
	}

	if port == 0 {
		if isTLS {
			port = 443
		} else {
			port = 80
		}
	}

	url = "http://" + server
	if isTLS {
		url = "https://" + server
	}

	if !isTLS && port != 80 {
		url = url + ":" + strconv.Itoa(port)
	}

	if isTLS && port != 443 {
		url = url + ":" + strconv.Itoa(port)
	}

	if uri[0:1] == "/" {
		url = url + uri
	} else {
		url = url + "/" + uri
	}

	return url, nil
}

// SendTlsRequest build HTTP query data, header information, BasicAuth, create a new http request and then
// send the request to server by client.
// return response body and nil if successful, otherwise return empty []byte and an error
func NewSendRequest(r *RequestParams, client *http.Client, bodyReader io.Reader) ([]byte, error) {
	var body []byte

	if r == nil {
		return body, fmt.Errorf("can not handle http request without any request request parameters")
	}

	queryData, err := newHandleQueryData(r)
	if err != nil {
		return body, err
	}

	r.Url = strings.TrimSpace(r.Url)
	if r.Url == "" {
		return body, fmt.Errorf("HTTP request Url must not empty")
	}

	if queryData != "" {
		r.Url = r.Url + "?" + queryData
	}

	if client == nil {
		return body, fmt.Errorf("http client must not nil")
	}

	if !CheckHttpMethod(r.Method) {
		return body, fmt.Errorf("HTTP method is not valid")
	}

	req, err := http.NewRequest(strings.ToUpper(r.Method), r.Url, bodyReader)
	if err != nil {
		return body, fmt.Errorf("create new HTTP request error %s", err)
	}

	if err := newAddReqestHeader(r, req); err != nil {
		return body, fmt.Errorf("add request header information onto request  error %s", err)
	}

	if err := newSetBasicAuth(r, req); err != nil {
		return body, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}
