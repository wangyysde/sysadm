package v1beta1

import (
	"fmt"
	"go/constant"
	"go/token"
	tc "go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

func WalkType(u Universe, useName *Name, in tc.Type) *Type {
	name := TcNameToName(in.String())
	if useName != nil {
		name = *useName
	}
	switch t := in.(type) {
	case *tc.Struct:
		out := u.Type(name)
		if out.Kind != Unknown {
			return out
		}
		out = &Type{}
		out.Name = name
		out.Kind = Struct
		members := make(map[Name]*Type)
		for i := 0; i < t.NumFields(); i++ {
			f := t.Field(i)
			// we are pick up exported Field only
			if !f.Exported() {
				continue
			}
			ft := f.Type()
			fName := Name{Path: name.Path, Package: name.Package, Name: f.Name()}
			fOut := WalkType(u, nil, ft)
			if fOut.Kind == Unknown || fOut.Kind == Unsupported {
				continue
			}
			members[fName] = fOut

		}
		out.Members = members
		return out
	case *tc.Map:
		out := u.Type(name)
		if out.Kind != Unknown {
			return out
		}
		out = &Type{}
		out.Name = name
		out.Kind = Map
		elemType := t.Elem()
		et := WalkType(u, nil, elemType)
		if et.Kind == Unknown || et.Kind == Unsupported {
			out.Kind = et.Kind
			return out
		}
		keyType := t.Key()
		kt := WalkType(u, nil, keyType)
		if kt.Kind == Unknown || kt.Kind == Unsupported {
			out.Kind = kt.Kind
			return out
		}
		out.Elem = et
		out.Key = kt

		return out
	case *tc.Pointer:
		out := u.Type(name)
		if out.Kind != Unknown {
			return out
		}
		out = &Type{}
		out.Name = name
		out.Kind = Pointer
		out.Elem = WalkType(u, nil, t.Elem())
		return out
	case *tc.Slice:
		out := u.Type(name)
		if out.Kind != Unknown {
			return out
		}
		out = &Type{}
		out.Name = name
		out.Kind = Slice
		out.Elem = WalkType(u, nil, t.Elem())
		return out
	case *tc.Array:
		out := u.Type(name)
		if out.Kind != Unknown {
			return out
		}
		out = &Type{}
		out.Name = name
		out.Kind = Array
		out.Elem = WalkType(u, nil, t.Elem())
		return out
	case *tc.Chan:
		out := u.Type(name)
		if out.Kind != Unknown {
			return out
		}
		out = &Type{}
		out.Name = name
		out.Kind = Chan
		out.Elem = WalkType(u, nil, t.Elem())
		return out
	case *tc.Basic:
		out := u.Type(name)
		if out.Kind != Unknown {
			return out
		}
		out = &Type{}
		out.Name = name
		out.Kind = Builtin
		return out
	case *tc.Signature:
		out := u.Type(name)
		if out.Kind != Unknown {
			return out
		}

		out = &Type{}
		if t.Recv() != nil {
			out.Kind = Unsupported
			return out
		}
		out.Kind = Func
		out.Name = name
		out.Signature = convertSignature(u, t)

		return out
	case *tc.Interface:
		out := u.Type(name)
		if out.Kind != Unknown {
			return out
		}
		out = &Type{}
		out.Name = name
		out.Kind = Interface
		t.Complete()
		for i := 0; i < t.NumMethods(); i++ {
			if out.Methods == nil {
				out.Methods = make(map[Name]*Type)
			}
			method := t.Method(i)
			name := TcNameToName(method.String())
			mt := WalkType(u, &name, method.Type())
			out.Methods[name] = mt
		}
		return out
	case *tc.Named:
		if !token.IsExported(name.Name) {
			out := &Type{}
			out.Kind = Unsupported
			return out
		}
		out := u.Type(name)
		if out.Kind != Unknown {
			return out
		}
		return WalkType(u, &name, t.Underlying())

	default:
		out := &Type{}
		out.Kind = Unknown
		return out
	}

	out := &Type{}
	out.Kind = Unknown
	return out
}

func TcNameToName(in string) Name {
	// Detect anonymous type names. (These may have '.' characters because
	// embedded types may have packages, so we detect them specially.)
	if strings.HasPrefix(in, "struct{") ||
		strings.HasPrefix(in, "<-chan") ||
		strings.HasPrefix(in, "chan<-") ||
		strings.HasPrefix(in, "chan ") ||
		strings.HasPrefix(in, "func(") ||
		strings.HasPrefix(in, "*") ||
		strings.HasPrefix(in, "map[") ||
		strings.HasPrefix(in, "[") {
		return Name{Name: in}
	}

	// Otherwise, if there are '.' characters present, the name has a
	// package path in front.
	nameParts := strings.Split(in, ".")
	name := Name{Name: in}
	if n := len(nameParts); n >= 2 {
		// The final "." is the name of the type--previous ones must
		// have been in the package path.
		name.Package, name.Name = strings.Join(nameParts[:n-1], "."), nameParts[n-1]
	}
	return name
}

func isBasicField(f *tc.Var) (*tc.Basic, bool) {
	ft := f.Type()
	switch t := ft.(type) {
	case *tc.Basic:
		return t, true
	}

	return nil, false
}

func convertSignature(u Universe, t *tc.Signature) *Signature {
	signature := &Signature{}
	for i := 0; i < t.Params().Len(); i++ {
		signature.Parameters = append(signature.Parameters, WalkType(u, nil, t.Params().At(i).Type()))
	}
	for i := 0; i < t.Results().Len(); i++ {
		signature.Results = append(signature.Results, WalkType(u, nil, t.Results().At(i).Type()))
	}
	if r := t.Recv(); r != nil {
		signature.Receiver = WalkType(u, nil, r.Type())
	}
	signature.Variadic = t.Variadic()
	return signature
}

func tcFuncNameToName(in string) Name {
	name := strings.TrimPrefix(in, "func ")
	nameParts := strings.Split(name, "(")
	return TcNameToName(nameParts[0])
}

func tcVarNameToName(in string) Name {
	nameParts := strings.Split(in, " ")
	// nameParts[0] is "var".
	// nameParts[2:] is the type of the variable, we ignore it for now.
	return TcNameToName(nameParts[1])
}

func AddFunction(u Universe, useName *Name, in *tc.Func) *Type {
	name := tcFuncNameToName(in.String())
	if useName != nil {
		name = *useName
	}
	out := u.Function(name)
	out.Signature = convertSignature(u, in.Type().(*tc.Signature))
	out.Kind = DeclarationOf

	return out
}

func AddConstant(u Universe, useName *Name, in *tc.Const) *Type {
	name := tcVarNameToName(in.String())
	if useName != nil {
		name = *useName
	}
	out := u.Constant(name)
	out.Kind = DeclarationOf
	out.Underlying = WalkType(u, nil, in.Type())

	var constval string

	// For strings, we use `StringVal()` to get the un-truncated,
	// un-quoted string. For other values, `.String()` is preferable to
	// get something relatively human readable (especially since for
	// floating point types, `ExactString()` will generate numeric
	// expressions using `big.(*Float).Text()`.
	switch in.Val().Kind() {
	case constant.String:
		constval = constant.StringVal(in.Val())
	default:
		constval = in.Val().String()
	}

	out.ConstValue = &constval
	return out
}

func AddVariable(u Universe, useName *Name, in *tc.Var) *Type {
	name := tcVarNameToName(in.String())
	if useName != nil {
		name = *useName
	}

	out := u.Variable(name)
	out.Kind = DeclarationOf
	out.Underlying = WalkType(u, nil, in.Type())

	return out
}

func GetPkgData(pkgDir string, u Universe, pkgData map[string]Package) error {
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedModule}
	pkgs, err := packages.Load(cfg, pkgDir)
	if err != nil {
		return err
	}

	if packages.PrintErrors(pkgs) > 0 {
		return fmt.Errorf("error(s) occurred when import %s files", pkgDir)
	}

	for _, pkg := range pkgs {
		pkgPath := pkg.PkgPath
		name := pkg.Name
		p := Package{}
		p.PkgPath = pkgPath
		p.Name = name
		t := pkg.Types
		pTypes := make(map[string]*Type)
		pFunctions := make(map[string]*Type)
		pPkgVars := make(map[string]*Type)
		pConstant := make(map[string]*Type)
		s := t.Scope()
		for _, n := range s.Names() {
			obj := s.Lookup(n)
			tn, ok := obj.(*tc.TypeName)
			if ok {
				objT := WalkType(u, nil, tn.Type())
				if objT.Kind == Unknown || objT.Kind == Unsupported {
					continue
				}
				pTypes[objT.Name.Name] = objT
			}
			tf, ok := obj.(*tc.Func)
			// We only care about functions, not concrete/abstract methods.
			if ok && tf.Type() != nil && tf.Type().(*tc.Signature).Recv() == nil {
				objF := AddFunction(u, nil, tf)
				pFunctions[objF.Name.Name] = objF
			}
			tv, ok := obj.(*tc.Var)
			if ok && !tv.IsField() {
				objVar := AddVariable(u, nil, tv)
				pPkgVars[objVar.Name.Name] = objVar
			}
			tconst, ok := obj.(*tc.Const)
			if ok {
				objConst := AddConstant(u, nil, tconst)
				pConstant[objConst.Name.Name] = objConst
			}

		}
		p.Types = pTypes
		p.Functions = pFunctions
		p.Variables = pPkgVars
		p.Constants = pConstant
		imports := make(map[string]*Package)
		for k, p := range pkg.Imports {
			impPkg := &Package{}
			impPkg.PkgPath = p.PkgPath
			impPkg.Name = p.Name
			imports[k] = impPkg
		}
		p.Imports = imports
		pkgData[pkgPath] = p
	}

	return nil
}
