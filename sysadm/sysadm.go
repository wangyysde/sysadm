/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2021 Bzhy Network. All rights reserved.
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

package main

import (
//	"fmt"
//	"net/http"

	"github.com/wangyysde/sysadm/sysadm/cmd"
//	"github.com/wangyysde/sysadmServer"
)

func main(){

    cmd.Execute()

    /*
    r := sysadmServer.New()

       // Define handlers
    r.GET("/", func(c *sysadmServer.Context) {
        c.String(http.StatusOK, "Hello World!")
    })  
    r.GET("/ping", func(c *sysadmServer.Context) {
        c.String(http.StatusOK, "echo ping message")
    })  

    // Listen and serve on defined port
    sysadmServer.Log(fmt.Sprintf("Listening on port %s", "8080"),"info")
    r.Run(":8080")

	*/
}