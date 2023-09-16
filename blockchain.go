package assignment01bca

import (
	"crypto/sha256"
	"fmt"
)

func CalculateHash(stringToHash string) string{
	return fmt.Sprintf("%x", sha256.Sum256([]byte(stringToHash)));
}

type Block struct{
	transaction string 
	nonce int
	prevPtr *Block
	previousHash string
	currentHash string 
}

func (blk *Block) Print() {
	fmt.Printf("-------------------------------------------\n")
	fmt.Printf("Transaction: %s\n", blk.transaction)
	fmt.Printf("Nonce: %d\n", blk.nonce)
	fmt.Printf("Previous hash: %s\n", blk.previousHash)
	fmt.Printf("Current hash: %s\n", blk.currentHash)
	fmt.Printf("-------------------------------------------\n")
}

func (blk *Block) Change(transaction string) {
	blk.transaction=transaction;
	blk.currentHash=CalculateHash(fmt.Sprintf("%s%d%s", blk.transaction, blk.nonce, blk.previousHash))
	
}

func NewBlock(transaction string, nonce int, prevPtr *Block, previousHash string) *Block{
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

type Blockchain struct{
	NewestBlock *Block
}

func (bc *Blockchain) AddBlock(transaction string, nonce int) *Block{
	if bc.NewestBlock==nil{
		bc.NewestBlock=NewBlock(transaction, nonce, nil, "")
	}

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