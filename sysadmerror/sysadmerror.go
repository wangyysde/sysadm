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

package sysadmerror

import (
	"fmt"
	"strings"
)

 var Levels = []string{
    "trace",
    "debug",
    "info",
    "warning",
    "error",
    "fatal",
    "panic",
}

type Sysadmerror struct {
	ErrorNo int 
	ErrorLevel int 
	ErrorMsg string 
}

// Get the index of error level.
// Return int of index if found otherwise return 1 which is the index of "debug"
func GetLevelNum(level string) int { 

	for key := range Levels {
		if strings.ToLower(Levels[key]) == strings.ToLower(level) {
			return key
		}
	}
	return 1
}

// Get the string of error level 
// Return the string of the error level if the level was found, otherwise return "debug"
func GetLevelString(level int) string{
	if level < 0 || level > 6 {
		return "debug"
	}

	return Levels[level]
}

// Return the error no of err if err is not nil, otherwise return 0
func GetErrorNo(err Sysadmerror) int {
	if err == (Sysadmerror{}) {
		return 0
	}
	return err.ErrorNo
}

// Return the error level of err if err is not nil, otherwise return 1("debug")
func GetErrorLevelNum(err Sysadmerror) int {
	if err == (Sysadmerror{}) {
		return 1
	}

	return err.ErrorLevel
}

// Return the error level of err if err is not nil, otherwise return 1("debug")
func GetErrorLevelString(err Sysadmerror) string {
	if err == (Sysadmerror{}) {
		return "debug"
	}

	return GetLevelString(err.ErrorLevel)
}

// NewErrorWithNumLevel create a new instance of Sysadmerror with errno,errLevel(int) and errMsg
// return Sysadmerror
func NewErrorWithNumLevel(errno int, errLevel int, errMsg string, args ...interface{}) Sysadmerror {
	errmsg := fmt.Sprintf(errMsg, args...)
	if errLevel < 0 || errLevel > 6 {
		errLevel = 1
	}

	err := Sysadmerror{
		ErrorNo: errno,
		ErrorLevel: errLevel,
		ErrorMsg: errmsg,
	}

	return err
}

// NewErrorWithStringLevel create a new instance of Sysadmerror with errno,errLevel(string) and errMsg
// return Sysadmerror. ErrorLevel will be set to 1("debug") if errLevel was not found in Levels
func NewErrorWithStringLevel(errno int, errLevel string, errMsg string, args ...interface{}) Sysadmerror {
	errmsg := fmt.Sprintf(errMsg, args...)
	level := GetLevelNum(errLevel)
	
	err := Sysadmerror{
		ErrorNo: errno,
		ErrorLevel: level,
		ErrorMsg: errmsg,
	}

	return err
}

// Get the maxLevels in []Sysadmerror
// return -1 if the length of []Sysadmerror less 1
// otherwise return the maxLevels in the []Sysadmerror
func GetMaxLevel(errs []Sysadmerror) int {
	if len(errs) < 1 {
		return -1
	}

	maxLevel := 0
	for _,v := range errs {
		l := v.ErrorLevel
		if l > maxLevel {
			maxLevel = l
		}
	}

	return maxLevel
}