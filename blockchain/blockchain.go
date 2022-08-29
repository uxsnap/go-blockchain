package blockchain

import (
	"bytes"
	"go-blockchain/block"
	"go-blockchain/transaction"
	"log"
)

type Blockchain struct {
	Chain          []block.Block
	UTXOs          map[string]transaction.TransactionOutput
	MinTransaction float64
	Difficulty     int
}

func CreateBlockchain() (blockchain Blockchain) {
	blockchain.Chain = []block.Block{}
	blockchain.Difficulty = 5
	blockchain.MinTransaction = 0.1
	blockchain.UTXOs = make(map[string]transaction.TransactionOutput)

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

	tempUTXOs := make(map[string]transaction.TransactionOutput)

	if !bytes.Equal(curBlock.GetHash(), curBlock.GenerateHash()) {
		return false
	}

	if !bytes.Equal(curBlock.GetPrevHash(), prevBlock.GetHash()) {
		return false
	}

	if target := bytes.Repeat([]byte{0}, b.Difficulty); !bytes.Equal(curBlock.GetHash()[0:b.Difficulty], target) {
		return false
	}

	var tempOutput transaction.TransactionOutput

	for _, currentTranscation := range curBlock.Transactions {
		if !currentTranscation.VerifySignature() {
			return false
		}

		if currentTranscation.GetInputsValue() != currentTranscation.GetOutputsValue() {
			return false
		}

		for _, input := range currentTranscation.Inputs {
			tempOutput = tempUTXOs[input.TransactionOutputId]

			if tempOutput == (transaction.TransactionOutput{}) {
				return false
			}

			if input.UTXO.Value != tempOutput.Value {
				return false
			}

			delete(tempUTXOs, input.TransactionOutputId)
		}

		for i, output := range currentTranscation.Outputs {
			tempUTXOs[output.Id] = output

			if currentTranscation.Outputs[0].Recipient != currentTranscation.Recipient {
				log.Printf("#Transaction(%d) output recipient is not who it should be\n", i)
				return false
			}

			if currentTranscation.Outputs[1].Recipient != currentTranscation.Sender {
				log.Printf("#Transaction(%d) output 'change' is not sender.\n", i)
				return false
			}
		}
	}

	log.Println("Blockchain is Valid")

	return true
}

func (b *Blockchain) AddTransaction(bk *block.Block, t transaction.Transaction) bool {
	if t.TransactionId == "" {
		return false
	}

	if !bytes.Equal(bk.GetPrevHash(), []byte("root")) && !b.ProcessTransaction(t) {
		log.Println("Transaction failed to process.")

		return false
	}

	bk.Transactions = append(bk.Transactions, t)

	log.Println("Transaction successfuly added to Block")
	return true
}

func (b *Blockchain) ProcessTransaction(t transaction.Transaction) bool {
	if !t.VerifySignature() {
		log.Println("#Transaction signature failed to verify")
		return false
	}

	for _, input := range t.Inputs {
		if input == (transaction.TransactionInput{}) {
			continue
		}

		input.UTXO = b.UTXOs[input.TransactionOutputId]
	}

	if t.GetInputsValue() < b.MinTransaction {
		log.Printf("#Transaction Inputs too small: %f\n", t.GetInputsValue())
		return false
	}

	leftOver := t.GetInputsValue() - t.Value
	t.TransactionId = string(t.CalculateHash())

	t.Outputs = append(t.Outputs,
		transaction.TransactionOutput{
			Recipient:           t.Recipient,
			Value:               t.Value,
			ParentTransactionId: t.TransactionId,
		},
		transaction.TransactionOutput{
			Recipient:           t.Sender,
			Value:               leftOver,
			ParentTransactionId: t.TransactionId,
		},
	)

	for _, output := range t.Outputs {
		b.UTXOs[output.Id] = output
	}

	for _, input := range t.Inputs {
		if input.UTXO == (transaction.TransactionOutput{}) {
			continue
		}

		delete(b.UTXOs, input.UTXO.Id)
	}

	return true
}
