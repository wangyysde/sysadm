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
	
	hostid,err := addHostToDB(tx,&requestData)
	errs = append(errs,err...)
	if hostid == 0 {
		_ = tx.Rollback()
		err := apiutils.SendResponseForErrorMessage(c,30301001,"add host error. check logs for details")
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

	port := data.Port
	if tmpPort,e := utils.CheckPort(port); tmpPort == 0{
		errs = append(errs, e...)
		return fmt.Errorf("SSH port(%d) is not valid",port),errs
	}

	user := strings.TrimSpace(data.User)
	if matched,e := regexp.MatchString(`^[a-zA-Z0-9_.][a-zA-Z0-9_.-]{1,31}$`, user); !matched{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020008,"error","os user (%s) is not valid %s.",user,e))
		return fmt.Errorf("os user (%s) is not valid.",user),errs
	}
	data.User =  user

	pubkeyUploaded := data.PubkeyUploaded
	if !pubkeyUploaded {
		password := strings.TrimSpace(data.Password)
		rePassword := strings.TrimSpace(data.RePassword)
		if password != rePassword {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020009,"error","the password(%s) is not matched twice password(%s).",password,rePassword))
			return fmt.Errorf("the password(%s) is not matched twice password(%s).",password,rePassword), errs
		}

		if matched,e := regexp.MatchString(`^.{1,64}$`, password); !matched{
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(3020010,"error","user password (%s) is not valid %s.",password,e))
			return fmt.Errorf("user password(%s) is not valid.",password),errs
		}

		data.Password = password
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
	insertData["sshPort"] = data.Port
	
	rows,err := tx.InsertData("host",insertData)
	errs = append(errs,err...)
	if rows == 0 {
		return 0, errs
	}

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