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
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	runtimeBetaV1 "sysadm/apimachinery/runtime/v1beta1"
)

func init() {
	flag.StringVar(&generatedRuntimeCodeFilename, "runtime-code-filename", defaultGeneratedRuntimeCodeFilename, "file name of runtime code which will be generated")
	flag.StringVar(&goHeadFileName, "head-file", defaultGoHeadFileName, "the path of the head file for go files")
	flag.StringVar(&moduleName, "module-name", defaultModuleName, "go module name of the project")
	flag.StringVar(&sysadmRoot, "sysadm-root", "", "root path of the project")

	flag.Parse()
}

var completed = make(chan bool, 1)

func main() {
	e := generateRuntimeFileForApiServer()
	if e != nil {
		fmt.Printf("generate code file for apiServer error: %s\n", e)
		os.Exit(-1)
	}
	destFile := filepath.Join(sysadmRoot, "apiserver", "app", generatedRuntimeCodeFilename)
	fmt.Printf("code for apiServer: %s has be generated\n", destFile)

	return
}

func generateRuntimeFileForApiServer() error {
	goHeadFile := filepath.Join(sysadmRoot, goHeadFileName)
	headContent, e := os.ReadFile(goHeadFile)
	if e != nil {
		return e
	}
	contentStr := string(headContent)
	contentStr = contentStr + "\n\npackage app\n\n"
	contentStr = contentStr + generateApiServerCodeContent()

	destFile := filepath.Join(sysadmRoot, "apiserver", "app", generatedRuntimeCodeFilename)
	return os.WriteFile(destFile, []byte(contentStr), 0600)
}

func generateApiServerCodeContent() string {
	listStr := "import (\n"
	listStr = listStr + "\truntimeBetaV1 \"" + moduleName + "/apimachinery/runtime/v1beta1\"\n"
	for _, rt := range registeredResources {
		listStr = listStr + "\t" + rt.Alias + " \"" + rt.ImportPath + "\"\n"
	}
	listStr = listStr + ")\n\n"

	listStr = listStr + "var registeredResources = []runtimeBetaV1.RegistryType{\n"
	for _, rt := range registeredResources {
		listStr = listStr + "\t{Gv: " + rt.Alias + ".SchemaGroupVersion, AddNewTypeFn: " + rt.Alias + ".TypeRegistryFunc, " +
			"ConversionRegFn: " + rt.Alias + ".ConversionRegistry,ImportPath: \"" + rt.ImportPath + "\", Alias: \"" + rt.Alias + "\"},\n"
	}
	listStr = listStr + "}\n"

	return listStr
}

func generateConversionFiles() error {
	var scheme = &runtimeBetaV1.Scheme{}
	registeredValue := reflect.ValueOf(&registeredResources).Elem().Interface()
	if registeredValue == nil || !reflect.DeepEqual(registeredValue, reflect.Zero(reflect.TypeOf(registeredValue)).Interface()) {
		return fmt.Errorf("type of resource is invalid")
	}

	if len(registeredResources) < 1 {
		return fmt.Errorf("no registered resource")
	}
	for _, r := range registeredResources {
		e := scheme.AddSchemeData(r)
		if e != nil {
			return e
		}
	}
	observedGVKs := scheme.GetObservedVersionKinds()
	aliasMapGVKs := make(map[string][]runtimeBetaV1.GroupVersionKind, 0)
	prepareAliasMapGVKs(observedGVKs, aliasMapGVKs)

	finished := 0
	shouldReturn := false
	ctx, cancelFunc := context.WithCancel(context.Background())

	for alias, gvks := range aliasMapGVKs {
		go generateConversionForTypes(scheme, cancelFunc, alias, gvks)
	}

	for {
		select {
		case <-ctx.Done():
			shouldReturn = true
		case <-completed:
			finished++
			if finished >= len(registeredResources) {
				shouldReturn = true
			}
		}
		if shouldReturn {
			break
		}
	}

	if finished < len(registeredResources) {
		return fmt.Errorf("some conversion files have not generated")
	}

	return nil
}

func generateConversionForTypes(scheme *runtimeBetaV1.Scheme, cancelFunc context.CancelFunc, alias string, gvks []runtimeBetaV1.GroupVersionKind) {

	goHeadFile := filepath.Join(sysadmRoot, goHeadFileName)
	headContent, e := os.ReadFile(goHeadFile)
	if e != nil {
		fmt.Printf("read head content for go files error: %s\n", e)
		cancelFunc()
		return
	}
	contentStr := string(headContent)
	contentStr = contentStr + "\n\n"
	contentStr = contentStr + "// Code generated by runtimegenerator. DO NOT EDIT.\n\n"
	contentStr = contentStr + "package " + gvks[0].Version + "\n\n"
	contentStr = contentStr + "import (\n"
	contentStr = contentStr + "\tunsafe \"unsafe\""
	contentStr = contentStr + "\truntimeBetaV1 \"sysadm/apimachinery/runtime/v1beta1\""
	unversionGV := runtimeBetaV1.GroupVersion{Group: gvks[0].Group, Version: runtimeBetaV1.APIVersionInternal}
	unversionAlias := ""
	unversionImportPath := ""
	for _, rt := range registeredResources {
		if unversionGV == rt.Gv {
			unversionAlias = rt.Alias
			unversionImportPath = rt.ImportPath
			break
		}
	}
	contentStr = contentStr + "\t" + unversionAlias + " \"" + unversionImportPath + "\"\n"
	contentStr = contentStr + ")\n\n"
	contentStr = contentStr + "var ConversionRegistry runtimeBetaV1.FuncRegistry = RegisterConversions\n\n"

	registeredFuns := []conversionRegistryFunc{}
	conversionFuns := []conversionFunc{}
	for _, gvk := range gvks {
		versionType := scheme.GetVersionedTypeByGVK(gvk)
		if versionType == nil {
			fmt.Printf("versioned type %s for group %s and version %s was not found\n", gvk.Kind, gvk.Group, gvk.Version)
			cancelFunc()
			return
		}

		unversionedGVK := runtimeBetaV1.GroupVersionKind{Group: gvk.Group, Version: runtimeBetaV1.APIVersionInternal, Kind: gvk.Kind}
		unversionType := scheme.GetUnversionTypeByGVK(unversionedGVK)
		if unversionType == nil {
			fmt.Printf("unversioned type %s for group %s was not found\n", gvk.Kind, gvk.Group)
			cancelFunc()
			return
		}

		versionTypeKind := getTypeKind(versionType)
		if versionTypeKind == Builtin {
			continue
		}
		if versionTypeKind == Struct {
			structName := versionType.Name()
			conversionFuncName := "Convert_" + gvk.Version + "_" + structName + "_To_" + unversionAlias + "_" + structName

		}

	}
	gv := rt.Gv
	addNewTypeFunc := rt.AddNewTypeFn

}

func prepareAliasMapGVKs(gvks []runtimeBetaV1.GroupVersionKind, mappedGvks map[string][]runtimeBetaV1.GroupVersionKind) {
	for _, gvk := range gvks {
		gv := runtimeBetaV1.GroupVersion{Group: gvk.Group, Version: gvk.Version}
		for _, rt := range registeredResources {
			if gv == rt.Gv {
				alias := rt.Alias
				kind := gvk.Kind
				addedGvk := runtimeBetaV1.GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: kind}
				if _, ok := mappedGvks[alias]; ok {
					mappedGvks[alias] = append(mappedGvks[alias], addedGvk)
				} else {
					mappedGvks[alias] = []runtimeBetaV1.GroupVersionKind{addedGvk}
				}
			}
		}
	}

	return
}

func prepareRegistryFunsForConversion(registeredFuns []conversionRegistryFunc, conversionFuns conversionFunc, versionType, unversionType reflect.Type,
	alias, unversionAlias string) error {

	registryFun := conversionRegistryFunc{}

}

func getTypeKind(t reflect.Type) TypeKind {
	var ret TypeKind = Builtin
	switch t.Kind() {
	case reflect.Struct:
		ret = Struct
	case reflect.Map:
		ret = Map
	case reflect.Slice:
		ret = Slice
	case reflect.Pointer:
		ret = Pointer
	}

	return ret
}
