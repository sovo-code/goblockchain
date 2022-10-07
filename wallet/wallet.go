package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"goblockchain/utils"
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
