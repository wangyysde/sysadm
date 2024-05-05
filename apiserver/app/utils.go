/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
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
	"github.com/wangyysde/sysadmServer"
	"net/http"
	"sysadm/httpclient"
	"sysadm/sysadmerror"
)

// log log messages to logfile or stdout
func logErrors(errs []sysadmerror.Sysadmerror) {

	for _, e := range errs {
		l := sysadmerror.GetErrorLevelString(e)
		no := e.ErrorNo
		sysadmServer.Logf(l, "erroCode: %d Msg: %s", no, e.ErrorMsg)
	}
}

// createHttpClient create a http client for a http request
func createHttpClient(tcpTimeOut, keepaliveProbe, tlsHandshakeTimeout, idleConnTimeout int, srcIP string, caPEM,
	certPEM, keyPEM []byte, insecureSkipVerify, isTLS bool) (*http.Client, error) {
	var rt http.RoundTripper = nil

	if tcpTimeOut == 0 {
		tcpTimeOut = defaultHttpTimeout
	}

	if keepaliveProbe == 0 {
		keepaliveProbe = defaultHttpKeepAliveProbe
	}

	if tlsHandshakeTimeout == 0 {
		tlsHandshakeTimeout = defaultTLSHandshakeTimeout
	}

	if idleConnTimeout == 0 {
		idleConnTimeout = defaultIdleConnTimeout
	}

	dialer, err := httpclient.BuildDailer(tcpTimeOut, keepaliveProbe, srcIP)
	if err != nil {
		return nil, err
	}

	if isTLS {
		tlsConf, err := httpclient.BuildTlsConfWithPEMBlock(caPEM, certPEM, keyPEM, insecureSkipVerify)
		if err != nil {
			return nil, err
		}

		rt, err = httpclient.BuildTlsRoundTripper(dialer, tlsConf, tlsHandshakeTimeout, idleConnTimeout,
			0, 0, 0, 0, 0, true,
			true, false)
		if err != nil {
			return nil, err
		}
	} else {
		rt, err = httpclient.NewBuildRoundTripper(dialer, idleConnTimeout, 0, 0, 0,
			0, 0, true, true, false)
		if err != nil {
			return nil, err
		}
	}

	client := httpclient.BuildHttpClient(rt, defaultHttpTimeout)

	return client, nil
}

/*
// buildSendCommandRequestParas build completion url and set request parameters for sending a command data to client.
func buildSendCommandRequestParas(data *commandDataBeSent) (bool, *httpclient.RequestParams, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	requestUrl, e := httpclient.BuildUrl(data.agentAddress, data.commandUri, data.agentPort, data.agentIsTls)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20041001, "error", "build request"+
			"url error %s", e))
		logErrors(errs)
		return false, nil, errs
	}

	requestParas := &httpclient.RequestParams{
		Method: http.MethodPost,
		Url:    requestUrl,
	}

	return true, requestParas, errs
}
*/

// buildClientRequestParas build completion url and set request parameters for connecting to a client.
func buildClientRequestParas(address, uri string, port int, isTls bool) (bool, *httpclient.RequestParams, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	requestUrl, e := httpclient.BuildUrl(address, uri, port, isTls)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20041002, "error", "build request"+
			"url error %s", e))
		logErrors(errs)
		return false, nil, errs
	}

	requestParas := &httpclient.RequestParams{
		Method: http.MethodPost,
		Url:    requestUrl,
	}

	return true, requestParas, errs
}
