package transaction

import (
	"crypto/ecdsa"
	"math/big"
)

type Transaction struct {
	TransactionId string
	Value         float64

	Signature *big.Int
	Sender    ecdsa.PublicKey
	Recipient ecdsa.PublicKey

	Inputs  TransactionInput
	Outputs TransactionOutput
}
