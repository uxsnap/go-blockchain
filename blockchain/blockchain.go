package blockchain

import (
	"log"
	"strings"

	"github.com/nuxxxcake/go-blockchain.git/block"
	"github.com/nuxxxcake/go-blockchain.git/transaction"
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

func (b *Blockchain) IsChainValid(genesisTransaction *transaction.Transaction) {
	for i := 1; i < len(b.Chain); i++ {
		if ok := isValid(b, genesisTransaction, i, i-1); !ok {
			log.Println("Blockchain is Invalid")
			return
		}
	}

	log.Println("Blockchain is Valid")
}

func (b *Blockchain) AddBlock(newBlock *block.Block) {
	newBlock.Mine(b.Difficulty)

	b.Chain = append(b.Chain, *newBlock)
}

func isValid(b *Blockchain, genesisTransaction *transaction.Transaction, curBlockIndex, prevBlockIndex int) bool {
	curBlock := b.Chain[curBlockIndex]
	prevBlock := b.Chain[prevBlockIndex]

	tempUTXOs := make(map[string]transaction.TransactionOutput)
	target := strings.Repeat("0", b.Difficulty)

	tempUTXOs[genesisTransaction.Outputs[0].Id] = genesisTransaction.Outputs[0]

	if curBlock.GetHash() != curBlock.GenerateHash() {
		log.Println("The hash of current block is invalid ")
		return false
	}

	if curBlock.GetPrevHash() != prevBlock.GetHash() {
		log.Println("The prevHash of current block and the hash of the prevBlock is not equal ")
		return false
	}

	if curBlock.GetHash()[0:b.Difficulty] != target {
		log.Println("The curBlock hash didn't meet the target requirements")
		return false
	}

	var tempOutput transaction.TransactionOutput

	for i, currentTranscation := range curBlock.Transactions {
		if !currentTranscation.VerifySignature() {
			log.Println("The transaction signature is not right")
			return false
		}

		if currentTranscation.GetInputsValue() != currentTranscation.GetOutputsValue() {
			log.Println("The inputs value is not equal the outputs value")
			return false
		}

		for _, input := range currentTranscation.Inputs {
			tempOutput = tempUTXOs[input.TransactionOutputId]

			if tempOutput == (transaction.TransactionOutput{}) {
				log.Printf("#Lost Transaction(%d)\n", i)
				return false
			}

			if input.UTXO.Value != tempOutput.Value {
				log.Printf("#Transaction(%d) value is not valid\n", i)
				return false
			}

			delete(tempUTXOs, input.TransactionOutputId)
		}

		for _, output := range currentTranscation.Outputs {
			tempUTXOs[output.Id] = output
		}

		if currentTranscation.Outputs[0].Recipient != currentTranscation.Recipient {
			log.Printf("#Transaction(%d) output recipient is not who it should be\n", i)
			return false
		}

		if currentTranscation.Outputs[1].Recipient != currentTranscation.Sender {
			log.Printf("#Transaction(%d) output 'change' is not sender.\n", i)
			return false
		}
	}

	return true
}

func (b *Blockchain) AddTransaction(bk *block.Block, t transaction.Transaction) bool {
	if t.Value == 0 {
		return false
	}

	if bk.GetPrevHash() != "root" {
		if !b.ProcessTransaction(&t) {
			log.Println("Transaction failed to process.")

			return false
		}
	}

	bk.Transactions = append(bk.Transactions, t)

	return true
}

func (b *Blockchain) ProcessTransaction(t *transaction.Transaction) bool {
	if !t.VerifySignature() {
		log.Println("#Transaction signature failed to verify")
		return false
	}

	for ind, input := range t.Inputs {
		if input == (transaction.TransactionInput{}) {
			continue
		}

		t.Inputs[ind].UTXO = b.UTXOs[input.TransactionOutputId]
	}

	if t.GetInputsValue() < b.MinTransaction {
		log.Printf("#Transaction Inputs too small: %f\n", t.GetInputsValue())
		return false
	}

	leftOver := t.GetInputsValue() - t.Value
	t.TransactionId = t.CalculateHash()

	t.Outputs = append(t.Outputs,
		transaction.CreateTransactionOutput(
			t.Recipient,
			t.Value,
			t.TransactionId,
		),
		transaction.CreateTransactionOutput(
			t.Sender,
			leftOver,
			t.TransactionId,
		),
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
