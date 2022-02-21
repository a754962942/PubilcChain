package BLC

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

//使用boltdb取代slice存储
const (
	//db文件名称
	dbFile = "blockchain.db"
	//bucket名称
	blocksBucket = "blocks"
	//定义矿工地址
	miner = "DDHW"
	//创世区块留言
	genesisCoinbaseData = "This is GenesisBlock"
)

type Blockchain struct {
	//记录区块hash值
	tip []byte
	db  *bolt.DB
}

//创建区块链结构，初始化只有创世区块
func CreateBlockchain() *Blockchain {
	//1.只能第一次创建
	if dbExists() {
		fmt.Println("BlockChain already exists.")
		os.Exit(1)
	}
	var tip []byte
	//2.没有则创建文件
	db, _ := bolt.Open("dbfile", 0600, nil)
	//接下来是更新数据库操作
	db.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(blocksBucket))
		if buck == nil {
			//2.2.1 第一次用，创建创世区块
			fmt.Println("No Existing BlockChain found. Creating a new one.")
			cbtx := NewCoinbaseTX(miner, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)
			//2.2.2区块链数据编码
			block_data := genesis.Serialize()
			//2.2.3创建新的bucket，存入区块信息
			bucket, _ := tx.CreateBucket([]byte(blocksBucket))
			bucket.Put(genesis.Hash, block_data)
			bucket.Put([]byte("last"), genesis.Hash)
			tip = genesis.Hash
		} else {
			//不是第一次使用，之前有块
			tip = buck.Get([]byte("last"))
		}
		return nil
	})
	return &Blockchain{tip, db}
}

func (bc *Blockchain) MinedBlock(transactions []*Transaction, data string) {
	var tip []byte
	//1.获取tip值,此时不能再打开数据库文件,要用区块的结构
	bc.db.View(func(tx *bolt.Tx) error {
		//获取buck
		buck := tx.Bucket([]byte(blocksBucket))
		//获取tip
		tip = buck.Get([]byte("last"))
		return nil
	})
	//2.更新数据库
	bc.db.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(blocksBucket))
		//创建Coinbase交易
		cbtx := NewCoinbaseTX(miner, data)
		transactions := append(transactions, cbtx)
		// fmt.Println("This is Coinbase.")
		block := NewBlock(transactions, tip)
		//将新区块放入db
		buck.Put(block.Hash, block.Serialize())
		buck.Put([]byte("last"), block.Hash)
		//覆盖tip值
		bc.tip = block.Hash
		return nil
	})
}

//判断区块链是否已经存在
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}
