package objectsUI

type ObjectTitle struct {
	// index of table header
	ID string

	// this title be show on table header
	Title string

	// uri for order
	OrderUri string

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

	// for SELECT,
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