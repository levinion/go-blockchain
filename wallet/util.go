package wallet

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"

	"github.com/mr-tron/base58"
)

func (w *Wallet) PrivateKey() []byte {
	prvKey, _ := ecdh.P256().GenerateKey(rand.Reader)
	return prvKey.Bytes()
}

func (w *Wallet) PublicKey() []byte {
	privateKey, err := ecdh.P256().NewPrivateKey(w.privateKey)
	if err != nil {
		panic(err)
	}
	return privateKey.PublicKey().Bytes()
}

// 获取公钥哈希
func PublicKeyHash(public []byte) []byte {
	hash1 := sha256.Sum256(public)
	hash2 := sha256.Sum256(hash1[:])

	return hash2[:]
}

func PublicKeyHash2WalletAddr(hash []byte) []byte {
	return []byte(base58.Encode(hash))
}

func WalletAddr2PublicKeyHash(addr []byte) []byte {
	hash, err := base58.Decode(string(addr))
	if err != nil {
		panic(err)
	}
	return hash
}
