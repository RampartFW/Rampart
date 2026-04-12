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

// GenerateNodeCert generates a new node certificate signed by the CA.
func GenerateNodeCert(nodeName, caDir string) error {
	caKeyPath := filepath.Join(caDir, "ca.key")
	caCertPath := filepath.Join(caDir, "ca.crt")

	// Load CA certificate
	caCertBytes, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %w", err)
	}
	block, _ := pem.Decode(caCertBytes)
	if block == nil {
		return fmt.Errorf("failed to decode CA certificate")
	}
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Load CA private key
	caKeyBytes, err := os.ReadFile(caKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read CA key: %w", err)
	}
	block, _ = pem.Decode(caKeyBytes)
	if block == nil {
		return fmt.Errorf("failed to decode CA key")
	}
	caKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA key: %w", err)
	}

	// Generate node private key
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate node Ed25519 key: %w", err)
	}

	// Create node certificate template
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Rampart"},
			CommonName:   nodeName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // 1 year
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// Add node name to SAN
	template.DNSNames = append(template.DNSNames, nodeName, "localhost")
	template.IPAddresses = append(template.IPAddresses, net.ParseIP("127.0.0.1"), net.IPv6loopback)

	// If nodeName is an IP address, add it as an IP SAN
	if ip := net.ParseIP(nodeName); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	}

	// Sign node certificate with CA
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, pub, caKey)
	if err != nil {
		return fmt.Errorf("failed to create node certificate: %w", err)
	}

	// Save node key
	nodeKeyPath := filepath.Join(caDir, nodeName+".key")
	keyBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}
	keyOut, err := os.OpenFile(nodeKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open node key file for writing: %w", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes}); err != nil {
		return fmt.Errorf("failed to encode node key: %w", err)
	}
	keyOut.Close()

	// Save node certificate
	nodeCertPath := filepath.Join(caDir, nodeName+".crt")
	certOut, err := os.Create(nodeCertPath)
	if err != nil {
		return fmt.Errorf("failed to open node certificate file for writing: %w", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to encode node certificate: %w", err)
	}
	certOut.Close()

	return nil
}

// VerifyCertificate verifies a certificate against the CA.
func VerifyCertificate(certPath, caCertPath string) error {
	caCertBytes, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %w", err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCertBytes) {
		return fmt.Errorf("failed to append CA certificate to pool")
	}

	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("failed to read certificate: %w", err)
	}
	block, _ := pem.Decode(certBytes)
	if block == nil {
		return fmt.Errorf("failed to decode certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	opts := x509.VerifyOptions{
		Roots: caPool,
	}

	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify certificate: %w", err)
	}

	return nil
}
