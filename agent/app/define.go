/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* https://www.sysadm.cn/licenses/apache-2.0.txt
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
*
 */

package app

var ver string ="1.0.0"
var author string ="Wayne Wang<net_use@bzhy.com>"
// the port of a server listen which agent send the reponses message to 
// this value used by cmd
var  DefaultServerPort int = 443
// default path of configuration file of agent. this path is relative to working directory
var DefaultConf string = "conf/agent.yaml"
// logfile for agent which is used to log runing log messages of agent to 
var DefaultLogFile string = "/var/log/sysad-agent.log"
//  log message with the format(kind) will be output. its value is one of "text" and "json". default value is text 
var DefaultLogKind = "text"
//  just the log messages will be output what the level of the log message is higher "logLevel".
var DefaultLogLevel = "debug"
// default the format of time in the log message
var DefaultTimeStampFormat = "2006-01-02 15:04:05"
// whether check the certs which got from a server isvalid
var DefaultskipVerifyCert bool = false
//specifies whether agent using TLS protocol when it is communicate with  a server (agent send the reponses message to the server)
var DefaultServerIsTls bool = false
//  the method of getting commands by agent. agent gets commands from the server periodically and run them if this value is true
var DefaultPassive bool = false
// listen port of agent using when agent running as daemon.
var  DefaultListenPort int = 5443
//period(second) for agent gets command from server when agent running as passive
var DefaultPeriod int = 60
// insecret specifies whether agent listen on a insecret port when it is runing as daemon
var DefaultInsecret bool = false
// insecret listen port of agent listening when it is running ad daemon 
var DefaultInsecretPort int = 5080
// set agent runs in debug mode defaultly
var DefalutDebugMode bool = true
// specifies a identifer of the node which agent running on it.
// It is any combination of the IP,HOSTNAME and MAC joined by commas  or a customize string what the leght of the string is less 63
// agent will get all IPs without not active and reponse these IPs in list to the server by nodeIdentifer.IPs filed if IP is included in NodeIdentifer
// agent will get hostname and reponse the hostname  to the server by nodeIdentifer.Hostname filed if hostname is included in NodeIdentifer
// agent will get all MACs without not active and reponse these MACs in list to the server by nodeIdentifer.MACs filed if MAC is included in NodeIdentifer
// customize string is reponse to the server directly .
// customize string is conflicted with IP,HOSTNAME and MAC. the nodeIdentifer can be changed by the server during agent communicate with the server
var DefaultNodeIdentifer string = "IP,HOSTNAME,MAC"
// in active mode, if the path where agent receive command fro is not set, then its value should be set to defaultReceiveCommandUri
var defaultReceiveCommandUri string = "/receiveCommand"
// Timeout is the maximum amount of time a dial will wait for a connect to complete. When using TCP and dialing a host name with multiple IP 
// addresses, the timeout may be divided between them. This value is for build net.Dialer for a http client.
var DefaultTcpTimeout int = 180
// DefaultKeepAliveProbeInterval specifies the interval between keep-alive probes for an active network connection. If negative, keep-alive probes are 
// disabled.
var DefaultKeepAliveProbeInterval int = 15
// TLSHandshakeTimeout specifies the maximum amount of time waiting to wait for a TLS handshake. Zero means no timeout.
var defaultTLSHandshakeTimeout int = 180
// IdleConnTimeout is the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself. Zero means no limit.
var defaultIdleConnTimeout int = 300
// MaxIdleConns controls the maximum number of idle (keep-alive) connections across all hosts. Zero means no limit.
var defaultMaxIdleConns int = 5
// MaxIdleConnsPerHost, if non-zero, controls the maximum idle (keep-alive) connections to keep per-host.
var defaultMaxIdleConnsPerHost  int = 2
// MaxConnsPerHost optionally limits the total number of connections per host, including connections in the dialing, active, and idle states. Zero means no limit.
var defaultMaxConnsPerHost int = 5
// ReadBufferSize specifies the size of the read buffer used when reading from the transport. If zero, a default (currently 4KB) is used.
var defaultReadBufferSize  int = 4096
// WriteBufferSize specifies the size of the write buffer used when writing to the transport.If zero, a default (currently 4KB) is used.
var defaultWriteBufferSize  int = 4096 
// DisableKeepAlives, if true, disables HTTP keep-alives and will only use the connection to the server for a single HTTP request
var defaultDisableKeepAives bool = false
// DisableCompression, if true, prevents the Transport from requesting compression with an "Accept-Encoding: gzip"
var defaultDisableCompression  bool = false
// Timeout specifies a time limit for requests made by this Client. The timeout includes connection time, any redirects, and reading the response body. The timer remains
var defaultHTTPTimeOut int =  30