package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// Fetches the certificates from the remote server
func getCertFromRemoteServer(hostname string) ([]*x509.Certificate, error) {
	conn, err := tls.Dial("tcp", hostname, &tls.Config{
		InsecureSkipVerify: true, // We're only interested in retrieving the cert chain, not verification here
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", hostname, err)
	}
	defer conn.Close()

	state := conn.ConnectionState()
	return state.PeerCertificates, nil
}

// Verifies the certificate chain presented by the remote server
func verifyChain(certs []*x509.Certificate) {
	for i := len(certs) - 1; i > 0; i-- {
		expiry := certs[i].NotAfter
		if err := certs[i-1].CheckSignatureFrom(certs[i]); err != nil {
			fmt.Printf("Certificate at position %d has abnormality. It's not signed by certificate at position %d. Error: %v - Expires: %v\n", i+1, i, err, expiry)
		} else {
			fmt.Printf("Certificate at position %d is normal. It's signed by certificate at position %d. Expires: %v\n", i+1, i, expiry)
		}
	}
	fmt.Printf("Certificate at position 1 is the Root Certificate. Expires: %v\n", certs[0].NotAfter)
}

// Validate a server certificate against its own chain
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
		fmt.Println("Usage: go run main.go foo.bar.com:443")
		os.Exit(1)
	}

	remoteHost := os.Args[1]

	// Fetching the certificate from the remote server
	fmt.Println("\nFetching certificate from the remote server:", remoteHost)
	remoteCerts, err := getCertFromRemoteServer(remoteHost)
	if err != nil {
		fmt.Printf("Failed to fetch certificate from remote server: %v\n", err)
		os.Exit(1)
	}

	// Display remote server certificate information
	fmt.Println("Remote Server Certificate Expiry Dates:")
	for i, cert := range remoteCerts {
		fmt.Printf("Certificate at position %d expires on: %v\n", i+1, cert.NotAfter)
	}

	// Verifying the server's certificate chain
	fmt.Println("\nVerifying remote server certificate chain:")
	verifyChain(remoteCerts)

	// Validate the first (leaf) certificate against the rest of the chain
	if len(remoteCerts) > 0 {
		fmt.Println("\nValidating remote server certificate against its own chain:")
		if err := validateCertAgainstChain(remoteCerts[0], remoteCerts); err != nil {
			fmt.Printf("Remote server certificate validation against chain failed: %v\n", err)
		}
	} else {
		fmt.Println("\nNo remote certificates found to validate.")
	}
}
