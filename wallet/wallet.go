package wallet

import (
	"go-blockchain/keys"
	"go-blockchain/transaction"
)

type Wallet struct {
	PublicKey  keys.PublicKey
	PrivateKey keys.PrivateKey
	UTXOs      map[string]transaction.TransactionOutput
}

func (w *Wallet) Create() {
	w.PrivateKey = keys.PrivateKey{}
	w.PublicKey = w.PrivateKey.GetPublicKey()
}

func (w *Wallet) GetBalance() {
	// var total float64

	// for id, UTXO := range w.UTXOs {
	// 	if UTXO.IsMine(ecdsa.PublicKey(w.PublicKey)) {
	// 		w.UTXOs[id] =
	// 	}
	// }
}
