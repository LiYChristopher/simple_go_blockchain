package blockchain

import (
  "fmt"
)

//VizChain displays blockchain in the terminal
func (b *Blockchain) VizChain() {
  for _, b := range(b.Chain) {
    fmt.Printf("# BLOCK IDX: %v ##############\n", b.Index)
    fmt.Printf("#        Timestamp: %v       \n", b.Timestamp)
    fmt.Printf("#        Transactions: %v    \n", len(b.Transactions))
    fmt.Printf("#        Proof (nonce): %v   \n", b.Proof)
    fmt.Printf("#        PrevHash: %v        \n", *b.PrevHash)
    fmt.Printf("#        BlockHash: %v       \n", *b.BlockHash)
    fmt.Printf("######################### END#\n\n")
  }
}
