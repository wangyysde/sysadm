package app

import (
	"fmt"
	"github.com/wangyysde/sysadmServer"
	"net/http"
	"strings"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	sysadmRegion "sysadm/region/app"
	"sysadm/sysadmLog"
	"sysadm/user"
	"sysadm/utils"
)

// order fields data of cluster list page
var allOrderFields = map[string]string{"TD1": "cnName", "TD2": "enName", "TD3": "country", "TD4": "province", "TD8": "status"}

// which field will be order if user has not selected
var defaultOrderField = "TD1"

// 1 for DESC 0 for ASC
var defaultOrderDirection = "1"

// define all list items(cols) name
var allListItems = map[string]string{"TD1": "名称", "TD2": "Name", "TD3": "国家", "TD4": "省", "TD5": "市", "TD6": "地址", "TD7": "值班电话", "TD8": "状态"}

// all popmenu items defined Format:
// item name, action name, action method
var allPopMenuItems = []string{"查看详情,detail,GET,page", "编辑数据中心,edit,GET,page", "删除数据中心,del,POST,tip", "启用数据中心,enable,POST,tip", "禁用数据中心,disable,POST,tip"}

func listHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(7000110001, "debug", "now handling datacenter list"))
	listTemplateFile := "objectlistNew.html"

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		outPutErrorMsg(c, 7000110002, errs, e)
		return
	}

	// get request data
	requestData, e := getRequestData(c)
	if e != nil {
		outPutErrorMsg(c, 7000110003, errs, e)
		return
	}

	// 初始化模板数据
	tplData, e := objectsUI.InitTemplateData("/"+DefaultObjectName+"/", "基础设施", "数据中心列表", "添加数据中心", "no",
		allPopMenuItems, []string{}, requestData)
	if e != nil {
		outPutErrorMsg(c, 7000110004, errs, e)
		return
	}

	// preparing select data
	region := sysadmRegion.New()
	conditions := make(map[string]string, 0)
	conditions["display"] = "= '" + sysadmRegion.CountryDisplay + "'"
	order := make(map[string]string, 0)
	var emptyString []string
	countryList, e := region.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		outPutErrorMsg(c, 7000110005, errs, e)
		return
	}
	delete(conditions, "display")
	provinceList, e := region.GetProvinceList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		outPutErrorMsg(c, 7000110006, errs, e)
		return
	}

	cityList, e := region.GetCityList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		outPutErrorMsg(c, 7000110007, errs, e)
		return
	}
	if e := buildSelectData(tplData, countryList, provinceList, cityList, requestData); e != nil {
		outPutErrorMsg(c, 7000110008, errs, e)
		return
	}

	// build table header for list objects
	objectsUI.BuildThData(requestData, allOrderFields, allListItems, tplData, defaultOrderField, defaultOrderDirection)

	searchContent := objectsUI.GetSearchContentFromRequest(requestData)
	ids := objectsUI.GetObjectIdsFromRequest(requestData)
	searchKeys := []string{"cnName", "address"}
	startPos := objectsUI.GetStartPosFromRequest(requestData)
	dcConditions := objectsUI.BuildCondition(requestData, "=0", "city")

	// get total number of list objects
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = New()
	dcCount, e := dcEntity.GetObjectCount(searchContent, ids, searchKeys, dcConditions)
	if e != nil || dcCount < 1 {
		outPutErrorMsg(c, 7000110009, errs, e)
		return
	}

	// get list data
	requestOrder := objectsUI.BuildOrderDataForQuery(requestData, allOrderFields, defaultOrderField, defaultOrderDirection)
	dcList, e := dcEntity.GetObjectList(searchContent, ids, searchKeys, dcConditions, startPos, runData.pageInfo.NumPerPage, requestOrder)
	if e != nil {
		outPutErrorMsg(c, 7000110010, errs, e)
		return
	}

	// prepare cluster list data
	objListData, e := prepareObjectData(countryList, provinceList, cityList, dcList)
	if e != nil {
		outPutErrorMsg(c, 7000110011, errs, e)
		return
	}
	tplData["objListData"] = objListData

	// prepare page number information
	objectsUI.BuildPageNumInfo(tplData, requestData, dcCount, startPos, runData.pageInfo.NumPerPage, defaultOrderField, defaultOrderDirection)

	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, listTemplateFile, tplData)
}

func getRequestData(c *sysadmServer.Context) (map[string]string, error) {
	requestData, e := utils.NewGetRequestData(c, []string{"groupSelectID", "searchContent", "start", "orderfield", "direction"})
	if e != nil {
		return requestData, e
	}

	objectIds := ""
	objectIDMap, _ := utils.GetRequestDataArray(c, []string{"objectid[]"})
	if objectIDMap != nil {
		objectIDSlice, ok := objectIDMap["objectid[]"]
		if ok {
			objectIds = strings.Join(objectIDSlice, ",")
		}
	}
	requestData["objectIds"] = objectIds
	if strings.TrimSpace(requestData["start"]) == "" {
		requestData["start"] = "0"
	}

	return requestData, nil
}

func outPutErrorMsg(c *sysadmServer.Context, errcode int, errs []sysadmLog.Sysadmerror, e error) {
	messageTemplateFile := "showmessage.html"
	messageTplData := make(map[string]interface{}, 0)
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(errcode, "error", "%s", e))
	runData.logEntity.LogErrors(errs)
	messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
	c.HTML(http.StatusOK, messageTemplateFile, messageTplData)

	return
}

func buildSelectData(tplData map[string]interface{}, countryList, provinceList, cityList []interface{}, requestData map[string]string) error {
	thirdOptions := make(map[string]string, 0)
	groupSelectID := requestData["groupSelectID"]
	thirdSelectedID := "0"
	secondSelectedID := "0"
	firstSelectedID := "0"
	var thirdSelectedOptions, secondSelectedOptions []objectsUI.SelectOption
	for _, line := range cityList {
		cityData, ok := line.(sysadmRegion.CitySchema)
		if !ok {
			return fmt.Errorf("the data is not province schema")
		}
		code := strings.TrimSpace(cityData.Code)
		provinceCode := strings.TrimSpace(cityData.ProvinceCode)
		name := strings.TrimSpace(cityData.Name)
		if groupSelectID == code {
			thirdSelectedID = code
			secondSelectedID = provinceCode
		}
		if secondSelectedID == provinceCode {
			selectedOption := objectsUI.SelectOption{
				Id:       code,
				Text:     name,
				ParentID: provinceCode,
			}
			thirdSelectedOptions = append(thirdSelectedOptions, selectedOption)
		}

		subOption := "['" + code + "','" + name + "']"
		addOption, ok := thirdOptions[provinceCode]
		if ok {
			addOption = addOption + "," + subOption
		} else {
			addOption = subOption
		}
		thirdOptions[provinceCode] = addOption
	}
	thirdSelect := objectsUI.SelectData{Title: "市(区)", SelectedId: thirdSelectedID, SelectedOptions: thirdSelectedOptions}
	var thirdOptionList []objectsUI.SelectOption
	for code, value := range thirdOptions {
		option := objectsUI.SelectOption{
			ParentID:    code,
			OptionsList: value,
		}
		thirdOptionList = append(thirdOptionList, option)
	}
	thirdSelect.Options = thirdOptionList

	secondSelect := objectsUI.SelectData{Title: "省(市)"}
	secondOptions := make(map[string]string, 0)
	for _, line := range provinceList {
		provinceData, ok := line.(sysadmRegion.ProvinceSchema)
		if !ok {
			return fmt.Errorf("the data is not province schema")
		}
		code := strings.TrimSpace(provinceData.Code)
		countryCode := strings.TrimSpace(provinceData.CountryCode)
		name := strings.TrimSpace(provinceData.Name)
		subOption := "['" + code + "','" + name + "']"
		if code == secondSelectedID {
			firstSelectedID = countryCode

		}
		if firstSelectedID == countryCode {
			selectedOption := objectsUI.SelectOption{
				Id:       code,
				Text:     name,
				ParentID: countryCode,
			}
			secondSelectedOptions = append(secondSelectedOptions, selectedOption)
		}

		addOption, ok := secondOptions[countryCode]
		if ok {
			addOption = addOption + "," + subOption
		} else {
			addOption = subOption
		}
		secondOptions[countryCode] = addOption
	}
	var secondOptionList []objectsUI.SelectOption
	for code, value := range secondOptions {
		option := objectsUI.SelectOption{
			ParentID:    code,
			OptionsList: value,
		}
		secondOptionList = append(secondOptionList, option)
	}
	secondSelect.Options = secondOptionList

	var firstOptions []objectsUI.SelectOption
	firstSelect := objectsUI.SelectData{
		Title:      "国家",
		SelectedId: firstSelectedID,
	}
	for _, line := range countryList {
		countryData, ok := line.(sysadmRegion.CountrySchema)
		if !ok {
			return fmt.Errorf("the data is not country schema")
		}
		option := objectsUI.SelectOption{
			Id:       countryData.Code,
			Text:     countryData.ChineseName,
			ParentID: "0",
		}
		firstOptions = append(firstOptions, option)
	}
	firstSelect.Options = firstOptions

	firstSelect.SelectedId = firstSelectedID
	secondSelect.SelectedId = secondSelectedID
	secondSelect.SelectedOptions = secondSelectedOptions
	thirdSelect.SelectedId = thirdSelectedID
	thirdSelect.SelectedOptions = thirdSelectedOptions

	if e := objectsUI.BuildMultiSelectData(firstSelect, secondSelect, thirdSelect, tplData); e != nil {
		return e
	}

	return nil
}

func prepareObjectData(countryList, provinceList, cityList, dcList []interface{}) ([]map[string]interface{}, error) {
	var dataList []map[string]interface{}
	for _, line := range dcList {
		dcData, ok := line.(DatacenterSchema)
		if !ok {
			return dataList, fmt.Errorf("data is not datacenter schema")
		}

		lineMap := make(map[string]interface{}, 0)
		lineMap["TD1"] = dcData.CnName
		lineMap["TD2"] = dcData.EnName
		lineMap["TD3"] = getCountryName(dcData.Country, countryList)
		lineMap["TD4"] = getProvinceName(dcData.Province, provinceList)
		lineMap["TD5"] = getCityName(dcData.City, cityList)
		lineMap["TD6"] = dcData.Address
		lineMap["TD7"] = dcData.DutyTel
		statusStr := ""
		popmenuitems := ""
		switch dcData.Status {
		case StatusUnused:
			statusStr = "未启用"
			popmenuitems = "0,1,2,3"
		case StatusEnabled:
			statusStr = "启用"
			popmenuitems = "0,1,4"
		case StatusDisabled:
			statusStr = "禁用"
			popmenuitems = "0,1,2,3"
		}
		lineMap["TD8"] = statusStr
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return dataList, nil
}

func getCountryName(code string, countryList []interface{}) string {
	for _, line := range countryList {
		countryData, ok := line.(sysadmRegion.CountrySchema)
		if !ok {
			return "未知"
		}
		if code == countryData.Code {
			return countryData.ChineseName
		}
	}

	return "未知"
}

func getProvinceName(code string, provinceList []interface{}) string {
	for _, line := range provinceList {
		provinceData, ok := line.(sysadmRegion.ProvinceSchema)
		if !ok {
			return "未知"
		}
		if code == provinceData.Code {
			return provinceData.Name
		}
	}

	return "未知"
}

func getCityName(code string, cityList []interface{}) string {
	for _, line := range cityList {
		cityData, ok := line.(sysadmRegion.CitySchema)
		if !ok {
			return "未知"
		}
		if code == cityData.Code {
			return cityData.Name
		}
	}

	return "未知"
}
