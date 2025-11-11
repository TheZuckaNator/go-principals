package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
)

func generateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey
}

func main() {
	priv, pub := generateKeys()

	// Marshal private key to DER bytes
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		panic(err)
	}

	fmt.Println("Private Key (hex):", hex.EncodeToString(privBytes))
	fmt.Printf("Public Key:\n  X: %x\n  Y: %x\n", pub.X, pub.Y)
}
