#!/bin/bash

# =============================================================
# @Author:  Wayne Wang <net_use@bzhy.com>
#
# @Copyright (c) 2022 Bzhy Network. All rights reserved.
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

SYSADM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
BASE_IMG_VER="1.0.0"
REGISTRY_VER="1.0.0"

if [ ! -z $1 ]; then
  REGISTRY_VER=$1
fi

if [ ! -e ${SYSADM_ROOT}/_output/bin/registry ]
then
  echo "registry binary package file not exist in ${SYSADM_ROOT}/_output/bin. You should build registry binary package file first"
  exit -1
fi

TEMP=`mktemp -d ${TMPDIR-/tmp}/sysadm.XXXXXX`
cp ${SYSADM_ROOT}/build/registry/Dockerfile.base ${TEMP}
cd ${TEMP}
echo "Now building sysadm-registry-base:${BASE_IMG_VER} ..."
docker build -f Dockerfile.base -t sysadm-registry-base:${BASE_IMG_VER} .
rm -rf Dockerfile.base
cp ${SYSADM_ROOT}/build/registry/Dockerfile ${TEMP}
cp ${SYSADM_ROOT}/_output/bin/registry ${TEMP}/
cp ${SYSADM_ROOT}/build/registry/entrypoint.sh ${TEMP}/
cp ${SYSADM_ROOT}/build/registry/install_cert.sh ${TEMP}/
cp ${SYSADM_ROOT}/build/registry/config.yml ${TEMP}/
echo "Now building sysadm_registry:${REGISTRY_VER} ..."
docker build -f Dockerfile --build-arg sysadm_base_image_version=${BASE_IMG_VER} -t sysadm_registry:${REGISTRY_VER} .
rm -rf ${TEMP}
docker rmi sysadm-registry-base:${BASE_IMG_VER}
