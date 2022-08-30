package transaction

import (
	"crypto/sha256"
	"encoding/hex"
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

	headers := recipient.String() + strconv.Itoa(int(value)) + parentTransactionId

	h := sha256.New()
	h.Write([]byte(headers))
	hashed := h.Sum(nil)

	to.Id = hex.EncodeToString(hashed)

	return
}

func (to *TransactionOutput) IsMine(publicKey keys.PublicKey) bool {
	return publicKey.Equal(to.Recipient)
}
