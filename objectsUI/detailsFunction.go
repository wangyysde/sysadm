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

Note :
本文件定义了动态生成表单数据所需要函数和方法
*/

package objectsUI

import (
	"fmt"
	"strings"
)

func AddTitleValueData(id, title, value, actionUri, actionFun string, lineData *LineData) error {
	id = strings.TrimSpace(id)
	title = strings.TrimSpace(title)
	value = strings.TrimSpace(value)
	actionUri = strings.TrimSpace(actionUri)
	actionFun = strings.TrimSpace(actionFun)

	if lineData == nil {
		return fmt.Errorf("line data must be not nil")
	}

	data := lineData.Data
	titleValueData := TitleValue{
		ID:        id,
		Title:     title,
		Kind:      "TitleValue",
		Value:     value,
		ActionUri: actionUri,
		ActionFun: actionFun,
	}
	data = append(data, titleValueData)
	lineData.Data = data

	return nil
}
