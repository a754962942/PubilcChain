package main

import (
	"fmt"

	"github.com/a754962942/PubilcChain/part2-proof-of-work/BLC"
)

func main() {
	bc := BLC.NewBlockchain()
	bc.AddBlock("Send 1 BTC To A")
	bc.AddBlock("Send 3 BTC To B")
	//区块链遍历
	for _, block := range bc.Blocks {
		fmt.Printf("Prev.Hash:%x\n", block.PrevBlockHash)
		fmt.Printf("Data:%s\n", block.Data)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		pow := BLC.NewProofOfWork(block)
		fmt.Printf("Pow:%t\n", pow.Validate())
		fmt.Println()
	}
}
