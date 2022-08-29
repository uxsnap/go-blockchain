package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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

	w.PrivateKey = keys.PrivateKey{PublicKey: pk.PublicKey, D: pk.D}
	w.PublicKey = w.PrivateKey.GetPublicKey()
}

func (w *Wallet) GetBalance() float64 {
	var total float64

	for id, UTXO := range w.UTXOs {
		if UTXO.IsMine(w.PublicKey) {
			w.UTXOs[id] = UTXO
			total += UTXO.Value
		}
	}

	return total
}

func (w *Wallet) SendFunds(to keys.PublicKey, value float64) transaction.Transaction {
	if w.GetBalance() < value {
		log.Println("#Not enough funds to send transaction. Transaction Discarded")
		return transaction.Transaction{}
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

	return newTransaction
}
