package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	fmt.Println("🔍 Debugging Merkle Tree")
	fmt.Println("========================\n")

	// Create simple 2-tx tree
	txs := []*Transaction{
		{ID: "tx1", From: "Alice", To: "Bob", Amount: 10.0},
		{ID: "tx2", From: "Bob", To: "Charlie", Amount: 5.0},
	}

	tree, _ := NewMerkleTree(txs)

	fmt.Println("Step 1: Tree Structure")
	fmt.Println("----------------------")
	fmt.Printf("Root Hash: %s\n", hex.EncodeToString(tree.Root.Hash))
	fmt.Printf("Number of leaves: %d\n", len(tree.Leaves))
	fmt.Printf("Leaf 0 Hash: %s\n", hex.EncodeToString(tree.Leaves[0].Hash))
	fmt.Printf("Leaf 1 Hash: %s\n", hex.EncodeToString(tree.Leaves[1].Hash))

	fmt.Println("\nStep 2: Transaction Hash")
	fmt.Println("------------------------")
	tx0Hash := txs[0].Hash()
	fmt.Printf("TX0 String: %s\n", txs[0].String())
	fmt.Printf("TX0 Hash: %s\n", hex.EncodeToString(tx0Hash))
	fmt.Printf("Leaf 0 Hash (should match): %s\n", hex.EncodeToString(tree.Leaves[0].Hash))
	
	if hex.EncodeToString(tx0Hash) == hex.EncodeToString(tree.Leaves[0].Hash) {
		fmt.Println("✅ TX hash matches leaf hash!")
	} else {
		fmt.Println("❌ BUG: TX hash does NOT match leaf hash!")
	}

	fmt.Println("\nStep 3: Generate Proof for TX0")
	fmt.Println("-------------------------------")
	proof, err := tree.GenerateProof(0)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("Proof has %d hashes\n", len(proof.Hashes))
	for i, hash := range proof.Hashes {
		position := "left"
		if proof.Positions[i] {
			position = "right"
		}
		fmt.Printf("  Proof[%d] (%s): %s\n", i, position, hex.EncodeToString(hash)[:16]+"...")
	}

	fmt.Println("\nStep 4: Manual Verification")
	fmt.Println("----------------------------")
	
	// Manually verify the proof
	currentHash := tx0Hash
	fmt.Printf("Start with TX hash: %s\n", hex.EncodeToString(currentHash)[:16]+"...")

	for i, siblingHash := range proof.Hashes {
		fmt.Printf("\nLevel %d:\n", i+1)
		
		var combined []byte
		if proof.Positions[i] {
			fmt.Println("  Position: sibling on RIGHT")
			fmt.Printf("  Current:  %s\n", hex.EncodeToString(currentHash)[:16]+"...")
			fmt.Printf("  Sibling:  %s\n", hex.EncodeToString(siblingHash)[:16]+"...")
			combined = append(currentHash, siblingHash...)
		} else {
			fmt.Println("  Position: sibling on LEFT")
			fmt.Printf("  Sibling:  %s\n", hex.EncodeToString(siblingHash)[:16]+"...")
			fmt.Printf("  Current:  %s\n", hex.EncodeToString(currentHash)[:16]+"...")
			combined = append(siblingHash, currentHash...)
		}

		hash := sha256.Sum256(combined)
		currentHash = hash[:]
		fmt.Printf("  Result:   %s\n", hex.EncodeToString(currentHash)[:16]+"...")
	}

	fmt.Println("\nStep 5: Compare with Root")
	fmt.Println("-------------------------")
	fmt.Printf("Computed Root: %s\n", hex.EncodeToString(currentHash))
	fmt.Printf("Actual Root:   %s\n", hex.EncodeToString(tree.Root.Hash))
	
	if hex.EncodeToString(currentHash) == hex.EncodeToString(tree.Root.Hash) {
		fmt.Println("\n✅ MATCH! Manual verification works!")
	} else {
		fmt.Println("\n❌ MISMATCH! Bug in proof generation or tree building!")
	}

	fmt.Println("\nStep 6: Call VerifyProof Function")
	fmt.Println("----------------------------------")
	isValid := VerifyProof(tx0Hash, proof, tree.Root.Hash)
	fmt.Printf("VerifyProof returned: %v\n", isValid)
	
	if isValid {
		fmt.Println("✅ VerifyProof works correctly!")
	} else {
		fmt.Println("❌ Bug is in VerifyProof function!")
	}
}