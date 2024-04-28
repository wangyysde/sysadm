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

package app

import (
	"sysadm/sysadmapi/apiutils"
	"sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

/*
   add handlers for all api version. the path of these handlers for are /api/<version>/syssetting
   this function called in daemonServer
*/
func (s SysSetting)AddHandlers(r *sysadmServer.Engine)(errs []sysadmerror.Sysadmerror){
	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(4020001,"fatal","can not add handlers to nil" ))
		return errs
	}

	for v, h := range apiHandlers{
		e :=  h(r,v, s)
		errs =  append(errs,e...)
	}

	return errs
}


func addHandlersFor1Dot0(r *sysadmServer.Engine, version string ,s SysSetting)(errs []sysadmerror.Sysadmerror){
	moduleName := s.ModuleName
	groupPath := "/api/" + version + "/" + moduleName
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(4020002,"debug","add group handlers for %s", groupPath ))
		
	v1 := r.Group(groupPath)
	{
		v1.GET("/test",testHandler)
	}
	
	return errs
}


func testHandler(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(4030001,"debug","this is test interface"))
	err := apiutils.SendResponseForErrorMessage(c,4030002,"this is test interface")
	errs = append(errs,err...)
	logErrors(errs)
	
}