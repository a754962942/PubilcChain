package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

//交易输入结构
type TXInput struct {
	Txid     []byte //引用交易ID
	VoutIdx  int    //引用的交易输出编号
	FromAddr string //输入方验签（数组签名）
}

//交易输出结构
type TXOutput struct {
	Value  int    //输出金额
	ToAddr string //收方验签
}

//交易结构
type Transaction struct {
	ID   []byte     //交易ID
	Vin  []TXInput  //交易输入项
	Vout []TXOutput //交易输出项
}

//将交易信息序列化后计算hash值
func (tx *Transaction) SetID() {
	//Buffer是一个实现了读写方法的可变大小的字节缓冲。本类型的零值是一个空的可用于读写的缓冲。
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	enc.Encode(tx)
	hash := sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

//判断该输入是否可以被某账户使用
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.FromAddr == unlockingData
}

//判断某输出是否可以被账户使用
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ToAddr == unlockingData
}
