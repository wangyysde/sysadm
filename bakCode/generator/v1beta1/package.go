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

// Type gets the given Type in this Package.  If the Type is not already
// defined, this will add it and return the new Type value.  The caller is
// expected to finish initialization.
func (p *Package) Type(typeName string) *Type {
	if t, ok := p.Types[typeName]; ok {
		return t
	}
	if p.PkgPath == "" {
		// Import the standard builtin types!
		if t, ok := builtins.Types[typeName]; ok {
			p.Types[typeName] = t
			return t
		}
	}
	t := &Type{Name: Name{Package: p.PkgPath, Name: typeName}}
	p.Types[typeName] = t
	return t
}

// Function gets the given function Type in this Package. If the function is
// not already defined, this will add it.  If a function is added, it's the
// caller's responsibility to finish construction of the function by setting
// Underlying to the correct type.
func (p *Package) Function(funcName string) *Type {
	if t, ok := p.Functions[funcName]; ok {
		return t
	}
	t := &Type{Name: Name{Package: p.PkgPath, Name: funcName}}
	t.Kind = DeclarationOf
	p.Functions[funcName] = t
	return t
}

// Variable gets the given variable Type in this Package. If the variable is
// not already defined, this will add it. If a variable is added, it's the caller's
// responsibility to finish construction of the variable by setting Underlying
// to the correct type.
func (p *Package) Variable(varName string) *Type {
	if t, ok := p.Variables[varName]; ok {
		return t
	}
	t := &Type{Name: Name{Package: p.PkgPath, Name: varName}}
	t.Kind = DeclarationOf
	p.Variables[varName] = t
	return t
}

// Constant gets the given constant Type in this Package. If the constant is
// not already defined, this will add it. If a constant is added, it's the caller's
// responsibility to finish construction of the constant by setting Underlying
// to the correct type.
func (p *Package) Constant(constName string) *Type {
	if t, ok := p.Constants[constName]; ok {
		return t
	}
	t := &Type{Name: Name{Package: p.PkgPath, Name: constName}}
	t.Kind = DeclarationOf
	p.Constants[constName] = t
	return t
}
