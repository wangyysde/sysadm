#!/usr/bin/env bash

# * =============================================================
# * @Author:  Wayne Wang <net_use@bzhy.com>
# *
# * @Copyright (c) 2024 Bzhy Network. All rights reserved.
# * @HomePage http://www.sysadm.cn
# *
# * Licensed under the Apache License, Version 2.0 (the "License");
# * you may not use this file except in compliance with the License.
# * You may obtain a copy of the License at:
# * http://www.apache.org/licenses/LICENSE-2.0
# * Unless required by applicable law or agreed to in writing, software
# * distributed under the License is distributed on an "AS IS" BASIS,
# * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# * See the License for the specific language governing permissions and  limitations under the License.
#* @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html

set -o errexit
set -o nounset
set -o pipefail

GENERATED_FILE="zz_generated.runtime.go"
GOHEADERFILE_FILE="CopyRight"
APIDIRS="command syssetting"
MODULENAME="sysadm"
CONVERSIONTAG="// +sysadm:api-resource=true"

function generate::runtime::files(){
    SYSADM_ROOT=$1
    cd ${SYSADM_ROOT}

    APIROOTS=${APIROOTS:-$(git grep --files-with-matches -e "${CONVERSIONTAG}" ${APIDIRS} | \
    xargs -n 1 dirname | sort | uniq)}

    importPackages=""
    for item in ${APIROOTS}
    do
       rootPathName=$(echo ${item} |cut -d "/" -f1)
       package=$(echo ${item} |tr '/' '\n'  |tail -n1)
       aliasName="${rootPathName}${package}"
       if [ "X${importPackages}" == "X" ]; then
          importPackages="${MODULENAME}/${item}"
       else
          importPackages="${importPackages} ${MODULENAME}/${item}"
        fi
    done

    if [ "X${importPackages}" == "X" ]; then
       return
    fi

    GENERATEDFILE="${SYSADM_ROOT}/apiserver/app/${GENERATED_FILE}"
    [[ -f "${GENERATEDFILE}" ]] && rm -rf "${GENERATEDFILE}"
    cat "${SYSADM_ROOT}/${GOHEADERFILE_FILE}" >"${GENERATEDFILE}"
    echo -e "" >>"${GENERATEDFILE}"
    echo -e "" >>"${GENERATEDFILE}"
    echo "package app" >>"${GENERATEDFILE}"
    echo "" >>"${GENERATEDFILE}"
    echo "import (" >>"${GENERATEDFILE}"

    for item in ${importPackages}
    do
        echo -e "\t_ \"${item}\"" >>"${GENERATEDFILE}"
    done

    echo ")" >>"${GENERATEDFILE}"
    echo "" >>"${GENERATEDFILE}"
}

