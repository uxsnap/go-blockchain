package transaction

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"go-blockchain/keys"
)

type Transaction struct {
	TransactionId string
	Value         float64

	Signature string
	Sender    keys.PublicKey
	Recipient keys.PublicKey

	Inputs  []TransactionInput
	Outputs []TransactionOutput

	Sequence int
}

func (t *Transaction) GenerateSignature(privateKey keys.PrivateKey) {
	data := t.Sender.String() + t.Recipient.String() + fmt.Sprintf("%f", t.Value)

	t.Signature = data
}

func (t *Transaction) VerifySignature() bool {
	data := t.Sender.String() + t.Recipient.String() + fmt.Sprintf("%f", t.Value)

	return ecdsa.VerifyASN1(&ecdsa.PublicKey{Curve: t.Sender.Curve, X: t.Sender.X, Y: t.Sender.Y}, []byte(data), []byte(t.Signature))
}

func (t *Transaction) GetInputsValue() float64 {
	var total float64

	for _, input := range t.Inputs {
		if input.UTXO == (TransactionOutput{}) {
			continue
		}

		total += input.UTXO.Value
	}

	return total
}

func (t *Transaction) GetOutputsValue() float64 {
	var total float64

	for _, output := range t.Outputs {
		total += output.Value
	}

	return total
}

func (t *Transaction) CalculateHash() []byte {
	t.Sequence++

	data := t.Sender.String() + t.Recipient.String() + fmt.Sprintf("%f", t.Value) + fmt.Sprintf("%d", t.Sequence)

	hash := sha256.Sum256([]byte(data))

	return hash[:]
}
