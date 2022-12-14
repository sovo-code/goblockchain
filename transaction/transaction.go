package transaction

import (
	"bytes"
	"crypto/ecdsa"
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
	txIn := TxInput{[]byte{}, -1, []byte{}, nil}
	txOut := TxOutput{constcoe.InitCoin, toaddress}
	tx := Transaction{[]byte("This is the Base Transaction!"), []TxInput{txIn}, []TxOutput{txOut}}
	return &tx
}

// 判断是否是创世
func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}

func (tx *Transaction) PlainCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, txin := range tx.Inputs {
		inputs = append(inputs, TxInput{txin.TxID, txin.OutIdx, nil, nil})
	}

	for _, txout := range tx.Outputs {
		outputs = append(outputs, TxOutput{txout.Value, txout.HashPubKey})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

func (tx *Transaction) PlainHash(inidx int, prevHash []byte) []byte {
	txCopy := tx.PlainCopy()
	txCopy.Inputs[inidx].PubKey = prevHash
	return txCopy.TxHash()
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsBase() {
		return
	}
	for idx, input := range tx.Inputs {
		plainhash := tx.PlainHash(idx, input.PubKey)
		signature := utils.Sign(plainhash, privKey)
		tx.Inputs[idx].Sig = signature
	}
}

// verify
func (tx *Transaction) Verify() bool {
	for idx, input := range tx.Inputs {
		plainhash := tx.PlainHash(idx, input.PubKey)
		if !utils.Verify(plainhash, input.PubKey, input.Sig) {
			return false
		}
	}
	return true
}
