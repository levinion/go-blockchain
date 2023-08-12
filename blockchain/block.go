package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"gbc/constant"
	"gbc/transaction"
	"gbc/util"
	"math/big"
	"time"
)

type Block struct {
	//时间戳、哈希、前一个区块的哈希，三者共同构成头信息
	Timestamp    int64
	Hash         []byte
	PrevHash     []byte
	Target       []byte //目标哈希
	Nonce        int64
	Transactions []*transaction.Transaction
}

func (b *Block) SetHash() {
	b.Hash = b.GetHash(b.Nonce)
}

func (b *Block) GetTarget() []byte {
	target := big.NewInt(1)
	target.Lsh(target, 256-constant.Difficulty) //左移，难度值越大则范围越小，游戏越难
	return target.Bytes()
}

// 得到区块中存在的所有交易的哈希
func (b *Block) GetTransactionsHash() []byte {
	//拼接所有交易ID
	buf := make([][]byte, 0)
	for _, trans := range b.Transactions {
		buf = append(buf, trans.ID)
	}
	info := bytes.Join(buf, []byte{})
	hash := sha256.Sum256(info)
	return hash[:]
}

func (b *Block) FindNonce() int64 {
	nonce := int64(0)
	target := big.NewInt(0).SetBytes(b.Target)
	for {
		hash := big.NewInt(0).SetBytes(b.GetHash(nonce))
		if hash.Cmp(target) == -1 {
			break
		}
		nonce++
	}

	return nonce
}

// 获取该区块的哈希值，通过将其他所有字段拼接并计算sha256
func (b *Block) GetHash(nonce int64) []byte {
	// 字段拼接
	info := bytes.Join(
		[][]byte{
			util.Int64ToBytes(b.Timestamp),
			b.PrevHash,
			b.Target,
			util.Int64ToBytes(nonce),
			b.GetTransactionsHash(),
		}, []byte{})

	// 哈希计算
	hash := sha256.Sum256(info)
	return hash[:]
}

// 验证区块nonce
func (b *Block) Valid() bool {
	target := big.NewInt(0).SetBytes(b.Target)
	hash := big.NewInt(b.Nonce)
	return hash.Cmp(target) == -1
}

// 创建一个区块，使用前一个区块的哈希和data
func CreateBlock(prevHash []byte, transactions []*transaction.Transaction) *Block {
	block := &Block{
		Timestamp:    time.Now().Unix(),
		PrevHash:     prevHash,
		Transactions: transactions,
	}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()

	block.SetHash()
	return block
}

// 创世区块，为链的第一个区块
func GenesisBlock(addr []byte) *Block {
	return CreateBlock([]byte("genesis"),
		[]*transaction.Transaction{
			transaction.BaseTx(addr),
		})
}

func (b *Block) String() string {
	r := fmt.Sprintf("Timestamp: %d\n", b.Timestamp) +
		fmt.Sprintf("Hash: %x\n", b.Hash) +
		fmt.Sprintf("PrevHash: %x\n", b.PrevHash) +
		fmt.Sprintf("Nonce: %d\n", b.Nonce) +
		fmt.Sprintf("Proof Valid: %v\n", b.Valid())
	for _, tx := range b.Transactions {
		r += fmt.Sprintf("Transactions:\n %v\n", tx)
	}
	return r
}

func (b *Block) Marshal() []byte {
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(b)
	return buf.Bytes()
}

func UnMarshal(stru []byte) *Block {
	buf := bytes.NewBuffer(stru)
	var block Block
	gob.NewDecoder(buf).Decode(&block)
	return &block
}
