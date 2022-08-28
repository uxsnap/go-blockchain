package blockchain

import (
	"bytes"
	"go-blockchain/block"
)

type Blockchain struct {
	Chain []block.Block
	// UTXOs map[string]TransactionOutput
	MinTransaction float64
	Difficulty     int
}

func CreateBlockchain() (blockchain Blockchain) {
	blockchain.Chain = []block.Block{}
	blockchain.Difficulty = 5
	blockchain.MinTransaction = 0.1

	return blockchain
}

func (b Blockchain) GetChain() []block.Block {
	return b.Chain
}

func (b Blockchain) IsChainValid() bool {
	chain := b.Chain

	if string(chain[0].GetPrevHash()) != "start" {
		return false
	}

	for i := 2; i < len(chain); i++ {
		if isValid(&b, i, i-1) {
			return false
		}
	}

	return true
}

func (b *Blockchain) AddBlock(newBlock block.Block) {
	newBlock.Mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)
}

func isValid(b *Blockchain, curBlockIndex, prevBlockIndex int) bool {
	curBlock := b.Chain[curBlockIndex]
	prevBlock := b.Chain[prevBlockIndex]

	// tempUTXOs := map[string]
	if !bytes.Equal(curBlock.GetHash(), curBlock.GenerateHash()) {
		return false
	}

	if !bytes.Equal(curBlock.GetPrevHash(), prevBlock.GetHash()) {
		return false
	}

	if target := bytes.Repeat([]byte{0}, b.Difficulty); !bytes.Equal(curBlock.GetHash()[0:b.Difficulty], target) {
		return false
	}

	return true
}
