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

Ref: https://docs.docker.com/registry/spec/api/
*/

package server

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wangyysde/sysadm/registryctl/config"
	"github.com/wangyysde/sysadm/sysadmerror"
)

// TODO: the following parameters should be configurable in the future.
var (
	timeout time.Duration = 30 
	keepAlive time.Duration = 30
	tlshandshaketimeout time.Duration = 10
	disableKeepAlives bool = false
	disableCompression bool = false
	maxIdleConns int = 10 
	maxIdleConnsPerHost int = http.DefaultMaxIdleConnsPerHost
	maxConnsPerHost int = 0
	idleConnTimeout time.Duration = 90
)

type httpHeader struct {
	key string
	value string
}

type requestData struct {
	key string
	value string
}

type requestParams struct {
	headers []httpHeader
	data []*requestData
	method string
	url string
}

var headers []httpHeader
var defaultHeaders []httpHeader = []httpHeader{
	{
		key: "User-Agent", 
		value: ("registryctl-"+config.RegistryctlVer),
	},
}


// Ref: https://pkg.go.dev/net/http#Transport
var sysadmTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   timeout * time.Second,
		KeepAlive: keepAlive * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          maxIdleConns,
	IdleConnTimeout:       idleConnTimeout * time.Second,
	TLSHandshakeTimeout:   tlshandshaketimeout * time.Second,
	DisableKeepAlives: disableKeepAlives,
	DisableCompression: disableCompression,
	MaxIdleConnsPerHost: maxIdleConnsPerHost,
	MaxConnsPerHost: maxConnsPerHost,
	ExpectContinueTimeout: 1 * time.Second,
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}

func addReqestHeader(r *requestParams,req *http.Request)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202015,"debug","now handling the headers for the request"))
	if r == nil || req == nil{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202016,"fatal","can not handling the headers for nil request"))
		return errs
	}
	r.headers = append(r.headers,defaultHeaders...)
	
	for _,h := range headers {
		if h.key != "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(202017,"debug","adding key: %s value %s to the header of the request",h.key,h.value))
			req.Header.Set(h.key,h.value)
		}
	}

	return errs
}

func setBasicAuth(req *http.Request)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	definedConfig := RuntimeData.RuningParas.DefinedConfig
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202018,"debug","setting authorization for the request"))
	if req == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202019,"fatal","can not setting authorization for nil request"))
		return errs
	}

	if definedConfig.Registry.Credit.Username != "" && definedConfig.Registry.Credit.Password != "" {
		req.SetBasicAuth(definedConfig.Registry.Credit.Username,definedConfig.Registry.Credit.Password)
	} else {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202020,"warning","username or password for registry server is empty. we try to access registry  server without credit."))
	}

	return errs
}

func handleQueryData(r *requestParams)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202021,"debug","now handling the data for the request"))
	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202022,"fatal","can not handling the data for nil request"))
		return "",errs
	}
	data := r.data
	ret := ""
	i := 0 
	for _,d := range data {
		if d.key != "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(202023,"debug","adding key: %s value %s to the data of the request",d.key,d.value))
			if i == 0 {
				ret = ret + d.key + "=" + url.QueryEscape(d.value)
				i = 1
			} else {
				ret = ret + "&" + d.key + "=" + url.QueryEscape(d.value)
			}
		}
	}

	return ret, errs
}

func sendRequest(r *requestParams)([]byte, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var body []byte
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202024,"debug","now handling the request"))
	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202025,"fatal","can not handling a nil request"))
		return body, errs
	}
	
	fatalLevel := sysadmerror.GetLevelNum("fatal")

	var bodyReader *strings.Reader = nil
	if len(r.data) > 0 {
		query,err := handleQueryData(r) 
		maxLevel := sysadmerror.GetMaxLevel(err)
		errs = appendErrs(errs, err)
		if maxLevel >= fatalLevel {
			return body, errs
		}
		r.url = r.url + "?" + query
		//bodyReader = strings.NewReader(query)
	}

	client := &http.Client{
		Transport: sysadmTransport,
		Timeout: timeout * time.Second,
	}
	
	var req *http.Request
	var err error
	if bodyReader == nil {
		req,err = http.NewRequest(strings.ToUpper(r.method), r.url,nil)
	}else{
		req,err = http.NewRequest(strings.ToUpper(r.method), r.url,bodyReader)
	}
	
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202026,"fatal","can not create a new request, error: %s",err))
		return body, errs
	}
	e := addReqestHeader(r,req)
	errs = appendErrs(errs,e)
	maxLevel := sysadmerror.GetMaxLevel(errs)
	if maxLevel >= fatalLevel {
		return body, errs
	}

	e = setBasicAuth(req)
	errs = appendErrs(errs,e)
	maxLevel = sysadmerror.GetMaxLevel(errs)
	if maxLevel >= fatalLevel {
		return body, errs
	}

	resp, err := client.Do(req)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202027,"fatal","can not send request, error: %s",err))
		return body, errs
	}
	defer resp.Body.Close()

	body,err = ioutil.ReadAll(resp.Body)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202028,"fatal","can not gets reponse body contenet, error: %s",err))
		return body, errs
	}

	return body,errs
}