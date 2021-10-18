package blockchain

import (
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

// AddBlock
func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

// Blocks
func Blocks(b *blockchain) []*Block {
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

func recalculateDifiiculty(b *blockchain) int {
	allBlocks := Blocks(b)
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

func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate the difficulty
		return recalculateDifiiculty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func persistBlockchain(b *blockchain) {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	createTxs := make(map[string]bool)

	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					createTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, ok := createTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

// Blockchain returns sigleton blockchain instance
func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		checkpoint := db.CheckPoint()
		if checkpoint == nil {
			b.AddBlock()
		} else {
			b.restore(checkpoint)
		}
	})
	return b
}
