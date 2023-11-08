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
	"github.com/wangyysde/sysadmServer"
	"net/http"
	"strconv"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/user"
)

func addformHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(7000190001, "debug", "try to display add object form page"))

	addTemplateFile := "addObjectFormNew.html"
	var emptyString []string
	baseUri := "/" + DefaultModuleName + "/"
	enctype := ""
	postUri := "/api/" + DefaultApiVersion + "/" + DefaultModuleName + "/add"
	var tplDataLines []objectsUI.ObjLineData

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000190002, errs, e)
		return
	}

	name := objectsUI.ObjItemInfo{Title: "系统名称", ID: "name", Name: "name", Kind: "TEXT", Size: 30, ActionUri: "api/" + DefaultApiVersion + "/" + DefaultModuleName + "/" + "validName", Note: "操作系统名称，英文且不区分大小写"}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{name}})

	versionName := objectsUI.ObjItemInfo{Title: "版本", ID: "versionName", Name: "versionName", Kind: "TEXT", Size: 30}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{versionName}})

	var archOptions []objectsUI.SubItems
	archOptionData := objectsUI.SubItems{Value: "0", Text: "===选择架构===", Checked: true}
	archOptions = append(archOptions, archOptionData)

	archOptionData = objectsUI.SubItems{Value: OsArchAarch, Text: OsArchAarch, Checked: false}
	archOptions = append(archOptions, archOptionData)

	archOptionData = objectsUI.SubItems{Value: OsArchPpc, Text: OsArchPpc, Checked: false}
	archOptions = append(archOptions, archOptionData)

	archOptionData = objectsUI.SubItems{Value: OsArchx86, Text: OsArchx86, Checked: false}
	archOptions = append(archOptions, archOptionData)

	archSelect := objectsUI.ObjItemInfo{Title: "架构", ID: "architecture", Name: "architecture", Kind: "SELECT",
		ItemData: archOptions}
	var lineItems []objectsUI.ObjItemInfo
	lineItems = append(lineItems, archSelect)

	var bitOptions []objectsUI.SubItems
	bitOptionData := objectsUI.SubItems{Value: "0", Text: "===选择位数===", Checked: true}
	bitOptions = append(bitOptions, bitOptionData)
	bitOptionData = objectsUI.SubItems{Value: strconv.Itoa(OsBit32), Text: strconv.Itoa(OsBit32), Checked: false}
	bitOptions = append(bitOptions, bitOptionData)
	bitOptionData = objectsUI.SubItems{Value: strconv.Itoa(OsBit64), Text: strconv.Itoa(OsBit64), Checked: false}
	bitOptions = append(bitOptions, bitOptionData)
	bitSelect := objectsUI.ObjItemInfo{Title: "位数", ID: "bit", Name: "bit", Kind: "SELECT",
		ItemData: bitOptions}
	lineItems = append(lineItems, bitSelect)
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: lineItems})

	descriptions := objectsUI.ObjItemInfo{Title: "描述", ID: "descriptions", Name: "descriptions", Kind: "TEXTAREA", Size: 40, Rows: 5}
	tplDataLines = append(tplDataLines, objectsUI.ObjLineData{Items: []objectsUI.ObjItemInfo{descriptions}})

	tplData, _ := objectsUI.InitAddObjectFormTemplateData(baseUri, "基础设施", "操作系统",
		enctype, postUri, "list", "list", emptyString, emptyString)
	tplData["data"] = tplDataLines

	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, addTemplateFile, tplData)
}
