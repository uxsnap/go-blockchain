package keys

import (
	"crypto/ecdsa"
)

type PrivateKey ecdsa.PrivateKey

func (pk PrivateKey) GetPublicKey() PublicKey {
	return PublicKey(pk.PublicKey)
}
