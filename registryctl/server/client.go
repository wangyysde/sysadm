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
	"encoding/base64"
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
	"github.com/wangyysde/sysadm/utils"
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
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(202018,"debug"," the request"))
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

func getRegistryRootUrl(c *sysadmServer.Context) string {
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
	
	return regUrlRoot
}

func buildRoundTripper(){
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

	return func(r *http.Request) {
		definedConfig := RuntimeData.RuningParas.DefinedConfig
		authStr := strings.TrimSpace(definedConfig.Registry.Credit.Username)+":"+strings.TrimSpace(definedConfig.Registry.Credit.Password)
		authEncode := base64.StdEncoding.EncodeToString(utils.Str2bytes(authStr))
		r.Header.Set("Authorization", ("Basic "+ authEncode))
		auth := r.Header.Get("Authorization")
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(202042,"debug","auth info: %s.",auth))
		
		rawURL := getRegistryUrl(c)
		url,_ := url.Parse(rawURL)
		r.URL = url
		logErrors(errs)
		r.Host = definedConfig.Registry.Server.Host + ":" + strconv.Itoa(definedConfig.Registry.Server.Port)
		
	}
}

/*
	buildModifyReponse: modifies the location field value of headers of response come from registry server 
*/
func buildModifyReponse(c *sysadmServer.Context)(func(r *http.Response) error){
	var errs []sysadmerror.Sysadmerror

	return func(r *http.Response) error {
		locationUrl := r.Header.Get("Location")
		if locationUrl != "" {
			definedConfig := RuntimeData.RuningParas.DefinedConfig
			registryUrl := definedConfig.RegistryUrl
			req := r.Request
			if registryUrl == "" {
				requestHost := req.Host
				if strings.TrimSpace(requestHost) != "" {
					registryUrl = requestHost
					definedConfig.RegistryUrl = registryUrl
				} else {
					serverAddr := definedConfig.Server.Address
					serverPort := strconv.Itoa(definedConfig.Server.Port)
					registryUrl = serverAddr + ":" + serverPort
					definedConfig.RegistryUrl = registryUrl
				}
			}

			registryUrl = strings.ToLower(registryUrl)
			if strings.HasPrefix(registryUrl,"http://") || strings.HasPrefix(registryUrl,"https://"){
				if strings.HasPrefix(registryUrl,"http://") {
					registryUrl = strings.TrimPrefix(registryUrl,"http://")
				} else {
					registryUrl = strings.TrimPrefix(registryUrl,"https://")
				}

				definedConfig.RegistryUrl = registryUrl
			}

			pUrl, _ := url.Parse(locationUrl)
			uri := pUrl.RequestURI()
			scheme := req.URL.Scheme
			newUrl := scheme + "://" + registryUrl + uri
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(2021001,"debug","new location url: %s",newUrl))
			logErrors(errs)
			r.Header.Set("Location",newUrl)
			
		}

		if r.StatusCode == http.StatusOK {
			return verifyManifestAndLayerWithDB(r)
		}
		return nil
	}
}

func verifyManifestAndLayerWithDB(r *http.Response) error{
	uri := r.Request.RequestURI
	if strings.TrimSpace(uri) != "" {
		uriArray := strings.Split(uri, "/")
		uriLen := len(uriArray)
		if strings.ToLower(uriArray[(uriLen - 2)]) == "manifests" {    // pull manifests of a image by a client
			// Get imageName form uri
			imageName := ""
			for i := 2; i< (uriLen - 2); i++ {
				if imageName == "" {
					imageName = uriArray[i]
				} else {
					imageName = imageName + "/" + uriArray[i]
				}
			}

			imgSets,_ := getImageInfoFromDB("","",imageName,"",0,0)
			if len(imgSets) < 1 {   // if there is not the information of the image for imageName in DB
				reference := uriArray[(uriLen - 1)]
				manifest := getManifests(imageName,reference)
				username,_,_ := r.Request.BasicAuth()
				if username == "" {
					username = "admin"
				}
				image := image{
					username: username,
					name: imageName,
					size: 0,
					tag: manifest.Tag,
					architecture: manifest.Architecture,
					digest: "",
					blobs: []blob{},
				}

				processImages[imageName] = image
				updataImage(imageName)
				delete(processImages, imageName)
				return nil
			}else {
				imgLine := imgSets[0]
				imageid := utils.Interface2String(imgLine["imageid"])
				updatePulltimesForImage(imageid,"")
				reference := uriArray[(uriLen - 1)]
				manifest := getManifests(imageName,reference)
				tagSets, errs := getTagInfoFromDB("",imageid,manifest.Tag,"","",0,0)
				logErrors(errs)
				if len(tagSets) < 1 {
					username,_,_ := r.Request.BasicAuth()
					if username == "" {
						username = "admin"
					}
					image := image{
						username: username,
						name: imageName,
						size: 0,
						tag: manifest.Tag,
						architecture: manifest.Architecture,
						digest: "",
						blobs: []blob{},
					}
					definedConfig := RuntimeData.RuningParas.DefinedConfig 
					apiServerTls := definedConfig.Sysadm.Server.Tls
					apiServerAddress := definedConfig.Sysadm.Server.Host
					apiServerPort := definedConfig.Sysadm.Server.Port
					apiVersion := definedConfig.Sysadm.ApiVerion
					userid,_ := getUserIdByUsername(apiServerTls,apiServerAddress,apiServerPort,apiVersion,username)
					if userid == 0 {
						userid = 1
					}
					processImages[imageName] = image
					imageID,_ := strconv.Atoi(imageid)
					_ = addTagsToDB(imageName,imageID,userid)
					delete(processImages, imageName)
				}else {
					tagLine := tagSets[0]
					tagid := utils.Interface2String(tagLine["tagid"])
					updatePulltimesForTag(tagid,"")
				}
				return nil 
			}

		}

	}
	return nil
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
			digest := r.Header.Get("Docker-Content-Digest")
			if manifest != nil {
				image := processImages[imageName]
				image.name = imageName
				image.tag = manifest.Tag
				image.architecture = manifest.Architecture
				image.digest = digest
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