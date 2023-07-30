/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* https://www.sysadm.cn/licenses/apache-2.0.txt
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
*
 */

package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// snow flake format:
// 0 ~ 40      time stamp                     0~7(dc)     0~9(az)    (0~4)seq
// 000000000000000000000000000000000000000000 00000000   00000000000   00000
const (
	maxDatacenterID = 256  //最大数据中心id 0~31
	maxAzID         = 1024 //最大可用区ID 0~255
	maxSequenceID   = 16

	timeLeft       = uint8(23) //时间id向左移位的量
	datacenterLeft = uint8(15) //机器id向左移位的量
	azLeft         = uint8(5)  //对象id向左移位的量

	startime = uint64(1688354362000) //初始毫秒,时间是: Mon Jul  03 11:19:27 CST 2023
)

type Worker struct {
	sync.Mutex
	lastStamp    uint64
	datacenterID uint64 //数据中心id,0~15
	azID         uint64 //可用区id , 0~15
	sequenceID   uint64 //同一毫秒内的序号，最大值为1023
}

func NewWorker(datacenterID, azID uint64) (*Worker, error) {
	if datacenterID < 0 || datacenterID >= maxDatacenterID {
		return nil, fmt.Errorf("datacenter id is not valid. the value of datacenter id should be 0 ~ %d", maxDatacenterID)
	}

	if azID < 0 || azID >= maxAzID {
		return nil, fmt.Errorf("AZ id is not valid. the value of available zone id should be 0 ~ %d", maxAzID)
	}

	return &Worker{
		lastStamp:    0,
		datacenterID: datacenterID,
		azID:         azID,
		sequenceID:   0,
	}, nil
}

func (w *Worker) GetID() (string, error) {
	//多线程互斥
	w.Lock()
	defer w.Unlock()

	mill := uint64(time.Now().UnixMilli())
	if mill < w.lastStamp {
		return "0", errors.New("time is moving backwards,waiting until")
	}

	if mill == w.lastStamp {
		//w.sequenceID = (w.sequenceID + 1) & maxSequenceID
		w.sequenceID = w.sequenceID + 1
		if w.sequenceID >= maxSequenceID {
			w.sequenceID = 0
		}
		//当一个毫秒内分配的id数>sequenceID个时，只能等待到下一毫秒去分配。
		if w.sequenceID == 0 {
			for mill == w.lastStamp {
				mill = uint64(time.Now().UnixMilli())
			}
		}
	} else {
		w.sequenceID = 0
	}

	w.lastStamp = mill

	id := ((w.lastStamp - startime) << timeLeft) |
		(w.datacenterID << datacenterLeft) |
		(w.azID << azLeft) |
		w.sequenceID

	return fmt.Sprintf("%d", id), nil
}
