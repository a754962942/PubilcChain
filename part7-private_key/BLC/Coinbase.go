package BLC

import "fmt"

//定义奖励类型
const subsidy = 2

//创建Coinbase交易
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s", to)
	}
	//创建一个输入项
	txin := TXInput{[]byte{}, -1, data}
	//创建输出项
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].VoutIdx == -1
}
