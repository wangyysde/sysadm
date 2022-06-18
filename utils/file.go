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

ErrorCode: 500xxx
 */

package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/wangyysde/sysadm/sysadmerror"
)

// Checking a file if is exists.
func CheckFileExists(f string,cmdRunPath string ) (bool,[]sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	cmdRunPath = strings.TrimSpace(cmdRunPath)
	if cmdRunPath != "" {
		dir ,error := filepath.Abs(filepath.Dir(cmdRunPath))
		if error != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(500001,"error","get the root path of the program error: %s",error))
			return false,errs
		}

		if ! filepath.IsAbs(f) {
			tmpDir := filepath.Join(dir,"../")
			f = filepath.Join(tmpDir,f)
		}
	}
	
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false,errs
	}

	return true,errs
}

/*
  Converting relative path to absolute path of  file(such as socket, accesslog, errorlog) and return the  file path
  return "" and error if  file can not opened .
  Or return string and nil.
*/
func CheckFileRW(f string,cmdRunPath string, isRmTest bool)(string,error){
	dir ,error := filepath.Abs(filepath.Dir(cmdRunPath))
	if error != nil {
		return "",error
	}

	if ! filepath.IsAbs(f) {
		tmpDir := filepath.Join(dir,"../")
		f = filepath.Join(tmpDir,f)
	}

	fp, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR|os.O_APPEND,os.ModeAppend|os.ModePerm)
	if err != nil {
		return "",err
	}
	fp.Close()
	if isRmTest {
		_ = os.Remove(f)
	}
	return f,nil
}

/*
  Converting relative path to absolute path of file(such as socket, accesslog, errorlog) and return the  file path
  return "" and error if  file can not opened .
  Or return string and nil.
*/
func CheckFileIsRead(f string,cmdRunPath string)(string,error){
	dir ,error := filepath.Abs(filepath.Dir(cmdRunPath))
	if error != nil {
		return "",error
	}

	if ! filepath.IsAbs(f) {
		tmpDir := filepath.Join(dir,"../")
		f = filepath.Join(tmpDir,f)
	}

	fp, err := os.Open(f)
	if err != nil {
		return "",err
	}
	fp.Close()

	return f,nil
}

