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

// TODO: this file should be CREATED by generator(will be coded)

package v1beta1

import (
	runtime "sysadm/apimachinery/runtime/v1beta1"
	"sysadm/syssetting"
)

var conversionRegistryFunc runtime.FuncRegistry = RegisterConversions

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddConversionFunc((*Syssetting)(nil), (*syssetting.Syssetting)(nil), func(a, b interface{}) error {
		return Convert_v1beta1_Syssetting_To_syssetting_Syssetting(a.(*Syssetting), b.(*syssetting.Syssetting))
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*syssetting.Syssetting)(nil), (*Syssetting)(nil), func(a, b interface{}) error {
		return Convert_syssetting_Syssetting_To_v1beta1_Syssetting(a.(*syssetting.Syssetting), b.(*Syssetting))
	}); err != nil {
		return err
	}

	return nil
}

func Convert_v1beta1_Syssetting_To_syssetting_Syssetting(in *Syssetting, out *syssetting.Syssetting) error {
	out.ID = in.ID
	out.Scope = in.Scope
	out.ReferenceObject = in.ReferenceObject
	out.Key = in.Key
	out.DefaultValue = in.DefaultValue
	out.Value = in.Value
	out.LastModifiedBy = in.LastModifiedBy
	out.LastModifiedTime = in.LastModifiedTime
	out.LastModifiedReason = in.LastModifiedReason
	out.LastValue = in.LastValue

	return nil
}

func Convert_syssetting_Syssetting_To_v1beta1_Syssetting(in *syssetting.Syssetting, out *Syssetting) error {
	out.ID = in.ID
	out.Scope = in.Scope
	out.ReferenceObject = in.ReferenceObject
	out.Key = in.Key
	out.DefaultValue = in.DefaultValue
	out.Value = in.Value
	out.LastModifiedBy = in.LastModifiedBy
	out.LastModifiedTime = in.LastModifiedTime
	out.LastModifiedReason = in.LastModifiedReason
	out.LastValue = in.LastValue

	return nil
}
