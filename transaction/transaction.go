package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"gbc/constant"
)

type Transaction struct {
	ID      []byte     // 哈希值，标志当前交易的ID
	Inputs  []TxInput  //资金的（所有）来源
	Outputs []TxOutput //交易去向，可以有多个，目的包括支持找零、多收款方等
}

func (t *Transaction) Hash() []byte {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(t) //结构体序列化，类似json
	hash := sha256.Sum256(buf.Bytes())
	return hash[:]
}

func (t *Transaction) SetID() {
	t.ID = t.Hash()
}

// 指定去向，创建第一个交易
func BaseTx(toAddr []byte) *Transaction {
	return &Transaction{
		ID: []byte(""),
		Inputs: []TxInput{
			{
				TxID:     []byte{},
				OutIndex: -1, //不存在来源
				FromAddr: []byte{},
			},
		},
		Outputs: []TxOutput{
			{
				Value:  constant.InitCoin, //初始发币金额
				ToAddr: toAddr,            //发币去向
			},
		},
	}
}

func (t *Transaction) IsBase() bool {
	return len(t.Inputs) == 1 && t.Inputs[0].OutIndex == -1 //没有来源，标志资金是凭空产生的
}

func (t *Transaction) String() string {
	s := fmt.Sprintf("ID: %x\n", t.ID)
	for i, in := range t.Inputs {
		s += fmt.Sprintf("Input%d: %v\n", i+1, in)
	}
	for i, out := range t.Outputs {
		s += fmt.Sprintf("Output%d: %v\n", i+1, out)
	}
	return s
}
