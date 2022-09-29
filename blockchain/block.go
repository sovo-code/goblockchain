package blockchain

import (
	"bytes"
	"crypto/sha256"
	"goblockchain/utils"
	"time"
)

//定义区块结构体
type Block struct {
	Timestamp int64
	Hash      []byte
	PreHash   []byte
	Target    []byte //for POW difficuty
	Nonce     int64  //for POW envidence
	Data      []byte
}

//设置Hash值
func (b *Block) SetHash() {
	//通过bytes.Join将多个字节串连接，第二个参数是将字节串连接时的分隔符，这里设置为[]byte{}即为空，ToHexInt将int64转换为字节串类型。然后我们对information做哈希就可以得到区块的哈希值了。
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PreHash, b.Target, utils.ToHexInt(b.Nonce), b.Data}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

//创建区块
func CreateBlock(prehash []byte, data []byte) *Block {
	//创建
	block := Block{time.Now().Unix(), []byte{}, prehash, []byte{}, 0, data}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

//创世区块
func GenesisBlock() *Block {
	genesisblock := "Hello, BlockChain!"
	return CreateBlock([]byte{}, []byte(genesisblock))
}
