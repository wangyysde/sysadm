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

SYSADM_VER=$1
ISDEPLOY=$2


function create::dockerfile(){
	datetime=`date  '+%Y%m%d %H:%M:%S'`
	cat <<EOF >/${TEMP}/Dockerfile
FROM ${BASE_IMG}
LABEL Version=${SYSADM_VER} \\
	Maintainer="wayne.wang<${EMAIL}>" \\
	Built="${datetime}" \\
	Description="sysadm component of sysadm platform"
RUN mkdir -p /opt/sysadm/{bin,conf,logs,run,html,tmpls} && useradd sysadm && apt-get update && apt install -y curl
COPY ./sysadm /opt/sysadm/bin/
COPY ./entrypoint.sh /opt/sysadm/bin/
COPY ./config.yaml /opt/sysadm/conf/
COPY ./html /opt/sysadm/html/
COPY ./tmpls /opt/sysadm/tmpls
RUN chmod u+x /opt/sysadm/bin/entrypoint.sh && chmod u+x /opt/sysadm/bin/sysadm && chown -R sysadm:sysadm /opt/sysadm 
ENTRYPOINT ["/opt/sysadm/bin/entrypoint.sh"]
EOF

}


if [ ! -e ${SYSADM_ROOT}/_output/bin/sysadm ]
then
  echo "sysadm binary package file not exist in ${SYSADM_ROOT}/_output/bin. You should build sysadm binary package file first"
  exit -1
fi


cd ${TEMP}
create::dockerfile
cp ${SYSADM_ROOT}/_output/bin/sysadm ${TEMP}/
cp ${SYSADM_ROOT}/build/sysadm/entrypoint.sh ${TEMP}/
cp ${SYSADM_ROOT}/_output/conf/config.yaml ${TEMP}/
cp -r ${SYSADM_ROOT}/_output/html ${TEMP}/
cp -r ${SYSADM_ROOT}/_output/tmpls ${TEMP}/
echo "Now building sysadm:${SYSADM_VER} ..."
docker build -f Dockerfile  -t "${DEFAULT_REGISTRY_URL}sysadm:${SYSADM_VER}" .
if [ $? == 0 ]; then
  if [ "${ISDEPLOY}" == "y" -o "${ISDEPLOY}" == "Y" ]; then
       deploy::package "sysadm" ${SYSADM_VER}
  fi
else 
	echo "build sysadm:${SYSADM_VER} image error"
fi
