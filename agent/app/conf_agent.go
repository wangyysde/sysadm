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
	"strconv"
	"strings"

	"github.com/wangyysde/sysadm/sysadmerror"
)

/*
 handle configuration items set in agent block.
 these configuration items have be set by flags in  daemon subcommand or in agent block in configuration file
*/
func handleAgentBlock()([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082001,"debug","try to handle configuration items in agent block")) 

	var envAgentMap map[string]string 
	envMapP := getEnvDefineForBlock("agent")
	if envMapP == nil {
		envAgentMap = map[string]string{}
	} else {
		envAgentMap = *envMapP
	}

	passiveValue, err := validatePassive(envAgentMap)
	errs = append(errs, err...)
	RunConf.Agent.Passive = passiveValue

	tlsValue := validateTlsConf(CliOps.Agent.Tls,fileConf.Agent.Tls,"agent")
	RunConf.Agent.Tls = *tlsValue

	serverValue, err := validateServerConf(CliOps.Agent.Server,fileConf.Agent.Server,"agent",false)
	RunConf.Agent.Server = *serverValue
	errs = append(errs,err...)

	periodValue, err := validatePeriod(envAgentMap)
	errs = append(errs,err...)
	RunConf.Agent.Period = periodValue

	insecretValue, err := validateInsecret(envAgentMap)
	errs = append(errs, err...)
	RunConf.Agent.Insecret  = insecretValue

	insecretPortValue, err := validateInsecretPort(envAgentMap)
	errs = append(errs,err...)
	RunConf.Agent.InsecretPort  = insecretPortValue

	// agent gets commands from the server periodically and run them if RunConf.Agent.Passive is true
	if RunConf.Agent.Passive {
		// Period must greater than zero if agent running in passive mode
		if RunConf.Agent.Period == 0 {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082002,"warning","agent running as passive mode, but the check period has be set to zero. that is not valid, default value(%d) will be set",DefaultPeriod)) 
			RunConf.Agent.Period = DefaultPeriod
		}

		return errs
	}
	
	// agent running as active mode
	if RunConf.Agent.Tls.IsTls {
		// IsTls is true,but one of certs file not set
		if  RunConf.Agent.Tls.Cert == "" ||  RunConf.Agent.Tls.Key == "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082003,"warning","agent has be set to run with SSL,but the certification files is not valid,we try to disable TLS. ca: %s cert: %s key: %s",RunConf.Agent.Tls.Ca, RunConf.Agent.Tls.Cert,RunConf.Agent.Tls.Key)) 
			RunConf.Agent.Tls.IsTls = false
			if ! RunConf.Agent.Insecret {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082008,"warning","agent will be running without TLS,because of the certification files is not valid. ")) 
			}
			RunConf.Agent.Insecret = true
		}
	} 
	
	if RunConf.Agent.Tls.IsTls {
		if RunConf.Agent.Server.Port == 0 {   // agent's ssl port must be set when agent will be runing with SSL.
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082009,"fatal","agent's ssl port must be set when agent will be runing with SSL.")) 
			return errs
		}

		if RunConf.Agent.Server.Port ==  RunConf.Agent.InsecretPort && RunConf.Agent.Insecret {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082010,"warning","agent's insecret listen port can not same to SSL prot. insecret listen will be disabled"))
			RunConf.Agent.Insecret = false 
		}
	} else {
		if ! RunConf.Agent.Insecret {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082011,"warning","we can not disable both ssl mode and insecret.")) 
			RunConf.Agent.Insecret = true
		}
	}

	if RunConf.Agent.Insecret {
		if RunConf.Agent.InsecretPort == 0 {
			if RunConf.Agent.Server.Port != 0 {
				if ! RunConf.Agent.Tls.IsTls {
					errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082012,"warning","insecret port is not valid. the port setted in server block will be used to listened in insecret mode."))
			 		RunConf.Agent.InsecretPort = RunConf.Agent.Server.Port
				} else {
					errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082013,"warning","agent can not listen in insecret with invalid port"))
					RunConf.Agent.Insecret = false
				}
			} else {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082014,"fatal","agent can not listen in insecret with invalid port"))
				return errs
			}
		}
	}

	if ! RunConf.Agent.Insecret && ! RunConf.Agent.Tls.IsTls {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082015,"fatal","we can not disable both insecret and ssl."))
		return errs
	}

	return errs

}

/*
	validatePassive validate the passive values in cliConf (set by command line flags), fileConf(set by configuration file) and envAgentMap (set by environment)
	the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher 
	the priority of the defination in configuration file.
*/
func validatePassive( envAgentMap map[string]string)(ret bool, errs []sysadmerror.Sysadmerror){
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082101,"debug","try to validate passive value"))
	ret = false

	passiveName, okPassive := envAgentMap["Passive"]
	if okPassive {
		passiveValue := os.Getenv(passiveName)
		passiveValue = strings.ToLower(strings.TrimSpace(passiveValue))
		// because the default value of passive is false, so we should to check whether the value has be set to true.
		if passiveValue == "y" || passiveValue == "yes" || passiveValue == "on" || passiveValue == "1" {
			ret = true
		}
	}

	filePassive := fileConf.Agent.Passive
	// the default value is false. so set the value of ret to true if filePassive is not false,otherwise skip.
	if filePassive {
		ret = true
	}

	cliPassive := CliOps.Agent.Passive
	if cliPassive {
		ret =true
	}

	return ret, errs
}

/*
	validatePeriod validate the period values in cliConf (set by command line flags), fileConf(set by configuration file) and envAgentMap (set by environment)
	the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher 
	the priority of the defination in configuration file.
*/
func validatePeriod( envAgentMap map[string]string)(ret int, errs []sysadmerror.Sysadmerror){
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082201,"debug","try to validate period value"))
	ret = 0

	periodName, okPeriod := envAgentMap["Period"]
	if okPeriod {
		periodValue := os.Getenv(periodName)
		if strings.TrimSpace(periodValue) != "" {
			periodInt, e := strconv.Atoi(strings.TrimSpace(periodValue))
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082202,"debug","get period value from environment with name %s ,but the value is not valid.",periodName))
			} else {
				ret = periodInt
			}
		}
		
	}

	if fileConf.Agent.Period != 0 {
		ret = fileConf.Agent.Period
	}

	if CliOps.Agent.Period != 0 && CliOps.Agent.Period != DefaultPeriod {
		ret = CliOps.Agent.Period
	}

	return ret,errs 
}

/*
	validateInsecret validate the insecret values in cliConf (set by command line flags), fileConf(set by configuration file) and envAgentMap (set by environment)
	the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher 
	the priority of the defination in configuration file.
*/
func validateInsecret( envAgentMap map[string]string)(ret bool, errs []sysadmerror.Sysadmerror){
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082301,"debug","try to validate insecret value"))
	ret = true

	insecretName, okInsecret := envAgentMap["Insecret"]
	if okInsecret {
		insecretValue := os.Getenv(insecretName)
		insecretValue = strings.ToLower(strings.TrimSpace(insecretValue))
		// because the default value of insecret is false, so we should to check whether the value has be set to true.
		if insecretValue == "y" || insecretValue == "yes" || insecretValue == "on" || insecretValue == "1" {
			ret = true
		}
	}

	fileInsecret := fileConf.Agent.Insecret
	// the default value is false. so set the value of ret to true if file Insecret is not false,otherwise skip.
	if fileInsecret {
		ret = true
	}

	cliInsecret := CliOps.Agent.Insecret
	if cliInsecret {
		ret =true
	}

	return ret, errs
}

/*
	validateInsecretPort validate the Insecret Port  values in cliConf (set by command line flags), fileConf(set by configuration file) and envAgentMap (set by environment)the priority of the defination in configuration file is higher than enverionments. and  the priority of the defination in command line flags is higher the priority of the defination in configuration file.
*/
func validateInsecretPort( envAgentMap map[string]string)(ret int, errs []sysadmerror.Sysadmerror){
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082401,"debug","try to validate Insecret Port value"))
	ret = 0

	insecretPortName, okInsecretPort := envAgentMap["InsecretPort"]
	if okInsecretPort {
		insecretPortValue := os.Getenv(insecretPortName)
		if strings.TrimSpace(insecretPortValue) != "" {
			insecretPortInt, e := strconv.Atoi(strings.TrimSpace(insecretPortValue))
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082402,"debug","get Insecret Port  value from environment with name %s ,but the value is not valid.",insecretPortName))
			} else {
				ret = insecretPortInt
			}
		}
		
	}

	if fileConf.Agent.InsecretPort  != 0 {
		ret = fileConf.Agent.InsecretPort 
	}

	if CliOps.Agent.InsecretPort != 0 {
		ret = CliOps.Agent.InsecretPort
	}

	return ret,errs 
}