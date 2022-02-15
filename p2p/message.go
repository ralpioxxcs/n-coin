package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/ralpioxxcs/n-coin/blockchain"
	"github.com/ralpioxxcs/n-coin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksReponse
	MessageNewBlockNotify
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJson(payload),
	}
	return utils.ToJson(m)
}

func sendNewestBlock(p *peer) {
	fmt.Printf("Sending newest block to %s\n", p.key)

	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)

	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksReponse, blockchain.Blocks(blockchain.Blockchain()))
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)

	p.inbox <- m
}

func handleMsg(m *Message, p *peer) {
	// fmt.Printf("Peer: %s, Sent a message with kind of: %d\n", p.key, m.Kind)
	switch m.Kind {
	case MessageNewestBlock:
		fmt.Printf("Received the newest block from %s\n", p.key)

		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)
		if payload.Height >= b.Height {
			fmt.Printf("Requesting all blocks from %s\n", p.key)
			requestAllBlocks(p) // request all the blocks from 4000
		} else {
			fmt.Printf("Sending newest block to %s\n", p.key)
			sendNewestBlock(p) // send 4000 our block
		}
	case MessageAllBlocksRequest:
		fmt.Printf("%s wants all the blocks\n", p.key)

		sendAllBlocks(p)
	case MessageAllBlocksReponse:
		fmt.Printf("Received all the blocks from %s\n", p.key)

		var payload []*blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))

		blockchain.Blockchain().Replace(payload)
	case MessageNewBlockNotify:

	}
}
