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

package httpclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wangyysde/sysadm/utils"
)

// BuildDailer build net.Dialer for building RoundTripper. tcpTimeOut is for tcp timeout, keepAliveProbe is the time for probe keepalive alive.
// srcIP is the source IP address for the connection. If nil, a local address is automatically chosen.
// ref net.Dailer: https://pkg.go.dev/net#Dialer
func BuildDailer(tcpTimeOut,keepAliveProbe int,srcIP string) (*net.Dialer,error){

	timeOut := time.Duration(tcpTimeOut)
	keepAlive := time.Duration(keepAliveProbe)
	
	if strings.TrimSpace(srcIP) != "" {
		tcpAddr,err :=  net.ResolveTCPAddr("tcp",srcIP)
		if err != nil {
			return nil, err
		}

		localAddr = tcpAddr
	}
    

	return &net.Dialer{
		Timeout: timeOut * time.Second ,
		LocalAddr: localAddr,
		KeepAlive: keepAlive * time.Second,
	},nil
}

// BuildTlsClientConfig get the absolute path of ca,cert and key file. then read the content for certificate and key from the files.
// build tls.Config for a new https connection 
// return nil,error if any error occurred, otherwise return  *tls.Config,nil
func BuildTlsClientConfig(caFile,certFile, keyFile, workingDir string, insecureSkipVerify bool ) (*tls.Config, error){

	// get ca file absolute path
	if strings.TrimSpace(caFile) != "" {
		ca,err :=  utils.CheckFileIsRead(caFile,workingDir)
		if err != nil {
			return nil, fmt.Errorf("can not read ca file %s error %s",caFile,err)
		}
		caFile = ca
	}

	// get cert file absolute path
	if strings.TrimSpace(certFile) != "" {
		cert,err := utils.CheckFileIsRead(certFile,workingDir)
		if err != nil {
			return nil, fmt.Errorf("can not read cert file %s error %s",certFile,err)
		}
		certFile = cert

	}

	// get key file absolute path
	if strings.TrimSpace(keyFile) != "" {
		key,err := utils.CheckFileIsRead(keyFile,workingDir)
		if err != nil {
			return nil, fmt.Errorf("can not read keyt file %s error %s",keyFile,err)
		}
		keyFile = key
	}


    pool := x509.NewCertPool()
    var cert tls.Certificate

    if strings.TrimSpace(caFile) != "" {
        ca, err := ioutil.ReadFile(caFile)
        if err != nil {
            err = fmt.Errorf("ca has be specified %s but can not read it %s",caFile, err)
            return nil, err
        }
        pool.AppendCertsFromPEM(ca)
    }

    if strings.TrimSpace(certFile) != "" && strings.TrimSpace(keyFile) != "" {
        certPair, err :=  tls.LoadX509KeyPair(certFile,keyFile)
        if err != nil {
            return nil, fmt.Errorf("can not load certifaction pair %s",err)
        }
        cert = certPair
		return &tls.Config{
        	RootCAs: pool,
        	Certificates: []tls.Certificate{cert},
        	InsecureSkipVerify: insecureSkipVerify,
    	},nil
    } 

    return &tls.Config{
        InsecureSkipVerify: insecureSkipVerify,
    },nil
}

// BuildTlsRoundTripper build http.RoundTripper for creating https client. tlsHandshake: TLSHandshakeTimeout specifies the maximum amount of time waiting to wait for a TLS handshake. Zero means no timeout. 
// 
func BuildTlsRoundTripper(dialer *net.Dialer, tlsConf *tls.Config, tlsHandshake, idleConn,maxIdleConns,maxIdleConnsPerHost,maxConnsPerHost,readBuffer,writeBuffer int, disableKeepAlive, disableCompression, forceAttempHTTP2 bool )(http.RoundTripper,error){

	var dialerContext func(ctx context.Context, network string, addr string) (net.Conn,error) = nil
	if dialer != nil {
		dialerContext = dialer.DialContext
	}

	tlsHandshakeTimeout := time.Duration(tlsHandshake)
	idleConnTimeout := time.Duration(idleConn)

	if tlsConf == nil {
		return nil, fmt.Errorf("tls config must not nil for tls connection")
	}

	
    transport := &http.Transport{
        Proxy: http.ProxyFromEnvironment,
        DialContext: dialerContext,
        ForceAttemptHTTP2:     forceAttempHTTP2,
        MaxIdleConns:          maxIdleConns,
        IdleConnTimeout:       idleConnTimeout,
        TLSHandshakeTimeout:   tlsHandshakeTimeout,
        DisableKeepAlives:     disableKeepAlive,
        DisableCompression:    disableCompression,
        MaxIdleConnsPerHost:   maxIdleConnsPerHost,
        MaxConnsPerHost:       maxConnsPerHost,
        TLSClientConfig:       tlsConf,
		WriteBufferSize:       writeBuffer,
		ReadBufferSize: 	   readBuffer,
    }

    return transport,nil
}

// BuildHttpClient build a new http client for send http request.
func BuildHttpClient(roundTripper http.RoundTripper, timeout int) (*http.Client){
	httpTimeout := time.Duration(timeout)

	return &http.Client{
		Transport: roundTripper,
		Timeout: httpTimeout,
	}
}

/* 
   SendTlsRequest build HTTP query data, header inforation, BasicAuth, create a new http request and then send the request to server by client.
   return response body and nil if successful, otherwise return empty []byte and an error
*/
func NewSendRequest(r *RequestParams, client *http.Client, bodyReader io.Reader)([]byte, error){
	var body []byte

	if r == nil {
		return body, fmt.Errorf("can not handle http request without any request request parameters")
	}

	queryData, err := newHandleQueryData(r)
	if err != nil {
		return body, err
	}

	r.Url = strings.TrimSpace(r.Url)
	if r.Url == "" {
		return body, fmt.Errorf("HTTP request Url must not empty")
	}

	if queryData != "" {
		r.Url = r.Url + "?" + queryData
	}

	if client == nil {
		return body, fmt.Errorf("http client must not nil")
	}

	if ! CheckHttpMethod(r.Method) {
		return body, fmt.Errorf("HTTP method is not valid")
	}

	req,err := http.NewRequest(strings.ToUpper(r.Method), r.Url,bodyReader)
	if err != nil {
		return body, fmt.Errorf("create new HTTP request error %s",err)
	}

	if err := newAddReqestHeader(r,req); err != nil {
		return body, fmt.Errorf("add request header information onto request  error %s",err)
	}

	if err := newSetBasicAuth(r, req); err != nil {
		return body, err
	}

	/*
	client = &http.Client{
		Transport: sysadmTransport,
		Timeout: 30 * time.Second,
	}
	*/

	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	body,err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body,nil
}

/* 
   newHandleQueryData handle query data(build a string ) according to RequestParams.QueryData
*/
func newHandleQueryData(r *RequestParams)(string,error){

	if r == nil {
		return "",fmt.Errorf("can nog handle query data without any parameters")
	}
	data := r.QueryData
	ret := ""
	i := 0 
	for _,d := range data {
		if d.Key != "" {
			if i == 0 {
				ret = ret + d.Key + "=" + url.QueryEscape(d.Value)
				i = 1
			} else {
				ret = ret + "&" + d.Key + "=" + url.QueryEscape(d.Value)
			}
		}
	}

	return ret, nil
}

/*
    CheckHttpMethod Check that method is the correct HTTP method. return true if it is correct, otherewise return false
*/
func CheckHttpMethod(method string) bool{
	method = strings.ToUpper(strings.TrimSpace(method))
	if method == "" {
		return false
	}

	for _,v := range Methods {
		if strings.Compare(method,v) == 0 {
			return true
		}
	}

	return false
}

/* 
   newAddReqestHeader add default header data to request response
*/
func newAddReqestHeader(r *RequestParams,req *http.Request) error {
	
	if r == nil || req == nil{
		return fmt.Errorf("can not add HTTP header on a nil request or without any request parameteres")
	}
	r.Headers = defaultHeaders

	for _,h := range r.Headers {
		if h.Key != "" {
			req.Header.Set(h.Key,h.Value)
		}
	}

	return nil
}

func newSetBasicAuth(r *RequestParams,req *http.Request)(error){
	if r ==  nil || req == nil {
		return fmt.Errorf("can not set authorization information on a nil request or without any authorization parameters to be set")
	}

	authData := r.BasicAuthData
	if strings.EqualFold(authData["isBasicAuth"],"true") {
		if authData["username"] != "" && authData["password"] != "" {
			req.SetBasicAuth(authData["username"],authData["password"])
		} else {
			return fmt.Errorf("username or password  is empty")
		}
	}

	return nil
}
