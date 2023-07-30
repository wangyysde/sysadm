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

package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

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
		key = strconv.FormatUint(uint64(it), 10)
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.FormatUint(uint64(it), 10)
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.FormatUint(uint64(it), 10)
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.FormatUint(uint64(it), 10)
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

// convert interface{} to int64. return 0,error if converted with error
// otherwise return int64, nil
func Interface2Int64(value interface{}) (int64, error) {
	var ret int64
	var err error
	if value == nil {
		return int64(0), fmt.Errorf("can not convert empty value to int")

	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		ret = int64(ft)
	case float32:
		ft := value.(float32)
		ret = int64(ft)
	case int:
		i := value.(int)
		ret = int64(i)
	case uint:
		it := value.(uint)
		ret = int64(it)
	case int8:
		it := value.(int8)
		ret = int64(it)
	case uint8:
		it := value.(uint8)
		ret = int64(it)
	case int16:
		it := value.(int16)
		ret = int64(it)
	case uint16:
		it := value.(uint16)
		ret = int64(it)
	case int32:
		it := value.(int32)
		ret = int64(it)
	case uint32:
		it := value.(uint32)
		ret = int64(it)
	case int64:
		ret = value.(int64)
	case uint64:
		it := value.(uint64)
		ret = int64(it)
	case string:
		it := value.(string)
		iti, err := strconv.ParseInt(it, 10, 64)
		if err == nil {
			ret = iti
		}
	case []byte:
		key := string(value.([]byte))
		iti, err := strconv.Atoi(key)
		if err != nil {
			ret = int64(iti)
		}
	default:
		newValue, _ := json.Marshal(value)
		key := string(newValue)
		iti, err := strconv.ParseInt(key, 10, 64)
		if err == nil {
			ret = iti
		}
	}

	return ret, err
}

// convert interface{} to uint64. return 0,error if converted with error
// otherwise return int64, nil
func Interface2Uint64(value interface{}) (uint64, error) {
	var ret uint64
	var err error
	if value == nil {
		return uint64(0), fmt.Errorf("can not convert empty value to uint64")
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		ret = uint64(ft)
	case float32:
		ft := value.(float32)
		ret = uint64(ft)
	case int:
		i := value.(int)
		ret = uint64(i)
	case uint:
		it := value.(uint)
		ret = uint64(it)
	case int8:
		it := value.(int8)
		ret = uint64(it)
	case uint8:
		it := value.(uint8)
		ret = uint64(it)
	case int16:
		it := value.(int16)
		ret = uint64(it)
	case uint16:
		it := value.(uint16)
		ret = uint64(it)
	case int32:
		it := value.(int32)
		ret = uint64(it)
	case uint32:
		it := value.(uint32)
		ret = uint64(it)
	case int64:
		i := value.(int64)
		ret = uint64(i)
	case uint64:
		ret = value.(uint64)
	case string:
		it := value.(string)
		iti, err := strconv.ParseUint(it, 10, 64)
		if err == nil {
			ret = iti
		}
	case []byte:
		iStr := string(value.([]byte))
		iti, err := strconv.ParseUint(iStr, 10, 64)
		if err == nil {
			ret = iti
		}
	default:
		newValue, _ := json.Marshal(value)
		iStr := string(newValue)
		iti, err := strconv.ParseUint(iStr, 10, 64)
		if err == nil {
			ret = iti
		}
	}

	return ret, err
}

// convert interface{} to float64. return 0,error if converted with error
// otherwise return float64, nil
func Interface2Float64(value interface{}) (float64, error) {
	var ret float64
	var err error
	if value == nil {
		return float64(0), fmt.Errorf("can not convert empty value to int")

	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		ret = float64(ft)
	case float32:
		ft := value.(float32)
		ret = float64(ft)
	case int:
		i := value.(int)
		ret = float64(i)
	case uint:
		it := value.(uint)
		ret = float64(it)
	case int8:
		it := value.(int8)
		ret = float64(it)
	case uint8:
		it := value.(uint8)
		ret = float64(it)
	case int16:
		it := value.(int16)
		ret = float64(it)
	case uint16:
		it := value.(uint16)
		ret = float64(it)
	case int32:
		it := value.(int32)
		ret = float64(it)
	case uint32:
		it := value.(uint32)
		ret = float64(it)
	case int64:
		i := value.(int64)
		ret = float64(i)
	case uint64:
		i := value.(uint64)
		ret = float64(i)
	case string:
		it := value.(string)
		iti, err := strconv.ParseFloat(it, 10)
		if err == nil {
			ret = iti
		}
	case []byte:
		key := string(value.([]byte))
		iti, err := strconv.ParseFloat(key, 10)
		if err == nil {
			ret = iti
		}
	default:
		newValue, _ := json.Marshal(value)
		key := string(newValue)
		iti, err := strconv.ParseFloat(key, 10)
		if err == nil {
			ret = iti
		}
	}

	return ret, err
}

// Interface2Bool 将interface类型的数据转换成bool类型数据。当数据类型为数值类型的数据时，Interface2Bool将尝试将其转换成
// 对应类型的数值，并与0进行比较，如果大于0，则返回true，否则返回false. 当数据为字符串和[]bytes类型时，对应的值如果为以下值
// 则返回true否则返回false.key == "1" || key == "Y" || key == "YES" || key == "ON"
// 其它类型的则尝试使用json进行编码，并将编码后的数据转换为字符串，之后与对应字符串值为上述值，则返回true，否则返回false
func Interface2Bool(value interface{}) bool {
	ret := false
	if value == nil {
		return false
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		if ft > 0 {
			ret = true
		}
	case float32:
		ft := value.(float32)
		if ft > 0 {
			ret = true
		}
	case int:
		i := value.(int)
		if i > 0 {
			ret = true
		}
	case uint:
		it := value.(uint)
		if it > 0 {
			ret = true
		}
	case int8:
		it := value.(int8)
		if it > 0 {
			ret = true
		}
	case uint8:
		it := value.(uint8)
		if it > 0 {
			ret = true
		}
	case int16:
		it := value.(int16)
		if it > 0 {
			ret = true
		}
	case uint16:
		it := value.(uint16)
		if it > 0 {
			ret = true
		}
	case int32:
		it := value.(int32)
		if it > 0 {
			ret = true
		}
	case uint32:
		it := value.(uint32)
		if it > 0 {
			ret = true
		}
	case int64:
		it := value.(int64)
		if it > 0 {
			ret = true
		}
	case uint64:
		it := value.(uint64)
		if it > 0 {
			ret = true
		}
	case string:
		str := value.(string)
		str = strings.TrimSpace(strings.ToUpper(str))
		if str == "1" || str == "Y" || str == "YES" || str == "ON" {
			ret = true
		}
	case []byte:
		key := string(value.([]byte))
		key = strings.TrimSpace(strings.ToUpper(key))
		if key == "1" || key == "Y" || key == "YES" || key == "ON" {
			ret = true
		}

	default:
		newValue, _ := json.Marshal(value)
		key := string(newValue)
		key = strings.TrimSpace(strings.ToUpper(key))
		if key == "1" || key == "Y" || key == "YES" || key == "ON" {
			ret = true
		}

	}

	return ret
}
