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

package db

import (
	"fmt"
)

// start DB transaction
func NewBegin(e DbEntity) (*Tx, error) {
	if e == nil {
		return nil, fmt.Errorf("DB entity is nil")
	}

	dbConfig := e.GetDbConfig()
	dbConn := dbConfig.Connect
	sqlTx, err := dbConn.Begin()
	if err != nil {
		return nil, err
	}

	dbTx := &Tx{Tx: sqlTx, Entity: e}
	return dbTx, nil

}

// NewInsertData building query statement according to tb and data first, then add the operation of inserting data into DB to a transaction.
// return error when any error was occurred. otherwise return nil
func (t *Tx) NewInsertData(tb string, data FieldData) error {

	if len(tb) < 1 {
		return fmt.Errorf("Table name(%s) is not valid.", tb)
	}

	if len(data) < 1 {
		return fmt.Errorf("error", "Can not insert empty data into table.")
	}

	entity := t.Entity
	query, err := entity.NewBuildInsertQuery(tb, data)
	if err != nil {
		return err
	}

	fmt.Printf("query statement: %s\n ", query)
	tx := t.Tx
	_, e := tx.Exec(query)
	if e != nil {
		return e
	}

	return nil

}

// NewUpdateData building query statement according to tb and data first, then add the operation of update to a transaction.
// return error when any error was occurred. otherwise return nil
func (t *Tx) NewUpdateData(tb string, data FieldData, where map[string]string) error {
	entity := t.Entity
	query, err := entity.NewBuildUpdateQuery(tb, data, where)
	if err != nil {
		return err
	}

	fmt.Printf("update statement: %s\n", query)
	tx := t.Tx
	_, e := tx.Exec(query)
	if e != nil {
		return e
	}

	return nil

}

// NewDeleteData building query statement according to dd data first, then add the operation of delete to a transaction.
// return error when any error was occurred.otherwise return nil
func (t *Tx) NewDeleteData(dd *SelectData) error {

	entity := t.Entity
	query, err := entity.NewBuildDeleteQuery(dd)
	if err != nil {
		return err
	}

	tx := t.Tx
	_, e := tx.Exec(query)
	return e
}

// NewRollback aborts the transaction
func (t *Tx) NewRollback() error {
	tx := t.Tx
	return tx.Rollback()
}

// NewCommit commits the transaction
func (t *Tx) NewCommit() error {
	tx := t.Tx
	return tx.Commit()
}
