package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

const dbPath = "./tmp/blocks"

type BlockChain struct {
	LastHash []byte //Stores last hash of last block in chain
	Database *badger.DB
}

//Part 2 of Chapter 3: Create blockchain iterator (iterate through bc in database)
//Implemented at the bottom
type BlockChainIterator struct {
	CurrentHash []byte
	Database *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	/* This was the code in the tutorial
		Redundant due to BadgerDB new API
	
	opts := badger.DefaultOptions //Options struct
	opts.Dir = dbPath             //Where database will store keys
	opts.ValueDir = dbPath        //Where database will store values


	db, err := badger.Open(opts)  */


	db, err := badger.Open(badger.DefaultOptions("/tmp/badger")) //New BadgerDB API

  

	Handle(err)

	//Because we are initialising the blockchain, we need write capabilities: Hence, Update func
	err = db.Update(func(txn *badger.Txn) error {
		//Inside closure we have access to the transaction
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound { // Check if there is blockchain in db
			fmt.Println("No existing blockchain found")
			genesis := Genesis()                              //Creating genesis block
			fmt.Println("Genesis proved")                     //Proving it
			err := txn.Set(genesis.Hash, genesis.Serialize()) //Genesis hash is key for Genesis block, Serialize gen block
			//Put it into database using transaction set func
			Handle(err)
			//Gen block hash is the only hash in the database, so set it to the key of lh
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash //Save to var so we can put it into memory storage

			return err
		} else { //If we already have a db
			//Get the lh from the db
			item, err := txn.Get([]byte("lh")) //Call get on the key (lh). This gives pointer to item struct and err
			Handle(err)
			lastHash, err = item.ValueCopy(nil) //Changed this in accordance with new BadgerDB API
   //Get value from item struct & put into lastHash var
			return err                   //Return 1st err and handke 2nd err above
		}
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db} //Create new blockchain in memory
	return &blockchain                     //This way we can use it further in our app
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte
	//Execute a read-only type of txn on our BadgerDB (call the database from our blockchain by calling on chain.Database)
	//Call View (read-only), which takes in a closure with a pointer to a Badger transaction
	//Returns an err
   err := chain.Database.View(func(txn *badger.Txn) error {
	//Get the current last hash out of the database (the hash of the last block in our database)
	//Call txn.Get to get the item, then unwrap the value from item and put it into our lastHash var
	//Return error if there is one
	//Handle the first err if there is one
	item, err := txn.Get([]byte("lh"))
	Handle(err)
	//This code is updated code from Badger DB new API (and found in comments section of Part3 on YouTube Tut)
	
	err = item.Value(func(val []byte) error  {
   lastHash = append([]byte{}, val...)
  return nil
  //New code ends here
})
	return err
   }) 
   Handle(err)
	//Create a new block with the lastHash and the data we are passing in to the AddBlock func
   newBlock := CreateBlock(data, lastHash)
   //With our new block now created, we want to do a read/write type txn on our db
   //so we can put the new block into the database and assign the new block's hash to our last hash key
   err = chain.Database.Update(func(txn *badger.Txn) error {
	   //As we did with Genesis (use new block's hash as key and serialize newBlock itself and put that in as the value)
	   err := txn.Set(newBlock.Hash, newBlock.Serialize())
	   Handle(err)
	   //Set newBlock's hash as last hash value
	   err = txn.Set([]byte("lh"), newBlock.Hash)
	   //Grab blockchain (chain) and lastHash field and set it equal to new block's hash
	   chain.LastHash = newBlock.Hash
	   return err

   })
	Handle(err)
//Succesfully created a layer of persistence for our blockchain
//However, lost the ability to go through our blockchain and print it out like we have before
//All the blocks are in the data layer. Can't just print them out
}

//Part 2 of Chapter 3: 
//Implementing Iterator: 
//Convert Blockchain struct into blockchain iterator struct
//Create method called Iterator that has a receiver for pointer to blockchain
//returns a pointer to a BlockChainIterator
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database} //These will become the BCI's current hash and database
	return iter
}

//Because we are starting with the last hash of our bc we will be iterating backwards through the blocks
	//starting with the newest and working our way back to the genesis block
	//Create next func for this

	func (iter *BlockChainIterator) Next() *Block {
		var block *Block
		//Read only transaction
		err := iter.Database.View(func(txn *badger.Txn) error {
			item, err := txn.Get(iter.CurrentHash)
			Handle(err)
			//New code from updated DB
			var encodedBlock []byte
  			err = item.Value(func(val []byte) error { //Changed := to =
   			encodedBlock = append([]byte{}, val...)
    		return nil
  			})
			  //New code ends

			block = Deserialize(encodedBlock)
	
			return err
		})
		Handle(err)

		iter.CurrentHash = block.PrevHash //Change iter current hash to the block's prev hash
		
		return block

	//Part 3 of Chapter 3:
	//Create CLI to allow user to pass in a new block and print out the blockchain
	//Step 1 is to create a CommandLine struct in main.go
	}