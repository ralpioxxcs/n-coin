package blockchain

import (
	"fmt"
	"sync"

	"github.com/ralpioxxcs/n-coin/db"
	"github.com/ralpioxxcs/n-coin/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2 // minutes
	allowedRange       int = 2
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
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
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
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

func (b *blockchain) recalculateDifiiculty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastRecaulculatedBlock := allBlocks[difficultyInterval-1]

	// minutes
	actualTime := (newestBlock.Timestamp / 60) - (lastRecaulculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval

	//    8             10
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate the difficulty
		return b.recalculateDifiiculty()
	} else {
		return b.CurrentDifficulty
	}
}

// Blockchain returns sigleton blockchain instance
func Blockchain() *blockchain {
	if b == nil {
		// do only one time
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			// fmt.Printf("NewestHash : %s\nHeight:%d\n", b.NewestHash, b.Height)

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
