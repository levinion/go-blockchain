package wallet

import (
	"gbc/constant"
	"gbc/util"
	"path/filepath"
)

type Wallet struct {
	privateKey []byte
	publicKey  []byte
}

func CreateWallet() *Wallet {
	wallet := &Wallet{}
	wallet.SetKeyPair()
	return wallet
}

func NewWallet() *Wallet {
	wallet := &Wallet{}
	return wallet
}

func (w *Wallet) SetKeyPair() {
	w.privateKey = w.PrivateKey()
	w.publicKey = w.PublicKey()
}

func (w *Wallet) Addr() []byte {
	hash := PublicKeyHash(w.publicKey)
	return PublicKeyHash2WalletAddr(hash)
}

func (w *Wallet) Store() {
	filename := filepath.Join(constant.WalletPath, string(w.Addr()))
	util.NewFileDB(filename).Store(w)
}

// 通过钱包地址获取本地钱包
func (w *Wallet) Load(addr []byte) *Wallet {
	filename := filepath.Join(constant.WalletPath, string(addr))
	util.NewFileDB(filename).Load(w)
	return w
}
