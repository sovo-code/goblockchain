package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"goblockchain/constcoe"
	"goblockchain/transaction"
	"goblockchain/utils"
	"reflect"
	"runtime"

	"github.com/dgraph-io/badger"
)

// 定义区块链
type BlockChain struct {
	LastHash []byte
	DataBase *badger.DB
}

// 迭代器
type BlockChainIterator struct {
	CurrentHash []byte
	DataBase    *badger.DB
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{chain.LastHash, chain.DataBase}
	return &iterator
}

// 返回一个block，并且迭代器指向前一个区块
func (iterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := iterator.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = DeSerializeBlock(val)
			return nil
		})
		utils.Handle(err)

		return err
	})
	utils.Handle(err)

	iterator.CurrentHash = block.PreHash

	return block
}

// 创建辅助函数判断迭代器是否终止
func (chain *BlockChain) BackOgPrevHash() []byte {
	var ogprevhash []byte

	err := chain.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("ogprevhash"))
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			ogprevhash = val

			return err
		})
		utils.Handle(err)

		return err
	})
	utils.Handle(err)

	return ogprevhash
}

// 区块链根据信息创建区块，由于交易创建并没有再矿工里面，所以没有收取手续费
// 修改函数，将区块加入数据库里面 思路:，lasthash，首先创建区块，获取序列化数据， 获取hash，更新数据库，lasthash,prehash更新
func (bc *BlockChain) AddBlock(newBlock *Block) {
	var lastHash []byte

	err := bc.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)

		return err
	})
	utils.Handle(err)
	// 如果矿工所打包的区块的prehash与链上的lasthash不同说明，区块过时了，所以需要驳回
	if !bytes.Equal(newBlock.PreHash, lastHash) {
		fmt.Println("This block is out of age!!! you are late!")
		runtime.Goexit()
	}

	err = bc.DataBase.Update(func(txn *badger.Txn) error {
		// 新区块加入区块链
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)
		// 更新lashhash
		err = txn.Set([]byte("lh"), newBlock.Hash)
		utils.Handle(err)
		bc.LastHash = newBlock.Hash
		return err
	})
	utils.Handle(err)
}

// //传统内存式区块链
// func CreateBlockChain() *BlockChain {
// 	blockchain := BlockChain{}
// 	blockchain.Blocks = append(blockchain.Blocks, GenesisBlock())
// 	return &blockchain
// }

// 为了永久存储，建立数据库区块链
func InitBlockChain(address []byte) *BlockChain {
	var lastHash []byte

	if utils.FileExits(constcoe.BCFile) {
		fmt.Println("blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(constcoe.BCPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		genesis := GenesisBlock(address)
		fmt.Println("Genesis Created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		utils.Handle(err)
		// lasthash
		err = txn.Set([]byte("lh"), genesis.Hash)
		utils.Handle(err)
		// prehash == address []byet("sovo")
		err = txn.Set([]byte("ogprevhash"), genesis.PreHash)
		utils.Handle(err)
		lastHash = genesis.Hash
		return err
	})
	utils.Handle(err)
	blockchain := BlockChain{LastHash: lastHash, DataBase: db}
	return &blockchain
}

// 从数据库加载区块链
func ContinueBlockChain() *BlockChain {
	// 检查数据库是否存在
	if !utils.FileExits(constcoe.BCFile) {
		fmt.Println("No blockchain was found, please create one first!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(constcoe.BCPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)
		return err
	})

	utils.Handle(err)

	chain := BlockChain{LastHash: lastHash, DataBase: db}
	return &chain
}

// 创建交易信息
// 寻找可用的交易信息(根据地址)
// 修改寻找交易信息不能再使用循环了，为此需要创建一个迭代器方便遍历
func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {
	var unSpentTxs []transaction.Transaction
	spentTxs := make(map[string][]int) //string类型为key，元素为[]int型

	iter := bc.Iterator()

all:
	for {
		block := iter.Next()
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
				if reflect.DeepEqual(out.HashPubKey, utils.PublicKeyHash(address)) {
					unSpentTxs = append(unSpentTxs, *tx)
				}
			}
			if !tx.IsBase() {
				for _, in := range tx.Inputs {
					if reflect.DeepEqual(in.PubKey, address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIdx)
					}
				}
			}
		}
		if bytes.Equal(block.PreHash, bc.BackOgPrevHash()) {
			break all
		}
	}
	return unSpentTxs
}

// 寻找地址的UTXO
func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if reflect.DeepEqual(out.HashPubKey, utils.PublicKeyHash(address)) {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				continue Work
			}
		}
	}
	return accumulated, unspentOuts
}

// 寻找大于输出的UTXO
func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0
Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if reflect.DeepEqual(out.HashPubKey, utils.PublicKeyHash(address)) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				if accumulated >= amount {
					break Work
				}
				continue Work
			}
		}
	}
	return accumulated, unspentOuts
}

// 创建交易
func (bc *BlockChain) CreateTransaction(from_PubKey, to_HashPubKey []byte, amount int, privkey ecdsa.PrivateKey) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from_PubKey, amount)
	if acc < amount {
		fmt.Println("Not Enough coins!")
		return &transaction.Transaction{}, false
	}
	// 处理可用交易
	for txid, outidx := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.Handle(err)
		input := transaction.TxInput{TxID: txID, OutIdx: outidx, PubKey: from_PubKey, Sig: nil}
		inputs = append(inputs, input)
	}

	outputs = append(outputs, transaction.TxOutput{Value: amount, HashPubKey: to_HashPubKey})
	if acc > amount {
		outputs = append(outputs, transaction.TxOutput{Value: acc - amount, HashPubKey: utils.PublicKeyHash(from_PubKey)})
	}
	tx := transaction.Transaction{ID: nil, Inputs: inputs, Outputs: outputs}
	tx.SetID()

	tx.Sign(privkey)
	return &tx, true
}
