package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("🌳 Merkle Tree - Basic Example")
	fmt.Println("===============================\n")

	// Create sample transactions
	transactions := []*Transaction{
		{ID: "tx1", From: "Alice", To: "Bob", Amount: 100.0},
		{ID: "tx2", From: "Bob", To: "Charlie", Amount: 50.0},
		{ID: "tx3", From: "Charlie", To: "Dave", Amount: 75.0},
		{ID: "tx4", From: "Dave", To: "Eve", Amount: 25.0},
	}

	fmt.Println("📝 Transactions:")
	for i, tx := range transactions {
		fmt.Printf("  %d. %s\n", i+1, tx.String())
	}
	fmt.Println()

	// Build Merkle tree
	tree, err := NewMerkleTree(transactions)
	if err != nil {
		log.Fatal("Error creating tree:", err)
	}

	// Display root hash
	rootHash := tree.GetRootHash()
	fmt.Printf("✅ Tree created successfully!\n")
	fmt.Printf("📊 Root Hash: %s...\n\n", rootHash[:32])

	// Generate and verify proof
	txIndex := 2
	fmt.Printf("🔐 Verifying Transaction #%d (%s)\n", txIndex+1, transactions[txIndex].ID)
	fmt.Println("-------------------------------------------")

	proof, err := tree.GenerateProof(txIndex)
	if err != nil {
		log.Fatal("Error generating proof:", err)
	}

	fmt.Printf("Proof generated: %d hashes\n", len(proof.Hashes))

	txHash := transactions[txIndex].Hash()
	isValid := VerifyProof(txHash, proof, tree.Root.Hash)

	if isValid {
		fmt.Println("✅ Proof is VALID - Transaction exists in the tree!")
	} else {
		fmt.Println("❌ Proof is INVALID")
	}

	// Tamper detection
	fmt.Println("\n🛡️  Tamper Detection Demo")
	fmt.Println("------------------------")

	originalTx := transactions[1]
	proof, _ = tree.GenerateProof(1)

	fmt.Printf("Original: %s\n", originalTx.String())
	isValid = VerifyProof(originalTx.Hash(), proof, tree.Root.Hash)
	fmt.Printf("  Verification: %v ✅\n\n", isValid)

	// Try tampering
	tamperedTx := &Transaction{
		ID:     originalTx.ID,
		From:   originalTx.From,
		To:     "Hacker",
		Amount: 999999.99,
	}

	fmt.Printf("Tampered: %s\n", tamperedTx.String())
	isValid = VerifyProof(tamperedTx.Hash(), proof, tree.Root.Hash)
	
	if !isValid {
		fmt.Println("  Verification: false ❌ (Correctly detected!)")
	} else {
		fmt.Println("  Verification: true (Should be false!)")
	}

	fmt.Println("\n✨ Demo complete!")
}
