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

import (
	"fmt"
	"net/http"
	"strings"

	apiserver "github.com/wangyysde/sysadm/apiserver/app"
	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

func addHandlers(r *sysadmServer.Engine) (errs []sysadmerror.Sysadmerror) {

	if e := addRootHandler(r); e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803001, "fatal", "add root handler error %s", e))
		return errs
	}

	// we should build nodeIdentifier after start listen
	nodeIdentifier, e := apiserver.BuildNodeIdentifer(RunConf.Global.NodeIdentifer)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803001, "fatal", "build node identifier error %s", e))
		return errs
	}
	runData.nodeIdentifer = &nodeIdentifier

	
	if !RunConf.Agent.Passive {
		// add handler for the path of uri specifing if agent running in active mode 
		if e := addReceiveCommandHandler(r); e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803002, "fatal", "add receive command  handler error %s", e))
			return errs
		}

		if e := addGetCommandStatusHandler(r); e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803003, "fatal", "add get command status  handler error %s", e))
			return errs
		}

		if e := addGetLogs(r); e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803003, "fatal", "add get logs  handler error %s", e))
			return errs
		}
	}
	

	return errs
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

func addReceiveCommandHandler(r *sysadmServer.Engine) error {
	if strings.TrimSpace(RunConf.Global.Uri) == "" {
		RunConf.Global.Uri = defaultReceiveCommandUri
	}

	listenUri := RunConf.Global.Uri
	if listenUri[0:1] != "/" {
		listenUri = "/" + listenUri
	}

	r.POST(listenUri, receivedCommand)
	r.GET(listenUri, receivedCommand)

	return nil
}

func receivedCommand(c *sysadmServer.Context) {
	var cmd apiserver.CommandData = apiserver.CommandData{}
	var errs []sysadmerror.Sysadmerror

	err := c.BindJSON(&cmd)
	if err != nil {
		data := make(map[string]interface{},0)
		msg := fmt.Sprintf("receive command error %s", err)
		commandStatus, e := apiserver.BuildCommandStatus("",RunConf.Global.NodeIdentifer,msg, *runData.nodeIdentifer,apiserver.ComandStatusSendError,data,true)
		c.JSON(http.StatusOK,commandStatus)

		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082001, "error", msg))
		errs = append(errs,sysadmerror.NewErrorWithStringLevel(10082001, "error",fmt.Sprintf("build command status error: %s",e)))
		logErrors(errs)
		return
	}

	doRouteCommand(&cmd, c)
}


func addGetCommandStatusHandler(r *sysadmServer.Engine) error {
	if strings.TrimSpace(RunConf.Global.CommandStatusUri) == "" {
		RunConf.Global.CommandStatusUri = defaultGetCommandStatus
	}

	listenUri := RunConf.Global.CommandStatusUri
	if listenUri[0:1] != "/" {
		listenUri = "/" + listenUri
	}

	r.POST(listenUri, getCommandStatus)
	r.GET(listenUri, getCommandStatus)

	return nil
}

func getCommandStatus(c *sysadmServer.Context) {
	var cmdStatusReq apiserver.CommandStatusReq = apiserver.CommandStatusReq{}
	var errs []sysadmerror.Sysadmerror

	err := c.BindJSON(&cmdStatusReq)
	if err != nil {
		msg := fmt.Sprintf("the request for getting command status is not valid %s", err)
		_ = apiutils.SendResponseForErrorMessage(c, 10082001, msg)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082001, "error", msg))
		logErrors(errs)
		return
	}

	doRouteCommand(&cmdStatus, c)
}