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
*
 */

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var (
	redis_host     string      = "192.53.117.73:6379"
	redis_user     string      = ""
	redis_password string      = ""
	entity         RedisEntity = nil
)

var ctx = context.Background()

func connectRedis(t *testing.T) (RedisEntity, error) {
	conf := ClientConf{
		Mode:     1,
		Master:   "",
		Addrs:    redis_host,
		Username: redis_user,
		Password: redis_password,
	}

	path, e := os.Getwd()
	if e != nil {
		t.Errorf("get rooted path error %s", e)
		return nil, e
	}

	return NewClient(conf, path)
}

func TestNewClient(t *testing.T) {

	tmpEntity, e := connectRedis(t)
	if e != nil {
		t.Errorf("can not connect to redis server %s", e)
		os.Exit(1)
	}

	entity = tmpEntity
	t.Log("connect to redis server successful")
}

func TestSet(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	e := Set(entity, ctx, "/commandStatus/202303131527", "202303131527")
	if e != nil {
		t.Errorf("can not set key value %s", e)
		return
	}

	t.Log("the key value has be set successful")
}

func TestGet(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	value, e := Get(entity, ctx, "/commandStatus/202303131527")
	if e != nil {
		t.Errorf("can not get key value %s", e)
		return
	}

	fmt.Printf("we have got the value %s of key %s \n", "/commandStatus/202303131527", value)
}

func TestExist(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	b, e := Exists(entity, ctx, "/commandStatus/202303131527")
	if e != nil {
		t.Errorf("check the key /commandStatus/202303131527 exist error %s", e)
		return
	}

	if b {
		t.Log("check the key /commandStatus/202303131527 exist successful")
		return
	}

	fmt.Print("check the key /commandStatus/202303131527 exist is not correct\n")

	b, e = Exists(entity, ctx, "/commandStatus/202303131528")
	if e != nil {
		t.Errorf("check the key /commandStatus/202303131528 exist error %s", e)
		return
	}

	if !b {
		t.Log("check the key /commandStatus/202303131528 exist successful")
		return
	}

	fmt.Print("check the key /commandStatus/202303131528 exist is not correct\n")
}

func TestHSet(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	s := make(map[string]interface{}, 0)
	s["Name"] = "testName"
	s["Age"] = "18"
	s["score"] = 90
	sec := make(map[string]interface{}, 0)
	sec["field1"] = "1111111"
	third := make(map[string]interface{}, 0)
	third["f1"] = "AAAAAA"
	sec["field2"] = third

	var fieldStr string = ""
	fieldBytes, err := json.Marshal(&sec)
	if err != nil {
		fieldStr = ""
	} else {
		fieldStr = fmt.Sprintf("%s", fieldBytes)
	}
	s["test"] = fieldStr
	e := HSet(entity, ctx, "/commandStatus/202303131529", s)
	if e != nil {
		t.Errorf("set the value of key /commandStatus/202303131529 error %s", e)
		return
	}

	fmt.Print("set the value of key /commandStatus/202303131529 ok\n")
}

func TestHGet(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	name, e := HGet(entity, ctx, "/commandStatus/202303131529", "Name")
	if e != nil {
		t.Errorf("get the value of field of a hash error %s", e)
		return
	}

	fmt.Printf("we have got the value of field of key /commandStatus/202303131529 %s", name)
}

func TestHGetAll(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	result, e := HGetAll(entity, ctx, "/commandStatus/202303131529")
	if e != nil {
		t.Errorf("get the value of hash key error %s", e)
		return
	}

	fmt.Printf("we have got the value of hash %+v", result)
}

func TestHExists(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	b, e := HExists(entity, ctx, "/commandStatus/202303131529", "name")
	if e != nil {
		t.Errorf("check field of key /commandStatus/202303131529 exist error %s", e)
		return
	}

	if b {
		t.Log("check field of key /commandStatus/202303131529 exist successful")
		return
	}

	fmt.Print("check field of  key /commandStatus/202303131529 exist is not correct\n")

	b, e = HExists(entity, ctx, "/commandStatus/202303131529", "sex")
	if e != nil {
		t.Errorf("check field of of key /commandStatus/202303131529 exist error %s", e)
		return
	}

	if !b {
		t.Log("check field of  key /commandStatus/202303131529 exist successful")
		return
	}

	fmt.Print("check field of key /commandStatus/202303131529 exist is not correct\n")
}

func TestKeys(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	k, e := Keys(entity, ctx, "/commandStatus/*")
	if e != nil {
		t.Errorf("get all keys error %s", e)
		return
	}

	for i, v := range k {
		fmt.Printf("No: %d key: %s \n", i, v)
	}
}

func TestDel(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	e := Del(entity, ctx, "/commandStatus/202303131529")
	if e != nil {
		t.Errorf("delete key error %s", e)
		return
	}

	fmt.Print("delete key sucessful")
}

func TestLPush(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	s := make(map[string]interface{}, 0)
	s["Name"] = "testName"
	s["Age"] = "18"
	s["score"] = 90
	sec := make(map[string]interface{}, 0)
	sec["field1"] = "1111111"
	third := make(map[string]interface{}, 0)
	third["f1"] = "AAAAAA"
	sec["field2"] = third

	var fieldStr string = ""
	fieldBytes, err := json.Marshal(&sec)
	if err != nil {
		fieldStr = ""
	} else {
		fieldStr = string(fieldBytes)
	}
	e := LPush(entity, ctx, "/commandLogs/20230317", fieldStr)
	if e != nil {
		t.Errorf("can not push key value %s", e)
		return
	}

	t.Log("json data key has be pushed into list")

	e = LPush(entity, ctx, "/commandLogs/20230317", "the second element's value")
	if e != nil {
		t.Errorf("can not push second key %s", e)
		return
	}

	t.Log("the second key has be pushed into list")
}

func TestLLen(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	len, e := LLen(entity, ctx, "/commandLogs/20230317")
	if e != nil {
		t.Errorf("get the length of a list error %s", e)
		return
	}

	t.Logf("the length of the list is %d\n", len)
}

func TestLPop(t *testing.T) {
	if entity == nil {
		tmpEntiy, e := connectRedis(t)
		if e != nil {
			t.Log("can not connect to redis server")
			os.Exit(2)
		}
		entity = tmpEntiy
	}

	s, e := RPop(entity, ctx, "/commandLogs/20230317")
	if e != nil {
		t.Errorf("pop element from a list error %s", e)
		return
	}
	t.Logf("popped element from the list is %s\n", s)
}
