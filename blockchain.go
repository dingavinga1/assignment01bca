package dvcoin

import (
	"crypto/sha256"
	"fmt"
	"time"
	"encoding/json"
	"strings"
)

/*
	CalculateHash
	@Param1 stringToHash "Input data for hashing function"
	
	Returns string "Sha256 Hash Value"
*/
func CalculateHash(stringToHash string) string{
	return fmt.Sprintf("%x", sha256.Sum256([]byte(stringToHash)));
}

/*
	Transaction
	# Structure to hold transactions

	@ Decorators(JSONable)
	
	- TransactionID "Random ID for each transaction"
	- SenderBlockchainAddress "Unique wallet address of sender"
	- RecipientBlockchainAddress "Unique wallet address of recipient"
	- Value "Amount transferred in transaction"
*/
type Transaction struct{
	TransactionID string `json:"TransactionID"`
	SenderBlockchainAddress string `json:"senderAddress"`
	RecipientBlockchainAddress string `json:"recipientAddress"`
	Value float32 `json:"value"`
}

/*
	GetString 
	-> Belongs to Transaction
	@@ImplicitParam trans "Data of Self (Transaction)"
	
	Returns string "Concatenated initial for creation of ID"
*/
func (trans Transaction) GetString() string{
	return fmt.Sprintf("%s%s%f", trans.SenderBlockchainAddress, trans.RecipientBlockchainAddress, trans.Value);
}

/*
	NewTransaction
	@Param1 sender "Address of sender"
	@Param2 recipient "Address of recipient"
	@Param3 value "Value send in transaction"
	
	Returns *Transaction "Reference to new transaction"
*/
func NewTransaction(sender string, recipient string, value float32) *Transaction{
	trans:=new(Transaction);
	trans.Value=value;
	trans.SenderBlockchainAddress=sender;
	trans.RecipientBlockchainAddress=recipient;

	timestamp := time.Now().Unix(); //getting current time for randomness

	trans.TransactionID=CalculateHash(trans.GetString()+fmt.Sprintf("%d", timestamp)); //generating unique ID for transaction

	return trans;
}

/*
	Block
	# Structure to hold blocks

	@ Decorators(JSONable)
	
	- CurrentHash "Hash of the current block" (Excluded from current hash)
	- Nonce "Final nonce for current block"
	- Timestamp "Timestamp for when this block was added to the chain" (Excluded from current hash)
	- Transactions "Array of transactions proposed in current block"
	- PrevPtr "Pointer to the previous block in the chain" (Excluded from current hash)
	- PreviousHash "Hash of the previous block"
*/
type Block struct{
	CurrentHash string `json:"currHash"` 
	Nonce int64 `json:"nonce"`
	Timestamp int64 `json:"timestamp"`

	Transactions []* Transaction `json:"transactions"`
	PrevPtr *Block `json:"-"` 
	PreviousHash string `json:"prevHash"`
}

/*
	GetString 
	-> Belongs to Block
	@@ImplicitParam trans "Data of Self (Block)"
	
	Returns string "Concatenated initial for mining"
*/
func (blk Block) GetString() string {
	strToHash:=blk.PreviousHash //previous hash inclusion
	for _, val:=range blk.Transactions{
		strToHash+=val.TransactionID //inclusion of unique IDs of each transaction
	}

	return strToHash
}

/*
	Print
	-> Belongs to Block
	@@ImplicitParam trans "Data of Self (Block)"
	
	Returns string "Neatly formatted JSON to display details of current block"
*/
func (blk Block) Print() {
	jsonData, err:=json.MarshalIndent(blk, "", "	") //converting details to JSON
	if err!=nil{ //error handling
		fmt.Println("Error converting data to JSON: ", err)
		return
	}

	fmt.Println(string(jsonData));
}

/*
	NewBlock
	@Param1 transactions "Array of transactions to propose in block"
	@Param2 prevPtr "Pointer to the previous block in the chain"
	@Param3 previousHash "Hash of the previous block"
	
	Returns *Block "Reference to new block"
*/
func NewBlock(transactions []*Transaction, prevPtr *Block, previousHash string) *Block{
	blk:=new(Block)
	blk.Transactions=transactions
	blk.PrevPtr=prevPtr
	blk.PreviousHash=previousHash

	return blk
}

/*
	Blockchain
	# Actual blockchain structure
	
	@ Decorators (JSONable)

	- Difficulty "Mining difficulty"
	- Chain "List of blocks in the chain"
	- TransactionPool "List of transactions waiting to be added to the chain"
*/
type Blockchain struct{
	Difficulty int `json:"difficulty"` 
	Chain []* Block `json:"chain"` 
	TransactionPool []* Transaction `json:"transPool"` 
}

/*
	Print
	-> Belongs to Blockchain
	@@ImplicitParam trans "Data of Self (Blockchain)"
	
	Returns string "Neatly formatted JSON to display entire blockchain"
*/
func (bc Blockchain) Print() {
	jsonData, err:=json.MarshalIndent(bc, "", "	") //converting details to JSON
	if err!=nil{ //error handling
		fmt.Println("Error converting data to JSON: ", err)
		return
	}

	fmt.Println(string(jsonData));
}

/*
	GiveBlock
	-> Belongs to Blockchain
	@@ImplicitParam bc "Data of Self (Blockchain)"
	@Param1 transactions "Transactions to be proposed"
	
	Returns *Block "Mining ready block"
*/
func (bc Blockchain) GiveBlock(transactions []*Transaction) *Block{
	var blk *Block
	if len(bc.Chain)!=0{ //checking if genesis block
		blk=NewBlock(transactions, bc.Chain[len(bc.Chain)-1], bc.Chain[len(bc.Chain)-1].CurrentHash)
	} else{
		blk=NewBlock(transactions, nil, "")
	}
	return blk;
}

/*
	AddTransaction
	-> Belongs to Blockchain
	@@ImplicitParam bc "Reference to Self (Blockchain)"
	@Param1 sender "Sender address"
	@Param2 recipient "Recipient address"
	@Param3 value "Amount transferred"

	"Adds a transaction to the pool"
*/
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32){
	trans:=NewTransaction(sender, recipient, value)
	bc.TransactionPool=append(bc.TransactionPool, trans);
}

/*
	MineBlock
	-> Belongs to Blockchain
	@@ImplicitParam bc "Data of Self (Blockchain)"
	@Param1 blk "Block to mine"

	Returns int64 "Nonce calculated according to specified difficulty level"
*/
func (bc Blockchain) MineBlock(blk *Block) int64{	
	var nonce int64=0
	starting:= strings.Repeat("0", bc.Difficulty) //pattern of initial zeroes according to diff
	initialString:=blk.GetString() //getting block data as string

	for{ //hit and trial for nonce
		hash:=CalculateHash(fmt.Sprintf("%s%d", initialString, nonce));
		if strings.HasPrefix(hash[:bc.Difficulty], starting){
			return nonce
		}
		nonce+=1
	}
}

/*
	VerifyNonce
	-> Belongs to Blockchain
	@@ImplicitParam bc "Data of Self (Blockchain)"
	@Param1 blk "Block to mine"
	@Param2 nonce "Nonce mined by user"

	Returns bool "Validity of nonce"
*/
func (bc Blockchain) VerifyNonce(blk *Block, nonce int64) bool{
	starting:= strings.Repeat("0", bc.Difficulty) //pattern of initial zeroes according to diff
	initialString:=blk.GetString() //getting block data as string

	hash:=CalculateHash(fmt.Sprintf("%s%d", initialString, nonce));

	if strings.HasPrefix(hash[:bc.Difficulty], starting){ //validating nonce mined
		return true
	}
	
	return false
}

/*
	AddBlock
	-> Belongs to Blockchain
	@@ImplicitParam bc "Data of Self (Blockchain)"
	@Param1 blk "Block to mine"
	@Param2 nonce "Nonce mined by user"

	Returns bool "Success of addition"
*/
func (bc *Blockchain) AddBlock(blk *Block, nonce int64) bool{
	if bc.VerifyNonce(blk, nonce){ //double checking for verifying nonce
		var temp []*Transaction 
		for _, val := range bc.TransactionPool { //removing these transactions from the transaction pool
			found:=false
			for _, val2:= range blk.Transactions {
				if val2.TransactionID==val.TransactionID{
					found=true
					break
				}
			}
			if found==false{ 
				temp=append(temp, val)
			}
		}

		bc.TransactionPool=temp //updating transaction pool

		blk.Timestamp=time.Now().Unix();
		blk.Nonce=nonce;
		blk.CurrentHash=CalculateHash(fmt.Sprintf("%s%d", blk.GetString(), nonce));

		bc.Chain=append(bc.Chain, blk)
		return true
	} else {
		return false
	}
}	

/*
	Node
	# Representing a node in the blockchain

	- Chain "Copy of the blockchain"
*/
type Node struct {
	Chain *Blockchain
}

/*
	ChooseTopTransactions
	-> Belongs to Node
	@@ImplicitParam bc "Data of Self (Node)"

	Returns []*Transaction "Top <=5 transactions to add to a block"
*/
func (node Node) ChooseTopTransactions() []*Transaction{
	l:=len(node.Chain.TransactionPool) 
	if l>5 { //greater between 5 and length of trans pool
		l=5
	}
	return node.Chain.TransactionPool[:l]
}

/*
	BuildBlock
	-> Belongs to Node
	@@ImplicitParam bc "Data of Self (Node)"

	Returns *Block "Mine ready block containing chosen transactions"
*/
func (node Node) BuildBlock(transactions []*Transaction) *Block{
	return node.Chain.GiveBlock(transactions);
}

/*
	Mine
	-> Belongs to Node
	@@ImplicitParam bc "Data of Self (Node)"

	Returns int64 "Proposed nonce"
*/
func (node *Node) Mine(blk *Block) int64 {
	return node.Chain.MineBlock(blk);
}
