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

const (
	// application version
	ver = "v1.0.0"

	// author of the application
	author string = "Wayne Wang<net_use@bzhy.com>"

	// default path of log file
	defaultLogFile = "/var/log/agent.log"

	// Timeout is the maximum amount of time a dial will wait for a connect to complete. When using TCP and dialing a host name with multiple IP
	// addresses, the timeout may be divided between them. This value is for build net.Dialer for a http client.
	defaultTcpTimeout int = 180

	// DefaultKeepAliveProbeInterval specifies the interval between keep-alive probes for an active network connection. If negative, keep-alive probes are
	// disabled.
	defaultKeepAliveProbeInterval int = 15

	// TLSHandshakeTimeout specifies the maximum amount of time waiting to wait for a TLS handshake. Zero means no timeout.
	defaultTLSHandshakeTimeout int = 180

	// IdleConnTimeout is the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself. Zero means no limit.
	defaultIdleConnTimeout int = 300

	// MaxIdleConns controls the maximum number of idle (keep-alive) connections across all hosts. Zero means no limit.
	defaultMaxIdleConns int = 5

	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle (keep-alive) connections to keep per-host.
	defaultMaxIdleConnsPerHost int = 2

	// MaxConnsPerHost optionally limits the total number of connections per host, including connections in the dialing, active, and idle states. Zero means no limit.
	defaultMaxConnsPerHost int = 5

	// ReadBufferSize specifies the size of the read buffer used when reading from the transport. If zero, a default (currently 4KB) is used.
	defaultReadBufferSize int = 4096

	// WriteBufferSize specifies the size of the write buffer used when writing to the transport.If zero, a default (currently 4KB) is used.
	defaultWriteBufferSize int = 4096

	// DisableKeepAlives, if true, disables HTTP keep-alives and will only use the connection to the server for a single HTTP request
	defaultDisableKeepAives bool = false

	// DisableCompression, if true, prevents the Transport from requesting compression with an "Accept-Encoding: gzip"
	defaultDisableCompression bool = false

	// ForceAttemptHTTP2 controls whether HTTP/2 is enabled when a non-zero Dial, DialTLS, or DialContext func or TLSClientConfig is provided.
	// By default, use of any those fields conservatively disables HTTP/2.To use a custom dialer or TLS config and still attempt HTTP/2
	// upgrades, set this to true.
	defaultForceAttemptHTTP2 bool = true

	// Timeout specifies a time limit for requests made by this Client. The timeout includes connection time, any redirects, and reading the response body. The timer remains
	defaultHTTPTimeOut int = 30

	//period(second) for agent gets command from apiServer
	defaultGetCommandInterval int = 60
)
