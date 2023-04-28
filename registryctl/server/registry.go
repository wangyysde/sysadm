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
	https://datatracker.ietf.org/doc/rfc7235/
*/

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	sysadm "sysadm/sysadm/server"
	"sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

type BodyError struct {
	Code string `json:"code"` 
	Message string `json:"message"`
	Detail string `json:"detail"`
}

var RegistryErrs = map[string]BodyError{
	"blob_unknown": {
		Code: "BLOB_UNKNOWN",
		Message: "blob unknown to registry",
		Detail: "blob is unknown to the registry or manifest references an unknown layer",
	},
	"blob_upload_invalid": {
		Code: "BLOB_UPLOAD_INVALID",
		Message: "blob upload invalid",
		Detail: "The blob upload encountered an error and can no longer proceed.",
	},
	"blob_upload_unknow": { 
		Code: "BLOB_UPLOAD_UNKNOWN",
		Message: "blob upload unknown to registry",
		Detail: "blob upload has been cancelled or was never started",
	},
	"digest_invalid": {
		Code: "DIGEST_INVALID",
		Message: "provided digest did not match uploaded content",
		Detail: "manifest includes an invalid layer digest",
	},
	"manifest_blob_unknown": {
		Code: "MANIFEST_BLOB_UNKNOWN",
		Message: "blob unknown to registry",
		Detail: " manifest blob is unknown to the registry.",
	},
	"manifest_invalid": {
		Code: "MANIFEST_INVALID",
		Message: "manifest invalid",
		Detail: " manifest checks fail.",
	},
	"manifest_unknow": {
		Code: "MANIFEST_UNKNOWN",
		Message: "manifest unknown",
		Detail: "the manifest, identified by name and tag is unknown to the repository.",
	},
	"manifest_unverified": {
		Code: "MANIFEST_UNVERIFIED",
		Message: "manifest failed signature verification",
		Detail: "the manifest fails signature verification",
	},
	"name_invalid": {
		Code: "NAME_INVALID",
		Message: "invalid repository name",
		Detail: "invalid repository name encountered either during manifest validation",
	},
	"name_unknown": {
		Code: "NAME_UNKNOWN",
		Message: "manifest tag did not match URI",
		Detail: "repository name not known to registry",
	},
	"size_invalid": {
		Code: "SIZE_INVALID",
		Message: "provided length did not match content length",
		Detail: "provided length did not match content length",
	},
	"tag_invalid": {
		Code: "TAG_INVALID",
		Message: "manifest tag did not match URI",
		Detail: "the tag in the manifest does not match the uri tag",
	},
	"unauthorized": {
		Code: "UNAUTHORIZED",
		Message: "authentication required",
		Detail: "authentication required",
	},
	"denied": {
		Code: "DENIED",
		Message: "requested access to the resource is denied",
		Detail: "The access controller denied access for the operation on a resource.",
	},
	"unsupported": {
		Code: "UNSUPPORTED",
		Message: "The operation is unsupported.",
		Detail: "The operation was unsupported due to a missing implementation or invalid set of parameters.",
	},
	"internal_error": {
		Code: "INTERNALERROR",
		Message: "internal error has occurred.",
		Detail: "internal error has occurred.",
	},
}
type ReponseError struct {
	Errors []BodyError `json:"errors"`
}

/*
   adding handlers for registry. the path of this handlers for is /v2/xxxxxx
   this function called in daemonServer
*/
func addRegistryHandlers(r *sysadmServer.Engine)(([]sysadmerror.Sysadmerror)){
	var errs []sysadmerror.Sysadmerror
	
	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030014,"fatal","we can not add handler to a nil router."))
		return errs
	}
	
	v1 := r.Group("/api/v1.0/registryctl")
	{
		v1.POST("/:action", apiV1PostHandlers)
		v1.DELETE("/:action", apiV1DeleteHandlers)
	}

	r.GET("/v2/*path", getHandlers)
	r.POST("/v2/*path", postHandlers)
	r.HEAD("/v2/*path",headHandlers)
	r.PUT("/v2/*path",putHandlers)
	r.PATCH("/v2/*path",patchHandlers)
	r.DELETE("/v2/*path",deleteHandlers)
	
	return errs
}

/**
	/v2/<name>/manifests/<reference>: Pulling an Image Manifest
	/v2/<name>/blobs/<digest>：  Pulling a layer by digest. Layers are stored in the blob portion of the registry
	/v2/<name>/blobs/uploads/<uuid>: get the status of an upload
	/v2/_catalog: Listing Repositories()
	/v2/_catalog?n=<integer>: Listing Repositories with Pagination
	/v2/<name>/tags/list: Listing Image Tags
	/v2/<name>/tags/list?n=<integer>: Listing Image Tags with Pagination
*/
func getHandlers(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	path := strings.TrimSpace(c.Param("path")) 
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20304001,"debug","handling GET handlers for PATH %s",path))
	if path == "/" {
		logErrors(errs)
		v2RootGetHandler(c)
		return 
	}

	reverseProxyDirector := buildReverseProxyDirector(c)
	if reverseProxyDirector == nil {
		logErrors(errs)
		responseErrorToClient("internal_error",c)
		return 
	}
	modifyResponse :=  buildModifyReponse(c)

	if roundTripper == nil {
		buildRoundTripper()
	}

	registryProxy := httputil.ReverseProxy{
		Director: reverseProxyDirector,
		Transport: roundTripper,
		ModifyResponse: modifyResponse,
	}

	registryProxy.ServeHTTP(c.Writer,c.Request)
	logErrors(errs)
}

/**
	/v2/<name>/blobs/uploads/: Starting An Upload
	/v2/<name>/blobs/uploads/?mount=<digest>&from=<repository name>: Cross Repository Blob Mount
*/
func postHandlers(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	path := strings.TrimSpace(c.Param("path")) 
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20304004,"debug","handling POST handlers for PATH %s",path))
	logErrors(errs)
	passProxy(c)
}


/**
	/v2/<name>/manifests/<reference>: The image manifest can be checked for existence
	/v2/<name>/blobs/<digest>:  The existence of a layer can be checked
*/
func headHandlers(c *sysadmServer.Context) {

	var errs []sysadmerror.Sysadmerror
	// gets image name from uri of the request 
	path := c.Param("path")
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20304002,"debug","handling HEAD handlers for PATH %s",path))
	pathArray := strings.Split(path,"/")
	arrayLen := len(pathArray)

	if arrayLen < 4 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20304002,"error","image name unknow"))
		logErrors(errs)
		responseErrorToClient("name_invalid",c)
		return
	}

	// gets image name  from RequestURI
	var imageName string =""
	for i := 1; i< arrayLen-2; i++ {
		if imageName == "" {
			imageName = pathArray[i]
		}else{
			imageName = imageName + "/" + pathArray[i]
		}
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20304003,"debug","gets path is %s imagename is %s",path,imageName))
	logErrors(errs)
	errs = errs[0:0]
	_,ok := processImages[imageName]
	if !ok {
		r := c.Request
		username,_,_ := r.BasicAuth()
		processImages[imageName] = image{
			name: imageName,
			username: username,
		}

	}

	action := pathArray[(arrayLen - 2)]
	if strings.TrimSpace(strings.ToLower(action)) == "manifests"{
		err := checkManifestsExist(imageName,c)
		errs = append(errs, err...)
		logErrors(errs)
		return 
	}

	digest := pathArray[(arrayLen - 1)]
	recordBlob(imageName, digest)

	err := checkBlobExist(imageName,digest,c)
	errs = append(errs, err...)
	logErrors(errs)
}

/**
	/v2/<name>/blobs/uploads/<uuid>?digest=<digest>: Monolithic Upload
	/v2/<name>/blobs/uploads/<uuid>?digest=<digest>: Completed Upload
	/v2/<name>/manifests/<reference>: Pushing an Image Manifest
*/
func putHandlers(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror

	var director func (r *http.Request) = nil
	var modifyResponse func (r *http.Response) error = nil
	path := c.Param("path")
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20304005,"debug","handling PUT handlers for PATH %s",path))
	logErrors(errs)
	pathArray := strings.Split(path,"/")
	arrayLen := len(pathArray)
	action := pathArray[(arrayLen -2)]
	if strings.TrimSpace(strings.ToLower(action)) == "manifests"{
		director = buildReverseProxyDirector(c)
		modifyResponse = putManifestsResponse(c)
		if director == nil {
			responseErrorToClient("internal_error",c)
			return 
		}
	}
	
	if director == nil {
		director = buildReverseProxyDirector(c)
	}
	if director == nil {
		responseErrorToClient("internal_error",c)
		return 
	}

	if modifyResponse == nil {
		modifyResponse = buildModifyReponse(c)
	}
	

	if roundTripper == nil {
		buildRoundTripper()
	}

	registryProxy := httputil.ReverseProxy{
		Director: director,
		Transport: roundTripper,
		ModifyResponse: modifyResponse,
	}

	registryProxy.ServeHTTP(c.Writer,c.Request)
}

/**
	/v2/<name>/blobs/uploads/<uuid>: Chunked Upload
	handlers for patch method 
*/
func patchHandlers(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	path := c.Param("path")
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20304006,"debug","handling PATCH handlers for PATH %s",path))
	logErrors(errs)
	passProxy(c)
}

/**
	/v2/<name>/blobs/uploads/<uuid>: Canceling an Upload
	/v2/<name>/blobs/<digest>: Deleting a Layer
	/v2/<name>/manifests/<reference>: Deleting an Image
	handlers for DELETE method 
*/
func deleteHandlers(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	path := c.Param("path")
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20304007,"debug","handling DELETE handlers for PATH %s",path))
	logErrors(errs)
	
	reverseProxyDirector := buildReverseProxyDirector(c)
	if reverseProxyDirector == nil {
		responseErrorToClient("internal_error",c)
		return 
	}
	modifyResponse :=  buildModifyReponse(c)

	if roundTripper == nil {
		buildRoundTripper()
	}

	registryProxy := httputil.ReverseProxy{
		Director: reverseProxyDirector,
		Transport: roundTripper,
		ModifyResponse: modifyResponse,
	}

	registryProxy.ServeHTTP(c.Writer,c.Request)
}

func getRepositories()([]sysadmerror.Sysadmerror){
	var requestParams requestParams = requestParams{}
	var regUrl string = ""
	definedConfig := RuntimeData.RuningParas.DefinedConfig
	if definedConfig.Registry.Server.Tls {
		if definedConfig.Registry.Server.Port == 443 {
			regUrl = "https://" + definedConfig.Registry.Server.Host + "/v2/_catalog"
		} else {
			regUrl = "https://" + definedConfig.Registry.Server.Host + ":" + strconv.Itoa(definedConfig.Registry.Server.Port) + "/v2/_catalog"
		}
	}else {
		if definedConfig.Registry.Server.Port == 80 {
			regUrl = "http://" + definedConfig.Registry.Server.Host + "/v2/_catalog"
		} else {
			regUrl = "http://" + definedConfig.Registry.Server.Host + ":" + strconv.Itoa(definedConfig.Registry.Server.Port) + "/v2/_catalog"
		}
	}

	requestParams.url = regUrl
	requestParams.method = "GET"
	body,err := sendRequest(&requestParams)
	logErrors(err)
	fmt.Println(string(body))
	
	return err
}


/*
  v2RootGetHandler is for get /v2/
*/
func v2RootGetHandler(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	r := c.Request
	username,password,_ := r.BasicAuth()
	// checking user login
	ok := isLogin(username,password)
	
	if !ok {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030011,"debug","responsing to client with user has not login"))
		c.Header("Docker-Distribution-API-Version","registry/2.0")
		c.Header("WWW-Authenticate","Basic realm=\"basic-realm\"")
		be := []BodyError{{
			Code: "UNAUTHORIZED",
			Message: "unauthorized:unauthorized",
			Detail: "",},
		}
		
		var re = ReponseError{
			Errors: be,
		}
		reponseBody,_ := json.Marshal(re)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030012,"debug","the content responsed to the client is: %s",reponseBody))
		c.JSON(http.StatusUnauthorized,re)
		logErrors(errs)
		return 
	}

	reverseProxyDirector := buildReverseProxyDirector(c)
	if reverseProxyDirector == nil {
		responseErrorToClient("internal_error",c)
		return 
	}
	modifyResponse :=  buildModifyReponse(c)

	if roundTripper == nil {
		buildRoundTripper()
	}

	registryProxy := httputil.ReverseProxy{
		Director: reverseProxyDirector,
		Transport: roundTripper,
		ModifyResponse: modifyResponse,
	}

	registryProxy.ServeHTTP(c.Writer,c.Request)

	
}

func isLogin(username string, password string) bool {
	var errs []sysadmerror.Sysadmerror
	if username == "" || password == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030004,"error","username or password are empty."))
		logErrors(errs)
		return false
	}

	var reqUrl string = ""
	m := sysadm.Modules
	definedConfig := RuntimeData.RuningParas.DefinedConfig 
	if definedConfig.Sysadm.Server.Tls {
		if definedConfig.Sysadm.Server.Port == 443 {
			reqUrl = "https://" + definedConfig.Sysadm.Server.Host + "/api/" + definedConfig.Sysadm.ApiVerion + "/" + m["user"].Path +"/login" 
		} else {
			reqUrl = "https://" + definedConfig.Sysadm.Server.Host + ":" + strconv.Itoa(definedConfig.Sysadm.Server.Port) + "/api/" + definedConfig.Sysadm.ApiVerion + "/" + m["user"].Path +"/login" 
		}
	}else {
		if definedConfig.Sysadm.Server.Port == 80 {
			reqUrl = "http://" + definedConfig.Sysadm.Server.Host + "/api/" + definedConfig.Sysadm.ApiVerion  + "/" + m["user"].Path +"/login"
		} else {
			reqUrl = "http://" + definedConfig.Sysadm.Server.Host + ":" + strconv.Itoa(definedConfig.Sysadm.Server.Port) + "/api/" + definedConfig.Sysadm.ApiVerion +  "/" + m["user"].Path +"/login"
		}
	}

	var requestParams requestParams = requestParams{}
	requestParams.url = reqUrl
	requestParams.method = "POST"
	requestParams.data = append(requestParams.data,&requestData{key: "username", value: username})
	requestParams.data = append(requestParams.data,&requestData{key: "password", value: password})
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030005,"debug","try to execute the request with:%s",reqUrl))
	body,err := sendRequest(&requestParams)
	errs = append(errs, err...)

	if len(body) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030006,"error","the response from  the server is empty"))
		logErrors(errs)
		return false
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030007,"debug","got response body is: %s",string(body)))
	ret := &sysadm.ApiResponseStatus{}
	e := json.Unmarshal(body,ret)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030009,"error","can not parsing reponse body to json. error: %s",e))
		logErrors(errs)
		return false
	}

	if ret.Errorcode  != 0  {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2030004,"debug","can not login with errorcode: %d message: %s",ret.Errorcode,ret.Message))
		logErrors(errs)
		return false
	}
	
	logErrors(errs)
	return ret.Status
}

/*
	passProxy: 1. set BasicAuth on request and
	2. change the host and port of the request to registry server
	3. change the url of the request to registry server
	4. pass the request to registry server
*/
func passProxy(c *sysadmServer.Context) {
	
	// build proxy director
	reverseProxyDirector := buildReverseProxyDirector(c)
	if reverseProxyDirector == nil {
		responseErrorToClient("internal_error",c)
		return 
	}

	// build roundTripper
	if roundTripper == nil {
		buildRoundTripper()
	}

	modifyResponse :=  buildModifyReponse(c)
	// set ReverseProxy
	registryProxy := httputil.ReverseProxy{
		Director: reverseProxyDirector,
		Transport: roundTripper,
		ModifyResponse: modifyResponse,
	}

	registryProxy.ServeHTTP(c.Writer,c.Request)
	

}


func responseErrorToClient(errorCode string,c *sysadmServer.Context){
	value,ok := RegistryErrs[errorCode]
	if !ok {
		value = BodyError{
			Code: "UNKNOWN_ERROR",
			Message: "a unknow error has occurred",
			Detail: "a unknow error has occurred",
		}
	}
	var errField []BodyError
	errField = append(errField,value)
	retErr := ReponseError{
		Errors: errField,
	}
	c.Header("Docker-Distribution-API-Version","registry/2.0")
	c.Header("WWW-Authenticate","Basic realm=\"basic-realm\"")
	reponseBody,_ := json.Marshal(retErr)
	//c.JSON(http.StatusUnauthorized,reponseBody)
	c.JSON(http.StatusOK,reponseBody)
	
}

func checkManifestsExist(imageName string, c *sysadmServer.Context) []sysadmerror.Sysadmerror{
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20306001,"debug","checking manifest exist"))

	reverseProxyDirector := buildReverseProxyDirector(c)
	if reverseProxyDirector == nil {
		responseErrorToClient("internal_error",c)
		return errs
	}
	modifyResponse :=  buildModifyReponse(c)

	if roundTripper == nil {
		buildRoundTripper()
	}

	registryProxy := httputil.ReverseProxy{
		Director: reverseProxyDirector,
		Transport: roundTripper,
		ModifyResponse: modifyResponse,
	}

	registryProxy.ServeHTTP(c.Writer,c.Request)
	//TODO
	return errs
}

func checkBlobExist(imageName string, digest string, c *sysadmServer.Context) []sysadmerror.Sysadmerror{
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20306001,"debug","checking blob exist"))

	reverseProxyDirector := buildReverseProxyDirector(c)
	if reverseProxyDirector == nil {
		responseErrorToClient("internal_error",c)
		return errs
	}
	modifyResponse :=  modifyReponseForCheckBlobExist(c,imageName,digest)

	if roundTripper == nil {
		buildRoundTripper()
	}

	registryProxy := httputil.ReverseProxy{
		Director: reverseProxyDirector,
		Transport: roundTripper,
		ModifyResponse: modifyResponse,
	}

	registryProxy.ServeHTTP(c.Writer,c.Request)
	//TODO
	return errs
}

/*
	write image name and blob information into  processImages
	imageName is the name of the image 
	digest is the digest of a blob
*/
func recordBlob(imageName string, digest string){
	if strings.TrimSpace(imageName) == "" || strings.TrimSpace(digest) == "" {
		return 
	}
	image,ok := processImages[imageName]
	if !ok {
		return 
	}

	blobs := image.blobs
	for _,blob := range blobs {
		d := blob.digest
		if strings.TrimSpace(strings.ToLower(d)) == strings.TrimSpace(strings.ToLower(digest)){
			return
		}
	}

	b := blob {
		digest: digest,
		size: 0,
	}

	blobs = append(blobs, b)
	image.blobs = blobs
	processImages[imageName] = image

}



