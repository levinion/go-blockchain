package blockchain

import "fmt"

func (bc *BlockChain) GetBalance(users ...string) {
	fmt.Printf("block number: %d\n", bc.Len())
	for _, user := range users {
		money, _ := bc.FindUTXOs([]byte(user))
		fmt.Println(user+":", money)
	}
}

func (bc *BlockChain) Len() int {
	n := 0
	bc.Range(func(block *Block) {
		n++
	})
	return n
}

func (bc *BlockChain) String() string {
	s := fmt.Sprintf("block number: %d\n", bc.Len())
	bc.Range(func(block *Block) {
		s += fmt.Sprintf("%v\n", *block)
	})
	return s
}
