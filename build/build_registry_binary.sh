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
SYSADM_OUTPUT="${SYSADM_ROOT}/_output/bin"

set +e

if [ -z $1 ]; then
  error "Please set the 'registry version' variable"
  exit 1
fi

VERSION="$1"

set -e

# RM old registry binary package if it is exist
rm -rf ${SYSADM_OUTPUT}/registry || true

TEMP=`mktemp -d ${TMPDIR-/tmp}/distribution.XXXXXX`
echo "Trying to clone registry code from github......"
git clone -b $VERSION https://github.com/distribution/distribution.git $TEMP

# add patch 2879
if [ -e ${SYSADM_ROOT}/patches/registry/2879.patch ]
then
  echo 'add patch https://github.com/distribution/distribution/pull/2879 ...'
  cp ${SYSADM_ROOT}/patches/registry/2879.patch $TEMP
  cd $TEMP
  git apply 2879.patch
fi 
 
# add patch 3370
if [ -e ${SYSADM_ROOT}/patches/registry/3370.patch ]
then
 echo 'add patch https://github.com/distribution/distribution/pull/3370 ...'
 cp ${SYSADM_ROOT}/patches/registry/3370.patch $TEMP
 cd $TEMP
 git apply 3370.patch
fi

# add patch redis
if [ -e ${SYSADM_ROOT}/patches/registry/redis.patch ]
then
  echo 'add patch redis patch ...'
  cp ${SYSADM_ROOT}/patches/registry/redis.patch $TEMP
  cd $TEMP
  git apply redis.patch
fi

echo 'build the registry binary ...'
cp ${SYSADM_ROOT}/build/registry/Dockerfile.binary ${TEMP}
docker build -f ${TEMP}/Dockerfile.binary -t registry-golang $TEMP

echo 'copy the registry binary to local...'
ID=$(docker create registry-golang)
docker cp $ID:/go/src/github.com/docker/distribution/bin/registry ${SYSADM_OUTPUT}/registry

docker rm -f $ID
docker rmi -f registry-golang

echo "Build registry binary success, then to build photon image..."
cp ${TEMP}/cmd/registry/config-example.yml ${SYSADM_ROOT}/build/registry/config.yml
rm -rf $TEMP
