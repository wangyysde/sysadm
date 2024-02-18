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
	sysadmDB "sysadm/db"
)

func BeginTx(dbEntity sysadmDB.DbEntity, entity ObjectEntity) (ObjectTx, error) {
	objecttx := ObjectTx{}
	if dbEntity == nil {
		dbEntity = runData.dbConf.Entity
	}

	if dbEntity == nil {
		return objecttx, fmt.Errorf("transaction can not be begin for empty DB entity")
	}

	tx, e := sysadmDB.NewBegin(dbEntity)
	if e != nil {
		return objecttx, e
	}

	if entity != nil {
		objecttx.Entity = entity
	}
	objecttx.Tx = tx

	return objecttx, nil
}

func (o ObjectTx) AddObject(data interface{}) error {

	if o.Entity == nil {
		return fmt.Errorf("objection entity is nil")
	}

	dbData, tbName, e := o.Entity.AddObjectByTx(data)
	tbName = strings.TrimSpace(tbName)
	if tbName == "" {
		return fmt.Errorf("object table is nil")
	}

	if e != nil {
		return e
	}

	dbFieldData := sysadmDB.FieldData(dbData)
	tx := o.Tx
	if tx == nil {
		return fmt.Errorf("transaction has not began")
	}

	return tx.NewInsertData(tbName, dbFieldData)
}

func (o ObjectTx) UpdateObject(data interface{}, conditions map[string]string, tbName string) error {
	if o.Entity == nil {
		return fmt.Errorf("objection entity is nil")
	}

	dbData, e := Marshal(data)
	if e != nil {
		return e
	}

	tbName = strings.TrimSpace(tbName)
	if tbName == "" {
		return fmt.Errorf("object table is nil")
	}
	dbFieldData := sysadmDB.FieldData(dbData)
	tx := o.Tx
	if tx == nil {
		return fmt.Errorf("transaction has not began")
	}

	return tx.NewUpdateData(tbName, dbFieldData, conditions)
}

func (o ObjectTx) AddObjectWithMap(tbName string, data map[string]interface{}) error {
	if o.Tx == nil {
		return fmt.Errorf("transaction has not began")
	}

	tbName = strings.TrimSpace(tbName)
	if tbName == "" || len(data) < 1 {
		return fmt.Errorf("table name is empty or no data should be added")
	}

	dbFieldData := sysadmDB.FieldData(data)
	tx := o.Tx
	return tx.NewInsertData(tbName, dbFieldData)
}

func (o ObjectTx) Rollback() error {
	tx := o.Tx
	if tx == nil {
		return fmt.Errorf("transaction has not began")
	}

	return tx.NewRollback()
}

func (o ObjectTx) Commit() error {
	tx := o.Tx
	if tx == nil {
		return fmt.Errorf("transaction has not began")
	}

	return tx.Commit()
}

func (o ObjectTx) UpdateObjectNextID() error {
	if o.Tx == nil {
		return fmt.Errorf("transaction has not began")
	}

	if o.Entity == nil {
		return fmt.Errorf("objection entity is nil")
	}

	objEntity := o.Entity
	tbName, idField, e := objEntity.GetObjectIDFieldName()
	if e != nil {
		return e
	}

	filedData, where, e := prepareUpdateObjNextIDData(tbName, idField, o.Tx.Entity)
	if e != nil {
		return e
	}

	tx := o.Tx
	if tx == nil {
		return fmt.Errorf("transaction has not began")
	}

	return tx.NewUpdateData("ids", filedData, where)
}
