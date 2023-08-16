/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
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
 */

package app

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sysadm/db"
	sysadmDB "sysadm/db"
	"sysadm/utils"
)

// SetRunDataForDBConf set dbConfig(*sysadmDB.DbConfig) the global variable runData
func SetRunDataForDBConf(dbConf *sysadmDB.DbConfig) error {
	if dbConf == nil {
		return fmt.Errorf("can not set nil to db config")
	}
	runData.dbConf = dbConf

	return nil
}

// SetWorkingRoot set working root path to global variable
func SetWorkingRoot(workingRoot string) error {
	workingRoot = strings.TrimSpace(workingRoot)
	if workingRoot == "" {
		return fmt.Errorf("can not set working root path with empty string")
	}

	runData.workingRoot = workingRoot

	return nil
}

// getObjectInfoByID get object information from DB by its id.
// this function should be called by an entity of an object
// the value of id should be the value of the primary key of object in DB
// success: return db.FieldData what can be Unmarshal and nil
// error: return nil and error
func GetObjectInfoByID(tableName, pkName, id string) (db.FieldData, error) {
	tableName = strings.TrimSpace(tableName)
	pkName = strings.TrimSpace(pkName)
	id = strings.TrimSpace(id)
	if tableName == "" || pkName == "" || id == "" {
		return nil, fmt.Errorf("table name, field name of primary key or id is empty")
	}

	whereMap := make(map[string]string, 0)
	whereMap[pkName] = "=" + id
	selectData := db.SelectData{
		Tb:        []string{tableName},
		OutFeilds: []string{"*"},
		Where:     whereMap,
	}

	dbEntity := runData.dbConf.Entity
	dbData, _ := dbEntity.NewQueryData(&selectData)
	if dbData == nil || len(dbData) < 1 {
		return nil, fmt.Errorf("can not get object information")
	}

	return dbData[0], nil
}

// getObjectCount get object count from db accroding to conditions
// this function should be called by an entity of an object
// success: return count of object and nil
// error: return -1 and an error
func GetObjectCount(tableName, pkName, searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	tableName = strings.TrimSpace(tableName)
	pkName = strings.TrimSpace(pkName)
	searchContent = strings.TrimSpace(searchContent)

	if tableName == "" || pkName == "" {
		return -1, fmt.Errorf("table name, field name of primary key or id is empty")
	}

	// don't change original conditions
	newConditions := make(map[string]string, 0)
	for k, v := range conditions {
		newConditions[k] = v
	}

	// preparing search condition
	searchSql := ""
	if len(searchKeys) > 0 && searchContent != "" {
		for _, k := range searchKeys {
			if searchSql == "" {
				searchSql = "=1 and (" + k + " like \"%" + searchContent + "%\""
			} else {
				searchSql = searchSql + " or " + k + " like \"%" + searchContent + "%\""
			}
		}
	}
	if searchSql != "" {
		searchSql = searchSql + ")"
		newConditions["1"] = searchSql
	}

	// preparing id list
	idsStr := strings.Join(ids, ",")
	idsStr = strings.TrimSpace(idsStr)
	if idsStr != "" {
		newConditions[pkName] = "in (" + idsStr + ")"
	}

	outSql := "count(" + pkName + ") as num"
	selectData := db.SelectData{
		Tb:        []string{tableName},
		OutFeilds: []string{outSql},
		Where:     newConditions,
	}
	dbEntity := runData.dbConf.Entity
	dbData, e := dbEntity.NewQueryData(&selectData)
	if e != nil || len(dbData) < 1 {
		return -1, fmt.Errorf("can not get object count")
	}

	row := dbData[0]
	numTmp, ok := row["num"]
	if !ok {
		return -1, fmt.Errorf("can not get object count")
	}

	num, e := utils.Interface2Int(numTmp)
	if e != nil {
		return -1, e
	}

	return num, nil
}

// getObjectList get object list from DB
// this function should be called by an entity of an object. searchKeys should be  the fields name of object table
// the value of key of conditions should be the fields name of object table, and the value of it should be a sql statement
// for where.
func GetObjectList(tableName, pkName, searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]map[string]interface{}, error) {
	var ret []map[string]interface{}

	tableName = strings.TrimSpace(tableName)
	pkName = strings.TrimSpace(pkName)
	searchContent = strings.TrimSpace(searchContent)
	if tableName == "" || pkName == "" {
		return ret, fmt.Errorf("table name, field name of primary key or id is empty")
	}

	// don't change original conditions
	newConditions := make(map[string]string, 0)
	for k, v := range conditions {
		newConditions[k] = v
	}

	// preparing search condition
	searchSql := ""
	if len(searchKeys) > 0 && searchContent != "" {
		for _, k := range searchKeys {
			if searchSql == "" {
				searchSql = "=1 and (" + k + " like \"%" + searchContent + "%\""
			} else {
				searchSql = searchSql + " or " + k + " like \"%" + searchContent + "%\""
			}
		}
	}

	if searchSql != "" {
		searchSql = searchSql + ")"
		newConditions["1"] = searchSql
	}

	// preparing id list
	idsStr := strings.Join(ids, ",")
	idsStr = strings.TrimSpace(idsStr)
	if idsStr != "" {
		newConditions[pkName] = "in (" + idsStr + ")"
	}

	var sqlLimit []int
	if step > 0 {
		sqlLimit = append(sqlLimit, startPos)
		sqlLimit = append(sqlLimit, step)
	}

	// preparing order field
	var order []db.OrderData
	for k, v := range orders {
		// 0 for ASC 1 for DESC
		directStr := 0
		if v == "1" {
			directStr = 1
		}
		line := db.OrderData{Key: k, Order: directStr}
		order = append(order, line)
	}

	selectData := db.SelectData{
		Tb:        []string{tableName},
		OutFeilds: []string{"*"},
		Where:     newConditions,
		Limit:     sqlLimit,
		Order:     order,
	}

	dbEntity := runData.dbConf.Entity
	dbData, e := dbEntity.NewQueryData(&selectData)

	if e != nil {
		return ret, fmt.Errorf("can not get object list. error %s", e)
	}

	return dbData, nil
}

// validKeysInSchema check the keys are the fields of obj
// success: return true
// error: return false
func ValidKeysInSchema(keys []string, obj any) (bool, error) {
	if len(keys) < 1 {
		return true, nil
	}

	oT := reflect.TypeOf(obj)
	if oT.Kind() != reflect.Pointer || (oT.Elem().Kind() != reflect.Struct) {
		return false, fmt.Errorf("the object is not a pointer or the destination of the pointer is not struct")
	}

	oTElem := oT.Elem()
	var objFields []string
	for i := 0; i < oTElem.NumField(); i++ {
		field := oTElem.Field(i)
		if !field.IsExported() {
			continue
		}
		tag, okTag := field.Tag.Lookup("db")
		if !okTag || tag == "" {
			continue
		}
		objFields = append(objFields, tag)
	}

	for _, k := range keys {
		found := false
		for _, o := range objFields {
			if k == o {
				found = true
				break
			}
		}

		if !found {
			return false, fmt.Errorf("object fields are: %+v  keys are %+v ", objFields, keys)
		}
	}

	return true, nil
}

func Unmarshal(data map[string]interface{}, dst any) error {
	dT := reflect.TypeOf(dst)
	if dT.Kind() != reflect.Pointer || dT.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("the type of dst is not a pointer or the destination where dst point to is not a struct")
	}

	dTElem := dT.Elem()
	dV := reflect.ValueOf(dst).Elem()

	for i := 0; i < dTElem.NumField(); i++ {
		field := dTElem.Field(i)
		if !field.IsExported() {
			continue
		}

		tag, okTag := field.Tag.Lookup("db")
		if !okTag || tag == "" {
			continue
		}

		v, _ := data[tag]
		switch fieldType := field.Type.Kind(); fieldType {
		case reflect.Bool:
			value := utils.Interface2Bool(v)
			dV.Field(i).SetBool(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var value = int64(0)
			vStr := utils.Interface2String(v)
			if vStr != "" {
				tmpValue, e := strconv.ParseInt(vStr, 10, 64)
				if e == nil {
					value = tmpValue
				} else {
					return fmt.Errorf("can not umarshal feild %s for %s", tag, e)
				}
			}

			dV.Field(i).SetInt(value)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var value uint64 = 0
			valueStr := utils.Interface2String(v)
			if valueStr != "" {
				tmpValue, e := strconv.ParseUint(valueStr, 10, 64)
				if e == nil {
					value = tmpValue
				} else {
					return fmt.Errorf("can not umarshal feild %s for %s", tag, e)
				}
			}
			dV.Field(i).SetUint(value)
		case reflect.Float32, reflect.Float64:
			var value float64 = 0
			vStr := utils.Interface2String(v)
			if vStr != "" {
				tmpValue, e := strconv.ParseFloat(vStr, 10)
				if e == nil {
					value = tmpValue
				} else {
					return fmt.Errorf("can not umarshal feild %s for %s", tag, e)
				}
			}
			dV.Field(i).SetFloat(value)
		case reflect.String:
			value := utils.Interface2String(v)
			dV.Field(i).SetString(value)
		default:
			continue
		}
	}

	return nil
}

func Marshal(s any) (map[string]interface{}, error) {
	sT := reflect.TypeOf(s)
	var sV reflect.Value
	if sT.Kind() == reflect.Pointer {
		sT = sT.Elem()
		sV = reflect.ValueOf(s).Elem()
	} else {
		sV = reflect.ValueOf(s)
	}
	if sT.Kind() != reflect.Struct {
		return nil, fmt.Errorf("we can only marshal struct to map")
	}

	data := make(map[string]interface{}, 0)
	for i := 0; i < sT.NumField(); i++ {
		field := sT.Field(i)
		if !field.IsExported() {
			continue
		}

		tag, okTag := field.Tag.Lookup("db")
		if !okTag || tag == "" {
			continue
		}

		var value interface{} = nil
		switch fieldType := field.Type.Kind(); fieldType {
		case reflect.Bool:
			value = sV.Field(i).Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = sV.Field(i).Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = sV.Field(i).Uint()
		case reflect.Float32, reflect.Float64:
			value = sV.Field(i).Float()
		case reflect.Complex64, reflect.Complex128:
			value = sV.Field(i).Complex()
		case reflect.String:
			value = sV.Field(i).String()
			if strings.TrimSpace(value.(string)) == "" {
				continue
			}
		default:
			continue
		}
		data[tag] = value
	}

	return data, nil
}

// AddObject insert  the data of object into the DB
// return error if any error has occurred. otherwise return nil
func AddObject(tableName, pkName string, data map[string]interface{}) error {
	tableName = strings.TrimSpace(tableName)

	// if primary key of the object is auto increment key, then we should delete the item
	pkName = strings.TrimSpace(pkName)
	if pkName != "" {
		pkValue, ok := data[pkName]
		if ok {
			pkValueStr := utils.Interface2String(pkValue)
			if strings.TrimSpace(pkValueStr) == "0" {
				delete(data, pkName)
			}
		}
	}
	dbData := db.FieldData(data)
	dbEntity := runData.dbConf.Entity

	return dbEntity.NewInsertData(tableName, dbData)
}
