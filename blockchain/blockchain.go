package blockchain

import (
	"github.com/dgraph-io/badger/v3"
)

const dbPath = "./tmp/blocks"

type BlockChain struct {
	LastHash []byte //Stores last hash of last block in chain
	Database *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions //Options struct
	opts.Dir = dbPath             //Where database will store keys
	opts.ValueDir = dbPath        //Where database will store values
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}
