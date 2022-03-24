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
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	sysadmDB "github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/sysadm/config"
	"github.com/wangyysde/sysadm/sysadmerror"
	log "github.com/wangyysde/sysadmLog"
	"github.com/wangyysde/sysadmServer"
)

var err error
var exitChan chan os.Signal

func DaemonStart(cmd *cobra.Command, cmdPath string){
	var definedConfig *config.Config 
	// parsing the configurations and get configurations from environment, then set them to definedConfig after checked.
	definedConfig, err = config.HandleConfig(CliData.ConfigPath, cmdPath)
	if err != nil {
		sysadmServer.Logf("error","error:%s",err)
		os.Exit(1)
	}
	RuntimeData.RuningParas.DefinedConfig = definedConfig

	// Get the install dir path of  sysadm 
	if _,err = getSysadmRootPath(cmdPath); err != nil {
		sysadmServer.Logf("error","error:%s",err)
		os.Exit(2)
	}

	// set loggers according to the configurations.
	setLogger()
	defer closeLogger()

	// building configuration of DB , checking the configurations are available, instance DB instance
	dbConfig,errs := buildDBConfig(definedConfig, cmdPath)	
	if len(errs) >0 {
		logErrors(errs)
	}
	RuntimeData.RuningParas.DBConfig = dbConfig
	
	// open a DB connection
	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	errs = dbEntity.OpenDbConnect()
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("fatal"){
		os.Exit(4)
	}
	
	defer dbEntity.CloseDB()

	// newing an instance of sysadmServer
	r := sysadmServer.New()
	r.Use(sysadmServer.Logger(),sysadmServer.Recovery())
    
	// loading template files from the system.
	tmplPath := RuntimeData.StartParas.SysadmRootPath + "/tmpls/*.html" 
	r.LoadHTMLGlob(tmplPath)
	r.Delims(templateDelimLeft,templateDelimRight)

	/*
	r.SetFuncMap(template.FuncMap{
        "safe": func(str string) template.HTML {
            return template.HTML(str)
        },
    })
	*/

	// initating session
	if err = initSession(r); err != nil {
		sysadmServer.Logf("error","error:%s",err)
		os.Exit(5)
	}

	// adding all handlers
	err = addFormHandler(r,cmdPath)
	if err != nil {
		sysadmServer.Logf("error","error:%s",err)
		os.Exit(6)
	}
	addApiHandler(r,cmdPath)

	// adding project handlers
	errs = addProjectsHandler(r,cmdPath) 
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("fatal"){
		os.Exit(9)
	}

	// adding Root handlers
	err = addRootHandler(r,cmdPath)
	if err != nil {
		sysadmServer.Logf("error","error:%s",err)
		os.Exit(7)
	}
	
	// adding static handlers
	if err = addStaicRoute(r,cmdPath); err != nil {
		sysadmServer.Logf("error","%s",err)
		os.Exit(8)
	}

	// setting channel for passing signal to other process
	exitChan = make(chan os.Signal)
	signal.Notify(exitChan, syscall.SIGHUP,os.Interrupt,os.Kill)

	//removeSocketFile deleting the socket file when server exit
	go removeSocketFile()
   
	// Listen and serve on defined socket
	if RuntimeData.RuningParas.DefinedConfig.Server.Socket !=  "" {
		go func (sock string){
			err := r.RunUnix(sock)
			if err != nil {
				sysadmServer.Logf("error","We can not listen to %s, error: %s", sock, err)
				os.Exit(10)
			}
			defer removeSocketFile()
		}(RuntimeData.RuningParas.DefinedConfig.Server.Socket)
	}

	// starting listening
	defer removeSocketFile()
	liststr := fmt.Sprintf("%s:%d",RuntimeData.RuningParas.DefinedConfig.Server.Address ,RuntimeData.RuningParas.DefinedConfig.Server.Port)
	sysadmServer.Logf("info","We are listen to %s,port:%d", liststr,RuntimeData.RuningParas.DefinedConfig.Server.Port)
    r.Run(liststr)

}

// removeSocketFile deleting the socket file when server exit
func removeSocketFile(){
	<-exitChan
	_,err := os.Stat(RuntimeData.RuningParas.DefinedConfig.Server.Socket)
	if err == nil {
		os.Remove(RuntimeData.RuningParas.DefinedConfig.Server.Socket)
	}

	os.Exit(9)
}

// set parameters to accessLogger and errorLooger
func setLogger(){
	sysadmServer.SetLoggerKind(RuntimeData.RuningParas.DefinedConfig.Log.Kind)
	sysadmServer.SetLogLevel(RuntimeData.RuningParas.DefinedConfig.Log.Level)
	sysadmServer.SetTimestampFormat(RuntimeData.RuningParas.DefinedConfig.Log.TimeStampFormat)
	if RuntimeData.RuningParas.DefinedConfig.Log.AccessLog != ""{
		_,fp,err := sysadmServer.SetAccessLogFile(RuntimeData.RuningParas.DefinedConfig.Log.AccessLog)
		if err != nil {
			sysadmServer.Logf("error","%s",err)
		}else{
			 RuntimeData.RuningParas.AccessLogFp = fp
		}
		
	}

	if RuntimeData.RuningParas.DefinedConfig.Log.SplitAccessAndError && RuntimeData.RuningParas.DefinedConfig.Log.ErrorLog != "" {
		_,fp, err := sysadmServer.SetErrorLogFile(RuntimeData.RuningParas.DefinedConfig.Log.ErrorLog)
		if err != nil {
			sysadmServer.Logf("error","%s",err)
		}else{
			RuntimeData.RuningParas.ErrorLogFp  = fp
		}
	}
	sysadmServer.SetIsSplitLog(RuntimeData.RuningParas.DefinedConfig.Log.SplitAccessAndError)
	
	level, e := log.ParseLevel(RuntimeData.RuningParas.DefinedConfig.Log.Level)
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
	if RuntimeData.RuningParas.AccessLogFp != nil {
		fp := RuntimeData.RuningParas.AccessLogFp 
		fp.Close()
		sysadmServer.LoggerConfigVar.AccessLogger = nil
		sysadmServer.LoggerConfigVar.AccessLogFile = ""
	}

	if RuntimeData.RuningParas.ErrorLogFp != nil {
		fp := RuntimeData.RuningParas.ErrorLogFp 
		fp.Close()
		sysadmServer.LoggerConfigVar.ErrorLogger = nil
		sysadmServer.LoggerConfigVar.ErrorLogFile = ""
	}
}

// Get the install dir path of  sysadm 
func getSysadmRootPath(cmdPath string) (string,error){
	dir ,error := filepath.Abs(filepath.Dir(cmdPath))
	if error != nil {
		return "",error
	}
	 
	dir = filepath.Join(dir,"../")
	RuntimeData.StartParas.SysadmRootPath  = dir

	return dir, nil
}

// addRootHandler adding handler for root path
func addRootHandler(r *sysadmServer.Engine,cmdRunPath string) error {
	if r == nil {
		return fmt.Errorf("router is nil.")
	}

	r.Any("/",handleRootPath)

	return nil
}

// handleRootPath is the handler for root path
func handleRootPath(c *sysadmServer.Context){
	isLogin,_ := getSessionValue(c,"isLogin")
	if isLogin == nil  {
		formUri := formBaseUri + formsData["login"].formUri
		c.Redirect(http.StatusTemporaryRedirect, formUri)
	}

	tplData := map[string] interface{}{
		"htmlTitle": mainTitle,
	}

	c.HTML(http.StatusOK, "index.html", tplData)
}

// logging all logs in errs to access log file or error log file
func logErrors(errs []sysadmerror.Sysadmerror){

	for _,e := range errs {
		l := sysadmerror.GetErrorLevelString(e)
		no := e.ErrorNo
		sysadmServer.Logf(l,"erroCode: %d Msg: %s",no,e.ErrorMsg)
	}
	
	return
}

// buildDBConfig getting the configuratios of DB from sysadm configuration and checking the validity of them.
func buildDBConfig(definedConfig *config.Config, cmdPath string)(*sysadmDB.DbConfig,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs,sysadmerror.NewErrorWithStringLevel(100001,"debug","try to build DB configurations"))
	dbConf := sysadmDB.DbConfig {
		Type: definedConfig.DB.Type,
		Host: definedConfig.DB.Host, 
		Port: definedConfig.DB.Port, 
		User: definedConfig.DB.User,
		Password: definedConfig.DB.Password,
		DbName: definedConfig.DB.Dbname,
		SslMode: definedConfig.DB.Sslmode,
		SslCa: definedConfig.DB.Sslrootcert,
		SslCert: definedConfig.DB.Sslcert,
		SslKey: definedConfig.DB.Sslkey,
		MaxOpenConns: definedConfig.DB.DbMaxConnect,
		MaxIdleConns: definedConfig.DB.DbIdleConnect,
	}
	dbConfig,errs := sysadmDB.InitDbConfig(&dbConf,cmdPath)
	
	return dbConfig,errs
}