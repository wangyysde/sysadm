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
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"sysadm/utils"
)

// InsertData build insert SQL statement and execute a query using the SQL statement.
// return error if teh SQL statement is be execute successful.
// Or return nil
func (p MySQL) NewInsertData(tb string, data FieldData) error {

	if len(tb) < 1 {
		return fmt.Errorf("Table name(%s) is not valid.", tb)
	}

	if len(data) < 1 {
		return fmt.Errorf("error", "Can not insert empty data into table.")
	}

	insertStr := "INSERT INTO `" + tb + "`("
	valueStr := "Values ("
	i := 1
	for key, value := range data {
		if i == 1 {
			insertStr = insertStr + "`" + key + "`"
			valueStr = valueStr + "\"" + utils.Interface2String(value) + "\""

		} else {
			insertStr = insertStr + ",`" + key + "`"
			valueStr = valueStr + ",\"" + utils.Interface2String(value) + "\""
		}
		i = i + 1
	}
	insertStr = insertStr + ") "
	valueStr = valueStr + ")"

	dbConnect := p.Config.Connect
	if p.Config.RunModeDebug {
		fmt.Printf("query statement: %s\n", (insertStr + valueStr))
	}

	stmt, err := dbConnect.Prepare((insertStr + valueStr))
	if err != nil {
		return fmt.Errorf("Prepare SQL(%s) error: %s.", err, (insertStr + valueStr))
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("exec SQL(%s) error: %s.", err, (insertStr + valueStr))
	}

	return nil
}

// execute a DB query according selectdata I
// return a set of the result and nil if teh SQL statement is be execute successful.
// Or return nil and error
func (p MySQL) NewQueryData(sd *SelectData) ([]map[string]interface{}, error) {
	var ret []map[string]interface{}

	if len(sd.Tb) < 1 || len(sd.OutFeilds) < 1 {
		return ret, fmt.Errorf("tables or output feilds is empty")
	}

	querySQL := "select "
	first := true
	for _, key := range sd.OutFeilds {
		if first {
			querySQL = querySQL + key
			first = false
		} else {
			querySQL = querySQL + "," + key
		}
	}

	querySQL = querySQL + " from "

	first = true
	for _, t := range sd.Tb {
		tArray := strings.Split(t, " ")
		tbStr := ""
		if len(tArray) > 1 {
			tbStr = "`" + tArray[0] + "` " + tArray[1]
		} else {
			tbStr = "`" + tArray[0] + "`"
		}

		if first {
			querySQL = querySQL + tbStr
			first = false
		} else {
			querySQL = querySQL + "," + tbStr
		}
	}

	first = true
	for key, value := range sd.Where {
		if first {
			querySQL = querySQL + " where " + key + value
			first = false
		} else {
			querySQL = querySQL + " and " + key + value
		}
	}

	first = true
	for _, key := range sd.Group {
		if first {
			querySQL = querySQL + " group by " + key
			first = false
		} else {
			querySQL = querySQL + "," + key
		}
	}

	first = true
	for _, key := range sd.Order {
		if first {
			querySQL = querySQL + " order by " + key.Key
			if key.Order == 0 {
				querySQL = querySQL + " ASC"
			} else {
				querySQL = querySQL + " DESC"
			}
			first = false
		} else {
			querySQL = querySQL + "," + key.Key
			if key.Order == 0 {
				querySQL = querySQL + " ASC"
			} else {
				querySQL = querySQL + " DESC"
			}
		}
	}

	if len(sd.Limit) == 1 {
		querySQL = querySQL + " limit " + strconv.Itoa(sd.Limit[0])
	}

	if len(sd.Limit) == 2 {
		querySQL = querySQL + " limit " + strconv.Itoa(sd.Limit[0]) + ", " + strconv.Itoa(sd.Limit[1])
	}

	dbConnect := p.Config.Connect
	if p.Config.RunModeDebug {
		fmt.Printf("Sql: %s\n", querySQL)
	}
	rows, err := dbConnect.Query(querySQL)
	if err != nil {
		return ret, fmt.Errorf("SQL query error: %s", err)
	}

	cols, err := rows.Columns()
	if err != nil {
		return ret, fmt.Errorf("get column data error: %s", err)
	}

	colsLen := len(cols)
	cache := make([]interface{}, colsLen)
	for i := range cache {
		var value interface{}
		cache[i] = &value
	}

	//var resData []FieldData
	for rows.Next() {
		_ = rows.Scan(cache...)

		line := make(map[string]interface{})
		for i, data := range cache {
			line[cols[i]] = *data.(*interface{})
			//line[cols[i]] = data
		}

		ret = append(ret, line)
	}

	_ = rows.Close()

	return ret, nil
}

// UpdateData: update data (map[string] interface{}) into the database according where(map[string]string)
// return nil if teh SQL statement is be execute successful.
// Or return 0 and []sysadmerror.Sysadmerror
func (p MySQL) NewUpdateData(tb string, data FieldData, where map[string]string) error {

	if tb == "" || len(data) < 1 {
		return fmt.Errorf("tables or feilds which will be update is empty")
	}

	querySQL := "update `" + tb + "` set "
	first := true
	for key, value := range data {
		if first {
			querySQL = querySQL + "`" + key + "`=" + utils.Interface2String(value) + ""
			first = false
		} else {
			querySQL = querySQL + "," + "`" + key + "`=" + utils.Interface2String(value) + ""
		}
	}

	first = true
	for key, value := range where {
		if first {
			querySQL = querySQL + " where `" + key + "`='" + value + "'"
			first = false
		} else {
			querySQL = querySQL + " and `" + key + "`='" + value + "'"
		}
	}

	dbConnect := p.Config.Connect
	if p.Config.RunModeDebug {
		fmt.Printf("query statement:%s \n", querySQL)
	}
	stmt, err := dbConnect.Prepare(querySQL)
	if err != nil {
		return fmt.Errorf("Prepare SQL %s  error: %s.", querySQL, err)
	}

	_, err = stmt.Exec()

	return err
}

// delete data from DB according selectData
// return nil the SQL statement is be execute successful.
// Or return error
func (p MySQL) NewDeleteData(dd *SelectData) error {
	if len(dd.Tb) < 1 {
		fmt.Errorf("tables is empty")
	}

	querySQL := "delete from "
	first := true
	for _, t := range dd.Tb {
		if first {
			querySQL = querySQL + "`" + t + "`"
			first = false
		} else {
			querySQL = querySQL + "," + "`" + t + "`"
		}
	}

	first = true
	for key, value := range dd.Where {
		if first {
			querySQL = querySQL + " where `" + key + "`='" + value + "'"
			first = false
		} else {
			querySQL = querySQL + " and `" + key + "`='" + value + "'"
		}
	}

	dbConnect := p.Config.Connect
	stmt, err := dbConnect.Prepare(querySQL)
	if p.Config.RunModeDebug {
		fmt.Printf("query statement:%s \n", querySQL)
	}
	if err != nil {
		return fmt.Errorf("Prepare SQL error: %s.", err)
	}
	_, err = stmt.Exec()

	return err
}

// NewBuildInsertQuery  build insert SQL statement according to tb and data.
// return string what can be execute query  and nil  if without error.return "" and nil
func (p MySQL) NewBuildInsertQuery(tb string, data FieldData) (string, error) {

	if len(tb) < 1 {
		return "", fmt.Errorf("Table name(%s) is not valid.", tb)
	}

	if len(data) < 1 {
		return "", fmt.Errorf("Can not insert empty data into table.")
	}

	insertStr := "INSERT INTO `" + tb + "`("
	valueStr := "Values ("
	i := 1
	for key, value := range data {
		if i == 1 {
			insertStr = insertStr + "`" + key + "`"
			valueStr = valueStr + "\"" + utils.Interface2String(value) + "\""

		} else {
			insertStr = insertStr + ",`" + key + "`"
			valueStr = valueStr + ",\"" + utils.Interface2String(value) + "\""
		}
		i = i + 1
	}
	insertStr = insertStr + ") "
	valueStr = valueStr + ")"

	return (insertStr + valueStr), nil
}

// NewBuildUpdateQuery build update SQL statement according to tb and data.
// return string what can be execute query and nil  if without error . return "" and error
func (p MySQL) NewBuildUpdateQuery(tb string, data FieldData, where map[string]string) (string, error) {
	if tb == "" || len(data) < 1 {
		return "", fmt.Errorf("tables or feilds which will be update is empty")
	}

	querySQL := "update `" + tb + "` set "
	first := true
	for key, value := range data {
		if first {
			querySQL = querySQL + "`" + key + "`=" + utils.Interface2String(value) + ""
			first = false
		} else {
			querySQL = querySQL + "," + "`" + key + "`=" + utils.Interface2String(value) + ""
		}
	}

	first = true
	for key, value := range where {
		if first {
			querySQL = querySQL + " where `" + key + "`='" + value + "'"
			first = false
		} else {
			querySQL = querySQL + " and `" + key + "`='" + value + "'"
		}
	}

	return querySQL, nil

}

// NewBuildDeleteQuery build update SQL statement according to dd .
// return string what can be execute query and nil if without error.otherewise return "" and error
func (p MySQL) NewBuildDeleteQuery(dd *SelectData) (string, error) {
	if len(dd.Tb) < 1 {
		return "", fmt.Errorf("tables is empty")
	}

	querySQL := "delete from "
	first := true
	for _, t := range dd.Tb {
		if first {
			querySQL = querySQL + "`" + t + "`"
			first = false
		} else {
			querySQL = querySQL + "," + "`" + t + "`"
		}
	}

	first = true
	for key, value := range dd.Where {
		if first {
			querySQL = querySQL + " where `" + key + "`='" + value + "'"
			first = false
		} else {
			querySQL = querySQL + " and `" + key + "`='" + value + "'"
		}
	}

	return querySQL, nil
}
