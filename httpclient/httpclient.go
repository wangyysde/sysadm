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
ErrorCode: 106xxx
*/

package httpclient

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wangyysde/sysadm/sysadmerror"
)

// TODO: the following parameters should be configurable in the future.
var (
	timeout time.Duration = 30 
	localAddr net.Addr = nil // Ref https://pkg.go.dev/net#Addr
	keepAlive time.Duration = 30
	tlshandshaketimeout time.Duration = 10
	disableKeepAlives bool = false
	disableCompression bool = false
	maxIdleConns int = 10 
	maxIdleConnsPerHost int = http.DefaultMaxIdleConnsPerHost
	maxConnsPerHost int = 0
	idleConnTimeout time.Duration = 90
	httpClientVer = "v1.0"
)


type RequestData struct {
	Key string
	Value string
}

type RequestParams struct {
	Headers []RequestData
	QueryData []*RequestData
	BasicAuthData map[string]string
	Method string
	Url string
}

var headers []RequestData
var defaultHeaders []RequestData = []RequestData{
	{
		Key: "User-Agent", 
		Value: ("sysadmHttpClient-"+httpClientVer),
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

/* 
   addReqestHeader add default header data to request response
*/
func addReqestHeader(r *RequestParams,req *http.Request)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(106001,"debug","now handling the headers for the request"))
	if r == nil || req == nil{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(106002,"error","can not handling the headers for nil request"))
		return errs
	}
	r.Headers = append(r.Headers,defaultHeaders...)
	
	for _,h := range r.Headers {
		if h.Key != "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(106003,"debug","adding key: %s value %s to the header of the request",h.Key,h.Value))
			req.Header.Set(h.Key,h.Value)
		}
	}

	return errs
}

func setBasicAuth(r *RequestParams,req *http.Request)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(106004,"debug","setting authorization for the request"))
	if req == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(106005,"fatal","can not setting authorization for nil request"))
		return errs
	}

	authData := r.BasicAuthData
	if strings.EqualFold(authData["isBasicAuth"],"true") {
		if authData["username"] != "" && authData["password"] != "" {
			req.SetBasicAuth(authData["username"],authData["password"])
		} else {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(106006,"error","username or password  is empty."))
		}
	}

	return errs
}

func handleQueryData(r *RequestParams)(string,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(106007,"debug","now handling the data for the request"))
	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(106008,"error","can not handling the data for nil request"))
		return "",errs
	}
	data := r.QueryData
	ret := ""
	i := 0 
	for _,d := range data {
		if d.Key != "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(106009,"debug","adding key: %s value %s to the data of the request",d.Key,d.Value))
			if i == 0 {
				ret = ret + d.Key + "=" + url.QueryEscape(d.Value)
				i = 1
			} else {
				ret = ret + "&" + d.Key + "=" + url.QueryEscape(d.Value)
			}
		}
	}

	return ret, errs
}

func SendRequest(r *RequestParams)([]byte, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var body []byte
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(106010,"debug","now handling the request"))
	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(106011,"fatal","can not handling a nil request"))
		return body, errs
	}
	
	fatalLevel := sysadmerror.GetLevelNum("fatal")

	var bodyReader *strings.Reader = nil
	if len(r.QueryData) > 0 {
		query,err := handleQueryData(r) 
		maxLevel := sysadmerror.GetMaxLevel(err)
		errs = append(errs, err...)
		if maxLevel >= fatalLevel {
			return body, errs
		}
		r.Url = r.Url + "?" + query
	
	}

	client := &http.Client{
		Transport: sysadmTransport,
		Timeout: timeout * time.Second,
	}
	
	var req *http.Request
	var err error
	if bodyReader == nil {
		req,err = http.NewRequest(strings.ToUpper(r.Method), r.Url,nil)
	}else{
		req,err = http.NewRequest(strings.ToUpper(r.Method), r.Url,bodyReader)
	}
	
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(106012,"fatal","can not create a new request, error: %s",err))
		return body, errs
	}
	e := addReqestHeader(r,req)
	errs = append(errs,e...)
	maxLevel := sysadmerror.GetMaxLevel(errs)
	if maxLevel >= fatalLevel {
		return body, errs
	}

	e = setBasicAuth(r,req)
	errs = append(errs,e...)
	maxLevel = sysadmerror.GetMaxLevel(errs)
	if maxLevel >= fatalLevel {
		return body, errs
	}

	resp, err := client.Do(req)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(106013,"fatal","can not send request, error: %s",err))
		return body, errs
	}
	defer resp.Body.Close()

	body,err = ioutil.ReadAll(resp.Body)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(106014,"fatal","can not gets reponse body contenet, error: %s",err))
		return body, errs
	}

	return body,errs
}