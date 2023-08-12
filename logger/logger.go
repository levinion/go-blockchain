package logger

import (
	"fmt"
	"os"
)

const (
	sucess        = "Sucess!"
	existError    = "Chain already exists!"
	notExistError = "Chain not found!"
)

func Success() {
	fmt.Println(sucess)
}

func ExistError() {
	fmt.Println(existError)
	os.Exit(1)
}

func NotExistError() {
	fmt.Println(notExistError)
	os.Exit(1)
}
