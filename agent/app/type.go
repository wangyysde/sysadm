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

import (
	"os"
	"net/http"

	"github.com/wangyysde/sysadm/config"
	"github.com/wangyysde/sysadm/httpclient"
)

type RunConfig struct {
	// workingDir is X/../ ,X is the path of the directory which binary package of agent locate in it. 
	WorkingDir string
	Version config.Version
	// for agnet configuration file path
	CfgFile string 
	// descriptor of  log file which will be used to close logger when system exit
	LogFileFp *os.File
	// the following parameters can be set through command flags(persistent flags), configuration file (global section) or ENV variables
	Global GlobalConf `form:"global" json:"global" yaml:"global" xml:"global"`
	// the following parameters can be set through command flags(daemon subcommand flags), configuration file (agent section) or ENV variables
	// this struct is for daemon subcommand and agent block in configuration file
	Agent AgentConf `form:"agent" json:"agent" yaml:"agent" xml:"agent"`
}

// this for command flags and options
type CliOptions struct {
	Version config.Version
	// for agnet configuration file path
	CfgFile string 
	// the following parameters can be set through command flags(persistent flags), configuration file (global section) or ENV variables
	Global GlobalConf `form:"global" json:"global" yaml:"global" xml:"global"`
	// the following parameters can be set through command flags(daemon subcommand flags), configuration file (agent section) or ENV variables
	Agent AgentConf `form:"agent" json:"agent" yaml:"agent" xml:"agent"`
}

// this for configuration file content 
type FileConf struct {
	// the following parameters can be set through command flags(persistent flags), configuration file (global section) or ENV variables
	Global GlobalConf `form:"global" json:"global" yaml:"global" xml:"global"`
	// the following parameters can be set through command flags(daemon subcommand flags), configuration file (agent section) or ENV variables
	Agent AgentConf `form:"agent" json:"agent" yaml:"agent" xml:"agent"`
}

// the following parameters can be set through command flags(persistent flags), configuration file (global section) or ENV variables
type GlobalConf struct {
	// tls parameters which agent will used to connect to a server(agent send the reponses message to the server)
	Tls config.Tls `form:"tls" json:"tls" yaml:"tls" xml:"tls"`

	// server parameters  which agent will used to connect to a server(agent send the reponses message to the server)
	Server config.Server `form:"server" json:"server" yaml:"server" xml:"server"`

	// where the results of a command running will be send to. one of server: a server receiving the results; stdout, file
	Output string `form:"output" json:"output" yaml:"output" xml:"output"`

	// the path of output file. this value must not empty if output be set to "file"
	OutputFile string `form:"outputFile" json:"outputFile" yaml:"outputFile" xml:"outputFile"`

	// log setting block
	Log config.Log `form:"log" json:"log" yaml:"log" xml:"log"`

	// set whether agent running in debug mode
	DebugMode bool `form:"debug" json:"debug" yaml:"debug" xml:"debug"`

	// specifies a identifer of the node which agent running on it.
	// It is any combination of the IP,HOSTNAME and MAC joined by commas  or a customize string what the leght of the string is less 63
	// agent will get all IPs without not active and reponse these IPs in list to the server by nodeIdentifer.IPs filed if IP is included in NodeIdentifer
	// agent will get hostname and reponse the hostname  to the server by nodeIdentifer.Hostname filed if hostname is included in NodeIdentifer
	// agent will get all MACs without not active and reponse these MACs in list to the server by nodeIdentifer.MACs filed if MAC is included in NodeIdentifer
	// customize string is reponse to the server directly .
	// customize string is conflicted with IP,HOSTNAME and MAC. the nodeIdentifer can be changed by the server during agent communicate with the server
	NodeIdentifer string `form:"nodeIdentifer" json:"nodeIdentifer" yaml:"nodeIdentifer" xml:"nodeIdentifer"`

	// specifies the uri where agent get commands to run when agent runing as daemon in passive mode. 
	// agent will send the requests to "/" on the server if GetUri is empty.
	// Uri is the path where agent will send result message to when is running as command.
	// Uri is the listen path where agent receives commands to run when  agent runing as daemon in active mode. 
	// the length of this value shoule less 63
	Uri string `form:"uri" json:"uri" yaml:"uri" xml:"uri"`

	// sourceIP specifies the source IP address which will be use to connect to a server by agent. this ip address must be configurated on one of the 
	// interfaces  on the host where agent running on.  agent will get a source IP address from host  automatically if the value of this field is "". 
	SourceIP string `form:"sourceIP" json:"sourceIP" yaml:"sourceIP" xml:"sourceIP"`
}

// the following parameters can be set through command flags(daemon subcommand flags), configuration file (agent section) or ENV variables
// this struct is for daemon subcommand and agent block in configuration file
type AgentConf struct {
	// the method of getting commands by agent. agent gets commands from the server periodically and run them if this value is true
	// otherwise the server send a command to a agent when it  want the agent to run the command on a host. 
	Passive bool `form:"passive" json:"passive" yaml:"passive" xml:"passive"`

	// tls parameters for agent when agent running as daemon.
	Tls config.Tls `form:"tls" json:"tls" yaml:"tls" xml:"tls"`

	// listen parameters   of agent using when agent running as daemon.
	Server config.Server `form:"server" json:"server" yaml:"server" xml:"server"`

	// period(second) for agent gets command from server when agent running as passive
	Period int `form:"period" json:"period" yaml:"period" xml:"period"`

	// insecret specifies whether agent listen on a insecret port when it is runing as daemon
	Insecret bool `form:"insecret" json:"insecret" yaml:"insecret" xml:"insecret"`

	// insecret listen port of agent listening when it is running ad daemon 
	InsecretPort int `form:"insecretPort" json:"insecretPort" yaml:"insecretPort" xml:"insecretPort"`
}

// nodeIdentifer is used to save node identifer information what agent will send to server
// customize string is conflicted with IP,HOSTNAME and MAC. the nodeIdentifer can be changed by the server during agent communicate with the server
// customize is the higher priority than others
type NodeIdentifer struct {
	Ips []string
	Macs []string
	Hostname string 
	Customize string
}

var RunConf RunConfig = RunConfig{
	WorkingDir: "",
	Version: config.Version{},
	CfgFile: "",
	Global: GlobalConf{},
	Agent: AgentConf{},
}
var CliOps CliOptions = CliOptions{
	Version: config.Version{},
	CfgFile: "",
	Global: GlobalConf{},
	Agent: AgentConf{},
}

// runTimeData used to record the data what are used frequently in runtime.
type runTimeData struct {
	// In passive mode server can tell agent change nodeIdentifer. this field used to record the node identifer what has be built by getNodeIdentifer func
	nodeIdentifer *NodeIdentifer 

	// the complete url address where agent get a command to execute from server. this url address can be changed when the server ask agent to do so.
	getCommandUrl string

	// getCommandParames to save the parameters of the request of get command.
	getCommandParames *httpclient.RequestParams

	// keep http or https client for reuse. we should recreate http client if the value of this field is nil
	httpClient *http.Client 
}

/* 
	Command is used to save command data received from or got from server.
*/
type Command struct {
	// command name , agent will route handler according to the value of this filed.
	Command string `form:"command" json:"command" yaml:"command" xml:"command"`
	// the value of this feild should not empty if server want to specify the node identifer.
	// specifies a identifer of the node which agent running on it.
	// It is any combination of the IP,HOSTNAME and MAC joined by commas  or a customize string what the leght of the string is less 63
	// agent will get all IPs without not active and reponse these IPs in list to the server by nodeIdentifer.IPs filed if IP is included in NodeIdentifer
	// agent will get hostname and reponse the hostname  to the server by nodeIdentifer.Hostname filed if hostname is included in NodeIdentifer
	// agent will get all MACs without not active and reponse these MACs in list to the server by nodeIdentifer.MACs filed if MAC is included in NodeIdentifer
	// customize string is reponse to the server directly .
	// customize string is conflicted with IP,HOSTNAME and MAC. the nodeIdentifer can be changed by the server during agent communicate with the server
	NodeIdentifer string `form:"nodeIdentifer" json:"nodeIdentifer" yaml:"nodeIdentifer" xml:"nodeIdentifer"`
	// where agent will reponse the result of the command execution,otherwise agent will reponse the result of the command execution to the url where
	// it got the command.
	ResponseUri string `form:"responseUri" json:"responseUri" yaml:"responseUri" xml:"responseUri"`
	// parameters what will be used to execute the command. key is parameter name and value is the value of the parameter.
	Parameters map[string]string `form:"parameter" json:"parameter" yaml:"parameter" xml:"parameter"`

}
var runData runTimeData = runTimeData{}