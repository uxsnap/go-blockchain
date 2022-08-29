package main

import (
	"go-blockchain/block"
	"go-blockchain/blockchain"
	"go-blockchain/transaction"
	"go-blockchain/wallet"
	"log"
)

func main() {
	blockchain := blockchain.CreateBlockchain()

	walletA := wallet.Wallet{}
	walletB := wallet.Wallet{}
	coinBase := wallet.Wallet{}

	walletA.Create()
	walletB.Create()
	coinBase.Create()

	genesisTransaction := transaction.Transaction{
		Sender:    coinBase.PublicKey,
		Recipient: walletA.PublicKey,
		Value:     100,
		Inputs:    []transaction.TransactionInput{},
	}

	genesisTransaction.GenerateSignature(coinBase.PrivateKey)
	genesisTransaction.TransactionId = "root"
	genesisTransaction.Outputs = append(genesisTransaction.Outputs, transaction.CreateTransactionOutput(
		genesisTransaction.Recipient,
		genesisTransaction.Value,
		genesisTransaction.TransactionId,
	))

	blockchain.UTXOs[genesisTransaction.Outputs[0].Id] = genesisTransaction.Outputs[0]

	log.Println("Creating and mining genesis block")

	genesis := block.GenerateBlock("root")

	blockchain.AddTransaction(&genesis, genesisTransaction)
	// blockchain.AddBlock(genesis)

	// log.Println(blockchain.GetChain())
}
