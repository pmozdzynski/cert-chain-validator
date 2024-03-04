package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

func parseCerts(filename string) ([]*x509.Certificate, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var certs []*x509.Certificate
	for len(data) > 0 {
		var block *pem.Block
		block, data = pem.Decode(data)

		if block != nil {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate: %v", err)
			}
			certs = append(certs, cert)
		}
	}
	return certs, nil
}

func verifyChain(certs []*x509.Certificate) {
	for i := len(certs) - 1; i > 0; i-- {
		if err := certs[i-1].CheckSignatureFrom(certs[i]); err != nil {
			fmt.Printf("Certificate at position %d has abnormality. It's not signed by certificate at position %d. Error: %v\n", i+1, i, err)
		} else {
			fmt.Printf("Certificate at position %d is normal. It's signed by certificate at position %d.\n", i+1, i)
		}
	}
	fmt.Println("Certificate at position 1 is the Root Certificate.")
}

func validateCertAgainstChain(cert *x509.Certificate, chain []*x509.Certificate) error {
	roots := x509.NewCertPool()
	intermediates := x509.NewCertPool()
	for _, c := range chain {
		if c.IsCA {
			roots.AddCert(c)
		} else {
			intermediates.AddCert(c)
		}
	}

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
	}

	_, err := cert.Verify(opts)
	if err != nil {
		return fmt.Errorf("verification failed: %v", err)
	}

	fmt.Println("Certificate verification successful against the provided chain.")
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go cert.pem chain.pem")
		os.Exit(1)
	}

	certFile := os.Args[1]
	chainFile := os.Args[2]

	certs, err := parseCerts(certFile)
	if err != nil {
		fmt.Printf("Failed to parse certificate: %s - Error: %v\n", certFile, err)
		os.Exit(1)
	}

	chain, err := parseCerts(chainFile)
	if err != nil {
		fmt.Printf("Failed to parse certificate chain: %s - Error: %v\n", chainFile, err)
		os.Exit(1)
	}

	// Restores the detailed output for the chain's validation
	fmt.Println("Verifying certificate chain:")
	verifyChain(chain) // Note: this uses the original verifyChain logic

	// The added step for validating an end-entity certificate against the chain
	if len(certs) > 0 {
		fmt.Println("Validating certificate against the chain:")
		if err := validateCertAgainstChain(certs[0], chain); err != nil {
			fmt.Printf("Certificate validation against chain failed: %v\n", err)
		}
	} else {
		fmt.Println("No certificates found to validate.")
	}
}
