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

package main

import (
	"os"
	"strings"
	v1beta12 "sysadm/bakCode/generator/v1beta1"
	"sysadm/generator/v1beta1"
)

func generateConversionFiles() error {
	copyRight, e := os.ReadFile(copyrightFile)
	if e != nil {
		return e
	}
	dirSlice := strings.Split(dirs, ",")
	u := make(v1beta12.Universe)
	pkgData := make(map[string]v1beta12.Package)
	for _, d := range dirSlice {
		e := v1beta12.GetPkgData(d, u, pkgData)
		if e != nil {
			return e
		}
	}

	for _, p := range pkgData {
		content := string(copyRight)
		content = content + "\n"
		content = content + "package " + p.Name + "\n\n"
		//	for key, t := range p.Types {
		//
		//	}
	}

	return nil
}

func generateConversionFile(pkg v1beta1.Package, copyright string) error {
	content := copyright
	content = content + "\n"
	content = content + "// Code generated by conversion-gen. DO NOT EDIT.\n\n"
	content = content + "package " + pkg.Name + "\n\n"
	internalPkgPaths := strings.Split(pkg.InternalTypePath, "/")
	internalPkgName := internalPkgPaths[len(internalPkgPaths)-1]
	content = content + "import (\n"
	content = content + "\t runtime \"sysadm/apimachinery/runtime/v1beta1\"\n"
	content = content + "\t " + internalPkgName + " \"" + pkg.InternalTypePath + "\"\n"
	content = content + ")\n\n"

	content = content + "var ConversionRegistry runtime.FuncRegistry = RegisterConversions\n\n"
	content = content + "// RegisterConversions adds conversion functions to the given scheme.\n"
	content = content + "// Public to allow building arbitrary schemes.\n"
	content = content + "func RegisterConversions(s *runtime.Scheme) error {\n"
	for n, t := range pkg.Types {
		if t.Kind != v1beta12.Struct {
			continue
		}
		versionedName, internalName := buildConversionFunctionName(n, pkg.Name, internalPkgName)
		pkgFuncs := pkg.Functions
		if _, ok := pkgFuncs[versionedName]; ok {
			content = content + "\t if err := s.AddConversionFunc((*" + n + ")(nil), (*" + internalPkgName + "." + n + ")(nil), func(a,b interface{}) error {\n"
			content = content + "\t\t" + "return " + versionedName + "(a.(*" + n + "), b.(*" + internalPkgName + "." + n + ")\n"
			content = content + "\t}); err != nil { \n"
			content = content + "\t\t return err \n"
			content = content + "}\n"
		} else {
			content = content + "\t if err := s.AddGeneratedConversionFunc((*" + n + ")(nil), (*" + internalPkgName + "." + n + ")(nil), func(a,b interface{}) error {\n"
			content = content + "\t\t" + "return " + versionedName + "(a.(*" + n + "), b.(*" + internalPkgName + "." + n + ")\n"
			content = content + "\t}); err != nil { \n"
			content = content + "\t\t return err \n"
			content = content + "}\n"
		}
		if _, ok := pkgFuncs[internalName]; ok {
			content = content + "\t if err := s.AddConversionFunc((*" + internalPkgName + "." + n + ")(nil), (*" + n + ")(nil), func(a,b interface{}) error {\n"
			content = content + "\t\t" + "return " + internalName + "(a.(*" + internalPkgName + "." + n + "), b.(*" + n + ")\n"
			content = content + "\t}); err != nil { \n"
			content = content + "\t\t return err \n"
			content = content + "}\n"
		} else {
			content = content + "\t if err := s.AddGeneratedConversionFunc((*" + internalPkgName + "." + n + ")(nil), (*" + n + ")(nil), func(a,b interface{}) error {\n"
			content = content + "\t\t" + "return " + internalName + "(a.(*" + internalPkgName + "." + n + "), b.(*" + n + ")\n"
			content = content + "\t}); err != nil { \n"
			content = content + "\t\t return err \n"
			content = content + "}\n"
		}

		members := t.Members

	}

}

func buildConversionFunctionName(objName, versioned, internalTypeName string) (string, string) {
	versionedToInternal := "Convert_" + versioned + "_" + objName + "_To_" + internalTypeName + "_" + objName
	internalToVersioned := "Convert_" + internalTypeName + "_" + objName + "_To_" + versioned + "_" + objName

	return versionedToInternal, internalToVersioned
}