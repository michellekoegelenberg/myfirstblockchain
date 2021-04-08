package blockchain

type BlockChain struct {
	Blocks []*Block
}

//Modify Block to store nonce so validation func can be implemented
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

//Add 0 for initial nonce (after prevHash)
//Modify CreateBlock so that it runs the pow algo on each block we create
//Execute the run func on that pow which will return the nonce and hash
//Put the nonce and hash into the block structure
//Return the block
//Now we can pass around data properly :)
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0} //Block gets initjsvdjksvbd
	pow := NewProof(block)                               //Init pow with NewProof allows us to pair a target with each new block that gets created
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
