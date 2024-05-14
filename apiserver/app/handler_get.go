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
	"fmt"
	"github.com/wangyysde/sysadmServer"
	"reflect"
	"strings"
	runtime "sysadm/apimachinery/runtime/v1beta1"
	objects "sysadm/objects/app"
	"sysadm/utils"
)

func getResourceHandler(c *sysadmServer.Context) {
	uri := c.FullPath()
	gvk, e := getGvkBasedPath(uri)
	if e != nil {
		// TODO
	}
	queryParas := c.Request.URL.Query()
	queryData := runtime.RequestQuery(queryParas)

	unversionGvk := runtime.GroupVersionKind{Group: gvk.Group, Version: runtime.APIVersionInternal, Kind: gvk.Kind}
	// respData, e := getResource(unversionGvk, queryData)
	_, e = getResource(unversionGvk, queryData)
	// TODO

	return
}

func listResourceHandler(c *sysadmServer.Context) {
	// TODO

}

func watchResourceHandler(c *sysadmServer.Context) {
	// TODO

}

func getResource(gvk runtime.GroupVersionKind, queryData runtime.RequestQuery) ([]interface{}, error) {
	obj := scheme.GetUnversionTypeByGVK(gvk)
	if obj == nil {
		return nil, fmt.Errorf("resource with GVK %+v was not found", gvk)
	}

	if obj.Kind() == reflect.Pointer {
		obj = obj.Elem()
	}

	var condition map[string]string = nil
	// we call the method of the resource if there is a method named CreateGetCondition on the resource
	// the form of the CreateGetCondition is func(*Receiver)CreateGetCondition(queryData runtime.RequestQuery)(map[string]string,error)
	createGetCondition, ok := obj.MethodByName("CreateGetCondition")
	if ok {
		methodParas := make([]reflect.Value, 2)
		methodParas[0] = reflect.Zero(createGetCondition.Type.In(0))
		methodParas[1] = reflect.ValueOf(queryData)
		values := createGetCondition.Func.Call(methodParas)
		err := values[1].Interface().(error)
		if err != nil {
			return nil, err
		}
		condition = values[0].Interface().(map[string]string)
	} else {
		tmpCondition, e := objects.CreateGetCondition(obj, queryData)
		if e != nil {
			return nil, e
		}
		condition = tmpCondition
	}

	resourceData, e := objects.GetResource(gvk, obj, condition)
	if e != nil {
		return nil, e
	}

	if len(resourceData) < 1 {
		return resourceData, nil
	}

	// try to set relation resource data to every line data
	rr, e := getRelationGvks(obj, gvk)
	if e != nil {
		return nil, e
	}
	if len(rr) > 0 {
		e := setRelationResource(resourceData, rr)
		if e != nil {
			return nil, e
		}
	}

	if isReferenced(obj) {
		e := setReferenceResource(resourceData, gvk)
		if e != nil {
			return nil, e
		}
	}

	return resourceData, nil

}

func getRelationGvks(obj reflect.Type, objGvk runtime.GroupVersionKind) ([]resourceRelation, error) {
	if obj == nil || obj.Kind() != reflect.Struct {
		return nil, fmt.Errorf("the type of parent resource must be struct")
	}

	ret := make([]resourceRelation, 0)
	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		fieldT := field.Type
		var underlyingT reflect.Type = nil
		switch fieldT.Kind() {
		case reflect.Struct:
			underlyingT = fieldT
		case reflect.Array, reflect.Pointer, reflect.Slice:
			underlyingT = fieldT.Elem()
		default:
			continue
		}
		kind := strings.TrimSpace(strings.ToLower(underlyingT.Name()))
		group := scheme.GetGroupByKind(kind)
		if group == "" {
			continue
		}
		gvk := runtime.GroupVersionKind{Group: group, Version: runtime.APIVersionInternal, Kind: kind}
		rr := resourceRelation{parentGvk: objGvk, childGvk: gvk, parentFieldName: field.Name}
		ret = append(ret, rr)
	}

	return ret, nil
}

func setRelationResource(data []interface{}, rr []resourceRelation) error {
	if len(data) < 1 || len(rr) < 1 {
		return nil
	}

	for _, d := range data {
		id, e := objects.GetFeildValueByName(d, runtime.ResourcePkFieldName)
		if e != nil {
			return e
		}
		for _, r := range rr {
			parentGvk := r.parentGvk
			childGvk := r.childGvk
			childType := scheme.GetUnversionTypeByGVK(childGvk)
			if childType == nil {
				return fmt.Errorf("there is no child resource")
			}
			ids := make([]interface{}, 0)
			ids = append(ids, id)
			childData, e := objects.GetRelatedResource(parentGvk, childGvk, childType, ids)
			if e != nil {
				return e
			}
			// set zero value for field instead of nil
			if len(childData) < 1 {
				e := objects.SetFieldValueZeroByName(d, r.parentFieldName)
				if e != nil {
					return e
				}
			}

			// because the child data we just got is slice of pointers, we try to convert it to slice of struct
			newChildData, e := objects.ReplacePointerWithStructForSlice(childData)
			if e != nil {
				return e
			}

			e = objects.SetFeildValueByName(d, newChildData, r.parentFieldName)
			if e != nil {
				return e
			}
		}
	}

	return nil
}

func setReferenceResource(data []interface{}, gvk runtime.GroupVersionKind) error {
	if len(data) < 1 {
		return nil
	}

	referenceTbName := objects.GetResourceReferenceTablesName(gvk)
	if referenceTbName == "" {
		return fmt.Errorf("resource GVK %+v is invalid", gvk)
	}

	for _, d := range data {
		idTmp, e := objects.GetFeildValueByName(d, runtime.ResourcePkFieldName)
		if e != nil {
			return e
		}
		id, e := utils.Interface2Int(idTmp)
		if e != nil {
			return e
		}

		referenceResources, e := objects.GetReferencedResources(referenceTbName, id)
		if e != nil {
			return e
		}
		e = objects.SetFeildValueByName(d, referenceResources, runtime.ReferenceFieldName)
		if e != nil {
			return e
		}
	}

	return nil
}

func isReferenced(obj reflect.Type) bool {
	if obj.Kind() == reflect.Pointer {
		obj = obj.Elem()
	}

	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		if field.Name == runtime.ReferenceFieldName {
			var fieldUnderlyT reflect.Type = nil
			switch field.Type.Kind() {
			case reflect.Struct:
				fieldUnderlyT = field.Type
			case reflect.Pointer, reflect.Slice:
				fieldUnderlyT = field.Type.Elem()
			default:
				continue

			}
			referenceT := reflect.TypeOf(runtime.ReferenceInfo{})
			if fieldUnderlyT == referenceT {
				return true
			}
		}
	}

	return false
}
