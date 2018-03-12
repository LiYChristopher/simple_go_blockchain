package blockchain

import (
  "fmt"
  "time"
  "crypto/sha256"
  "encoding/hex"
  "strconv"
  "os"
  "internal/merkletree"
  "errors"
)

//Transaction is a datastructure for a single piece of data that may exist in a block.
type Transaction struct {
  ID string `json:"id,omitempty"`
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

func getTXIDs(transactions []Transaction) []string {
  var txIDs []string
  for _, tx := range(transactions) {
    txIDs = append(txIDs, tx.ID)
  }
  return txIDs
}

//Blockchain is data-structure to contain all related operations.
type Blockchain struct {
  Chain []Block `json:"chain"`
  CurrentTX []Transaction `json:"-"`
  Nodes map[string]bool `json:"-"`
  Length int `json:"length"`
}

//NewBlockchain creates a new chain with an auto-generated 'Genesis' block.
func NewBlockchain() *Blockchain {
  nbc := Blockchain{}
  nbc.Nodes = make(map[string]bool, 0)

  // instantiate genesis block
  genPrevHash := "0"
  genBlockHash := "1"
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
  b.Length = len(b.Chain)
}

//NewTransaction creates a new transaction, adds to current transactions on Blockchain.
func (b *Blockchain) NewTransaction(sender string, rcpt string, amount float64) {
  //hash content to create TXID
  concatData := sender + rcpt + strconv.FormatFloat(amount, 'f', 2, 64)
  nTID := encodeSHA256(concatData)
  nT := Transaction{ID: nTID, Sender: sender, Recipient: rcpt, Amount: amount}
  b.CurrentTX = append(b.CurrentTX, nT)
}

func (b *Blockchain) LastBlock() *Block {
  return &b.Chain[len(b.Chain)-1]
}

//HashBlock hashes a block - transactions, previous hash, proof.
func (b *Blockchain) HashBlock(block *Block) string {
  hashedBlockData := b.hashBlockData(block)
  hashedBlockData += strconv.Itoa(int(block.Proof))
  return encodeSHA256(hashedBlockData)
}

//hashBlockData calculates the SHA-256 hash of a given block's transactions and previous hash.
//this is the intermediate step prior to calculation of the block's hash using proof.
func (b *Blockchain) hashBlockData(block *Block) string {
  blockData := b.HashTX(getTXIDs(block.Transactions))
  blockData += *block.PrevHash
  blockData = encodeSHA256(blockData)
  return blockData
}

//HashTX - recursively hash transactions until one hash remains
func (b *Blockchain) HashTX(txIDs []string) string {
  mt := merkletree.MerkleTree{TransactionIDs: txIDs}
  mt.GetRoot()
  return *mt.Root
}

//Mine performs sha256(timestamp, txHash, prevHash, proof) to get Block hash
func (b *Blockchain) Mine() {
  var proof int64

  //concatenate block data (intermediate data - sha256)
  blockData := b.HashTX(getTXIDs(b.CurrentTX))  //merkle root of transactions
  prevHash := *b.LastBlock().BlockHash
  blockData += prevHash
  pendingBlockData := encodeSHA256(blockData)

  //calculate proof of work
  proof, newBlockHash := b.proveAndHash(pendingBlockData)
  b.NewBlock(proof, &prevHash, &newBlockHash)
}

//proveAndHash performs POW, and then returns proof (nonce) and resulting hash
func (b *Blockchain) proveAndHash(blockData string) (proof int64, proofHash string) {

  for {
    //check validity
    if b.isValidPoW(blockData, proof) {
      proofHash = encodeSHA256(blockData + strconv.Itoa(int(proof)))
      break
    }
    proof ++
  }
  return proof, proofHash
}

//isValidPoW hashes block data and proof, returns bool if matches difficulty.
func (b *Blockchain) isValidPoW(data string, proof int64) bool {
  //calculate proofHash
  proofHash := encodeSHA256(data + strconv.Itoa(int(proof)))
  if proofHash[:3] == os.Getenv("DIFFICULTY") {
    return true
  }
  return false
}

//NewNode registers a new node to the blockchain.
func (b *Blockchain) NewNode(addr string) (err error) {
  if _, ok := b.Nodes[addr]; ok {
    err = errors.New("Node already registered")
  }
  b.Nodes[addr] = true
  return err
}

//isValidChain determines if the chain is valid - used in consensus.
func (b *Blockchain) isValidChain() bool {
  var lastBlock Block
  var hashedLastBlock string

  lastBlock = b.Chain[1]
  cur := 2

  for cur < len(b.Chain) {
      block := b.Chain[cur]
      fmt.Printf(" -- Currently iterating through block %v\n", *block.BlockHash)
      hashedLastBlock = b.HashBlock(&lastBlock)

      //Check that the hash of the block is correct
      if *block.PrevHash != hashedLastBlock  {
        return false
      }

      //since previous hash is encoded in current block data, isValidPoW will simply verify the
      //the proof still stands given the previous hash hasn't changed
      hashedCurBlockData := b.hashBlockData(&block)

      //Check that the Proof of Work is correct
      if !b.isValidPoW(hashedCurBlockData, block.Proof) {
        return false
      }

      lastBlock = block
      cur ++
    }
  return true
}

//resolveConflicts is a consensus algo. that will set the chain == to longest in network.
func (b *Blockchain) ResolveConflicts() (replaced bool) {
  var newChain *Blockchain
  var replacementNode string
  //update Length
  b.Length = len(b.Chain)

  neighbors := b.Nodes
  maxLength := b.Length

  for node, _ := range(neighbors) {
    fmt.Printf("Currently checking chain of node '%v' ... \n", node)
    nodeChain := getNodeChain(node)
    fmt.Println("Node downloaded.")

    if nodeChain.Length > maxLength && nodeChain.isValidChain() {
      fmt.Println("Nodechain is greater than max-length / is valid.")
      maxLength = nodeChain.Length
      newChain = nodeChain
      replacementNode = node
    }
  }

  if newChain != nil {
    b.Chain = newChain.Chain
    replaced = true
    fmt.Printf("Current blockchain replaced by chain from host '%v'\n", replacementNode)
  } else {
    replaced = false
  }
  return
}

//encodeSHA256 wraps encoding/sha256 to convert string to SHA-256 sum.
func encodeSHA256(text string) string {
  shaEncoder := sha256.New()
  shaEncoder.Write([]byte(text))
  return hex.EncodeToString(shaEncoder.Sum(nil))
}
