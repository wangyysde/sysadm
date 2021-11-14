#!/usr/bin/env bash

# =============================================================
# @Author:  Wayne Wang <net_use@bzhy.com>
#
# @Copyright (c) 2021 Bzhy Network. All rights reserved.
# @HomePage http://www.sysadm.cn
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at: 
# http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and  limitations under the License.
# @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html
#

set -o errexit
set -o nounset
set -o pipefail

SYSADM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
SYSADM_OUTPUT="${SYSADM_ROOT}/_output/bin"

function all_info(){
	echo "# Build code."
	echo "#"
	echo "# Args:"
	echo "#   WHAT: Directory(Pacakge) names to build.  If any of these directories has a 'main'"
	echo "#     package, the build will produce executable files under $(SYSADM_OUTPUT)."
	echo "#     If not specified, \"everything\" will be built."
	echo "#   GOFLAGS: Extra flags to pass to 'go' when building."
	echo "#   GOLDFLAGS: Extra linking flags passed to 'go' when building."
	echo "#   GOGCFLAGS: Additional go compile flags passed to 'go' when building."
	echo "#"
	echo "# Example:"
	echo "#   make"
	echo "#   make all"
	echo "#   make all WHAT=sysadm GOFLAGS=-v"
	echo "#   make all GOLDFLAGS=\"\""
	echo "#     Note: Specify GOLDFLAGS as an empty string for building unstripped binaries, which allows"
	echo "#           you to use code debugging tools like delve. When GOLDFLAGS is unspecified, it defaults"
	echo "#           to \"-s -w\" which strips debug information. Other flags that can be used for GOLDFLAGS"
	echo "#           are documented at https://golang.org/cmd/link/"
}


function other_info() {
	echo "# Building, clean or install code."
	echo "#"
	echo "# Commands:"
	echo "#    make all building  all packages"
	echo "#    make all WHAT=<packagename> building specified package named packagename"
	echo "#    make clean  clean all cached files of last built"
}
case "$@" in
	all)
		all_info 
		;;
	*)
		other_info
		;;
esac

exit 0