package main

import (
    "crypto/sha256"
    "fmt"
    "time"
)

func main() {
    start := time.Now()

    for i := 0; i < 1_000_000; i++ {
        data := []byte(fmt.Sprintf("tx-%d", i))
        _ = sha256.Sum256(data)
    }

    fmt.Printf("Completed in %v\n", time.Since(start))
}
