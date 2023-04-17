/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* https://www.sysadm.cn/licenses/apache-2.0.txt
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
*
*/

package init

import (
	"fmt"
	"github.com/pkg/errors"
	certsphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/pkiutil"
)

func runCerts(data *initData) error {
	var lastCACert *certsphase.KubeadmCert
	for _, cert := range certsphase.GetDefaultCertList() {
		if cert.CAName == "" {
			if err := createCACertAndKeyFiles(cert, data); err != nil {
				fmt.Printf("certCreate createCACertAndKeyFiles err: %v\n", err)
				return err
			}
			lastCACert = cert
		} else {
			if err := createCertAndKeyFilesWithCA(cert, lastCACert, data); err != nil {
				fmt.Printf("certCreate createCertAndKeyFilesWithCA err: %v\n", err)
				return err
			}
		}
	}

	if err := createCertsSa(data); err != nil {
		fmt.Printf("certCreate createCertsSa err: %v\n", err)
		return err
	}
	fmt.Printf("[certs] Using certificateDir folder %q\n", data.CertificateWriteDir())

	return nil
}

func createCertsSa(data *initData) error {
	return certsphase.CreateServiceAccountKeyAndPublicKeyFiles(data.CertificateWriteDir(), data.Cfg().ClusterConfiguration.PublicKeyAlgorithm())
}

func createCertAndKeyFilesWithCA(cert *certsphase.KubeadmCert, caCert *certsphase.KubeadmCert, data *initData) error {
	if certData, intermediates, err := pkiutil.TryLoadCertChainFromDisk(data.CertificateDir(), cert.BaseName); err == nil {
		certsphase.CheckCertificatePeriodValidity(cert.BaseName, certData)

		caCertData, err := pkiutil.TryLoadCertFromDisk(data.CertificateDir(), caCert.BaseName)
		if err != nil {
			return errors.Wrapf(err, "couldn't load CA certificate %s", caCert.Name)
		}

		certsphase.CheckCertificatePeriodValidity(caCert.BaseName, caCertData)

		if err := pkiutil.VerifyCertChain(certData, intermediates, caCertData); err != nil {
			return errors.Wrapf(err, "[certs] certificate %s not signed by CA certificate %s", cert.BaseName, caCert.BaseName)
		}

		fmt.Printf("[certs] Using existing %s certificate and key on disk\n", cert.BaseName)
		return nil
	}

	cfg := data.Cfg()
	cfg.CertificatesDir = data.CertificateWriteDir()
	defer func() { cfg.CertificatesDir = data.CertificateDir() }()

	return certsphase.CreateCertAndKeyFilesWithCA(cert, caCert, cfg)
}

func createCACertAndKeyFiles(ca *certsphase.KubeadmCert, data *initData) error {
	if cert, err := pkiutil.TryLoadCertFromDisk(data.CertificateDir(), ca.BaseName); err == nil {
		certsphase.CheckCertificatePeriodValidity(ca.BaseName, cert)

		if _, err := pkiutil.TryLoadKeyFromDisk(data.CertificateDir(), ca.BaseName); err == nil {
			fmt.Printf("[certs] Using existing %s certificate authority\n", ca.BaseName)
			return nil
		}
		fmt.Printf("[certs] Using existing %s keyless certificate authority\n", ca.BaseName)
		return nil
	}

	cfg := data.Cfg()
	cfg.CertificatesDir = data.CertificateWriteDir()
	defer func() { cfg.CertificatesDir = data.CertificateDir() }()

	return certsphase.CreateCACertAndKeyFiles(ca, cfg)
}
