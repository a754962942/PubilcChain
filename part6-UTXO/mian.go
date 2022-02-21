package main

import (
	"fmt"

	"github.com/a754962942/PubilcChain/part6-UTXO/BLC"
)

func main() {
	//创世区块初始化区块链
	bc := BLC.CreateBlockchain()
	//获取Balance
	bc.GetBalance("DDHW")
	//发送8个给黑小雨
	bc.Send("DDHW", "HHXY", "MONEY", 2)
	//查询余额
	bc.GetBalance("DDHW")
	bc.GetBalance("HHXY")
	// bc.AddBlock("Send 100BTC to miner")
	fmt.Println()
	//遍历
	bci := bc.Iterator()
	for {
		block, next := bci.PreBlock()
		fmt.Printf("Prev.Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data :%s\n", block.Transactions[0].Vin[0].FromAddr)
		fmt.Printf("Hash :%x\n", block.Hash)
		fmt.Printf("Nonce :%d\n", block.Nonce)
		pow := BLC.NewProofOfWork(block)
		fmt.Printf("POW : %t\n", pow.Validate())
		fmt.Println()
		if !next {
			//next为假，代表当前区块为创世区块，可以结束循环
			break
		}
	}
}
