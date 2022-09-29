package main

import (
	"fmt"
	"goblockchain/blockchain"
	"time"
)

func main() {
	blockChain := blockchain.CreateBlockChain()
	time.Sleep(time.Second)
	blockChain.AddBlock("After genesis, I have something to say.")
	time.Sleep(time.Second)
	blockChain.AddBlock("sovo is awesome!")
	time.Sleep(time.Second)
	blockChain.AddBlock("I can't wait to follow his github!")
	time.Sleep(time.Second)

	for _, block := range blockChain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PreHash)
		fmt.Printf("Target: %x\n", block.Target)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("data: %s\n", block.Data)
	}
}
