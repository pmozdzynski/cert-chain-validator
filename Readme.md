# Certificate Validator

This program is a simple command-line tool written in Go to validate certificates and certificate chains. It reads certificate and certificate chain files, parses them, and then verifies the chain of trust by checking the signatures between the certificates in the chain.

## Usage

1. Clone the repository or copy the `main.go` file to your local machine.
2. Run the program using the `go run` command with the following syntax:
   ```sh
   go run validate.go cert.pem chain.pem
   ```
   - `cert.pem`: The path to the certificate file to be validated.
   - `chain.pem`: The path to the certificate chain file containing intermediate certificates.

## Features

- Parses certificate files in PEM format.
- Validates the chain of trust by checking the signatures between the certificates.
- Outputs the validation status of each certificate in the chain.

## Dependencies

- Go Standard Library (No external dependencies required)
