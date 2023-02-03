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

package sysadmapi

import (
	"fmt"
	"strings"

	"encoding/json"
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

/*
 * 检查代码是否是合法的命令状态代码，如果是合法的状态代码则返回true, 否则返回false
 */
func IsCommandStatusCodeValid(code CommandStatusCode) bool {
	for _, v := range AllCommandStatusCode {
		if v == code {
			return true
		}
	}

	return false
}
