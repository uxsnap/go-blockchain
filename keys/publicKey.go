package keys

import (
	"crypto/ecdsa"
	"strconv"
)

type PublicKey ecdsa.PublicKey

func (pk *PublicKey) String() string {
	xString := strconv.Itoa(int(pk.X.Int64()))
	yString := strconv.Itoa(int(pk.Y.Int64()))

	return xString + yString
}

func (pk *PublicKey) Equal(otherKey PublicKey) bool {
	return pk.X.Cmp(otherKey.X) == 0 && pk.Y.Cmp(otherKey.Y) == 0 &&
		pk.Curve == otherKey.Curve
}
