package transaction

import (
	"bytes"
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

func CreateTransactionOutput(recipient keys.PublicKey, value float64, parentTransactionId string) (to TransactionOutput) {
	to.Recipient = recipient
	to.Value = value
	to.ParentTransactionId = parentTransactionId

	id := bytes.Join([][]byte{[]byte(recipient.String()), []byte(strconv.Itoa(int(value))), []byte(parentTransactionId)}, []byte{})
	hash := sha256.Sum256(id)

	to.Id = string(hash[:])

	return
}

func (to *TransactionOutput) IsMine(publicKey keys.PublicKey) bool {
	return publicKey.Equal(to.Recipient)
}
