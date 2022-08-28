package main

import (
	"go-blockchain/block"
	"go-blockchain/blockchain"
)

func main() {
	blockchain := blockchain.CreateBlockchain()

	genesis := block.GenerateBlock("root")

	// keys.Check()
	// keys.Check2()
}
