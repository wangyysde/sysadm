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
	"database/sql"
	"fmt"
	"regexp"
	"strconv"

	_ "github.com/lib/pq"

	"sysadm/sysadmerror"
)

type Postgre struct {
	Config *DbConfig `json:"config"`
}


// OpenDbConnect open a new connection to postgre server with configuration parameters
// return errors with fatal level if there is any error occurred
// otherwise return errors with the levels lower fatal
// set the new connection to config.Connect
// and set tMaxOpenConns and MaxIdleConns for the connection
// p.CloseDB should be defer called after called this method
func (p Postgre)OpenDbConnect() []sysadmerror.Sysadmerror {
	config := p.Config
	var errs []sysadmerror.Sysadmerror
	if config == nil {
		errs = append(errs,sysadmerror.NewErrorWithStringLevel(101000,"fatal","DB Configuration is nil"))
		return errs
	}
	
	errs = append(errs,sysadmerror.NewErrorWithStringLevel(101001,"debug","Try to open a connection to %s server",config.Type))
	dbDsnstr := ""
	if config.SslMode == "disable" {
		dbDsnstr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password,config.DbName)
	} else {
		//TODO :for connection with ssl
		dbDsnstr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password,config.DbName)
	}
	
	dbConnect, err := sql.Open("postgres", dbDsnstr)
	if err != nil {
		errs = append(errs,sysadmerror.NewErrorWithStringLevel(101002,"fatal","Can not connect to postgre server with host: %s and Port: %d. error message is :%s",config.Host,config.Port,err))
		return errs
	}

	errs = append(errs,sysadmerror.NewErrorWithStringLevel(101003,"debug","Connect to postgre server with host: %s and Port: %d. successful",config.Host,config.Port))

	p.Config.Connect = dbConnect
	dbConnect.SetMaxOpenConns(p.Config.MaxOpenConns)
	dbConnect.SetMaxIdleConns(p.Config.MaxIdleConns)

	err = dbConnect.Ping()
	if err != nil {
		errs = append(errs,sysadmerror.NewErrorWithStringLevel(101004,"fatal","we can connect to the postgre server while we can not ping the server with the connection.Error is:%s",err))
		return errs
	}

	return errs
}

// InsertData build insert SQL statement and execute a query using the SQL statement.
// return affected rows and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
// Or return 0 and []sysadmerror.Sysadmerror
func (p Postgre)InsertData(tb string,data FieldData) (int, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(101005,"debug","Checking insert data is valid."))
	if len(tb) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(101006,"error","Table name(%s) is not valid.",tb))
		return 0, errs
	}
	
	if len(data) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(101007,"error","Can not insert empty data into table."))
		return 0, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(101008,"debug","Preparing insert SQL for inert into table (%s).",tb))
    insertStr := "INSERT INTO \"" + tb + "\"("
    placeHoldStr := "Values ("
    var values []interface{}
    i := 1
    for key,value := range data {
        if i == 1 {
           insertStr = insertStr + "\"" + key + "\""
           placeHoldStr = fmt.Sprintf("%s$%d",placeHoldStr,i)
        } else {
           insertStr = insertStr + ",\"" + key + "\""
           placeHoldStr = fmt.Sprintf("%s,$%d",placeHoldStr,i)
        }
        values = append(values,value)
        i = i + 1
    }
    insertStr = insertStr + ") "
    placeHoldStr = placeHoldStr + ")"

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(101009,"debug","Insert SQL: %s.",(insertStr + placeHoldStr)))
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(101010,"debug","Insert Data: %v.",values))

	dbConnect := p.Config.Connect
    stmt, err := dbConnect.Prepare((insertStr + placeHoldStr))
    if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(101011,"error","Prepare SQL error: %s.",err))
		return 0, errs
    }
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(101012,"debug","Prepare SQL ok."))
    res, err := stmt.Exec(values...)
   if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(101013,"error","exec SQL error: %s.",err))
		return 0, errs
    }
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(101014,"debug","execute SQL query ok."))

    id, err := res.RowsAffected()
    if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(101015,"error","fetch rows of affected error: %s.",err))
        return 0, errs
    }
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(101016,"debug","fetch rows of affected ok."))

    return int(id), errs
}

// CloseDB try to close the connection to the DB server 
// return []sysadmerror.Sysadmerror
func (p Postgre)CloseDB()([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	dbConnect := p.Config.Connect
	if dbConnect == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(101017,"warning","The connection to DB server is nil. "))
	}

	dbConnect.Close()

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(101018,"debug","The connection to DB server has be closed. "))

	return errs
}

/*
   execute a DB query according selectdata I
   return a set of the result and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
   Or return nil and []sysadmerror.Sysadmerror
*/
func (p Postgre)QueryData(sd *SelectData) ([]FieldData, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(101019,"debug","now preparing db query."))
	if len(sd.Tb) <1 || len(sd.OutFeilds) <1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(101020,"error","tables or output feilds is empty"))
		return nil,errs
	}
	
	querySQL := "select "
	first := true
	for _,key := range sd.OutFeilds {
		if first {
			querySQL = querySQL + key
			first = false
		} else {
			querySQL = querySQL + "," + key
		}
	}

	querySQL = querySQL + " from " 

	first = true
	for _,t := range sd.Tb {
		if first {
			querySQL = querySQL + "\"" + t + "\""
			first = false
		} else {
			querySQL = querySQL + "," + "\"" + t + "\""
		}
	}

	first = true
	for key,value := range sd.Where {
		if first {
			querySQL = querySQL + " where " + key + value
			first = false
		} else {
			querySQL = querySQL + " and " + key + value
		}
	}

	first = true
	for _,key := range sd.Group {
		if first {
			querySQL = querySQL + " group by " + key 
			first = false
		} else {
			querySQL = querySQL + "," + key 
		}
	}

	first = true
	for _,key := range sd.Order {
		if first {
			querySQL = querySQL + " order by " + key.Key
			if key.Order == 0 {
				querySQL =querySQL + " ASC"
			}else {
				querySQL =querySQL + " DESC"
			} 
			// Because postgre can not execute the statement with multiple order fields, so break here
			break
		} 
	}

	if len(sd.Limit) == 1 {
		querySQL = querySQL + " limit " + strconv.Itoa(sd.Limit[0]) 
	}

	if len(sd.Limit) == 2 {
		querySQL = querySQL + " limit " + strconv.Itoa(sd.Limit[0]) + " OFFSET " + strconv.Itoa(sd.Limit[1]) 
	}

	dbConnect := p.Config.Connect
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(101021,"debug","now execute the SQL query: %s",querySQL))
	rows, err := dbConnect.Query(querySQL)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(101022,"error","SQL query error: %s",err))
		return nil,errs
	}

	cols, err := rows.Columns()
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(101023,"error","get column data error: %s",err))
	}

	colsLen := len(cols)
	cache := make([]interface{},colsLen)
	for i := range cache {
		var value interface{}
		cache[i] = &value
	}
    
	var resData []FieldData
	for rows.Next(){
		_ = rows.Scan(cache...)

		line :=  make(map[string]interface{})
		for i, data := range cache {
			line[cols[i]] = *data.(*interface{})
		}

		resData = append(resData,line)
	}

	_ = rows.Close()
	
	return resData, errs
}

/*
   execute a DB query according selectdata I
   return a set of the result and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
   Or return nil and []sysadmerror.Sysadmerror
   TODO
*/
func (p Postgre)DeleteData(dd *SelectData) (int64, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	return 0, errs
}
/*
   execute a DB query according selectdata I
   return a set of the result and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
   Or return nil and []sysadmerror.Sysadmerror
   TODO
*/
func (p Postgre)UpdateData(tb string, data FieldData, where map[string]string) (int, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	return 0, errs
}


/*
   BuildWhereFieldExact build the value of Where Field for key with value. 
   
   TODO
*/
func (p Postgre)BuildWhereFieldExact(value string) string{
	return ""
}

/*
   BuildWhereFieldExactWithSlice build the value of Where Field for key with value. 
   
   TODO
*/
func (p Postgre)BuildWhereFieldExactWithSlice(value []string) string{
	return ""
}



/*
   BuildWhereFieldExact build the value of Where Field for key with value. 
   
   TODO
*/
func (p Postgre)Identifier(identifier string) bool{
	matched,err := regexp.MatchString("^[a-zA-Z0-9]{1,64}",identifier)
	if !matched || err != nil {
		return false
	}

	return matched
}

/* 
  BuildInsertQuery  build insert SQL statement according to tb and data.
  return string what can be execute query  and []sysadmerror.Sysadmerror if without error .
  Or return "" and []sysadmerror.Sysadmerror

  TODO
*/
func (p Postgre)BuildInsertQuery(tb string,data FieldData) (string, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	return "", errs
}

/*
  BuildUpdateQuery build update SQL statement according to tb and data.
  return string what can be execute query and []sysadmerror.Sysadmerror if without error .
  Or return "" and []sysadmerror.Sysadmerror

  TODO
*/
func (p Postgre)BuildUpdateQuery(tb string, data FieldData, where map[string]string) (string, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	return "", errs
}

/*
  BuildDeleteQuery build update SQL statement according to dd .
  return string what can be execute query and []sysadmerror.Sysadmerror if without error .
  Or return "" and []sysadmerror.Sysadmerror
*/
func (p Postgre)BuildDeleteQuery(dd *SelectData) (string, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	return "", errs
}

/*
	Get the DBConfig
*/
func (p Postgre)GetDbConfig()(*DbConfig){
	return p.Config
}