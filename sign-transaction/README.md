# ECDSA Transaction Signing in Go

This project demonstrates how to hash, sign, and verify blockchain-style transactions using ECDSA (Elliptic Curve Digital Signature Algorithm) in Go.

It also includes notes on a bug found in a code example from Go for Blockchain Development by José Gobert, which incorrectly handled ECDSA signatures.

## Overview

The code:

1. Defines a simple Transaction struct (From, To, Amount)

2. Hashes the transaction using SHA-256

3. Signs the hash using an ECDSA private key (P-256 curve)

4. Outputs the resulting ASN.1 encoded signature in hexadecimal

### Example Output
```bash
$ go run main.go
3045022100a9c8eac8a1f52d4f41...<snip>...b021b
```

That long hex string is your digital signature, encoded in ASN.1 DER format — the same encoding standard used by Bitcoin, Ethereum, and SSL/TLS.

###  How It Works

Hashing the transaction

```go
hash := sha256.Sum256([]byte(fmt.Sprintf("%s%s%f", tx.From, tx.To, tx.Amount)))
```

### Signing the hash

```go
sig, err := ecdsa.SignASN1(rand.Reader, priv, hash[:])
```

Output

### The signature bytes are printed as a hex string.

### Why ECDSA?

ECDSA (Elliptic Curve Digital Signature Algorithm) is the same cryptographic method used in:

Ethereum private/public keys

Bitcoin’s transaction signatures

Modern blockchain wallets and smart contract frameworks

# Bug Found in Go for Blockchain Development (José Gobert)

The book’s version of signTransaction had multiple issues:

## Original:

```go
func signTransaction(tx Transaction, priv *ecdsa.PrivateKey) []byte {
    hash := hashTransaction(tx)
    r, s, _ := ecdsa.Sign(rand.Reader, priv, hash)
    signature := append(r.Bytes(), s.Bytes() ...)
    return signature
}
```

##  Problems

| Issue | Description |
|--------|-------------|
| Syntax | `append(r.Bytes(), s.Bytes() ...)` → invalid syntax (extra space before `...`) |
| Encoding | Simply concatenating `r` and `s` (`r || s`) is *not reversible* — can’t safely extract them later |
| Error handling | Uses `_` to ignore errors from `ecdsa.Sign` — unsafe for cryptographic code |
| No verification | Without ASN.1 encoding, verifying the signature later is difficult |
| Non-standard | Real systems use **ASN.1 DER** (what `ecdsa.SignASN1` does automatically) |


## Fixed Version

```go
sig, err := ecdsa.SignASN1(rand.Reader, priv, hash)
if err != nil {
    panic(err)
}
fmt.Println(hex.EncodeToString(sig))
```

This correctly produces a portable, verifiable signature.

## Files
File	Description
main.go	Generates a key pair, signs a transaction, and prints the signature
sign.go (optional)	Contains reusable signTransaction and verifyTransaction helpers

### Dependencies

All dependencies are from Go’s standard library:

crypto/ecdsa

crypto/elliptic

crypto/rand

crypto/sha256

encoding/hex

fmt

No external packages are needed.

### Run It

```bash
go run main.go
```

You’ll get a valid ECDSA signature in hex format.

## Summary

* Correct use of ecdsa.SignASN1 for encoding

* Proper error handling

* Reproducible and verifiable results

--- 
# Fixed a real-world bug from Go for Blockchain Development (José Gobert)

