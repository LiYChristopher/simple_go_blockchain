package merkletree

import (
  "crypto/sha256"
  "encoding/hex"
)

//MerkleTree struct hashes tX IDs, and calculates MerkleRoot.
type MerkleTree struct {
  Root *string
  TXIDs []string
  layers [][]string //illustrates leaves --> root
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
	if len(mt.TXIDs) == 1 {
		root := mt.TXIDs[0]
		mt.Root = &root
	} else {
    //fmt.Printf("Concat Hashing Tx '%v' with '%v'\n", mt.TXIDs[0], mt.TXIDs[1])
  	concat = mt.hashPair(mt.TXIDs[0], mt.TXIDs[1])
  	mt.TXIDs = append([]string{concat}, mt.TXIDs[2:]...)
  	//fmt.Println("Completed 1 iteration of tree")
  	//fmt.Printf(" -- %v\n\n", mt.TXIDs)
    mt.layers = append(mt.layers, mt.TXIDs)
  	mt.GetRoot()
  }
}
