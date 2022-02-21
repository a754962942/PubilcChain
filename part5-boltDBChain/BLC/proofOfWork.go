package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

type ProofOfWork struct {
	//当前需要验证的区块
	block *Block
	//大数存储,区块难度
	target *big.Int
}

//难度值
const targetBits = 24

//Nonce上限
var maxNonce = math.MaxInt64

//准备数据
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			//父区块Hash
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToByte(pow.block.Timestamp),
			IntToByte(int64(targetBits)),
			IntToByte(nonce),
		},
		[]byte{},
	)
	return data
}
func (pow *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf("Mining the block containing %s,maxNonce = %d\n", pow.block.Data, maxNonce)
	for nonce < maxNonce {
		//数据准备
		data := pow.prepareData(int64(nonce))
		//计算hash
		hash = sha256.Sum256(data)

		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("\r%x", hash)
	fmt.Print("\n\n")
	return int64(nonce), hash[:]
}

//构造函数
func NewProofOfWork(block *Block) *ProofOfWork {
	//target为最终难度值
	target := big.NewInt(1)
	// fmt.Println("-----------------")
	//target为1向左移256-targetBits位
	target.Lsh(target, uint(256-targetBits))
	// fmt.Println(target)
	//生成pow结构
	pow := &ProofOfWork{block: block, target: target}
	return pow
}

//校验区块正确性
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1

}
