package BLC

import (
	"bytes"
	"encoding/gob"
)

type Block struct {
	// //时间戳，创建区块时的时间
	// Timestamp int64
	// //上一个区块的hash，父哈希
	// PrevBlockHash []byte
	//Data交易数据
	Data []byte
	// //Hash 当前区块的哈希
	// Hash []byte
	//Nonce
	Nonce int64
}

//将Block对象序列化成[]byte
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}
func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return &block
}
