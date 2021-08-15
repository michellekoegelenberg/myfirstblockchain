package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

/* Because a blockchain is an open and public database, we don't want to store private/sens info inside of
the blockchain. No accounts, no balances, no addresses in our transactions.
Everything is done through the inputs and outputs and we can derive everything else
from the inputs and outputs. */

type Transaction struct {
	ID      []byte
	Inputs  []TxInput  //array of inputs
	Outputs []TxOutput //array of outputs

}

type TxOutput struct {
	Value  int    // value in tokens which is assigned and locked inside this output
	PubKey string // Unlocks the (tokens inside the) Value field. Usually derived via script (lang). Keeping it simple for now. Arb key to repres user address
}

/*Outputs: Indivisible. Can't reference a part of an output.
If there are 10 tokens inside our output we need to create two new outputs,
one with 5 tokens inside and another with another 5.
*/

//Inputs are just references to prev outputs

type TxInput struct {
	ID  []byte // ID references the transaction the output is inside of
	Out int    // Index of the output (within the transaction)
	Sig string // Provides data used in output's pubkey ("Jack" unlock the output being referenced by the input)
}

/* In Genesis block we have our first transaction (Coinbase Transaction)
In this transaction: One input and one output.
Input references an empty output (no older outputs).
Doesn't store sig. Stores arb data.
Reward attached to it. Released to a single account when that individual mines the coinbase.
Just adding a const to make things simple for now
*/

//Create a Coinbase func. The goal of the f is to return a reference to a Transaction
//in order to do that it needs a 'to' account address and the 'data' we want to put in the coinbase

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data} //Empty []b{} because no ID. Out is -1 because it is referencing no output
	txout := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}} //nil for now (SetID)
	tx.SetID()

	return &tx
}

//Add NewT func after blockchain.go updates
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

//NewT func ends (now we need to create some commands inside our CLI)

//SetID Meth: Want to create a 32 b hash based on the bytes that represent our transaction
//So, we create a bytes.Buffer called encoded
// Then use the gob NewEncoder on the address of encoded and save it in 'encode'
//Then we use the encoder (encode) to encode (Encode) our tx (tx)
//Save in err
//Handle err
//Hash the Bytes portion of our encoded bytes (using sha and sum) on the bytes buffer which is called encoded
//Set that hash into the tx ID

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]

}

//Create util funcs to allow us to determine what kind of tx we have and
// whether or not we can unlock an output or an input that is ref an output

//Meth to help us determine whether or not a tx is a CoinBase tx
func (tx *Transaction) IsCoinbase() bool {
	//Coinbase has 1 input (length is 1), 1st input in tx Id is 0 (length is 0), inputs of out index is -1
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

//Create 2 methods we need to unlock the data inside of our outputs and inputs

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

/* If these come back as true it means the account (data) owns the information inside of the output
or it owns the info inside of the output that is referenced by the input */

/* Now that we have our transaction struct set up and we have the ability to
create a coinbase transaction, we need add txs to our blocks */

/*Edit Block Struct (block.go): Replace Data with an array of Txns */
