package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"go-blockchain/transaction"
	"strconv"
	"time"
)

type Block struct {
	PrevHash     []byte
	Hash         []byte
	Timestamp    int
	Nonce        int
	Transactions []transaction.Transaction
}

func convertIntToBytes(data int) []byte {
	return []byte(strconv.Itoa(data))
}

func (b Block) GetHash() []byte {
	return b.Hash
}

func (b Block) GetPrevHash() []byte {
	return b.PrevHash
}

func (b Block) GenerateHash() []byte {
	headers := bytes.Join([][]byte{b.PrevHash, convertIntToBytes(b.Nonce), convertIntToBytes(b.Timestamp)}, []byte{})

	hash := sha256.Sum256(headers)

	return hash[:]
}

func GenerateBlock(prevHash string) (b Block) {
	b.PrevHash = []byte(prevHash)
	b.Hash = b.GenerateHash()
	b.Nonce = 0
	b.Timestamp = int(time.Now().Unix())

	fmt.Println(b)

	return
}

func (b *Block) Mine(difficulty int) {
	target := bytes.Repeat([]byte("0"), difficulty)

	for ok := bytes.Equal(b.Hash[0:difficulty], target); !ok; {
		fmt.Println(string(b.Hash), target)
		b.Nonce++
		b.Hash = b.GenerateHash()
	}

	fmt.Println("Mined: " + string(b.Hash))
}
