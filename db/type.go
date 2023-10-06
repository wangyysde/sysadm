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
	"net"

	"sysadm/sysadmerror"
)

type DbConfig struct {
	Type         string   `json:"type"`
	Host         string   `json:"host"`
	HostIP       net.IP   `json:"hostIP"`
	Port         int      `json:"port"`
	User         string   `json:"user"`
	Password     string   `json:"password"`
	DbName       string   `json:"dbname"`
	SslMode      string   `json:"sslmode"`
	SslCa        string   `json:"sslca"`
	SslCert      string   `json:"sslcert"`
	SslKey       string   `json:"sslkey"`
	MaxOpenConns int      `json:"maxopenconns"`
	MaxIdleConns int      `json:"maxidleconns"`
	Connect      *sql.DB  `json:"connect"`
	RunModeDebug bool     `json:"runModeDebug"`
	Entity       DbEntity `json:"entity"`
}

type DbEntity interface {
	OpenDbConnect() []sysadmerror.Sysadmerror
	CloseDB() []sysadmerror.Sysadmerror
	InsertData(string, FieldData) (int, []sysadmerror.Sysadmerror)
	QueryData(sd *SelectData) ([]FieldData, []sysadmerror.Sysadmerror)
	DeleteData(dd *SelectData) (int64, []sysadmerror.Sysadmerror)
	UpdateData(string, FieldData, map[string]string) (int, []sysadmerror.Sysadmerror)
	BuildWhereFieldExact(string) string
	BuildWhereFieldExactWithSlice([]string) string
	BuildInsertQuery(tb string, data FieldData) (string, []sysadmerror.Sysadmerror)
	BuildUpdateQuery(tb string, data FieldData, where map[string]string) (string, []sysadmerror.Sysadmerror)
	BuildDeleteQuery(dd *SelectData) (string, []sysadmerror.Sysadmerror)
	GetDbConfig() *DbConfig
	NewInsertData(tb string, data FieldData) error
	NewQueryData(sd *SelectData) ([]map[string]interface{}, error)
	NewUpdateData(tb string, data FieldData, where map[string]string) error
	NewDeleteData(dd *SelectData) error
	NewBuildInsertQuery(tb string, data FieldData) (string, error)
	NewBuildUpdateQuery(tb string, data FieldData, where map[string]string) (string, error)
	NewBuildDeleteQuery(dd *SelectData) (string, error)
}

// key is the filed name and value is the value that will be set to the field.
type FieldData map[string]interface{}

type OrderData struct {
	Key   string
	Order int
}

type SelectData struct {
	Tb        []string
	OutFeilds []string
	Where     map[string]string
	Order     []OrderData
	Group     []string
	Limit     []int
}

type Tx struct {
	Entity DbEntity `json:"entity"`
	Tx     *sql.Tx
}
