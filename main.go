package main

import (
	"fmt"
	"sync"
	"time"

	dvcoin "github.com/dingavinga1/dvcoin"
)

func HoldCompetition(nodes []*dvcoin.Node, numNodes int){
	var mu sync.Mutex //Mutex for thread sync

	wg := new(sync.WaitGroup) //Wait group for joining
 
	nope:=false //flag for winner
	for i := 0; i < numNodes; i++ { //creating threads for each miner
		wg.Add(1) //adding new thread to workgroup
		
		go func(id int) { //function for participation in the mining competition
			defer wg.Done()
			time.Sleep(1 * time.Second) //waiting for all threads initial sync

			/* Mining */
			blk := nodes[id].BuildBlock(nodes[id].ChooseTopTransactions()) 
			nonce := nodes[id].Mine(blk)

			verified := nodes[id].Chain.VerifyNonce(blk, nonce) //system check for validating proposed nonce

			if verified { 
				mu.Lock() //Acquiring lock

				if nope{ //Loser check
					fmt.Printf("Sorry, Node#%d. You Lost!\n", id)
					mu.Unlock()
					return
				} else{ //Winner
					nope=true
					nodes[id].Chain.AddBlock(blk, nonce) //Adding block to the chain
					fmt.Printf("Congratulations, Node#%d! You won and your proposed block has been added.\n", id)
					mu.Unlock()
					return
				}	
			} else { //Check for invalid nonce
				fmt.Printf("Node#%d, you're a cheater.\n", id)
				return
			}
		}(i)
	}

	wg.Wait() //waiting for all threads to finish execution
}	

func main() {
	var bc dvcoin.Blockchain //blockchain
	bc.Difficulty = 2 //setting difficulty

	numNodes := 3 //number of total nodes holding copies of the blockchain

	var nodes []*dvcoin.Node
	for i := 0; i < numNodes; i++ { //creating 3 nodes and holding copy of chain
		nodes = append(nodes, &dvcoin.Node{})
		nodes[i].Chain = &bc
	}

	/* Providing 6 initial transactions in the transaction pool */
	bc.AddTransaction("Abdullah", "Aisha", 2)
	bc.AddTransaction("Huzaifa", "Aisha", 3)
	bc.AddTransaction("Aisha", "Huzaifa", 1)
	bc.AddTransaction("Aisha", "Abdullah", 4)
	bc.AddTransaction("Usman", "Abdullah", 8)
	bc.AddTransaction("Abdullah", "Usman", 5)

	fmt.Println("\t\t\t========================================")
	fmt.Println("\t\t\t           Before competition")
	fmt.Println("\t\t\t========================================")
	bc.Print()

	fmt.Println("\n\n\n\n\t\t\t========================================")
	fmt.Println("\t\t\t         After 1st competition")
	fmt.Println("\t\t\t========================================")
	HoldCompetition(nodes, numNodes);
	bc.Print()

	/* Providing 2 more transactions in the transaction pool */
	bc.AddTransaction("Usman", "Aisha", 10)
	bc.AddTransaction("Huzaifa", "Usman", 8)

	fmt.Println("\n\n\n\n\t\t\t========================================")
	fmt.Println("\t\t\t        After final competition")
	fmt.Println("\t\t\t========================================")
	HoldCompetition(nodes, numNodes);
	bc.Print()

}
