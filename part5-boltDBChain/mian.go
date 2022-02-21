package main

import (
	"fmt"

	"github.com/a754962942/PubilcChain/part5-boltDBChain/BLC"
)

func main() {
	bc := BLC.NewBlockchain()
	bc.AddBlock("Send 200 ATOM To EE")
	bc.AddBlock("Send 100000 ETH To CC")
	// //区块链遍历
	// for _, block := range bc {
	// 	fmt.Printf("Prev.Hash:%x\n", block.PrevBlockHash)
	// 	fmt.Printf("Data:%s\n", block.Data)
	// 	fmt.Printf("Hash:%x\n", block.Hash)
	// 	fmt.Printf("Nonce:%d\n", block.Nonce)
	// 	pow := BLC.NewProofOfWork(block)
	// 	fmt.Printf("Pow:%t\n", pow.Validate())
	// 	fmt.Println()
	// }
	bci := bc.Iterator()
	for {
		block, next := bci.PreBlock()
		fmt.Printf("Prev.Hash : %x\n", block.PrevBlockHash)
		fmt.Printf("Data : %s\n", block.Data)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		pow := BLC.NewProofOfWork(block)
		fmt.Printf("POW : %t\n", pow.Validate())
		fmt.Printf("next: %t\n", next)
		fmt.Println()
		if !next {
			//next为假则代表已经到创世区块
			break
		}
	}

}
