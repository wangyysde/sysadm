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
* @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html

 */

package server

import (
	"net/http"

	"github.com/wangyysde/sysadmServer"
	"github.com/wangyysde/sysadm/sysadmerror"
)

// add handlers to *sysadmServer.Engine
func addHandlers(r *sysadmServer.Engine)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	e := addTestHandlers(r)
	errs = append(errs,e...)

	return errs

}

// this function add test handlers to R
func addTestHandlers(r *sysadmServer.Engine) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	r.GET("/ping", func(c *sysadmServer.Context) {
        c.String(http.StatusOK, "echo ping message")
    })  

	return errs
}