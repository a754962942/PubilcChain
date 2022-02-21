package BLC

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/boltdb/bolt"
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
type BlockchainIterator struct {
	currentHash []byte   //当前区块Hash
	db          *bolt.DB //已经打开的数据库
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
func NewGenesisBlock() *Block {
	return NewBlock("Genenis Block", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

}
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	//编码器
	encoder := gob.NewEncoder(&result)
	//编码
	encoder.Encode(b)
	return result.Bytes()
}
func DeserializeBlock(d []byte) *Block {
	var block Block
	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(d))
	//解析区块数据
	decoder.Decode(&block)
	return &block
}
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}
func (i *BlockchainIterator) PreBlock() (*Block, bool) {
	var block *Block
	//根据hash获取块数据
	i.db.View(func(tx *bolt.Tx) error {
		//获取bucket
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		//解码当前块数据
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	//当前hash变更为前一块hash
	i.currentHash = block.PrevBlockHash
	//返回区块
	return block, len(i.currentHash) > 0
}
