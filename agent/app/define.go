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
// logfile for agent which is used to log runing log messages of agent to 
var LogFile string = "/var/log/sysad-agent.log"
// configuration file path
var CfgFile string = "conf/agent.yaml"
// whether check the certs which got from a server isvalid
var DefaultskipVerifyCert bool = false
//specifies whether agent using TLS protocol when it is communicate with  a server (agent send the reponses message to the server)
var DefaultServerIsTls bool = true
//  the method of getting commands by agent. agent gets commands from the server periodically and run them if this value is true
var DefaultPassive bool = true
// listen port of agent using when agent running as daemon.
var  DefaultListenPort int = 5443
//period(second) for agent gets command from server when agent running as passive
var DefaultPeriod int = 60
// insecret specifies whether agent listen on a insecret port when it is runing as daemon
var DefaultInsecret bool = false
// insecret listen port of agent listening when it is running ad daemon 
var DefaultInsecretPort int = 5080

