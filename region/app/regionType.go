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
	sessions "github.com/wangyysde/sysadmSessions"
	sysadmDB "sysadm/db"
	sysadmObjects "sysadm/objects/app"
	"sysadm/sysadmLog"
	sysadmSetting "sysadm/syssetting/app"
)

type Region struct {
	// object name. the value should be "region"
	Name string
	// table name which hold host data in DB
	TableName string
	// field name of primary key in the table
	PkName string
}

// 国家表结构
type CountrySchema struct {
	// 中文国家名称
	ChineseName string `form:"chineseName" json:"chineseName" yaml:"chineseName" xml:"chineseName" db:"chineseName"`
	// 英文国家名称
	EnglishName string `form:"englishName" json:"englishName" yaml:"englishName" xml:"englishName" db:"englishName"`
	// 国家代码
	Code string `form:"code" json:"code" yaml:"code" xml:"code" db:"code"`
	// 是否显示在前台,0表示否，1表示是
	Display int `form:"display" json:"display" yaml:"display" xml:"display" db:"display"`
}

// 全国省市列表
type ProvinceSchema struct {
	// 省级代码
	Code string `form:"code" json:"code" yaml:"code" xml:"code" db:"code"`
	// 省级名称
	Name string `form:"name" json:"name" yaml:"name" xml:"name" db:"name"`
	// 下级市级代码列表，用,号分割
	CityList string `form:"cityList" json:"cityList" yaml:"cityList" xml:"cityList" db:"cityList"`
	// 所属于的国家代码
	CountryCode string `form:"countryCode" json:"countryCode" yaml:"countryCode" xml:"countryCode" db:"countryCode"`
}

// 全国地级市列表
type CitySchema struct {
	// 市级代码
	Code string `form:"code" json:"code" yaml:"code" xml:"code" db:"code"`
	// 市级名称
	Name string `form:"name" json:"name" yaml:"name" xml:"name" db:"name"`
	// 上级省市代码
	ProvinceCode string `form:"provinceCode" json:"provinceCode" yaml:"provinceCode" xml:"provinceCode" db:"provinceCode"`
	// 下级县区代码列表，用,号分割
	CountyList string `form:"countyList" json:"countyList" yaml:"countyList" xml:"countyList" db:"countyList"`
}

// 全国区县表
type CountySchema struct {
	// 区县代码
	Code string `form:"code" json:"code" yaml:"code" xml:"code" db:"code"`
	// 区县名称
	Name string `form:"name" json:"name" yaml:"name" xml:"name" db:"name"`
	// 上级地级市代码
	CityCode string `form:"cityCode" json:"cityCode" yaml:"cityCode" xml:"cityCode" db:"cityCode"`
}

// 存储运行期数据
type runingData struct {
	dbConf        *sysadmDB.DbConfig
	logEntity     *sysadmLog.LoggerConfig
	workingRoot   string
	sessionName   string
	sessionOption sessions.Options
	pageInfo      sysadmSetting.PageInfo
	objectEntiy   sysadmObjects.ObjectEntity
}
