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

ErrorCode: 500xxx
*/

package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"sysadm/sysadmerror"
)

// Checking a file if is exists.
func CheckFileExists(f string, cmdRunPath string) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	cmdRunPath = strings.TrimSpace(cmdRunPath)
	if cmdRunPath != "" {
		dir, error := filepath.Abs(filepath.Dir(cmdRunPath))
		if error != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(500001, "error", "get the root path of the program error: %s", error))
			return false, errs
		}

		if !filepath.IsAbs(f) {
			tmpDir := filepath.Join(dir, "../")
			f = filepath.Join(tmpDir, f)
		}
	}

	_, err := os.Stat(f)
	if err != nil {
		if os.IsExist(err) {
			return true, errs
		}
		return false, errs
	}

	return true, errs
}

/*
Converting relative path to absolute path of  file(such as socket, accesslog, errorlog) and return the  file path
return "" and error if  file can not opened .
Or return string and nil.
*/
func CheckFileRW(f string, cmdRunPath string, isRmTest bool) (string, error) {
	dir, error := filepath.Abs(filepath.Dir(cmdRunPath))
	if error != nil {
		return "", error
	}

	if !filepath.IsAbs(f) {
		tmpDir := filepath.Join(dir, "../")
		f = filepath.Join(tmpDir, f)
	}

	fp, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		return "", err
	}
	fp.Close()
	if isRmTest {
		_ = os.Remove(f)
	}
	return f, nil
}

/*
Converting relative path to absolute path of file(such as socket, accesslog, errorlog) and return the  file path
return "" and error if  file can not opened .
Or return string and nil.
*/
func CheckFileIsRead(f string, cmdRunPath string) (string, error) {
	var dir string = ""
	var err error
	if strings.TrimSpace(cmdRunPath) != "" {
		dir, err = filepath.Abs(filepath.Dir(cmdRunPath))
		if err != nil {
			return "", err
		}
	}

	if !filepath.IsAbs(f) {
		tmpDir := filepath.Join(dir, "../")
		f = filepath.Join(tmpDir, f)
	}

	fp, err := os.Open(f)
	if err != nil {
		return "", err
	}
	fp.Close()

	return f, nil
}

/*
Converting relative path to absolute path of file and return the  file path
return "" and error if  file can not opened . Or return string and nil.
workingDir should be a absolute path and f is a path relative to workingDir or a absolute path.
*/
func CheckFileIsReadable(f string, workingDir string) (string, error) {
	if strings.TrimSpace(f) == "" {
		return "", fmt.Errorf("file is empty")
	}

	if !filepath.IsAbs(f) {
		if strings.TrimSpace(workingDir) == "" {
			return "", fmt.Errorf("working directory is not valid")
		}

		f = filepath.Join(workingDir, f)
	}

	fi, e := os.Stat(f)
	if e != nil {
		return "", e
	}

	if fi.IsDir() {
		return "", fmt.Errorf("path %s is a directory, not a regular file", f)
	}

	fp, err := os.Open(f)
	if err != nil {
		return "", err
	}
	fp.Close()

	return f, nil
}

// IsFileReadable check a regular file is readable. it returns true if the file can open with RDONLY
// otherwise it returns false.
// f will be join with current directory if f is not a absolute
func IsFileReadable(f string) bool {
	f = strings.TrimSpace(f)
	if f == "" {
		return false
	}
	if !filepath.IsAbs(f) {
		currentDir, e := os.Getwd()
		if e != nil {
			return false
		}
		f = filepath.Join(currentDir, f)
	}

	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	if fi.IsDir() {
		return false
	}

	fp, e := os.Open(f)
	if e != nil {
		return false
	}
	fp.Close()

	return true
}

/*
Converting relative path to absolute path of file and return the  file path
return "" and error if  file can not opened . Or return string and nil.
workingDir should be a absolute path and f is a path relative to workingDir or a absolute path.
*/
func CheckFileWritable(f string, workingDir string, isCreate, isRmCreate bool) (string, error) {
	if strings.TrimSpace(f) == "" {
		return "", fmt.Errorf("file is empty")
	}

	if !filepath.IsAbs(f) {
		if strings.TrimSpace(workingDir) == "" {
			return "", fmt.Errorf("working directory is not valid")
		}

		f = filepath.Join(workingDir, f)
	}

	var fp *os.File = nil
	var err error = nil
	var notExist bool = false
	_, e := os.Stat(f)
	if e != nil && os.IsNotExist(e) {
		notExist = true
	}

	if isCreate {
		fp, err = os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0644)
	} else {
		fp, err = os.OpenFile(f, os.O_RDWR, 0644)
	}
	if err != nil {
		return "", err
	}

	fp.Close()
	if isCreate && isRmCreate && notExist {
		_ = os.Remove(f)
	}

	return f, nil
}

// ReadUploadedFile read all content of upload file .
func ReadUploadedFile(file *multipart.FileHeader) ([]byte, error) {
	var ret []byte

	size := file.Size
	if strings.TrimSpace(file.Filename) == "" || size == 0 {
		return ret, fmt.Errorf("no uploaded file can be read")
	}

	src, err := file.Open()
	if err != nil {
		return ret, err
	}
	defer src.Close()

	buf := make([]byte, size)
	for {
		_, e := src.Read(buf)
		if e == io.EOF {
			ret = append(ret, buf...)
			break
		}

		if e != nil {
			return ret, e
		}
	}
	return ret, nil
}

func ValidPath(path string) bool {
	var pathRegex = regexp.MustCompile(`^[a-zA-Z]:\([w-]+\)*w([w-]+.)*w+$|^(/[w-]+)*(/[w-]+.w+)$`)
	return pathRegex.MatchString(path)
}
