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
	https://datatracker.ietf.org/doc/rfc7235/

	errorCode: 70001xxx
*/

package httpclient

import(
	"net"
	"net/http"
	"time"
)

// TODO: the following parameters should be configurable in the future.
var (
	defaultTimeout time.Duration = 30 
	defaultKeepAlive time.Duration = 30
	defaultTlshandshaketimeout time.Duration = 10
	defaultDisableKeepAlives bool = false
	defaultDisableCompression bool = false
	defaultMaxIdleConns int = 10 
	defaultMaxIdleConnsPerHost int = http.DefaultMaxIdleConnsPerHost
	defaultMaxConnsPerHost int = 0
	defaultIdleConnTimeout time.Duration = 90
	defaultInsecureSkipVerify bool = false
)

// TODO: the following parameters should be configurable in the future.
var (
	timeout time.Duration = 30 
	localAddr net.Addr = nil // Ref https://pkg.go.dev/net#Addr
	keepAlive time.Duration = 30
	tlshandshaketimeout time.Duration = 10
	disableKeepAlives bool = false
	disableCompression bool = false
	maxIdleConns int = 10 
	maxIdleConnsPerHost int = http.DefaultMaxIdleConnsPerHost
	maxConnsPerHost int = 0
	idleConnTimeout time.Duration = 90
	httpClientVer = "v1.0"
)
