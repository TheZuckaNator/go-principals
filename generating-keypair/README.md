# Generating ECDSA Key Pairs in Go

A simple Go program that generates an Elliptic Curve (ECDSA) private and public key pair using the P-256 curve.

---

# Overview

## This project demonstrates how to:

1. Generate a cryptographic key pair with Go’s standard library

2. Marshal (convert) the private key into bytes

3. Display both private and public keys in a readable hex format

* It’s a minimal example for learning how Go handles key generation and the crypto and x509 packages.

    * When you generate a key pair using ecdsa.GenerateKey, you get a Go object — not something you can save to a file yet.

    * x509 lets you serialize that key into a standard format:

    ```go
    bytes, _ := x509.MarshalECPrivateKey(privateKey)
    ```

   * Bytes is the DER-encoded version of your private key — basically the binary form that other systems understand.

    If you wanted to save it as a .pem file, you’d wrap it like this:

    ```go
    pem.Encode(os.Stdout, &pem.Block{
        Type:  "EC PRIVATE KEY",
        Bytes: bytes,
    })
    ```


### How it Works

```go
privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
if err != nil {
    panic(err)
}
publicKey := &privateKey.PublicKey


ecdsa.GenerateKey() returns two values:
```

### The private key


The private key is then encoded with:

```bash
x509.MarshalECPrivateKey(privateKey)
```

and printed as hex bytes.

### Run It

Clone the repo and run:

```bash
go run main.go
```

You’ll see output like:

```yaml
Private Key (hex): 307702010104205f...
Public Key:
  X: 3c02f8a6d1f8b1...
  Y: 02a63f987f8e23...
```

### Files
File	Description
main.go	Contains the ECDSA key generation logic

### Notes

* Uses Go’s built-in crypto/ecdsa, crypto/elliptic, crypto/rand, and crypto/x509 packages.

* Curve used: P-256 (aka secp256r1).

You can export the keys to PEM files for real-world use (e.g., signing, wallets, JWTs).

