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
	"reflect"
)

type VerbKind int

// TypeMeta describes an individual object in an API response or request
// with strings representing the type of the object and its API schema version.
// Structures that are versioned or persisted should inline TypeMeta.
type TypeMeta struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	// Cannot be updated.
	// In CamelCase.
	// +optional
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	// APIVersion defines the versioned schema of this representation of an object.
	// Servers should convert recognized schemas to the latest internal value, and
	// may reject unrecognized values.
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
}

// GroupVersion contains the "group" and the "version", which uniquely identifies the API.
type GroupVersion struct {
	Group   string
	Version string
}

// GroupVersionKind unambiguously identifies a kind.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

// ObservedVersionKind identifies a kind with permit verb list.
type ObservedVersionKind struct {
	Gvk GroupVersionKind
	// Verbs is the OR result of all action values that a resource is allowed to perform
	Verbs VerbKind
}

type typePair struct {
	source reflect.Type
	dest   reflect.Type
}

// ConversionFunc converts the object a into the object b, reusing arrays or objects
// or pointers if necessary. It should return an error if the object cannot be converted
// or if some data is invalid. If you do not wish a and b to share fields or nested
// objects, you must copy a before calling this function.
type ConversionFunc func(a, b interface{}) error

type ConversionFuncs struct {
	untyped map[typePair]ConversionFunc
}

// Converter knows how to convert one type to another.
type Converter struct {
	// Map from the conversion pair to a function which can
	// do the conversion.
	ConversionFuncs
}

// Scheme defines methods for serializing and deserializing API objects, a type
// registry for converting group, version, and kind information to and from Go
// schemas, and mappings between Go schemas of different versions. A scheme is the
// foundation for a versioned API and versioned configuration over time.
//
// In a Scheme, a Type is a particular Go struct, a Version is a point-in-time
// identifier for a particular representation of that Type (typically backwards
// compatible), a Kind is the unique name for that Type within the Version, and a
// Group identifies a set of Versions, Kinds, and Types that evolve over time. An
// Unversioned Type is one that is not yet formally bound to a type and is promised
// to be backwards compatible (effectively a "v1" of a Type that does not expect
// to break in the future).
//
// Schemes are not expected to change at runtime and are only threadsafe after
// registration is complete.
type Scheme struct {
	// versionMap allows one to figure out the go type of an object with
	// the given version and name.
	gvkToType map[GroupVersionKind]reflect.Type

	// typeToGroupVersion allows one to find metadata for a given go object.
	// The reflect.Type we index by should *not* be a pointer.
	typeToGVK map[reflect.Type]GroupVersionKind

	// unversionedTypeToGVK is an internal version for an object.
	unversionedTypeToGVK map[reflect.Type]GroupVersionKind

	// unversionedGvkToType  is an internal version for an object.
	unversionedGvkToType map[GroupVersionKind]reflect.Type

	// converter stores all registered conversion functions. It also has
	// default converting behavior.
	converter *Converter

	// observedVersionKinds keeps track of the order we've seen versions during type registration
	observedVersionKinds []ObservedVersionKind
}

type FuncRegistry func(scheme *Scheme) error

type RegistryType struct {
	Gv              GroupVersion
	AddNewTypeFn    FuncRegistry
	ConversionRegFn FuncRegistry
}

type RequestQuery map[string][]string

type ReferenceInfo struct {

	// ID of the resource referencing this resource
	ReferenceId int `json:"referenceId" xml:"referenceId" yaml:"referenceId" db:"referenceId"`

	// Group of the resource referencing this resource
	ReferenceGroup string `json:"referenceGroup" xml:"referenceGroup" yaml:"referenceGroup" db:"referenceGroup"`

	// Kind of the resource referencing this resource
	ReferenceKind string `json:"referenceKind" xml:"referenceKind" yaml:"referenceKind" db:"referenceKind"`

	// Version of the resource referencing this resource
	ReferenceVersion string `json:"referenceVersion" xml:"referenceVersion" yaml:"referenceVersion" db:"referenceVersion"`
}

// ServerError represent an error response to the client
type ServerError struct {
	// Message represent an error detail
	Message string `json:"message" xml:"message" yaml:"message" db:"message"`
}
