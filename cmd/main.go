package main

import (
  //"fmt"
  "../internal/blockchain"
)

func main() {

  //init Blockchain
  bc := blockchain.NewBlockchain()

  //create new TX (phase 1)
  bc.NewTransaction("Rick", "Morty", 69)
  bc.NewTransaction("Morty", "Rick", 33)
  bc.NewTransaction("Rick", "Morty", 36)
  bc.NewTransaction("Morty", "Mr. Meeseecks", 10)

  //Mine new block
  bc.Mine()

  //create new TX (phase 2)
  bc.NewTransaction("Morty", "Squanchy", 10)
  bc.NewTransaction("Squanchy", "Morty", 20)
  bc.NewTransaction("Morty", "Squanchy", 20)
  bc.NewTransaction("Squanchy", "Mr. Meeseecks", 3)

  //Mine new block
  bc.Mine()

  //create new TX (phase 3)
  bc.NewTransaction("Squanchy", "Rick", 5)
  bc.NewTransaction("Rick", "Morty", 30)
  bc.NewTransaction("Rick", "Squanchy", 5)
  bc.NewTransaction("Squanchy", "Mr. Meeseecks", 50)

  //Mine new block
  bc.Mine()

  //create new TX (phase 3)
  bc.NewTransaction("Pickle Rick", "Rick", 20)
  bc.NewTransaction("Jerry", "Pickle Rick", 1)
  bc.NewTransaction("Mr. Meeseecks", "Pickle Rick", 40)
  bc.NewTransaction("Pickle Rick", "Squanchy", 80)
  bc.NewTransaction("Gasorpasorp", "Morty", 3)
  bc.NewTransaction("Morty", "Squanchy", 3)

  //Mine new block
  bc.Mine()

  //visualize Blockchain
  bc.VizChain()
}
