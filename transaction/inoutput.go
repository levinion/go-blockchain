package transaction

import (
	"bytes"
	"fmt"
)

// 收款信息，表明资金的来源，每个输入绑定一个或多个交易的支出，可以定位到上一笔交易
type TxInput struct {
	TxID     []byte //支出方ID
	OutIndex int    //支出方的Output序号
	FromAddr []byte //付款者hash
}

func (in *TxInput) CheckFromAddr(addr []byte) bool {
	return bytes.Equal(in.FromAddr, addr)
}

func (in *TxInput) String() string {
	return fmt.Sprintf("Transaction ID: %s, OutIndex: %d, From: %s\n",
		in.TxID, in.OutIndex, in.FromAddr)
}

// 付款信息
type TxOutput struct {
	Value  int    //金额
	ToAddr []byte //收款者hash
}

func (out *TxOutput) CheckToAddr(addr []byte) bool {
	return bytes.Equal(out.ToAddr, addr)
}

func (out *TxOutput) String() string {
	return fmt.Sprintf("Value: %d, To: %s\n", out.Value, out.ToAddr)
}
