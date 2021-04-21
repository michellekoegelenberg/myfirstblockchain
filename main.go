package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"github.com/michellekoegelenberg/myfirstblockchain/blockchain"
)

//Part 3 of Chapter 3:
//Create CLI to allow user to pass in a new block and print out the blockchain
//Step 1 is to create a CommandLine struct
type CommandLine struct {
	blockchain *blockchain.BlockChain
}

//Add instructions for the user
//method

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print – Prints the blocks in the chain")
}

//Create another method to validate arguments

func (cli *CommandLine) validateArgs() {
 if len(os.Args) < 2 {
	 cli.printUsage() //User hasn't entered in a command
	 runtime.Goexit() //Exits app by shutting down the go-routine (badger needs to properly garbage collect values and keys before shutdown)
 }
}

//Create AddBlock meth for user

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data) //Call to the bc inside our cli, call the AddBlock method on the data passed in from the cli
	fmt.Println("Added block!")
}

//Add method to print out blockchain
func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator() //Access our iterator by calling cli.blockchain.Iterator() which converst bc struct to iter struct
	//Create for loop and add print statements
	for {
		block := iter.Next()

		//Cut out print statements from below (main func) and paste them here
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate())) //print string format by converting bool to string
		fmt.Println()

		if len(block.PrevHash) == 0 { //Genesis block doesn't have prevHash, so it will be zero
			break
		}
	}
}

//Create run method to call all other methods
func (cli *CommandLine) run() {
	cli.validateArgs()
	//Create some flags for our cli tool to operate
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", " ", "Block data") //Add subset to add block cmd (Adding '-block' will alow them to add data string)

	//Get the args by calling a switch statement. Call on fist arg after the original call to the program.
	switch os.Args[1] {
		case "add":
			err := addBlockCmd.Parse(os.Args[2:]) //parse all the args that come after the first arg in our arg list
			blockchain.Handle(err) //call to blockchain library and calling handle meth we created

		case "print":
			err := printChainCmd.Parse(os.Args[2:]) //parse all the args that come after the first arg in our arg list
			blockchain.Handle(err)

		default: 
		cli.printUsage()
		runtime.Goexit() //To gracefully shut down the channel and then the program
	} 

	//When we parse flags: If they don't give error, they will give bolean value 'true'
	//Check this by calling parsed method on the flags

	if addBlockCmd.Parsed() { //Check to see if parsed
		if *addBlockData == "" { //Check to see if empty string
			cli.printUsage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData) //If not empty take data and make new block in our bc
	}

	if printChainCmd.Parsed() { //Check if parsed
		cli.printChain()

	}
}


//Part 3 of Chapter 3: Fix up main func
func main() {
	defer os.Exit(0) //Give time to garbage collect the keys and values
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close() //Properly close the database before the main func ends
	//Defer only executes if the go channel is able to exit properly.That's why we're using runtime goexits rather than os.exit
	//Add failsafe: Add defer os.Exit(0) at the top – done here

	//Create the cli struct by passing in our blockchain

	cli := CommandLine{chain}
	cli.run()
	}
