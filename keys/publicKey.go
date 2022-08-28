package keys

import (
	"crypto/ecdsa"
	"strconv"
)

type PublicKey ecdsa.PublicKey

func (pk PublicKey) String() string {
	xString := strconv.Itoa(int(pk.X.Int64()))
	yString := strconv.Itoa(int(pk.Y.Int64()))

	return xString + yString
}
