package BLC

import (
	"bytes"
	"crypto/sha256"
	"strconv"
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
}

func (block *Block) SetHash() {
	// 1.将时间戳转字节数组
	//(1).格式化时间戳
	timestr := strconv.FormatInt(block.Timestamp, 2)
	timestamp := []byte(timestr)
	//2.将除了Hash以外的其他属性，以字节数组的形式全拼接起来
	header := bytes.Join([][]byte{block.PrevBlockHash, block.Data, timestamp}, []byte{})
	//3.将拼接起来的数据进行256 hash
	hash := sha256.Sum256(header)
	//4.将hash赋给Hash属性字节
	block.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), prevBlockHash, []byte(data), []byte{}}
	block.SetHash()
	return block
}
func NewGenenisBlock() *Block {
	return NewBlock("Genenis Block", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

}
