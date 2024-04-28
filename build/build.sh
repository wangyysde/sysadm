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
if [ -f "${SYSADM_ROOT}/build/build_common.sh" ]; then
    . "${SYSADM_ROOT}/build/build_common.sh"
else
    echo "${SYSADM_ROOT}/build/build_common.sh was not found"
    exit 1
fi

if [ -f "${SYSADM_ROOT}/build/build_tok8s.sh" ]; then
    . "${SYSADM_ROOT}/build/build_tok8s.sh"
else
    echo "${SYSADM_ROOT}/build/build_tok8s.sh was not found"
    exit 1
fi

echo "getting build information......"
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

  if [ "X${package_name}" == "Xapiserver"  ];then
      . "${SYSADM_ROOT}/build/update_generated_runtime.sh"
      generate::runtime::files "${SYSADM_ROOT}"
  fi

	cd "${SYSADM_ROOT}" 
	goFiles="$(ls ${package_name}/*.go |tr "\n" " ")"
	echo  "Now building ${package_name} package. ${package_name} binary file will be placed into ${SYSADM_OUTPUT}/....."
    echo -n "go build -o ${SYSADM_OUTPUT}/${package_name} ${goFiles}"
        # static compiling 
	# GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o "${SYSADM_OUTPUT}/${package_name}" ${goFiles}
	# static and cross compiling 
	# GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o "${SYSADM_OUTPUT}/${package_name}" ${goFiles}
 	go build -o "${SYSADM_OUTPUT}/${package_name}" ${goFiles}
	[ $? -eq 0 ] && echo "[ Success ]" || echo "[ False ]"
}

BUILD_LIST=""
WHAT=""
BUILD_IMAGE="n"
DEPLOY_BY_DOCKECOMPOSE="n"
DEPLOYTYPE=""
IMAGEVER=${DEFAULT_IMAGE_VER}
REGISTRY_URL=${DEFAULT_REGISTRY_URL}
if [ $# != 0 ]; then
  WHAT=$1
  shift
fi

if [ $# != 0 ]; then
  BUILD_IMAGE=$1
  [ "X${BUILD_IMAGE}" == "X" ] && BUILD_IMAGE="n"
  shift
fi

if [ $# != 0 ]; then
  IMAGEVER=$1
  [ "X${IMAGEVER}" == "X" ] && IMAGEVER=${DEFAULT_IMAGE_VER}
  shift
fi

if [ $# != 0 ]; then
  DEPLOY=$1
  [ "X${DEPLOY}" == "X" ] && DEPLOY="n"
  shift
fi

if [ $# != 0 ]; then
  DEPLOYTYPE=$1
  shift
fi
[ "X${DEPLOYTYPE}" == "X" ] && DEPLOYTYPE=${DEFAULT_DEPLOY_TYPE}

# if [ "${DEPLOY}" == "y" -o "${DEPLOY}" == "Y" ]; then
#  BUILD_IMAGE="y"
#  DEPLOY_BY_DOCKECOMPOSE="y"
# else
#   BUILD_IMAGE="n"
# fi

if [ "X${DEPLOYTYPE}" == "X${DEFAULT_DEPLOY_TYPE}" ]; then
  DEPLOY_BY_DOCKECOMPOSE="n"
fi

if [ "X${WHAT}" == "X" ]; then
  BUILD_LIST=${PACKAGE_LIST}
elif [ "X${WHAT}" == "Xall" ]; then
  BUILD_LIST=${PACKAGE_LIST}
else 
  BUILD_LIST=${WHAT} 
fi

[ "X${IMAGEVER}" == "X" ] && IMAGEVER=${DEFAULT_IMAGE_VER} 

if [ "X${DEPLOYTYPE}" == "Xrsync" ]; then
   . "${SYSADM_ROOT}/build/build_tok8s.sh"
   deploy::staticfile
   exit 0
fi


BUILD_LIST_ARRAY=(${BUILD_LIST//,/ })
for p in ${BUILD_LIST_ARRAY[@]}
do
   create::build::infofile ${p}
   [ $? -eq  1 ] && exit 1
   build::package ${p}
   [ $? -eq  1 ] && exit 1

  if [ "${BUILD_IMAGE}" == "y" -o "${BUILD_IMAGE}" == "Y" ]; then
	if [ -e "${SYSADM_ROOT}/build/build_${p}_image.sh" ]; then
		"${SYSADM_ROOT}/build/build_${p}_image.sh" "${IMAGEVER}" "${DEPLOY_BY_DOCKECOMPOSE}"
		if [ $? -ne 0 ]; then
			echo "building ${p} image error"
			exit 1
		fi

	else
		echo "${SYSADM_ROOT}/build/build_${p}_image.sh script file not exist"
		exit 1
	fi
  fi

  if [ "${DEPLOY}" == "y" -o "${DEPLOY}" == "Y" ]; then
      if [ "X${DEPLOYTYPE}" == "X${DEFAULT_DEPLOY_TYPE}" -o "X${DEPLOYTYPE}" == "X${DEPLOY_TYPE_K8SCURL}" ]; then
          deploy::to::k8s "${p}" "${IMAGEVER}" "${DEPLOYTYPE}"
      fi

      if [ "X${DEPLOYTYPE}" == "X${DEPLOY_TYPE_COMPOSE}" ]; then
          deploy::package "${p}" "${IMAGEVER}"
      fi
  fi
done

exit 0
