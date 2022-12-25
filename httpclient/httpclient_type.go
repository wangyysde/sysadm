/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
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

Ref: https://docs.docker.com/registry/spec/api/
ErrorCode: 106xxx
*/

package httpclient

import (
	"time"
)

type RoundTripperData struct {
	Timeout time.Duration
	KeepAlive time.Duration
	Tlshandshaketimeout time.Duration
	DisableKeepAlives bool
	DisableCompression bool
	MaxIdleConns int
	MaxIdleConnsPerHost int 
	MaxConnsPerHost int
	IdleConnTimeout time.Duration
	InsecureSkipVerify bool
}

type RequestData struct {
	Key string
	Value string
}

type RequestParams struct {
	Headers []RequestData
	QueryData []*RequestData
	BasicAuthData map[string]string
	Method string
	Url string
}

var defaultHeaders []RequestData = []RequestData{
	{
		Key: "User-Agent", 
		Value: ("sysadmHttpClient-"+httpClientVer),
	},
}

