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
if [ -f "${SYSADM_ROOT}/build/build_common.sh" ]; then
    . "${SYSADM_ROOT}/build/build_common.sh"
else
    echo "${SYSADM_ROOT}/build/build_common.sh was not found"
    exit -2
fi

TEMP=`mktemp -d ${TMPDIR-/tmp}/sysadm.XXXXXX`

INFRASTRUCTURE_VER=$1
ISDEPLOY=$2

function create::dockerfile(){
	datetime=`date  '+%Y%m%d %H:%M:%S'`
	cat <<EOF >/${TEMP}/Dockerfile
FROM ${BASE_IMG}
LABEL Version=${INFRASTRUCTURE_VER} \\
	Maintainer="wayne.wang<${EMAIL}>" \\
	Built="${datetime}" \\
	Description="infrastructure component of sysadm platform"
RUN mkdir -p /opt/infrastructure/{bin,conf,logs,run} && useradd sysadm 
COPY ./infrastructure /opt/infrastructure/bin/
COPY ./entrypoint.sh /opt/infrastructure/bin/
COPY ./infrastructure.yaml /opt/infrastructure/conf/
RUN chmod u+x /opt/infrastructure/bin/entrypoint.sh && chmod u+x /opt/infrastructure/bin/infrastructure && chown -R sysadm:sysadm /opt/infrastructure 
ENTRYPOINT ["/opt/infrastructure/bin/entrypoint.sh"]
EOF

}


if [ ! -e ${SYSADM_ROOT}/_output/bin/infrastructure ]
then
  echo "infrastructure binary package file not exist in ${SYSADM_ROOT}/_output/bin. You should build infrastructure binary package file first"
  exit -1
fi


cd ${TEMP}
create::dockerfile
cp ${SYSADM_ROOT}/_output/bin/infrastructure ${TEMP}/
cp ${SYSADM_ROOT}/build/infrastructure/entrypoint.sh ${TEMP}/
cp ${SYSADM_ROOT}/_output/conf/infrastructure.yaml ${TEMP}/
echo "Now building infrastructure:${INFRASTRUCTURE_VER} ..."
docker build -f Dockerfile  -t "${DEFAULT_REGISTRY_URL}infrastructure:${INFRASTRUCTURE_VER}" .
if [ $? == 0 ]; then
  if [ "${ISDEPLOY}" == "y" -o "${ISDEPLOY}" == "Y" ]; then
     deploy::package "infrastructure" ${INFRASTRUCTURE_VER}
  fi
#	docker push ${REGISTRY_URL}infrastructure:${INFRASTRUCTURE_VER}
else 
	echo "build infrastructure:${INFRASTRUCTURE_VER} image error"
fi
