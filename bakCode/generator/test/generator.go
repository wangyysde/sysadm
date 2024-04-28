package main

import (
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
	sysadmGenerator "sysadm/generator/v1beta1"
)

func main() {
	var pkgData = make(map[string]*sysadmGenerator.Package)
	var u = make(sysadmGenerator.Universe)

	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedModule}
	pkgs, err := packages.Load(cfg, "./data")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		pkgPath := pkg.PkgPath
		name := pkg.Name
		sysadmPkg := sysadmGenerator.Package{}
		sysadmPkg.PkgPath = pkgPath
		sysadmPkg.Name = name
		t := pkg.Types
		sysadmPkgTypes := make(map[string]*sysadmGenerator.Type)
		sysadmPkgFunctions := make(map[string]*sysadmGenerator.Type)
		sysadmPkgVars := make(map[string]*sysadmGenerator.Type)
		sysadmConstant := make(map[string]*sysadmGenerator.Type)
		s := t.Scope()
		for _, n := range s.Names() {
			obj := s.Lookup(n)
			tn, ok := obj.(*types.TypeName)
			if ok {
				objT := sysadmGenerator.WalkType(u, nil, tn.Type())
				if objT.Kind == sysadmGenerator.Unknown || objT.Kind == sysadmGenerator.Unsupported {
					continue
				}
				sysadmPkgTypes[objT.Name.Name] = objT
			}
			tf, ok := obj.(*types.Func)
			// We only care about functions, not concrete/abstract methods.
			if ok && tf.Type() != nil && tf.Type().(*types.Signature).Recv() == nil {
				objF := sysadmGenerator.AddFunction(u, nil, tf)
				sysadmPkgFunctions[objF.Name.Name] = objF
			}
			tv, ok := obj.(*types.Var)
			if ok && !tv.IsField() {
				objVar := sysadmGenerator.AddVariable(u, nil, tv)
				sysadmPkgVars[objVar.Name.Name] = objVar
			}
			tconst, ok := obj.(*types.Const)
			if ok {
				objConst := sysadmGenerator.AddConstant(u, nil, tconst)
				sysadmConstant[objConst.Name.Name] = objConst
			}

		}
		sysadmPkg.Types = sysadmPkgTypes
		sysadmPkg.Functions = sysadmPkgFunctions
		sysadmPkg.Variables = sysadmPkgVars
		sysadmPkg.Constants = sysadmConstant
		imports := make(map[string]*sysadmGenerator.Package)
		for k, p := range pkg.Imports {
			impPkg := &sysadmGenerator.Package{}
			impPkg.PkgPath = p.PkgPath
			impPkg.Name = p.Name
			imports[k] = impPkg
		}
		sysadmPkg.Imports = imports
		pkgData[pkgPath] = &sysadmPkg
	}

	//fmt.Printf("pkgData: %+v\n", pkgData)
	for _, p := range pkgData {
		//	fmt.Printf("pkgPath:%s\n", pkgPath)
		//	fmt.Printf("package: %+v\n", p)
		for tName, c := range p.Types {
			fmt.Printf("type Name: %s type: %+v \n", tName, *c)
			if c.Kind == sysadmGenerator.Struct {
				//		fmt.Printf("struct members:\n")
				for n, m := range c.Members {
					fmt.Printf("Name:%+v member: %+v\n", n, *m)
				}
			}
		}
	}
}
