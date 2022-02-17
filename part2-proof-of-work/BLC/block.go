package BLC

import (
	"time"
)

type Block struct {
	//时间戳，创建区块时的时间
	Timestamp int64
	//上一个区块的hash，父哈希
	PrevBlockHash []byte
	//Data交易数据
	Data []byte
	//Hash 当前区块的哈希
	Hash []byte
	//Nonce
	Nonce int64
}

//构造函数
func NewBlock(data string, prevBlockHash []byte) *Block {
	//创建区块
	block := &Block{time.Now().Unix(), prevBlockHash, []byte(data), []byte{}, 0}
	//将block作为参数，创建一个pow对象
	pow := NewProofOfWork(block)
	//Run(执行一次工作量证明)
	nonce, hash := pow.Run()
	//设置区块Hash
	block.Hash = hash[:]
	// //设置Nonce
	block.Nonce = nonce
	return block
}
func NewGenenisBlock() *Block {
	return NewBlock("Genenis Block", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

}
