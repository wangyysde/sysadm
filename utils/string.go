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

Errorcode: 120xxx
*/

package utils

import (
	"strings"
	"unsafe"
)

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// FoundStrInSlice check if slice of string sourceStr included subStr string.
// return true if slice of string sourceStr included subStr string,otherwise return false
func FoundStrInSlice(sourceStr []string, subStr string, insensitive bool) bool {
	for _, v := range sourceStr {
		if insensitive {
			v = strings.ToUpper(v)
			subStr = strings.ToUpper(subStr)
		}

		if strings.Compare(v, subStr) == 0 {
			return true
		}
	}

	return false
}

/*
* convert a slice to a string with comma(,) joined.
* retrun a string
 */
func ConvSlice2String(sourceSlice []string) string {
	ret := ""
	for _, v := range sourceSlice {
		if ret == "" {
			ret = v
		} else {
			ret = ret + "," + v
		}
	}

	return ret
}

/*
* return a slice which from sourceStr using comma(,)split
 */
func ConvString2Slice(sourceStr string) []string {

	return strings.Split(sourceStr, ",")

}
