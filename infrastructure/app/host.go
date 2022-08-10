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
	"fmt"
	"regexp"
	"strings"

	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
	"github.com/wangyysde/sysadmServer"
)

func addHost(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror

	var requestData ApiHost
	if e := c.ShouldBind(&requestData); e != nil {
		msg := fmt.Sprintf("get host data err %s", e) 
		_ = apiutils.SendResponseForErrorMessage(c,3030001,msg)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3030001,"error",msg))
		logErrors(errs)
		return 
	}
	
	if e,err := checkAddHostData(&requestData); e != nil {
		errs=append(errs,err...)
		msg := fmt.Sprintf("%s", e) 
		err = apiutils.SendResponseForErrorMessage(c,3020012,msg)
		errs = append(errs,err...)
		logErrors(errs)
		return 

	}

	
	if requestData.Port == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3030002,"warn","can not get SSH Service Port. Default port will be used"))
		requestData.Port = 22
	}

	if strings.TrimSpace(requestData.User) == "" || strings.TrimSpace(requestData.Password) == "" {
		msg := "user account and password must be set" 
		_ = apiutils.SendResponseForErrorMessage(c,3030003,msg)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020005,"error",msg))
		return
	}

	msg := "host has be added successful"
	err := apiutils.SendResponseForSuccessMessage(c,msg)	
	errs=append(errs,err...)
	logErrors(errs)
}

/*
	validating the valid of the data what have received on http context
	return nil and sysadmerror.Sysadmerror if the data validated.
	otherwise return error and  sysadmerror.Sysadmerror
*/
func checkAddHostData(data *ApiHost)(error,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	// hostname should be large 1 and less 64
	hostname := strings.TrimSpace(data.Hostname) 
	if(len(hostname) < 1 || len(hostname) >63 ){
		errMsg := fmt.Sprintf("hostname %s is not valid.",hostname)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020006,"error",errMsg))
		return  fmt.Errorf(errMsg), errs
	}
	data.Hostname = hostname

	// checking ip
	ip := strings.TrimSpace(data.Ip)
	if tmpIP,e := utils.CheckIpAddress(ip,false); tmpIP == nil {
		return fmt.Errorf("host ip %s is not valid.",ip),e
	}
	data.Ip = ip

	ipType := strings.TrimSpace(data.Iptype)
	if ipType != "4" && ipType != "6"{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020007,"warn","ip type (%s) is not valid. we try to look node IP as ipv4.",ipType))
		ipType = "4"
	}
	data.Iptype = ipType

	port := data.Port
	if tmpPort,e := utils.CheckPort(port); tmpPort == 0{
		errs = append(errs, e...)
		return fmt.Errorf("SSH port(%d) is not valid",port),errs
	}

	user := strings.TrimSpace(data.User)
	if matched,e := regexp.MatchString(`^[a-zA-Z0-9_.][a-zA-Z0-9_.-]{1,31}$`, user); !matched{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020008,"error","os user (%s) is not valid %s.",user,e))
		return fmt.Errorf("os user (%s) is not valid.",user),errs
	}
	data.User =  user

	pubkeyUploaded := data.PubkeyUploaded
	if !pubkeyUploaded {
		password := strings.TrimSpace(data.Password)
		rePassword := strings.TrimSpace(data.RePassword)
		if password != rePassword {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020009,"error","the password(%s) is not matched twice password(%s).",password,rePassword))
			return fmt.Errorf("the password(%s) is not matched twice password(%s).",password,rePassword), errs
		}

		if matched,e := regexp.MatchString(`^.{1,64}$`, password); !matched{
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020010,"error","user password (%s) is not valid %s.",password,e))
			return fmt.Errorf("user password(%s) is not valid.",password),errs
		}

		data.Password = password
	}

	osID := data.OsID
	osVersionID := data.OsVersionID

	if osID == 0 || osVersionID == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020011,"error","os and version of the os for node should be selected"))
		return fmt.Errorf("os and version of the os for node should be selected"),errs
	}

	yumIDs := data.YumID
	if len(yumIDs) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020012,"warn","no yum configuration has be selected"))
	}

	return nil, errs
}
