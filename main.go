package main

import (
	"fmt"
	"github.com/dingavinga1/assignment01bca"
)

func main(){
	var bc assignment01bca.Blockchain; // creating a blockchain object
	bc.NewestBlock=nil // initializing blockchain

	// adding various transactions
	bc.AddBlock("I am Aisha", 0)
	bc.AddBlock("I am ok", 1)
	bc.AddBlock("I am Abdullah", 2)
	toBeChanged:=bc.AddBlock("I am ok", 3)
	bc.AddBlock("I am Huzaifa", 4)
	bc.AddBlock("I am ok", 5)
	bc.AddBlock("I am Usman", 6)
	bc.AddBlock("I am ok", 7)

	bc.ListBlocks() //listing blockchain before change

	// verifying blockchain before change
	isOk:=bc.Verify()
	if isOk==false {
		fmt.Println("Blockchain has been compromised!")
	} else{
		fmt.Println("Blockchain is intact")
	}

	toBeChanged.Change("I am not ok") //changing a block
	bc.ListBlocks() // listing blockchain after change

	// verifying blockchain after change
	isOk=bc.Verify()
	if isOk==false {
		fmt.Println("Blockchain has been compromised!")
	} else{
		fmt.Println("Blockchain is intact")
	}

}