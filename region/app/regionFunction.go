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
	"strings"
	sysadmObjects "sysadm/objects/app"
)

func New() Region {
	ret := Region{}
	ret.Name = defaultObjectName
	ret.TableName = defaultTableName
	ret.PkName = defaultPkName
	return ret
}

func (r Region) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("can not get region information with empty ID")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(r.TableName, r.PkName, id)
	if e != nil {
		return nil, e
	}

	countryData := CountrySchema{}
	e = sysadmObjects.Unmarshal(dbData, &countryData)

	return countryData, e
}

func (r Region) GetProvinceByCode(code string) (interface{}, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("can not get province information with empty code")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(defaultProvinceTable, defaultProvincePkName, code)
	if e != nil {
		return nil, e
	}

	provinceData := ProvinceSchema{}
	e = sysadmObjects.Unmarshal(dbData, &provinceData)

	return provinceData, e
}

func (r Region) GetCityByCode(code string) (interface{}, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("can not get city information with empty code")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(defaultCityTable, defaultCityPkName, code)
	if e != nil {
		return nil, e
	}

	cityData := CitySchema{}
	e = sysadmObjects.Unmarshal(dbData, &cityData)

	return cityData, e
}

func (r Region) GetCountyByCode(code string) (interface{}, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("can not get county information with empty code")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(defaultCountyTable, defaultCountyPkName, code)
	if e != nil {
		return nil, e
	}

	countyData := CountySchema{}
	e = sysadmObjects.Unmarshal(dbData, &countyData)

	return countyData, e
}

func (r Region) GetObjectCount(searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &CountrySchema{}); !ok {
		return -1, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &CountrySchema{}); !ok {
		return -1, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	return sysadmObjects.GetObjectCount(r.TableName, r.PkName, searchContent, ids, searchKeys, conditions)
}

func (r Region) GetObjectList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &CountrySchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &CountrySchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &CountrySchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(r.TableName, r.PkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		country := &CountrySchema{}
		if e := sysadmObjects.Unmarshal(v, country); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *country)
	}

	return tmpRes, nil
}

func (r Region) GetProvinceList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &ProvinceSchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &ProvinceSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &ProvinceSchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(defaultProvinceTable, defaultProvincePkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		province := &ProvinceSchema{}
		if e := sysadmObjects.Unmarshal(v, province); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *province)
	}

	return tmpRes, nil
}

func (r Region) GetCityList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &CitySchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &CitySchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &CitySchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(defaultCityTable, defaultCityPkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		city := &CitySchema{}
		if e := sysadmObjects.Unmarshal(v, city); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *city)
	}

	return tmpRes, nil
}

func (r Region) GetCountyList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &CountySchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &CountySchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &CountySchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(defaultProvinceTable, defaultProvincePkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		county := &CountySchema{}
		if e := sysadmObjects.Unmarshal(v, county); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *county)
	}

	return tmpRes, nil
}

func (r Region) AddObject(data interface{}) error {
	osData, ok := data.(CountrySchema)
	if !ok {
		return fmt.Errorf("the data try to add is not valid")
	}

	insertData, e := sysadmObjects.Marshal(osData)
	if e != nil {
		return e
	}

	return sysadmObjects.AddObject(r.TableName, "", insertData)
}

func (r Region) AddObjectByTx(data interface{}) (map[string]interface{}, string, error) {
	addData := make(map[string]interface{}, 0)

	schemaData, ok := data.(CountrySchema)
	if !ok {
		return addData, "", fmt.Errorf("there is an error occurred when coverting data to OS Schema schema")
	}

	addData, e := sysadmObjects.Marshal(schemaData)
	if e != nil {
		return addData, "", e
	}

	return addData, r.TableName, nil
}

func (r Region) GetObjectIDFieldName() (string, string, error) {
	return r.TableName, r.PkName, nil
}
