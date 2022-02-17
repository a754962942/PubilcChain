package BLC

type Blockchain struct {
	//存储有序的区块
	Blocks []*Block
}

//新增区块
func (blockchain *Blockchain) AddBlock(data string) {
	//1.查询父区块
	preBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	//2.新建区块
	newBlock := NewBlock(data, preBlock.Hash)
	//3.将新建区块添加进BlockChain数组
	blockchain.Blocks = append(blockchain.Blocks, newBlock)
}

//创建一个带有创世区块的区块链
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenenisBlock()}}
}
