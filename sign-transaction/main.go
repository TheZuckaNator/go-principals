package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Transaction struct {
	From   string
	To     string
	Amount float64
}

func hashTransaction(tx Transaction) []byte {
	data := fmt.Sprintf("%s%s%f", tx.From, tx.To, tx.Amount)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

func signTransaction(tx Transaction, priv *ecdsa.PrivateKey) ([]byte, error) {
	hash := hashTransaction(tx)

	// ASN.1 encoded ECDSA signature (r,s) -> []byte
	sig, err := ecdsa.SignASN1(rand.Reader, priv, hash)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func main() {
	// generate a keypair
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	tx := Transaction{
		From:   "alice",
		To:     "bob",
		Amount: 42.0,
	}

	sig, err := signTransaction(tx, priv)
	if err != nil {
		panic(err)
	}

	// print ONLY the signature, hex encoded
	fmt.Println(hex.EncodeToString(sig))
}
