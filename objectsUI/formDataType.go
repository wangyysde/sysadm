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
本文件定义了动态生成表单数据所需要数据结构，本文件中定义的数据结构是对globalType.go文件中定义的数据结构的优化版本.
*/

package objectsUI

type FormData struct {
	// 指定的值会设置在form id 部分.如果没有设置,则以name值作为ID值
	ID string

	// 指定的值会设置在form name 部分，此部分的值不能为空
	Name string

	// 提交方法，默认是post
	Method string

	// 提交后的target, 默认是"_self"
	Target string

	// 是否有多部分数据提交，通常为空或者"multipart/form-data"
	Enctype string

	// form 行数据
	Data []interface{}

	// 如果本字段的值不为零时，则forumDataSubmitForm函数会调用本字段值命名的JS函数提交表单。
	// 否则forumDataSubmitForm会将表单提交到指定地址,具体地址见addObjTabsSubmitForm函数的实现方法
	SubmitFn string

	// 如果本字段的值不为零时，当用户点击表单页面上的取消按键时，由forumDataCancelForm函数调用本字段值命名的JS函数取消表单提交，退回到指定页面。
	// 否则forumDataCancelForm函数会重新load /pageUri/module/list页面
	CancelFn string
}

type LineData struct {
	// 行ID，也就div的id
	ID string

	// 如果为真，则表示当前行作为一个容器的开始行，即在行div标签前再加一个div标签
	StartContainer bool

	// 如果为真，则表示当前行为一个容器的结束行，即在该行的后面会多加一个div的结束标签
	EndContainer bool

	// 如果为true，表示默认情况下不显示
	NoDisplay bool

	// 本行对应的对象数据列表
	Data []interface{}
}

type LabelInput struct {
	// Label的ID，不能为空.其作为span ID的一部分
	ID string

	// 对象的类型,此字段的值固定为LABEL
	Kind string

	// 如果为true，表示默认情况下不显示
	NoDisplay bool

	// Label显示的大小，其值只能是big,mid和small之一
	Size string

	// 是否显示边线, 其值只能是空、Left, Top, Right,Bottom 和 All
	WithLine string

	// Label所需要显示的文本内容
	Title string
}

type TextInput struct {
	// 指定的值会设置在<input id 部分.如果没有设置,则以name值作为ID值
	ID string

	// 指定的值会设置在<input name 部分，此部分的值不能为空
	Name string

	// 对象的类型,此字段的值固定为TEXT
	Kind string

	// 默认值
	DefaultValue string

	// 文本框的长度
	Size int

	// 是否禁用
	Disable bool

	// 如果为true，表示默认情况下不显示
	NoDisplay bool

	// 当本字段的值不为空时，设置本文本框在失去焦点调用formDataTextInputValueChange函数。在调用此函数时
	// 本字段的值将与pageUrl变量的值拼接在一起组成Ajax调用的uri地址。且此请会以当前text内的值作为objvalue参数的值附在URI上
	ActionUri string

	// 当本字段的值不为空时，设置本文本框在失去焦点调用formDataTextInputValueChange函数。如本字段的值不为空
	// 则formDataTextInputValueChange函数先优先调用以本字段值为名字的函数(如果存在),否则调用Ajax请求ActionUri地址
	ActionFun string

	// Text Input 的 Title
	Title string

	// 用于显示在输入文本框后面的注意提示信息
	Note string
}

type Option struct {
	// Select Radio 或 checkbox 的选项文本
	Text string

	// Select Radio 或 checkbox 的选项值
	Value string

	// 对于radio和checkbox,指示option是否被checked
	Checked bool

	// 对于radio和checkbox,指示option是否被disabled
	Disabled bool
}

type SelectInput struct {
	// 指定的值会设置在<input id 部分.如果没有设置,则以name值作为ID值
	ID string

	// 指定的值会设置在<input name 部分，此部分的值不能为空
	Name string

	// 对象的类型,此字段的值固定为SELECT
	Kind string

	// 选中的值
	SelectedValue string

	// 文本框的长度
	Size int

	// 是否禁用
	Disable bool

	// 如果为true，表示默认情况下不显示
	NoDisplay bool

	// 当本字段的值不为空时，设置当Select 选中的选项变化时调用formDataSelectInputValueChange函数。在调用此函数时
	// 本字段的值将与pageUrl变量的值拼接在一起组成Ajax调用的uri地址。且此请求会以当前选中选项的值作为objvalue参数的值附在URI上
	ActionUri string

	// 当本字段的值不为空时，设置当Select 选中的选项变化时调用formDataSelectInputValueChange函数。如本字段的值不为空
	// 则formDataSelectInputValueChange函数先优先调用以本字段值为名字的函数(如果存在),否则调用Ajax请求ActionUri地址
	ActionFun string

	// Text Input 的 Title
	Title string

	// 存放select的options数据
	Options []Option

	// 用于显示在输入文本框后面的注意提示信息
	Note string
}

type CheckBoxInput struct {
	// 指定的值会设置在<input id 部分.如果没有设置,则以name值作为ID值
	ID string

	// 指定的值会设置在<input name 部分，此部分的值不能为空
	Name string

	// checkbox 组的title
	Title string

	// 对象的类型,此字段的值固定为SELECT
	Kind string

	// 如果为true，表示默认情况下不显示
	NoDisplay bool

	// 当本字段的值不为空时，设置当checkbox 某项被点击时调用formDataCheckBoxClick函数。
	// 则formDataCheckBoxClick函数检测本字段所指定的函数是否存在，如果存在则调用以本字段值命名的函数
	ActionFun string

	// 存放radio的options数据
	Options []Option
}

type RadioInput struct {
	// 指定的值会设置在<input id 部分.如果没有设置,则以name值作为ID值
	ID string

	// 指定的值会设置在<input name 部分，此部分的值不能为空
	Name string

	// Radio 组的title
	Title string

	// 对象的类型,此字段的值固定为RADIO
	Kind string

	// 如果为true，表示默认情况下不显示
	NoDisplay bool

	// 当本字段的值不为空时，设置当radio 某项被点击时调用formDataRadioClick函数。
	// 则formDataRadioClick函数检测本字段所指定的函数是否存在，如果存在则调用以本字段值命名的函数
	ActionFun string

	// 存放radio的options数据
	Options []Option
}

type FileInput struct {
	// Note: 为了美化File类型的Input按钮，这个类型有固定的Action行为，故File类型的Input不支持定义ActionUri 和 ActionFun
	// 指定的值会设置在<input id 部分.如果没有设置,则以name值作为ID值
	ID string

	// 指定的值会设置在<input name 部分，此部分的值不能为空
	Name string

	// 对象的类型,此字段的值固定为FILE
	Kind string

	// 默认值
	DefaultValue string

	// 是否禁用
	Disable bool

	// 如果为true，表示默认情况下不显示
	NoDisplay bool

	// Text Input 的 Title
	Title string

	// 用于显示在输入文本框后面的注意提示信息
	Note string
}

type TextareaInput struct {
	// Note: 为了美化File类型的Input按钮，这个类型有固定的Action行为，故File类型的Input不支持定义ActionUri 和 ActionFun
	// 指定的值会设置在<input id 部分.如果没有设置,则以name值作为ID值
	ID string

	// 指定的值会设置在<input name 部分，此部分的值不能为空
	Name string

	// 对象的类型,此字段的值固定为FILE
	Kind string

	// 默认值
	DefaultValue string

	// 列宽数
	ColNum int

	// 行数
	RowNum int

	// 是否只读
	ReadOnly bool

	// 如果为true，表示默认情况下不显示
	NoDisplay bool

	// Text Input 的 Title
	Title string

	// 当本字段的值不为空时，设置本文本框在失去焦点调用formDataTextareaInputValueChange函数。在调用此函数时
	// 本字段的值将与pageUrl变量的值拼接在一起组成Ajax调用的uri地址。且此请会以当前text内的值作为objvalue参数的值附在URI上
	ActionUri string

	// 当本字段的值不为空时，设置本文本框在失去焦点调用formDataTextareaInputValueChange函数。如本字段的值不为空
	// 则formDataTextareaInputValueChange函数先优先调用以本字段值为名字的函数(如果存在),否则调用Ajax请求ActionUri地址
	ActionFun string

	// 用于显示在输入文本框后面的注意提示信息
	Note string
}

type HiddenInput struct {
	// 指定的值会设置在<input id 部分.如果没有设置,则以name值作为ID值
	ID string

	// 指定的值会设置在<input name 部分，此部分的值不能为空
	Name string

	// 对象的类型,此字段的值固定为TEXT
	Kind string

	// 默认值
	DefaultValue string
}

type WordsInput struct {
	// 指定的值会设置在<input id 部分.如果没有设置,则以name值作为ID值
	ID string

	// 指定的值会设置在span name 部分，此部分的值不能为空
	Name string

	// 对象的类型,此字段的值固定为WORDS
	Kind string

	// 文字是否是awesome字体
	Awesome bool

	// 文本内容,显示在页面上的内容
	Word string

	// 如果本字段的值不为空，则表示是一个带超级连接的文本
	Url string

	// 如果为true，表示默认情况下不显示
	NoDisplay bool

	// 当本字段的值不为空,且是有超级连接时，设置超级连接onClick事件调用formDataWordsInputClick函数。
	// formDataWordsInputClick函数调用以本字段值为名字的函数(如果存在)
	ActionFun string
}
