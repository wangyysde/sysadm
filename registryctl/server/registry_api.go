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

	errorCode: 206xxxxx
*/

package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/wangyysde/sysadm/httpclient"
	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
	"github.com/wangyysde/sysadmServer"
)



func (r RegistryCtl)GetModuleName()string{
	return "registryctl"
}

func (r RegistryCtl)GetActionList()[]string{
	return registryctlActions
}

func apiV1PostHandlers(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror
	entity := RegistryCtl{}

	action := strings.TrimSuffix(strings.TrimPrefix(c.Param("action"),"/"),"/")
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20600001,"debug","handling the request for registryctl module  with action %s.",action))
	switch strings.ToLower(action){
	case "imagelist":
		err := entity.listImage(c)
		errs = append(errs, err...)
	case "getcount":
		err := entity.getCount(c)
		errs = append(errs, err...)
	case "taglist":
		err := entity.tagList(c)
		errs = append(errs,err...)
	default: 
		err := entity.ActionNotFound(c,action)
		errs = append(errs,err...)
	}

	logErrors(errs)

}

func apiV1DeleteHandlers(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror
	entity := RegistryCtl{}

	action := strings.TrimSuffix(strings.TrimPrefix(c.Param("action"),"/"),"/")
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20700001,"debug","handling the request for registryctl module  with action %s in DELETE method.",action))
	switch strings.ToLower(action){
	case "imagedel":
		err := entity.imageDelHandler(c)
		errs = append(errs, err...)
	case "tagdel":
		err := entity.tagDelHandler(c)
		errs = append(errs, err...)
	default: 
		err := entity.ActionNotFound(c,action)
		errs = append(errs,err...)
	}

	logErrors(errs)
}

/*
	ListImage list images information acccording "imageid","projectid","name","ownerid","start","numperpage"
	imageid: image "id, like id1,id2,id3.... 
	projectid: project id like "id1,id2,id3...." 
	name: image name. 
	ownerid：user id like "id1,id2,id3...." 
*/
func (r RegistryCtl)listImage(c *sysadmServer.Context) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	dataMap,err := utils.GetRequestData(c,[]string{"imageid","projectid","name","ownerid","start","numperpage"})
	errs = append(errs,err...)

	imageid := strings.TrimSpace(dataMap["imageid"])
	projectid := strings.TrimSpace(dataMap["projectid"])
	name := strings.TrimSpace(dataMap["name"])
	ownerid  := strings.TrimSpace(dataMap["ownerid"])
	startStr  := strings.TrimSpace(dataMap["start"])
	numperpageStr  := strings.TrimSpace(dataMap["numperpage"])
	start := 0
	numperpage := 0 
	if startStr != "" {
		start,_ = strconv.Atoi(startStr)
	}
	if numperpageStr != "" {
		numperpage,_ = strconv.Atoi(numperpageStr)
	}
	
	dataSet,err := getImageInfoFromDB(imageid,projectid,name,ownerid,start, numperpage)
	errs = append(errs,err...)
	err = apiutils.SendResponseForMap(c,dataSet)
	errs = append(errs,err...)

	return errs
}

/*
	imageDelHandler is handler for delete image throught registryctl api 
*/
func (r RegistryCtl)imageDelHandler(c *sysadmServer.Context) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	dataMap,err := utils.GetRequestData(c,[]string{"imageid"})
	errs = append(errs,err...)
	imageid,ok := dataMap["imageid"]
	if !ok {
		msg := "imageid which will be deleted has not found "
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20600004,"error","imageid which will be deleted has not found"))
		err := apiutils.SendResponseForErrorMessage(c,20600004,msg)
		errs = append(errs,err...)
		return errs
	}

	imageidArray := strings.Split(imageid,",")
	for _,id := range imageidArray {
		imgSets,err := getImageInfoFromDB(id,"","","",0,0)
		errs = append(errs,err...)
		if len(imgSets) < 1 {
			msg := "image for id " + id + "was not found"
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20600005,"error","image for id %s was not found",id))
			err := apiutils.SendResponseForErrorMessage(c,20600005,msg)
			errs = append(errs,err...)
			return errs
		}

		for _,imageLine := range imgSets {
			imageID := utils.Interface2String(imageLine["imageid"])
			res,err := getTagInfoFromDB("",imageID,"","","",0,0)
			errs = append(errs,err...)
			tagIDStr := ""
			for _,tagLine := range res {
				tmpID := utils.Interface2String(tagLine["tagid"])
				if tagIDStr == "" {
					tagIDStr = tmpID
				}else {
					tagIDStr = tagIDStr + "," + tmpID
				}
			}
			e := r.tagDel(c, tagIDStr,"")
			if e != nil {
				msg := "can not delete tags belonging to the image which will be deleted"
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(20600010,"error","can not delete tags belonging to the image which will be deleted.tagID %s error %s",tagIDStr,e))
				err := apiutils.SendResponseForErrorMessage(c,20600015,msg)
				errs = append(errs,err...)
				return errs
			}
			resNum := delImagesFromDB(imageID,"")
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20600011,"debug","%d rows data for image has be deleted from DB with image id: %s",resNum,imageID))
		}

	}

	msg := "image with id " + imageid +" have been deleted."
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20600007,"debug",msg))
	err = apiutils.SendResponseForSuccessMessage(c,msg)	
	errs=append(errs,err...)
	return errs
}

/*
	ListImage list images information acccording "imageid","projectid","name","ownerid","start","numperpage"
	imageid: image "id, like id1,id2,id3.... 
	projectid: project id like "id1,id2,id3...." 
	name: image name. 
	ownerid：user id like "id1,id2,id3...." 
*/
func (r RegistryCtl)tagDelHandler(c *sysadmServer.Context) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	dataMap,err := utils.GetRequestData(c,[]string{"tagid","digest"})
	errs = append(errs,err...)
	tagid,okTagid := dataMap["tagid"]
	digest,okDigest := dataMap["digest"]
	
	if !okTagid && !okDigest {
		msg := "tag id or digest is invalid" 
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20600007,"error",msg))
		err := apiutils.SendResponseForErrorMessage(c,20600007,msg)
		errs = append(errs,err...)
		return errs
	}
	
	e := r.tagDel(c,tagid,digest)
	if e != nil {
		msg := fmt.Sprintf("delete tags with id %s error %s",tagid,e) 
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20600007,"error",msg))
		err := apiutils.SendResponseForErrorMessage(c,20600008,msg)
		errs = append(errs, err...)
		
	}
	err = apiutils.SendResponseForSuccessMessage(c,"tag has be deleted.")
	errs = append(errs,err...)
	return errs
}



func (r RegistryCtl)tagDel(c *sysadmServer.Context, tagID string, digest string) (error){
	if strings.TrimSpace(tagID) == "" && strings.TrimSpace(digest) == "" {
		return fmt.Errorf("both tagid  and digest are empty")
	}
	
	var res []map[string]interface{}
	if strings.TrimSpace(tagID) != "" {
		res,_ = getTagInfoFromDB(tagID,"","","","",0,0)
	} else {
		res,_ = getTagInfoFromDB("","","","",digest,0,0)
	}
	
	for _,line := range res {
		tmpTagid := utils.Interface2String(line["tagid"])
		tmpImageid := utils.Interface2String(line["imageid"])
		tmpDigest := utils.Interface2String(line["digest"])
		imageInfo,_ := getImageInfoFromDB(tmpImageid,"","","",0,0)
		if len(imageInfo) < 1 {
			return fmt.Errorf("can not got image name when deleting tag")
		}
		imageLine := imageInfo[0]
		imageName := utils.Interface2String(imageLine["name"])
		blobInfo,_ := getBlobInfoFromDB("",tmpTagid,"")
		blobDigest := ""
		for _,blobLine := range blobInfo {
			tmpDigest := utils.Interface2String(blobLine["digest"])
			if strings.TrimSpace(tmpDigest) != "" {
				if blobDigest == "" {
					blobDigest = strings.TrimSpace(tmpDigest)
				}else{
					blobDigest = blobDigest + "," + strings.TrimSpace(tmpDigest)
				}
			}
		}
		
		if strings.TrimSpace(imageName) != "" && strings.TrimSpace(blobDigest) != "" {
			err := r.layerDel(c,imageName,blobDigest)
			if err != nil{
				return fmt.Errorf("delete layer %s(digest) of image %s error %s",blobDigest,imageName,err)
			}
		}

		if strings.TrimSpace(tmpDigest) == "" {
			return fmt.Errorf("manifests digest is empty")
		}
		
		registryUrl  := getRegistryRootUrl(c)
		requestParas :=  httpclient.RequestParams{}
		definedConfig := RuntimeData.RuningParas.DefinedConfig
		requestParasPr,_ := httpclient.AddBasicAuthData(&requestParas,true,definedConfig.Registry.Credit.Username,definedConfig.Registry.Credit.Password)
		requestParas = *requestParasPr
		delLayerUrl := registryUrl + "/v2/" + imageName + "/manifests/" + tmpDigest
		requestParas.Url = delLayerUrl
		requestParas.Method = http.MethodDelete
		_,_ = httpclient.SendRequest(&requestParas)
		rowNum := delTagsFromDB(tmpTagid,"","")
		
		if rowNum == 0 {
			return fmt.Errorf("no tag has be deleted")
		}
	}
	return nil
}

func (r RegistryCtl)layerDel(c *sysadmServer.Context,imageName string, digest string)error{

	if strings.TrimSpace(imageName) == "" || strings.TrimSpace(digest) == "" {
		return fmt.Errorf("can not delete layer with name of image or digest of layer is empty")
	}
	registryUrl  := getRegistryRootUrl(c)
	requestParas :=  httpclient.RequestParams{}
	
	definedConfig := RuntimeData.RuningParas.DefinedConfig
	requestParasPr,_ := httpclient.AddBasicAuthData(&requestParas,true,definedConfig.Registry.Credit.Username,definedConfig.Registry.Credit.Password)
	requestParas = *requestParasPr

	digestArray := strings.Split(digest, ",")
	for _,item := range digestArray {
		delLayerUrl := registryUrl + "/v2/" + imageName + "/blobs/" + item
		requestParas.Url = delLayerUrl
		requestParas.Method = http.MethodDelete
		_,_ = httpclient.SendRequest(&requestParas)
	}

	_ = delBlobFromDB("","",digest)
	
	return nil
}


/*
	getCount  get the total number of image acccording "imageid","projectid","name","ownerid"
	imageid: image "id, like id1,id2,id3.... 
	projectid: project id like "id1,id2,id3...." 
	name: image name. 
	ownerid：user id like "id1,id2,id3...." 
*/
func (r RegistryCtl)getCount(c *sysadmServer.Context) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	dataMap,err := utils.GetRequestData(c,[]string{"imageid","projectid","name","ownerid"})
	errs = append(errs,err...)

	imageid := strings.TrimSpace(dataMap["imageid"])
	projectid := strings.TrimSpace(dataMap["projectid"])
	name := strings.TrimSpace(dataMap["name"])
	ownerid  := strings.TrimSpace(dataMap["ownerid"])
	
	dataSet,err := getImageCountFromDB(imageid,projectid,name,ownerid)
	errs = append(errs,err...)
	err = apiutils.SendResponseForMap(c,dataSet)
	errs = append(errs,err...)

	return errs
}


/*
	ListImage list images information acccording "imageid","tagid","name","ownerid","start","numperpage"
	imageid: image "id, like id1,id2,id3.... 
	tagid: tag id like "id1,id2,id3...." 
	name: tag name. 
	ownerid： user id like "id1,id2,id3...." 
*/
func (r RegistryCtl)tagList(c *sysadmServer.Context) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	dataMap,err := utils.GetRequestData(c,[]string{"tagid","imageid","name","ownerid","start","numperpage"})
	errs = append(errs,err...)

	tagid := strings.TrimSpace(dataMap["tagid"])
	imageid := strings.TrimSpace(dataMap["imageid"])
	name := strings.TrimSpace(dataMap["name"])
	ownerid  := strings.TrimSpace(dataMap["ownerid"])
	startStr  := strings.TrimSpace(dataMap["start"])
	numperpageStr  := strings.TrimSpace(dataMap["numperpage"])
	start := 0
	numperpage := 0 
	if startStr != "" {
		start,_ = strconv.Atoi(startStr)
	}
	if numperpageStr != "" {
		numperpage,_ = strconv.Atoi(numperpageStr)
	}
	
	dataSet,err := getTagInfoFromDB(tagid,imageid,name,ownerid,"",start, numperpage)
	errs = append(errs,err...)
	err = apiutils.SendResponseForMap(c,dataSet)
	errs = append(errs,err...)

	return errs
}

func (r RegistryCtl)ActionNotFound(c *sysadmServer.Context,action string) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	msg := "request data error: module registryctl has not action " + action + "with POST method"
	errs = apiutils.SendResponseForErrorMessage(c,20600002,msg)

	return errs
}