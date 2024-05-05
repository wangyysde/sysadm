#!/bin/bash

# =============================================================
# @Author:  Wayne Wang <net_use@bzhy.com>
#
# @Copyright (c) 2023 Bzhy Network. All rights reserved.
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
INFRASTRUCTURE_VER="1.0"
TEMP=`mktemp -d ${TMPDIR-/tmp}/sysadm.XXXXXX`
EMAIL="net_use@bzhy.com"

if [ -f "${SYSADM_ROOT}/build/build_common.sh" ]; then
    . "${SYSADM_ROOT}/build/build_common.sh"
else
    echo "${SYSADM_ROOT}/build/build_common.sh was not found"
    exit 1
fi

APISERVER_VER=$1
ISDEPLOY=$2

function create::dockerfile(){
	datetime=`date  '+%Y%m%d %H:%M:%S'`
	cat <<EOF >/${TEMP}/Dockerfile
FROM ${BASE_IMG}
LABEL Version=${AGENT_VER} \\
	Maintainer="wayne.wang<${EMAIL}>" \\
	Built="${datetime}" \\
	Description="apiserver component of sysadm platform"
RUN mkdir -p /opt/apiserver/bin /opt/apiserver/conf /opt/apiserver/logs /opt/apiserver/pki && useradd sysadm
COPY ./apiserver /opt/apiserver/bin/
COPY ./entrypoint.sh /opt/apiserver/bin/
COPY ./apiserver.yaml /opt/apiserver/conf/
RUN chmod u+x /opt/apiserver/bin/entrypoint.sh && chmod u+x /opt/apiserver/bin/apiserver && chown -R sysadm:sysadm /opt/apiserver 
ENTRYPOINT ["/opt/apiserver/bin/entrypoint.sh"]
EOF

}


if [ ! -e ${SYSADM_ROOT}/_output/bin/apiserver ]
then
  echo "apiserver binary package file not exist in ${SYSADM_ROOT}/_output/bin. You should build apiserver binary package file first"
  exit -1
fi


cd ${TEMP}
create::dockerfile
cp ${SYSADM_ROOT}/_output/bin/apiserver ${TEMP}/
cp ${SYSADM_ROOT}/build/apiserver/entrypoint.sh ${TEMP}/
cp ${SYSADM_ROOT}/_output/conf/apiserver.yaml ${TEMP}/
echo "Now building apiserver:${APISERVER_VER} ..."
docker build -f Dockerfile  -t ${DEFAULT_REGISTRY_URL}apiserver:${APISERVER_VER} .
if [ $? == 0 ]; then
	docker push ${DEFAULT_REGISTRY_URL}apiserver:${APISERVER_VER}
else 
	echo "build apiserver:${APISERVER_VER} image error"
fi
