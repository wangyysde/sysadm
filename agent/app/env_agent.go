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
*
 */

package app

var env_agent map[string]string = map[string]string{
	"Passive": "AGENT_PASSIVE",
	"IsTls": "AGENT_ISTLS",
	"Ca": "AGENT_CA",
	"Cert": "AGENT_CERT",
	"Key": "AGENT_KEY",
	"InsecureSkipVerify": "AGENT_INSECURESKIPVERIFY",
}

var env_agentRedis map[string]string = map[string]string{
	"IsTls": "AGENT_REDIS_ISTLS",
	"Ca": "AGENT_REDIS_CA",
	"Cert": "AGENT_REDIS_CERT",
	"Key": "AGENT_REDIS_KEY",
	"InsecureSkipVerify": "AGENT_REDIS_INSECURESKIPVERIFY",
	"RedisMode": "AGENT_REDIS_MODE",
	"RedisMaster": "AGENT_REDIS_MASTER",
	"RedisAddrs": "AGENT_REDIS_ADDRS",
}