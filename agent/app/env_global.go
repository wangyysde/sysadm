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

var env_global map[string]string = map[string]string{
	"IsTls": "GLOBAL_ISTLS",
	"Ca": "GLOBAL_CA",
	"Cert": "GLOBAL_CERT",
	"Key": "GLOBAL_KEY",
	"InsecureSkipVerify": "GLOBAL_INSECURESKIPVERIFY",
	"AccessLog": "GLOBAL_ACCESSLOG",
	"ErrorLog": "GLOBAL_ERRORLOG",
	"Kind": "GLOBAL_LOGKIND",
	"Level": "GLOBAL_LOGLEVEL",
	"TimeFormat": "GLOBAL_LOGTIMEFORMAT",
	"NodeIdentifer": "GLOBAL_NODEIDENTIFER",
	"Uri": "GLOBAL_URI",
	"SourceIP": "GLOBAL_SOURCEIP",
	"commandStatusUri": "GLOBAL_COMMANDSTATUSURI",
	"commandLogsUri": "GLOBAL_COMMANDLOGSURI",
}