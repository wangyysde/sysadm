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

PACKAGE_LIST="sysadm registryctl infrastructure"

echo "getting build information......"
SYSADM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
SYSADM_OUTPUT="${SYSADM_ROOT}/_output/bin"
GIT_COMMITID="$(git log --pretty=format:"%H" -1)"
BRANCH_NAME="$(git rev-parse --abbrev-ref HEAD)"
GITTREESTATS=""
TIMESTR="$(date '+%Y-%m-%d%Z%H:%M:%S')"
GOVERSION="$(go env GOVERSION)"
COMPILER="$(go env CC)"
ARCH="$(go env GOARCH)"
OS="$(go env GOOS)"

STATS="$(git status -s)"
[ -z "${STATS}" ] && GITTREESTATS="clean" || GITTREESTATS="noclean"

if [  ! -e ${SYSADM_OUTPUT} ]; then
	mkdir -p ${SYSADM_OUTPUT}
fi


function create::build::infofile(){
	package_name=$1
	build_file="${SYSADM_ROOT}/${package_name}/cmd/buildInfo.go"
	if [ "X${package_name}" == "X" ]; then
		echo "Package name  is not valid"
		return 1
	fi

	if [ ! -e "${SYSADM_ROOT}/${package_name}" ]; then
		echo "Package ${package_name} is not exist"
		return 1
	fi
	
    echo "Creating buildInfo for package ${package_name} ..."
	[ ! -e "${build_file}" ] && rm -rf "${build_file}"
	cat "${SYSADM_ROOT}/CopyRight" >"${build_file}"
	echo "" >> "${build_file}"
	echo "package cmd" >> "${build_file}"
	echo "" >> "${build_file}"
	echo "var gitCommitId = \"${GIT_COMMITID}\"" >> "${build_file}"
	echo "var branchName = \"${BRANCH_NAME}\"" >> "${build_file}"
	echo "var gitTreeStatus = \"${GITTREESTATS}\"" >> "${build_file}"
	echo "var buildDateTime = \"${TIMESTR}\"" >> "${build_file}"
	echo "var goVersion = \"${GOVERSION}\"" >> "${build_file}"
	echo "var compiler = \"${COMPILER}\"" >> "${build_file}"
	echo "var arch = \"${ARCH}\"" >> "${build_file}"
	echo "var hostos = \"${OS}\"" >> "${build_file}"
	return 0
}

function build::package(){
	package_name=$1
	if [ "X${package_name}" == "X"  ];then
		eccho "Package ${package_name} is not exist"
		return 1
	fi

	cd "${SYSADM_ROOT}" 
	goFiles="$(ls ${package_name}/*.go |tr "\n" " ")"
	echo  "Now building ${package_name} package. ${package_name} binary file will be placed into ${SYSADM_OUTPUT}/....."
    echo -n "go build -o ${SYSADM_OUTPUT}/${package_name} ${goFiles}"
	go build -o "${SYSADM_OUTPUT}/${package_name}" ${goFiles}
	[ $? -eq 0 ] && echo "[ Success ]" || echo "[ False ]"
}

BUILD_LIST=""
[ "X$@" == "X" ] && BUILD_LIST=${PACKAGE_LIST} || BUILD_LIST=$@

for p in ${BUILD_LIST}
do
   create::build::infofile ${p}
   [ $? -eq  1 ] && exit 1
   build::package ${p}
   [ $? -eq  1 ] && exit 1
done

exit 0
