package v1beta1

// Type returns the canonical type for the given fully-qualified name. Builtin
// types will always be found, even if they haven't been explicitly added to
// the map. If a non-existing type is requested, this will create (a marker for)
// it.
func (u Universe) Type(n Name) *Type {
	return u.Package(n.Package).Type(n.Name)
}

// Function returns the canonical function for the given fully-qualified name.
// If a non-existing function is requested, this will create (a marker for) it.
// If a marker is created, it's the caller's responsibility to finish
// construction of the function by setting Underlying to the correct type.
func (u Universe) Function(n Name) *Type {
	return u.Package(n.Package).Function(n.Name)
}

// Variable returns the canonical variable for the given fully-qualified name.
// If a non-existing variable is requested, this will create (a marker for) it.
// If a marker is created, it's the caller's responsibility to finish
// construction of the variable by setting Underlying to the correct type.
func (u Universe) Variable(n Name) *Type {
	return u.Package(n.Package).Variable(n.Name)
}

// Constant returns the canonical constant for the given fully-qualified name.
// If a non-existing constant is requested, this will create (a marker for) it.
// If a marker is created, it's the caller's responsibility to finish
// construction of the constant by setting Underlying to the correct type.
func (u Universe) Constant(n Name) *Type {
	return u.Package(n.Package).Constant(n.Name)
}

// Package returns the Package for the given path.
// If a non-existing package is requested, this will create (a marker for) it.
// If a marker is created, it's the caller's responsibility to finish
// construction of the package.
func (u Universe) Package(packagePath string) *Package {
	if p, ok := u[packagePath]; ok {
		return p
	}
	p := &Package{
		PkgPath:   packagePath,
		Types:     map[string]*Type{},
		Functions: map[string]*Type{},
		Variables: map[string]*Type{},
		Constants: map[string]*Type{},
		Imports:   map[string]*Package{},
	}
	u[packagePath] = p
	return p
}
