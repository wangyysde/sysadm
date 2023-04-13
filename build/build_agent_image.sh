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
INFRASTRUCTURE_VER="1.0"
BASE_IMG="harbor.bzhy.com/os/centos:centos7.9.2009"
TEMP=`mktemp -d ${TMPDIR-/tmp}/sysadm.XXXXXX`
EMAIL="net_use@bzhy.com"

AGENT_VER=$1
REGISTRY_URL=$2

function create::dockerfile(){
	datetime=`date  '+%Y%m%d %H:%M:%S'`
	cat <<EOF >/${TEMP}/Dockerfile
FROM ${BASE_IMG}
LABEL Version=${AGENT_VER} \\
	Maintainer="wayne.wang<${EMAIL}>" \\
	Built="${datetime}" \\
	Description="agent component of sysadm platform"
RUN mkdir -p /opt/agent/{bin,conf,logs,run} && useradd sysadm 
COPY ./agent /opt/agent/bin/
COPY ./entrypoint.sh /opt/agent/bin/
COPY ./agent.yaml /opt/agent/conf/
RUN chmod u+x /opt/agent/bin/entrypoint.sh && chmod u+x /opt/agent/bin/agent && chown -R sysadm:sysadm /opt/agent 
ENTRYPOINT ["/opt/agent/bin/entrypoint.sh"]
EOF

}


if [ ! -e ${SYSADM_ROOT}/_output/bin/agent ]
then
  echo "agent binary package file not exist in ${SYSADM_ROOT}/_output/bin. You should build agent binary package file first"
  exit -1
fi


cd ${TEMP}
create::dockerfile
cp ${SYSADM_ROOT}/_output/bin/agent ${TEMP}/
cp ${SYSADM_ROOT}/build/agent/entrypoint.sh ${TEMP}/
cp ${SYSADM_ROOT}/_output/conf/agent.yaml ${TEMP}/
echo "Now building agent:${AGENT_VER} ..."
docker build -f Dockerfile  -t ${REGISTRY_URL}agent:${AGENT_VER} .
if [ $? == 0 ]; then
	docker push ${REGISTRY_URL}agent:${AGENT_VER}
else 
	echo "build agent:${AGENT_VER} image error"
fi
