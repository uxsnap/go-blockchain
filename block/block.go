package block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-blockchain/transaction"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	PrevHash     string
	Hash         string
	Timestamp    int
	Nonce        int
	Transactions []transaction.Transaction
}

func (b Block) GetHash() string {
	return b.Hash
}

func (b Block) GetPrevHash() string {
	return b.PrevHash
}

func (b Block) GenerateHash() string {
	headers := b.PrevHash + strconv.Itoa(b.Nonce) + strconv.Itoa(b.Timestamp)
	h := sha256.New()
	h.Write([]byte(headers))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

func GenerateBlock(prevHash string) (b Block) {
	b.PrevHash = prevHash
	b.Hash = b.GenerateHash()
	b.Nonce = 0
	b.Timestamp = int(time.Now().Unix())

	return
}

func (b *Block) Mine(difficulty int) {
	target := strings.Repeat("0", difficulty)

	for b.Hash[0:difficulty] != target {
		b.Nonce++
		b.Hash = b.GenerateHash()
	}

	fmt.Println("Mined: " + string(b.Hash))
}
