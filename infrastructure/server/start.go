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
	"os"
	"fmt"
	"syscall"
	"os/signal"
	"github.com/spf13/cobra"

	"sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

var exitChan chan os.Signal

func Start(cmd *cobra.Command, cmdPath string){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10010001,"debug","try to start infrastructure server"))

	// parsing the configurations and get configurations from environment, then set them to CurrentRuningData after checked.
	err := handleConfig(cmdPath)
	errs = append(errs,err...)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("fatal"){
		logErrors(errs)
		os.Exit(9)
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10010002,"debug","configurations have been handled and it is ok"))

	// openning  loggers and set log format to loggers 
	err = setLogger()
	errs = append(errs, err...)
	defer closeLogger()
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10010003,"debug","loggers have been set"))
	logErrors(errs)
	errs =errs[0:0]

	// initating connections to DB server
	entity,e := initDB() 
	errs = append(errs,e...)
	if entity == nil {
		logErrors(errs)
		os.Exit(8)
	}
	defer entity.CloseDB()
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10010004,"debug","connections to DB server have be openned"))

	// initating server
	r := sysadmServer.New()
	r.Use(sysadmServer.Logger(),sysadmServer.Recovery())

	e = addHandlers(r)
	errs =  append(errs,e...)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("fatal"){
		logErrors(errs)
		os.Exit(7)
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10010005,"debug","handlers have be added."))

	definedConf := &CurrentRuningData.Config
	exitChan = make(chan os.Signal)
	signal.Notify(exitChan, syscall.SIGHUP,os.Interrupt,syscall.SIGTERM)
	go removeSocketFile(definedConf)

    // Listen and serve on defined port
	if definedConf.Server.Socket !=  "" {
		go func (sock string){
			err := r.RunUnix(sock)
			if err != nil {
				sysadmServer.Logf("error","We can not listen to %s, error: %s", sock, err)
				os.Exit(6)
			}
			defer removeSocketFile(definedConf)
		}(definedConf.Server.Socket)
	}
	defer removeSocketFile(definedConf)
	liststr := fmt.Sprintf("%s:%d",definedConf.Server.Address,definedConf.Server.Port)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10010006,"debug","listening to %s", liststr))
	logErrors(errs)
    r.Run(liststr)

}

/*
	log log messages to logfile or stdout
*/
func logErrors(errs []sysadmerror.Sysadmerror){

	for _,e := range errs {
		l := sysadmerror.GetErrorLevelString(e)
		no := e.ErrorNo
		sysadmServer.Logf(l,"erroCode: %d Msg: %s",no,e.ErrorMsg)
	}
	
}

// removeSocketFile deleting the socket file when server exit
func removeSocketFile(definedConfig *Config){
	<-exitChan
	_,err := os.Stat(definedConfig.Server.Socket)
	if err == nil {
		os.Remove(definedConfig.Server.Socket)
	}

	os.Exit(1)
}
