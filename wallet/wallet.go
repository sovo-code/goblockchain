package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"goblockchain/constcoe"
	"goblockchain/utils"
	"io/ioutil"
)

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	// ecc所使用的曲线
	curve := elliptic.P256()
	// 传进所用的曲线，和一个随机数
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.Handle(err)
	// https://www.cnblogs.com/baiyuxiong/p/4334266.html 这里是append的第二个用法只支持两个slice的连接
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

type Wallet struct {
	Privtekey ecdsa.PrivateKey
	PublicKey []byte
}

func NewWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()
	Wallet := Wallet{Privtekey: privateKey, PublicKey: publicKey}
	return &Wallet
}

// 创建钱包地址
func (w *Wallet) Address() []byte {
	pubHash := utils.PublicKeyHash(w.PublicKey)
	return utils.PubHash2Address(pubHash)
}

// 保存钱包
func (w *Wallet) Save() {
	filename := constcoe.Wallets + string(w.Address()) + ".wlt"
	var content bytes.Buffer
	// 因为w包含ecdsa.PrivateKey，而其包含curve接口所以为了能够序列化及反序列化w所以要首先注册elliptic.P256()
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	utils.Handle(err)
	err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

// 加载钱包
func LoadWallet(address string) *Wallet {
	filename := constcoe.Wallets + address + ".wlt"
	if !utils.FileExits(filename) {
		utils.Handle(errors.New("no wallet with such address"))
	}
	var w Wallet
	gob.Register(elliptic.P256())
	fileContent, err := ioutil.ReadFile(filename)
	utils.Handle(err)
	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&w)
	utils.Handle(err)
	return &w
}
