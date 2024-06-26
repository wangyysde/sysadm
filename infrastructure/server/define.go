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
 */

package server

var ver="1.0.0"
var author="Wayne Wang<net_use@bzhy.com>"
var serverAddress = "127.0.0.1"
var serverPort = 11050
var serverSocket = "/var/run/infrastructure.sock"    
var serverIsTls = false
var serverCa = ""
var serverCert = ""
var serverKey = ""
var serverInsecureSkipVerify = true
var serverAccessLog = "logs/infrastructure-access.log"
var serverErrorLog = "logs/infrastructure-error.log"
var serverLogKind = "json"
var serverLogLevel = "debug"
var serverLogSplitAccessAndError = true
var serverLogTimeStampFormat = "2006-01-02 15:04:05"
var dbType = "mysql"
var dbName = "infrastructure"
var dbServerAddress = "172.28.1.10"
var dbServerPort = 30306
var dbServerSocket = ""
var dbServerIsTls = false
var dbServerCa = ""
var dbServerCert = ""
var dbServerKey = ""
var dbServerInsecureSkipVerify = true
var dbMaxOpenConns = 10
var dbMaxIdeleConns = 5
var apiServerAddress = "apiserver"
var apiServerPort = 8081
var apiServerIsTls = false
var apiServerCa = ""
var apiServerCert = ""
var apiServerKey = ""
var apiServerInsecureSkipVerify = true
var apiVersion = "v1.0"



