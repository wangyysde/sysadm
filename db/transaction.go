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
 */

package db

import (
	"sysadm/sysadmerror"
)

/*
start DB transaction
*/
func Begin(e DbEntity) (*Tx, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	if e == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071001, "error", "DB entity is nil."))
		return nil, errs
	}

	dbConfig := e.GetDbConfig()
	dbConn := dbConfig.Connect

	tx, err := dbConn.Begin()
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071002, "error", "start a db transation error  %s.", err))
		return nil, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071003, "debug", "start a db transation."))
	return &Tx{Entity: e, Tx: tx}, errs
}

/*
InsertData： building query statement according to tb and data first, then add the operation of inserting data into DB to a transaction.
return 0 and []sysadmerror.Sysadmerror when any error was occurred.
otherwise return RowsAffected and  []sysadmerror.Sysadmerror
*/
func (t *Tx) InsertData(tb string, data FieldData) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	entity := t.Entity
	query, err := entity.BuildInsertQuery(tb, data)
	errs = append(errs, err...)
	if query == "" {
		return 0, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071019, "debug", "insert query %s", query))
	tx := t.Tx
	res, e := tx.Exec(query)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071004, "error", "insert data into db error: %s", e))
		return 0, errs
	}

	rowsInerted, e := res.RowsAffected()
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071005, "error", "no data has be insert into db %s", e))
		return 0, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071006, "debug", "data has be insert into db"))

	return int(rowsInerted), errs

}

/*
UpdateData building query statement according to tb and data first, then add the operation of update to a transaction.
return -1 and []sysadmerror.Sysadmerror when any error was occurred.
otherwise return RowsAffected and  []sysadmerror.Sysadmerror
*/
func (t *Tx) UpdateData(tb string, data FieldData, where map[string]string) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	entity := t.Entity
	query, err := entity.BuildUpdateQuery(tb, data, where)
	errs = append(errs, err...)
	if query == "" {
		return -1, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071013, "debug", "update query statement %s", query))
	errs = append(errs, err...)

	tx := t.Tx
	res, e := tx.Exec(query)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071007, "error", "update data error: %s", e))
		return -1, errs
	}

	rowsInerted, e := res.RowsAffected()
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071008, "error", "no data has be update %s", e))
		return -1, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071009, "debug", "data has be update"))

	return int(rowsInerted), errs
}

/*
DeleteData building query statement according to dd data first, then add the operation of delete to a transaction.
return -1 and []sysadmerror.Sysadmerror when any error was occurred.
otherwise return RowsAffected and  []sysadmerror.Sysadmerror
*/
func (t *Tx) DeleteData(dd *SelectData) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	entity := t.Entity
	query, err := entity.BuildDeleteQuery(dd)
	errs = append(errs, err...)
	if query == "" {
		return -1, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071020, "debug", "delete query %s", query))
	tx := t.Tx
	res, e := tx.Exec(query)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071010, "error", "delete data error: %s", e))
		return -1, errs
	}

	rowsInerted, e := res.RowsAffected()
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071011, "error", "delete data error %s", e))
		return -1, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1071012, "debug", "data has be deleted"))

	return int(rowsInerted), errs
}

/*
rollback aborts the transaction
*/
func (t *Tx) Rollback() error {

	tx := t.Tx
	return tx.Rollback()

}

/*
Commit commits the transaction
*/
func (t *Tx) Commit() error {

	tx := t.Tx
	return tx.Commit()

}
