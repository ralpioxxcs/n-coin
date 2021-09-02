package blockchain

import (
	"fmt"
	"sync"

	"github.com/ralpioxxcs/n-coin/db"
	"github.com/ralpioxxcs/n-coin/utils"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

// AddBlock
func (b *blockchain) AddBlock(data string) {
	block := CreateBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	// slices of block pointer
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			// update hash from previous hash
			hashCursor = block.PrevHash
		} else {
			// when genesis block reached
			break
		}
	}
	return blocks
}

// Blockchain returns sigleton blockchain instance
func Blockchain() *blockchain {
	if b == nil {
		// do only one time
		once.Do(func() {
			b = &blockchain{"", 0}
			fmt.Printf("NewestHash : %s\nHeight:%d\n", b.NewestHash, b.Height)
			// search for checkpoint from db
			checkpoint := db.CheckPoint()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				// restore b from bytes
				fmt.Println("Restoring ...")
				b.restore(checkpoint)
			}

		})
	}
	fmt.Printf("NewestHash : %s\nHeight:%d\n", b.NewestHash, b.Height)
	return b
}
