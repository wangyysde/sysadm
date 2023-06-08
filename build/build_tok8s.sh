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

KUBECTLBIN="${SYSADM_ROOT}/build/kubenetes/kubectl"
KUBECONF="${SYSADM_ROOT}/build/kubenetes/kubeconf.conf"
DEPLOYNANESPACE="sysadm"

declare -A APPLICATIONS
APPLICATIONS["sysadm"]="deploy/sysadm"
APPLICATIONS["registryctl"]="deploy/registryctl"
APPLICATIONS["infrastructure"]="deploy/infrastructure"
APPLICATIONS["apiserver"]="deploy/apiserver"

if [ "X${DEFAULT_REGISTRY_URL}" == "X" ]; then
   if [ -f "${SYSADM_ROOT}/build/build_common.sh" ]; then
       . "${SYSADM_ROOT}/build/build_common.sh"
   else
       echo "${SYSADM_ROOT}/build/build_common.sh was not found"
       exit 1
   fi
fi

function deploy::to::k8s(){
  package=$1
  imagever=$2

  appName=${APPLICATIONS[${package}]}
  if [ "X${appName}" == "X" ]; then
      echo "package name is not valid"
      exit 1
  fi

  echo "try to stop application ${appName}"
  ${KUBECTLBIN} scale --kubeconfig=${KUBECONF} --replicas=0 -n ${DEPLOYNANESPACE} ${appName}
  if [ $? != 0 ]; then
      echo "deploy application error"
      exit 2
  fi

  echo "pushing image to registry"
  ${DOCKER_BIN_PATH} push "${DEFAULT_REGISTRY_URL}${package}:${imagever}"
  if [ $? != 0 ]; then
    echo "push image to registry error"
    exit 3
  fi

  echo "pushing image to registry"
  ${KUBECTLBIN} scale --kubeconfig=${KUBECONF} --replicas=1 -n ${DEPLOYNANESPACE} ${appName}
  if [ $? != 0 ]; then
      echo "deploy application error"
      exit 4
  fi

  echo "deploy application ${appName} sucessful"
}