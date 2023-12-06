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

package objectsUI

const (
	// 当选中某个radio选项时，改变与该radio关联的object是否显示的属性.
	// 当前不支持对多个对象修改其显示属性
	JsActionKind_Radio_ChangeSubDisplay string = "1"

	// 表示当某个radio选项被选中时，调用addObjRadioCustizeAction 自定义JS函数。
	// 该函数通过页面自定义的addtionalJs文件定义，即相同的函数名，但是所实现的功能可能不同。
	// 自定义函数在被调用时，会传入actionUrl,objID, subObjID 和 radioOption .
	// 自定义函数是在addObjectJs.js文件的addObjRadioClick函数中被调用的.
	JsActionKind_Radio_CustomizeAction string = "2"

	// 当选择下拉菜单的菜单项时，修改关联的文本对象的值
	JsActionKind_Select_Change_TextValue string = "1"
	// 当选择下拉菜单的菜单项时，修改关联的select的下拉菜单项列表
	JsActionKind_Select_Change_SelectOptions string = "2"
	// 当选择下拉菜单的菜单项时，执行addObjSelectCustizeAction 自定义JS函数。
	// 该函数通过页面自定义的addtionalJs文件定义，即相同的函数名，但是所实现的功能可能不同。
	// 自定义函数在被调用时，会传入actionUrl,objID, subObjID 和 被选中的Option值 .
	JsActionKind_Select_CustomizeAction string = "4"

	// 默认的在前端显示的时间格式
	DefaultTimeStampFormat string = "2006-01-02 15:04:05"

	defaultResourceDetailTemplateFile = "resourceDetail.html"
)
