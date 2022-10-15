package blockchain

import (
	"bytes"
	"crypto/sha256"
	"goblockchain/constcoe"
	"goblockchain/utils"
	"math"
	"math/big"
)

func (b *Block) GetTarget() []byte {
	target := big.NewInt(1)                          //创建一个值为1的*Int
	target.Lsh(target, uint(256-constcoe.Difficuty)) //左移位运算,增加难度
	return target.Bytes()
}

// 将nonce与数据连接起来
func (b *Block) GetBase4Nonce(nonce int64) []byte {
	data := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PreHash, utils.ToHexInt(nonce), b.Target, b.BackTransactionSummary()}, []byte{})
	return data
}

// 寻找nonce
func (b *Block) FindNonce() int64 {
	var intHash big.Int
	var intTraget big.Int
	var hash [32]byte
	var nonce int64
	nonce = 0
	intTraget.SetBytes(b.Target)

	for nonce < math.MaxInt64 {
		data := b.GetBase4Nonce(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(&intTraget) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce
}

// 验证正确性
func (b *Block) ValidataPow() bool {
	var intHash big.Int
	var intTraget big.Int
	var hash [32]byte
	intTraget.SetBytes(b.Target)
	data := b.GetBase4Nonce(b.Nonce)
	hash = sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	if intHash.Cmp(&intTraget) == -1 {
		return true
	} else {
		return false
	}
}
