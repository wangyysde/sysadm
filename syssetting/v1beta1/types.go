/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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

package v1beta1

import runtime "sysadm/apimachinery/runtime/v1beta1"

type Syssetting struct {
	// 配置项ID,数据库中自增整数值.字段名必须是ID，且tag值一定为id
	ID int `form:"id" json:"id" yaml:"id" xml:"id" db:"id"`
	// 配置项的应用范围，0 全局 1 k8s集群 2 节点级别 3 项目级别 4 用户组级别 5 用户级别
	// 对应的值参见syssettingDefined.go文件中的定义
	Scope int `form:"scope" json:"scope" yaml:"scope" xml:"scope" db:"scope"`
	// 配置了硒配置项的资源信息引用，字段名和切片类型是固定的。
	ReferenceObject []runtime.ReferenceInfo
	// 配置的key值, 不能为空。同一个级别的不能重复，大小写不敏感
	Key string `form:"key" json:"key" yaml:"key" xml:"key" db:"key"`
	// 配置项的默认值
	DefaultValue string `form:"defaultValue" json:"defaultValue" yaml:"defaultValue" xml:"defaultValue" db:"defaultValue"`
	// 配置项值
	Value string `form:"value" json:"value" yaml:"value" xml:"value" db:"value"`
	// 上次修改配置项值的用户,对应user表的userid.0表示系统自动设置的
	LastModifiedBy int `form:"lastModifiedBy" json:"lastModifiedBy" yaml:"lastModifiedBy" xml:"lastModifiedBy" db:"lastModifiedBy"`
	// 上次修改配置项值的时间截
	LastModifiedTime int `form:"lastModifiedTime" json:"lastModifiedTime" yaml:"lastModifiedTime" xml:"lastModifiedTime" db:"lastModifiedTime"`
	// 上次修改配置项值的原因
	LastModifiedReason string `form:"lastModifiedReason" json:"lastModifiedReason" yaml:"lastModifiedReason" xml:"lastModifiedReason" db:"lastModifiedReason"`
	// 上次修改前的值
	LastValue string `form:"lastValue" json:"lastValue" yaml:"lastValue" xml:"lastValue" db:"lastValue"`
}
