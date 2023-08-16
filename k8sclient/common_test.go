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

package k8sclient

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"k8s.io/client-go/rest"
	"strings"
	"sysadm/utils"
)

var restConf *rest.Config = nil
var dbUser = "sysadm"
var dbPasswd = "Sysadm12345"
var dbHost = "k8s.sysadm.cn"
var dbPort = 30306
var dbName = "k8ssysadm"
var clusterID = "26932263920893984"

func getClusterData() (string, string, string, string, string, string, error) {
	dbDsnstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPasswd, dbHost, dbPort, dbName)
	dbConnect, e := sql.Open("mysql", dbDsnstr)
	if e != nil {
		return "", "", "", "", "", "", e
	}

	e = dbConnect.Ping()
	if e != nil {
		return "", "", "", "", "", "", e
	}

	query := "select * from k8scluster where id=" + clusterID
	rows, e := dbConnect.Query(query)
	if e != nil {
		return "", "", "", "", "", "", e
	}

	cols, e := rows.Columns()
	if e != nil {
		return "", "", "", "", "", "", e
	}

	colsLen := len(cols)
	cache := make([]interface{}, colsLen)
	for i := range cache {
		var value interface{}
		cache[i] = &value
	}

	_ = rows.Next()
	_ = rows.Scan(cache...)
	line := make(map[string]interface{})
	for i, data := range cache {
		line[cols[i]] = *data.(*interface{})
	}

	id := utils.Interface2String(line["id"])
	apiserver := utils.Interface2String(line["apiserver"])
	clusterUser := utils.Interface2String(line["clusterUser"])
	ca := utils.Interface2String(line["ca"])
	cert := utils.Interface2String(line["cert"])
	key := utils.Interface2String(line["key"])

	apiserver = strings.ToLower(apiserver)
	if !strings.HasPrefix(apiserver, "https:") {
		apiserver = "https://" + apiserver
	}
	
	return id, apiserver, clusterUser, ca, cert, key, nil
}