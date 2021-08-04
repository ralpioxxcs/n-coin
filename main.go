package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	data     string
	hash     string
	prevHash string
}

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

func main() {

	genesisBlock := block{"Gensis Block", "", ""}
	//genesisBlock.hash = sha256.Sum256([]byte(genesisBlock.data + genesisBlock.prevHash))

	hash := sha256.Sum256([]byte(genesisBlock.data + genesisBlock.prevHash))
	hexHash := fmt.Sprintf("%x", hash)
	genesisBlock.hash = hexHash

	fmt.Println(genesisBlock)

}
