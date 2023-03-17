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
	"time"

	redis "github.com/go-redis/redis/v8"

	"github.com/wangyysde/sysadm/config"
)

type RedisEntity interface {
	Close() error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	HExists(ctx context.Context, key, field string) *redis.BoolCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
	HKeys(ctx context.Context, key string) *redis.StringSliceCmd
	LPush(ctx context.Context,key string, values ...interface{}) *redis.IntCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LPop(ctx context.Context, key string) *redis.StringCmd
	RPop(ctx context.Context, key string) *redis.StringCmd
	LLen(ctx context.Context, key string) *redis.IntCmd
 }

// connection parameters for client.
type ClientConf struct {
	// connection mode 1 for single server; 2 for cluster; 3 for sentinel mode
    Mode int `form:"mode" json:"mode" yaml:"mode" xml:"mode"`

	// master server name. the value of this field is empty when mode are 1 and 2
    Master string `form:"master" json:"master" yaml:"master" xml:"master"`

	// a string join with semicolon for the addresses of server
	// this is redis server address and port,like as localhost:6379 when mode is 1
	// these are addresses and ports of redis servers like as localhost:6379;192.168.1.10:6379;x.x.x.x:6379 when mode is 2
	// these are addresses and ports of sentinel like as localhost:6379;192.168.1.10:6379;x.x.x.x:6379 when mode is 3
    Addrs string  `form:"addrs" json:"addrs" yaml:"addrs" xml:"addrs"`

	//redis server username
    Username string `form:"username" json:"username" yaml:"username" xml:"username"`

	// redis server password
    Password string `form:"password" json:"password" yaml:"password" xml:"password"`

	// sentinel username
    SentinelUsername string `form:"sentinelUsername" json:"sentinelUsername" yaml:"sentinelUsername" xml:"sentinelUsername"`

	// sentinel password
    SentinelPassword string `form:"sentinelPassword" json:"sentinelPassword" yaml:"sentinelPassword" xml:"sentinelPassword"`

	// db 
	DB int `form:"db" json:"db" yaml:"db" xml:"db"`

	// tls parameters for agent when agent running as daemon.
	Tls config.Tls `form:"tls" json:"tls" yaml:"tls" xml:"tls"`
}
