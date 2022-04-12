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
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/wangyysde/sysadm/registryctl/config"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
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

var roundTripper http.RoundTripper = nil

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

type proxyParams struct {
	header http.Header
	method string
	url *url.URL
	contentLength int64
	transferEncoding []string
	host string
	trailer http.Header
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

func getRegistryUrl(c *sysadmServer.Context) string {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202029,"debug","preparing registry root url"))
	definedConfig := RuntimeData.RuningParas.DefinedConfig
	registryHost := definedConfig.Registry.Server.Host
	registryPort := definedConfig.Registry.Server.Port
	registryTls := definedConfig.Registry.Server.Tls
	var regUrlRoot string = ""
	if registryTls {
		if  registryPort  == 443 {
			regUrlRoot = "https://" + registryHost
		} else {
			regUrlRoot = "https://" + registryHost + ":" + strconv.Itoa(registryPort)
		}
	}else {
		if registryPort == 80 {
			regUrlRoot = "http://" + registryHost 	
		} else {
			regUrlRoot = "http://" + registryHost + ":" + strconv.Itoa(registryPort) 
		}
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202030,"debug","got registry root url is :%s",regUrlRoot))
	r := c.Request
	requestURI := r.RequestURI
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202031,"debug","got request uri is :%s",requestURI))
	registryURL := regUrlRoot +  requestURI
	logErrors(errs)
	return registryURL
}

func setProxyHeader(c *sysadmServer.Context, p *proxyParams){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202032,"debug","setting headers for requesting registry server"))

	r := c.Request
	p.header = r.Header
	logErrors(errs)
}

func setProxyMethod(c *sysadmServer.Context, p *proxyParams){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202033,"debug","setting method for requesting registry server"))

	r := c.Request
	method := r.Method
	if strings.TrimSpace(method) == "" {
		method = "GET"
	}

	p.method = method
	logErrors(errs)
}

func setProxyURL(c *sysadmServer.Context, p *proxyParams) error {
	var errs []sysadmerror.Sysadmerror
	rawURL := getRegistryUrl(c)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202034,"debug","setting URL(%s) for requesting registry server",rawURL))
	url,err := url.Parse(rawURL)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202035,"error","parse proxy url(%s) error: %s",rawURL,err))
		logErrors(errs)
		return err
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202036,"debug","URL: %#v",url))
	p.url= url
	logErrors(errs)
	return nil
}

func setContentLength(c *sysadmServer.Context, p *proxyParams){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202036,"debug","setting ContentLength for requesting registry server"))

	r := c.Request
	p.contentLength = r.ContentLength
	logErrors(errs)
}

func settransferEncoding(c *sysadmServer.Context, p *proxyParams){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202036,"debug","setting TransferEncoding for requesting registry server"))
	r := c.Request
	p.transferEncoding  = r.TransferEncoding
	logErrors(errs)
}

func setHost(c *sysadmServer.Context, p *proxyParams){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202037,"debug","setting host for requesting registry server"))

	definedConfig := RuntimeData.RuningParas.DefinedConfig
	registryHost := definedConfig.Registry.Server.Host
	registryPort := definedConfig.Registry.Server.Port
	host := registryHost + ":" + strconv.Itoa(registryPort)
	p.host = host
	logErrors(errs)

}

func setTrailer(c *sysadmServer.Context, p *proxyParams){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202038,"debug","setting trailer for requesting registry server"))

	r := c.Request
	p.trailer = r.Trailer
	logErrors(errs)
}

func buildRoundTripper(){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202038,"debug","building RoundTripper for requesting registry server"))
	definedConfig := RuntimeData.RuningParas.DefinedConfig

	var transport http.RoundTripper = &http.Transport{
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
			InsecureSkipVerify: definedConfig.Registry.Server.InsecureSkipVerify,
		},
	}

	roundTripper = transport
}

func buildReverseProxyDirector(c *sysadmServer.Context)(func(r *http.Request)) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202039,"debug","building proxy request parameters."))
	
	p := proxyParams{}
	setProxyHeader(c,&p)
	setProxyMethod(c,&p)
	err := setProxyURL(c,&p)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202040,"error","set proxy url error %s.",err))
		logErrors(errs)
		return nil
	}
	setContentLength(c,&p)
	settransferEncoding(c,&p)
	size := c.GetHeader("Content-Length")
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202041,"error","Content-Length %s.",size))
	logErrors(errs)
	setHost(c,&p)
	setTrailer(c,&p)
	return func(r *http.Request) {
		errs := setBasicAuth(r)
		logErrors(errs)
		r.Host = p.host
		r.URL = p.url
	}
}

func buildModifyReponse(c *sysadmServer.Context)(func(r *http.Response) error){
	return func(r *http.Response) error {
		var errs []sysadmerror.Sysadmerror
		header := r.Header
		//r.StatusCode = http.StatusOK
	//	body,_ := ioutil.ReadAll(r.Body)
	//	ret := &ReponseError{}
	//	_ = json.Unmarshal(body,ret)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202039,"debug","reponse header: %+v",header))
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202040,"debug","reponse StatusCode: %+v",r.StatusCode))
	//	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202040,"debug","reponse body: %+v",ret))

		logErrors(errs)
		
		return nil
	}
}

/*
	modifyReponseForCheckBlobExist: recording the infromation of blob ,such as digest, size to global variable processImages
	imageName: the name of image
	digest: the digest of a blob
*/
func modifyReponseForCheckBlobExist(c *sysadmServer.Context,imageName string,digest string)(func(r *http.Response) error){
	return func(r *http.Response) error {
		var errs []sysadmerror.Sysadmerror

		if r.StatusCode == http.StatusOK {
			recordBlob(imageName,digest)
			image := processImages[imageName]
			blobs := image.blobs
			for k,blob := range blobs {
				d := blob.digest
				if strings.TrimSpace(strings.ToLower(d)) == strings.TrimSpace(strings.ToLower(digest)){
					header := r.Header
					contentLength := header.Get("Content-Length")

					contentLengthInt,_ := strconv.Atoi(contentLength)
					blob.size = int64(contentLengthInt)
					blobs[k] = blob
					errs = append(errs, sysadmerror.NewErrorWithStringLevel(202039,"debug","record the information of  blob %s size exceeded%d ",digest,contentLengthInt))
				}
			}

			image.blobs = blobs 
			processImages[imageName] = image
			
		}

		logErrors(errs)
		return nil
	}
}



func putManifestsResponse(c *sysadmServer.Context)(func(r *http.Response) error){
	// get the name of the image from path
	path := c.Param("path")
	pathArray := strings.Split(path,"/")
	arrayLen := len(pathArray)
	// gets image name  from RequestURI
	var imageName string =""
	for i := 1; i< arrayLen-2; i++ {
		if imageName == "" {
			imageName = pathArray[i]
		}else{
			imageName = imageName + "/" + pathArray[i]
		}
	}
	reference := pathArray[(arrayLen - 1)]
	image := processImages[imageName]

	// sum the size of the all blobs of the image and set the information to the image
	var size int64 = 0
	for _,blob := range image.blobs {
		size += blob.size
	}
	image.size = size
	processImages[imageName] = image


	return func(r *http.Response) error {
		var errs []sysadmerror.Sysadmerror

		if r.StatusCode == http.StatusCreated {
			manifest := getManifests(imageName,reference)
			if manifest != nil {
				image := processImages[imageName]
				image.name = imageName
				image.tag = manifest.Tag
				image.architecture = manifest.Architecture
				processImages[imageName] = image
			}
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(202040,"debug","image information: %#v",processImages))
			logErrors(errs)
		}
		updataImage(imageName)
		delete(processImages,imageName)
		return nil
	}
}


func getManifests(name string, reference string ) *Manifest{
	var ret *Manifest = nil
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(name) == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202040,"error","image name is empty, can not get the manifest for empty name of image"))
		logErrors(errs)
		return ret
	}
	if strings.TrimSpace(reference) == "" {
		reference = "lastest"
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202039,"debug","building proxy request parameters."))

	// gets the configurations from RuntimeData
	definedConfig := RuntimeData.RuningParas.DefinedConfig
	registryHost := definedConfig.Registry.Server.Host
	registryPort := definedConfig.Registry.Server.Port
	registryTls := definedConfig.Registry.Server.Tls
	var urlStr string = ""
	if registryTls {
		if  registryPort  == 443 {
			urlStr = "https://" + registryHost
		} else {
			urlStr = "https://" + registryHost + ":" + strconv.Itoa(registryPort)
		}
	}else {
		if registryPort == 80 {
			urlStr = "http://" + registryHost 	
		} else {
			urlStr = "http://" + registryHost + ":" + strconv.Itoa(registryPort) 
		}
	}

	urlStr = urlStr + "/v2/" + name + "/manifests/" + reference

	var requestParams requestParams = requestParams{}
	requestParams.url = urlStr
	requestParams.method = "GET"
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202040,"debug","try to execute the request with url :%s",urlStr))
	body,err := sendRequest(&requestParams)
	errs = append(errs, err...)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202041,"debug","got the reponse body is: %s",body))

	if len(body) < 1 {
		logErrors(errs)
		return ret
	}

	ret = &Manifest{}
	e := json.Unmarshal(body,ret)
	if e != nil { 
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202041,"error","can not unmarshal body: %s",e))
		logErrors(errs)
		return nil
	}

	return ret
}