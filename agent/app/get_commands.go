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
	"net"
	"os"
	"strings"

	apiserver "sysadm/apiserver/app"
	"sysadm/sysadmerror"
	"sysadm/utils"
	"github.com/wangyysde/sysadmLog"
	"github.com/wangyysde/sysadmServer"
)

// this function is for "gethostip" command which is get all ip address on the localhost
func runGetHostIP(gotCommand *apiserver.CommandData, c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror

	// we should save command status to redis when command is asynchronous or agent running in passive mode
	commandStatusData := make(map[string]interface{}, 0)
	tmpC, err := handleCommandStatus(c, gotCommand, "", commandStatusData, apiserver.ComandStatusReceived, false, false)
	c = tmpC
	errs = append(errs, err...)

	ints, e := net.Interfaces()
	var retMap []map[string]interface{}

	if e != nil {
		data := make(map[string]interface{}, 0)
		tmpC, err := handleCommandStatus(c, gotCommand, "can not got interface data on localhost", data, apiserver.CommandStatusError, false, true)
		c = tmpC
		errs = append(errs, err...)
		_, err = setCommandLogIntoRedis(gotCommand.CommandSeq, "can not got interface data on localhost", sysadmLog.ErrorLevel)
		errs = append(errs, err...)
		handleCommandLogs(gotCommand.CommandSeq, 0, 0, c)
		logErrors(errs)
		return
	}

	cmdParameters := gotCommand.Parameters
	withMac := false
	withMask := false

	_, okWithMac := cmdParameters["withmac"]
	if okWithMac {
		withMac = true
	}

	_, okWithMask := cmdParameters["withmask"]
	if okWithMask {
		withMask = true
	}

	for _, dev := range ints {
		devP := &dev
		nicsData := map[string]interface{}{
			"name": dev.Name,
		}

		if withMac {
			nicsData["mac"] = dev.HardwareAddr
		}

		addrs, err := devP.Addrs()
		if err != nil {
			data := make(map[string]interface{}, 0)
			tmpC, err := handleCommandStatus(c, gotCommand, fmt.Sprintf("get ip of interface %s error %s", dev.Name, err), data, apiserver.CommandStatusError, false, true)
			c = tmpC
			errs = append(errs,err...)
			_, err = setCommandLogIntoRedis(gotCommand.CommandSeq, "can not got interface data on localhost", sysadmLog.ErrorLevel)
			errs = append(errs, err...)
			handleCommandLogs(gotCommand.CommandSeq, 0, 0, c)
			logErrors(errs)
			return
		}
		var ipMaps []map[string]interface{}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok {
				ipstr := ipnet.IP.String()
				if withMask {
					ipstr = ipstr + "/" + ipnet.Mask.String()
				}
				ipMap := map[string]interface{}{
					"ip": ipstr,
				}

				ipMaps = append(ipMaps, ipMap)
			}
		}

		nicsData["ips"] = ipMaps
		retMap = append(retMap, nicsData)
	}

	retData := make(map[string]interface{}, 0)
	retData["ips"] = retMap
	tmpC, err = handleCommandStatus(c, gotCommand, "", retData, apiserver.CommandStatusOK, false, true)
	c = tmpC
	errs = append(errs,err...)
	_, err = setCommandLogIntoRedis(gotCommand.CommandSeq, "the command of gethostip has be executed successfully", sysadmLog.InfoLevel)
	errs = append(errs, err...)
	handleCommandLogs(gotCommand.CommandSeq, 0, 0, c)
	logErrors(errs)
}

// this function is for "addyum" command which is add an yum configuration onto localhost
func runAddyum(gotCommand *apiserver.CommandData, c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror

	// we should save command status to redis when command is asynchronous or agent running in passive mode
	commandStatusData := make(map[string]interface{}, 0)
	tmpC, err := handleCommandStatus(c, gotCommand, "", commandStatusData, apiserver.ComandStatusReceived, false, false)
	c = tmpC
	errs = append(errs, err...)

	cmdParameters := gotCommand.Parameters
	yumName, okYumName := cmdParameters["yumName"]
	yumCatalog, okYumCatalog := cmdParameters["yumCatalog"]
	base_url, okBase_url := cmdParameters["base_url"]
	gpgcheck, okGpgcheck := cmdParameters["gpgcheck"]
	gpgkey, okGpgkey := cmdParameters["gpgkey"]
	yumName = strings.TrimSpace(yumName)
	yumCatalog = strings.TrimSpace(yumCatalog)
	base_url = strings.TrimSpace(base_url)
	if yumName == "" || yumCatalog == "" || base_url == "" || !okYumName || !okYumCatalog || !okBase_url {
		data := make(map[string]interface{}, 0)
		tmpC, err := handleCommandStatus(c, gotCommand, "parameters of command is invalid", data, apiserver.CommandStatusError, false, true)
		c = tmpC
		errs = append(errs, err...)
		_, err = setCommandLogIntoRedis(gotCommand.CommandSeq, "parameters of command is invalid", sysadmLog.ErrorLevel)
		errs = append(errs, err...)
		handleCommandLogs(gotCommand.CommandSeq, 0, 0, c)
		logErrors(errs)
		return
	}

	if !okGpgcheck || !okGpgkey {
		gpgcheck = "0"
	}

	gpgcheck = strings.TrimSpace(gpgcheck)
	gpgkey = strings.TrimSpace(gpgkey)
	if gpgcheck == "1" && gpgkey == "" {
		data := make(map[string]interface{}, 0)
		tmpC, err := handleCommandStatus(c, gotCommand, "gpcheck has be set to true but gpgkey has not be set", data, apiserver.CommandStatusError, false, true)
		c = tmpC
		errs = append(errs, err...)
		_, err = setCommandLogIntoRedis(gotCommand.CommandSeq, "gpcheck has be set to true but gpgkey has not be set", sysadmLog.ErrorLevel)
		errs = append(errs, err...)
		handleCommandLogs(gotCommand.CommandSeq, 0, 0, c)
		logErrors(errs)
		return
	}

	if strings.HasPrefix(gpgkey, "file:") && gpgcheck == "1" {
		gpgkeyfile := strings.TrimPrefix(gpgkey, "file://")
		_, e := utils.CheckFileIsReadable(gpgkeyfile, "")
		if e != nil {
			data := make(map[string]interface{}, 0)
			tmpC, err := handleCommandStatus(c, gotCommand, fmt.Sprintf("gpgkey has  be set, but gpgkey %s is not readable", gpgkeyfile), data, apiserver.CommandStatusError, false, true)
			c = tmpC
			errs = append(errs, err...)
			_, err = setCommandLogIntoRedis(gotCommand.CommandSeq, fmt.Sprintf("gpgkey has  be set, but gpgkey %s is not readable", gpgkeyfile), sysadmLog.ErrorLevel)
			errs = append(errs, err...)
			handleCommandLogs(gotCommand.CommandSeq, 0, 0, c)
			logErrors(errs)
			return
		}
	}

	var yumContent string = "\n"
	yumCatalog = yumContent + "[" + yumCatalog + "]\n"
	yumContent = yumContent + "name=" + yumName + "-" + yumCatalog + "\n"
	yumContent = yumContent + "baseurl=" + base_url + "\n"
	yumContent = yumContent + "gpgcheck=" + gpgcheck + "\n"
	yumContent = yumContent + "gpgkey=" + gpgkey + "\n"

	yumfile := yumConfRootPath + yumName + ".repo"
	fp, e := os.OpenFile(yumfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		data := make(map[string]interface{}, 0)
		tmpC, err := handleCommandStatus(c, gotCommand, fmt.Sprintf("can not write yum configuration file (%s) to the disk %s", yumfile, e), data, apiserver.CommandStatusError, false, true)
		c = tmpC
		errs = append(errs, err...)
		_, err = setCommandLogIntoRedis(gotCommand.CommandSeq, fmt.Sprintf("can not write yum configuration file (%s) to the disk %s", yumfile, e), sysadmLog.ErrorLevel)
		errs = append(errs, err...)
		handleCommandLogs(gotCommand.CommandSeq, 0, 0, c)
		logErrors(errs)
		return
	}

	_, e = fp.WriteString(yumContent)
	data := make(map[string]interface{}, 0)

	if e != nil {
		tmpC, err := handleCommandStatus(c, gotCommand, fmt.Sprintf("can not write yum configuration content to disk error:%s", e), data, apiserver.CommandStatusError, false, true)
		c = tmpC
		errs = append(errs, err...)
		_, err = setCommandLogIntoRedis(gotCommand.CommandSeq, fmt.Sprintf("can not write yum configuration content to disk error:%s", e), sysadmLog.ErrorLevel)
		errs = append(errs, err...)
	} else {
		_, err := handleCommandStatus(c, gotCommand, "", data, apiserver.CommandStatusOK, false, true)
		errs = append(errs, err...)
		_, err = setCommandLogIntoRedis(gotCommand.CommandSeq, "the yum configuration has be configurated", sysadmLog.InfoLevel)
		errs = append(errs, err...)
	}

	handleCommandLogs(gotCommand.CommandSeq, 0, 0, c)
	logErrors(errs)
}
