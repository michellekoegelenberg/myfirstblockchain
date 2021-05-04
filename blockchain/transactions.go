package blockchain

/* Because a blockchain is an open and public database, we don't want to store private/sens info inside of
the blockchain. No accounts, no balances, no addresses in our transactions.
Everything is done through the inputs and outputs and we can derive everything else
from the inputs and outputs. */

type Transactions struct {
	ID      []byte
	Outputs []TxOutput
	Input   []TxInput
}

type TxOutput struct {
	Value  int    // value in tokens which is assigned and locked inside this output
	PubKey string // Unlocks the (tokens inside the) Value field. Usually derived via script. Keeping it simple for now. Arb key to repres usr address
}

//Inputs are just references to prev outputs

type TxInput struct {
	ID  []byte // ID references the transaction the output is inside of
	Out int    // Index of the output (within the transaction)
	Sig string // Provides data used in output's pubkey ("Jack" unlock the output being referenced by the input)
}
