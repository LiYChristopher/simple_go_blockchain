package blockchain

import (
  "fmt"
  "time"
  "crypto/sha256"
  "encoding/hex"
  "strconv"
  "os"
  "internal/merkletree"
)

//Blockchain is data-structure to contain all related operations.
type Blockchain struct {
  Chain []Block `json:"chain"`
  CurrentTX []Transaction `json:"-"`

}

//NewBlockChain creates a new chain with an auto-generated 'Genesis' block.
func NewBlockchain() *Blockchain {
  nbc := Blockchain{}

  // instantiate genesis block
  genPrevHash := "250000"
  genBlockHash := "250001"
  nbc.NewBlock(0, &genPrevHash, &genBlockHash)
  return &nbc
}

//NewBlock appends to the chain (ledger) once proof has been solved
func (b *Blockchain) NewBlock(proof int64, prevHash *string, blockHash *string) {
  nB := Block{
    Index: len(b.Chain) + 1,
    Timestamp: time.Now().Unix(),
    Transactions: b.CurrentTX,
    Proof: proof,
    PrevHash: prevHash,
    BlockHash: blockHash,
  }

  // New block appended - delete all transactions
  b.CurrentTX = make([]Transaction, 0)
  b.Chain = append(b.Chain, nB)
}

//NewTransaction creates a new transaction, adds to current transactions on Blockchain.
func (b *Blockchain) NewTransaction(sender string, rcpt string, amount float64) {
  //hash content to create TXID
  shaEncoder := sha256.New()
  concatData := sender + rcpt + strconv.FormatFloat(amount, 'f', 2, 64) + strconv.FormatInt(time.Now().Unix(), 10)
  shaEncoder.Write([]byte(concatData))
  nTID := hex.EncodeToString(shaEncoder.Sum(nil))
  nT := Transaction{ID: nTID, Sender: sender, Recipient: rcpt, Amount: amount}
  b.CurrentTX = append(b.CurrentTX, nT)
}

func (b *Blockchain) LastBlock() *Block {
  return &b.Chain[len(b.Chain)-1]
}

func (b *Blockchain) getTXIDs() []string {
  var txIDs []string
  for _, tx := range(b.CurrentTX) {
    txIDs = append(txIDs, tx.ID)
  }
  return txIDs
}

//HashTX - recursively hash transactions until one hash remains
func (b *Blockchain) HashTX() string {
  txIDs := b.getTXIDs()
  mt := merkletree.MerkleTree{TXIDs: txIDs}
  mt.GetRoot()
  return *mt.Root
}

//Mine performs sha256(timestamp, txHash, prevHash, proof) to get Block hash
func (b *Blockchain) Mine() {
  var proof int64
  shaEncoder := sha256.New()

  //concatenate block data (intermediate data - sha256)
  blockData := strconv.Itoa(int(time.Now().Unix()))
  blockData += b.HashTX()  //merkle root of transactions
  prevHash := *b.LastBlock().BlockHash
  blockData += prevHash
  shaEncoder.Write([]byte(blockData))
  intermediateData := hex.EncodeToString(shaEncoder.Sum(nil))

  //calculate proof of work
  proof, newBlockHash := b.proveAndHash(intermediateData)
  b.NewBlock(proof, &prevHash, &newBlockHash)
}

//proveAndHash performs POW, and then returns proof (nonce) and resulting hash
func (b *Blockchain) proveAndHash(intermediateData string) (proof int64, proofHash string) {
  var d string
  shaEncoder := sha256.New()

  for {
    //calculate proofHash
    d = intermediateData + strconv.Itoa(int(proof))
  	shaEncoder.Write([]byte(d))
  	proofHash = hex.EncodeToString(shaEncoder.Sum(nil))

    //check validity
    if b.isValidPoW(proofHash) {
      break
    }
    proof ++
  }
  return proof, proofHash
}

//isValidPoW performs simple bool check for proofHash.
func (b *Blockchain) isValidPoW(proofHash string) bool {
  if proofHash[:3] == os.Getenv("DIFFICULTY") {
    return true
  }
  return false
}

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

//Transaction is a datastructure for a single piece of data that may exist in a block.
type Transaction struct {
  ID string `json:",omitempty"`
  Sender string `json:"sender"`
  Recipient string `json:"recipient"`
  Amount float64 `json:"amount"`
}

//Block represents a unit that may be appended to a blockchain.
type Block struct {
  Index int `json:"index"`
  Timestamp int64 `json:"timestamp"`
  Transactions []Transaction `json:"transactions"`
  Proof int64 `json:"proof"`
  BlockHash *string `json:"hash"`
  PrevHash *string `json:"previous_hash"`
}
