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
*
* NOTE:
* 本文件定义的是用于处理服务器端要求客户执行的指令相关的处理函数或方法
 */

package app

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/wangyysde/sysadmLog"
	"sysadm/utils"
)

/*
 * MarshalCommandData marshal CommandData to json encode. return []byte and nil
 * otherwise return  []byte and error
 */
func MarshalCommandData(commandData *CommandData) ([]byte, error) {
	var ret []byte

	if commandData == nil {
		return ret, fmt.Errorf("can not marshal a nil to byte")
	}

	if strings.TrimSpace(commandData.CommandSeq) == "" || strings.TrimSpace(commandData.Command.Command) == "" {
		return ret, fmt.Errorf("command sequence and command name must not be empty")
	}

	return json.Marshal(commandData)
}

/*
 * UnMarshalCommandData unmarshal json encoded data to *CommandData
 * return  *CommandData and nil if successful, otherwise return nil and error
 */
func UnMarshalCommandData(data []byte) (*CommandData, error) {
	var retP *CommandData = nil

	if len(data) < 1 {
		return retP, fmt.Errorf("the data to be unMarshal must not be empty")
	}

	var ret CommandData = CommandData{}
	err := json.Unmarshal(data, &ret)

	return &ret, err
}

/*
 * MarshalCommandReq try to marshal commandReq to json encoded. return []byte and nil
 * otherwise return  []byte and error
 */
func MarshalCommandReq(commandReq *CommandReq) ([]byte, error) {
	var ret []byte

	if commandReq == nil {
		return ret, fmt.Errorf("can not marshal a nil to byte")
	}

	if len(commandReq.Ips) < 1 && len(commandReq.Macs) < 1 && strings.TrimSpace(commandReq.Hostname) == "" && strings.TrimSpace(commandReq.Customize) == "" {
		return ret, fmt.Errorf("at least of fields IP, MAC, HOSTNAME or Customize must be not empty")
	}

	return json.Marshal(commandReq)
}

/*
* UnmarshalCommandReq parsing the body to *CommandReq. Normally body is a request body what a client get command to run
* return  *CommandReq and nil if successful, otherwise return nil and error
 */
func UnmarshalCommandReq(body []byte) (*CommandReq, error) {
	var retP *CommandReq = nil

	if len(body) < 1 {
		return retP, fmt.Errorf("the body to be parsed is empty")
	}

	var ret CommandReq = CommandReq{}
	err := json.Unmarshal(body, &ret)

	return &ret, err
}

/*
 * MarshalCommandStatusReq marshal commandStatusReq to json encode. return []byte and nil
 * otherwise return  []byte and error
 */
func MarshalCommandStatusReq(commandStatusReq *CommandStatusReq) ([]byte, error) {
	var ret []byte

	if commandStatusReq == nil {
		return ret, fmt.Errorf("can not marshal a nil to byte")
	}

	if strings.TrimSpace(commandStatusReq.CommandSeq) == "" {
		return ret, fmt.Errorf("command sequence must not be empty")
	}

	return json.Marshal(commandStatusReq)
}

/*
 * UnMarshalCommandStatusReq unmarshal json encoded data to *CommandStatusReq
 * return  *CommandStatusReq and nil if successful, otherwise return nil and error
 */
func UnMarshalCommandStatusReq(data []byte) (*CommandStatusReq, error) {
	var retP *CommandStatusReq = nil

	if len(data) < 1 {
		return retP, fmt.Errorf("the data to be unMarshal must not be empty")
	}

	var ret CommandStatusReq = CommandStatusReq{}
	err := json.Unmarshal(data, &ret)

	return &ret, err
}

/*
 * MarshalCommandStatus marshal CommandStatus to json encode. return []byte and nil
 * otherwise return  []byte and error
 */
func MarshalCommandStatus(commandStatus *CommandStatus) ([]byte, error) {
	var ret []byte

	if commandStatus == nil {
		return ret, fmt.Errorf("can not marshal a nil to byte")
	}

	if strings.TrimSpace(commandStatus.CommandSeq) == "" {
		return ret, fmt.Errorf("command sequence must not be empty")
	}

	if len(commandStatus.Ips) < 1 && len(commandStatus.Macs) < 1 && strings.TrimSpace(commandStatus.Hostname) == "" && strings.TrimSpace(commandStatus.Customize) == "" {
		return ret, fmt.Errorf("at least of fields IP, MAC, HOSTNAME or Customize must be not empty")
	}

	if !IsCommandStatusCodeValid(commandStatus.StatusCode) {
		return ret, fmt.Errorf("command status code is not valid")
	}

	return json.Marshal(commandStatus)
}

/*
 * UnMarshalCommandStatus unmarshal json encoded data to *CommandStatus
 * return  *CommandStatus and nil if successful, otherwise return nil and error
 */
func UnMarshalCommandStatus(data []byte) (*CommandStatus, error) {
	var retP *CommandStatus = nil

	if len(data) < 1 {
		return retP, fmt.Errorf("the data to be unMarshal must not be empty")
	}

	var ret CommandStatus = CommandStatus{}
	err := json.Unmarshal(data, &ret)

	return &ret, err
}

/*
 * MarshalLogReq marshal LogReq to json encode. return []byte and nil
 * otherwise return  []byte and error
 */
func MarshalLogReq(logReq *LogReq) ([]byte, error) {
	var ret []byte

	if logReq == nil {
		return ret, fmt.Errorf("can not marshal a nil to byte")
	}

	if strings.TrimSpace(logReq.CommandSeq) == "" || strings.TrimSpace(logReq.CommandSeq) == "0000000000000000000" {
		return ret, fmt.Errorf("command sequence must not be empty")
	}

	if strings.TrimSpace(logReq.StartSeq) == "" {
		return ret, fmt.Errorf("start log seqence must not be empty")
	}
	if logReq.Num == 0 {
		return ret, fmt.Errorf("max count of server can received must be not zero")
	}

	return json.Marshal(logReq)
}

/*
 * UnMarshallLogReq unmarshal json encoded data to *LogReq
 * return  *LogReq and nil if successful, otherwise return nil and error
 */
func UnMarshallLogReq(data []byte) (*LogReq, error) {
	var retP *LogReq = nil

	if len(data) < 1 {
		return retP, fmt.Errorf("the data to be unMarshal must not be empty")
	}

	var ret LogReq = LogReq{}
	err := json.Unmarshal(data, &ret)

	return &ret, err
}

/*
 * MarshalLogData marshal LogData to json encode. return []byte and nil
 * otherwise return  []byte and error
 */
func MarshalLogData(logData *LogData) ([]byte, error) {
	var ret []byte

	if logData == nil {
		return ret, fmt.Errorf("can not marshal a nil to byte")
	}

	if strings.TrimSpace(logData.CommandSeq) == "" || strings.TrimSpace(logData.CommandSeq) == "0000000000000000000" {
		return ret, fmt.Errorf("command sequence must not be empty")
	}

	if len(logData.Ips) < 1 && len(logData.Macs) < 1 && strings.TrimSpace(logData.Hostname) == "" && strings.TrimSpace(logData.Customize) == "" {
		return ret, fmt.Errorf("at least of fields IP, MAC, HOSTNAME or Customize must be not empty")
	}

	return json.Marshal(logData)
}

/*
 * UnMarshalCommandStatus unmarshal json encoded data to *LogData
 * return  *LogData and nil if successful, otherwise return nil and error
 */
func UnMarshallLogData(data []byte) (*LogData, error) {
	var retP *LogData = nil

	if len(data) < 1 {
		return retP, fmt.Errorf("the data to be unMarshal must not be empty")
	}

	var ret LogData = LogData{}
	err := json.Unmarshal(data, &ret)

	return &ret, err
}

/*
 * MarshalRepStatus marshal RepStatus to json encode. return []byte and nil
 * otherwise return  []byte and error
 */
func MarshalRepStatus(repStatus *RepStatus) ([]byte, error) {
	var ret []byte

	if repStatus == nil {
		return ret, fmt.Errorf("can not marshal a nil to byte")
	}

	if strings.TrimSpace(repStatus.CommandSeq) == "" {
		return ret, fmt.Errorf("command sequence must not be empty")
	}

	if !IsCommandStatusCodeValid(repStatus.StatusCode) {
		return ret, fmt.Errorf("command status code is not valid")
	}

	return json.Marshal(repStatus)
}

/*
 * UnMarshalRepStatus unmarshal json encoded data to *RepStatus
 * return  *RepStatus and nil if successful, otherwise return nil and error
 */
func UnMarshalRepStatus(data []byte) (*RepStatus, error) {
	var retP *RepStatus = nil

	if len(data) < 1 {
		return retP, fmt.Errorf("the data to be unMarshal must not be empty")
	}

	var ret RepStatus = RepStatus{}
	err := json.Unmarshal(data, &ret)

	return &ret, err
}

// 检查代码是否是合法的命令状态代码，如果是合法的状态代码则返回true, 否则返回false
func IsCommandStatusCodeValid(code CommandStatusCode) bool {
	for _, v := range AllCommandStatusCode {
		if v == code {
			return true
		}
	}

	return false
}

// GetCommandStatusCodeByInt get command status code by int.
// return CommandStatusCode if int is in AllCommandStatusCode
// otherwise return CommandStatusUnkown
func GetCommandStatusCodeByInt(code int) CommandStatusCode {
	for _, v := range AllCommandStatusCode {
		if int(v) == code {
			return v
		}
	}

	return CommandStatusUnkown
}

func BuildCommandStatus(commandSeq, nodeIdentifierStr, statusMessage string, nodeIdentifier NodeIdentifier, statusCode CommandStatusCode,
	data map[string]interface{}, notCommand bool) (CommandStatus, error) {
	var e error = nil
	if strings.TrimSpace(commandSeq) == "" {
		commandSeq = "0000000000000000000"
	}

	if reflect.DeepEqual(NodeIdentifier{}, nodeIdentifier) {
		if strings.TrimSpace(nodeIdentifierStr) == "" {
			nodeIdentifierStr = DefaultNodeIdentifer
		}

		nodeIdentifier, e = BuildNodeIdentifer(nodeIdentifierStr)
	}

	return CommandStatus{
		CommandSeq:     commandSeq,
		NodeIdentifier: nodeIdentifier,
		StatusCode:     statusCode,
		StatusMessage:  statusMessage,
		Data:           data,
		NotCommand:     notCommand,
	}, e

}

func BuildNodeIdentifer(nodeIdentiferStr string) (NodeIdentifier, error) {
	if strings.TrimSpace(nodeIdentiferStr) == "" || len(strings.TrimSpace(nodeIdentiferStr)) > MaxCustomizeNodeIdentiferLen {
		nodeIdentiferStr = DefaultNodeIdentifer
	}

	ret := NodeIdentifier{}

	if !IsNodeIdentiferStrValid(nodeIdentiferStr) && len(nodeIdentiferStr) <= MaxCustomizeNodeIdentiferLen {
		ret.Customize = strings.TrimSpace(nodeIdentiferStr)
		return ret, nil
	}

	identiferSlice := strings.Split(nodeIdentiferStr, ",")
	for _, value := range identiferSlice {
		switch {
		case strings.ToUpper(strings.TrimSpace(value)) == "IP":
			ips, err := utils.GetLocalIPs()
			if err != nil {
				return ret, fmt.Errorf("get local host ip address error %s", err)
			}
			ret.Ips = ips
		case strings.ToUpper(strings.TrimSpace(value)) == "MAC":
			macs, err := utils.GetLocalMacs()
			if err != nil {
				return ret, fmt.Errorf("get local host mac information error %s", err)
			}
			ret.Macs = macs
		case strings.ToUpper(strings.TrimSpace(value)) == "HOSTNAME":
			hostname, err := os.Hostname()
			if err != nil {
				return ret, fmt.Errorf("can not get hostname %s", err)
			}
			ret.Hostname = hostname
		}
	}

	return ret, nil
}

/*
* IsNodeIdentiferStrValid check whether nodeIdentifierStr is a valid node identifer string
* nodeIdentifierStr should be a combination of "IP", "HOSTNAME" or "MAC"
 */
func IsNodeIdentiferStrValid(nodeIdentifierStr string) bool {
	if strings.TrimSpace(nodeIdentifierStr) == "" {
		return false
	}

	identiferSlice := strings.Split(nodeIdentifierStr, ",")
	for _, v := range identiferSlice {
		tmpV := strings.ToUpper(strings.TrimSpace(v))
		if tmpV != "IP" && tmpV != "HOSTNAME" && tmpV != "MAC" {
			return false
		}
	}

	return true
}

// IsCommandSeqValid check whether commandSeq is a valid command sequence.Note: "0000000000000000000" considered as a not valid command sequence
// return true if commandSeq is a valid command sequence
// otherwise return faluse
func IsCommandSeqValid(commandSeq string) bool {
	commandSeq = strings.TrimSpace(commandSeq)
	if len(commandSeq) != 19 {
		return false
	}
	if commandSeq == "0000000000000000000" {
		return false
	}

	return true
}

// BuildCommandSeqByID build command sequence with id.
// return 0000000000000000000 if id is empty.
// otherewise return YYYYMMDD0000ID
func BuildCommandSeqByID(id string) string {
	if strings.TrimSpace(id) == "" {
		return "0000000000000000000"
	}

	dateStr := time.Now().Local().Format("20060102")
	for i := len(id); i < 11; i++ {
		dateStr = dateStr + "0"
	}

	return (dateStr + id)
}

// GetCommandIDFromSeq get command id from a command sequence.
// return command id and nil if successful
// otherwise return "" and an error
func GetCommandIDFromSeq(seq string) (int, error) {
	if !IsCommandSeqValid(seq) {
		return 0, fmt.Errorf("command sequence %s is not valid", seq)
	}

	idStr := seq[8:]
	for {
		if idStr[0:] == "0" {
			idStr = idStr[1:]
		} else {
			break
		}
	}

	id, e := strconv.Atoi(idStr)
	if e != nil {
		return 0, e
	}
	return id, nil
}

// ConvCommandStatus2Map convert CommandStatus to map[string]interface{} for saving into redis server
// because redis can not save nested data, so ConvCommandStatus2Map extend all sub struct to one level
func ConvCommandStatus2Map(commandStatus *CommandStatus) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	if strings.TrimSpace(commandStatus.CommandSeq) == "" {
		return ret, fmt.Errorf("command sequerence is empty")
	}

	ipStr := strings.Join(commandStatus.NodeIdentifier.Ips, ",")
	macStr := strings.Join(commandStatus.Macs, ",")
	var dataStr string = ""
	if len(commandStatus.Data) > 0 {
		dataBytes, e := json.Marshal(commandStatus.Data)
		if e != nil {
			return ret, e
		}
		dataStr = string(dataBytes)
	}

	notCommand := 0
	if commandStatus.NotCommand {
		notCommand = 1
	}

	ret["commandSeq"] = commandStatus.CommandSeq
	ret["ips"] = ipStr
	ret["macs"] = macStr
	ret["hostname"] = commandStatus.NodeIdentifier.Hostname
	ret["customize"] = commandStatus.NodeIdentifier.Customize
	ret["statusCode"] = commandStatus.StatusCode
	ret["statusMessage"] = commandStatus.StatusMessage
	ret["data"] = dataStr
	ret["notCommand"] = notCommand

	return ret, nil
}

// ConvMap2CommandStatus convert map[string]interface{} to  CommandStatus
//
//	return CommandStatus{}, error if any error was occurred
//
// otherwise return CommandStatus, nil
func ConvMap2CommandStatus(data map[string]interface{}) (CommandStatus, error) {
	var ret CommandStatus = CommandStatus{}

	commandSeq, ok := data["commandSeq"]
	if !ok {
		return ret, fmt.Errorf("command sequence is empty")
	}

	ipsStr := data["ips"]
	ipSlice := strings.Split(utils.Interface2String(ipsStr), ",")

	macStr := data["macs"]
	macSlice := strings.Split(utils.Interface2String(macStr), ",")

	hostname := data["hostname"]
	customize := data["customize"]
	statusCode := data["statusCode"]

	statusMessage := data["statusMessage"]
	dataTmp := data["data"]
	dataStr := utils.Interface2String(dataTmp)
	dataBytes := utils.Str2bytes(dataStr)
	dataMap := make(map[string]interface{}, 0)
	e := json.Unmarshal(dataBytes, &dataMap)
	if e != nil {
		return ret, fmt.Errorf("can not unmarshal json to map")
	}

	notCommandTmp := data["notCommand"]
	notCommand, e := utils.Interface2Int(notCommandTmp)
	notCommandBool := false
	if notCommand > 0 {
		notCommandBool = true
	}
	if e != nil {
		return ret, fmt.Errorf("can not convert notCommand to int")
	}

	statusCodeInt, e := utils.Interface2Int(statusCode)
	if e != nil {
		return ret, fmt.Errorf("can not convert statusCode to int")
	}

	ret.CommandSeq = utils.Interface2String(commandSeq)
	ret.NodeIdentifier.Ips = ipSlice
	ret.NodeIdentifier.Macs = macSlice
	ret.NodeIdentifier.Hostname = utils.Interface2String(hostname)
	ret.NodeIdentifier.Customize = utils.Interface2String(customize)
	ret.StatusCode = GetCommandStatusCodeByInt(statusCodeInt)
	ret.StatusMessage = utils.Interface2String(statusMessage)
	ret.NotCommand = notCommandBool

	return ret, nil
}

// BuildLog  build an instance of Log used logID, message and level
func BuildLog(logID int, message string, level sysadmLog.Level) Log {
	logSeq := BuildLogSeqByID(logID)

	return Log{
		LogSeq:  logSeq,
		Level:   level,
		Message: message,
	}
}

// BuildLogSeqByID build log sequence with id.
// return YYYYMMDD111111 if id is less 1 or id is bigger 111111.
// otherewise return YYYYMMDD0000ID
func BuildLogSeqByID(id int) string {
	if id < 1 || id >= 111111 {
		id = 111111
	}

	dateStr := time.Now().Local().Format("20060102")
	idStr := strconv.Itoa(id)

	for i := len(idStr); i < 6; i++ {
		dateStr = dateStr + "0"
	}

	return (dateStr + idStr)
}

// GetLogIDFromSeq get log id from a log sequence.
// return date string in the sequence, id,and nil if successful
// otherwise return "","" and error
func GetLogIDFromSeq(seq string) (string, int, error) {
	if !IsLogSeqValid(seq) {
		return "", 0, fmt.Errorf("sequence %s is not valid", seq)
	}

	dateStr := seq[0:8]
	idStr := seq[8:]
	id, e := strconv.Atoi(idStr)
	if e != nil {
		return "", 0, fmt.Errorf("sequence %s is not valid", seq)
	}

	return dateStr, id, nil
}

// IsLogSeqValid check whether seq is a valid log sequence.Note: "00000000000000" considered as a not valid log sequence
// return true if seq is a valid log sequence
// otherwise return false
func IsLogSeqValid(seq string) bool {
	seq = strings.TrimSpace(seq)
	if len(seq) != 14 {
		return false
	}
	if seq == "00000000000000" {
		return false
	}

	return true
}
