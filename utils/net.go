/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
* @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html

ErrorCode: 501xxx
*/

package utils

import (
	"fmt"
	"net"
	"strings"

	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

/*
  CheckIpAddress checking IP address.
  return net.IP if address is a valid IP address or address can be resolve to an IP when isLocal is false.
  CheckIpAddress will check the address is a valid IP address and it is one of the IP address of the local NIC if  isLocal is true.
  return IP(net.IP) if the ip address is valid
  Or return nil with error
*/
func CheckIpAddress(address string, isLocal bool) (net.IP, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(501001,"debug","now checking IP address %s ",address))
	if len(address) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501002,"error","The address(%s) is empty or the length of it is less 1",address))
		return nil, errs
	}

	if ip := net.ParseIP(address); ip != nil {
		if !isLocal {
			return ip, errs
		}

		if address == "0.0.0.0" || address == "::" {
			return ip, errs
		}

		adds,err := net.InterfaceAddrs()
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(501003,"error","Get interface address error: %s",err))
			return nil, errs
		}

		for _,v := range adds {
			ipnet,ok := v.(*net.IPNet)
			if !ok {
				continue
			}
			if ip.Equal(ipnet.IP) {
				return ip, errs
			}
		}

		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501004,"error","The address(%s) is not any of the addresses of host interfaces.",address))
		return nil, errs
	}

	ips,err := net.LookupIP(address)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501005,"error","Lookup the IP of address(%s) error %s.",address,err))
		return nil , errs
	}

	if !isLocal {
		return ips[0], errs
	}

	adds,err := net.InterfaceAddrs()
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501006,"error","Get interface address error: %s",err))
		return nil, errs
	}

	for _,ip := range ips {
		for _,v := range adds {
			ipnet,ok := v.(*net.IPNet)
			if !ok {
				continue
			}
			if ip.Equal(ipnet.IP) {
				return ip, errs
			}
		}
	}
	
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(501007,"error","The IP(%v) to the address(hostname:%s) is not any the IP address of host interfaces.",ips,address))
	return nil, errs
}

// check the validity of port 
// return port with nozero if the port is valid 
// Or return 0 with error
func CheckPort(port int)(int, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(501008,"debug","now checking port number %d ",port))
	if port > 0 && port <= 65535 {
		return port,errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(501009,"error","The port should be great than 1024 and less than 65535. Now is :%d",port))
	return 0, errs
}

/*
	GetRequestData get query data or postform data on a Context for keys([]string)
	the data on Query will be returned if there are the same key on both query and postform.
	return a pointer point to map[string]string which including the data that have be found.
	return nil and errs if there is not any data has be found.
*/
func GetRequestData(c *sysadmServer.Context, keys []string)(map[string]string, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	if c == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501010,"error","can not get data on nil context"))
		return nil, errs
	}

	if len(keys) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501011,"error","no data should be got"))
		return nil, errs
	}

	ret := make(map[string]string,0)
	for _,k := range keys {
		queryData,okQuery := c.GetQuery(k)
		if okQuery {
			ret[k] = queryData
		}else{
			formData,okForm := c.GetPostForm(k)
			if okForm {
				ret[k] = formData
			}else{
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(501013,"debug","the data for key %s was not found",k))
			}
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501012,"debug","try to get data for key: %s value: %s",k,ret[k]))
	}

	if len(ret) < 1{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501014,"debug","not found any data for keys"))
		return nil, errs
	}

	return ret, errs
}

/*
	GetRequestDataArray get query array data or postform data on a Context for keys([]string)
	the data on Query will be returned if there are the same key on both query and postform.
	return a pointer point to map[string]string which including the data that have be found.
	return nil and errs if there is not any data has be found.
*/
func GetRequestDataArray(c *sysadmServer.Context, keys []string)(map[string][]string, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	if c == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501015,"error","can not get data on nil context"))
		return nil, errs
	}

	if len(keys) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501016,"error","no data should be got"))
		return nil, errs
	}

	ret := make(map[string][]string,0)
	for _,k := range keys {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501017,"debug","try to get data for key: %s",k))
		queryData,okQuery := c.GetQueryArray(k)
		if okQuery {
			ret[k] = queryData
		}else{
			formData,okForm := c.GetPostFormArray(k)
			if okForm {
				ret[k] = formData
			}else{
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(501018,"debug","the data for key %s was not found",k))
			}
		}
	}

	if len(ret) < 1{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501019,"debug","not found any data for keys"))
		return nil, errs
	}

	return ret, errs
}

/*
	GetRequestDataMap get query array data or postform data on a Context for keys([]string)
	the data on Query will be returned if there are the same key on both query and postform.
	return a pointer point to map[string]string which including the data that have be found.
	return nil and errs if there is not any data has be found.
*/
func GetRequestDataMap(c *sysadmServer.Context, keys []string)(map[string]map[string]string, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror
	if c == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501020,"error","can not get data on nil context"))
		return nil, errs
	}

	if len(keys) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501021,"error","no data should be got"))
		return nil, errs
	}

	ret := make(map[string]map[string]string,0)
	for _,k := range keys {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501022,"debug","try to get data for key: %s",k))
		queryData,okQuery := c.GetQueryMap(k)
		if okQuery {
			ret[k] = queryData
		}else{
			formData,okForm := c.GetPostFormMap(k)
			if okForm {
				ret[k] = formData
			}else{
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(501023,"debug","the data for key %s was not found",k))
			}
		}
	}

	if len(ret) < 1{
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(501024,"debug","not found any data for keys"))
		return nil, errs
	}

	return ret, errs
}

/*
	check the existence of key in dataSet, return the value of the dataSet[key] with trimspaced if the key is exist in dataSet 
	otherwrise return ""  
*/
func GetKeyData(dataSet map[string]string, key string)string{
	ret := ""

	if strings.TrimSpace(key) == "" {
		return ""
	}

	value, ok := dataSet[key]
	if !ok {
		ret = value
	} else {
		ret = strings.TrimSpace(value)
	}

	return ret
}

// GetLocalIPs get ip address set on the interfaces on the local host. and return ip list as []string if get successfule.
// otherwise return empty slice and an error.
func GetLocalIPs()([]string, error){
	var ips []string

	adds,err := net.InterfaceAddrs()
	if err != nil {
		return ips, fmt.Errorf("can not get host unicast IP list %s", err)
	}
	
	for _,v := range adds {
		ipnet,ok := v.(*net.IPNet)
		if !ok {
			continue
		}
		ipstr := Bytes2str(ipnet.IP)
		if FoundStrInSlice(ips,ipstr,true) {
			continue
		}
		ips = append(ips, ipstr)
	}

	return ips, nil
}

// GetLocalMacs get mac information  on the interfaces on the local host. and return mac list as []string if get successfule.
// otherwise return empty slice and an error.
func GetLocalMacs()([]string,error){
	var macs []string

	ints,err := net.Interfaces()
	if err != nil {
		return macs, fmt.Errorf("can not get host interfaces information %s", err)
	}

	for _, dev := range ints {
		mac := Bytes2str(dev.HardwareAddr)
		macs = append(macs, mac)
	}

	return macs, nil
}