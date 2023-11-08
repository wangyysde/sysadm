/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
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
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sysadm/user"

	"github.com/wangyysde/sysadmServer"
	"sysadm/db"
	"sysadm/httpclient"
	"sysadm/sysadmapi/apiutils"
	"sysadm/sysadmerror"
	"sysadm/utils"
)

func addHost(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror

	requestData := ApiHost{}

	if e, err := checkAddHostData(&requestData, c); e != nil {
		errs = append(errs, err...)
		msg := fmt.Sprintf("%s", e)
		err = apiutils.SendResponseForErrorMessage(c, 3030002, msg)
		errs = append(errs, err...)
		logErrors(errs)
		return

	}

	// adding host information into DB
	dbConf := WorkingData.dbConf
	dbEntity := dbConf.Entity
	tx, err := db.Begin(dbEntity)
	errs = append(errs, err...)
	if tx == nil {
		err := apiutils.SendResponseForErrorMessage(c, 3030003, "start transaction on DB query error. check logs for details.")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	// insert host information into host table
	hostid, err := addHostToDB(tx, &requestData)
	errs = append(errs, err...)
	if hostid == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c, 30301001, "add host error. check logs for details")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	// insert IP information into hostIP table
	ipID, err := addManageIPToDB(tx, &requestData, hostid)
	errs = append(errs, err...)
	if ipID == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c, 30301002, fmt.Sprintf("host %s manage IP %s add into DB error", requestData.Hostname, requestData.Ip))
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	nextCommandID, err := getNextID("command", "commandID", tx)
	errs = append(errs, err...)
	if nextCommandID == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c, 30303017, fmt.Sprintf("add the information of  host %s into DB error", requestData.Hostname))
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	// insert command information into command
	cID, err := addCommandToDB(tx, "gethostip", 0, hostid, nextCommandID)
	errs = append(errs, err...)
	if cID == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c, 30301004, fmt.Sprintf("add the information of host %s into DB error", requestData.Hostname))
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	parameters := make(map[string]string, 0)
	parameters["withmac"] = "1"
	parameters["withmask"] = "1"
	pid, err := addCommandParametersToDB(tx, cID, parameters)
	errs = append(errs, err...)
	if pid == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c, 30301005, fmt.Sprintf("add the information of host %s into DB error", requestData.Hostname))
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	nextCommandID++
	yid, err := addYumHostToDB(tx, &requestData, hostid, nextCommandID)
	errs = append(errs, err...)
	if yid == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c, 30301006, fmt.Sprintf("add the information of host %s into DB error", requestData.Hostname))
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	e := tx.Commit()
	if e != nil {
		err := apiutils.SendResponseForErrorMessage(c, 30301007, fmt.Sprintf("add the information of host %s into DB error %s", requestData.Hostname, e))
		errs = append(errs, err...)
		logErrors(errs)
	}
	msg := "host has be added successful"
	err = apiutils.SendResponseForSuccessMessage(c, msg)
	errs = append(errs, err...)
	logErrors(errs)
}

/*
validating the valid of the data what have received on http context
return nil and sysadmerror.Sysadmerror if the data validated.
otherwise return error and  sysadmerror.Sysadmerror
*/
func checkAddHostData(data *ApiHost, c *sysadmServer.Context) (error, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	// hostname should be large 1 and less 6
	multiform, e := c.MultipartForm()
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020005, "error", "get multipart form data error %s", e))
		return e, errs
	}
	formData := multiform.Value

	hostnameSlice, ok := formData["hostname"]
	errMsg, hostname := "", ""
	if !ok {
		errMsg = fmt.Sprintf("can not got hostname")
	} else {
		hostname = strings.TrimSpace(hostnameSlice[0])
		if len(hostname) < 1 || len(hostname) > 255 {
			errMsg = fmt.Sprintf("hostname %s is not valid.", hostname)
		}
	}
	if errMsg != "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020006, "error", errMsg))
		return fmt.Errorf("%s", errMsg), errs
	}
	data.Hostname = hostname

	// check userid
	userid := 0
	errMsg = ""
	useridSlice, ok := formData["userid"]
	if ok {
		tmpid := useridSlice[0]
		tmpUserid, e := strconv.Atoi(tmpid)
		if e == nil {
			userid = tmpUserid
		}
	}
	if userid == 0 {
		tmpid, e := user.GetSessionValue(c, "userid", WorkingData.sessionOption.sessionName)
		if e == nil {
			tmpUserid, e := utils.Interface2Int(tmpid)
			if e == nil {
				userid = tmpUserid
			} else {
				errMsg = fmt.Sprintf("添加主机信息需要登陆之后操作")
			}
		} else {
			errMsg = fmt.Sprintf("添加主机信息需要登陆之后操作")
		}
	}
	if errMsg != "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020007, "error", errMsg))
		return fmt.Errorf("%s", errMsg), errs
	}
	data.Userid = userid

	// checking agent running mode
	passiveModeSlice, ok := formData["passiveMode"]
	if !ok {
		data.PassiveMode = 0
	} else {
		passiveMode := strings.TrimSpace(passiveModeSlice[0])
		if passiveMode == "1" {
			data.PassiveMode = 1
		} else {
			data.PassiveMode = 0
		}
	}

	// checking ip
	ip := ""
	ipSlice, ok := formData["ip"]
	if ok {
		tmpIpStr := ipSlice[0]
		tmpIP, _ := utils.CheckIpAddress(tmpIpStr, false)
		if tmpIP != nil {
			ip = tmpIpStr
		}
	}
	if data.PassiveMode == 0 && ip == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020008, "error", "host ip must be specified when agent run in active mode"))
		return fmt.Errorf("host ip must be specified when agent run in active mode"), errs
	}
	data.Ip = ip

	// checking ip type
	ipType := ""
	ipTypeSlice, ok := formData["iptype"]
	if ok {
		tmpIpType := strings.TrimSpace(ipTypeSlice[0])
		if tmpIpType == "4" || tmpIpType == "6" {
			ipType = tmpIpType
		}
	}
	if data.Ip != "" && (ipType != "4" || ipType != "6") {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020009, "error", "host ip type(%s) is not valid", ipType))
		return fmt.Errorf("host ip type(%s) is not valid", ipType), errs
	}
	data.Iptype = ipType

	if data.PassiveMode == 1 {
		data.AgentPort = 0
		data.AgentIsTls = 0
		data.InsecureSkipVerify = 0
		data.CommandUri = ""
		data.CommandStatusUri = ""
		data.CommandLogsUri = ""
		data.AgentCa = ""
		data.AgentCert = ""
		data.AgentKey = ""
	} else {
		// get agent port if passive mode is false
		portSlice, ok := formData["agentPort"]
		portStr := portSlice[0]
		port, e := strconv.Atoi(portStr)
		if !ok || port < 1 || port > 65535 || e != nil {
			err := fmt.Errorf("Agent port(%s) is not valid", portStr)
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020009, "warn", "Agent port(%s) is not valid", portStr))
			return err, errs
		}
		data.AgentPort = port

		// get uris if passive mode is false
		commandUriSlice, commandOk := formData["commandUri"]
		if !commandOk {
			data.CommandUri = ""
		} else {
			data.CommandUri = strings.TrimSpace(commandUriSlice[0])
		}

		statusUriSlice, statusOk := formData["commandStatusUri"]
		if !statusOk {
			data.CommandStatusUri = ""
		} else {
			data.CommandStatusUri = strings.TrimSpace(statusUriSlice[0])
		}

		logSlice, logOk := formData["commandLogsUri"]
		if !logOk {
			data.CommandLogsUri = ""
		} else {
			data.CommandLogsUri = strings.TrimSpace(logSlice[0])
		}

		agentIsTlsSlice, ok := formData["agentIsTls"]
		if !ok {
			data.AgentIsTls = 0
		}
		agentIsTlsStr := strings.TrimSpace(agentIsTlsSlice[0])
		if agentIsTlsStr == "1" {
			data.AgentIsTls = 1
		} else {
			data.AgentIsTls = 0
		}

		if data.AgentIsTls == 1 {
			agentCa, _ := c.FormFile("agentCa")
			caContent, e := utils.ReadUploadedFile(agentCa)
			if e == nil {
				data.AgentCa = string(caContent)
			} else {
				data.AgentCa = ""
			}

			agentCert, _ := c.FormFile("agentCert")
			certContent, e := utils.ReadUploadedFile(agentCert)
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020011, "error", "upload certification file error %s", e))
				return e, errs
			}
			data.AgentCert = string(certContent)

			agentKey, _ := c.FormFile("agentKey")
			keyContent, e := utils.ReadUploadedFile(agentKey)
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020012, "error", "upload key file error %s", e))
				return e, errs
			}
			data.AgentKey = string(keyContent)

			insecureSkipVerifySlice, ok := formData["insecureSkipVerify"]
			if !ok {
				data.InsecureSkipVerify = 0
			} else {
				insecureSkipVerify := strings.TrimSpace(insecureSkipVerifySlice[0])
				if insecureSkipVerify == "1" {
					data.InsecureSkipVerify = 1
				} else {
					data.InsecureSkipVerify = 0
				}
			}
		} else {
			data.AgentCa = ""
			data.AgentCert = ""
			data.AgentKey = ""
			data.InsecureSkipVerify = 0
		}
	}

	// get os ID
	osIDSlice, ok := formData["osID"]
	if !ok {
		err := fmt.Errorf("parsing post data error")
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020012, "err", "%s", err))
		return err, errs
	}
	osIDStr := strings.TrimSpace(osIDSlice[0])
	osID, e := strconv.Atoi(osIDStr)
	if osID == 0 || e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020013, "error", "the os for node should be selected"))
		return fmt.Errorf("the os for node should be selected"), errs
	}
	data.OsID = osID

	// get os version ID
	osVersionIDSlice, ok := formData["osversionid"]
	if !ok {
		err := fmt.Errorf("parsing post data error")
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020014, "err", "%s", err))
		return err, errs
	}
	osVersionIDStr := strings.TrimSpace(osVersionIDSlice[0])
	osVersionID, e := strconv.Atoi(osVersionIDStr)
	if osVersionID == 0 || e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020015, "error", "the version of the os for node should be selected"))
		return fmt.Errorf("the version of the os for node should be selected"), errs
	}
	data.OsVersionID = osVersionID

	yumIDs, okForm := formData["yumid[]"]
	if okForm && len(yumIDs) > 0 {
		data.YumID = yumIDs
	} else {
		data.YumID = []string{}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020013, "warn", "no yum configuration has be selected"))
	}

	return nil, errs
}

/*
	addHostToDB add host information into DB

return hostid(which just added) and []sysadmerror.Sysadmerror
otherwise return 0  and []sysadmerror.Sysadmerror
*/
func addHostToDB(tx *db.Tx, data *ApiHost) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	// adding host information into DB
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(3030004, "debug", "now try to insert host information into host table"))

	// get the last hostid which will be added
	nextHostid, err := getNextID("host", "hostid", tx)
	errs = append(errs, err...)
	if nextHostid == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303018, "error", "can not get next hostid"))
		return 0, errs
	}

	insertData := make(db.FieldData, 0)
	insertData["userid"] = data.Userid
	insertData["projectid"] = 0
	insertData["hostname"] = data.Hostname
	insertData["osID"] = data.OsID
	insertData["osversionid"] = data.OsVersionID
	insertData["status"] = "run"
	insertData["ip"] = data.Ip
	insertData["iptype"] = data.Iptype
	insertData["passiveMode"] = data.PassiveMode
	insertData["commandUri"] = data.CommandUri
	insertData["commandStatusUri"] = data.CommandStatusUri
	insertData["commandLogsUri"] = data.CommandLogsUri
	insertData["agentIsTls"] = data.AgentIsTls
	insertData["agentCa"] = data.AgentCa
	insertData["agentCert"] = data.AgentCert
	insertData["agentKey"] = data.AgentKey
	insertData["insecureSkipVerify"] = data.InsecureSkipVerify
	insertData["agentPort"] = data.AgentPort
	insertData["hostid"] = nextHostid

	_, err = tx.InsertData("host", insertData)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
		return 0, errs
	}

	// update hostid value id ids table using transaction
	tmpID, err := updateNextID("host", "hostid", tx, (nextHostid + 1))
	errs = append(errs, err...)
	if tmpID == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303019, "error", "can not update next hostid"))
		return 0, errs
	}

	return nextHostid, errs
}

/*
	addManageIPToDB add manage IP of host into DB

return ipID(which just added) and []sysadmerror.Sysadmerror
otherwise return 0  and []sysadmerror.Sysadmerror
the IP information(such as devName, mask) of the host should be updated when we got the details of the IP information of the host
*/
func addManageIPToDB(tx *db.Tx, data *ApiHost, hostid int) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303010, "debug", "try to add host %s manager IP to DB", data.Hostname))
	insertData := make(db.FieldData, 0)

	insertData["devName"] = ""
	if data.Iptype == "4" {
		insertData["ipv4"] = data.Ip
		insertData["maskv4"] = ""
		insertData["ipv6"] = ""
		insertData["maskv6"] = ""
	} else {
		insertData["ipv4"] = ""
		insertData["maskv4"] = ""
		insertData["ipv6"] = data.Ip
		insertData["maskv6"] = ""
	}
	insertData["hostid"] = hostid
	insertData["status"] = 1
	insertData["isManage"] = 1

	_, err := tx.InsertData("hostIP", insertData)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
		return 0, errs
	}

	return 1, errs
}

/*
addYumHostToDB add relation information between host and yum into DB
return 1 and []sysadmerror.Sysadmerror
otherwise return 0  and []sysadmerror.Sysadmerror
*/
func addYumHostToDB(tx *db.Tx, data *ApiHost, hostid, nextCommandID int) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303010, "debug", "try to add yum host relations to DB"))
	insertData := make(db.FieldData, 0)

	insertData["hostid"] = hostid
	yumIDS := data.YumID
	for _, yID := range yumIDS {
		insertData["yumid"] = yID
		_, err := tx.InsertData("hostYum", insertData)
		errs = append(errs, err...)
		if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
			return 0, errs
		}

		cid, err := addYumConfigCommandToDB(tx, hostid, nextCommandID, yID, data)
		errs = append(errs, err...)
		if cid == 0 {
			return 0, errs
		}
		nextCommandID++
	}

	nextID, err := updateNextID("command", "commandID", tx, nextCommandID)
	errs = append(errs, err...)
	if nextID == 0 {
		return 0, errs
	}

	return 1, errs
}

func addYumConfigCommandToDB(tx *db.Tx, hostid, nextCommandID int, yumid string, data *ApiHost) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303011, "debug", "try to add yum configuration command  to DB"))

	apiServer := WorkingData.apiServer
	apiVersion := apiServer.ApiVersion
	tls := apiServer.Tls.IsTls
	address := apiServer.Server.Address
	port := apiServer.Server.Port
	ca := apiServer.Tls.Ca
	cert := apiServer.Tls.Cert
	key := apiServer.Tls.Key

	// get yum information list
	moduleName := "yum"
	actionName := "yumlist"
	apiServerData := apiutils.BuildApiServerData(moduleName, actionName, apiVersion, tls, address, port, ca, cert, key)
	urlRaw, _ := apiutils.BuildApiUrl(apiServerData)

	requestParas := httpclient.RequestParams{}
	requestParas.Url = urlRaw
	requestParas.Method = http.MethodPost
	requestParasP, err := httpclient.AddQueryData(&requestParas, "yumid", yumid)
	errs = append(errs, err...)
	body, err := httpclient.SendRequest(requestParasP)
	errs = append(errs, err...)
	ret, err := apiutils.ParseResponseBody(body)
	errs = append(errs, err...)
	if !ret.Status {
		message := ret.Message
		messageLine := message[0]
		msg := utils.Interface2String(messageLine["msg"])
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303012, "error", "get yum information list error %s", msg))
		return 0, errs
	}

	// add configurate yum command into DB
	line := ret.Message[0]
	yumName := utils.Interface2String(line["name"])
	yumCatalog := utils.Interface2String(line["catalog"])
	base_url := utils.Interface2String(line["base_url"])
	enabled, _ := utils.Interface2Int(line["enabled"])
	gpgcheck := utils.Interface2String(line["gpgcheck"])
	gpgkey := utils.Interface2String(line["gpgkey"])
	if enabled == 1 {
		// insert command information into command
		cID, err := addCommandToDB(tx, "addyum", 0, hostid, nextCommandID)
		errs = append(errs, err...)
		if cID == 0 {
			return 0, errs
		}

		parameters := make(map[string]string, 0)
		parameters["yumName"] = yumName
		parameters["yumCatalog"] = yumCatalog
		parameters["base_url"] = base_url
		parameters["gpgcheck"] = gpgcheck
		parameters["gpgkey"] = gpgkey
		pid, err := addCommandParametersToDB(tx, cID, parameters)
		errs = append(errs, err...)
		if pid == 0 {
			return 0, errs
		}
	}

	return 1, errs
}

// get next ID from ids table for tableName with fieldName
// return the ID value and []sysadmerror.Sysadmerror if successfule
// otherewise return 0 and []sysadmerror.Sysadmerror
func getNextID(tableName, fieldName string, tx *db.Tx) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303013, "debug", "try to get table %s next ID with field %s", tableName, fieldName))

	if strings.TrimSpace(tableName) == "" || strings.TrimSpace(fieldName) == "" {
		return 0, append(errs, sysadmerror.NewErrorWithStringLevel(30303014, "error", "table name %s or field name %s is empty", tableName, fieldName))
	}

	whereStatement := make(map[string]string, 0)
	whereStatement["tableName"] = " = '" + tableName + "'"
	whereStatement["fieldName"] = " = '" + fieldName + "'"
	selectData := db.SelectData{
		Tb:        []string{"ids"},
		OutFeilds: []string{"nextValue"},
		Where:     whereStatement,
	}

	dbEntity := WorkingData.dbConf.Entity
	retData, err := dbEntity.QueryData(&selectData)
	errs = append(errs, err...)
	if retData == nil {
		return 0, errs
	}

	nextID := 0
	line := retData[0]
	if v, ok := line["nextValue"]; ok {
		id, e := utils.Interface2Int(v)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303015, "error", "can not get next ID with table %s and field %s error %s", tableName, fieldName, e))
			return 0, errs
		} else {
			nextID = id
		}
	} else {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303016, "error", "can not get next ID with table %s and field %s", tableName, fieldName))
		return 0, errs
	}

	return nextID, errs
}

// update next ID into ids table for tableName with fieldName
// return the 1 and []sysadmerror.Sysadmerror if successfule
// otherewise return 0 and []sysadmerror.Sysadmerror
func updateNextID(tableName, fieldName string, tx *db.Tx, nextID int) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303016, "debug", "now try to update nextID(%d) for table %s with field %s ", nextID, tableName, fieldName))

	// update commandID value in ids table using transaction
	updateData := make(db.FieldData, 0)
	updateData["nextValue"] = nextID
	whereStatement := make(map[string]string, 0)
	whereStatement["tableName"] = tableName
	whereStatement["fieldName"] = fieldName

	_, err := tx.UpdateData("ids", updateData, whereStatement)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
		return 0, errs
	}

	return 1, errs
}
