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
	"net"

	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadmServer"
)

func runGetHostIP(gotCommand *Command, c *sysadmServer.Context){
	ints,err := net.Interfaces()
	var retMap []map[string]interface{}  
	
	if err != nil {
		outPutMessage(c, gotCommand, "get interface list error %s\n", err)
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
   		nicsData := map[string]interface{} {
      		"name": dev.Name,	
   		}

   		if withMac {
	   		nicsData["mac"] = dev.HardwareAddr
   		}


   		addrs,err := devP.Addrs()
   		if err != nil {
	  		outPutMessage(c, gotCommand,"get ip of interface %s error %s\n",dev.Name,err)	
	  		return 
   		}
   		var ipMaps []map[string]interface{}  

   		for _,addr := range addrs {
      		ipnet,ok := addr.(*net.IPNet)
      		if ok {
				ipMap := map[string]interface{}{
					"ip": ipnet.IP.String(),
				}
		
				if withMask {
					ipMap["mask"] = ipnet.Mask.String()
				}
				ipMaps = append(ipMaps, ipMap)
      		}
    	}

		nicsData["ips"] = ipMaps
		retMap = append(retMap,nicsData)
	}

	retData := apiutils.NewBuildResponseDataForMap(true, 0, retMap)
	outPutResult(c, gotCommand, retData)
}
