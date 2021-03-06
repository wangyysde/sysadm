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
	"strconv"
	"strings"
	"time"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
)

type MySQL struct {
	Config *DbConfig `json:"config"`
}


/* 
  OpenDbConnect open a new connection to MySQL server with configuration parameters
  return errors with fatal level if there is any error occurred
  otherwise return errors with the levels lower fatal
  set the new connection to config.Connect
  and set tMaxOpenConns and MaxIdleConns for the connection
  p.CloseDB should be defer called after called this method
*/
func (p MySQL)OpenDbConnect() []sysadmerror.Sysadmerror {
	config := p.Config
	var errs []sysadmerror.Sysadmerror
	if config == nil {
		errs = append(errs,sysadmerror.NewErrorWithStringLevel(107000,"fatal","DB Configuration is nil"))
		return errs
	}
	
	errs = append(errs,sysadmerror.NewErrorWithStringLevel(107001,"debug","Try to open a connection to %s server",config.Type))
	dbDsnstr := ""
	if config.SslMode == "disable" {
		dbDsnstr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",config.User,config.Password, config.Host, config.Port, config.DbName)
	} else {
		//TODO :for connection with ssl
		//dbDsnstr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password,config.DbName)
	}
	
	dbConnect, err := sql.Open("mysql", dbDsnstr)
	if err != nil {
		errs = append(errs,sysadmerror.NewErrorWithStringLevel(107002,"fatal","Can not connect to postgre server with host: %s and Port: %d. error message is :%s",config.Host,config.Port,err))
		return errs
	}

	errs = append(errs,sysadmerror.NewErrorWithStringLevel(107003,"debug","Connect to postgre server with host: %s and Port: %d. successful",config.Host,config.Port))

	p.Config.Connect = dbConnect
	dbConnect.SetMaxOpenConns(p.Config.MaxOpenConns)
	dbConnect.SetMaxIdleConns(p.Config.MaxIdleConns)
	dbConnect.SetConnMaxLifetime(time.Minute * 5)   // TODO: this 5 should be configuratable

	err = dbConnect.Ping()
	if err != nil {
		errs = append(errs,sysadmerror.NewErrorWithStringLevel(107004,"fatal","we can connect to the postgre server while we can not ping the server with the connection.Error is:%s",err))
		return errs
	}

	return errs
}

/* 
  InsertData build insert SQL statement and execute a query using the SQL statement.
  return affected rows and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
  Or return 0 and []sysadmerror.Sysadmerror
*/
func (p MySQL)InsertData(tb string,data FieldData) (int, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if len(tb) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107006,"error","Table name(%s) is not valid.",tb))
		return 0, errs
	}
	
	if len(data) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107007,"error","Can not insert empty data into table."))
		return 0, errs
	}

    insertStr := "INSERT INTO `" + tb + "`("
    valueStr := "Values ("
      i := 1
    for key,value := range data {
        if i == 1 {
           insertStr = insertStr + "`" + key + "`"
		   valueStr = valueStr + "\""  + utils.Interface2String(value) + "\""

        } else {
           insertStr = insertStr + ",`" + key + "`"
           valueStr = valueStr + ",\""  +utils.Interface2String(value) + "\""
        }
        i = i + 1
    }
    insertStr = insertStr + ") "
    valueStr = valueStr + ")"

	dbConnect := p.Config.Connect
    stmt, err := dbConnect.Prepare((insertStr + valueStr))
    if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107011,"error","Prepare SQL(%s) error: %s.",err,(insertStr + valueStr)))
		return 0, errs
    }
    res, err := stmt.Exec()
   if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107013,"error","exec SQL(%s) error: %s.",err,(insertStr + valueStr)))
		return 0, errs
    }
    id, err := res.RowsAffected()
    if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107015,"error","fetch rows of affected error: %s.",err))
        return 0, errs
    }

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(107016,"debug","exec SQL(%s) successful.",(insertStr + valueStr)))
    return int(id), errs
}

/* 
  CloseDB try to close the connection to the DB server 
  return []sysadmerror.Sysadmerror
*/
func (p MySQL)CloseDB()([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	dbConnect := p.Config.Connect
	if dbConnect == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107017,"warning","The connection to DB server is nil. "))
	}

	dbConnect.Close()

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(107018,"debug","The connection to DB server has be closed. "))

	return errs
}

/*
   execute a DB query according selectdata I
   return a set of the result and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
   Or return nil and []sysadmerror.Sysadmerror
*/
func (p MySQL)QueryData(sd *SelectData) ([]FieldData, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(107019,"debug","now preparing db query."))
	if len(sd.Tb) <1 || len(sd.OutFeilds) <1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107020,"error","tables or output feilds is empty"))
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
			first = false
		} else {
			querySQL = querySQL + "," + key.Key
			 if key.Order == 0 {
				querySQL =querySQL + " ASC"
			}else {
				querySQL =querySQL + " DESC"
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
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(107021,"debug","now execute the SQL query: %s",querySQL))
	rows, err := dbConnect.Query(querySQL)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107022,"error","SQL query error: %s",err))
		return nil,errs
	}

	cols, err := rows.Columns()
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107023,"error","get column data error: %s",err))
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
   UpdateData: update data (map[string] interface{}) into the database according where(map[string]string)
   return affectRows and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
   Or return 0 and []sysadmerror.Sysadmerror
*/
func (p MySQL)UpdateData(tb string, data FieldData, where map[string]string) (int, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	if tb == ""  || len(data) <1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107032,"error","tables or feilds which will be update is empty"))
		return 0,errs
	}
	
	querySQL := "update `" + tb + "` set "
	first := true
	for key,value := range data {
		if first {
			querySQL = querySQL + "`" + key +"`=" + utils.Interface2String(value) + ""
			first = false
		} else {
			querySQL = querySQL + "," + "`" + key +"`=" + utils.Interface2String(value) + ""
		}
	}

	first = true
	for key,value := range where {
		if first {
			querySQL = querySQL + " where " + key + value
			first = false
		} else {
			querySQL = querySQL + " and " + key + value
		}
	}

	dbConnect := p.Config.Connect
	stmt, err := dbConnect.Prepare(querySQL)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1070333,"debug","try to execute SQL:%s",querySQL))
    if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107034,"error","Prepare SQL %s  error: %s.",querySQL,err))
		return 0, errs
    }
	
    res, err := stmt.Exec()
   if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107035,"error","exec SQL error: %s.",err))
		return 0, errs
    }
	
	ret,err := res.RowsAffected()
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107036,"error","can not get rowsaffected: %s.",err))
		return 0, errs
    }

	return int(ret), errs
}


/*
   execute a DB query according selectdata I
   return a set of the result and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
   Or return nil and []sysadmerror.Sysadmerror
*/
func (p MySQL)DeleteData(dd *SelectData) (int64, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(107024,"debug","now preparing db query."))
	if len(dd.Tb) <1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107025,"error","tables is empty"))
		return 0,errs
	}
	
	querySQL := "delete from "
	first := true
	for _,t := range dd.Tb {
		if first {
			querySQL = querySQL + "`" + t + "`"
			first = false
		} else {
			querySQL = querySQL + "," + "`" + t + "`"
		}
	}

	first = true
	for key,value := range dd.Where {
		if first {
			querySQL = querySQL + " where " + key + value
			first = false
		} else {
			querySQL = querySQL + " and " + key + value
		}
	}

	dbConnect := p.Config.Connect
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(107026,"debug","now execute the SQL query: %s",querySQL))
	stmt, err := dbConnect.Prepare(querySQL)
    if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107027,"error","Prepare SQL error: %s.",err))
		return 0, errs
    }
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(107028,"debug","Prepare SQL ok."))
    res, err := stmt.Exec()
   if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107029,"error","exec SQL error: %s.",err))
		return 0, errs
    }
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(107030,"debug","execute SQL query ok."))

	ret,err := res.RowsAffected()
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(107031,"error","can not get rowsaffected: %s.",err))
		return 0, errs
    }

	return ret, errs
}

/*
   BuildWhereFieldExact build the value of Where Field for key with value. 
*/
func (m MySQL)BuildWhereFieldExact(value string) string{
	if strings.TrimSpace(value) == ""{
		return ""
	}

	var ret = ""
	valueArray := strings.Split(value, ",")
	if len(valueArray) >1 {
		ret = " in ("
		first := true
		for _,v := range valueArray {
			if first {
				ret = ret + "\"" + v + "\""
				first = false
			} else {
				ret = ret + ",\"" + v + "\""
			}
		}
		ret += ")"
	} else {
		ret = ret + "=\"" + value + "\""
	}
	
	return ret
}

func (m MySQL)Identifier(identifier string) bool{
	matched,err := regexp.MatchString("^[a-zA-Z0-9]{1,64}",identifier)
	if !matched || err != nil {
		return false
	}

	return matched
}
