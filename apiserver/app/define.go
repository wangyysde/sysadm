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
*
* NOTE:
* defined some functions are related to handle configurations.
 */

package app

import (
	"crypto/x509"
)

const (
	DefaultTlsPort               = 9443
	DefaultPort                  = 9080
	ApiVersion            string = "v1beta1"
	GetCommandUri         string = "getCommand"
	publicKeyAlgorithm           = x509.RSA
	pkiPath                      = "pki"
	caFile                       = "ca.crt"
	caKeyFile                    = "ca.key"
	apiServerCertFile            = "apiserver.crt"
	apiServerFullCertFile        = "apiserver-full.crt"
	apiServerCertKeyFile         = "apiserver.key"

	apiServerCertCommonName = "sysadm-apiserver"
	agentCertCommonName     = "sysadm-agent"
)

var caOrgnaization = []string{"sysadm.cn", "www.sysadm.cn"}
var apiServerCertOrgnaization = []string{"sysadm.cn", "www.sysadm.cn"}
var agentCertOrgnaization = []string{"sysadm.cn", "www.sysadm.cn"}
var caCommonName = "sysadm"
var caPeriodDays = 10 * 365
var apiServerCertPeriodDays = 5 * 365
var agentCertPeriodDays = 365

// current version of apiserver
var appVer string = "1.0.1"

// author of apiserver
var appAuthor string = "Wayne Wang<net_use@bzhy.com>"

// default path of configuration file
var confFilePath string = "conf/apiserver.yaml"

// default address of apiserver listening
var apiserverAddress string = "0.0.0.0"

// default port of apiserver listening
var apiserverPort int = 9443

// default insecret port of apiserver listening
var apiserverInsecretPort int = 8085

// default access log file path
var accessLogFile string = "logs/apiserver-accesslog.log"

// default error log file path
var errorLogFile string = "logs/apiserver-accesslog.log"

// default log kind
var defaultLogKind string = "text"

// default log level
var defaultLogLevel string = "debug"

// default the format of time in the log message
var defaultTimeStampFormat string = "2006-01-02 15:04:05"

// max connection number of apiserver connect to DB server
var defaultMaxDBOpenConns int = 20

// max number of idle connections
var defaultMaxDBIdleConns int = 5

// over time of command execution, second
var defaultMaxExecuteTime int = 3600

// max try time of a command try to execute
var defaultCommandExecuteMaxTryTimes int = 3

// concurrency number of apiserver sending command data to agent when apiserver is running in active mode
var defaultConcurrencySendCommand int = 10

// concurrency number of apiserver get command status from agent when apiserver is running in active mode
var defaultConcurrencyGetCommandStatus int = 10

// concurrency number of apiserver get command log from agent when apiserver is running in active mode
var defaultConcurrencyGetCommandLog int = 10

// Timeout is the maximum amount of time a dial will wait for a connect to complete.
var defaultHttpTimeout int = 180

// KeepAlive specifies the interval between keep-alive probes for an active network connection.
var defaultHttpKeepAliveProbe int = 15

// specifies the maximum amount of time waiting to wait for a TLS handshake. Zero means no timeout.
var defaultTLSHandshakeTimeout int = 5

// IdleConnTimeout is the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself
var defaultIdleConnTimeout int = 60

// 命令的最大重试次数
var defaultMaxCommandTrytimes int = 3

// 日志信息在redis里存储的路径
var defaultLogRootPathInRedis = "/sysadm/apiserver/logs"

// 每次获取命令日志的最大条数
var defaultMaxGetLogNumPerTime = 10
