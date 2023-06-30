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
	"strings"
	"sysadm/db"
	sysadmDB "sysadm/db"
)

// set dbConfig(*sysadmDB.DbConfig) the global variable runData
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

func InitObjEntity(objName string) (ObjectEntity, error) {
	objName = strings.TrimSpace(objName)
	if objName == "" {
		return nil, fmt.Errorf("object name is not valid")
	}

	var ret ObjectEntity = nil
	switch strings.ToLower(objName) {
	case "project":
		ret = Project{}
		ret.setDefaultForObject()
	}

	if ret == nil {
		return nil, fmt.Errorf("object %s was not found", objName)
	}

	return ret, nil
}

func getObjectInfoByID(tableName, pkName, id string) (db.FieldData, error) {
	tableName = strings.TrimSpace(tableName)
	pkName = strings.TrimSpace(pkName)
	id = strings.TrimSpace(id)
	if tableName == "" || pkName == "" || id == "" {
		return nil, fmt.Errorf("table name, field name of primary key or id is empty")
	}

	whereMap := make(map[string]string, 0)
	whereMap[pkName] = id
	selectData := db.SelectData{
		Tb:        []string{tableName},
		OutFeilds: []string{"*"},
		Where:     whereMap,
	}

	dbEntity := runData.dbConf.Entity
	dbData, _ := dbEntity.QueryData(&selectData)
	if dbData == nil || len(dbData) < 1 {
		return nil, fmt.Errorf("can not get object information")
	}

	return dbData[0], nil
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

		v, ok := data[tag]
		switch fieldType := field.Type.Kind(); fieldType {
		case reflect.Bool:
			value := false
			if ok {
				tmpV, tmpOk := v.(bool)
				if tmpOk && tmpV {
					value = true
				}
			}
			dV.Field(i).SetBool(value)
		case reflect.Int:
			var value = 0
			if ok {
				tmpV, tmpOk := v.(int)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetInt(int64(value))
		case reflect.Int8:
			var value int8 = 0
			if ok {
				tmpV, tmpOk := v.(int8)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetInt(int64(value))
		case reflect.Int16:
			var value int16 = 0
			if ok {
				tmpV, tmpOk := v.(int16)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetInt(int64(value))
		case reflect.Int32:
			var value int32 = 0
			if ok {
				tmpV, tmpOk := v.(int32)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetInt(int64(value))
		case reflect.Int64:
			var value int64 = 0
			if ok {
				tmpV, tmpOk := v.(int64)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetInt(value)
		case reflect.Uint:
			var value uint = 0
			if ok {
				tmpV, tmpOk := v.(uint)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetUint(uint64(value))
		case reflect.Uint8:
			var value uint8 = 0
			if ok {
				tmpV, tmpOk := v.(uint8)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetUint(uint64(value))
		case reflect.Uint16:
			var value uint16 = 0
			if ok {
				tmpV, tmpOk := v.(uint16)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetUint(uint64(value))
		case reflect.Uint32:
			var value uint32 = 0
			if ok {
				tmpV, tmpOk := v.(uint32)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetUint(uint64(value))
		case reflect.Uint64:
			var value uint64 = 0
			if ok {
				tmpV, tmpOk := v.(uint64)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetUint(value)
		case reflect.Float32:
			var value float32 = 0
			if ok {
				tmpV, tmpOk := v.(float32)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetFloat(float64(value))
		case reflect.Float64:
			var value float64 = 0
			if ok {
				tmpV, tmpOk := v.(float64)
				if tmpOk {
					value = tmpV
				}
			}
			dV.Field(i).SetFloat(value)
		case reflect.String:
			var value string = ""
			if ok {
				tmpV, tmpOk := v.(string)
				if tmpOk {
					value = tmpV
				}
			}
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
		default:
			continue
		}
		data[tag] = value
	}

	return data, nil
}
