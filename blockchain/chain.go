package blockchain

import (
	"bytes"
	"encoding/gob"
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
	decoder := gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(b)
}

func (b *blockchain) persist() {
	db.SaveBlockChain(utils.ToBytes(b))
}

// AddBlock
func (b *blockchain) AddBlock(data string) {
	block := CreateBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
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
