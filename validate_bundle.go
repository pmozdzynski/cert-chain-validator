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

func verifyChain(chain []*x509.Certificate) {
	for i := len(chain) - 1; i > 0; i-- {
		expiry := chain[i].NotAfter
		if err := chain[i-1].CheckSignatureFrom(chain[i]); err != nil {
			fmt.Printf("Certificate at position %d has abnormality. It's not signed by certificate at position %d. Error: %v - Expires: %v\n", i, i+1, err, expiry)
		} else {
			fmt.Printf("Certificate at position %d is normal. It's signed by certificate at position %d. Expires: %v\n", i, i+1, expiry)
		}
	}
	fmt.Printf("Certificate at position %d is the Root Certificate. Expires: %v\n", 1, chain[len(chain)-1].NotAfter)
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
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go combined.pem")
		os.Exit(1)
	}

	combinedFile := os.Args[1]

	certs, err := parseCerts(combinedFile)
	if err != nil {
		fmt.Printf("Failed to parse certificates: %s - Error: %v\n", combinedFile, err)
		os.Exit(1)
	}

	if len(certs) < 2 {
		fmt.Println("The combined file should contain at least one certificate and one CA.")
		os.Exit(1)
	}

	// Separate the end-entity certificate and the chain
	cert := certs[0]
	chain := certs[1:]

	fmt.Println("Certificate Expiry Dates:")
	fmt.Printf("End-Entity Certificate expires on: %v\n", cert.NotAfter)
	for i, cert := range chain {
		fmt.Printf("CA Certificate at position %d expires on: %v\n", i+1, cert.NotAfter)
	}

	// Verifying the chain
	fmt.Println("\nVerifying certificate chain:")
	verifyChain(chain)

	// Validating the end-entity certificate against the chain
	fmt.Println("\nValidating end-entity certificate against the chain:")
	if err := validateCertAgainstChain(cert, chain); err != nil {
		fmt.Printf("Certificate validation against chain failed: %v\n", err)
	}
}
