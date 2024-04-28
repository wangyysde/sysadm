package v1beta1

// Package holds package-level information.
// Fields are public, as everything in this package, to enable consumption by
// code.
type Package struct {
	// Canonical name of this package-- its path.
	// such as sysadm/generator/v1beta1, github.com/kubernetes
	PkgPath string

	// The location this package was loaded from the disk
	// the value of this field is a directory point to the location of package source files
	SourcePath string

	// Short name of this package; the name that appears in the
	// 'package x' line.
	Name string

	// InternalTypePath is a package path of internal type for a versioned package
	InternalTypePath string

	// Types within this package, indexed by their name (*not* including
	// package name).
	Types map[string]*Type

	// Functions within this package, indexed by their name (*not* including
	// package name).
	Functions map[string]*Type

	// Global variables within this package, indexed by their name (*not* including
	// package name).
	Variables map[string]*Type

	// Global constants within this package, indexed by their name (*not* including
	// package name).
	Constants map[string]*Type

	// Packages imported by this package, indexed by (canonicalized)
	// package path.
	Imports map[string]*Package
}

// Type represents a subset of possible go types.
type Type struct {
	// There are two general categories of types, those explicitly named
	// and those anonymous. Named ones will have a non-empty package in the
	// name field.
	//
	// An exception: If Kind == DeclarationOf, then this name is the name of a
	// top-level function, variable, or const, and the type can be found in Underlying.
	// We do this to allow the naming system to work against these objects, even
	// though they aren't strictly speaking types.
	Name Name

	// The general kind of this type.
	Kind Kind

	// If there are comment lines immediately before the type definition,
	// they will be recorded here.
	CommentLines []string

	// If there are comment lines preceding the `CommentLines`, they will be
	// recorded here. There are two cases:
	// ---
	// SecondClosestCommentLines
	// a blank line
	// CommentLines
	// type definition
	// ---
	//
	// or
	// ---
	// SecondClosestCommentLines
	// a blank line
	// type definition
	// ---
	SecondClosestCommentLines []string

	// If Kind == Struct
	Members map[Name]*Type

	// If Kind == Map, Slice, Pointer, or Chan
	Elem *Type

	// If Kind == Map, this is the map's key type.
	Key *Type

	// If Kind == Alias, this is the underlying type.
	// If Kind == DeclarationOf, this is the type of the declaration.
	Underlying *Type

	// If Kind == Interface, this is the set of all required functions.
	// Otherwise, if this is a named type, this is the list of methods that
	// type has. (All elements will have Kind=="Func")
	Methods map[Name]*Type

	// If Kind == func, this is the signature of the function.
	Signature *Signature

	Funcs map[Name]*Type

	// ConstValue contains a stringified constant value if
	// Kind == DeclarationOf and this is a constant value
	// declaration. For string constants, this field contains
	// the entire, un-quoted value. For other types, it contains
	// a human-readable literal.
	ConstValue *string
}

// A type name may have a package qualifier.
type Name struct {
	// Empty if embedded or builtin. This is the package path unless Path is specified.
	Package string
	// The type name.
	Name string
	// An optional location of the type definition for languages that can have disjoint
	// packages and paths.
	Path string
}

// Signature is a function's signature.
type Signature struct {
	// TODO: store the parameter names, not just types.

	// If a method of some type, this is the type it's a member of.
	Receiver   *Type
	Parameters []*Type
	Results    []*Type

	// True if the last in parameter is of the form ...T.
	Variadic bool

	// If there are comment lines immediately before this
	// signature/method/function declaration, they will be recorded here.
	CommentLines []string
}

// Universe is a map of all packages. The key is the package name, but you
// should use Package(), Type(), Function(), or Variable() instead of direct
// access.
type Universe map[string]*Package
