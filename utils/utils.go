package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"goblockchain/constcoe"
	"log"
	"math/big"
	"os"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

// 错误处理函数,注意：函数名大写开头才能被导出使用
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToHexInt(num int64) []byte {
	//创建一段缓存
	buff := new(bytes.Buffer)
	//采用大端存储将num写入buff，返回报错信息
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// 检测文件是否存在
func FileExits(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}

// 构造公钥hash和钱包地址的函数
// 将公钥转化为公钥hash
func PublicKeyHash(publicKey []byte) []byte {
	hashPublicKey := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(hashPublicKey[:])
	Handle(err)
	publicRipeMd := hasher.Sum(nil)

	return publicRipeMd
}

// 检查位生成函数
func CheckSum(ripeMdHash []byte) []byte {
	// 取公钥hash两次sha256前constcoe.CheckLength位数据作为检查位
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:constcoe.ChecksumLength]
}

// base256转base58及其反函数
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Handle(err)
	return decode
}

// 公钥hash生成钱包地址
func PubHash2Address(pubKeyHash []byte) []byte {
	networkVersionedHash := append([]byte{constcoe.NetWorkVersion}, pubKeyHash...)
	checkSum := CheckSum(networkVersionedHash)
	finalHash := append(networkVersionedHash, checkSum...)
	address := Base58Encode(finalHash)
	return address
}

func Address2PubHash(address []byte) []byte {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-constcoe.ChecksumLength]
	return pubKeyHash
}

// 签名
func Sign(msg []byte, privKey ecdsa.PrivateKey) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, msg)
	Handle(err)
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

// 认证
func Verify(msg []byte, pubkey []byte, signature []byte) bool {
	curve := elliptic.P256()
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubkey)
	x.SetBytes(pubkey[:(keyLen / 2)])
	y.SetBytes(pubkey[(keyLen / 2):])

	rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

	return ecdsa.Verify(&rawPubKey, msg, &r, &s)
}
