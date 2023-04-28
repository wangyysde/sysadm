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
	"net/http"
	"strconv"
	"strings"
	"time"

	"sysadm/db"
	"sysadm/httpclient"
	"sysadm/sysadmapi"
	"sysadm/sysadmerror"
	"sysadm/utils"
)


func updataImage(imageName string){
	var errs []sysadmerror.Sysadmerror

	// checking image name
	image,ok := processImages[imageName]
	if strings.TrimSpace(imageName) == "" || !ok {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2022001,"error","can not found image's information or image name is empty"))
		logErrors(errs)
		return 
	}

	// getting the userid 
	username := image.username
	if strings.TrimSpace(username) == "" {
		username = "admin"
	}
	definedConfig := RuntimeData.RuningParas.DefinedConfig 
	apiServerTls := definedConfig.Sysadm.Server.Tls
	apiServerAddress := definedConfig.Sysadm.Server.Host
	apiServerPort := definedConfig.Sysadm.Server.Port
	apiVersion := definedConfig.Sysadm.ApiVerion
	userid,err := getUserIdByUsername(apiServerTls,apiServerAddress,apiServerPort,apiVersion,username)
	errs = append(errs,err...)
	if userid == 0 {
		userid = 1
	}
	
	imgSets,err := getImageInfoFromDB("","",imageName,"",0,0)
	imageID := 0
	errs = append(errs,err...)
	// if the image named imageName is not exist, then add it into the DB
	if len(imgSets) < 1 {
		nameArray := strings.Split(imageName,"/")
		projectName := nameArray[0]
		projectid,err := getProjectIdByName(apiServerTls,apiServerAddress,apiServerPort,apiVersion,projectName)
		errs = append(errs,err...)
		if projectid == 0 {
			count,err := addProject(apiServerTls,apiServerAddress,apiServerPort,apiVersion,projectName,userid)
			errs = append(errs,err...)
			if count == 0 {
				logErrors(errs)
				return
			}
		}
		projectid,err = getProjectIdByName(apiServerTls,apiServerAddress,apiServerPort,apiVersion,projectName)
		errs = append(errs,err...)
		if projectid == 0 {
			logErrors(errs)
			return
		}

		imageID = addImageToDB(imageName,projectid,userid)
		if imageID == 0 {
			logErrors(errs)
			return 
		}
	} else {			// if the image named imageName is exist, then update the information for the image
		img := imgSets[0]
		imageIDStr := utils.Interface2String(img["imageid"])
		id,e := strconv.Atoi(imageIDStr)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(2022002,"error","can not get image ID %s",e))
			logErrors(errs)
			return 
		}
		imageID = id

		data := make(db.FieldData,0)
		data["tagsnum"] = "tagsnum + 1"
		data["lasttag"] =  "\"" + image.tag + "\""
		update_time := time.Now().Unix()
		data["update_time"] = update_time
		sizeStr := strconv.Itoa(int(image.size))
		data["size"] = "size + " + sizeStr 

		where := make(map[string]string,0)
		where["imageid"] = "=" + imageIDStr

		dbEntity := RuntimeData.RuningParas.DBConfig.Entity
		_,err := dbEntity.UpdateData("image",data,where)
		errs = append(errs,err...)
		if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
			logErrors(errs)
			return 
		} 
	}

	_ = delTagsFromDB("",image.tag,strconv.Itoa(imageID))

	_ = addTagsToDB(imageName,imageID,userid)
	
	lastTagId := getObjectMaxID("tag")
	if lastTagId == 0 {
		return 
	}

	blobs := image.blobs
	for _,b := range blobs {
		_ = addBlobToDB(b,lastTagId)
	}
}

func updatePulltimesForImage(imageid string,imageName string)  {
	var errs []sysadmerror.Sysadmerror

	where := make(map[string]string,0)
	if strings.TrimSpace(imageid) != "" {
		where["imageid"] = "=" + imageid
	}

	if strings.TrimSpace(imageName) != "" {
		where["name"] = "=" + imageName
	}

	data := make(db.FieldData,0)
	data["pulltimes"] = "pulltimes + 1"

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	_,err := dbEntity.UpdateData("image",data,where)
	errs = append(errs,err...)
	logErrors(errs)
	
}
	
func updatePulltimesForTag(tagid string, digest string)  {
	var errs []sysadmerror.Sysadmerror

	where := make(map[string]string,0)
	if strings.TrimSpace(tagid) != "" {
		where["tagid"] = "=" + tagid
	}

	if strings.TrimSpace(digest) != "" {
		where["digest"] = "=" + digest
	}

	data := make(db.FieldData,0)
	data["pulltimes"] = "pulltimes + 1"

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	_,err := dbEntity.UpdateData("tag",data,where)
	errs = append(errs,err...)
	logErrors(errs)
	
}

/* 
   getUserIdByUsername get userid by username from api server
   userid will be return if username was found , otherwise 0 will be returned 
*/
func getUserIdByUsername(tls bool, address string, port int, apiVersion string,username string)(int, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	userApiData,err := sysadmapi.BuildApiRequestData("user","getinfo",tls,address,port,apiVersion,http.MethodPost)
	errs = append(errs,err...)
	if userApiData == nil {
		return 0,errs
	}

	queryData := make(map[string]string,0)
	queryData["username"] = username
	err = sysadmapi.SetQueryData(userApiData,queryData)
	errs = append(errs,err...)
	body,e := httpclient.SendRequest(&userApiData.RequestParas)
	errs = append(errs,e...)
	ret, err := sysadmapi.ParseResponseBody(body)
	errs = append(errs, err...)
	if ret.Status {
		retsMap := ret.DataSet
		lineMap := retsMap[0]
		userid,err := strconv.Atoi(lineMap["userid"])
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(2022002,"error","userid %s can not convert to int:%s",lineMap["userid"],err))
			return 0,errs
		}
		return userid, errs
	}

	return 0,errs
}

/* 
   getProjectIdByName get projectid by projectname from api server
   projectid will be return if projectname was found , otherwise 0 will be returned 
*/
func getProjectIdByName(tls bool, address string, port int, apiVersion string,name string)(int, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	projectApiData,err := sysadmapi.BuildApiRequestData("project","getinfo",tls,address,port,apiVersion,http.MethodPost)
	errs = append(errs,err...)
	if projectApiData == nil {
		return 0,errs
	}

	queryData := make(map[string]string,0)
	queryData["name"] = name
	err = sysadmapi.SetQueryData(projectApiData,queryData)
	errs = append(errs,err...)
	body,e := httpclient.SendRequest(&projectApiData.RequestParas)
	errs = append(errs,e...)
	ret, err := sysadmapi.ParseResponseBody(body)
	errs = append(errs, err...)
	if ret.Status {
		retsMap := ret.DataSet
		lineMap := retsMap[0]
		projectid,err := strconv.Atoi(lineMap["projectid"])
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(2022003,"error","projectid %s can not convert to int:%s",lineMap["projectid"],err))
			return 0,errs
		}
		return projectid, errs
	}

	return 0,errs
}

/* 
   addProject add project information into the system via  api server
   projectid will be return if projectname was found , otherwise 0 will be returned 
*/
func addProject(tls bool, address string, port int, apiVersion string,name string,ownerid int)(int, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	projectApiData,err := sysadmapi.BuildApiRequestData("project","add",tls,address,port,apiVersion,http.MethodPost)
	errs = append(errs,err...)
	if projectApiData == nil {
		return 0,errs
	}

	queryData := make(map[string]string,0)
	queryData["name"] = name
	queryData["comment"] = "project has be added automatic when client pushed a image"
	queryData["ownerid"] = strconv.Itoa(ownerid)
	err = sysadmapi.SetQueryData(projectApiData,queryData)
	errs = append(errs,err...)
	body,e := httpclient.SendRequest(&projectApiData.RequestParas)
	errs = append(errs,e...)
	ret, err := sysadmapi.ParseResponseBody(body)
	errs = append(errs, err...)
	if ret.Status {
		return 1, errs
	}

	return 0,errs
}

/* 
	getImageInfoFromDB: get image information from DB server accroding to imageid,projectid,name, ownerid 
	return []map[string]string and []sysadmerror.Sysadmerror
*/
func getImageInfoFromDB(imageid string, projectid string, name string,ownerid string,start int, num int)([]map[string]interface{},[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var rets []map[string]interface{}

	// Qeurying data from DB
	whereMap :=  make(map[string]string,0)
	if imageid != "" {
		var ids = ""
		imageids := strings.Split(imageid, ",")
		if len(imageids) >1 {
			ids = " in ("
			first := true
			for _,id := range imageids {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\"" + id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + imageid + "\""
		}
		whereMap["imageid"] = ids
	}

	if projectid != "" {
		var ids = ""
		projectids := strings.Split(projectid, ",")
		if len(projectids) >1 {
			ids = " in ("
			first := true
			for _,id := range projectids {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\""+ id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + projectid + "\""
		}
		whereMap["projectid"] = ids
	}

	if name != "" {
		whereMap["name"] = " like \"%" + name +"%\""
	}

	if ownerid != "" {
		var ids = ""
		ownerids := strings.Split(ownerid, ",")
		if len(ownerids) >1 {
			ids = " in ("
			first := true
			for _,id := range ownerids {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\"" + id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + ownerid + "\""
		}
		whereMap["ownerid"] = ids
	}
	
	var limit []int
	if num > 0 {
		if start < 0 {
			start = 0
		}
		limit = append(limit,start)
		limit = append(limit,num)
	}

	selectData := db.SelectData{
		Tb: []string{"image"},
		OutFeilds: []string{"imageid","projectid","name", "ownerid","description","tagsnum","lasttag","architecture","pulltimes","creation_time","update_time","size"},
		Where: whereMap,
		Limit: limit,
	}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	errs = append(errs,err...)
	if retData == nil {
		return rets,errs
		
	} 

	
	for _,line := range retData {
		lineData := make(map[string]interface{},0)
		for k,v := range line {
			//value := utils.Interface2String(v)
			lineData[k] = v
		}

		rets = append(rets,lineData)

	}

	return rets,errs
}

/* 
	getImageInfoFromDB: get image information from DB server accroding to imageid,projectid,name, ownerid 
	return []map[string]string and []sysadmerror.Sysadmerror
*/
func getImageCountFromDB(imageid string, projectid string, name string,ownerid string)([]map[string]interface{},[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var rets []map[string]interface{}

	// Qeurying data from DB
	whereMap :=  make(map[string]string,0)
	if imageid != "" {
		var ids = ""
		imageids := strings.Split(imageid, ",")
		if len(imageids) >1 {
			ids = " in ("
			first := true
			for _,id := range imageids {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\"" + id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + imageid + "\""
		}
		whereMap["imageid"] = ids
	}

	if projectid != "" {
		var ids = ""
		projectids := strings.Split(projectid, ",")
		if len(projectids) >1 {
			ids = " in ("
			first := true
			for _,id := range projectids {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\""+ id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + projectid + "\""
		}
		whereMap["projectid"] = ids
	}

	if name != "" {
		whereMap["name"] = " like\" %" + name +"%\""
	}

	if ownerid != "" {
		var ids = ""
		ownerids := strings.Split(ownerid, ",")
		if len(ownerids) >1 {
			ids = " in ("
			first := true
			for _,id := range ownerids {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\"" + id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + ownerid + "\""
		}
		whereMap["ownerid"] = ids
	}
	

	selectData := db.SelectData{
		Tb: []string{"image"},
		OutFeilds: []string{"count(imageid) as num"},
		Where: whereMap,
	}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	errs = append(errs,err...)
	if retData == nil {
		return rets,errs
		
	} 

	
	for _,line := range retData {
		lineData := make(map[string]interface{},0)
		for k,v := range line {
			//value := utils.Interface2String(v)
			lineData[k] = v
		}

		rets = append(rets,lineData)

	}

	return rets,errs
}



/* 
	getTagInfoFromDB: get tag information from DB server accroding to tagid, imageid,name, ownerid 
	return []map[string]string and []sysadmerror.Sysadmerror
*/
func getTagInfoFromDB(tagid string,imageid string,name string,ownerid string,digest string, start int, num int)([]map[string]interface{},[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var rets []map[string]interface{}

	// Qeurying data from DB
	whereMap :=  make(map[string]string,0)
	if tagid != "" {
		var ids = ""
		tagids := strings.Split(tagid, ",")
		if len(tagids) >1 {
			ids = " in ("
			first := true
			for _,id := range tagids {
				if first {
					ids += id
					first = false
				} else {
					ids = ids + "," +id
				}
			}
			ids += ")"
		} else {
			ids = ids + " =" + tagid
		}
		whereMap["tagid"] = ids
	}

	if imageid != "" {
		var ids = ""
		imageids := strings.Split(imageid, ",")
		if len(imageids) >1 {
			ids = " in ("
			first := true
			for _,id := range imageids {
				if first {
					ids += id
					first = false
				} else {
					ids = ids + "," +id
				}
			}
			ids += ")"
		} else {
			ids = ids + " =" + imageid
		}
		whereMap["imageid"] = ids
	}

	
	if digest != "" {
		var ids = ""
		digests := strings.Split(digest, ",")
		if len(digests) >1 {
			ids = " in ("
			first := true
			for _,id := range digests {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else { 
					ids = ids + ",\"" +id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + " =\"" + digest + "\""
		}
		whereMap["digest"] = ids
	}

	if name != "" {
		whereMap["name"] = " like \"%" + name +"%\""
	}

	if ownerid != "" {
		var ids = ""
		ownerids := strings.Split(ownerid, ",")
		if len(ownerids) >1 {
			ids = " in ("
			first := true
			for _,id := range ownerids {
				if first {
					ids += id
					first = false
				} else {
					ids = ids + "," +id
				}
			}
			ids += ")"
		} else {
			ids = ids + " =" + ownerid
		}
		whereMap["ownerid"] = ids
	}

	var limit []int
	if num > 0 {
		if start < 0 {
			start = 0
		}
		limit = append(limit,start)
		limit = append(limit,num)
	}

	var order []db.OrderData 
	order = append(order, db.OrderData{Key: "name", Order: 1})
	selectData := db.SelectData{
		Tb: []string{"tag"},
		OutFeilds: []string{"tagid","imageid","name", "description","description","pulltimes","ownerid","creation_time","update_time","size","digest"},
		Where: whereMap,
		Order: order,
		Limit: limit,
	}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	errs = append(errs,err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
		return rets,errs
	} 

	// if the user is not exist in DB 
	if len(retData) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1040015,"debug","no data"))
		return rets,errs
	}
	for _,line := range retData {
		lineData := make(map[string]interface{},0)
		for k,v := range line {
		//	value := utils.Interface2String(v)
			lineData[k] = v
		}

		rets = append(rets,lineData)

	}
	
	return rets,errs

}

func getBlobInfoFromDB(blobid string, tagid string,digest string)([]map[string]interface{},[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var rets []map[string]interface{}

	// Qeurying data from DB
	whereMap :=  make(map[string]string,0)
	if blobid != "" {
		var ids = ""
		blobids := strings.Split(blobid, ",")
		if len(blobids) >1 {
			ids = " in ("
			first := true
			for _,id := range blobids {
				if first {
					ids += id
					first = false
				} else {
					ids = ids + "," +id
				}
			}
			ids += ")"
		} else {
			ids = ids + " =" + blobid
		}
		whereMap["blobid"] = ids
	}

	if tagid != "" {
		var ids = ""
		tagids := strings.Split(tagid, ",")
		if len(tagids) >1 {
			ids = " in ("
			first := true
			for _,id := range tagids {
				if first {
					ids += id
					first = false
				} else {
					ids = ids + "," +id
				}
			}
			ids += ")"
		} else {
			ids = ids + " =" + tagid
		}
		whereMap["tagid"] = ids
	}
	
	if digest != "" {
		var ids = ""
		digests := strings.Split(digest, ",")
		if len(digests) >1 {
			ids = " in ("
			first := true
			for _,id := range digests {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\"" +id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + digest + "\""
		}
		whereMap["digest"] = ids
	}

	selectData := db.SelectData{
		Tb: []string{"blob"},
		OutFeilds: []string{"blobid","tagid","digest", "size","creation_time","update_time"},
		Where: whereMap,
	}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	errs = append(errs,err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
		return rets,errs
	} 

	// if the blob is not exist in DB 
	if len(retData) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1040015,"debug","no data"))
		return rets,errs
	}
	for _,line := range retData {
		lineData := make(map[string]interface{},0)
		for k,v := range line {
			lineData[k] = v
		}

		rets = append(rets,lineData)

	}
	
	return rets,errs	
}

/*
	addImageToDB: Insert the data of image into the database.
	return the last imageid  if execute successfully otherwise return zero 
*/
func addImageToDB(imageName string, projectid int, ownerid int)(int){
	var errs []sysadmerror.Sysadmerror
	
	image := processImages[imageName]
	data := make(db.FieldData,0)
	data["projectid"] = projectid
	data["name"] = imageName 
	data["ownerid"] = ownerid
	data["description"] = ""
	data["tagsnum"] = 1 
	data["lasttag"] = image.tag 
	data["architecture"] = image.architecture 
	data["pulltimes"] = 1
	creation_time := time.Now().Unix()
	update_time := creation_time
	data["creation_time"] = creation_time 
	data["update_time"] = update_time
	data["size"] = image.size
	
	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	_,err := dbEntity.InsertData("image",data)
	errs = append(errs, err...)
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error"){
		return 0
	}
	
	lastID := getObjectMaxID("image")

	return lastID
}

/*
	addTagsToDB: Insert the data of tag into the database.
	return affected rows if execute successfully otherwise return zero 
*/
func addTagsToDB(imageName string, imageId int, ownerid int)(int){
	var errs []sysadmerror.Sysadmerror

	image := processImages[imageName]
	data := make(map[string]interface{},0)
	data["imageid"] = imageId
	data["name"] = image.tag 
	data["description"] = ""
	data["pulltimes"] = 1
	data["ownerid"] = ownerid
	creation_time := time.Now().Unix()
	update_time := creation_time
	data["creation_time"] = creation_time
	data["update_time"] = update_time
	data["size"] = image.size
	data["digest"] = image.digest

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	rows,err := dbEntity.InsertData("tag",data)
	errs = append(errs, err...)
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error"){
		rows = 0
	}

	return rows
}

/*
	addBlobToDB: Insert the data of blob into the database.
	return affected rows if execute successfully otherwise return zero 
*/
func addBlobToDB(blob blob, tagid int)int{
	var errs []sysadmerror.Sysadmerror

	data := make(map[string]interface{},0)
	data["tagid"] = tagid
	data["digest"] = blob.digest 
	data["size"] = blob.size
	creation_time := time.Now().Unix()
	update_time := creation_time
	data["creation_time"] = creation_time
	data["update_time"] = update_time

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	rows,err := dbEntity.InsertData("blob",data)
	errs = append(errs, err...)
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error"){
		rows = 0
	}

	return rows
}

/*
	delImagesFromDB: delete images from DB according  to imageid  or imageName
	return affected rows if executation is successful, otherwise returns zero
*/
func delImagesFromDB(imageid string, imageName string,)(int) {
	var errs []sysadmerror.Sysadmerror

	// Qeurying data from DB
	whereMap :=  make(map[string]string,0)
	if imageid != "" {
		var ids = ""
		imageids := strings.Split(imageid, ",")
		if len(imageids) >1 {
			ids = " in ("
			first := true
			for _,id := range imageids {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\"" +id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + imageid + "\""
		}
		whereMap["imageid"] = ids
	}

	// preparing tag name for where options
	if imageName != "" {
		whereMap["name"] = "=\"" + imageName + "\""
	}

	delData := db.SelectData{
		Tb: []string{"image"},
		OutFeilds: []string{},
		Where: whereMap,
	}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	rows,err := dbEntity.DeleteData(&delData)
	errs = append(errs, err...)
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error"){
		rows = 0
	}

	return int(rows)
}

/*
	delTagsFromDB: delete tags from DB according  to (tagid and imageid ) or (tagname and imageid)
	return affected rows if executation is successful, otherwise returns zero
*/
func delTagsFromDB(tagid string, tagname string, imageid string)(int) {
	var errs []sysadmerror.Sysadmerror

	// Qeurying data from DB
	whereMap :=  make(map[string]string,0)
	if tagid != "" {
		var ids = ""
		tagids := strings.Split(tagid, ",")
		if len(tagids) >1 {
			ids = " in ("
			first := true
			for _,id := range tagids {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\"" +id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + tagid + "\""
		}
		whereMap["tagid"] = ids
	}

	// preparing tag name for where options
	if tagname != "" {
		var ids = ""
		tagnames := strings.Split(tagname, ",")
		if len(tagnames) >1 {
			ids = " in ("
			first := true
			for _,id := range tagnames {
				if first {
					ids = ids + "\"" + id + "\"" 
					first = false
				} else {
					ids = ids + ",\"" + id + "\"" 
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + tagname + "\"" 
		}
		whereMap["name"] = ids
	}

	// preparing imageid for where options
	if imageid != "" {
		var ids = ""
		imageids := strings.Split(imageid, ",")
		if len(imageids) >1 {
			ids = " in ("
			first := true
			for _,id := range imageids {
				if first {
					ids = "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\"" + id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + imageid + "\""
		}
		whereMap["imageid"] = ids
	}


	delData := db.SelectData{
		Tb: []string{"tag"},
		OutFeilds: []string{},
		Where: whereMap,
	}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	rows,err := dbEntity.DeleteData(&delData)
	errs = append(errs, err...)
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error"){
		rows = 0
	}

	return int(rows)
}

/*
	delBlobFromDB: delete the layers of a image from DB according  to (blobid,tagid,or digest)
	both blobid,tagid and digest can be a string joined with ","
	return affected rows if executation is successful, otherwise returns zero
*/
func delBlobFromDB(blobid string, tagid string, digest string)(int) {
	var errs []sysadmerror.Sysadmerror

	// Qeurying data from DB
	whereMap :=  make(map[string]string,0) 
	if blobid != "" {
		var ids = ""
		blobids := strings.Split(blobid, ",")
		if len(blobids) >1 {
			ids = " in ("
			first := true
			for _,id := range blobids {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\"" +id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + blobid + "\""
		}
		whereMap["blobid"] = ids
	}

	if tagid != "" {
		var ids = ""
		tagids := strings.Split(tagid, ",")
		if len(tagids) >1 {
			ids = " in ("
			first := true
			for _,id := range tagids {
				if first {
					ids = ids + "\"" + id + "\""
					first = false
				} else {
					ids = ids + ",\"" +id + "\""
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + tagid + "\""
		}
		whereMap["tagid"] = ids
	}

	// preparing tag name for where options
	if digest != "" {
		var ids = ""
		digests := strings.Split(digest, ",")
		if len(digests) >1 {
			ids = " in ("
			first := true
			for _,id := range digests {
				if first {
					ids = ids + "\"" + id + "\"" 
					first = false
				} else {
					ids = ids + ",\"" + id + "\"" 
				}
			}
			ids += ")"
		} else {
			ids = ids + "=\"" + digest + "\"" 
		}
		whereMap["digest"] = ids
	}

	delData := db.SelectData{
		Tb: []string{"blob"},
		OutFeilds: []string{},
		Where: whereMap,
	}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	rows,err := dbEntity.DeleteData(&delData)
	errs = append(errs, err...)
	logErrors(errs)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error"){
		rows = 0
	}

	return int(rows)
}


/*
	getObjectMaxID: get maxid of the object from DB.
	return zero if any error occurs, otherwise return the maxid
*/
func getObjectMaxID(object string)(int){
	var errs []sysadmerror.Sysadmerror
	var tb []string 
	var outFeilds []string
	switch object {
		case "image":
			tb = append(tb,"image")
			outFeilds = append(outFeilds,"max(imageid) as id")
		case "tag":
			tb = append(tb,"tag")
			outFeilds = append(outFeilds,"max(tagid) as id")
		default:
			return 0
	}


	selectData := db.SelectData{
		Tb: tb,
		OutFeilds: outFeilds,
	}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	errs = append(errs,err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
		logErrors(errs)
		return 0
	} 

	lineData := retData[0]
	v, ok := lineData["id"]
	if !ok {
		return 0
	}
	value := utils.Interface2String(v)
	id,e := strconv.Atoi(value)
	if e != nil {
		return 0
	}

	return id
}