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
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	//sysadmDB "github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/registryctl/config"
	"github.com/wangyysde/sysadm/sysadmerror"
	log "github.com/wangyysde/sysadmLog"
	"github.com/wangyysde/sysadmServer"
)

type StartParameters struct {
	// Point to configuration file path of server 
	ConfigPath  string 		
	OldConfigPath string
	accessLogFp *os.File
	errorLogFp *os.File
	sysadmRootPath string
}

var StartData = &StartParameters{
	ConfigPath: "",
	OldConfigPath: "",
	accessLogFp: nil,
	errorLogFp: nil,
	sysadmRootPath: "",
}

var definedConfig *config.Config 
var err error
var exitChan chan os.Signal

func DaemonStart(cmd *cobra.Command, cmdPath string){
	var errs []sysadmerror.Sysadmerror
	definedConfig, errs = config.HandleConfig(StartData.ConfigPath,cmdPath)
	e := setLogger()
	defer closeLogger()
	if len(e) > 0 {
		errs = appendErrs(errs,e)
	}
	logErrors(errs)
	maxLevel := sysadmerror.GetMaxLevel(errs)
	fatalLevel := sysadmerror.GetLevelNum("fatal")

	if maxLevel >= fatalLevel {
		os.Exit(202002)
	}
	
	//errs = errs[:0]
	
	if _,err = getSysadmRootPath(cmdPath); err != nil {
		sysadmServer.Logf("fatal","erroCode: 202003 Msg: can not get the root path of the program,err:%s ",err)
		os.Exit(202003)
	}
	
	r := sysadmServer.New()
	r.Use(sysadmServer.Logger(),sysadmServer.Recovery())
    
	/*
	dbConf := sysadmDB.DbConfig {
		Type: "postgre",
		Host: config.ConfigDefined.DB.Host, 
		Port: config.ConfigDefined.DB.Port, 
		User: config.ConfigDefined.DB.User,
		Password: config.ConfigDefined.DB.Password,
		DbName: config.ConfigDefined.DB.Dbname,
		SslMode: "disable",
		SslCa: "",
		SslCert: "",
		SslKey: "",
		MaxOpenConns: 50,
		MaxIdleConns: 20,
	}
	dbConfig,errs := sysadmDB.InitDbConfig(&dbConf,cmdPath)
	if len(errs) >0 {
		logErrors(errs)
	}

	dbEntity := dbConfig.Entity
	if sysadmerror.GetMaxLevel(errs) < sysadmerror.GetLevelNum("fatal") {
		errs := dbEntity.OpenDbConnect()
		if len(errs) >0 {
			logErrors(errs)
		}
		if sysadmerror.GetMaxLevel(errs) < sysadmerror.GetLevelNum("fatal") {
			dbData := sysadmDB.FieldData {
				"username": "test3User",
      			"email": "test3User@bzhy.com",
      			"password": "123456",
      			"realname": "Real Name",
			}
			rows,e := dbEntity.InsertData("test_user",dbData)
			if len(e) >0 {
				logErrors(e)
			} else {
				sysadmServer.Logf("info","rows %d data have be inserted",rows)
			}
		}

	}

	if err = ininDBConnect(); err != nil {
		sysadmServer.Logf("error","Open connection to postgre error:%s",err)
		os.Exit(4)
	}
	defer dbConnect.Close()

	err = addFormHandler(r,cmdPath)
	if err != nil {
		sysadmServer.Logf("error","error:%s",err)
		os.Exit(1)
	}
	// Define handlers
  //  r.GET("/", func(c *sysadmServer.Context) {
  //      c.String(http.StatusOK, "Hello World!")
  //  })  
	
	err = addRootHandler(r,cmdPath)
	if err != nil {
		sysadmServer.Logf("error","error:%s",err)
		os.Exit(1)
	}
	
    r.GET("/ping", func(c *sysadmServer.Context) {
        c.String(http.StatusOK, "echo ping message")
    })  

	if err = addStaicRoute(r,cmdPath); err != nil {
		sysadmServer.Logf("error","%s",err)
		os.Exit(2)
	}
*/
	r.GET("/ping", func(c *sysadmServer.Context) {
        c.String(http.StatusOK, "echo ping message")
    })  
	exitChan = make(chan os.Signal)
	signal.Notify(exitChan, syscall.SIGHUP,os.Interrupt,os.Kill)
	go removeSocketFile()
    // Listen and serve on defined port
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
	sysadmServer.Logf("info","We are listen to %s,port:%d", liststr,definedConfig.Server.Port)
    r.Run(liststr)

}

// removeSocketFile deleting the socket file when server exit
func removeSocketFile(){
	<-exitChan
	_,err := os.Stat(definedConfig.Server.Socket)
	if err == nil {
		os.Remove(definedConfig.Server.Socket)
	}

	os.Exit(1)
}

// set parameters to accessLogger and errorLooger
func setLogger()([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	sysadmServer.SetLoggerKind(definedConfig.Log.Kind)
	sysadmServer.SetLogLevel(definedConfig.Log.Level)
	sysadmServer.SetTimestampFormat(definedConfig.Log.TimeStampFormat)
	if definedConfig.Log.AccessLog != ""{
		_,fp,err := sysadmServer.SetAccessLogFile(definedConfig.Log.AccessLog)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(202000,"error","can not set access log file(%s) error: %s",definedConfig.Log.AccessLog,err))
		}else{
			StartData.accessLogFp = fp
		}
		
	}

	if definedConfig.Log.SplitAccessAndError && definedConfig.Log.ErrorLog != "" {
		_,fp, err := sysadmServer.SetErrorLogFile(definedConfig.Log.ErrorLog)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(202001,"error","can not set error log file(%s) error: %s",definedConfig.Log.ErrorLog,err))
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

	return errs
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

// Get the install dir path of  sysadm 
func getSysadmRootPath(cmdPath string) (string,error){
	dir ,error := filepath.Abs(filepath.Dir(cmdPath))
	if error != nil {
		return "",error
	}
	 
	dir = filepath.Join(dir,"../")
	StartData.sysadmRootPath = dir

	return dir, nil
}

/*
func addRootHandler(r *sysadmServer.Engine,cmdRunPath string) error {
	if r == nil {
		return fmt.Errorf("router is nil.")
	}

	r.Any("/",handleRootPath)
//	r.POST("/",handleRootPath)

	return nil
}

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
*/
func logErrors(errs []sysadmerror.Sysadmerror){

	for _,e := range errs {
		l := sysadmerror.GetErrorLevelString(e)
		no := e.ErrorNo
		sysadmServer.Logf(l,"erroCode: %d Msg: %s",no,e.ErrorMsg)
	}
	
}

func appendErrs(dst []sysadmerror.Sysadmerror,from []sysadmerror.Sysadmerror)([]sysadmerror.Sysadmerror){

//	for _,e := range from {
		dst = append(dst,from...)
//	}
	
	return dst
}