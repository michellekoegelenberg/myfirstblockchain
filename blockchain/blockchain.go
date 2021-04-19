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

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions //Options struct
	opts.Dir = dbPath             //Where database will store keys
	opts.ValueDir = dbPath        //Where database will store values

	 //Could try this:   db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))

	db, err := badger.Open(opts) //Changed this in accordance with new BadgerDB API

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

			return (err)
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
   err := chain.Database.View(func(txn *badger.Txn)) error {
	//Get the current last hash out of the database (the hash of the last block in our database)
	//Call txn.Get to get the item, then unwrap the value from item and put it into our lastHash var
	//Return error if there is one
	//Handle the first err if there is one
	item, err := txn.Get(byte[]("lh"))
	Handle(err)
	lastHash, err = item.Value()
	return err
   }
	Handle(err)
	//Create a new block with the lastHash and the data we are passing in to the AddBlock func
   newBlock := CreateBlock(data, lastHash)
   //With our new block now created, we want to do a read/write type txn on our db
   //so we can put the new block into the database and assign the new block's hash to our last hash key
   err := chain.Database.Update(func(txn *badger.Txn) error) {
	   //As we did with Genesis (use new block's hash as key and serialize newBlock itself and put that in as the value)
	   err := txn.Set(newBlock.Hash, newBlock.Serialize())
	   Handle(err)
	   //Set newBlock's hash as last hash value
	   err = txn.Set([]byte("lh"), newBlock.Hash)
	   //Grab blockchain (chain) and lastHash field and set it equal to new block's hash
	   chain.lastHash = newBlock.Hash
	   return err

   }
	Handle(err)
//Succesfully created a layer of persistence for our blockchain
//However, lost the ability to go through our blockchain and print it out like we have before
//All the blocks are in the data layer. Can't just print them out
}
