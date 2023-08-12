package blockchain

import (
	"encoding/hex"
	"fmt"
	"gbc/constant"
	"gbc/db"
	"gbc/logger"
	"gbc/transaction"
	"gbc/util"
)

type BlockChain struct {
	//只需要知道最后一个节点的哈希即可完成遍历
	Database *db.DB
	LastHash []byte
}

// 创建区块链
func CreateBlockChain(addr []byte) *BlockChain {
	if util.FileIsExit(constant.DBPath) {
		logger.ExistError()
	}
	db := db.InitDB()
	var lastHash []byte
	genesis := GenesisBlock(addr)
	lastHash = genesis.Hash
	db.Set(genesis.Hash, genesis.Marshal())
	db.Set([]byte("lastHash"), lastHash)
	return &BlockChain{db, lastHash}
}

// 获取链
func GetBlockChain() *BlockChain {
	if !util.FileIsExit(constant.DBPath) {
		logger.NotExistError()
	}
	db := db.InitDB()
	lastHash := db.Get([]byte("lastHash"))
	return &BlockChain{Database: db, LastHash: lastHash}
}

func InitBlockChain(addr []byte) *BlockChain {
	if util.FileIsExit(constant.DBPath) {
		return GetBlockChain()
	}
	return CreateBlockChain(addr)
}

// 添加链上区块
func (bc *BlockChain) AddBlock(transaction []*transaction.Transaction) *BlockChain {
	//create new block
	newBlock := CreateBlock(bc.LastHash, transaction)
	//add new block's hash
	bc.Database.Set(newBlock.Hash, newBlock.Marshal())
	//updata last hash
	bc.LastHash = newBlock.Hash
	bc.Database.Set([]byte("lastHash"), bc.LastHash)
	return bc
}

// 寻找未花费交易，返回一个未花费交易（未花费交易金额之和即为账户余额）的列表，表示该账户所持有的utxo
func (bc *BlockChain) FindUnspentTransactions(addr []byte) []transaction.Transaction {
	unSpentTxs := make([]transaction.Transaction, 0) //未花费交易
	spentTxs := make(map[string][]int, 0)            //已花费交易
	bc.Range(func(block *Block) {                    //从后往前遍历区块
		for _, tx := range block.Transactions { //遍历交易
			txID := hex.EncodeToString(tx.ID)
		RangeOut:
			for outIndex, out := range tx.Outputs { //遍历支出
				for _, spentOut := range spentTxs[txID] { //遍历已经花费的交易
					if spentOut == outIndex {
						continue RangeOut //若支出已花费，跳过该支出
					}
				}
				if out.CheckToAddr(addr) {
					unSpentTxs = append(unSpentTxs, *tx) //若支付对象是自己，则加入未花费交易
				}
			}
			//将当前交易的上一个交易的支出标记为已花费的交易
			if !tx.IsBase() {
				for _, in := range tx.Inputs { //遍历来源
					if in.CheckFromAddr(addr) { //若交易发起方是自己，表示上一笔交易（自己作为支出方）金额已经不属于自己，该交易已消费
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIndex) //将OutIndex加入到前置交易的收入ID中
					}
				}
			}
		}
	})
	return unSpentTxs
}

// 找到所有UTXO，返回账户所拥有的总金额
func (bc *BlockChain) FindUTXOs(addr []byte) (int, map[string]int) {
	unSpentOutputs := make(map[string]int)
	unSpentTxs := bc.FindUnspentTransactions(addr)
	sum := 0
RangeTx:
	for _, tx := range unSpentTxs { //遍历未花费的交易
		txID := hex.EncodeToString(tx.ID)
		for outIndex, out := range tx.Outputs { //遍历交易输出
			if out.CheckToAddr(addr) { //若对象是自己（即自己的余额的一部分）
				sum += out.Value                //将该部分金额支出
				unSpentOutputs[txID] = outIndex //待支出的交易ID和output序号
				continue RangeTx
			}
		}
	}
	return sum, unSpentOutputs
}

// 对一笔支出，寻找账户可支出的utxo，返回支出的总金额和所有支出的交易列表（交易ID和Output序号）
func (bc *BlockChain) FindSpendableOutputs(addr []byte, amount int) (int, map[string]int) {
	unSpentOutputs := make(map[string]int)
	unSpentTxs := bc.FindUnspentTransactions(addr)
	sum := 0
RangeTx:
	for _, tx := range unSpentTxs { //遍历未花费的交易
		txID := hex.EncodeToString(tx.ID)
		for outIndex, out := range tx.Outputs { //遍历交易输出
			if out.CheckToAddr(addr) && sum < amount { //若对象是自己（即自己的余额的一部分）且统计值不超过交易额
				sum += out.Value                //将该部分金额支出
				unSpentOutputs[txID] = outIndex //待支出的交易ID和output序号
				if sum >= amount {
					break RangeTx //若金额足够则退出循环
				}
				continue RangeTx //否则继续统计下一个交易
			}
		}
	}
	return sum, unSpentOutputs
}

// 创建交易
func (bc *BlockChain) CreateTransactions(from []byte, to []byte, amount int) (*transaction.Transaction, bool) {
	trans := new(transaction.Transaction)
	money, outPuts := bc.FindSpendableOutputs(from, amount)
	if money < amount {
		fmt.Println("Not enough coins!")
		return &transaction.Transaction{}, false
	}
	// 使用查找到的outputs，形成inputs
	for txID, outIndex := range outPuts {
		id, _ := hex.DecodeString(txID)
		trans.Inputs = append(trans.Inputs, transaction.TxInput{
			TxID:     id,
			OutIndex: outIndex,
			FromAddr: from,
		})
	}
	//创建outputs
	trans.Outputs = append(trans.Outputs, transaction.TxOutput{Value: amount, ToAddr: to})
	//找零
	if money > amount {
		trans.Outputs = append(trans.Outputs, transaction.TxOutput{Value: money - amount, ToAddr: from})
	}
	//设置哈希ID
	trans.SetID()
	return trans, true
}
