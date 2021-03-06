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

package config

// Deinfining default value of configuration file
var DefaultConfigFile = "conf/registryctl.yaml"
var SupportVers = [...]string{"v0.1", "v0.2","v21.0.0"}
var SysadmVersion = ""
var ApiVersion = "v1.0"                              // API version of sysadm server
var SysadmHost = "sysadmserver"
var SysadmPort = 8081
var SysadmTls = false
var SysadmInsecureSkipVerify = false 
var RegistryctlVer = "v21.0.0"
var RegistryApiVer = "" 
var DefaultIP = "0.0.0.0"
var DefaultPort = 50001
var DefaultSocket = "/var/run/sysadm.sock"
var DefaultAccessLog = "logs/sysadm-access.log"
var DefaultErrorLog = "logs/sysadm-error.log"
var DefaultLogKind = "text"
var DefaultLogLevel = "debug"
var DefaultUser = "admin"
var DefaultPasswd = "Sysadm12345"
var DefaultDbType =  "postgre"
var DefaultDbHost = "localhost"
var DefaultDbPort = 5432
var DefaultDbUser = "Sysadm"
var DefaultDbPassword = "sysadm12345"
var DefaultDbDbName = "sysadm"
var DefaultDbMaxConnect = 100
var DefaultDbIdleConnect = 20
var DefaultDbSslmode = "disable"
var DefaltDbSslrootcert = ""
var DefaultDbSslkey = ""
var DefaultDbSslcert = ""
var RegistryHost = "sysadm_registry"
var RegistryPort = 5000
var RegistryTls = false
var RegistryInsecureSkipVerify = false 
var RegistryCa = ""
var RegistryKey = ""
var RegistryCert = ""
var RegistryUsername = "sysadm_registry_user"
var RegistryPassword = "sysadm_registry_password"
