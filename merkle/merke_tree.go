package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	ID     string
	From   string
	To     string
	Amount float64
}

// String returns a string representation of the transaction
func (t *Transaction) String() string {
	return fmt.Sprintf("%s:%s->%s:%.2f", t.ID, t.From, t.To, t.Amount)
}

// Hash returns the SHA256 hash of the transaction
func (t *Transaction) Hash() []byte {
	data := []byte(t.String())
	hash := sha256.Sum256(data)
	return hash[:]
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

// MerkleTree represents the complete Merkle tree
type MerkleTree struct {
	Root         *MerkleNode
	Transactions []*Transaction
	Leaves       []*MerkleNode
}

// MerkleProof represents a proof that a transaction exists in the tree
type MerkleProof struct {
	Hashes    [][]byte
	Positions []bool // true = right, false = left
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
    node := &MerkleNode{}

    if left == nil && right == nil {
        // Leaf node - data is already hashed
        node.Hash = data  // ✅ FIX: Use hash directly!
    } else {

		// Internal node - hash the concatenation of children
		var prevHashes []byte
		if left != nil {
			prevHashes = append(prevHashes, left.Hash...)
		}
		if right != nil {
			prevHashes = append(prevHashes, right.Hash...)
		}
		hash := sha256.Sum256(prevHashes)
		node.Hash = hash[:]
	}

	node.Left = left
	node.Right = right

	return node
}

// NewMerkleTree creates a new Merkle tree from a list of transactions
func NewMerkleTree(transactions []*Transaction) (*MerkleTree, error) {
	if len(transactions) == 0 {
		return nil, errors.New("cannot create Merkle tree with no transactions")
	}

	var nodes []*MerkleNode
	var leaves []*MerkleNode

	// Create leaf nodes from transactions
	for _, tx := range transactions {
		node := NewMerkleNode(nil, nil, tx.Hash())
		nodes = append(nodes, node)
		leaves = append(leaves, node)
	}

	// If odd number of nodes, duplicate the last one
	if len(nodes)%2 != 0 {
		nodes = append(nodes, nodes[len(nodes)-1])
	}

	// Build the tree level by level
	for len(nodes) > 1 {
		var level []*MerkleNode

		for i := 0; i < len(nodes); i += 2 {
			node := NewMerkleNode(nodes[i], nodes[i+1], nil)
			level = append(level, node)
		}

		// If odd number of nodes at this level, duplicate the last one
		if len(level)%2 != 0 && len(level) > 1 {
			level = append(level, level[len(level)-1])
		}

		nodes = level
	}

	tree := &MerkleTree{
		Root:         nodes[0],
		Transactions: transactions,
		Leaves:       leaves,
	}

	return tree, nil
}

// GetRootHash returns the hex-encoded root hash
func (mt *MerkleTree) GetRootHash() string {
	if mt.Root == nil {
		return ""
	}
	return hex.EncodeToString(mt.Root.Hash)
}

// GenerateProof generates a Merkle proof for a transaction at given index
func (mt *MerkleTree) GenerateProof(txIndex int) (*MerkleProof, error) {
	if txIndex < 0 || txIndex >= len(mt.Transactions) {
		return nil, errors.New("transaction index out of bounds")
	}

	proof := &MerkleProof{
		Hashes:    [][]byte{},
		Positions: []bool{},
	}

	// Start from the leaf level
	nodes := make([]*MerkleNode, len(mt.Leaves))
	copy(nodes, mt.Leaves)
	currentIndex := txIndex

	// Handle odd number of leaves
	if len(nodes)%2 != 0 {
		nodes = append(nodes, nodes[len(nodes)-1])
	}

	// Traverse up the tree
	for len(nodes) > 1 {
		var level []*MerkleNode

		for i := 0; i < len(nodes); i += 2 {
			if i == currentIndex || i+1 == currentIndex {
				// Found our node - add sibling to proof
				if i == currentIndex {
					// Current node is left, add right sibling
					proof.Hashes = append(proof.Hashes, nodes[i+1].Hash)
					proof.Positions = append(proof.Positions, true) // sibling is on right
					currentIndex = i / 2
				} else {
					// Current node is right, add left sibling
					proof.Hashes = append(proof.Hashes, nodes[i].Hash)
					proof.Positions = append(proof.Positions, false) // sibling is on left
					currentIndex = i / 2
				}
			}

			// Build parent node
			node := NewMerkleNode(nodes[i], nodes[i+1], nil)
			level = append(level, node)
		}

		// Handle odd number at this level
		if len(level)%2 != 0 && len(level) > 1 {
			level = append(level, level[len(level)-1])
		}

		nodes = level
	}

	return proof, nil
}

// VerifyProof verifies a Merkle proof
func VerifyProof(txHash []byte, proof *MerkleProof, rootHash []byte) bool {
	currentHash := txHash

	for i, siblingHash := range proof.Hashes {
		var combined []byte

		if proof.Positions[i] {
			// Sibling is on the right
			combined = append(currentHash, siblingHash...)
		} else {
			// Sibling is on the left
			combined = append(siblingHash, currentHash...)
		}

		hash := sha256.Sum256(combined)
		currentHash = hash[:]
	}

	return hex.EncodeToString(currentHash) == hex.EncodeToString(rootHash)
}

// PrintTree prints the tree structure (for debugging)
func (mt *MerkleTree) PrintTree() {
	fmt.Println("Merkle Tree Structure:")
	fmt.Println("=====================")
	mt.printNode(mt.Root, "", true)
	fmt.Printf("\nRoot Hash: %s\n", mt.GetRootHash())
}

func (mt *MerkleTree) printNode(node *MerkleNode, prefix string, isLast bool) {
	if node == nil {
		return
	}

	fmt.Print(prefix)
	if isLast {
		fmt.Print("└── ")
		prefix += "    "
	} else {
		fmt.Print("├── ")
		prefix += "│   "
	}

	fmt.Printf("%s\n", hex.EncodeToString(node.Hash)[:16]+"...")

	if node.Left != nil || node.Right != nil {
		if node.Right != nil {
			mt.printNode(node.Right, prefix, false)
		}
		if node.Left != nil {
			mt.printNode(node.Left, prefix, true)
		}
	}
}