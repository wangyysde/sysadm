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
	"strconv"
	"strings"
	"time"

	apiServerApp "sysadm/apiserver/app"
	"sysadm/db"
	"sysadm/sysadmerror"
	"sysadm/utils"
)

/*
addCommandToDB add the information of a command into DB
return 1 and []sysadmerror.Sysadmerror
otherwise return 0  and []sysadmerror.Sysadmerror
*/
func addCommandToDB(tx *db.Tx, command string, synchronized int, hostid, nextCommandID int) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30304001, "debug", "try to add the information of command %s to DB", command))

	insertData := make(db.FieldData, 0)
	insertData["commandID"] = nextCommandID
	insertData["command"] = command
	insertData["hostID"] = hostid
	insertData["synchronized"] = synchronized

	now := time.Now()
	nowInt64 := now.Unix()
	insertData["createTime"] = nowInt64
	insertData["tryTimes"] = 0
	insertData["status"] = apiServerApp.CommandStatusCreated

	_, err := tx.InsertData("command", insertData)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error") {
		return 0, errs
	}

	return nextCommandID, errs

}

/*
addCommandParametersToDB add the information of command parameters into DB
return 1 and []sysadmerror.Sysadmerror
otherwise return 0  and []sysadmerror.Sysadmerror
*/
func addCommandParametersToDB(tx *db.Tx, commandID int, paras map[string]string) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30304002, "debug", "try to add the information of parameters for command to DB"))
	for key, value := range paras {
		insertData := make(db.FieldData, 0)
		insertData["name"] = key
		insertData["value"] = value
		insertData["commandID"] = commandID

		_, err := tx.InsertData("commandParameters", insertData)
		errs = append(errs, err...)
		if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error") {
			return 0, errs
		}
	}

	return 1, errs
}

func GetCommand(ips, macs []string, hostname, customize string) (*apiServerApp.Command, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	commandSeq := buildCommandSeq("")
	ret := &apiServerApp.Command{
		CommandSeq: commandSeq,
		Command:    "",
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30304003, "debug", "get commands from DB with ips %+v macs %+v hostname %s customize %s", ips, macs, hostname, customize))

	if len(ips) < 1 && len(macs) < 1 && strings.TrimSpace(hostname) == "" && strings.TrimSpace(customize) == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30304004, "error", "can not get commands for a node without ip,mac,hostname and customize"))
		return ret, errs
	}

	dbEntity := WorkingData.dbConf.Entity

	// priority of customize is highest
	customize = strings.TrimSpace(customize)
	if customize != "" {
		_, e := strconv.Atoi(customize)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(30304005, "error", "customize is not empty,but it can not change to hostid, error %s", e))
			return ret, errs
		}

		commands, err := getCommandForhHostid(customize)
		errs = append(errs, err...)
		if commands != nil {
			return commands, errs
		}
	}

	// priority of mac is second
	if len(macs) > 0 {
		whereStatement := make(map[string]string, 0)
		macWhereStatement := dbEntity.BuildWhereFieldExactWithSlice(macs)
		whereStatement["mac"] = macWhereStatement

		selectData := db.SelectData{
			Tb:        []string{"hostMAC"},
			OutFeilds: []string{"hostid"},
			Where:     whereStatement,
		}
		retData, err := dbEntity.QueryData(&selectData)
		errs = append(errs, err...)
		if len(retData) > 0 {
			lineData := retData[0]
			hostid := utils.Interface2String(lineData["hostid"])
			commands, err := getCommandForhHostid(hostid)
			errs = append(errs, err...)
			if commands != nil {
				return commands, errs
			}
		}
	}

	// priority of mac is third
	if len(ips) > 0 {
		var ipv4Str, ipv6Str string = "", ""
		for _, ip := range ips {
			_, ipv4OrIpv6 := utils.JudgeIpv4OrIpv6(ip)
			if ipv4OrIpv6 == 0 {
				continue
			}
			if ipv4OrIpv6 == 4 {
				if ipv4Str == "" {
					ipv4Str = ip
				} else {
					ipv4Str = ipv4Str + "," + ip
				}
			} else {
				if ipv6Str == "" {
					ipv6Str = ip
				} else {
					ipv6Str = ipv6Str + "," + ip
				}
			}
		}

		whereStatement := make(map[string]string, 0)
		if ipv4Str != "" {
			whereStatement["ipv4"] = dbEntity.BuildWhereFieldExact(ipv4Str)
			whereStatement["status"] = "='1'"
			whereStatement["isManage"] = "='1'"
			selectData := db.SelectData{
				Tb:        []string{"hostIP"},
				OutFeilds: []string{"hostid"},
				Where:     whereStatement,
			}
			retData, err := dbEntity.QueryData(&selectData)
			errs = append(errs, err...)
			if len(retData) > 0 {
				lineData := retData[0]
				hostid := utils.Interface2String(lineData["hostIP"])
				commands, err := getCommandForhHostid(hostid)
				errs = append(errs, err...)
				if commands != nil {
					return commands, errs
				}
			}
		}

		if ipv6Str != "" {
			whereStatement["ipv6"] = dbEntity.BuildWhereFieldExact(ipv6Str)
			whereStatement["status"] = "='1'"
			whereStatement["isManage"] = "='1'"
			selectData := db.SelectData{
				Tb:        []string{"hostIP"},
				OutFeilds: []string{"hostid"},
				Where:     whereStatement,
			}
			retData, err := dbEntity.QueryData(&selectData)
			errs = append(errs, err...)
			if len(retData) > 0 {
				lineData := retData[0]
				hostid := utils.Interface2String(lineData["hostIP"])
				commands, err := getCommandForhHostid(hostid)
				errs = append(errs, err...)
				if commands != nil {
					return commands, errs
				}
			}
		}
	}

	// the last one is hostname
	hostname = strings.TrimSpace(hostname)
	if hostname != "" {
		whereStatement := make(map[string]string, 0)
		whereStatement["hostname"] = " = '" + hostname + "'"
		whereStatement["statusID"] = " = 1 "
		selectData := db.SelectData{
			Tb:        []string{"host"},
			OutFeilds: []string{"hostid"},
			Where:     whereStatement,
		}
		retData, err := dbEntity.QueryData(&selectData)
		errs = append(errs, err...)
		if len(retData) > 0 {
			lineData := retData[0]
			hostid := lineData["hostid"]
			commands, err := getCommandForhHostid(utils.Interface2String(hostid))
			errs = append(errs, err...)
			if commands != nil {
				return commands, errs
			}
		}
	}

	return ret, errs
}

func getCommandForhHostid(hostid string) (*apiServerApp.Command, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	var tmpRet apiServerApp.Command
	var ret *apiServerApp.Command = nil

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30304010, "debug", "get commands from DB with hostid %s", hostid))

	hostid = strings.TrimSpace(hostid)
	if hostid == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30304011, "error", "can not get command with no host"))
		return ret, errs
	}

	dbEntity := WorkingData.dbConf.Entity
	whereStatement := make(map[string]string, 0)
	whereStatement["hostID"] = " = '" + hostid + "'"
	whereStatement["tryTimes"] = " <3"
	whereStatement["status"] = " = " + strconv.Itoa(int(apiServerApp.CommandStatusCreated))
	selectData := db.SelectData{
		Tb:        []string{"command"},
		OutFeilds: []string{"commandID", "command", "synchronized"},
		Where:     whereStatement,
		Order:     []db.OrderData{{Key: "commandID", Order: 0}},
	}
	retData, err := dbEntity.QueryData(&selectData)
	errs = append(errs, err...)
	if retData == nil || len(retData) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30304012, "error", "can not get command informations with hostid %s", hostid))
		return ret, errs
	}

	lineData := retData[0]
	commandID := utils.Interface2String(lineData["commandID"])
	command := utils.Interface2String(lineData["command"])
	s, e := utils.Interface2Int(lineData["synchronized"])
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30304012, "error", "synchronized of command is not valid. error %s", e))
		return ret, errs
	}
	synchronized := true
	if s != 0 {
		synchronized = false
	}

	commandSeq := buildCommandSeq(commandID)
	tmpRet = apiServerApp.Command{
		CommandSeq:   commandSeq,
		Command:      command,
		Synchronized: synchronized,
	}

	pwhere := make(map[string]string, 0)
	pwhere["commandID"] = " = '" + commandID + "'"
	selectData = db.SelectData{
		Tb:        []string{"commandParameters"},
		OutFeilds: []string{"name", "value"},
		Where:     pwhere,
	}
	pData, err := dbEntity.QueryData(&selectData)
	errs = append(errs, err...)
	paras := make(map[string]string, 0)
	for _, l := range pData {
		key := utils.Interface2String(l["name"])
		value := utils.Interface2String(l["value"])
		paras[key] = value
	}
	tmpRet.Parameters = paras

	ret = &tmpRet

	return ret, errs
}

func buildCommandSeq(commandID string) string {
	commandID = strings.TrimSpace(commandID)
	if commandID == "" {
		return "0000000000000000000"
	}
	commandSeq := time.Now().Format("20060102")
	for i := len(commandID); i > 0; i++ {
		commandSeq = commandSeq + "0"
	}

	commandSeq = commandSeq + commandID

	return commandSeq
}
