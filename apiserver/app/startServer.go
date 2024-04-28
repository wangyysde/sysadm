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
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	certutil "k8s.io/client-go/util/cert"

	runtime "sysadm/apimachinery/runtime/v1beta1"
	sysadmPki "sysadm/apiserver/pki"
	"sysadm/sysadmerror"
	sysadmSysSetting "sysadm/syssetting"
)

var exitChan chan os.Signal
var shouldExit = false
var startedInSecret = false
var startedSecret = false

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

	exitChan = make(chan os.Signal, 1)
	logErrors(errs)

	for {
		e := startDaemon()
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080004, "error", "%s", e))
			if shouldExit {
				logErrors(errs)
				os.Exit(0)
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10080005, "info", "server restarting"))
		}
		if shouldExit {
			os.Exit(0)
		}
		signal.Notify(exitChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		s := <-exitChan
		if s == syscall.SIGHUP || s == syscall.SIGINT || s == syscall.SIGTERM {
			shouldExit = true
		}
		logErrors(errs)
		errs = errs[0:0]

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

	r := sysadmServer.New()
	r.Use(sysadmServer.Logger(), sysadmServer.Recovery())
	e = addResourceHanders(r)
	if e != nil {
		shouldExit = true
		return fmt.Errorf("add resources handlers error: %s", e)
	}

	// listen insecret port
	if runData.runConf.ConfServer.Insecret && runData.runConf.ConfServer.InsecretPort != 0 {
		go startInsecret(r)

	}

	if runData.runConf.ConfServer.IsTls {
		go startSecret(r)

	}

	if !startedInSecret && !startedSecret {
		shouldExit = true
	}

	return nil
}

func getCertAndKey(certType, scope int) (string, string, error) {
	gv := sysadmSysSetting.SchemaGroupVersion

	gvk := runtime.GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: runtime.APIVersionInternal}
	certKey, keyKey, e := sysadmSysSetting.GetCertAndKeyName(certType)
	if e != nil {
		return "", "", e
	}

	certQueryData, e := sysadmSysSetting.BuildQueryData(0, scope, 0, certKey)
	if e != nil {
		return "", "", e
	}
	certData, e := getResource(gvk, certQueryData)
	if e != nil {
		return "", "", e
	}
	keyQueryData, e := sysadmSysSetting.BuildQueryData(0, scope, 0, keyKey)
	if e != nil {
		return "", "", e
	}
	keyData, e := getResource(gvk, keyQueryData)
	if e != nil {
		return "", "", e
	}

	if len(certData) < 1 || len(keyData) < 1 {
		return "", "", fmt.Errorf("certification or key was not exist")
	}

	certs, e := sysadmSysSetting.GetValue(certData)
	if e != nil {
		return "", "", e
	}

	keys, e := sysadmSysSetting.GetValue(keyData)
	if e != nil {
		return "", "", e
	}

	if len(certs) > 1 || len(keys) > 1 {
		return "", "", fmt.Errorf("certification or key is duplicated in the system ")
	}

	return certs[0], keys[0], nil

}

func startInsecret(engine *sysadmServer.Engine) {
	var errs []sysadmerror.Sysadmerror
	listenStr := fmt.Sprintf("%s:%d", runData.runConf.ConfServer.Address, runData.runConf.ConfServer.InsecretPort)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030001, "debug", "listening  service to %s", listenStr))
	logErrors(errs)
	errs = errs[0:0]
	startedInSecret = true
	e := engine.Run(listenStr)
	if e != nil {
		startedInSecret = false
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030002, "error", "can not listent service. error %s", e))
		logErrors(errs)
	}
}

func startSecret(engine *sysadmServer.Engine) {
	var errs []sysadmerror.Sysadmerror
	e := prepareApiServerCerts()
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030003, "error", "prepare certification for apiServer error %s", e))
		logErrors(errs)
		return
	}

	tlsStr := fmt.Sprintf("%s:%d", runData.runConf.ConfServer.Address, runData.runConf.ConfServer.Port)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030004, "debug", "listening TLS service to %s", tlsStr))
	logErrors(errs)
	errs = errs[0:0]

	certPath := filepath.Join(runData.workingRoot, pkiPath)
	certFile := filepath.Join(certPath, apiServerCertFile)
	keyFile := filepath.Join(certPath, apiServerCertKeyFile)
	startedSecret = true
	e = engine.RunTLS(tlsStr, certFile, keyFile)
	if e != nil {
		startedSecret = false
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20030005, "error", "can not listent TLS service. error %s", e))
		logErrors(errs)
	}

	return
}

func prepareApiServerCerts() error {
	ca, caKey, e := getCertAndKey(sysadmSysSetting.CertTypeCa, sysadmSysSetting.SettingScopeGlobal)
	if e != nil {
		return e
	}
	caCerts, e := certutil.ParseCertsPEM([]byte(ca))
	if e != nil {
		return e
	}
	caCert := caCerts[0]

	e = sysadmPki.ValidateCertPeriod(caCert, time.Duration(0))
	if e != nil {
		return e
	}
	_, e = sysadmPki.ParseKeyPEM(caKey)
	if e != nil {
		return e
	}

	cert, key, e := getCertAndKey(sysadmSysSetting.CertTypeApiServer, sysadmSysSetting.SettingScopeGlobal)
	if e != nil {
		return e
	}
	certs, e := certutil.ParseCertsPEM([]byte(cert))
	if e != nil {
		return e
	}
	apiServerCert := certs[0]
	e = sysadmPki.ValidateCertPeriod(apiServerCert, time.Duration(0))
	if e != nil {
		return e
	}

	_, e = sysadmPki.ParseKeyPEM(key)
	if e != nil {
		return e
	}

	certPath := filepath.Join(runData.workingRoot, pkiPath)
	if err := os.MkdirAll(certPath, os.FileMode(0755)); err != nil {
		return err
	}

	fullCert := cert + "\n" + ca
	if err := os.WriteFile(filepath.Join(certPath, apiServerCertFile), []byte(fullCert), os.FileMode(0600)); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(certPath, apiServerCertKeyFile), []byte(key), os.FileMode(0600))
}
