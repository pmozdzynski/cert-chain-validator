# Certificate Validator

This repository contains few simple command-line tools written in Go to validate certificates and certificate chains.

## Tools

1. `validate.go`: Validates separate certificate and chain files.
2. `validate_bundle.go`: Validates a combined certificate and CA bundle in a single file.
3. `validate_remote.go`: Validates a combined certificate and the chain on remote server.
4. `cert_downloader.go`: Download the certificate and its chain from remote server.

## Usage

You can either run the tools using `go run` as described below, or compile them into binaries and use them directly.

### Compiling and Installing

To build all tools as standalone binaries:

```sh
# Build all tools
for f in validate.go validate_bundle.go validate_remote.go cert_downloader.go; do
    go build -o "${f%.go}" "$f"
done
```

You can then move the binaries to a directory in your `PATH` (e.g., `~/bin` or `/usr/local/bin`) to use them from anywhere:

```sh
mv validate validate_bundle validate_remote cert_downloader ~/bin/
```

Make sure `~/bin` is in your `PATH`:

```sh
export PATH="$HOME/bin:$PATH"
```

Now you can run the tools directly, e.g.:

```sh
validate cert.pem chain.pem
validate_bundle combined.pem
validate_remote foo.bar.com:443
cert_downloader foo.bar.com:443 server_cert.pem chain.pem
```

---

### validate.go

1. Clone the repository or copy the `validate.go` file to your local machine.
2. Run the program using the `go run` command with the following syntax, or use the compiled binary:
   ```sh
   go run validate.go cert.pem chain.pem
   # or, if compiled:
   ./validate cert.pem chain.pem
   ```
   - `cert.pem`: The path to the certificate file to be validated.
   - `chain.pem`: The path to the certificate chain file containing intermediate certificates.

### validate_bundle.go

1. Clone the repository or copy the `validate_bundle.go` file to your local machine.
2. Run the program using the `go run` command with the following syntax, or use the compiled binary:
   ```sh
   go run validate_bundle.go combined.pem
   # or, if compiled:
   ./validate_bundle combined.pem
   ```
   - `combined.pem`: The path to the file containing both the certificate and the CA bundle.

### validate_remote.go

1. Clone the repository or copy the `validate_remote.go` file to your local machine.
2. Run the program using the `go run` command with the following syntax, or use the compiled binary:
   ```sh
   go run validate_remote.go foo.bar.com:443
   # or, if compiled:
   ./validate_remote foo.bar.com:443
   ```
   - where foo.bar.com:443 - The address and port of the remote server from which the certificate should be fetched and validated.

#### Running validate_remote.go in Kubernetes Pod

If you need to run `validate_remote.go` inside a Kubernetes pod, follow these steps:

1. **Build the binary for Linux (`amd64`)**
   ```sh
   GOOS=linux GOARCH=amd64 go build -o validate_remote validate_remote.go
   ```

2. **Copy the binary to the Kubernetes Pod**
3. **Run the validation inside the pod**

### cert_downloader.go

1. Clone the repository or copy the `cert_downloader.go` file to your local machine.
2. Run the program using the `go run` command with the following syntax, or use the compiled binary:
   ```sh
   go run cert_downloader.go foo.bar.com:443 server_cert.pem chain.pem
   # or, if compiled:
   ./cert_downloader foo.bar.com:443 server_cert.pem chain.pem
   ```
   - `foo.bar.com:443` - The address and port of the remote server (e.g., example.com:443).
   - `cert.pem`: The file path where the server's certificate will be saved.
   - `chain.pem`: The file path where the certificate chain (intermediates and root) will be saved.

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

### cert_downloader.go
Downloads the server's certificate and its chain from a remote server and saves them as separate files.

### Dependencies
Go Standard Library (No external dependencies required)

