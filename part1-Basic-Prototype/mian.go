package main

import (
	"fmt"
	"time"

	"github.com/a754962942/PubilcChain/part1-Basic-Prototype/BLC"
)

func main() {
	blockchain := BLC.NewBlockchain()
	blockchain.AddBlock("Send 30 BTC to a From admin")
	blockchain.AddBlock("Send 20 BTC to b From admin")
	blockchain.AddBlock("Send 10 BTC to c From admin")
	for i, block := range blockchain.Blocks {
		if i == 0 {
			fmt.Printf("创世区块\n")
		} else {
			fmt.Printf("第%d个区块\n", i)
		}
		fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05 "))
		fmt.Printf("Data:%s\n", block.Data)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("Hash:%x\n", block.Hash)

	}
}
