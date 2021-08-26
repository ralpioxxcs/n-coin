package blockchain

import (
	"crypto/sha256"
	"fmt"

	"github.com/ralpioxxcs/n-coin/db"
	"github.com/ralpioxxcs/n-coin/utils"
)

/*
Block1
	b1Hash = (data + "")

Block2
	b2Hash = (data + b1Hash)

Block3
	b3Hash = (data + b2Hash)

	.
	.

if Block1's hash is changed, block2 hash is changed,,
*/

// Block
type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func CreateBlock(data, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()

	return block
}
