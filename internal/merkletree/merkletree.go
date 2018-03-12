package merkletree

import (
	"crypto/sha256"
	"encoding/hex"
)

//MerkleTree struct hashes tX IDs, and calculates MerkleRoot.
type MerkleTree struct {
	Root           *string
	TransactionIDs []string
	layers         [][]string //illustrates leaves --> root
}

func (mt *MerkleTree) hashPair(tx1 string, tx2 string) string {
	shaEncoder := sha256.New()
	txConcat := tx1 + tx2
	shaEncoder.Write([]byte(txConcat))
	hash := hex.EncodeToString(shaEncoder.Sum(nil))
	return hash
}

//GetRoot recursively calculates Merkleroot from txIDs.
func (mt *MerkleTree) GetRoot() {
	var concat string
	// if no current transactions, set root to ''
	if len(mt.TransactionIDs) == 0 {
		root := ""
		mt.Root = &root
		return
	}

	if len(mt.TransactionIDs) == 1 {
		root := mt.TransactionIDs[0]
		mt.Root = &root
	} else {
		concat = mt.hashPair(mt.TransactionIDs[0], mt.TransactionIDs[1])
		mt.TransactionIDs = append([]string{concat}, mt.TransactionIDs[2:]...)
		mt.layers = append(mt.layers, mt.TransactionIDs)
		mt.GetRoot()
	}
}
