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
	"fmt"
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type RedisSentinel struct {
	Client *redis.Client
	ClientConf ClientConf
}

func (r RedisSentinel) Close() error {
	c := r.Client
	if c != nil {
		e := r.Close()
		return e
	}

	return fmt.Errorf("redis connection has be closed")
}

func (r RedisSentinel) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd{
	return r.Client.Set(ctx,key,value,expiration)
}

func (r RedisSentinel) Get(ctx context.Context, key string) *redis.StringCmd{
	return r.Client.Get(ctx,key)
}

func (r RedisSentinel) Del(ctx context.Context, keys ...string) *redis.IntCmd{
	return r.Client.Del(ctx,keys...)
}

func (r RedisSentinel) Exists(ctx context.Context, keys ...string) *redis.IntCmd{
	return r.Client.Exists(ctx,keys...)
}

func (r RedisSentinel) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd{
	return r.Client.HSet(ctx,key,values...)
}

func (r RedisSentinel) HGet(ctx context.Context, key, field string) *redis.StringCmd{
	return r.Client.HGet(ctx,key,field)
}

func (r  RedisSentinel) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd{
	return r.Client.HGetAll(ctx,key)
}

func (r RedisSentinel) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd{
	return r.Client.HDel(ctx,key,fields...)
}

func (r RedisSentinel) HExists(ctx context.Context, key, field string) *redis.BoolCmd{
	return r.Client.HExists(ctx,key,field)
}

func (r RedisSentinel) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd{
	return r.Client.Keys(ctx,pattern)
}

func (r RedisSentinel) HKeys(ctx context.Context, key string) *redis.StringSliceCmd{
	return r.Client.Keys(ctx,key)
}

func (r RedisSentinel) LPush(ctx context.Context,key string, values ...interface{}) *redis.IntCmd{
	return r.Client.LPush(ctx,key,values...) 
}

func (r RedisSentinel) RPush(ctx context.Context,key string, values ...interface{}) *redis.IntCmd{
	return r.Client.RPush(ctx,key,values...) 
}

func (r RedisSentinel) LPop(ctx context.Context, key string) *redis.StringCmd{
	return r.Client.LPop(ctx,key)
}

func (r RedisSentinel) RPop(ctx context.Context, key string) *redis.StringCmd{
	return r.Client.RPop(ctx,key)
}

func (r RedisSentinel) LLen(ctx context.Context, key string) *redis.IntCmd{
	return r.Client.LLen(ctx,key)
}