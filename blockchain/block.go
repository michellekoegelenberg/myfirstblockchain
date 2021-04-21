package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

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
	block := &Block{[]byte{}, []byte(data), prevHash, 0} //Block gets init
	pow := NewProof(block)                               //Init pow with NewProof allows us to pair a target with each new block that gets created
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// The Serialize meth gives us a bytes representation of our block
// Encode block into bytes

func (b *Block) Serialize() []byte {
	var res bytes.Buffer //Create result (of type bytes.Buffer)
	encoder := gob.NewEncoder(&res) //Create new encoder on results bytes buffer

	err := encoder.Encode(b) //Call 'Encode' on the block itself. Passes back an err. Need to do err handling

	Handle(err) //Handle err func added later

	return res.Bytes() // Return the bytes portion of our result. Gives us a bytes representation of our block
}



//Deserialize func is basically the opposite of serialize
func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
