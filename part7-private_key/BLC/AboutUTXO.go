package BLC

import (
	"encoding/hex"
	"fmt"
	"log"
)

//查找账户可解锁的交易
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	//已经花出去的UTXO，构建tx->VOutIdx 的map
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for {
		block, next := bci.PreBlock()
		for _, tx := range block.Transactions {
			//将字节数组转换为字符串
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				fmt.Printf("FindUnspentTransactions:Outidx:%d,out:%v\n", outIdx, out)
				//如果交易已经被花出去了，直接跳过此交易
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							fmt.Println("Continue Outputs.spentTXOs[txID]:", spentTXOs[txID])
							continue Outputs
						}
					}
				}
				//可以被address解锁，就代表属于address的utxo再此交易中
				if out.IsLockedWithKey() {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			//用来维护spentTXOs,已经被引用过了，代表被使用
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.UsesKey(HashPubKey(in.PubKey)) {
						inTxID := hex.EncodeToString([]byte(in.Txid))
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.VoutIdx)
					}
				}
			}
		}
		if !next {
			break
		}
	}
	return unspentTXs
}
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	//先找所有交易
	unspentTransactions := bc.FindUnspentTransactions(address)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			//可解锁代表用户的资产
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}
func (bc *Blockchain) GetBalance(address string) {
	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance Of '%s':%d\n", address, balance)
}
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	//获取可使用交易
	unspentTXs := bc.FindUnspentTransactions(address)
	//记录余额
	accumulated := 0
Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Vout {
			fmt.Printf("FindSpendableOutputs:Outidx:%d,out:%v\n", outIdx, out)
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				//UTXO足够了就跳出循环，break可以跳出多重循环
				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutputs
}

//创建普通交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	//1.需要组合输入项和输出项
	var inputs []TXInput
	var outputs []TXOutput
	//2.查询最小UTXO
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("ERROR:Not enough funds.")
	}
	//构建输入项
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)
		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	//4.构建输出项
	outputs = append(outputs, TXOutput{amount, to})
	//需要找零
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) //a change

	}
	//5.交易生成
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}

//交易发送
func (bc *Blockchain) Send(from, to, data string, amount int) {
	//创建普通交易
	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MinedBlock([]*Transaction{tx}, data)
	fmt.Println("Success")
}
