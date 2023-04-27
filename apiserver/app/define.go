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

// current version of apiserver 
var appVer string = "1.0.1"

// author of apiserver
var appAuthor string = "Wayne Wang<net_use@bzhy.com>"

// default path of configuration file
var confFilePath string = "conf/apiserver.yaml"

// specifies the default uri where client get commands to run when apiserver runing as daemon in passive mode
var passiveResponseCommandUri string = "/getCommand"

// specifies the default uri for apiserver to send command data to client when apiserver is run in active mode
var activeSendCommandUri string = "/receiveCommand"

// specifies the default uri where client send command status to when apiserver running as daemon in passive mode
var passiveResponseCommandStatusUri string = "/receiveCommandStatus"

// specifies the default uri for apiserver to send command status to client when apiserver is run in active mode
var activeSendCommandStatusUri string = "/getCommandStatus"

// specifies the default uri where client send command logs to when apiserver running as daemon in passive mode
var passiveResponseCommandLogsUri string = "/receiveLogs"

// specifies the default uri for apiserver to send command status to client when apiserver is run in active mode
var activeSendCommandLogsUri string = "/getLogs"

// interval of checking new command for client by apiserver when apiserver is running actively. default is 5 second.
var checkCommandInterval int = 5

// interval of try to get command status from client by apiserver when apiserver is running actively. default is 5 second
var getStatusInterval int = 5

// interval of try to get command log from client by apiserver when apiserver is running actively. default is 5 second
var getLogInterval int = 5

// default address of apiserver listening
var apiserverAddress string = "0.0.0.0"

// default port of apiserver listening
var apiserverPort int = 8085

// default insecret port of apiserver listening
var apiserverInsecretPort int = 5085

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
