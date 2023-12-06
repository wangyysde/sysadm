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

func InitFormData(id, name, method, target, enctype, submitFn, cancelFn string) (*FormData, error) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	method = strings.TrimSpace(method)
	target = strings.TrimSpace(target)
	enctype = strings.TrimSpace(enctype)
	submitFn = strings.TrimSpace(submitFn)
	cancelFn = strings.TrimSpace(cancelFn)

	if name == "" {
		return nil, fmt.Errorf("the form name must be not empty")
	}

	if id == "" {
		id = name
	}

	if method == "" {
		method = "post"
	}

	if target == "" {
		target = "_self"
	}

	var newData []interface{}
	return &FormData{
		ID:       id,
		Name:     name,
		Method:   method,
		Target:   target,
		Enctype:  enctype,
		Data:     newData,
		SubmitFn: submitFn,
		CancelFn: cancelFn,
	}, nil
}

func InitLineData(id string, startContainer, endContainer, noDisplay bool) *LineData {
	var data []interface{}
	return &LineData{
		ID:             id,
		StartContainer: startContainer,
		EndContainer:   endContainer,
		NoDisplay:      noDisplay,
		Data:           data,
	}
}

func AddLabelData(id, size, withLine, title string, noDisplay bool, lineData *LineData) error {
	id = strings.TrimSpace(id)
	size = strings.ToLower(strings.TrimSpace(size))
	withLine = strings.TrimSpace(withLine)
	title = strings.TrimSpace(title)
	if lineData == nil {
		return fmt.Errorf("line data must be not nil")
	}

	if id == "" {
		return fmt.Errorf("label ID must be not empty")
	}

	if size != "big" && size != "mid" && size != "small" {
		return fmt.Errorf("Label size must be one of big, mid or small")
	}

	if withLine != "" && withLine != "Left" && withLine != "Top" && withLine != "Right" &&
		withLine != "Bottom" && withLine != "All" {
		return fmt.Errorf("with line must be one of empty, Left, Top, Right,Bottom or All")
	}
	if withLine != "" {
		withLine = "With" + withLine
	}

	data := lineData.Data
	labelData := LabelInput{
		ID:        id,
		Kind:      "LABEL",
		NoDisplay: noDisplay,
		Size:      size,
		WithLine:  withLine,
		Title:     title,
	}
	data = append(data, labelData)
	lineData.Data = data

	return nil
}

func AddTextData(id, name, defaultValue, title, actionUri, actionFun, note string, size int, disabled, noDisplay bool, lineData *LineData) error {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	defaultValue = strings.TrimSpace(defaultValue)
	title = strings.TrimSpace(title)
	actionUri = strings.TrimSpace(actionUri)
	actionFun = strings.TrimSpace(actionFun)
	note = strings.TrimSpace(note)
	data := lineData.Data

	if lineData == nil {
		return fmt.Errorf("line data must be not nil")
	}

	if name == "" {
		return fmt.Errorf("text input name must be not empty")
	}

	if id == "" {
		id = name
	}

	if size < 2 {
		return fmt.Errorf("text input size must be less 2")
	}

	textData := TextInput{
		ID:           id,
		Name:         name,
		Kind:         "TEXT",
		DefaultValue: defaultValue,
		Size:         size,
		Disable:      disabled,
		NoDisplay:    noDisplay,
		ActionUri:    actionUri,
		ActionFun:    actionFun,
		Note:         note,
		Title:        title,
	}

	data = append(data, textData)
	lineData.Data = data

	return nil
}

func AddSelectData(id, name, selectedData, actionUri, actionFun, title, note string, size int, disabled, noDisplay bool,
	options map[string]string, lineData *LineData) error {

	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	selectedData = strings.TrimSpace(selectedData)
	actionUri = strings.TrimSpace(actionUri)
	actionFun = strings.TrimSpace(actionFun)
	title = strings.TrimSpace(title)
	note = strings.TrimSpace(note)

	if lineData == nil {
		return fmt.Errorf("line data must be not nil")
	}

	if name == "" {
		return fmt.Errorf("text input name must be not empty")
	}

	if id == "" {
		id = name
	}

	if size < 2 {
		return fmt.Errorf("text input size must be less 2")
	}

	selectOptions := []Option{}
	i := 0
	for k, v := range options {
		if selectedData == "" && i == 0 {
			selectedData = k
		}
		o := Option{Text: v, Value: k}
		selectOptions = append(selectOptions, o)
	}

	data := lineData.Data
	selectData := SelectInput{
		ID:            id,
		Name:          name,
		Kind:          "SELECT",
		SelectedValue: selectedData,
		Size:          size,
		Disable:       disabled,
		NoDisplay:     noDisplay,
		ActionUri:     actionUri,
		ActionFun:     actionFun,
		Title:         title,
		Options:       selectOptions,
		Note:          note,
	}
	data = append(data, selectData)
	lineData.Data = data

	return nil
}

func AddCheckBoxOption(text, value string, checked, disabled bool, options []Option) ([]Option, error) {
	text = strings.TrimSpace(text)
	value = strings.TrimSpace(value)

	option := Option{
		Text:     text,
		Value:    value,
		Checked:  checked,
		Disabled: disabled,
	}
	options = append(options, option)

	return options, nil
}

func AddCheckBoxData(id, name, title, actionFun string, noDisplay bool, options []Option, lineData *LineData) error {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	actionFun = strings.TrimSpace(actionFun)
	title = strings.TrimSpace(title)

	if lineData == nil {
		return fmt.Errorf("line data must be not nil")
	}

	if name == "" {
		return fmt.Errorf("checkbox input name must be not empty")
	}

	if id == "" {
		id = name
	}

	if len(options) < 1 {
		return fmt.Errorf("options count of checkbox must be large zero")
	}

	data := lineData.Data
	checkBoxData := CheckBoxInput{
		ID:        id,
		Name:      name,
		Title:     title,
		Kind:      "CHECKBOX",
		NoDisplay: noDisplay,
		ActionFun: actionFun,
		Options:   options,
	}
	data = append(data, checkBoxData)

	return nil
}

func AddRadioData(id, name, title, actionFun string, noDisplay bool, options []Option, lineData *LineData) error {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	actionFun = strings.TrimSpace(actionFun)
	title = strings.TrimSpace(title)

	if lineData == nil {
		return fmt.Errorf("line data must be not nil")
	}

	if name == "" {
		return fmt.Errorf("radio input name must be not empty")
	}

	if id == "" {
		id = name
	}

	if len(options) < 1 {
		return fmt.Errorf("options count of checkbox must be large zero")
	}

	checkedNum := 0
	for _, v := range options {
		if v.Checked {
			checkedNum++
		}
	}
	if checkedNum > 1 {
		return fmt.Errorf("number of checked options should be not large one")
	}

	data := lineData.Data
	radioData := RadioInput{
		ID:        id,
		Name:      name,
		Title:     title,
		Kind:      "RADIO",
		NoDisplay: noDisplay,
		ActionFun: actionFun,
		Options:   options,
	}
	data = append(data, radioData)

	return nil
}

func AddFileData(id, name, defaultValue, title, note string, disabled, noDisplay bool, lineData *LineData) error {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	defaultValue = strings.TrimSpace(defaultValue)
	title = strings.TrimSpace(title)
	note = strings.TrimSpace(note)
	data := lineData.Data

	if lineData == nil {
		return fmt.Errorf("line data must be not nil")
	}

	if name == "" {
		return fmt.Errorf("file input name must be not empty")
	}

	if id == "" {
		id = name
	}

	fileData := FileInput{
		ID:           id,
		Name:         name,
		Kind:         "FILE",
		DefaultValue: defaultValue,
		Disable:      disabled,
		NoDisplay:    noDisplay,
		Note:         note,
		Title:        title,
	}

	data = append(data, fileData)
	lineData.Data = data

	return nil
}

func AddTextareaData(id, name, defaultValue, title, actionUri, actionFun, note string, colNum, rowNum int, readonly, noDisplay bool, lineData *LineData) error {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	defaultValue = strings.TrimSpace(defaultValue)
	title = strings.TrimSpace(title)
	actionUri = strings.TrimSpace(actionUri)
	actionFun = strings.TrimSpace(actionFun)
	note = strings.TrimSpace(note)
	data := lineData.Data

	if lineData == nil {
		return fmt.Errorf("line data must be not nil")
	}

	if name == "" {
		return fmt.Errorf("text input name must be not empty")
	}

	if id == "" {
		id = name
	}

	if colNum < 1 {
		colNum = 40
	}

	if rowNum < 1 {
		rowNum = 5
	}

	textareaData := TextareaInput{
		ID:           id,
		Name:         name,
		Kind:         "TEXTAREA",
		DefaultValue: defaultValue,
		ColNum:       colNum,
		RowNum:       rowNum,
		ReadOnly:     readonly,
		NoDisplay:    noDisplay,
		ActionUri:    actionUri,
		ActionFun:    actionFun,
		Note:         note,
		Title:        title,
	}

	data = append(data, textareaData)
	lineData.Data = data

	return nil
}

func AddHiddenData(id, name, defaultValue string, lineData *LineData) error {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	defaultValue = strings.TrimSpace(defaultValue)
	data := lineData.Data

	if lineData == nil {
		return fmt.Errorf("line data must be not nil")
	}

	if name == "" {
		return fmt.Errorf("hidden input name must be not empty")
	}

	if id == "" {
		id = name
	}

	hiddenData := HiddenInput{
		ID:           id,
		Name:         name,
		Kind:         "HIDDEN",
		DefaultValue: defaultValue,
	}

	data = append(data, hiddenData)
	lineData.Data = data

	return nil
}

func AddWordsInputData(id, name, word, url, actionFun string, noDisplay, awesome bool, lineData *LineData) error {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	word = strings.TrimSpace(word)
	url = strings.TrimSpace(url)
	actionFun = strings.TrimSpace(actionFun)

	if lineData == nil {
		return fmt.Errorf("line data must be not nil")
	}

	if name == "" {
		return fmt.Errorf("words input name must be not empty")
	}

	if id == "" {
		id = name
	}

	data := lineData.Data
	wordsData := WordsInput{
		ID:        id,
		Name:      name,
		Awesome:   awesome,
		Word:      word,
		Kind:      "WORDS",
		Url:       url,
		NoDisplay: noDisplay,
		ActionFun: actionFun,
	}
	data = append(data, wordsData)
	lineData.Data = data

	return nil
}
