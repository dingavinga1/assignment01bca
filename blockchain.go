package assignment01bca

import (
	"crypto/sha256"
	"fmt"
)

func CalculateHash(stringToHash string) string{ // hashing function takes string as input and returns sha256 checksum
	return fmt.Sprintf("%x", sha256.Sum256([]byte(stringToHash)));
}

type Block struct{ //structure that defines a single block in the chain
	transaction string 
	nonce int
	prevPtr *Block // pointer to last block in the chain
	previousHash string // hash of the previous block
	currentHash string  // hash of transation+nonce+previousHash
}

func (blk *Block) Print() { // prints details for a single block in the chain
	fmt.Printf("-------------------------------------------\n")
	fmt.Printf("Transaction: %s\n", blk.transaction)
	fmt.Printf("Nonce: %d\n", blk.nonce)
	fmt.Printf("Previous hash: %s\n", blk.previousHash)
	fmt.Printf("Current hash: %s\n", blk.currentHash)
	fmt.Printf("-------------------------------------------\n")
}

// Made a member function for Block instead of ChangeBlock function
func (blk *Block) Change(transaction string) { // updates transaction within a single block to void integrity 
	blk.transaction=transaction;
	blk.currentHash=CalculateHash(fmt.Sprintf("%s%d%s", blk.transaction, blk.nonce, blk.previousHash))
	
}

func NewBlock(transaction string, nonce int, prevPtr *Block, previousHash string) *Block{ // allocating a new Block in memory and returns its pointer
	stream:=fmt.Sprintf("%s%d%s", transaction, nonce, previousHash);
	hash:=CalculateHash(stream);

	blk:=new(Block)
	blk.transaction=transaction
	blk.nonce=nonce
	blk.prevPtr=prevPtr
	blk.previousHash=previousHash
	blk.currentHash=hash

	return blk
}

type Blockchain struct{ // structure that defines a blockchain
	NewestBlock *Block // pointer to the latest block
}

func (bc *Blockchain) AddBlock(transaction string, nonce int) *Block{ // adds a new block to the chain and returns a pointer to that block
	if bc.NewestBlock==nil{ // check for adding first block
		bc.NewestBlock=NewBlock(transaction, nonce, nil, "")
		return bc.NewestBlock
	}
	// adding 
	blk:=NewBlock(transaction, nonce, bc.NewestBlock, bc.NewestBlock.currentHash)
	bc.NewestBlock=blk;

	return blk;
}

func (bc *Blockchain) ListBlocks(){
	iter:=bc.NewestBlock;
	for iter.prevPtr!=nil {
		iter.Print();
		iter=iter.prevPtr
	}
}

func (bc *Blockchain) Verify() bool{
	iter:=bc.NewestBlock;
	for iter.prevPtr!=nil {
		if iter.prevPtr.currentHash!=iter.previousHash {
			return false
		}
		iter=iter.prevPtr
	}
	return true
}