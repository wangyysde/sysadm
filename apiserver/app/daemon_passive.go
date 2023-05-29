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
* Note: this file holds functions of daemon for passive mode
 */

package app

import (
	"fmt"
	"github.com/wangyysde/sysadmServer"
	"sysadm/sysadmerror"
)

// RunDaemonPassive is the main function for passive mode
func RunDaemonPassive() (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	r := sysadmServer.New()
	r.Use(sysadmServer.Logger(), sysadmServer.Recovery())

	ok, err := addHandlers(r)
	errs = append(errs, err...)
	if !ok {
		return false, errs
	}

	// try to listen insecret port
	if runData.runConf.ConfServer.InsecretPort != 0 {
		go func() {
			listenStr := fmt.Sprintf("%s:%d", runData.runConf.ConfServer.Address, runData.runConf.ConfServer.InsecretPort)
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030001, "debug", "listening  service to %s", listenStr))
			logErrors(errs)
			errs = errs[0:0]
			e := r.Run(listenStr)
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030002, "error", "can not listent service. error %s", e))
				logErrors(errs)
			}
		}()
	}

	// try to listen TLS
	if runData.runConf.ConfServer.IsTls {
		tlsStr := fmt.Sprintf("%s:%d", runData.runConf.ConfServer.Address, runData.runConf.ConfServer.Port)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030002, "debug", "listening TLS service to %s", tlsStr))
		logErrors(errs)
		errs = errs[0:0]
		e := r.RunTLS(tlsStr, runData.runConf.ConfServer.Cert, runData.runConf.ConfServer.Key)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030003, "fatal", "can not listen TLS service. error %s", e))
			return false, errs
		}

	} else {
		tlsStr := fmt.Sprintf("%s:%d", runData.runConf.ConfServer.Address, runData.runConf.ConfServer.Port)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030004, "debug", "listening HTTP service to %s", tlsStr))
		logErrors(errs)
		errs = errs[0:0]
		e := r.Run(tlsStr)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030005, "error", "can not listent service. error %s", e))
			logErrors(errs)
		}
	}

	return true, errs
}

// addHandlers add handlers for getCommand,receiveCommandStatus and receiveLogs when apiserver running in passive mode
func addHandlers(r *sysadmServer.Engine) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	if e := addRootHandler(r); e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030006, "fatal", "add root handler error %s", e))
		return false, errs
	}

	// add a handler for client get command to run when apiserver is running in passive mode
	getCommandUri := runData.runConf.ConfGlobal.CommandUri
	r.POST(getCommandUri, getCommand)

	// add a handler for receiving command status what were sent by  client
	r.POST(runData.runConf.ConfGlobal.CommandStatusUri, receiveCommandStatus)

	// add a handler for receiving command logs what were sent by  client
	r.POST(runData.runConf.ConfGlobal.CommandLogsUri, receiveCommandLogs)

	return true, errs
}

// addRootHandler adding handler for root path
func addRootHandler(r *sysadmServer.Engine) error {
	if r == nil {
		return fmt.Errorf("router is nil")
	}

	r.Any("/", func(c *sysadmServer.Context) {
		c.JSON(200, sysadmServer.H{
			"status": "ok",
		})
	})

	return nil
}

func getCommand(c *sysadmServer.Context) {
	// TODO

}

func receiveCommandStatus(c *sysadmServer.Context) {
	// TODO

}

func receiveCommandLogs(c *sysadmServer.Context) {
	// TODO

}