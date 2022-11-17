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
	"github.com/wangyysde/sysadm/config"
)

type RunConfig struct {
	Version config.Version

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

	// logfile for agent which is used to log runing log messages of agent to 
	LogFile  string `form:"logFile" json:"logFile" yaml:"logFile" xml:"logFile"`
}

// the following parameters can be set through command flags(daemon subcommand flags), configuration file (agent section) or ENV variables
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

var RunConf RunConfig = RunConfig{}
var CliOps CliOptions = CliOptions{
	Version: config.Version{},
	CfgFile: "",
	Global: GlobalConf{},
	Agent: AgentConf{},
}