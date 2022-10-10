package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"goblockchain/transaction"
	"goblockchain/utils"
	"time"
)

// 定义区块结构体
type Block struct {
	Timestamp    int64
	Hash         []byte
	PreHash      []byte
	Target       []byte //for POW difficuty
	Nonce        int64  //for POW envidence
	Transactions []*transaction.Transaction
}

// 返回交易信息汇总
func (b *Block) BackTransactionSummary() []byte {
	txIDs := make([][]byte, 0)
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}
	summary := bytes.Join(txIDs, []byte{})
	return summary
}

// 设置Hash值
func (b *Block) SetHash() {
	//通过bytes.Join将多个字节串连接，第二个参数是将字节串连接时的分隔符，这里设置为[]byte{}即为空，ToHexInt将int64转换为字节串类型。然后我们对information做哈希就可以得到区块的哈希值了。
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PreHash, b.Target, utils.ToHexInt(b.Nonce), b.BackTransactionSummary()}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

// 创建区块
func CreateBlock(prehash []byte, txs []*transaction.Transaction) *Block {
	//创建
	block := Block{time.Now().Unix(), []byte{}, prehash, []byte{}, 0, txs}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

// 创世区块
func GenesisBlock(address []byte) *Block {
	tx := transaction.BaseTx(address)
	genesis := CreateBlock([]byte("sovo"), []*transaction.Transaction{tx})
	genesis.SetHash()
	return genesis
}

// Badger的键值对只支持字节串存储形式
// 序列化
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	utils.Handle(err)
	return res.Bytes()
}

// 反序列化
func DeSerializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.Handle(err)
	return &block
}
