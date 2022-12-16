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
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/wangyysde/sysadm/httpclient"
	"github.com/wangyysde/sysadm/sysadmerror"
	sysadmutils "github.com/wangyysde/sysadm/utils"
)

func run_DaemonPassive() ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	confNodeIdentifer := RunConf.Global.NodeIdentifer
	nodeIdentifer, err := getNodeIdentifer(confNodeIdentifer)
	if err != nil {
		errs := append(errs, sysadmerror.NewErrorWithStringLevel(50090101,"fatal","get node identifier error %s",err))
		return errs
	}
	runData.nodeIdentifer = nodeIdentifer
	
	url := buildGetCommandUrl()
	runData.getCommandUrl  = url 

	runData.getCommandParames = &httpclient.RequestParams{}
	runData.getCommandParames.Method = http.MethodPost
	runData.getCommandParames.Url = url
	
	signal.Notify(exitChan, syscall.SIGHUP,os.Interrupt,syscall.SIGTERM)
	defer forkNewProcess()
	getCommandLoop()
	
	return errs
}

func forkNewProcess(){
   var errs []sysadmerror.Sysadmerror

   s := <-exitChan
   if s == syscall.SIGHUP || s ==  syscall.SIGINT  || s == syscall.SIGTERM {
     os.Exit(0)
   }
   
   pid, _, err := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
   if  err != 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090102,"fatal","can not fork a new process %s.",err))
		logErrors(errs)
		os.Exit(10)
   }
     
    // exit normally for parent process
    if pid > 0 {
        os.Exit(0)
    }
    
    // set permission mask for child process 
    _ = syscall.Umask(0) 
     // set a new session for child process. 
    sid, s_errno := syscall.Setsid()
    if s_errno != nil || sid < 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090103,"fatal","syscall.Setsid error.errorno %d .",s_errno))
		logErrors(errs)
		os.Exit(11)
    }
    
   getCommandLoop() 
}

func getCommandLoop(){
	var errs []sysadmerror.Sysadmerror
	
	if runData.httpClient == nil {
		if err := buildHttpClient(); err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090104,"fatal","build http client error %s.",err))
			logErrors(errs)
			os.Exit(12)
		}

	}
	
	getCommandInterval := time.Duration(RunConf.Agent.Period)
	for {
		if runData.httpClient == nil {
			if err := buildHttpClient(); err != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090105,"warning","build http client error %s.we will try it again",err))
				logErrors(errs)
				errs = errs[0:0]
				continue
			}

		}

		nodeIdentifer :=  runData.nodeIdentifer
		nodeIdentiferJson, err :=  json.Marshal(nodeIdentifer)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090106,"warning","encoding node identifer to json string error %s. we will try it again",err))
			logErrors(errs)
			errs = errs[0:0]
			continue
		}

		requestParas := runData.getCommandParames
		body,err := httpclient.NewSendRequest(requestParas,runData.httpClient,strings.NewReader(sysadmutils.Bytes2str(nodeIdentiferJson)))
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090107,"warning","can not get command from server error %s",err))
			logErrors(errs)
			errs = errs[0:0]
			continue
		}
	
		go handleHTTPBody(body)
		time.Sleep(getCommandInterval * time.Second)
	}
}

func buildHttpClient() error{
	var rt http.RoundTripper = nil 
	dailer, err :=  httpclient.BuildDailer(DefaultTcpTimeout,DefaultKeepAliveProbeInterval,RunConf.Global.SourceIP)
	if err != nil {
		return err
	}

	if RunConf.Global.Tls.IsTls {
		tlsConf, err := httpclient.BuildTlsClientConfig(RunConf.Global.Tls.Ca,RunConf.Global.Tls.Cert,RunConf.Global.Tls.Key,RunConf.WorkingDir,RunConf.Global.Tls.InsecureSkipVerify)
		if err != nil {
			return err
		}

		rt, err = httpclient.BuildTlsRoundTripper(dailer,tlsConf,defaultTLSHandshakeTimeout,defaultIdleConnTimeout,defaultMaxIdleConns,defaultMaxIdleConnsPerHost,defaultMaxConnsPerHost,defaultReadBufferSize,defaultWriteBufferSize,defaultDisableKeepAives,defaultDisableCompression,true)
		if err != nil {
			return err
		}
	} else {
		rt, err = httpclient.NewBuildTlsRoundTripper(dailer,defaultIdleConnTimeout,defaultMaxIdleConns,defaultMaxIdleConnsPerHost,defaultMaxConnsPerHost,defaultMaxConnsPerHost,defaultWriteBufferSize,defaultDisableKeepAives,defaultDisableCompression,true)
		if err != nil {
			return err
		}
	}

	client := httpclient.BuildHttpClient(rt,defaultHTTPTimeOut)
	runData.httpClient =  client

	return nil 
}

