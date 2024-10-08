package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
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

// Saves a certificate to a file in PEM format
func saveCertToFile(cert *x509.Certificate, filename string) error {
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})

	err := ioutil.WriteFile(filename, certPEM, 0644)
	if err != nil {
		return fmt.Errorf("failed to write certificate to file: %v", err)
	}

	fmt.Printf("Certificate saved to: %s\n", filename)
	return nil
}

// Saves the certificate chain (excluding the first cert) to a file
func saveChainToFile(certs []*x509.Certificate, filename string) error {
	var chainPEM []byte
	for _, cert := range certs {
		chainPEM = append(chainPEM, pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})...)
	}

	err := ioutil.WriteFile(filename, chainPEM, 0644)
	if err != nil {
		return fmt.Errorf("failed to write certificate chain to file: %v", err)
	}

	fmt.Printf("Certificate chain saved to: %s\n", filename)
	return nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run cert_downloader.go foo.bar.com:443 server_cert.pem chain.pem")
		os.Exit(1)
	}

	remoteHost := os.Args[1]
	serverCertFile := os.Args[2]
	chainFile := os.Args[3]

	// Fetching the certificate from the remote server
	fmt.Println("Fetching certificates from the remote server:", remoteHost)
	remoteCerts, err := getCertFromRemoteServer(remoteHost)
	if err != nil {
		fmt.Printf("Failed to fetch certificates from remote server: %v\n", err)
		os.Exit(1)
	}

	if len(remoteCerts) == 0 {
		fmt.Println("No certificates found on the remote server.")
		os.Exit(1)
	}

	// Save the server's certificate (the first in the chain)
	fmt.Println("Saving the server certificate...")
	if err := saveCertToFile(remoteCerts[0], serverCertFile); err != nil {
		fmt.Printf("Error saving server certificate: %v\n", err)
		os.Exit(1)
	}

	// Save the rest of the certificates as the chain
	if len(remoteCerts) > 1 {
		fmt.Println("Saving the certificate chain...")
		if err := saveChainToFile(remoteCerts[1:], chainFile); err != nil {
			fmt.Printf("Error saving certificate chain: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("No certificate chain found on the remote server.")
	}
}
