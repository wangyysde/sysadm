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
	"encoding/json"
	"fmt"
	"strconv"
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

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Interface2String(value interface{}) string {
	var key string
	if value == nil {
		return ""
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

// convert interface{} to int. return 0,error if converted with error
// otherwise return int, nil
func Interface2Int(value interface{}) (int, error) {
	var ret int
	var err error
	if value == nil {
		return 0, fmt.Errorf("can not convert empty value to int")

	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		ret = int(ft)
	case float32:
		ft := value.(float32)
		ret = int(ft)
	case int:
		ret = value.(int)
	case uint:
		it := value.(uint)
		ret = int(it)
	case int8:
		it := value.(int8)
		ret = int(it)
	case uint8:
		it := value.(uint8)
		ret = int(it)
	case int16:
		it := value.(int16)
		ret = int(it)
	case uint16:
		it := value.(uint16)
		ret = int(it)
	case int32:
		it := value.(int32)
		ret = int(it)
	case uint32:
		it := value.(uint32)
		ret = int(it)
	case int64:
		it := value.(int64)
		ret = int(it)
	case uint64:
		it := value.(uint64)
		ret = int(it)
	case string:
		it := value.(string)
		ret, err = strconv.Atoi(it)
	case []byte:
		key := string(value.([]byte))
		ret, err = strconv.Atoi(key)
	default:
		newValue, _ := json.Marshal(value)
		key := string(newValue)
		ret, err = strconv.Atoi(key)
	}

	return ret, err
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
