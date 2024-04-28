/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"

	sysadmObjects "sysadm/objects/app"
)

func New() Syssetting {
	ret := Syssetting{}
	ret.Name = defaultObjectName
	ret.TableName = defaultTableName
	ret.PkName = defaultPkName
	return ret
}

func (s Syssetting) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("can not get sysSetting information with empty ID")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(s.TableName, s.PkName, id)
	if e != nil {
		return nil, e
	}

	settingData := SysSettingSchema{}
	e = sysadmObjects.Unmarshal(dbData, &settingData)

	return settingData, e
}

func (s Syssetting) GetObjectCount(searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &SysSettingSchema{}); !ok {
		return -1, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &SysSettingSchema{}); !ok {
		return -1, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	return sysadmObjects.GetObjectCount(s.TableName, s.PkName, searchContent, ids, searchKeys, conditions)
}

func (s Syssetting) GetObjectList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &SysSettingSchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error: %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &SysSettingSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &SysSettingSchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(s.TableName, s.PkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		sysSettingData := &SysSettingSchema{}
		if e := sysadmObjects.Unmarshal(v, sysSettingData); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *sysSettingData)
	}

	return tmpRes, nil
}

func (s Syssetting) AddObject(data interface{}) error {
	ssSchemaData, ok := data.(SysSettingSchema)
	if !ok {
		return fmt.Errorf("the data try to add is not valid")
	}

	insertData, e := sysadmObjects.Marshal(ssSchemaData)
	if e != nil {
		return e
	}

	return sysadmObjects.AddObject(s.TableName, s.PkName, insertData)
}

func (s Syssetting) AddObjectByTx(data interface{}) (map[string]interface{}, string, error) {
	addData := make(map[string]interface{}, 0)

	schemaData, ok := data.(SysSettingSchema)
	if !ok {
		return addData, "", fmt.Errorf("there is an error occurred when coverting data to sysSetting Schema data")
	}

	addData, e := sysadmObjects.Marshal(schemaData)
	if e != nil {
		return addData, "", e
	}

	return addData, s.TableName, nil
}

func (s Syssetting) GetObjectIDFieldName() (string, string, error) {
	return s.TableName, s.PkName, nil
}

func (s Syssetting) GetGlobalValueByKey(key string) (defaultValue, value []string, e error) {
	return s.GetValueByKey(key, "", SettingScopeGlobal)
}

func (s Syssetting) GetK8sValueByKey(key, k8sID string) (defaultValue, value []string, e error) {
	defaultValue, value = []string{}, []string{}
	e = nil
	k8sID = strings.TrimSpace(k8sID)
	if k8sID == "0" || k8sID == "" {
		e = fmt.Errorf("ID of K8S cluster is not valid")
		return
	}

	return s.GetValueByKey(key, k8sID, SettingScopeK8sCluster)
}

func (s Syssetting) GetNodeValueByKey(key, hostID string) (defaultValue, value []string, e error) {
	defaultValue, value = []string{}, []string{}
	e = nil
	hostID = strings.TrimSpace(hostID)
	if hostID == "0" || hostID == "" {
		e = fmt.Errorf("ID of host is not valid")
		return
	}

	return s.GetValueByKey(key, hostID, SettingScopeNode)
}

func (s Syssetting) GetProjectValueByKey(key, projectID string) (defaultValue, value []string, e error) {
	defaultValue, value = []string{}, []string{}
	e = nil
	projectID = strings.TrimSpace(projectID)
	if projectID == "0" || projectID == "" {
		e = fmt.Errorf("ID of project is not valid")
		return
	}

	return s.GetValueByKey(key, projectID, SettingScopeProject)
}

func (s Syssetting) GetUserGroupValueByKey(key, groupID string) (defaultValue, value []string, e error) {
	defaultValue, value = []string{}, []string{}
	e = nil
	groupID = strings.TrimSpace(groupID)
	if groupID == "0" || groupID == "" {
		e = fmt.Errorf("ID of user group is not valid")
		return
	}

	return s.GetValueByKey(key, groupID, SettingScopeUserGroup)
}

func (s Syssetting) GetUserValueByKey(key, userID string) (defaultValue, value []string, e error) {
	defaultValue, value = []string{}, []string{}
	e = nil
	userID = strings.TrimSpace(userID)
	if userID == "0" || userID == "" {
		e = fmt.Errorf("ID of user  is not valid")
		return
	}

	return s.GetValueByKey(key, userID, SettingScopeUser)
}

func (s Syssetting) GetValueByKey(key, objectID string, scope int) (defaultValue, value []string, e error) {
	defaultValue, value = []string{}, []string{}
	e = nil
	key = strings.TrimSpace(strings.ToLower(key))
	if key == "" {
		e = fmt.Errorf("key of setting item should not empty")
		return
	}

	conditions := make(map[string]string, 0)
	scopeStr := strconv.Itoa(scope)
	conditions["scope"] = "='" + scopeStr + "'"
	if scope != SettingScopeGlobal {
		conditions["objectID"] = "='" + objectID + "'"
	}
	conditions["key"] = "='" + key + "'"

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, err := sysadmObjects.ValidKeysInSchema(conditionKeys, &SysSettingSchema{}); !ok {
		e = errors.Wrap(err, "the keys of condition must be the object field name")
		return
	}

	dbData, err := sysadmObjects.GetObjectList(s.TableName, s.PkName, "", []string{}, []string{}, conditions, 0, 0, nil)
	if e != nil {
		e = errors.Wrap(err, "query data error")
		return
	}

	for _, v := range dbData {
		sysSettingData := &SysSettingSchema{}
		if err := sysadmObjects.Unmarshal(v, sysSettingData); err != nil {
			e = errors.Wrap(err, "unmarshal DB data to sysSetting Schema data error")
			return
		}
		defaultValue = append(defaultValue, sysSettingData.DefaultValue)
		value = append(value, sysSettingData.Value)
	}

	return

}

func (s Syssetting) GetCertAndKey(certType int) (string, string, error) {
	certKey, keyKey, e := getCertAndKeyName(certType)
	if e != nil {
		return "", "", e
	}

	_, certs, e := s.GetGlobalValueByKey(certKey)
	if e != nil {
		return "", "", e
	}

	_, keys, e := s.GetGlobalValueByKey(keyKey)
	if e != nil {
		return "", "", e
	}

	if len(certs) == 0 || len(keys) == 0 {
		return "", "", nil
	}

	if len(certs) > 1 || len(keys) > 1 {
		return "", "", fmt.Errorf("there are more than one cert or key of cert in the system.")
	}

	certBase64 := certs[0]
	keyBase64 := keys[0]

	cert, e := base64.StdEncoding.DecodeString(certBase64)
	if e != nil {
		return "", "", e
	}

	key, e := base64.StdEncoding.DecodeString(keyBase64)
	if e != nil {
		return "", "", e
	}

	return string(cert), string(key), nil
}

func (s Syssetting) SaveCertAndKey(certType int, cert, key, reason string) error {
	scope := SettingScopeGlobal
	certName, keyName, e := getCertAndKeyName(certType)
	if e != nil {
		return e
	}

	certBase64 := base64.StdEncoding.EncodeToString([]byte(cert))
	keyBase64 := base64.StdEncoding.EncodeToString([]byte(key))
	certSchemaData := SysSettingSchema{
		Scope:              scope,
		ObjectID:           "0",
		Key:                certName,
		Value:              certBase64,
		LastModifiedBy:     0,
		LastModifiedTime:   int(time.Now().Unix()),
		LastModifiedReason: reason,
	}
	tx, e := sysadmObjects.BeginTx(runData.dbConf.Entity, s)
	if e != nil {
		return e
	}
	e = tx.AddObject(certSchemaData)
	if e != nil {
		return e
	}

	keySchemmaData := SysSettingSchema{
		Scope:              scope,
		ObjectID:           "0",
		Key:                keyName,
		Value:              keyBase64,
		LastModifiedBy:     0,
		LastModifiedTime:   int(time.Now().Unix()),
		LastModifiedReason: reason,
	}
	e = tx.AddObject(keySchemmaData)
	if e != nil {
		tx.Rollback()
		return e
	}

	e = tx.Commit()
	if e != nil {
		tx.Rollback()
		return e
	}

	return nil
}

func getCertAndKeyName(certType int) (string, string, error) {
	certKey := ""
	keyKey := ""

	switch certType {
	case CertTypeCa:
		certKey = SettingKeyForCA
		keyKey = SettingKeyForCaKey
	case CertTypeApiServer:
		certKey = SettingKeyForApiServerCert
		keyKey = SettingKeyForApiServerCertKey
	case CertTypeAgent:
		certKey = SettingKeyForAgentCert
		keyKey = SettingKeyFroAgentCertKey
	default:
		return "", "", fmt.Errorf("type of certificate was not found")
	}

	return certKey, keyKey, nil
}