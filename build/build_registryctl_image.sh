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

REGISTRYCTL_VER=$1
ISDEPLOY=$2

function create::dockerfile(){
	datetime=`date  '+%Y%m%d %H:%M:%S'`
	cat <<EOF >/${TEMP}/Dockerfile
FROM ${BASE_IMG}
LABEL Version=${REGISTRYCTL_VER} \\
	Maintainer="wayne.wang<${EMAIL}>" \\
	Built="${datetime}" \\
	Description="registryctl component of sysadm platform"
RUN mkdir -p /opt/registryctl/{bin,conf,logs,run} && useradd sysadm  
COPY ./registryctl /opt/registryctl/bin/
COPY ./entrypoint.sh /opt/registryctl/bin/
COPY ./registryctl.yaml /opt/registryctl/conf/
RUN chmod u+x /opt/registryctl/bin/entrypoint.sh && chmod u+x /opt/registryctl/bin/registryctl && chown -R sysadm:sysadm /opt/registryctl
ENTRYPOINT ["/opt/registryctl/bin/entrypoint.sh"]
EOF

}


if [ ! -e ${SYSADM_ROOT}/_output/bin/registryctl ]
then
  echo "registryctl binary package file not exist in ${SYSADM_ROOT}/_output/bin. You should build registryctl binary package file first"
  exit -1
fi


cd ${TEMP}
create::dockerfile
cp ${SYSADM_ROOT}/_output/bin/registryctl ${TEMP}/
cp ${SYSADM_ROOT}/build/registryctl/entrypoint.sh ${TEMP}/
cp ${SYSADM_ROOT}/_output/conf/registryctl.yaml ${TEMP}/
echo "Now building registryctl:${REGISTRYCTL_VER} ..."
docker build -f Dockerfile  -t ${DEFAULT_REGISTRY_URL}registryctl:${REGISTRYCTL_VER} .
if [ $? == 0 ]; then
  if [ "${ISDEPLOY}" == "y" -o "${ISDEPLOY}" == "Y" ]; then
       deploy::package "registryctl" ${REGISTRYCTL_VER}
  fi
#	docker push ${REGISTRY_URL}registryctl:${REGISTRYCTL_VER}
else 
	echo "build registryctl:${REGISTRYCTL_VER} image error"
fi
