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

// The possible classes of types.
type Kind string

const (
	// Builtin is a primitive, like bool, string, int.
	Builtin Kind = "Builtin"
	Struct  Kind = "Struct"
	Map     Kind = "Map"
	Slice   Kind = "Slice"
	Pointer Kind = "Pointer"

	// Alias is an alias of another type, e.g. in:
	//   type Foo string
	//   type Bar Foo
	// Bar is an alias of Foo.
	//
	// In the real go type system, Foo is a "Named" string; but to simplify
	// generation, this type system will just say that Foo *is* a builtin.
	// We then need "Alias" as a way for us to say that Bar *is* a Foo.
	Alias Kind = "Alias"

	// Interface is any type that could have differing types at run time.
	Interface Kind = "Interface"

	// The remaining types are included for completeness, but are not well
	// supported.
	Array Kind = "Array" // Array is just like slice, but has a fixed length.
	Chan  Kind = "Chan"
	Func  Kind = "Func"

	// DeclarationOf is different from other Kinds; it indicates that instead of
	// representing an actual Type, the type is a declaration of an instance of
	// a type. E.g., a top-level function, variable, or constant. See the
	// comment for Type.Name for more detail.
	DeclarationOf Kind = "DeclarationOf"
	Unknown       Kind = ""
	Unsupported   Kind = "Unsupported"

	// Protobuf is protobuf type.
	Protobuf Kind = "Protobuf"
)

var (
	String = &Type{
		Name: Name{Name: "string"},
		Kind: Builtin,
	}
	Int64 = &Type{
		Name: Name{Name: "int64"},
		Kind: Builtin,
	}
	Int32 = &Type{
		Name: Name{Name: "int32"},
		Kind: Builtin,
	}
	Int16 = &Type{
		Name: Name{Name: "int16"},
		Kind: Builtin,
	}
	Int = &Type{
		Name: Name{Name: "int"},
		Kind: Builtin,
	}
	Uint64 = &Type{
		Name: Name{Name: "uint64"},
		Kind: Builtin,
	}
	Uint32 = &Type{
		Name: Name{Name: "uint32"},
		Kind: Builtin,
	}
	Uint16 = &Type{
		Name: Name{Name: "uint16"},
		Kind: Builtin,
	}
	Uint = &Type{
		Name: Name{Name: "uint"},
		Kind: Builtin,
	}
	Uintptr = &Type{
		Name: Name{Name: "uintptr"},
		Kind: Builtin,
	}
	Float64 = &Type{
		Name: Name{Name: "float64"},
		Kind: Builtin,
	}
	Float32 = &Type{
		Name: Name{Name: "float32"},
		Kind: Builtin,
	}
	Float = &Type{
		Name: Name{Name: "float"},
		Kind: Builtin,
	}
	Bool = &Type{
		Name: Name{Name: "bool"},
		Kind: Builtin,
	}
	Byte = &Type{
		Name: Name{Name: "byte"},
		Kind: Builtin,
	}

	builtins = &Package{
		Types: map[string]*Type{
			"bool":    Bool,
			"string":  String,
			"int":     Int,
			"int64":   Int64,
			"int32":   Int32,
			"int16":   Int16,
			"int8":    Byte,
			"uint":    Uint,
			"uint64":  Uint64,
			"uint32":  Uint32,
			"uint16":  Uint16,
			"uint8":   Byte,
			"uintptr": Uintptr,
			"byte":    Byte,
			"float":   Float,
			"float64": Float64,
			"float32": Float32,
		},
		Imports: map[string]*Package{},
		PkgPath: "",
		Name:    "",
	}
)

var (
	generationKind    = ""
	dirs              = ""
	generatedFileName = ""
)

const (
	defaultGenerationKind    = "conversion"
	defaultGeneratedFileName = "zz_generated.conversion.go"
	conversionTag            = "+sysadm:conversion-gen"
)
