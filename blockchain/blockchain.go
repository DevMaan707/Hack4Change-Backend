package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"time"
)

type Block struct {
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	Data         string
}

type Blockchain struct {
	Blocks []Block
}

func NewBlock(data string, previousHash string) Block {
	block := Block{
		Hash:         "",
		PreviousHash: previousHash,
		Timestamp:    time.Now(),
		Data:         data,
	}
	block.Hash = block.calculateHash()
	return block
}

func (b *Block) calculateHash() string {
	record := b.PreviousHash + b.Timestamp.String() + b.Data
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]Block{NewBlock("Genesis Block", "")}}
}

func (bc *Blockchain) AddBlock(data string) {
	previousBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, previousBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Fatal(err)
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Fatal(err)
	}
	return &block
}
