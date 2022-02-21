package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

//钱包结构
type Wallet struct {
	PrivateKey ecdsa.PrivateKey //私钥
	PublicKey  []byte           //公钥
}

const Checksum = 4
const version = byte(0x00)

//生成私钥和公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	//生成私钥
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//利用私钥推导出公钥
	pubKey := append(private.PublicKey.X.Bytes(), private.Y.Bytes()...)
	return *private, pubKey
}
func NewWallet() *Wallet {
	//随机生成密钥对
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

//计算公钥hash
func HashPubKey(pubKey []byte) []byte {
	// 1.先hash一次
	publicSHA256 := sha256.Sum256(pubKey)
	//2.计算ripemd160
	RIPEMD160Hasher := ripemd160.New()
	RIPEMD160Hasher.Write(publicSHA256[:])
	publicRIPEMD160 := RIPEMD160Hasher.sum(nil)
	return publicRIPEMD160
}

const ChecksumLen = 4

//计算校验和
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:Checksum]
}
func (w *Wallet) GetAddress() []byte{
	//1.计算公钥hash
	pubKeyhash := HashPubKey(w.PublicKey)
	//2.计算校验和
	versionedPayload := append([]byte{version}, pubKeyhash...)
	checksum := checksum(versionedPayload)
	//3.计算base58编码
	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)
	return address

}
