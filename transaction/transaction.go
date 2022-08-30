package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-blockchain/keys"
	"log"
	"math/big"
)

type signature struct {
	r *big.Int
	s *big.Int
}

type Transaction struct {
	TransactionId string
	Value         float64
	Signature     signature
	Sender        keys.PublicKey
	Recipient     keys.PublicKey

	Inputs  []TransactionInput
	Outputs []TransactionOutput

	Sequence int
}

func (t *Transaction) GenerateSignature(privateKey keys.PrivateKey) {
	data := t.Sender.String() + t.Recipient.String() + fmt.Sprintf("%f", t.Value)

	r, s, err := ecdsa.Sign(rand.Reader, &ecdsa.PrivateKey{PublicKey: privateKey.PublicKey, D: privateKey.D}, []byte(data))

	if err != nil {
		log.Fatal(err)
	}

	t.Signature = signature{r, s}
}

func (t *Transaction) VerifySignature() bool {
	data := t.Sender.String() + t.Recipient.String() + fmt.Sprintf("%f", t.Value)

	return ecdsa.Verify(
		&ecdsa.PublicKey{Curve: t.Sender.Curve, X: t.Sender.X, Y: t.Sender.Y},
		[]byte(data),
		t.Signature.r,
		t.Signature.s,
	)
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

func (t *Transaction) CalculateHash() string {
	t.Sequence++

	h := sha256.New()

	headers := t.Sender.String() + t.Recipient.String() + fmt.Sprintf("%f", t.Value) + fmt.Sprintf("%d", t.Sequence)

	h.Write([]byte(headers))

	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}
