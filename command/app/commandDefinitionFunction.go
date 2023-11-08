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
	"strconv"
	"strings"
	"time"

	sysadmApiServerApp "sysadm/apiserver/app"
	sysadmObjects "sysadm/objects/app"
	sysadmUtils "sysadm/utils"
)

func (c Command) GetDefinitionListForCreateObj(objectName, dataFromObject string, osID, osversionid int) ([]interface{}, error) {
	var ret []interface{}

	objectName = strings.TrimSpace(objectName)
	dataFromObject = strings.TrimSpace(dataFromObject)

	if objectName == "" {
		return ret, fmt.Errorf("object name is nul")
	}

	if osID == 0 || osversionid == 0 {
		return ret, fmt.Errorf("os or os's version is nul")
	}
	conditions := make(map[string]string, 0)
	conditions["executionType"] = "='" + strconv.Itoa(int(ExecutionTypeAuto)) + "'"
	conditions["automationKind"] = "='" + strconv.Itoa(int(AutomationKindObjectCreate)) + "'"
	conditions["objectName"] = "='" + objectName + "'"
	if dataFromObject != "" {
		conditions["dataFromObject"] = "='" + dataFromObject + "'"
	}
	conditions["osID"] = "='" + strconv.Itoa(osID) + "'"
	conditions["osversionid"] = "='" + strconv.Itoa(osversionid) + "'"
	conditions["deprecated"] = "='" + strconv.Itoa(CommandDefinedUnDeprecated) + "'"

	var emptyString []string
	return c.GetObjectList("", emptyString, emptyString, conditions, 0, 0, nil)
}

func (c Command) AddCommandForHostByTx(tx sysadmObjects.ObjectTx, hostid uint, commandDefs CommandDefinedSchema) (string, string, error) {
	relatedObjPkValues := ""

	if tx.Tx == nil {
		return "", relatedObjPkValues, fmt.Errorf("transaction is not begin")
	}

	commandID, e := GenerateCommandID()
	if e != nil {
		return "", relatedObjPkValues, e
	}

	dependID := commandDefs.Dependent
	dependendID := "0"
	if dependID != 0 {
		dependCommand, e := c.GetObjectInfoByID(string(dependID))
		if e != nil {
			return "", relatedObjPkValues, e
		}
		dependCommandData, ok := dependCommand.(CommandDefinedSchema)
		if !ok {
			return "", relatedObjPkValues, fmt.Errorf("internal error")
		}

		pkValues := ""
		dependendID, pkValues, e = c.AddCommandForHostByTx(tx, hostid, dependCommandData)
		if e != nil {
			return "", relatedObjPkValues, e
		}

		if relatedObjPkValues != "" {
			relatedObjPkValues = relatedObjPkValues + "," + pkValues
		}
	}

	commandData := CommandSchema{
		CommandID:        commandID,
		DefinedID:        commandDefs.ID,
		DependendID:      dependendID,
		Type:             commandDefs.Type,
		TransactionScope: commandDefs.TransactionScope,
		MustParas:        commandDefs.MustParas,
		Command:          commandDefs.Command,
		HostID:           hostid,
		Crontab:          commandDefs.Crontab,
		Synchronized:     commandDefs.Synchronized,
		CreateTime:       int(time.Now().Unix()),
		Status:           int(sysadmApiServerApp.CommandStatusCreated),
	}
	addCommandData, e := sysadmObjects.Marshal(commandData)
	e = tx.AddObjectWithMap(defaultCommandTableName, addCommandData)
	if e != nil {
		return "", relatedObjPkValues, e
	}

	paraKind := commandDefs.ParaKind
	if paraKind == int(ParaKindNo) {
		return commandID, relatedObjPkValues, nil
	}

	pkValues := ""
	pkValues, e = c.AddParaForHostByTx(tx, hostid, commandID, commandDefs)
	if e != nil {
		return "", relatedObjPkValues, e
	}

	if relatedObjPkValues == "" {
		relatedObjPkValues = pkValues
	} else {
		relatedObjPkValues = relatedObjPkValues + "," + pkValues
	}
	return commandID, relatedObjPkValues, nil
}

func (c Command) AddParaForHostByTx(tx sysadmObjects.ObjectTx, hostid uint, commandID string, commandDefs CommandDefinedSchema) (string, error) {
	pkValues := ""

	if tx.Tx == nil {
		return pkValues, fmt.Errorf("transaction is not begin")
	}

	definedParas, e := GetParaDefinedListByCommandID(commandDefs.ID)
	if e != nil {
		return pkValues, e
	}

	for _, para := range definedParas {
		switch ParaKind(para.ParaKind) {
		case ParaKindFixed:
			paraData := CommandParametersSchema{
				Name:         para.Key,
				Value:        para.Value,
				CommandID:    commandID,
				ParaKind:     para.ParaKind,
				SubCommandID: "",
			}

			addParaData, e := sysadmObjects.Marshal(paraData)
			if e != nil {
				return pkValues, e
			}
			e = tx.AddObjectWithMap(defaultParasTableName, addParaData)
			if e != nil {
				return pkValues, e
			}
		case ParaKindObjFieldValue:
			objPkValue := para.Value
			objTbName := para.TableName
			objPkName := para.PkName
			fieldName := para.FieldName

			if pkValues == "" {
				pkValues = objPkValue
			} else {
				pkValues = pkValues + "," + objPkValue
			}
			conditions := make(map[string]string, 0)
			conditions[objPkName] = "='" + objPkValue + "'"
			var emptyString []string
			dbData, e := sysadmObjects.GetObjectList(objTbName, objPkName, "", emptyString, emptyString, conditions, 0, 0, nil)
			if e != nil {
				return pkValues, e
			}

			for _, line := range dbData {
				value, ok := line[fieldName]
				if !ok {
					return pkValues, fmt.Errorf("object %s has not attribute %s", objTbName, fieldName)
				}
				valueStr := sysadmUtils.Interface2String(value)
				paraData := CommandParametersSchema{
					Name:         para.Key,
					Value:        valueStr,
					CommandID:    commandID,
					ParaKind:     para.ParaKind,
					SubCommandID: "",
				}

				addParaData, e := sysadmObjects.Marshal(paraData)
				if e != nil {
					return pkValues, e
				}
				e = tx.AddObjectWithMap(defaultParasTableName, addParaData)
				if e != nil {
					return pkValues, e
				}
			}
		case ParaKindGetByCommand:
			defindedID := para.SubCommandID
			subCommandData, e := c.GetObjectInfoByID(string(defindedID))
			if e != nil {
				return pkValues, e
			}

			subCommandDefs, ok := subCommandData.(CommandDefinedSchema)
			if !ok {
				return pkValues, fmt.Errorf("can not convert data to CommandDefinedSchema")
			}

			subCommandID, tmpPkValues, e := c.AddCommandForHostByTx(tx, hostid, subCommandDefs)
			if e != nil {
				return pkValues, e
			}
			tmpPkValues = strings.TrimSpace(tmpPkValues)
			if tmpPkValues != "" {
				if pkValues == "" {
					pkValues = tmpPkValues
				} else {
					pkValues = pkValues + "," + tmpPkValues
				}
			}
			paraData := CommandParametersSchema{
				Name:         para.Key,
				Value:        "",
				CommandID:    commandID,
				ParaKind:     para.ParaKind,
				SubCommandID: subCommandID,
			}
			addParaData, e := sysadmObjects.Marshal(paraData)
			if e != nil {
				return pkValues, e
			}
			e = tx.AddObjectWithMap(defaultParasTableName, addParaData)
			if e != nil {
				return pkValues, e
			}
		default:
			return "", fmt.Errorf("kind of parameter's is not valid")
		}
	}

	return pkValues, nil
}

func GetParaDefinedListByCommandID(commandID uint) ([]CommandParasDefinedSchema, error) {
	var ret []CommandParasDefinedSchema
	if commandID < 0 {
		return ret, fmt.Errorf("command ID(defined) is not valid")
	}

	conditions := make(map[string]string, 0)
	conditions["commandID"] = "='" + strconv.Itoa(int(commandID)) + "'"
	var emptyString []string
	dbData, e := sysadmObjects.GetObjectList(defaultCommandParasDefinedTableName, defaultCommandParasDefinedPkName, "",
		emptyString, emptyString, conditions, 0, 0, nil)
	if e != nil {
		return ret, e
	}

	for _, line := range dbData {
		paraDefined := &CommandParasDefinedSchema{}
		if e := sysadmObjects.Unmarshal(line, paraDefined); e != nil {
			return ret, e
		}
		ret = append(ret, *paraDefined)
	}

	return ret, nil
}
