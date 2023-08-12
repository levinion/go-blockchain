package blockchain

import (
	"bytes"
	"gbc/db"
)

type BlockIter struct {
	block *Block
	db    *db.DB
}

func (bc *BlockChain) Iter() *BlockIter {
	ob := bc.Database.Get(bc.LastHash)
	return &BlockIter{block: UnMarshal(ob), db: bc.Database}
}

func (iter *BlockIter) Next() {
	ob := iter.db.Get(iter.block.PrevHash)
	block := UnMarshal(ob)
	iter.block = block
}

func (iter *BlockIter) End() bool {
	return bytes.Equal(iter.block.PrevHash, []byte("genesis"))
}

func (bc *BlockChain) Range(f func(block *Block)) {
	iter := bc.Iter()
	for ; !iter.End(); iter.Next() {
		f(iter.block)
	}
	//now iter point to genesis block
	f(iter.block)
}
