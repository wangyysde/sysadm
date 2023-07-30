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
	_ "github.com/lib/pq"
)

// InsertData build insert SQL statement and execute a query using the SQL statement.
// return affected rows and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
// Or return 0 and []sysadmerror.Sysadmerror
func (p Postgre) NewInsertData(tb string, data FieldData) error {
	if len(tb) < 1 {
		return fmt.Errorf("Table name(%s) is not valid.", tb)
	}

	// TODO
	return nil
}

/*
execute a DB query according selectdata I
return a set of the result and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
Or return nil and []sysadmerror.Sysadmerror
*/
func (p Postgre) NewQueryData(sd *SelectData) ([]map[string]interface{}, error) {
	var resData []map[string]interface{}
	// TODO
	return resData, nil
}

/*
execute a DB query according selectdata I
return a set of the result and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
Or return nil and []sysadmerror.Sysadmerror
TODO
*/
func (p Postgre) NewDeleteData(dd *SelectData) error {
	// TODO
	return nil
}

/*
execute a DB query according selectdata I
return a set of the result and []sysadmerror.Sysadmerror if teh SQL statement is be execute successful.
Or return nil and []sysadmerror.Sysadmerror
TODO
*/
func (p Postgre) NewUpdateData(tb string, data FieldData, where map[string]string) error {
	// TODO
	return nil
}
