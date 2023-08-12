package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWallet(t *testing.T) {
	wallet := CreateWallet()
	pk := wallet.PrivateKey()
	assert.Equal(t, 32, len(pk))
}
