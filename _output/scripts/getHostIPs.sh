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

# gets the interfaces information in /etc/sysconfig/network-scripts/ifcfg-* files 
# on a host and retrun json format data
# 
[ -e ~/.bash_profile ] &&  . ~/.bash_profile || . /etc/profile

/usr/bin/echo -n "["
first="yes"
for file in `/usr/bin/ls  /etc/sysconfig/network-scripts/ifcfg-*`
do
	deviceLine=`/usr/bin/cat ${file} |/usr/bin/grep -i "DEVICE"`
	deviceName=`/usr/bin/echo ${deviceLine} |/usr/bin/cut -d "=" -f2 |/usr/bin/tr -d '"'`
	if [ "X${deviceName}" ==  "X"  -o "${deviceLine}" == "${deviceName}" ]; then 
		deviceName=`/usr/bin/echo ${deviceLine} |/usr/bin/cut -d "-" -f3 |/usr/bin/tr -d '"'`
	fi

	# interface name must not empty.	
	if [ "X${deviceName}" ==  "X"  -o "${deviceLine}" == "${deviceName}" ]; then
		continue
    if

    onbootLine=`/usr/bin/cat ${file} |/usr/bin/grep -i "ONBOOT"`
    onbootValue=`/usr/bin/echo ${onbootLine} |/usr/bin/cut -d "=" -f2 |/usr/bin/tr -d '"'`
	# skip the interface which not onboot
	if [ "X${onbootValue}" ==  "X"  -o "${onbootValue}" == "${onbootLine}" ]; then
		continue
	fi
	
	ipaddrLine=`/usr/bin/cat ${file} |/usr/bin/grep -i "IPADDR"`
	ipaddrValue=`/usr/bin/echo ${ipaddrLine} |/usr/bin/cut -d "=" -f2 |/usr/bin/tr -d '"'`
	# skip the interface which has dynamic IP
	if [ "X${ipaddrValue}" == "X" ]; then
		continue
    fi	

	prefixLine=`/usr/bin/cat ${file} |/usr/bin/grep -i "PREFIX"`
	prefixValue=`/usr/bin/echo ${prefixLine} |/usr/bin/cut -d "=" -f2 |/usr/bin/tr -d '"'`
    # Sysconfig.txt says that PREFIX takes precedence over
	# NETMASK when both are present.
	if [ "X${prefixValue}" ==  "X"  -o "${prefixValue}" == "${prefixLine}" ]; then
		networkLine=`/usr/bin/cat ${file} |/usr/bin/grep -i "NETWORK"`
        networkValue=`/usr/bin/echo ${networkLine} |/usr/bin/cut -d "=" -f2 |/usr/bin/tr -d '"'`
		if [ "X${networkValue}" ==  "X"  -o "${networkValue}" == "${networkLine}" ]; then
			prefixValue=${networkValue}
		fi
    fi
	
	ipv6addrLine=`/usr/bin/cat ${file} |/usr/bin/grep -i "IPV6ADDR"`
	ipv6addrValue=`/usr/bin/echo ${ipv6addrLine} |/usr/bin/cut -d "=" -f2 |/usr/bin/tr -d '"'`
	
	ipv6initLine=`/usr/bin/cat ${file} |/usr/bin/grep -i "IPV6INIT"`
    ipv6initValue=`/usr/bin/echo ${ipv6initLine} |/usr/bin/cut -d "=" -f2 |/usr/bin/tr -d '"'`

    hwaddrLine=`/usr/bin/cat ${file} |/usr/bin/grep -i "HWADDR"`
    hwaddrValue=`/usr/bin/echo ${hwaddrLine} |/usr/bin/cut -d "=" -f2 |/usr/bin/tr -d '"'`

	
    if [  "X${first}" == "Xyes" ]; then
		/usr/bin/echo -n "{ \"devname\": \"${deviceName}\", \"onboot\": \"${onbootValue}\", \"ipv4\": \"${ipaddrValue}\", \"prefix\": \"${prefixValue}\", \"ipv6\": \"${ipv6addrValue}\"}"
        first="no"
	else
		/usr/bin/echo -n ",{ \"devname\": \"${deviceName}\", \"onboot\": \"${onbootValue}\", \"ipv4\": \"${ipaddrValue}\", \"prefix\": \"${prefixValue}\", \"ipv6\": \"${ipv6addrValue}\"}"
	fi
done
/usr/bin/echo "]"
	
	
