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

	"sysadm/command"
)

var conversionRegistryFunc runtime.FuncRegistry = RegisterConversions

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddConversionFunc((*Command)(nil), (*command.Command)(nil), func(a, b interface{}) error {
		return Convert_v1beta1_Command_To_command_Command(a.(*Command), b.(*command.Command))
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*command.Command)(nil), (*Command)(nil), func(a, b interface{}) error {
		return Convert_command_Command_To_v1beta1_Command(a.(*command.Command), b.(*Command))
	}); err != nil {
		return err
	}

	return nil
}

func Convert_v1beta1_Command_To_command_Command(in *Command, out *command.Command) error {
	out.ID = in.ID
	out.Command = in.Command
	out.Name = in.Name

	return nil
}

func Convert_command_Command_To_v1beta1_Command(in *command.Command, out *Command) error {
	out.ID = in.ID
	out.Command = in.Command
	out.Name = in.Name

	return nil
}