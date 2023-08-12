package blockchain

import (
	"gbc/constant"
	"gbc/transaction"
	"gbc/util"
)

type TransactionPool struct {
	Transactions []*transaction.Transaction
}

func InitTransactionPool() *TransactionPool {
	pool := &TransactionPool{}
	if !util.FileIsExit(constant.TxPoolPath) {
		pool.Store()
	}
	pool.Load()
	return pool
}

func (p *TransactionPool) Load() {
	util.NewFileDB(constant.TxPoolPath).Load(p)
}

func (p *TransactionPool) Store() {
	util.NewFileDB(constant.TxPoolPath).Store(p)
}

func (p *TransactionPool) Remove() {
	util.NewFileDB(constant.TxPoolPath).Remove()
}

func (p *TransactionPool) Append(trans *transaction.Transaction) {
	p.Transactions = append(p.Transactions, trans)
}
