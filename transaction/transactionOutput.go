package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"go-blockchain/keys"
	"strconv"
)

type TransactionOutput struct {
	Id                  string
	Recipient           keys.PublicKey
	Value               float64
	ParentTransactionId string
}

func (to *TransactionOutput) Create(recipient keys.PublicKey, value float64, parentTransactionId string) {
	to.Recipient = recipient
	to.Value = value
	to.ParentTransactionId = parentTransactionId

	id := bytes.Join([][]byte{[]byte(recipient.String()), []byte(strconv.Itoa(int(value))), []byte(parentTransactionId)}, []byte{})
	hash := sha256.Sum256(id)

	to.Id = string(hash[:])
}

func (to *TransactionOutput) IsMine(publicKey ecdsa.PublicKey) bool {
	return publicKey.Equal(to.Recipient)
}
