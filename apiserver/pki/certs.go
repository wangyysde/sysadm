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
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"k8s.io/apimachinery/pkg/util/validation"
	netutils "k8s.io/utils/net"
	"math"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"
)

// CreateCertificateAuthority creates new certificate and private key for the certificate authority.
// keyType is one of x509.ECDSA or x509.RSA.
// return *x509.Certificate, crypto.Signer, caCert in PEM, caKey in PEM and error
func CreateCertificateAuthority(commonName string, orgnaization []string, periodDays int,
	keyType x509.PublicKeyAlgorithm) (*x509.Certificate, crypto.Signer, []byte, []byte, error) {
	key, err := GeneratePrivateKey(keyType)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "unable to create private key while generating CA certificate")
	}

	cert, err := CreateSelfSignedCACert(commonName, orgnaization, periodDays, key)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "unable to create self-signed CA certificate")
	}

	keyPem, err := MarshalPrivateKeyToPEM(key)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "unable to marshal private key to pem")
	}

	certPem := EncodeCertPEM(cert)

	return cert, key, certPem, keyPem, nil
}

func GeneratePrivateKey(keyType x509.PublicKeyAlgorithm) (crypto.Signer, error) {
	if keyType == x509.ECDSA {
		return ecdsa.GenerateKey(elliptic.P256(), cryptorand.Reader)
	}

	return rsa.GenerateKey(cryptorand.Reader, rsaKeySize)
}

// CreateSelfSignedCACert creates a CA certificate
func CreateSelfSignedCACert(commonName string, orgnaization []string, periodDays int, key crypto.Signer) (*x509.Certificate, error) {
	now := time.Now()
	caDuration := time.Duration(defaultCaPeriodDays) * time.Hour * 24
	if periodDays > 0 {
		caDuration = time.Duration(periodDays) * time.Hour * 24
	}
	tmpl := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: orgnaization,
		},
		DNSNames:              []string{commonName},
		NotBefore:             now.UTC(),
		NotAfter:              now.Add(caDuration).UTC(),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

// MarshalPrivateKeyToPEM converts a known private key type of RSA or ECDSA to
// a PEM encoded block or returns an error.
func MarshalPrivateKeyToPEM(privateKey crypto.PrivateKey) ([]byte, error) {
	switch t := privateKey.(type) {
	case *ecdsa.PrivateKey:
		derBytes, err := x509.MarshalECPrivateKey(t)
		if err != nil {
			return nil, err
		}
		block := &pem.Block{
			Type:  ECPrivateKeyBlockType,
			Bytes: derBytes,
		}
		return pem.EncodeToMemory(block), nil
	case *rsa.PrivateKey:
		block := &pem.Block{
			Type:  RSAPrivateKeyBlockType,
			Bytes: x509.MarshalPKCS1PrivateKey(t),
		}
		return pem.EncodeToMemory(block), nil
	default:
		return nil, fmt.Errorf("private key is not a recognized type: %T", privateKey)
	}
}

// EncodeCertPEM returns PEM-endcoded certificate data
func EncodeCertPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}

// CreateCertAndKey creates new certificate and key by passing the certificate authority certificate and key
func CreateCertAndKey(keyType x509.PublicKeyAlgorithm, ca, caKey string, IPs []net.IP,
	periodDays int, commonName string, orgnaization, dnsNames []string, usages []x509.ExtKeyUsage) (*x509.Certificate, crypto.Signer, []byte, []byte, error) {
	key, err := GeneratePrivateKey(keyType)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "unable to create private key")
	}

	altNames := &certutil.AltNames{DNSNames: dnsNames, IPs: IPs}

	caCert, e := ParseCertPEM(ca)
	if e != nil {
		return nil, nil, nil, nil, errors.Wrap(e, "parse ca error")
	}
	keySigner, e := ParseKeyPEM(caKey)
	if e != nil {
		return nil, nil, nil, nil, errors.Wrap(e, "ca key is not valid")
	}

	cert, err := CreateSignedCert(altNames, periodDays, commonName, orgnaization, usages, key, caCert, keySigner)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "unable to sign certificate")
	}

	keyPem, err := MarshalPrivateKeyToPEM(key)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "unable to marshal private key to PEM")
	}
	certPem := EncodeCertPEM(cert)

	return cert, key, certPem, keyPem, nil

}

// CreateSignedCert creates a signed certificate using the given CA certificate and key
func CreateSignedCert(altNames *certutil.AltNames, periodDays int, commonName string, orgnaization []string,
	usages []x509.ExtKeyUsage, key crypto.Signer, caCert *x509.Certificate, caKey crypto.Signer) (*x509.Certificate, error) {
	serial, err := cryptorand.Int(cryptorand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	if len(commonName) == 0 {
		return nil, errors.New("must specify a CommonName")
	}

	keyUsage := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature

	RemoveDuplicateAltNames(altNames)

	certDuration := time.Duration(defaultCertPeriodDays) * time.Hour * 24
	if periodDays > 0 {
		certDuration = time.Duration(periodDays) * time.Hour * 24
	}
	notAfter := time.Now().Add(certDuration).UTC()

	certTmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: orgnaization,
		},
		DNSNames:              altNames.DNSNames,
		IPAddresses:           altNames.IPs,
		SerialNumber:          serial,
		NotBefore:             caCert.NotBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           usages,
		BasicConstraintsValid: true,
	}
	certDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

// RemoveDuplicateAltNames removes duplicate items in altNames.
func RemoveDuplicateAltNames(altNames *certutil.AltNames) {
	if altNames == nil {
		return
	}

	if altNames.DNSNames != nil {
		altNames.DNSNames = sets.NewString(altNames.DNSNames...).List()
	}

	ipsKeys := make(map[string]struct{})
	var ips []net.IP
	for _, one := range altNames.IPs {
		if _, ok := ipsKeys[one.String()]; !ok {
			ipsKeys[one.String()] = struct{}{}
			ips = append(ips, one)
		}
	}
	altNames.IPs = ips
}

// ValidateCertPeriod checks if the certificate is valid relative to the current time
// (+/- offset)
func ValidateCertPeriod(cert *x509.Certificate, offset time.Duration) error {
	period := fmt.Sprintf("NotBefore: %v, NotAfter: %v", cert.NotBefore, cert.NotAfter)
	now := time.Now().Add(offset)
	if now.Before(cert.NotBefore) {
		return errors.Errorf("the certificate is not valid yet: %s", period)
	}
	if now.After(cert.NotAfter) {
		return errors.Errorf("the certificate has expired: %s", period)
	}
	return nil
}

func ParseCertPEM(cert string) (*x509.Certificate, error) {
	cert = strings.TrimSpace(cert)
	if cert == "" {
		return nil, fmt.Errorf("content of certification in PEM is empty")
	}
	certs, err := certutil.ParseCertsPEM([]byte(cert))
	if err != nil {
		return nil, errors.Wrap(err, "can not parse certification")
	}

	return certs[0], nil
}

func ParseCertsChainPem(cert string) (*x509.Certificate, []*x509.Certificate, error) {
	cert = strings.TrimSpace(cert)
	if cert == "" {
		return nil, []*x509.Certificate{}, fmt.Errorf("content of certification in PEM is empty")
	}
	certs, err := certutil.ParseCertsPEM([]byte(cert))
	if err != nil {
		return nil, []*x509.Certificate{}, errors.Wrap(err, "can not parse certification")
	}

	parsedCert := certs[0]
	intermediates := certs[1:]

	return parsedCert, intermediates, nil
}

func ParseKeyPEM(keyStr string) (crypto.Signer, error) {
	keyStr = strings.TrimSpace(keyStr)
	if keyStr == "" {
		return nil, fmt.Errorf("content of key in PEM is empty")
	}

	privKey, err := keyutil.ParsePrivateKeyPEM([]byte(keyStr))
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse key")
	}

	// Allow RSA and ECDSA formats only
	var key crypto.Signer
	switch k := privKey.(type) {
	case *rsa.PrivateKey:
		key = k
	case *ecdsa.PrivateKey:
		key = k
	default:
		return nil, errors.Errorf("the private key file %s is neither in RSA nor ECDSA format")
	}

	return key, nil
}

// AppendSANsToAltNames parses SANs from as list of strings and adds them to altNames for use on a specific cert
// altNames is passed in with a pointer, and the struct is modified
// valid IP address strings are parsed and added to altNames.IPs as net.IP's
// RFC-1123 compliant DNS strings are added to altNames.DNSNames as strings
// RFC-1123 compliant wildcard DNS strings are added to altNames.DNSNames as strings
func AppendSANsToAltNames(altNames *certutil.AltNames, SANs []string) error {
	for _, altname := range SANs {
		if ip := netutils.ParseIPSloppy(altname); ip != nil {
			altNames.IPs = append(altNames.IPs, ip)
		} else if len(validation.IsDNS1123Subdomain(altname)) == 0 {
			altNames.DNSNames = append(altNames.DNSNames, altname)
		} else if len(validation.IsWildcardDNS1123Subdomain(altname)) == 0 {
			altNames.DNSNames = append(altNames.DNSNames, altname)
		} else {
			return fmt.Errorf(
				"'%s' was not added to the SAN, because it is not a valid IP or RFC-1123 compliant DNS entry\n",
				altname,
			)
		}
	}

	return nil
}
