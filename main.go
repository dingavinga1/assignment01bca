package main

import (
	"fmt"
	"github.com/dingavinga1/assignment01bca"
)

func main(){
	var bc assignment01bca.Blockchain;
	bc.NewestBlock=nil

	bc.AddBlock("I am Aisha", 0)
	bc.AddBlock("I am ok", 1)
	bc.AddBlock("I am Abdullah", 2)
	toBeChanged:=bc.AddBlock("I am ok", 3)
	toBeChanged.Print()
	bc.AddBlock("I am Huzaifa", 4)
	bc.AddBlock("I am ok", 5)
	bc.AddBlock("I am Usman", 6)
	bc.AddBlock("I am ok", 7)

	bc.ListBlocks()

	toBeChanged.Change("I am not ok")
	bc.ListBlocks()

	isOk:=bc.Verify()
	if isOk==false {
		fmt.Println("Blockchain has been compromised!")
	} else{
		fmt.Println("Blockchain is intact")
	}

}