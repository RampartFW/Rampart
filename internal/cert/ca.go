package cert

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// GenerateCA, cluster için bir kök CA oluşturur.
func GenerateCA(dir string) error {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Organization: []string{"Rampart Cluster Authority"},
			CommonName:   "Rampart Root CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, pub, priv)
	if err != nil {
		return err
	}

	if err := savePEM(filepath.Join(dir, "ca.crt"), "CERTIFICATE", certDER); err != nil {
		return err
	}
	
	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}
	return savePEM(filepath.Join(dir, "ca.key"), "PRIVATE KEY", privDER)
}

// GenerateNodeCert, belirli bir node için CA tarafından imzalanmış sertifika oluşturur.
func GenerateNodeCert(nodeName, caDir, outDir string, ips []net.IP) error {
	caCertPEM, err := os.ReadFile(filepath.Join(caDir, "ca.crt"))
	if err != nil {
		return err
	}
	caKeyPEM, err := os.ReadFile(filepath.Join(caDir, "ca.key"))
	if err != nil {
		return err
	}

	block, _ := pem.Decode(caCertPEM)
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	block, _ = pem.Decode(caKeyPEM)
	caPriv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Organization: []string{"Rampart Cluster"},
			CommonName:   nodeName,
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		IPAddresses:  ips,
		DNSNames:     []string{nodeName, "localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, pub, caPriv)
	if err != nil {
		return err
	}

	if err := savePEM(filepath.Join(outDir, nodeName+".crt"), "CERTIFICATE", certDER); err != nil {
		return err
	}
	
	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}
	return savePEM(filepath.Join(outDir, nodeName+".key"), "PRIVATE KEY", privDER)
}

func savePEM(path, typeStr string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return pem.Encode(f, &pem.Block{Type: typeStr, Bytes: data})
}
