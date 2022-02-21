package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

//交易输入结构
type TXInput struct {
	Txid      []byte //引用交易ID
	VoutIdx   int    //引用的交易输出编号
	Signature []byte //签名信息
	PubKey    []byte // 公钥
}

//交易输出结构
type TXOutput struct {
	Value      int    //输出金额
	PubkeyHash []byte //公钥Hash值
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

//交易修剪
//将交易中公钥私钥全部置为nil
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	//将原交易内的签名和公钥都置为空
	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.VoutIdx, nil, nil})
	}
	//复制输入项
	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubkeyHash})
	}
	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//1.Coinbase交易无需签名
	if tx.IsCoinbase() {
		return
	}
	//2.修剪交易
	txCopy := tx.TrimmedCopy()
	//3.循环向输入项签名
	for inID, vin := range txCopy.Vin {
		//找到输入项引用的交易
		prevTX := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTX.Vout[vin.VoutIdx].PubkeyHash
		txCopy.SetID()
		//txid生成后再把PubKey置空
		txCopy.Vin[inID].PubKey = nil
		//使用ecsda签名获得r和s
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		//形成签名数据
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].Signature = signature
	}
}


//lock sign the output
func (out *TXOutput) Lock(address []byte) {
	//Base58解密地址，返回public Key Hash+checksum
	pubKeyHash := Base58Decode(address)
	//去除checksum
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4]
}
//判断out.Hash是否与input.hash相等，即判断out是否可以接受该Transaction
func (out *TXOutput)IsLockedWithKey(pubkeyHash byte[])  bool{
	return bytes.Compare(out.PubkeyHash,pubkeyHash) == 0
}
//新建TXOutput将address转换为pubkeyhash
func  NewTXOutput(value int ,address string)*TXOutput  {
	txo:=&TXOutput{
		value,nil,
	}
	txo.Lock([]byte(address))
	return txo
}
//判断当前input.publickeyHash是否与传入的publickey相等
func (in *TXInput)UsesKey(pubKeyHash []byte)bool  {
	lockingHash:=HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash,pubKeyHash) ==0

}