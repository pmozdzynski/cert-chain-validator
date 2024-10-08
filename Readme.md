# Certificate Validator

This repository contains two simple command-line tools written in Go to validate certificates and certificate chains. One tool validates separate certificate and chain files, while the other validates a combined certificate and CA bundle in a single file.

## Tools

1. `validate.go`: Validates separate certificate and chain files.
2. `validate_bundle.go`: Validates a combined certificate and CA bundle in a single file.

## Usage

### validate.go

1. Clone the repository or copy the `validate.go` file to your local machine.
2. Run the program using the `go run` command with the following syntax:
   ```sh
   go run validate.go cert.pem chain.pem
   ```
   - `cert.pem`: The path to the certificate file to be validated.
   - `chain.pem`: The path to the certificate chain file containing intermediate certificates.

### validate_bundle.go

1. Clone the repository or copy the `validate_bundle.go` file to your local machine.
2. Run the program using the `go run` command with the following syntax:
   ```sh
   go run validate_bundle.go combined.pem
   ```
   - `combined.pem`: The path to the file containing both the certificate and the CA bundle.

### validate_remote.go

1. Clone the repository or copy the `validate_remote.go` file to your local machine.
2. Run the program using the `go run` command with the following syntax:
   ```sh
   go run validate_remote.go foo.bar.com:443
   ```
   - where foo.bar.com:443 - The address and port of the remote server from which the certificate should be fetched and validated.


## Features

### validate.go

- Parses certificate files in PEM format.
- Validates the chain of trust by checking the signatures between the certificates.
- Outputs the validation status of each certificate in the chain.
- Checks the expiry dates of the certificates.

### validate_bundle.go

- Parses a combined certificate and CA bundle file in PEM format.
- Validates the chain of trust by checking the signatures between the certificates.
- Outputs the validation status of each certificate in the bundle.
- Checks the expiry dates of the certificates.

### validate_remote.go
Connects to a remote server over TLS and fetches the server's certificate chain.
Verifies the certificate chain by checking the signatures between certificates.
Outputs the validation status of each certificate in the chain.
Checks the expiry dates of the certificates.
Validates the server's certificate against the rest of the chain.

### Dependencies
Go Standard Library (No external dependencies required)
