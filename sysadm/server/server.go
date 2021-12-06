/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2021 Bzhy Network. All rights reserved.
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
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/wangyysde/sysadm/sysadm/config"
	"github.com/wangyysde/sysadmServer"
	log "github.com/wangyysde/sysadmLog"
)

type StartParameters struct {
	// Point to configuration file path of server 
	ConfigPath  string 		
	OldConfigPath string
}

var StartData = &StartParameters{
	ConfigPath: "",
	OldConfigPath: "",
}

var definedConfig *config.Config 
var err error
func DaemonStart(cmd *cobra.Command, cmdPath string){
	definedConfig, err = config.HandleConfig(StartData.ConfigPath,cmdPath)
	if err != nil {
		sysadmServer.Logf("error","error:%s",err)
		os.Exit(1)
	}
	
	sysadmServer.SetLoggerKind(definedConfig.Log.Kind)
	sysadmServer.SetLogLevel(definedConfig.Log.Level)
	sysadmServer.SetTimestampFormat(definedConfig.Log.TimeStampFormat)
	if definedConfig.Log.AccessLog != ""{
		fn,fp,err := sysadmServer.SetAccessLogFile(definedConfig.Log.AccessLog)
		if err != nil {
			sysadmServer.Logf("error","%s",err)
		}else{
			defer fn(fp)
		}
		
	}

	if definedConfig.Log.SplitAccessAndError && definedConfig.Log.ErrorLog != "" {
		fu,fp, err := sysadmServer.SetErrorLogFile(definedConfig.Log.ErrorLog)
		if err != nil {
			sysadmServer.Logf("error","%s",err)
		}else{
			defer fu(fp)
		}
	}
	sysadmServer.SetIsSplitLog(definedConfig.Log.SplitAccessAndError)
	
	level, e := log.ParseLevel(definedConfig.Log.Level)
	if e != nil {
		sysadmServer.SetMode(sysadmServer.DebugMode)
	}else {
		if level >= log.DebugLevel {
			sysadmServer.SetMode(sysadmServer.DebugMode)
		}else{
			sysadmServer.SetMode(sysadmServer.ReleaseMode)
		}
	}

	r := sysadmServer.New()
	// Define handlers
    r.GET("/", func(c *sysadmServer.Context) {
        c.String(http.StatusOK, "Hello World!")
    })  
    r.GET("/ping", func(c *sysadmServer.Context) {
        c.String(http.StatusOK, "echo ping message")
    })  

    // Listen and serve on defined port
    sysadmServer.Log(fmt.Sprintf("Listening on port %s", "8080"),"info")
	if definedConfig.Server.Socket !=  "" {
		go func (sock string){
			err := r.RunUnix(sock)
			if err != nil {
				sysadmServer.Logf("error","We can not listen to %s, error: %s", sock, err)
				os.Exit(2)
			}
		}(definedConfig.Server.Socket)
	}
	liststr := fmt.Sprintf("%s:%d",definedConfig.Server.Address,definedConfig.Server.Port)
	sysadmServer.Logf("info","We are listen to %s,port:%d", liststr,definedConfig.Server.Port)
    r.Run(liststr)

}

