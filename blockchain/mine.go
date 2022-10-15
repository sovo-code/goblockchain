// mine.go
package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"goblockchain/transaction"
	"goblockchain/utils"
	"log"
)

func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()
	//In the near future, we'll have to validate the transactions first here.
	if !bc.VerifyTransactions(transactionPool.PubTx) {
		log.Println("falls in transactions verification")
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	}
	candidateBlock := CreateBlock(bc.LastHash, transactionPool.PubTx) //PoW has been done here.
	if candidateBlock.ValidataPow() {
		bc.AddBlock(candidateBlock)
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	} else {
		fmt.Println("Block has invalid nonce.")
		return
	}
}

func (bc *BlockChain) VerifyTransactions(txs []*transaction.Transaction) bool {
	if len(txs) == 0 {
		return true
	}

	spentOutPuts := make(map[string]int)
	for _, tx := range txs {
		pubKey := tx.Inputs[0].PubKey
		unspentOutPuts := bc.FindUnspentTransactions(pubKey)
		inputAmount := 0
		OutputAmount := 0

		for _, input := range tx.Inputs {
			if outidx, ok := spentOutPuts[hex.EncodeToString(input.TxID)]; ok && outidx == input.OutIdx {
				return false
			}
			ok, amount := isInputRight(unspentOutPuts, input)
			if !ok {
				return false
			}
			inputAmount += amount
			spentOutPuts[hex.EncodeToString(input.TxID)] = input.OutIdx
		}

		for _, output := range tx.Outputs {
			OutputAmount += output.Value
		}
		if inputAmount != OutputAmount {
			return false
		}

		if !tx.Verify() {
			return false
		}
	}
	return true
}

func isInputRight(txs []transaction.Transaction, in transaction.TxInput) (bool, int) {
	for _, tx := range txs {
		if bytes.Equal(tx.ID, in.TxID) {
			return true, tx.Outputs[in.OutIdx].Value
		}
	}
	return false, 0
}
