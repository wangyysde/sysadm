package objectsUI

type ObjectTitle struct {
	// index of table header
	ID string

	// this title be show on table header
	Title string

	// uri for order
	OrderUri string

	// 指示对应的字段是否可以被排序
	IsOrder bool

	// direction for order. 1 for DESC 0 for ASC
	OrderDirection string

	// mark whether the field has be ordered
	OrderSelected string
}

// to save select options checkbox's items or radio items
type SubItems struct {
	// option value
	Value string

	// option text
	Text string

	// whether selected or checked
	Checked bool

	// 用于记录与之关联的对象是否显示。本字段主要应用于radio类型对象，配合addObjRadioClick JS方法用于控制
	// 与选项相关联的对象的显示属性。
	RelatedObjectIsDisplay bool
}

type ObjItemInfo struct {
	// message will be displayed at before input
	Title string

	// ID
	ID string

	// Name
	Name string

	// kind of input such as TEXT, SELECT,CHECKBOX,RADIO,FILE, TEXTAREA
	Kind string

	// size for TEXT,SELECT, FILE. this value is for textarea cols
	Size int

	// rows for textarea
	Rows int

	// whether disabled
	Disable bool

	// whether display in default
	NoDisplay bool

	// 当这个字段值不为空时，在前端会根据类型，会根据对象类型添加相应的JS事件调用
	JsActionKind string

	// for SELECT,RADIO,
	SubObjID string

	// define ajax uri
	ActionUri string

	// default value for itme
	DefaultValue string

	// ItemData for select checkbox and radio
	ItemData []SubItems

	// note
	Note string
}

type ObjLineData struct {
	Items []ObjItemInfo
}

type ItemForDetail struct {
	Label string

	Value string

	ActionUrl string

	ActionType string

	IsSeparator bool
}

type LineDataForDetail struct {
	Items []ItemForDetail
}

// 用于存储下拉菜单选择数据,最多支持三级菜单
type SelectOption struct {
	// 存储第一级菜单option的value.如果对应的菜单为非第一级菜单，则本字段值为0
	Id string `form:"id" json:"id" yaml:"id" xml:"id" db:"id"`
	// 存储第一组长菜单option的 text内容，如果对应的菜单为非第一级菜单，则本字段值为空
	Text string `form:"text" json:"text" yaml:"text" xml:"text" db:"text"`
	// 如果菜单项属于非第一级菜单的，则指明当前菜单项所属于的上一级菜单项的ID
	ParentID string `form:"parentID" json:"parentID" yaml:"parentID" xml:"parentID" db:"parentID"`
	// 如果菜单项属于非第一级菜单的，则指明当前父菜单项所对应的下一级菜单项列表，格式为[[value1,text1],[value2,text2]....]形式
	OptionsList string `form:"optionsList" json:"optionsList" yaml:"optionsList" xml:"optionsList" db:"optionsList"`
}

// 用于存储下拉菜单数据,最多支持三级菜单
type SelectData struct {
	// 本级菜单的Title,用于显示在下拉框的前面，可以为空
	Title string `form:"title" json:"title" yaml:"title" xml:"title" db:"title"`
	// 记录默认选中的菜单项ID
	SelectedId string `form:"selectedId" json:"selectedId" yaml:"selectedId" xml:"selectedId" db:"selectedId"`
	// 当上一级默认选中一个选项时，则本字段记录被选中选项对应的本级菜单选项
	SelectedOptions []SelectOption `form:"selectedOptions" json:"selectedOptions" yaml:"selectedOptions" xml:"selectedOptions" db:"selectedOptions"`
	// 菜单项列
	Options []SelectOption `form:"options" json:"options" yaml:"options" xml:"options" db:"options"`
}

// 定义一个用于定义sort的自定义函数类型
type SortBy func(p, q interface{}) bool

// 用于对结构体数据进行排序
type SortData struct {
	Data []interface{}
	By   SortBy
}
