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

package server

import(
	"fmt"
	"os"
	"path/filepath"

	"github.com/wangyysde/sysadm/sysadm/config"
	"github.com/wangyysde/sysadmServer"
)

// add handlers for static html page,image,css and js to routers
// return nil if success or return error
func addStaicRoute(r *sysadmServer.Engine,cmdRunPath string) error {
	if r == nil {
		return fmt.Errorf("route is nil.")
	}

	// Index page
	path,err := getStaticPath(config.DefaultPath,cmdRunPath)
	if err != nil {
		return err
	}
	
	// images directory
	path,err = getStaticPath(config.ImagesDir,cmdRunPath)
	if err != nil {
		return err
	}
	r.Static("/images", path)
	
	// css directory
	path,err = getStaticPath(config.CssDir,cmdRunPath)
	if err != nil {
		return err
	}
	r.Static("/css", path)

	// css directory
	path,err = getStaticPath(config.JsDir,cmdRunPath)
	if err != nil {
		return err
	}
	r.Static("/js", path)

	return nil
}

func getStaticPath(path string, cmdRunPath string) (string, error){
	dir ,error := filepath.Abs(filepath.Dir(cmdRunPath))
	if error != nil {
		return "",error
	}

	if path == "" {
		return dir,nil
	}

	if filepath.IsAbs(path) {
		return path,nil
	}

	tmpDir := filepath.Join(dir,"../")
	tmpDir = filepath.Join(tmpDir,config.DefaultHtmlPath)
	tmpDir = filepath.Join(tmpDir,"/")
	tmpDir = filepath.Join(tmpDir,path)

	_,err := os.Stat(tmpDir)
	if err != nil {
		return "",err
	}

	return tmpDir,nil
}