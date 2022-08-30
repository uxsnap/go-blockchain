package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"go-blockchain/blockchain"
	"go-blockchain/keys"
	"go-blockchain/transaction"
	"log"
)

type Wallet struct {
	PublicKey  keys.PublicKey
	PrivateKey keys.PrivateKey
	UTXOs      map[string]transaction.TransactionOutput
}

func (w *Wallet) Create() {
	pk, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)

	if err != nil {
		log.Fatal(err)
	}

	w.UTXOs = make(map[string]transaction.TransactionOutput)
	w.PrivateKey = keys.PrivateKey{PublicKey: pk.PublicKey, D: pk.D}
	w.PublicKey = w.PrivateKey.GetPublicKey()
}

func (w *Wallet) SendFunds(b *blockchain.Blockchain, to keys.PublicKey, value float64) (transaction.Transaction, error) {
	if w.GetBalance(b) < value {
		return transaction.Transaction{}, errors.New("#Not enough funds to send transaction. Transaction Discarded")
	}

	inputs := []transaction.TransactionInput{}
	var total float64

	for id, UTXO := range w.UTXOs {
		total += UTXO.Value
		inputs = append(inputs, transaction.TransactionInput{TransactionOutputId: id})
		if total > value {
			break
		}
	}

	newTransaction := transaction.Transaction{
		Sender:    w.PublicKey,
		Recipient: to,
		Value:     value,
		Inputs:    inputs,
	}

	newTransaction.GenerateSignature(w.PrivateKey)

	for _, input := range inputs {
		delete(w.UTXOs, input.TransactionOutputId)
	}

	return newTransaction, nil
}

func (w *Wallet) GetBalance(b *blockchain.Blockchain) float64 {
	var total float64

	for id, UTXO := range b.UTXOs {
		if UTXO.IsMine(w.PublicKey) {
			w.UTXOs[id] = UTXO
			total += UTXO.Value
		}
	}

	return total
}
