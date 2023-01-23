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
	"regexp"
	"strconv"
	"strings"

	"github.com/wangyysde/sshclient/sftp"
	"github.com/wangyysde/sshclient/sshclient"
	"github.com/wangyysde/sshclient/sshcopyid"
	"github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadm/sysadmerror"
	syssetting "github.com/wangyysde/sysadm/syssetting/app"
	"github.com/wangyysde/sysadm/utils"
	"github.com/wangyysde/sysadmServer"
)

func addHost(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror

	var requestData ApiHost
	if e := c.ShouldBind(&requestData); e != nil {
		msg := fmt.Sprintf("get host data err %s", e) 
		_ = apiutils.SendResponseForErrorMessage(c,3030001,msg)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3030001,"error",msg))
		logErrors(errs)
		return 
	}
	
	if e,err := checkAddHostData(&requestData); e != nil {
		errs=append(errs,err...)
		msg := fmt.Sprintf("%s", e) 
		err = apiutils.SendResponseForErrorMessage(c,3030002,msg)
		errs = append(errs,err...)
		logErrors(errs)
		return 

	}

	// adding host information into DB
	dbConf := WorkingData.dbConf
	dbEntity := dbConf.Entity
	tx,err := db.Begin(dbEntity)
	errs = append(errs,err...)
	if tx == nil {
		err := apiutils.SendResponseForErrorMessage(c,3030003,"start transaction on DB query error. check logs for details.")
		errs = append(errs,err...)
		logErrors(errs)
		return 
	}
	
	// insert host information into host table	
	hostid,err := addHostToDB(tx,&requestData)
	errs = append(errs,err...)
	if hostid == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c,30301001,"add host error. check logs for details")
		errs = append(errs,err...)
		logErrors(errs)
		return 
	}

	// insert IP information into hostIP table
	ipID, err := addManageIPToDB(tx, &requestData, hostid)
	errs = append(errs,err...)
	if ipID == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c,30301002,"host %s manage IP %s add into DB error", requestData.Hostname,requestData.Ip)
		errs = append(errs,err...)
		logErrors(errs)
		return 
	}

	// insert the relation between yum and host into hostYum
	rID, err := addYumHostToDB(tx, &requestData, hostid)
	if rID == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c,30301003,"add the information of  host %s into DB error", requestData.Hostname)
		errs = append(errs,err...)
		logErrors(errs)
		return 
	}

	// insert command information into command
	cID, err := addCommandToDB(tx,"gethostip",requestData.PassiveMode, hostid)
	if cID == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c,30301004,"add the information of host %s into DB error", requestData.Hostname)
		errs = append(errs,err...)
		logErrors(errs)
		return 
	}

	parameters := make(map[string]string,0)
	parameters["withmac"] = "1"
	parameters["withmask"] = "1"
	pid, err := addCommandParametersToDB(tx,cID,parameters)
	if pid == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c,30301005,"add the information of host %s into DB error", requestData.Hostname)
		errs = append(errs,err...)
		logErrors(errs)
		return 
	}





	

	hostid,err = addIPToDB(tx,&requestData)

	


	msg := "host has be added successful"
	err = apiutils.SendResponseForSuccessMessage(c,msg)	
	errs=append(errs,err...)
	logErrors(errs)
}

/*
	validating the valid of the data what have received on http context
	return nil and sysadmerror.Sysadmerror if the data validated.
	otherwise return error and  sysadmerror.Sysadmerror
*/
func checkAddHostData(data *ApiHost)(error,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	// hostname should be large 1 and less 64
	hostname := strings.TrimSpace(data.Hostname) 
	if(len(hostname) < 1 || len(hostname) >63 ){
		errMsg := fmt.Sprintf("hostname %s is not valid.",hostname)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020006,"error",errMsg))
		return  fmt.Errorf(errMsg), errs
	}
	data.Hostname = hostname

	// checking ip
	ip := strings.TrimSpace(data.Ip)
	if tmpIP,e := utils.CheckIpAddress(ip,false); tmpIP == nil {
		return fmt.Errorf("host ip %s is not valid.",ip),e
	}
	data.Ip = ip

	ipType := strings.TrimSpace(data.Iptype)
	if ipType != "4" && ipType != "6"{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020007,"warn","ip type (%s) is not valid. we try to look node IP as ipv4.",ipType))
		ipType = "4"
	}
	data.Iptype = ipType

	passiveMode := data.PassiveMode
	if passiveMode {
		data.AgentPort = 0
		data.ReceiveCommandUri = ""
	} else {
		port := data.AgentPort
		if port < 1 || port > 65535 {
			err := fmt.Errorf("Agent port(%d) is not valid",port)
			errs = append(errs, err))
			return err,errs
		}

		receiveCommandUri := strings.TrimSpace(data.ReceiveCommandUri)
		if len(receiveCommandUri) < 1 {
			err := fmt.Errorf("agent listen uri path %s  is not valid",receiveCommandUri )
			errs = append(errs, err)
			return err,errs
		}
	}
	
	osID := data.OsID
	osVersionID := data.OsVersionID

	if osID == 0 || osVersionID == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020011,"error","os and version of the os for node should be selected"))
		return fmt.Errorf("os and version of the os for node should be selected"),errs
	}

	yumIDs := data.YumID
	if len(yumIDs) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020012,"warn","no yum configuration has be selected"))
	}

	return nil, errs
}

/*
   addHostToDB add host information into DB
   return hostid(which just added) and []sysadmerror.Sysadmerror
   otherwise return 0  and []sysadmerror.Sysadmerror
*/
func addHostToDB(tx *db.Tx, data *ApiHost)(int,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	// adding host information into DB
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(3030004,"debug","now try to insert host information into host table"))
	insertData := make(db.FieldData,0)
	insertData["hostname"] = data.Hostname 
	insertData["osID"] = data.OsID
	insertData["versionID"] = data.OsVersionID
	insertData["statusID"] = 1
	insertData["agentPassive"] = data.PassiveMode
	insertData["agentPort"] = data.AgentPort
	insertData["receiveCommandUri"] = data.ReceiveCommandUri
		
	_,err := tx.InsertData("host",insertData)
	errs = append(errs,err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
		return 0, errs
	}

	// get the last hostid which is the last added
	selectData := db.SelectData{
		Tb: []string{"host"},
		OutFeilds: []string{"max(hostid) as hostid"},
	}
	entity := tx.Entity
	retData,err := entity.QueryData(&selectData)
	errs = append(errs,err...)
	if retData == nil {
		return 0,errs
	} 
	
	hostid := 0
	line := retData[0]
	if v,ok := line["hostid"]; ok {
		id,e := utils.Interface2Int(v)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(3030005,"error","got hostid error %s",e))
		} else {
			hostid = id
		}
	} else {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3030006,"error","got hostid error"))
	}

	if hostid == 0 {
		return 0,errs
	}

	return hostid,errs
}

/*
   addManageIPToDB add manage IP of host into DB
   return ipID(which just added) and []sysadmerror.Sysadmerror
   otherwise return 0  and []sysadmerror.Sysadmerror
   the IP information(such as devName, mask) of the host should be updated when we got the details of the IP information of the host
*/
func addManageIPToDB(tx *db.Tx, data *ApiHost, hostid int)(int,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303010,"debug","try to add host %s manager IP to DB", data.Hostname))
	insertData := make(db.FieldData,0)

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

	_,err := tx.InsertData("hostIP",insertData)
	errs = append(errs,err...)
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
func addYumHostToDB(tx *db.Tx, data *ApiHost, hostid int)(int,[]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303010,"debug","try to add yum host relations to DB"))
	insertData := make(db.FieldData,0)

	insertData["hostid"] = hostid
	yumIDS := data.YumID
	for _, yID := range yumIDS {
		insertData["yumid"] = yID
		rows,err := tx.InsertData("hostYum",insertData)
		errs = append(errs,err...)
		if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
			return 0, errs
		}
	}

	return 1, errs
}

func addIPToDB(tx *db.Tx, data *ApiHost)(maxIPid int,errs []sysadmerror.Sysadmerror){
	// SSH Public key should be update to host
	if !data.PubkeyUploaded {
		pubkeyFile, e :=  syssetting.GetFilePath(moduleName,"publicKeyFile")
		if pubkeyFile == "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(30302001,"error","%s",e))
			return 0, errs
		}
		pubkeyFile = WorkingData.workingRoot + "/" + pubkeyFile

		privateKeyFile, e := syssetting.GetFilePath(moduleName,"privateKeyFile")
		if privateKeyFile == "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(30301002,"error","%s",e))
			return 0, errs
		}
		privateKeyFile = WorkingData.workingRoot + "/" + privateKeyFile

		copyOK, err := sshcopyid.SshCopyId(data.Ip,data.Port,data.User,data.Password,privateKeyFile, pubkeyFile)
		if !copyOK{
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(30301003,"error","can not copy pubic key to node %s",err))
			return 0, errs
		}

		//hostIPs, err := getHostIPs(data.Ip,data.Port,data.User)
		
	}

	
	return 0, errs
}

func getHostIPs(hostip string, hostPort int, user string) (retIps []HostIP, errs []sysadmerror.Sysadmerror){

	scriptsPath, e := syssetting.GetFilePath(moduleName,"scriptsPath")
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303001,"error","%s",e))
		return retIps,errs
	}

	scriptGetHostIPs,e := syssetting.GetFilePath(moduleName,"scriptGetHostIPs")
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303002,"error","%s",e))
		return retIps,errs
	}

	privateKeyFile, e := syssetting.GetFilePath(moduleName,"privateKeyFile")
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303003,"error","%s",e))
		return retIps,errs
	}
	privateKeyFile = WorkingData.workingRoot + "/" + privateKeyFile

	sftpClient,err := sftp.ConnectSFTP(hostip,user,"",privateKeyFile,hostPort,true,true)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303004,"error","%s",err))
		return retIps,errs
	}

	getHostIpsScripts := WorkingData.workingRoot + "/" + scriptsPath + "/" + scriptGetHostIPs
	dstScript := "/tmp/" + scriptGetHostIPs
	e = sftp.Put(sftpClient,getHostIpsScripts,dstScript)
	sftpClient.Close()
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303005,"error","%s",err))
		return retIps,errs
	}

	addr := hostip + ":" + strconv.Itoa(hostPort)
	sshClient, err := sshclient.DialWithKey(addr, user, privateKeyFile)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303006,"error","%s",err))
		return retIps,errs
	}

	out,err := sshClient.Cmd("chmod +x " + dstScript).SmartOutput()
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303007,"error","%s:%s",err,out))
		return retIps,errs
	}

	ipstr,err := sshClient.Cmd(dstScript).SmartOutput()
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(30303008,"error","get host ip error %s:%s",err,ipstr))
		return retIps,errs
	}
	
	return retIps,errs
}