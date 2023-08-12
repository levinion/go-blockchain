package main

import (
	"fmt"
	"gbc/blockchain"
	"gbc/flax"
	"gbc/logger"
	"gbc/wallet"
	"os"
	"strconv"
)

func main() {
	// create block chain
	flax.App("create").
		Func(func(c *flax.Context) {
			//create chain
			blockchain.CreateBlockChain([]byte(c.Args[0]))
			logger.Success()
		}).ExactArgs(1)

	// send some money
	flax.App("send").
		Var(flax.String, "f", "", "from").
		Var(flax.String, "t", "", "to").
		Var(flax.Bool, "r", false, "if send by refname").
		Func(func(c *flax.Context) {
			from := c.Get("f").(string)
			to := c.Get("t").(string)
			amount, _ := strconv.Atoi(c.Args[0])
			bc := blockchain.GetBlockChain()
			if c.Get("r").(bool) {
				var err error
				ref := wallet.NewRefList().Load()
				from, err = ref.FindRef(from)
				if err != nil {
					panic(err) //fuck u err!
				}
				to, err = ref.FindRef(to)
				if err != nil {
					panic(err)
				}
			}
			trans, ok := bc.CreateTransactions([]byte(from), []byte(to), amount)
			if ok {
				pool := blockchain.InitTransactionPool()
				pool.Append(trans)
				pool.Store()
			}
		}).Rule(func(c *flax.Context) bool {
		return c.ExactArgs(1) && c.Exists("f") && c.Exists("t")
	})

	// show balance
	flax.App("balance").Func(func(c *flax.Context) {
		bc := blockchain.GetBlockChain()
		bc.GetBalance(c.Args...)
	}).Rule(func(c *flax.Context) bool {
		return len(c.Args) >= 1
	})

	flax.App("commit").Func(func(c *flax.Context) {
		bc := blockchain.GetBlockChain()
		pool := blockchain.InitTransactionPool()
		bc.AddBlock(pool.Transactions)
	}).ExactArgs(0)

	flax.App("reset").Func(func(c *flax.Context) {
		os.RemoveAll("./temp")
		fmt.Println("reset success!")
	}).ExactArgs(0)

	flax.App("create_wallet").
		Var(flax.String, "r", "", "by refname").
		Func(func(c *flax.Context) {
			w := wallet.CreateWallet()
			w.SetKeyPair()
			w.Store()
			ref := wallet.NewRefList().Load()
			if c.Exists("r") {
				r := c.Get("r").(string)
				ref.BindRef(string(w.Addr()), r)
				ref.Store()
			}
			fmt.Println("success!")
		}).ExactArgs(0)

	flax.App("info").
		Var(flax.Bool, "r", false, "by refname").
		Func(func(c *flax.Context) {
			addr := c.Args[0]
			ref := wallet.NewRefList().Load()
			if c.Get("r").(bool) {
				addr, _ = ref.FindRef(addr)
				fmt.Println(c.Args[0])
			}
			w := wallet.NewWallet().Load([]byte(addr))
			fmt.Printf("wallet address: %x\npublic key: %x\nrefname: %x\n",
				w.Addr(), w.PublicKey(), (*ref)[string(w.Addr())])
		}).ExactArgs(1)

	flax.Run()
}
