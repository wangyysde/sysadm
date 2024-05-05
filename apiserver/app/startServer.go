/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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
	"github.com/spf13/cobra"
	"github.com/wangyysde/sysadmServer"
	"os"
	"path/filepath"
	"sysadm/sysadmerror"
)

// var exitChan chan os.Signal
var shouldExit = false
var falseStartInSecret chan bool
var falseStartSecret chan bool

func StartServer(cmd *cobra.Command, args []string) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20010001, "debug", "starting  apiserver....."))

	ok, err := handlerConfig()
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		os.Exit(-1)
	}

	// initating loggers
	ok, err = setLogger()
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		os.Exit(-1)
	}
	defer closeLogger()
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080003, "debug", "loggers have been set"))
	logErrors(errs)
	errs = errs[0:0]

	// initating redis entity
	ok, err = initRedis()
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		os.Exit(-1)
	}
	defer closeRedisEntity()

	// initating DB entity
	ok, err = initDBEntity()
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		os.Exit(-1)
	}
	defer closeDBEntity()

	logErrors(errs)

	e := startDaemon()
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080004, "error", "%s", e))
		logErrors(errs)
		if shouldExit {
			logErrors(errs)
			os.Exit(0)
		}
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080006, "error", "unknow error"))
	logErrors(errs)
	os.Exit(-1)
}

func startDaemon() error {
	e := prepareSchema()
	if e != nil {
		shouldExit = true
		return fmt.Errorf("prepare resource schema data error: %s", e)
	}

	if !runData.runConf.ConfGlobal.Debug {
		sysadmServer.SetMode(sysadmServer.ReleaseMode)
	}
	r := sysadmServer.New()
	r.Use(sysadmServer.Logger(), sysadmServer.Recovery())
	e = addResourceHanders(r)
	if e != nil {
		shouldExit = true
		return fmt.Errorf("add resources handlers error: %s", e)
	}

	// listen insecret port

	falseStartInSecret = make(chan bool, 1)
	falseStartSecret = make(chan bool, 1)
	if runData.runConf.ConfServer.Insecret && runData.runConf.ConfServer.InsecretPort != 0 {
		go startInsecret(r)
	}

	if runData.runConf.ConfServer.IsTls {
		go startSecret(r)
	}

	s1 := <-falseStartInSecret
	s2 := <-falseStartSecret
	if s1 && s2 {
		shouldExit = true
		return fmt.Errorf("both secret and insecret service start error")
	}

	shouldExit = false
	return nil
}

func startInsecret(engine *sysadmServer.Engine) {
	var errs []sysadmerror.Sysadmerror
	listenStr := fmt.Sprintf("%s:%d", runData.runConf.ConfServer.Address, runData.runConf.ConfServer.InsecretPort)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030001, "debug", "listening  service to %s", listenStr))
	logErrors(errs)
	errs = errs[0:0]
	e := engine.Run(listenStr)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030002, "error", "can not listent service. error %s", e))
		logErrors(errs)
		falseStartInSecret <- true
	}
}

func startSecret(engine *sysadmServer.Engine) {
	var errs []sysadmerror.Sysadmerror

	tlsStr := fmt.Sprintf("%s:%d", runData.runConf.ConfServer.Address, runData.runConf.ConfServer.Port)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030004, "debug", "listening TLS service to %s", tlsStr))
	logErrors(errs)
	errs = errs[0:0]

	certPath := filepath.Join(runData.workingRoot, pkiPath)
	certFile := filepath.Join(certPath, apiServerFullCertFile)
	keyFile := filepath.Join(certPath, apiServerCertKeyFile)
	e := engine.RunTLS(tlsStr, certFile, keyFile)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030005, "error", "can not listent TLS service. error %s", e))
		logErrors(errs)
		falseStartSecret <- true
	}

	return
}
