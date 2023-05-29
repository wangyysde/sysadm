/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
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
* Note: this file holds main functions for start subcommand
 */

package server

import (
	//"context"
	//"fmt"
	//"os"

	"github.com/spf13/cobra"
	"os"

	"github.com/wangyysde/sysadmServer"
	//"sysadm/redis"
	apiserverApp "sysadm/apiserver/app"
	"sysadm/sysadmerror"
)

var exitChan chan os.Signal

func Start(cmd *cobra.Command, args []string) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20010001, "debug", "starting  apiserver....."))
	ok, err := apiserverApp.HandlerConfig()
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		os.Exit(1)
	}

	// initating loggers
	ok, err = apiserverApp.SetLogger()
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		os.Exit(2)
	}
	defer apiserverApp.CloseLogger()
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080003, "debug", "loggers have been set"))
	logErrors(errs)
	errs = errs[0:0]

	// initating redis entity
	ok, err = apiserverApp.InitRedis()
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		os.Exit(3)
	}
	defer apiserverApp.CloseRedisEntity()

	// initating DB entity
	ok, err = apiserverApp.InitDBEntity()
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		os.Exit(4)
	}
	defer apiserverApp.CloseDBEntity()

	exitChan = make(chan os.Signal)
	if apiserverApp.IsPassiveMode() {
		ok, err := apiserverApp.RunDaemonPassive()
		errs = append(errs, err...)
		if !ok {
			logErrors(errs)
			os.Exit(5)
		}
	}
	/*


		if RunConf.Agent.Passive {
			err = run_DaemonPassive()
			logErrors(err)
			if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("fatal") {
				os.Exit(4)
			}
			os.Exit(0)
		}

		// initating server
		r := sysadmServer.New()
		r.Use(sysadmServer.Logger(), sysadmServer.Recovery())

		err = addHandlers(r)
		errs = append(errs, err...)
		if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("fatal") {
			logErrors(errs)
			os.Exit(3)
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080004, "debug", "handlers have be added."))

		if !RunConf.Agent.Insecret {
			tlsStr := fmt.Sprintf("%s:%d", RunConf.Agent.Server.Address, RunConf.Agent.Server.Port)
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080005, "debug", "listening TLS service to %s", tlsStr))
			logErrors(errs)
			errs = errs[0:0]
			e := r.RunTLS(tlsStr, RunConf.Agent.Tls.Cert, RunConf.Agent.Tls.Key)
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080006, "fatal", "can not listen TLS service. error %s", e))
				logErrors(errs)
				os.Exit(4)
			}

		} else {
			if RunConf.Agent.Tls.IsTls {
				go func() {
					tlsStr := fmt.Sprintf("%s:%d", RunConf.Agent.Server.Address, RunConf.Agent.Server.Port)
					errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080007, "debug", "listening TLS service to %s", tlsStr))
					logErrors(errs)
					errs = errs[0:0]
					e := r.RunTLS(tlsStr, RunConf.Agent.Tls.Cert, RunConf.Agent.Tls.Key)
					if e != nil {
						errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080008, "error", "can not listen TLS service. error %s", e))
						logErrors(errs)
					}
				}()
			}

			listenStr := fmt.Sprintf("%s:%d", RunConf.Agent.Server.Address, RunConf.Agent.InsecretPort)
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080009, "debug", "listening  service to %s", listenStr))
			logErrors(errs)
			errs = errs[0:0]
			e := r.Run(listenStr)
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080010, "error", "can not listent service. error %s", e))
				logErrors(errs)
			}
		}
	*/
	logErrors(errs)
}

// log log messages to logfile or stdout
func logErrors(errs []sysadmerror.Sysadmerror) {

	for _, e := range errs {
		l := sysadmerror.GetErrorLevelString(e)
		no := e.ErrorNo
		sysadmServer.Logf(l, "erroCode: %d Msg: %s", no, e.ErrorMsg)
	}
}
