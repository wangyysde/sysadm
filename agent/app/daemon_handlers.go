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
	"strings"

	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

func addHandlers(r *sysadmServer.Engine) (errs []sysadmerror.Sysadmerror) {

	if e := addRootHandler(r); e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803001, "fatal", "add root handler error %s", e))
		return errs
	}

	if e := addReceiveCommandHandler(r); e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803002, "fatal", "add receive command  handler error %s", e))
		return errs
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

	return nil
}

func receivedCommand(c *sysadmServer.Context) {
	var cmd Command = Command{}
	var errs []sysadmerror.Sysadmerror

	err := c.BindJSON(&cmd)
	if err != nil {
		msg := fmt.Sprintf("receive command error %s", err)
		_ = apiutils.SendResponseForErrorMessage(c, 10082001, msg)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082001, "error", msg))
		logErrors(errs)
		return
	}

	doRouteCommand(&cmd, c)
}
