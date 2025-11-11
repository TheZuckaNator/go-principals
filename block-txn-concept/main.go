package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type TransactionType string

const (
	Credit TransactionType = "credit"
	Debit  TransactionType = "debit"
)

type Transaction struct {
	ID          int
	Hash        string
	From        string
	To          string
	Time        time.Time
	Description string
	Amount      float64
	Type        TransactionType
}

type Account struct {
	Address      string
	Owner        string
	Balance      float64
	Transactions []Transaction
}

func (a *Account) ApplyTransaction(t Transaction) error {
	switch t.Type {
	case Credit:
		a.Balance += t.Amount
	case Debit:
		if t.Amount > a.Balance {
			return fmt.Errorf("insufficient funds for tx %d", t.ID)
		}
		a.Balance -= t.Amount
	default:
		return fmt.Errorf("unknown transaction type: %s", t.Type)
	}

	a.Transactions = append(a.Transactions, t)
	return nil
}

func (a *Account) PrintStatement() {
	fmt.Printf("\n=== Account Statement =====================================\n")
	fmt.Printf("Owner   : %s\n", a.Owner)
	fmt.Printf("Address : %s\n\n", a.Address)

	for _, t := range a.Transactions {
		sign := "+"
		if t.Type == Debit {
			sign = "-"
		}
		fmt.Printf("Tx %d (%s)\n", t.ID, t.Hash[:16]+"...")
		fmt.Printf("  Time   : %s\n", t.Time.Format(time.RFC3339))
		fmt.Printf("  From   : %s\n", t.From)
		fmt.Printf("  To     : %s\n", t.To)
		fmt.Printf("  Type   : %s\n", t.Type)
		fmt.Printf("  Amount : %s%.2f\n", sign, t.Amount)
		fmt.Printf("  Note   : %s\n\n", t.Description)
	}

	fmt.Printf("Final balance: %.2f\n", a.Balance)
	fmt.Println("===========================================================\n")
}

// Block represents a simple block in the chain.
type Block struct {
	Index        int
	Timestamp    time.Time
	Nonce        uint64
	PrevHash     string
	Hash         string
	Transactions []Transaction
}

// computeTxHash returns a hash for a transaction (for display only).
func computeTxHash(t Transaction) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%d", t.ID)))
	h.Write([]byte(t.From))
	h.Write([]byte(t.To))
	h.Write([]byte(t.Time.Format(time.RFC3339Nano)))
	h.Write([]byte(t.Description))
	h.Write([]byte(fmt.Sprintf("%f", t.Amount)))
	h.Write([]byte(t.Type))
	return "0x" + hex.EncodeToString(h.Sum(nil))
}

// hashBlock computes the hash of the block based on:
// index, nonce, previous hash, timestamp, and tx hashes.
func hashBlock(b Block) string {
	h := sha256.New()

	// Order: Index -> Nonce -> PrevHash -> Timestamp -> Tx hashes
	h.Write([]byte(fmt.Sprintf("%d", b.Index)))
	h.Write([]byte(fmt.Sprintf("%d", b.Nonce)))
	h.Write([]byte(b.PrevHash))
	h.Write([]byte(b.Timestamp.Format(time.RFC3339Nano)))

	for _, tx := range b.Transactions {
		h.Write([]byte(tx.Hash))
	}

	return "0x" + hex.EncodeToString(h.Sum(nil))
}

// MineBlock finds a nonce such that the hash has `difficulty` leading zeros.
func MineBlock(b *Block, difficulty int) {
	target := "0x" + strings.Repeat("0", difficulty)

	for {
		hash := hashBlock(*b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			return
		}
		b.Nonce++
	}
}

func NewGenesisBlock(difficulty int) Block {
	b := Block{
		Index:        0,
		Timestamp:    time.Now(),
		Nonce:        0,
		PrevHash:     "0x0000000000000000000000000000000000000000000000000000000000000000",
		Transactions: nil,
	}
	MineBlock(&b, difficulty)
	return b
}

func NewBlock(prev Block, txs []Transaction, difficulty int) Block {
	b := Block{
		Index:        prev.Index + 1,
		Timestamp:    time.Now(),
		Nonce:        0,
		PrevHash:     prev.Hash,
		Transactions: txs,
	}
	MineBlock(&b, difficulty)
	return b
}

func printChain(chain []Block) {
	fmt.Println("=== Blockchain ============================================")
	for _, b := range chain {
		fmt.Printf("Block #%d\n", b.Index)
		fmt.Printf("  Timestamp : %s\n", b.Timestamp.Format(time.RFC3339))
		fmt.Printf("  Nonce     : %d\n", b.Nonce)
		fmt.Printf("  PrevHash  : %s\n", b.PrevHash[:20]+"...")
		fmt.Printf("  Hash      : %s\n", b.Hash[:20]+"...")
		fmt.Printf("  Tx count  : %d\n", len(b.Transactions))

		for _, tx := range b.Transactions {
			fmt.Printf("    - Tx %d: %s -> %s | %.2f (%s)\n",
				tx.ID,
				tx.From[:10]+"...",
				tx.To[:10]+"...",
				tx.Amount,
				tx.Type,
			)
		}
		fmt.Println()
	}
	fmt.Println("===========================================================")
}

func main() {
	account := &Account{
		Address: "0xAbC1234567890defABC1234567890defABC12345",
		Owner:   "Devon",
	}

	now := time.Now()

	// Example "addresses"
	alice := "0xA1cE000000000000000000000000000000000001"
	coffeeShop := "0xC0Ffee000000000000000000000000000000003"
	bookStore := "0xB00k000000000000000000000000000000000004"

	// Create raw txs
	rawTx1 := Transaction{
		ID:          1,
		From:        alice,
		To:          account.Address,
		Time:        now,
		Description: "Initial deposit",
		Amount:      1000.0,
		Type:        Credit,
	}
	rawTx2 := Transaction{
		ID:          2,
		From:        account.Address,
		To:          coffeeShop,
		Time:        now.Add(1 * time.Hour),
		Description: "Coffee",
		Amount:      4.50,
		Type:        Debit,
	}
	rawTx3 := Transaction{
		ID:          3,
		From:        account.Address,
		To:          bookStore,
		Time:        now.Add(2 * time.Hour),
		Description: "Book",
		Amount:      25.00,
		Type:        Debit,
	}

	// Compute tx hashes
	tx1 := rawTx1
	tx1.Hash = computeTxHash(rawTx1)

	tx2 := rawTx2
	tx2.Hash = computeTxHash(rawTx2)

	tx3 := rawTx3
	tx3.Hash = computeTxHash(rawTx3)

	difficulty := 3 // number of leading zeros required in hash

	genesis := NewGenesisBlock(difficulty)
	block1 := NewBlock(genesis, []Transaction{tx1, tx2}, difficulty)
	block2 := NewBlock(block1, []Transaction{tx3}, difficulty)

	chain := []Block{genesis, block1, block2}

	// Apply txs from blocks to account
	for _, b := range chain {
		for _, tx := range b.Transactions {
			if err := account.ApplyTransaction(tx); err != nil {
				fmt.Println("error applying tx:", err)
			}
		}
	}

	printChain(chain)
	account.PrintStatement()
}
