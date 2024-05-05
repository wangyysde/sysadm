/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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
 */

package pki

import (
	"crypto/x509"
	"fmt"
	"net"
	"testing"

	netutils "k8s.io/utils/net"
)

func TestCreateCA(t *testing.T) {
	_, _, caPem, caKeyPem, err := CreateCertificateAuthority("sysadm", []string{"sysadm", "bzhy.com"}, 90, x509.RSA)
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.Fatal("failed\n")
	}

	fmt.Printf("CA cert: %s\n", caPem)
	fmt.Printf("CA Key: %s \n", caKeyPem)
	fmt.Printf("CA certificate and private key has created\n")
}

func TestCreateServerCert(t *testing.T) {
	_, _, caPem, caKeyPem, err := CreateCertificateAuthority("sysadm", []string{"sysadm", "bzhy.com"}, defaultCaPeriodDays, x509.RSA)
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.Fatal("create CA failed\n")
	}

	ip1 := netutils.ParseIPSloppy("192.168.0.10")
	ip2 := netutils.ParseIPSloppy("172.28.1.103")
	ip3 := netutils.ParseIPSloppy("127.0.0.1")

	ips := []net.IP{ip1, ip2, ip3}
	_, _, certPem, keyPem, err := CreateCertAndKey(x509.RSA, string(caPem), string(caKeyPem), ips, 365, "testcert",
		[]string{"localhost", "WIN-20230717WDZ", "cp1"}, []string{"sysadm", "bzhy.com"}, []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth})
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.Fatal("create Certificate failed\n")
	}

	fmt.Printf("cert: %s\n", certPem)
	fmt.Printf("cert Key: %s \n", keyPem)
	fmt.Printf("server certificate and private key has created\n")
}

func TestCreateClientCert(t *testing.T) {
	_, _, caPem, caKeyPem, err := CreateCertificateAuthority("sysadm", []string{"sysadm", "bzhy.com"}, defaultCaPeriodDays, x509.RSA)
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.Fatal("create CA failed\n")
	}

	_, _, certPem, keyPem, err := CreateCertAndKey(x509.RSA, string(caPem), string(caKeyPem), []net.IP{}, 365, "testcert",
		[]string{}, []string{"sysadm", "bzhy.com"}, []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth})
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.Fatal("create Certificate failed\n")
	}

	fmt.Printf("cert: %s\n", certPem)
	fmt.Printf("cert Key: %s \n", keyPem)
	fmt.Printf("client certificate and private key has created\n")
}
