package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	fmt.Println("🌳 Advanced Merkle Tree Features")
	fmt.Println("=================================\n")

	// Create transactions
	transactions := []*Transaction{
		{ID: "tx001", From: "Alice", To: "Bob", Amount: 100.50},
		{ID: "tx002", From: "Bob", To: "Charlie", Amount: 50.25},
		{ID: "tx003", From: "Charlie", To: "Dave", Amount: 75.00},
		{ID: "tx004", From: "Dave", To: "Eve", Amount: 25.75},
		{ID: "tx005", From: "Eve", To: "Frank", Amount: 150.00},
		{ID: "tx006", From: "Frank", To: "Grace", Amount: 80.50},
		{ID: "tx007", From: "Grace", To: "Henry", Amount: 45.25},
		{ID: "tx008", From: "Henry", To: "Alice", Amount: 200.00},
	}

	tree, _ := NewMerkleTree(transactions)

	// Feature 1: Tree Statistics
	fmt.Println("📊 Tree Statistics")
	fmt.Println("------------------")
	fmt.Printf("Total Transactions: %d\n", len(transactions))
	fmt.Printf("Root Hash: %s...\n", tree.GetRootHash()[:32])
	fmt.Printf("Number of Leaves: %d\n\n", len(tree.Leaves))

	// Feature 2: Verify Multiple Transactions
	fmt.Println("🔄 Batch Verification")
	fmt.Println("--------------------")
	indices := []int{0, 2, 4, 6}
	
	for _, idx := range indices {
		proof, _ := tree.GenerateProof(idx)
		txHash := transactions[idx].Hash()
		isValid := VerifyProof(txHash, proof, tree.Root.Hash)
		
		status := "✅"
		if !isValid {
			status = "❌"
		}
		fmt.Printf("Transaction #%d (%s): %s\n", idx+1, transactions[idx].ID, status)
	}

	// Feature 3: Export Proof Details
	fmt.Println("\n📄 Proof Details for TX #4")
	fmt.Println("-------------------------")
	proof, _ := tree.GenerateProof(3)
	
	for i, hash := range proof.Hashes {
		position := "left"
		if proof.Positions[i] {
			position = "right"
		}
		fmt.Printf("Level %d (%s): %s...\n", i+1, position, hex.EncodeToString(hash)[:16])
	}

	// Feature 4: Tree Visualization
	fmt.Println("\n🌳 Tree Structure")
	fmt.Println("-----------------")
	tree.PrintTree()

	fmt.Println("\n✨ Advanced demo complete!")
}
