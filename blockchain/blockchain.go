package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

/*

[Blocks]

Block1
	b1Hash = (data + "")

Block2
	b2Hash = (data + b1Hash)

Block3
	b3Hash = (data + b2Hash)

	.
	.
	.

if Block1's hash is changed, block2 hash is changed,,

*/

type Block struct {
	Data     string
	Hash     string
	PrevHash string
}

type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash()}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

// GetBlockchain returns sigleton blockchain instance
func GetBlockchain() *blockchain {
	if b == nil {
		// do only one time
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis")
		})
	}
	return b
}

func (b *blockchain) AllBlocks() []*Block {
	return b.blocks
}
