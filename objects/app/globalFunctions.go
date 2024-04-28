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
	runtime "sysadm/apimachinery/runtime/v1beta1"
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

func GetRunDataForDBConf() *sysadmDB.DbConfig {
	return runData.dbConf
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

func GetWorkingRoot() string {
	return strings.TrimSpace(runData.workingRoot)
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
	whereMap[pkName] = "='" + id + "'"
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

// GetObjectCount get object count from db accroding to conditions
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

// GetObjectList get object list from DB
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

func GetObjectMaxID(dbEntity sysadmDB.DbEntity, tableName, idField string) (uint, error) {
	tableName = strings.TrimSpace(tableName)
	idField = strings.TrimSpace(idField)

	if tableName == "" || idField == "" {
		return 0, fmt.Errorf("table name, field name of id is empty")
	}

	if dbEntity == nil {
		dbEntity = runData.dbConf.Entity
	}

	if dbEntity == nil {
		return 0, fmt.Errorf("DB Entity is nil")
	}

	sql := "max(" + idField + ") as id"
	selectData := db.SelectData{
		Tb:        []string{tableName},
		OutFeilds: []string{sql},
	}
	dbData, e := dbEntity.NewQueryData(&selectData)
	if e != nil || len(dbData) < 1 {
		return 0, fmt.Errorf("there is an error occurred when building SQL statement")
	}

	row := dbData[0]
	numTmp, ok := row["id"]
	if !ok {
		return 0, fmt.Errorf("there is an error occurred when getting object ID")
	}

	id, e := utils.Interface2Uint64(numTmp)
	if e != nil {
		return 0, e
	}

	return uint(id), nil
}

func prepareUpdateObjNextIDData(tbName, idField string, dbEntity sysadmDB.DbEntity) (sysadmDB.FieldData, map[string]string, error) {
	updateData := make(sysadmDB.FieldData, 0)
	where := make(map[string]string, 0)

	tbName = strings.TrimSpace(tbName)
	idField = strings.TrimSpace(idField)
	if tbName == "" || idField == "" {
		return updateData, where, fmt.Errorf("table name or field name of id is empty")
	}

	if dbEntity == nil {
		dbEntity = runData.dbConf.Entity
	}

	if dbEntity == nil {
		return updateData, where, fmt.Errorf("DB Entity is nil")
	}

	updateData["nextValue"] = "nextValue + 1"

	where["tableName"] = tbName
	where["fieldName"] = idField

	return updateData, where, nil
}

func UpdateObjectNextID(dbEntity sysadmDB.DbEntity, tbName, idField string) error {
	fieldData, where, e := prepareUpdateObjNextIDData(tbName, idField, dbEntity)
	if e != nil {
		return e
	}

	return dbEntity.NewUpdateData(tbName, fieldData, where)
}

func GetCommandRelatedObjectList() ([]interface{}, error) {
	var ret []interface{}

	// don't change original conditions
	conditions := make(map[string]string, 0)
	conditions["isCommandRelated"] = "=1"
	conditions["deprecated"] = "=0"
	selectData := db.SelectData{
		Tb:        []string{DefautlObjectInfoTable},
		OutFeilds: []string{"*"},
		Where:     conditions,
	}

	dbEntity := runData.dbConf.Entity
	dbData, e := dbEntity.NewQueryData(&selectData)

	if e != nil {
		return ret, fmt.Errorf("can not get object list. error %s", e)
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		command := &ObjectInfoSchema{}
		if e := Unmarshal(v, command); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *command)
	}

	return tmpRes, nil
}

// CreateGetCondition build condition(map[string]string) for get resources from DB.
// obj is the type of resource,queryData is a map[string][]string with key is the value of the tag of the field and value
// is a []string
func CreateGetCondition(obj reflect.Type, queryData runtime.RequestQuery) (map[string]string, error) {
	if len(queryData) < 1 {
		return nil, fmt.Errorf("no request query data")
	}

	if obj.Kind() == reflect.Pointer {
		obj = obj.Elem()
	}
	if obj.Kind() != reflect.Struct {
		return nil, fmt.Errorf("object is not a valid resource")
	}

	condition := make(map[string]string, 0)
	for k, q := range queryData {
		if len(q) < 1 {
			continue
		}
		objElem := obj.Elem()
		for i := 0; i < objElem.NumField(); i++ {
			field := objElem.Field(i)
			if !field.IsExported() {
				continue
			}
			tag, okTag := field.Tag.Lookup("db")
			if !okTag || tag == "" {
				continue
			}
			if strings.TrimSpace(k) == strings.TrimSpace(tag) {
				if len(q) == 1 {
					condition[tag] = "='" + q[0] + "'"
				} else {
					qStr := ""
					for _, qv := range q {
						if qStr == "" {
							qStr = "'" + qv + "'"
						} else {
							qStr = qStr + ",'" + qv + "'"
						}
					}
					condition[tag] = "in (" + qStr + ")"
				}
			}
		}
	}

	if len(condition) < 1 {
		return nil, fmt.Errorf("request query data is not valid")
	}

	return condition, nil
}

func GetResource(gvk runtime.GroupVersionKind, obj reflect.Type, condition map[string]string) ([]interface{}, error) {
	tbName := getResourceTableName(gvk)
	dbData, e := getResourceFromDB(tbName, condition)
	if e != nil {
		return nil, e
	}

	var ret = make([]interface{}, 0)
	for _, line := range dbData {
		objValue := reflect.Zero(obj)
		e := Unmarshal(line, &objValue)
		if e != nil {
			return nil, e
		}
		ret = append(ret, &objValue)
	}

	return ret, nil
}

func getResourceTableName(gvk runtime.GroupVersionKind) string {
	group := gvk.Group
	groupUnderline := strings.Replace(group, ".", "_", -1)
	tbName := strings.ToLower("object_" + groupUnderline + "_" + gvk.Kind)

	return tbName
}

// GetFeildValueByName get a field name from a struct by the field name.
// return nil and an error if there is not a field named fieldname in the struct
func GetFeildValueByName(data any, fieldName string) (interface{}, error) {
	dT := reflect.TypeOf(data)
	dV := reflect.ValueOf(data)
	if dT.Kind() == reflect.Pointer {
		dT = dT.Elem()
		dV = dV.Elem()
	}
	if dT.Kind() != reflect.Struct {
		return nil, fmt.Errorf("we can not only get a struct field value")
	}

	for i := 0; i < dT.NumField(); i++ {
		field := dT.Field(i)
		if field.Name == fieldName {
			fieldValue := dV.Field(i)
			return fieldValue.Interface(), nil
		}
	}

	return nil, fmt.Errorf("field named %s was not found in struct %s", fieldName, dT.Name())
}

// SetFeildValueByName set value to a field named fieldName in a struct.
// data must be a pointer point to a struct.
// return an error if the field named fieldName is not find in the struct or the type of the value is not
// equal to the type of field.
func SetFeildValueByName(data, value any, fieldName string) error {
	dT := reflect.TypeOf(data)
	if dT.Kind() != reflect.Pointer || dT.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("data must be a pointer point to a struct")
	}
	fieldName = strings.TrimSpace(fieldName)
	if fieldName == "" {
		return fmt.Errorf("field name must not empty")
	}

	dTElem := dT.Elem()
	dV := reflect.ValueOf(data).Elem()
	for i := 0; i < dTElem.NumField(); i++ {
		field := dTElem.Field(i)
		if field.Name == fieldName {
			fieldType := field.Type
			valueType := reflect.TypeOf(value)
			if fieldType != valueType {
				return fmt.Errorf("the type of value is not equal to the type of the field of the data")
			}
			tmpValue := reflect.ValueOf(value)
			dV.Field(i).Set(tmpValue)
			return nil
		}
	}

	return fmt.Errorf("field named %s was not found in struct %s", fieldName, dTElem.Name())
}

func SetFieldValueZeroByName(data any, fieldName string) error {
	dT := reflect.TypeOf(data)
	if dT.Kind() != reflect.Pointer || dT.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("data must be a pointer point to a struct")
	}
	fieldName = strings.TrimSpace(fieldName)
	if fieldName == "" {
		return fmt.Errorf("field name must not empty")
	}

	dTElem := dT.Elem()
	dV := reflect.ValueOf(data).Elem()
	for i := 0; i < dTElem.NumField(); i++ {
		field := dTElem.Field(i)
		if field.Name == fieldName {
			fieldType := field.Type
			value := reflect.ValueOf(fieldType)
			dV.Field(i).Set(value)
			return nil
		}
	}

	return fmt.Errorf("field named %s was not found in struct %s", fieldName, dTElem.Name())
}

func createGetRelationResourceCondition(ids []interface{}) (map[string]string, error) {
	condition := make(map[string]string, 0)

	if len(ids) < 1 {
		return nil, fmt.Errorf("parent ID should not be empty")
	}

	if len(ids) == 1 {
		condition[runtime.ResourceRelationParentDBFieldName] = "'" + utils.Interface2String(ids[0]) + "'"
		return condition, nil
	}

	idsStr := ""
	for _, id := range ids {
		if idsStr == "" {
			idsStr = "in('" + utils.Interface2String(id) + "'"
		} else {
			idsStr = idsStr + ",'" + utils.Interface2String(id) + "'"
		}
	}
	idsStr = idsStr + ")"
	condition[runtime.ResourceRelationParentDBFieldName] = idsStr

	return condition, nil
}

func buildConditionByIds(ids []interface{}) (map[string]string, error) {
	condition := make(map[string]string, 0)
	if len(ids) < 1 {
		return nil, fmt.Errorf("resource ID should not be empty")
	}
	if len(ids) == 1 {
		condition[runtime.ResourcepKDbFieldName] = "'" + utils.Interface2String(ids[0]) + "'"
		return condition, nil
	}

	idsStr := ""
	for _, id := range ids {
		if idsStr == "" {
			idsStr = "in ('" + utils.Interface2String(id) + "'"
		} else {
			idsStr = idsStr + ",'" + utils.Interface2String(id) + "'"
		}
	}
	idsStr = idsStr + ")"
	condition[runtime.ResourcepKDbFieldName] = idsStr

	return condition, nil
}

func GetRelatedResource(parentGvk, childGvk runtime.GroupVersionKind, childObj reflect.Type, ids []interface{}) ([]interface{}, error) {
	relationTable := getResourceRelationTableName(parentGvk, childGvk)
	condition, e := createGetRelationResourceCondition(ids)
	if e != nil {
		return nil, e
	}
	dbData, e := getResourceFromDB(relationTable, condition)
	if e != nil {
		return nil, e
	}
	childIds := make([]interface{}, 0)
	for _, line := range dbData {
		v, ok := line[runtime.ResourceRelationChildDBFieldName]
		if !ok {
			return nil, fmt.Errorf("no filed name %s in table %s", runtime.ResourceRelationChildDBFieldName, relationTable)
		}
		childIds = append(childIds, v)
	}
	childCondition, e := buildConditionByIds(childIds)
	if e != nil {
		return nil, e
	}

	return GetResource(childGvk, childObj, childCondition)
}

func getResourceRelationTableName(parentGvk, childGvk runtime.GroupVersionKind) string {
	parentGroup := parentGvk.Group
	pG := strings.Replace(parentGroup, ".", "_", -1)
	childGroup := childGvk.Group
	cG := strings.Replace(childGroup, ".", "_", -1)

	return "relation_" + pG + "_" + parentGvk.Kind + "_to_" + cG + "_" + childGvk.Kind
}

func ReplacePointerWithStructForSlice(src []interface{}) ([]interface{}, error) {
	if src == nil {
		return nil, fmt.Errorf("source data is nil")
	}
	if len(src) < 1 {
		return nil, fmt.Errorf("source data is empty")
	}
	ret := make([]interface{}, 0)
	for _, v := range src {
		srcT := reflect.TypeOf(v)
		if srcT.Kind() != reflect.Pointer {
			return src, nil
		}
		srcD := reflect.ValueOf(v).Elem().Interface()
		ret = append(ret, srcD)
	}

	return ret, nil
}

func GetResourceReferenceTablesName(gvk runtime.GroupVersionKind) string {
	if gvk.Group == "" || gvk.Kind == "" {
		return ""
	}
	group := gvk.Group
	g := strings.Replace(group, ".", "_", -1)

	return "reference_" + g + "_" + gvk.Kind + "_objects"
}

// GetReferencedResources get referenced resources from DB by resource ID
// return nil and error if any error occurred. otherewise []runtime.ReferenceInfo and nil
func GetReferencedResources(tableName string, id int) ([]runtime.ReferenceInfo, error) {
	tableName = strings.TrimSpace(tableName)
	if tableName == "" {
		return nil, fmt.Errorf("table name of reference resource must not empty")
	}
	condition := make(map[string]string, 0)
	condition[runtime.ResourceReferenceDBObjectIdFieldName] = "='" + strconv.Itoa(id) + "'"

	dbData, e := getResourceFromDB(tableName, condition)
	if e != nil {
		return nil, e
	}

	var ret = make([]runtime.ReferenceInfo, 0)
	for _, line := range dbData {
		objValue := runtime.ReferenceInfo{}
		e := Unmarshal(line, &objValue)
		if e != nil {
			return nil, e
		}
		ret = append(ret, objValue)
	}

	return ret, nil
}
