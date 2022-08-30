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

	blockchain.AddBlock(genesis)

	block1 := block.GenerateBlock(genesis.Hash)

	log.Printf("\nWalletA balance is: %f\n", walletA.GetBalance(&blockchain))
	log.Println("\nWalletA is Attempting to send funds (40) to WalletB...")

	transaction, err := walletA.SendFunds(&blockchain, walletB.PublicKey, 40)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(transaction.Value, transaction.Outputs)

	blockchain.AddTransaction(&block1, transaction)

	// blockchain.AddBlock(block1)

	log.Printf("\nWalletA balance is: %f\n", walletA.GetBalance(&blockchain))
	log.Printf("\nWalletB balance is: %f\n", walletB.GetBalance(&blockchain))
}
