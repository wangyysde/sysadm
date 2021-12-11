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
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wangyysde/sysadm/sysadm/config"
	"github.com/wangyysde/sysadmServer"
	log "github.com/wangyysde/sysadmLog"
)

type StartParameters struct {
	// Point to configuration file path of server 
	ConfigPath  string 		
	OldConfigPath string
	accessLogFp *os.File
	errorLogFp *os.File
}

var StartData = &StartParameters{
	ConfigPath: "",
	OldConfigPath: "",
	accessLogFp: nil,
	errorLogFp: nil,
}

var definedConfig *config.Config 
var err error
var exitChan chan os.Signal

func DaemonStart(cmd *cobra.Command, cmdPath string){
	definedConfig, err = config.HandleConfig(StartData.ConfigPath,cmdPath)
	if err != nil {
		sysadmServer.Logf("error","error:%s",err)
		os.Exit(1)
	}
	
	setLogger()
	defer closeLogger()
	r := sysadmServer.New()
	r.Use(sysadmServer.Logger(),sysadmServer.Recovery())
	// Define handlers
  //  r.GET("/", func(c *sysadmServer.Context) {
  //      c.String(http.StatusOK, "Hello World!")
  //  })  
    r.GET("/ping", func(c *sysadmServer.Context) {
        c.String(http.StatusOK, "echo ping message")
    })  

	if err = addStaicRoute(r,cmdPath); err != nil {
		sysadmServer.Logf("error","%s",err)
		os.Exit(2)
	}

	exitChan = make(chan os.Signal)
	signal.Notify(exitChan, syscall.SIGHUP,os.Interrupt,os.Kill)
	go removeSocketFile()
    // Listen and serve on defined port
    sysadmServer.Log(fmt.Sprintf("Listening on port %s", "8080"),"info")
	if definedConfig.Server.Socket !=  "" {
		go func (sock string){
			err := r.RunUnix(sock)
			if err != nil {
				sysadmServer.Logf("error","We can not listen to %s, error: %s", sock, err)
				os.Exit(3)
			}
			defer removeSocketFile()
		}(definedConfig.Server.Socket)
	}
	defer removeSocketFile()
	liststr := fmt.Sprintf("%s:%d",definedConfig.Server.Address,definedConfig.Server.Port)
	sysadmServer.Logf("error","We are listen to %s,port:%d", liststr,definedConfig.Server.Port)
    r.Run(liststr)

}

// removeSocketFile deleting the socket file when server exit
func removeSocketFile(){
	<-exitChan
	_,err := os.Stat(definedConfig.Server.Socket)
	fmt.Print("This is runing\n")
	if err == nil {
		fmt.Print("aaaaaa\n")
		os.Remove(definedConfig.Server.Socket)
	}

	os.Exit(1)
}

// set parameters to accessLogger and errorLooger
func setLogger(){
	sysadmServer.SetLoggerKind(definedConfig.Log.Kind)
	sysadmServer.SetLogLevel(definedConfig.Log.Level)
	sysadmServer.SetTimestampFormat(definedConfig.Log.TimeStampFormat)
	if definedConfig.Log.AccessLog != ""{
		_,fp,err := sysadmServer.SetAccessLogFile(definedConfig.Log.AccessLog)
		if err != nil {
			sysadmServer.Logf("error","%s",err)
		}else{
			StartData.accessLogFp = fp
		}
		
	}

	if definedConfig.Log.SplitAccessAndError && definedConfig.Log.ErrorLog != "" {
		_,fp, err := sysadmServer.SetErrorLogFile(definedConfig.Log.ErrorLog)
		if err != nil {
			sysadmServer.Logf("error","%s",err)
		}else{
			StartData.errorLogFp = fp
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
}

// close access log file descriptor and error log file descriptor
// set AccessLogger  and ErrorLogger to nil
func closeLogger(){
	if StartData.accessLogFp != nil {
		fp := StartData.accessLogFp 
		fp.Close()
		sysadmServer.LoggerConfigVar.AccessLogger = nil
		sysadmServer.LoggerConfigVar.AccessLogFile = ""
	}

	if StartData.errorLogFp != nil {
		fp := StartData.errorLogFp 
		fp.Close()
		sysadmServer.LoggerConfigVar.ErrorLogger = nil
		sysadmServer.LoggerConfigVar.ErrorLogFile = ""
	}
}