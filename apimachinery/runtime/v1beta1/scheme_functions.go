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

package v1beta1

import (
	"fmt"
	"reflect"
	"strings"
)

func (s *Scheme) AddSchemeData(r *RegistryType) error {
	addNewTypeFunc := r.AddNewTypeFn
	return addNewTypeFunc(s)
}

func (s *Scheme) AddKnowTypes(gv GroupVersion, verbs VerbKind, types ...interface{}) error {
	if s.gvkToType == nil {
		s.gvkToType = make(map[GroupVersionKind]reflect.Type)
	}
	gvkToType := s.gvkToType
	if s.typeToGVK == nil {
		s.typeToGVK = make(map[reflect.Type]GroupVersionKind)
	}
	typeToGVK := s.typeToGVK
	if s.unversionedGvkToType == nil {
		s.unversionedGvkToType = make(map[GroupVersionKind]reflect.Type)
	}
	unversionedGvkToType := s.unversionedGvkToType
	if s.unversionedTypeToGVK == nil {
		s.unversionedTypeToGVK = make(map[reflect.Type]GroupVersionKind)
	}
	unversionTypeToGVK := s.unversionedTypeToGVK

	if len(gv.Group) == 0 || len(gv.Version) == 0 {
		return fmt.Errorf("version and group must required on all types: %v ", gv)
	}

	for _, obj := range types {
		t := reflect.TypeOf(obj)
		if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("all types must be pointers point  to structs")

		}
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return fmt.Errorf("all types must be pointers to structs")
		}

		objName := t.Name()
		kind := strings.TrimSpace(strings.ToLower(objName))
		gvk := GroupVersionKind{
			Group:   gv.Group,
			Version: gv.Version,
			Kind:    kind,
		}
		
		if gv.Version == APIVersionInternal {
			if oldT, found := unversionedGvkToType[gvk]; found && oldT != t {
				return fmt.Errorf("Double registration of different types for %v: old=%v.%v, new=%v.%v in scheme ", gvk, oldT.PkgPath(), oldT.Name(), t.PkgPath(), t.Name())
			}
			unversionedGvkToType[gvk] = t

			if _, found := unversionTypeToGVK[t]; !found {
				unversionTypeToGVK[t] = gvk
			}
		} else {
			if oldT, found := gvkToType[gvk]; found && oldT != t {
				return fmt.Errorf("Double registration of different types for %v: old=%v.%v, new=%v.%v in scheme ", gvk, oldT.PkgPath(), oldT.Name(), t.PkgPath(), t.Name())
			}
			gvkToType[gvk] = t

			if _, found := typeToGVK[t]; !found {
				typeToGVK[t] = gvk
			}

			s.addObservedVersion(gvk, verbs)
		}
	}
	s.gvkToType = gvkToType
	s.typeToGVK = typeToGVK
	s.unversionedGvkToType = unversionedGvkToType
	s.unversionedTypeToGVK = unversionTypeToGVK

	return nil
}

func (s *Scheme) addObservedVersion(gvk GroupVersionKind, verbs VerbKind) {
	if len(gvk.Version) == 0 || gvk.Version == APIVersionInternal {
		return
	}

	for _, observedVersionKind := range s.observedVersionKinds {
		if observedVersionKind.Gvk == gvk {
			return
		}
	}

	ovk := ObservedVersionKind{
		Gvk:   gvk,
		Verbs: verbs,
	}

	if s.observedVersionKinds == nil {
		s.observedVersionKinds = make([]ObservedVersionKind, 0)
	}
	s.observedVersionKinds = append(s.observedVersionKinds, ovk)

	return
}

func (s *Scheme) GetObservedVersionKinds() []ObservedVersionKind {
	return s.observedVersionKinds
}

func (s *Scheme) GetUnversionTypeByGVK(gvk GroupVersionKind) reflect.Type {
	unversionedGvkToType := s.unversionedGvkToType
	if unversionType, ok := unversionedGvkToType[gvk]; ok {
		return unversionType
	}

	return nil
}

func (s *Scheme) GetVersionedTypeByGVK(gvk GroupVersionKind) reflect.Type {
	gvkToType := s.gvkToType
	if versionedType, ok := gvkToType[gvk]; ok {
		return versionedType
	}

	return nil
}

func (s *Scheme) GetGroupByKind(k string) string {
	k = strings.TrimSpace(k)
	if k == "" {
		return ""
	}

	obVKs := s.observedVersionKinds
	for _, obVk := range obVKs {
		if obVk.Gvk.Kind == k {
			return obVk.Gvk.Group
		}
	}

	return ""
}

func (s *Scheme) AddConversionFunc(a, b interface{}, fn ConversionFunc) error {
	if s.converter == nil {
		s.converter = &Converter{}
		untyped := make(map[typePair]ConversionFunc)
		s.converter.untyped = untyped
	}
	c := s.converter

	return c.AddConversionFunc(a, b, fn)
}

func Register(gv GroupVersion, registryFn, conversionRegistryFn FuncRegistry) {
	if _, ok := registeredResources[gv]; ok {
		return
	}

	registryType := &RegistryType{}
	registryType.Gv = gv
	registryType.AddNewTypeFn = registryFn
	registryType.ConversionRegFn = conversionRegistryFn
	registeredResources[gv] = registryType

	return
}

func GetRegisteredResources() map[GroupVersion]*RegistryType {
	return registeredResources
}

func GetKindByType(obj interface{}) (string, error) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("all types must be pointers point  to structs")

	}
	objName := t.Name()
	kind := strings.TrimSpace(strings.ToLower(objName))

	return kind, nil
}
