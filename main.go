package main

import (
	"fmt"
	"goblockchain/blockchain"
	"goblockchain/transaction"
)

func main() {
	// 创建交易池
	txPool := make([]*transaction.Transaction, 0)
	var tempTx *transaction.Transaction
	var ok bool
	var property int
	// 创建区块链
	chain := blockchain.CreateBlockChain()
	// _的作用是忽略返回变量，这里的FIndUTXOs的作用是查找地址的账户余额，以及未曾使用的交易信息
	property, _ = chain.FindUTXOs([]byte("sovo"))
	fmt.Println("Balance of sovo: ", property)
	// 这里创建一个from sovo to sovo1的交易额为100的交易
	tempTx, ok = chain.CreateTransaction([]byte("sovo"), []byte("sovo1"), 100)
	// 创建成功，则放进交易池
	if ok {
		txPool = append(txPool, tempTx)
	}
	// 调用矿工去打包交易
	chain.Mine(txPool)
	txPool = make([]*transaction.Transaction, 0)
	property, _ = chain.FindUTXOs([]byte("sovo"))
	fmt.Println("Balance of sovo: ", property)

	tempTx, ok = chain.CreateTransaction([]byte("sovo"), []byte("Exia"), 200) // this transaction is invalid
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("sovo1"), []byte("Exia"), 50)
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("sovo"), []byte("Exia"), 100)
	if ok {
		txPool = append(txPool, tempTx)
	}
	chain.Mine(txPool)
	txPool = make([]*transaction.Transaction, 0)
	property, _ = chain.FindUTXOs([]byte("sovo"))
	fmt.Println("Balance of sovo: ", property)
	property, _ = chain.FindUTXOs([]byte("sovo1"))
	fmt.Println("Balance of sovo1: ", property)
	property, _ = chain.FindUTXOs([]byte("Exia"))
	fmt.Println("Balance of Exia: ", property)

	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PreHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Println("Proof of Work validation:", block.ValidataPow())
	}

	//I want to show the bug at this version.

	tempTx, ok = chain.CreateTransaction([]byte("sovo1"), []byte("Exia"), 30)
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("sovo1"), []byte("sovo"), 30)
	if ok {
		txPool = append(txPool, tempTx)
	}

	chain.Mine(txPool)

	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PreHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Println("Proof of Work validation:", block.ValidataPow())
	}

	property, _ = chain.FindUTXOs([]byte("sovo"))
	fmt.Println("Balance of sovo: ", property)
	property, _ = chain.FindUTXOs([]byte("sovo1"))
	fmt.Println("Balance of sovo1: ", property)
	property, _ = chain.FindUTXOs([]byte("Exia"))
	fmt.Println("Balance of Exia: ", property)
}
