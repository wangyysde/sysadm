/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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

package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"sysadm/httpclient"
	"sysadm/utils"
)

var exitChan chan os.Signal
var shouldExit = false

func startLoop() error {
	exitChan = make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-exitChan
		if s == syscall.SIGHUP || s == syscall.SIGINT || s == syscall.SIGTERM {
			shouldExit = true
		}
		return
	}()
	for {
		e := getCommandLoop()
		if shouldExit {
			return e
		}
		shouldExit = false
		signal.Notify(exitChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			s := <-exitChan
			if s == syscall.SIGHUP || s == syscall.SIGINT || s == syscall.SIGTERM {
				shouldExit = true
			}
			return
		}()
	}

}

func getCommandLoop() error {
	getCommandUrl := buildGetCommandUrl()
	getCommandParames := &httpclient.RequestParams{}
	getCommandParames.Method = http.MethodGet
	getCommandParames.Url = getCommandUrl
	systemUUID, e := utils.GetSystemUUID()
	if e != nil {
		return e
	}
	sendData := map[string]interface{}{
		"systemUUID": systemUUID,
	}
	sendDataJson, e := json.Marshal(sendData)
	if e != nil {
		return e
	}

	for {
		if shouldExit {
			return nil
		}
		if RunData.httpClient == nil {
			e := buildHttpClient()
			if e != nil {
				return e
			}
		}
		log.WithFields(log.Fields{"systemUUID": systemUUID, "url": getCommandUrl}).Debug("getting command from apiServer")
		body, e := httpclient.NewSendRequest(getCommandParames, RunData.httpClient, strings.NewReader(utils.Bytes2str(sendDataJson)))
		if e != nil {
			log.Errorf("%s", e)
		}
		log.WithField("data", body).Debug("got command data")
		e = handleHTTPBody(body)
		if e != nil {
			log.Errorf("%s", e)
		} else {
			log.WithField("data", body).Debug("command has be handled")
		}
		
		time.Sleep(time.Duration(defaultGetCommandInterval) * time.Second)
	}

	return fmt.Errorf("unknow error occurred")
}
