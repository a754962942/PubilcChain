package main

import (
	"fmt"

	"github.com/a754962942/PubilcChain/part3-Serialize-DeserializeBlock/BLC"
)

func main() {

	block := BLC.Block{
		[]byte("send 10000 BTC to Du"),
		100000,
	}
	fmt.Printf("Data:%s\n", block.Data)
	fmt.Printf("Nonce:%d\n", block.Nonce)
	fmt.Println()
	bytes := block.Serialize()
	fmt.Println(bytes)
	fmt.Println("")
	bc := BLC.DeserializeBlock(bytes)
	fmt.Printf("Data:%s\n", bc.Data)
	fmt.Printf("Nonce:%d\n", bc.Nonce)
	fmt.Println()

}
