package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"goblockchain/constcoe"
	"goblockchain/utils"
)

// 交易
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) TxHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte

	//NewEncoder返回一个将编码后数据写入w的*Encoder。
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx) //tx被编码,序列化结构体
	utils.Handle(err)         //错误处理

	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

func (tx *Transaction) SetID() {
	tx.ID = tx.TxHash()
}

// 创世交易信息，用于产生币
func BaseTx(toaddress []byte) *Transaction {
	txIn := TxInput{[]byte{}, -1, []byte{}}
	txOut := TxOutput{constcoe.InitCoin, toaddress}
	tx := Transaction{[]byte("This is the Base Transaction!"), []TxInput{txIn}, []TxOutput{txOut}}
	return &tx
}

//判断是否是创世
func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}
