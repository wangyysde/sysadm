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

errorCode: 204xxxx

*/

package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
	"github.com/wangyysde/sysadmServer"
)


func (y Yum)GetModuleName()string{
	return "yum"
}

func (y Yum)GetActionList()[]string{
	return yumActions
}

/*
   adding group handlers for yum. the path of this handlers for is /api/1.0/yum
   this function called in daemonServer
*/
func addYumHandlers(r *sysadmServer.Engine)(([]sysadmerror.Sysadmerror)){
	var errs []sysadmerror.Sysadmerror
	
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(2040001,"debug","add group handlers for  /api/1.0/yum" ))
		
	v1 := r.Group("/api/v1.0/yum")
	{
		v1.POST("/:action", apiV1YumPostHandlers)
	}
	
	return errs
}

func apiV1YumPostHandlers(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror
	entity := Yum{}

	action := strings.TrimSuffix(strings.TrimPrefix(c.Param("action"),"/"),"/")
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(2040002,"debug","handling the request for yum module  with action %s.",action))
	switch strings.ToLower(action){
	case "getosversion":
		err := entity.getOsVersionHandler(c)
		errs = append(errs, err...)
	case "getobject":
		err := entity.getObjectHandler(c)
		errs = append(errs, err...)
	case "yumlist":
		err := entity.getYumListHandler(c)
		errs = append(errs, err...)	
	case "getcount":
		err := entity.getYumCountHandler(c)
		errs = append(errs, err...)	
	default:
		err := apiutils.ActionNotFound(c,"yum",action,http.MethodPost)
		errs = append(errs,err...)
	}

	logErrors(errs)
}


/*
	getOsVersionHandler gets os version information from DB and response to the client.
*/
func (y Yum)getOsVersionHandler(c *sysadmServer.Context) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	osSets, err := 	getOsFromDB()
	errs= append(errs,err...)
	osVerSets, err := getOsVersionFromDB(osSets)
	errs = append(errs,err...)

	err = apiutils.SendResponseForMap(c,osVerSets)
	errs = append(errs,err...)

	return errs
}

/*
	getOsVersionHandler gets os version information from DB and response to the client.
*/
func (y Yum)getObjectHandler(c *sysadmServer.Context) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	typeSets, err := getType()
	errs= append(errs,err...)
	
	err = apiutils.SendResponseForMap(c,typeSets)
	errs = append(errs,err...)

	return errs
}

func getOsFromDB()([]map[string]interface{},[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var rets []map[string]interface{}

	selectData := db.SelectData{
		Tb: []string{"os"},
		OutFeilds: []string{"osID","name","description"},
	}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	errs = append(errs,err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
		return rets,errs
	} 

	// if the os information is not exist in DB 
	if len(retData) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2040003,"debug","no os information in DB"))
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


func getOsVersionFromDB(yumos []map[string]interface{})([]map[string]interface{},[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var rets []map[string]interface{}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	for _,value := range yumos {
		osid := value["osID"]
		whereMap :=  make(map[string]string,0)
		whereMap["osid"] = "=\"" + utils.Interface2String(osid) + "\""
		whereMap["typeID"] = "=1"  // os 
		selectData := db.SelectData{
			Tb: []string{"version"},
			OutFeilds: []string{"versionID","name","osid","description"},
			Where: whereMap,
		}
		retData,err := dbEntity.QueryData(&selectData)
		errs = append(errs,err...)
		if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
			return rets,errs
		}
		var verArray []map[string]interface{}
		for _,line := range retData {
			verLine := make(map[string]interface{})
			for k,v := range line {
				verLine[k] = v
			}
			verArray = append(verArray,verLine)
		}

		value["vers"] = verArray

		rets = append(rets, value)
	}

	return rets,errs
}

func getType()([]map[string]interface{},[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var rets []map[string]interface{}

	selectData := db.SelectData{
		Tb: []string{"type"},
		OutFeilds: []string{"typeID","name","comment"}, 
	}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	errs = append(errs,err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
		return rets,errs
	} 

	// if the type information is not exist in DB 
	if len(retData) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(2040003,"debug","no os information in DB"))
		return rets,errs
	}
	for _,line := range retData {
		lineType := make(map[string]interface{},0)
		lineType["typeID"] = utils.Interface2String(line["typeID"])
		lineType["name"] = utils.Interface2String(line["name"])
		lineType["comment"] = utils.Interface2String(line["comment"])
		rets = append(rets,lineType)
	}
	
	return rets,errs

}

/*
	getYumListHandler gets yum inforation from DB accroding "yumid","name","osid","typeid","kind","enabled","start","numperpage".
*/
func (y Yum)getYumListHandler(c *sysadmServer.Context) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	dataMap,err := utils.GetRequestData(c,[]string{"yumid","name","osid","typeid","kind","enabled","start","numperpage"})
	errs = append(errs,err...)

	yumid := utils.GetKeyData(dataMap, "yumid")
	name := utils.GetKeyData(dataMap, "name")
	osid := utils.GetKeyData(dataMap, "osid")
	typeid := utils.GetKeyData(dataMap, "typeid")
	kind := utils.GetKeyData(dataMap, "kind")
	enabled := utils.GetKeyData(dataMap, "enabled")
	start := utils.GetKeyData(dataMap, "start")
	numperpage := utils.GetKeyData(dataMap, "numperpage")

	dataSet,err := getYumListFromDB(yumid,name,osid,typeid,kind, enabled,start,numperpage)
	errs = append(errs,err...)
	err = apiutils.SendResponseForMap(c,dataSet)
	errs = append(errs,err...)

	return errs
}

/*
	getYumCountHandler get total number of  yum inforation from DB accroding "yumid","name","osid","typeid","kind","enabled","start","numperpage".
*/
func (y Yum)getYumCountHandler(c *sysadmServer.Context) ([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	dataMap,err := utils.GetRequestData(c,[]string{"yumid","name","osid","typeid","kind","enabled"})
	errs = append(errs,err...)

	yumid := utils.GetKeyData(dataMap, "yumid")
	name := utils.GetKeyData(dataMap, "name")
	osid := utils.GetKeyData(dataMap, "osid")
	typeid := utils.GetKeyData(dataMap, "typeid")
	kind := utils.GetKeyData(dataMap, "kind")
	enabled := utils.GetKeyData(dataMap, "enabled")

	dataSet,err := getCountFromDB(yumid,name,osid,typeid,kind, enabled)
	errs = append(errs,err...)
	err = apiutils.SendResponseForMap(c,dataSet)
	errs = append(errs,err...)

	return errs
}

/* 
	getImageInfoFromDB: get image information from DB server accroding to imageid,projectid,name, ownerid 
	return []map[string]string and []sysadmerror.Sysadmerror
*/
func getYumListFromDB(yumid string, name string, osid string,typeid string,kind string, enabled string ,start string, num string)([]map[string]interface{}, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var rets []map[string]interface{}

	whereMap := prepareWhereForListFromDB(yumid,name,osid,typeid,kind,enabled)

	var limit []int
	if strings.TrimSpace(num) != "" {
		numInt,err := strconv.Atoi(num)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(2040004,"error","internal error: convert string %s to int error %s",num,err))
		}else {
			if numInt > 0 {
				startInt,err := strconv.Atoi(start)
				if err != nil {
					startInt = 0
				}
				limit = append(limit,startInt)
				limit = append(limit,numInt)
			}
		}
	}
	
	selectData := db.SelectData{
		Tb: []string{"yum a","os b","version c","type d"},
		OutFeilds: []string{"a.yumid","a.name","a.osid","b.name as osName", "a.versionid","c.name as versionName", "a.typeid", "d.name as typeName","a.catalog","a.kind","base_url","enabled","gpgcheck","gpgkey"},
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
			lineData[k] = v
		}

		rets = append(rets,lineData)

	}

	return rets,errs
}

/* 
	prepareWhereForListFromDB: prepare where field for query yum infromation from DB accroding to yumid string, name string, osid string,typeid string,kind string, enabled string 
	return []map[string]string and []sysadmerror.Sysadmerror
*/
func prepareWhereForListFromDB(yumid string, name string, osid string,typeid string,kind string, enabled string )(map[string]string){
	
	whereMap :=  make(map[string]string,0)
	if strings.TrimSpace(yumid) != "" {
		whereMap["yumid"] = db.BuildWhereFieldExact(yumid)
	}

	if strings.TrimSpace(name) != "" {
		whereMap["name"] = db.BuildWhereFieldExact(name)
	}

	if strings.TrimSpace(osid) != "" {
		whereMap["osid"] = db.BuildWhereFieldExact(osid)
	}

	if strings.TrimSpace(typeid) != "" {
		whereMap["typeid"] = db.BuildWhereFieldExact(typeid)
	}

	if strings.TrimSpace(kind) != "" {
		whereMap["kind"] = db.BuildWhereFieldExact(kind)
	}

	whereMap["a.osid"] = "=b.osID"
	whereMap["a.versionid"] = "=c.versionID"
	whereMap["a.typedid"] = "=d.typeID"

	if strings.TrimSpace(enabled) != "" {
		if strings.ToLower(strings.TrimSpace(enabled)) == "true" || strings.ToLower(strings.TrimSpace(enabled)) == "yes" || strings.ToLower(strings.TrimSpace(enabled)) == "1" {
			whereMap["enabled"] = "1"
		} else {
			whereMap["enabled"] = "0"
		}
	}

	return whereMap
}

/* 
	getImageInfoFromDB: get image information from DB server accroding to imageid,projectid,name, ownerid 
	return []map[string]string and []sysadmerror.Sysadmerror
*/
func getCountFromDB(yumid string, name string, osid string,typeid string,kind string, enabled string)([]map[string]interface{}, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	var rets []map[string]interface{}

	whereMap := prepareWhereForListFromDB(yumid,name,osid,typeid,kind,enabled)

	selectData := db.SelectData{
		Tb: []string{"yum a","os b","version c","type d"},
		OutFeilds: []string{"count(a.yumid) as num"},
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
			lineData[k] = v
		}

		rets = append(rets,lineData)

	}

	return rets,errs
}