# Basic Chain - In Go

Created by following this [great tutorial](https://hackernoon.com/learn-blockchains-by-building-one-117428612f46), and also inspired in part by this existing Go implementation https://github.com/crisadamo/gochain.

Completed for learning purposes. How this differs from the original:

- Proof of Work involves hashing a block's data (transaction data, previous hash) with the brute-force calculated proof, rather than the previous block's and current block's proofs.
- A merkletree is implemented to hash pending transactions, where the merkle root becomes part of block hash calculation.

Stil some refactoring to do at this point.

## Setup

Simply run

```
cd path/to/app
go run cmd/start_node.go --port=<port_num>
```

to start a node.
