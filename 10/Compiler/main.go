package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type token struct {
	word string
	next *token
}

func main() {
	argv := os.Args
	argc := len(argv)
	if argc != 2 {
		fmt.Println("Error: argument number")
	}
	//引数がファイルかディレクトリか

	f, err := os.Open(argv[1])
	ErrOutput(err)
	defer f.Close()

	b, err := ioutil.ReadAll(f) // bをfor rangeでstring変換すると1文字ずつ取得できる
	ErrOutput(err)

	for i, word := range b {
		fmt.Print(string(word))
		if string(word) == "\n" {
			fmt.Printf("%d番目は改行です", i)
		}
	}

}

func ErrOutput(e error) {
	if e != nil {
		panic(e)
	}
}
