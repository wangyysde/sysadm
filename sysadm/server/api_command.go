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
 */

/*
* NOTE:
* This file include the functions what are  command sending acrcording client requested, command status receiving and command logs receiving
 */

package server

import (
	"strings"
	"net/http"

	apiServerApp "github.com/wangyysde/sysadm/apiserver/app"
	"github.com/wangyysde/sysadm/httpclient"
	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
	infrastructure "github.com/wangyysde/sysadm/infrastructure/app"
)

// adding command sending, command status receiving and command logs receiving  handlers
func addCommandHandlers(r *sysadmServer.Engine) []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror

	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030001, "fatal", "can not add handlers to nil"))
		return errs
	}

	// for client get commands to run when running at passive mode
	r.POST("getCommand", getCommandHandler)

	// for client send command running status when running at passive mode
	r.POST("receiveCommandStatus", receiveCommandStatusHandler)

	// for client set command running logs to the server when running at passive mode
	r.POST("receiveLogs", receiveLogsHandler)

	return errs
}

/*
handler for handling list of the project
Query parameters of request are below:
conditionKey: key name for DB query ,such as projectid, ownerid,name....
conditionValue: the value of the conditionKey.for projectid, ownereid using =, for name, comment using like.
deleted: 0 :normarl 1: deleted
start: start number of the result will be returned.
num: lines of the result will be returned.
*/
func getCommandHandler(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030001, "debug", "received a get command request"))

	body, err := httpclient.GetRequestBody(c.Request)
	errs = append(errs, err...)
	if len(body) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030002, "error", "parameters error"))
		err := apiutils.SendResponseForErrorMessage(c, 700030003, "parameters error")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	commandReq, e := apiServerApp.UnmarshalCommandReq(body)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030004, "error", "request command data is not valid"))
		err := apiutils.SendResponseForErrorMessage(c, 700030004, "request command data is not valid")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	if len(commandReq.Ips) < 1 && len(commandReq.Macs) < 1 && strings.TrimSpace(commandReq.Hostname) == "" && strings.TrimSpace(commandReq.Customize) == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700030005, "error", "all node identifier fields are empty"))
		err := apiutils.SendResponseForErrorMessage(c, 700030005, "all node identifier fields are empty")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	command,err := infrastructure.GetCommand(commandReq.Ips, commandReq.Macs, commandReq.Hostname, commandReq.Customize)
	errs = append(errs,err...)
	c.JSON(http.StatusOK, command)
	logErrors(errs)
}

/*
handler for handling list of the project
Query parameters of request are below:
conditionKey: key name for DB query ,such as projectid, ownerid,name....
conditionValue: the value of the conditionKey.for projectid, ownereid using =, for name, comment using like.
deleted: 0 :normarl 1: deleted
start: start number of the result will be returned.
num: lines of the result will be returned.
*/
func receiveCommandStatusHandler(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101004, "debug", "now handling project list"))
	// TODO
}

/*
handler for handling list of the project
Query parameters of request are below:
conditionKey: key name for DB query ,such as projectid, ownerid,name....
conditionValue: the value of the conditionKey.for projectid, ownereid using =, for name, comment using like.
deleted: 0 :normarl 1: deleted
start: start number of the result will be returned.
num: lines of the result will be returned.
*/
func receiveLogsHandler(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1101004, "debug", "now handling project list"))
	// TODO
}


