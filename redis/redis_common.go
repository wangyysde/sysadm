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
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"

	redis "github.com/go-redis/redis/v8"
	"sysadm/httpclient"
)

// IsValidMode check whether mode is a valid redis mode
func IsValidMode(mode int) bool {
	if mode < 1 || mode > 3 {
		return false
	}

	return true
}

// IsValidMode check whether master is a valid hostname
func IsValidMaster(mode int, master string) bool {
	// master for sentinel mode
	master = strings.TrimSpace(master)
	if mode != RedisModeSentinel {
		return true
	}

	if len(master) < 1 || len(master) > 63 {
		return false
	}

	return true
}

// IsValidAddrs check the address of redis server(or redis sentinel server) is a valid adds
// this is redis server address and port,like as localhost:6379 when mode is 1(single)
// these are addresses and ports of redis servers like as localhost:6379;192.168.1.10:6379;x.x.x.x:6379 when mode is 2(cluster)
// these are addresses and ports of sentinel like as localhost:6379;192.168.1.10:6379;x.x.x.x:6379 when mode is 3(sentinel)
func IsValidAddrs(mode int, addrs string) bool {
	addrs = strings.TrimSpace(addrs)

	// single mode
	if mode == RedisModeSingle {
		if len(addrs) < 5 {
			return false
		}
		addrsSlice := strings.Split(addrs, ":")
		if len(addrsSlice) != 2 {
			return false
		}

		if _, e := strconv.Atoi(addrsSlice[1]); e != nil {
			return false
		}

		return true
	}

	if len(addrs) < 5 {
		return false
	}
	addrSlice := strings.Split(addrs, ";")
	for _, v := range addrSlice {
		if len(v) < 5 {
			return false
		}
		vSlice := strings.Split(v, ":")
		if len(vSlice) != 2 {
			return false
		}

		if _, e := strconv.Atoi(vSlice[1]); e != nil {
			return false
		}
	}

	return true
}

// initating a entity according mode, then open a connection to redis server(s)
// return RedisEntity,nil if successful, otherwise return nil, error
func NewClient(conf ClientConf, workDir string) (RedisEntity, error) {
	var ret RedisEntity = nil
	var tlsClientConf *tls.Config = nil

	if conf.Tls.IsTls {
		tmpTlsClientConf, err := httpclient.BuildTlsClientConfig(conf.Tls.Ca, conf.Tls.Cert, conf.Tls.Key, workDir, conf.Tls.InsecureSkipVerify)
		if err != nil {
			return nil, err
		}

		tlsClientConf = tmpTlsClientConf
	}

	switch {
	// for single mode
	case conf.Mode == RedisModeSingle:
		client := redis.NewClient(&redis.Options{
			Addr:      conf.Addrs,
			Username:  conf.Username,
			Password:  conf.Password,
			DB:        conf.DB,
			TLSConfig: tlsClientConf,
		})

		var entity RedisEntity = RedisSingle{
			Client:     client,
			ClientConf: conf,
		}
		ret = entity
	// for cluster mode
	case conf.Mode == RedisModeCluster:
		addrSlice := strings.Split(conf.Addrs, ";")
		client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:     addrSlice,
			Username:  conf.Username,
			Password:  conf.Password,
			TLSConfig: tlsClientConf,
		})

		var entity RedisEntity = RedisCluster{
			Client:     client,
			ClientConf: conf,
		}

		ret = entity
	// for sentinel mode
	case conf.Mode == RedisModeSentinel:
		addrSlice := strings.Split(conf.Addrs, ";")
		client := redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       conf.Master,
			SentinelAddrs:    addrSlice,
			SentinelUsername: conf.SentinelUsername,
			SentinelPassword: conf.SentinelPassword,
			Username:         conf.Username,
			Password:         conf.Password,
			DB:               conf.DB,
			TLSConfig:        tlsClientConf,
		})

		var entity RedisEntity = RedisSentinel{
			Client:     client,
			ClientConf: conf,
		}

		ret = entity
	// other using single mode
	default:
		addrSlice := strings.Split(conf.Addrs, ";")
		addr := addrSlice[0]
		conf.Mode = 1
		client := redis.NewClient(&redis.Options{
			Addr:      addr,
			Username:  conf.Username,
			Password:  conf.Password,
			DB:        conf.DB,
			TLSConfig: tlsClientConf,
		})

		var entity RedisEntity = RedisSingle{
			Client:     client,
			ClientConf: conf,
		}
		ret = entity
	}

	return ret, nil
}

// Set set the value of a key
func Set(entity RedisEntity, ctx context.Context, key string, value interface{}) error {
	if entity == nil {
		return fmt.Errorf("can not set key value on a nil entity")
	}

	sc := entity.Set(ctx, key, value, 0)
	_, e := sc.Result()
	if e == redis.Nil {
		return fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return fmt.Errorf("can not set key value on the closed client")
	}

	return e
}

// Get get the of value of a key
func Get(entity RedisEntity, ctx context.Context, key string) (string, error) {
	if entity == nil {
		return "", fmt.Errorf("can not get value on  nil entity")
	}

	sc := entity.Get(ctx, key)
	str, e := sc.Result()
	if e == redis.Nil {
		return "", fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return "", fmt.Errorf("can not get key value on the closed client")
	}

	return str, e
}

// Del delete keys
func Del(entity RedisEntity, ctx context.Context, keys ...string) error {
	if entity == nil {
		return fmt.Errorf("can not delete key on  nil entity")
	}

	ic := entity.Del(ctx, keys...)
	_, e := ic.Result()
	if e == redis.Nil {
		return nil
	}

	if e == redis.ErrClosed {
		return fmt.Errorf("can not delete keys on the closed client")
	}

	return e
}

// Exists check whether keys is(are) exist
func Exists(entity RedisEntity, ctx context.Context, keys ...string) (bool, error) {
	if entity == nil {
		return false, fmt.Errorf("can not check existence on  nil entity")
	}

	ic := entity.Exists(ctx, keys...)
	result, e := ic.Result()
	if e == redis.Nil {
		return false, nil
	}

	if e == redis.ErrClosed {
		return false, fmt.Errorf("can not check existence on the closed client")
	}

	if result > 0 {
		return true, nil
	}
	return false, nil
}

// HSet set a hash key and the value
// values's format is one of the folloing:
// "key1", "value1", "key2", "value2"  -- pairs of key,value
// []string{"key1", "value1", "key2", "value2"} --- slice with pairs of key,value
// map[string]interface{}{"key1": "value1", "key2": "value2"} --- map
func HSet(entity RedisEntity, ctx context.Context, key string, values ...interface{}) error {
	if entity == nil {
		return fmt.Errorf("can not set key value on a nil entity")
	}

	ic := entity.HSet(ctx, key, values...)
	result, e := ic.Result()

	if e == redis.Nil {
		return fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return fmt.Errorf("can not set key value on the closed client")
	}

	if result > 0 {
		return nil
	}

	return fmt.Errorf("no key(s) has be set")
}

// HGet get a field's value of the hash
func HGet(entity RedisEntity, ctx context.Context, key, field string) (string, error) {
	if entity == nil {
		return "", fmt.Errorf("can not get value on  nil entity")
	}

	sc := entity.HGet(ctx, key, field)
	str, e := sc.Result()
	if e == redis.Nil {
		return "", fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return "", fmt.Errorf("can not get key value on the closed client")
	}

	return str, e
}

// HDel delete the field(s) in the hash
func HDel(entity RedisEntity, ctx context.Context, key string, fields ...string) error {
	if entity == nil {
		return fmt.Errorf("can not delete hash field on  nil entity")
	}

	ic := entity.HDel(ctx, key, fields...)
	_, e := ic.Result()
	if e == redis.Nil {
		return nil
	}

	if e == redis.ErrClosed {
		return fmt.Errorf("can not delete hash field on the closed client")
	}

	return e
}

// HExists whether the field's name is exists in a hash
func HExists(entity RedisEntity, ctx context.Context, key, field string) (bool, error) {
	if entity == nil {
		return false, fmt.Errorf("can not check existence on nil entity")
	}

	bc := entity.HExists(ctx, key, field)
	result, e := bc.Result()
	if e == redis.Nil {
		return false, nil
	}

	if e == redis.ErrClosed {
		return false, fmt.Errorf("can not check existence on the closed client")
	}

	return result, e
}

// Keys get all keys' name matched pattern
func Keys(entity RedisEntity, ctx context.Context, pattern string) ([]string, error) {
	var ret []string

	if entity == nil {
		return ret, fmt.Errorf("can not get keys on  nil entity")
	}

	ss := entity.Keys(ctx, pattern)
	result, e := ss.Result()
	if e == redis.Nil {
		return ret, fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return ret, fmt.Errorf("can not get key value on the closed client")
	}

	return result, e
}

// HKeys get all field's name in specified hash
func HKeys(entity RedisEntity, ctx context.Context, key string) ([]string, error) {
	var ret []string

	if entity == nil {
		return ret, fmt.Errorf("can not get keys on  nil entity")
	}

	key = strings.TrimSpace(key)
	if len(key) < 1 {
		return ret, fmt.Errorf("no key has be specified")
	}

	ss := entity.HKeys(ctx, key)
	result, e := ss.Result()
	if e == redis.Nil {
		return ret, fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return ret, fmt.Errorf("can not get key value on the closed client")
	}

	return result, e

}

// HGetAll get all field name and them value in a hash
func HGetAll(entity RedisEntity, ctx context.Context, key string) (map[string]string, error) {
	ret := make(map[string]string, 0)

	ss := entity.HGetAll(ctx, key)
	result, e := ss.Result()
	if e == redis.Nil {
		return ret, fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return ret, fmt.Errorf("can not get key value on the closed client")
	}

	return result, e
}

// prepend one or multiple elements to a list
func LPush(entity RedisEntity, ctx context.Context,key string, values ...interface{}) error {
	if entity == nil {
		return fmt.Errorf("can not delete hash field on  nil entity")
	}

	ic := entity.LPush(ctx,key,values...) 
	_, e := ic.Result()
	if e == redis.Nil {
		return fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return fmt.Errorf("can not push elementes on the closed client")
	}

	return e
}

// Append one or multiple elements to a list
func RPush(entity RedisEntity, ctx context.Context,key string, values ...interface{}) error {
	if entity == nil {
		return fmt.Errorf("can not push elementes on  nil entity")
	}

	ic := entity.RPush(ctx,key,values...) 
	_, e := ic.Result()
	if e == redis.Nil {
		return fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return fmt.Errorf("can not push elementes on the closed client")
	}

	return e
}

// Remove and get the first elements in a list
func LPop(entity RedisEntity, ctx context.Context, key string) (string, error) {
	if entity == nil {
		return "", fmt.Errorf("can not get value on  nil entity")
	}

	sc := entity.LPop(ctx,key)
	str, e := sc.Result()
	if e == redis.Nil {
		return "", fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return "", fmt.Errorf("can not get key value on the closed client")
	}

	return str, e
}

// Remove and get the last elements in a list
func RPop(entity RedisEntity, ctx context.Context, key string) (string, error) {
	if entity == nil {
		return "", fmt.Errorf("can not get value on  nil entity")
	}

	sc := entity.RPop(ctx,key)
	str, e := sc.Result()
	if e == redis.Nil {
		return "", fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return "", fmt.Errorf("can not get key value on the closed client")
	}

	return str, e
}

// Get the length of a list
func LLen(entity RedisEntity, ctx context.Context,key string)(int, error) {
	if entity == nil {
		return -1, fmt.Errorf("can not get length of a list on nil entity")
	}

	ic := entity.LLen(ctx,key) 
	len, e := ic.Result()
	if e == redis.Nil {
		return -1, fmt.Errorf("the key does not exist")
	}

	if e == redis.ErrClosed {
		return -1, fmt.Errorf("can not get length of a list  on the closed client")
	}

	return int(len), e
}