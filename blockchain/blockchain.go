package blockchain

import (
	"encoding/hex"
	"fmt"
	"goblockchain/transaction"
	"goblockchain/utils"
	"reflect"
)

//定义区块链
type BlockChain struct {
	Blocks []*Block
}

//区块链根据信息创建区块，由于交易创建并没有再矿工里面，所以没有收取手续费
func (bc *BlockChain) AddBlock(txs []*transaction.Transaction) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
	bc.Blocks = append(bc.Blocks, newBlock)
}
func CreateBlockChain() *BlockChain {
	blockchain := BlockChain{}
	blockchain.Blocks = append(blockchain.Blocks, GenesisBlock())
	return &blockchain
}

//创建交易信息
//寻找可用的交易信息(根据地址)
func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {
	var unSpentTxs []transaction.Transaction
	spentTxs := make(map[string][]int) //string类型为key，元素为[]int型
	for idx := len(bc.Blocks) - 1; idx >= 0; idx-- {
		block := bc.Blocks[idx]
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		IterOutputs:
			for outIdx, out := range tx.Outputs {
				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						if spentOut == outIdx {
							continue IterOutputs
						}
					}
				}
				if reflect.DeepEqual(out.ToAddress, address) {
					unSpentTxs = append(unSpentTxs, *tx)
				}
			}
			if !tx.IsBase() {
				for _, in := range tx.Inputs {
					if reflect.DeepEqual(in.FromAddress, address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIdx)
					}
				}
			}
		}
	}
	return unSpentTxs
}

//寻找地址的UTXO
func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if reflect.DeepEqual(out.ToAddress, address) {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				continue Work
			}
		}
	}
	return accumulated, unspentOuts
}

//寻找大于输出的UTXO
func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0
Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			accumulated += out.Value
			unspentOuts[txID] = outIdx
			if accumulated >= amount {
				break Work
			}
			continue Work
		}
	}
	return accumulated, unspentOuts
}

//创建交易
func (bc *BlockChain) CreateTransaction(from, to []byte, amount int) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		fmt.Println("Not Enough coins!")
		return &transaction.Transaction{}, false
	}
	// 处理可用交易
	for txid, outidx := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.Handle(err)
		input := transaction.TxInput{TxID: txID, OutIdx: outidx, FromAddress: from}
		inputs = append(inputs, input)
	}

	outputs = append(outputs, transaction.TxOutput{Value: amount, ToAddress: to})
	if acc > amount {
		outputs = append(outputs, transaction.TxOutput{Value: acc - amount, ToAddress: from})
	}
	tx := transaction.Transaction{ID: nil, Inputs: inputs, Outputs: outputs}
	tx.SetID()

	return &tx, true
}

//blockchain.go
func (bc *BlockChain) Mine(txs []*transaction.Transaction) {
	bc.AddBlock(txs)
}
