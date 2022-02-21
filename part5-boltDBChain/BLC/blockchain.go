package BLC

import (
	"fmt"

	"github.com/boltdb/bolt"
)

//使用boltdb取代slice存储
const (
	//db文件名称
	dbFile = "blockchain.db"
	//bucket名称
	blocksBucket = "blocks"
)

type Blockchain struct {
	//记录区块hash值
	tip []byte
	db  *bolt.DB
}

//创建一个带有创世区块的区块链
func NewBlockchain() *Blockchain {
	var tip []byte
	//1.打开数据库文件
	db, _ := bolt.Open(dbFile, 0600, nil)
	//2.更新数据库
	db.Update(func(tx *bolt.Tx) error {
		//2.1获取bucket
		buck := tx.Bucket([]byte(blocksBucket))
		if buck == nil {
			//2.2.1第一次使用,创建创世区块
			fmt.Println("No Existing BlockChain Found. Creating a new one ...")
			genesis := NewGenesisBlock()
			//2.2.2区块数据编码
			block_data := genesis.Serialize()
			//2.2.3创建新bucket,存入区块信息
			bucket, _ := tx.CreateBucket([]byte(blocksBucket))
			bucket.Put(genesis.Hash, block_data)
			bucket.Put([]byte("last"), genesis.Hash)
			tip = genesis.Hash
		} else {
			//2.3不是第一次使用,之前有块
			tip = buck.Get([]byte("last"))
		}
		return nil
	})
	//3.记录Blockchain信息
	return &Blockchain{tip, db}
}
func (bc *Blockchain)AddBlock(data string)  {
	var tip []byte
	//1.获取tip值,此时不能再打开数据库文件,要用区块的结构
	bc.db.View(func(tx *bolt.Tx)error{
		//获取buck
		buck:=tx.Bucket([]byte(blocksBucket))
		//获取tip
		tip = buck.Get([]byte("last"))
		return nil
	})
	//2.更新数据库
	bc.db.Update(func(tx *bolt.Tx)error{
		buck:=tx.Bucket([]byte(blocksBucket))
		block :=NewBlock(data,tip)
		//将新区块放入db
		buck.Put(block.Hash,block.Serialize())
		buck.Put([]byte("last"),block.Hash)
		//覆盖tip值
		bc.tip = block.Hash
		return nil
	})
}