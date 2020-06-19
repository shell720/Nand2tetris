package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

//Token そのトークンのワードと次へのポインタを所持する
type Token struct {
	word string
	next *Token
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

	//文字列を分割してトークンに
	//コメントは削除する
	var t *Token
	t = strToToken(b)
	fmt.Println(t.word)

}

//strToToken トークン処理
func strToToken(b []byte) *Token {
	l := len(b)
	symbol := []string{"{", "}", "(", ")", "[", "]", ".", ",",
		";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}
	var head Token // 開始を表すノードにしたいが返り値の調整のため次のcurが開始ノード
	cur := new(Token)
	head.next = cur
	for i := 0; i < l; {
		switch {
		case string(b[i:i+2]) == "//": //コメントアウト処理1
			for {
				i++
				if string(b[i]) == "\n" {
					break
				}
			}
			i++
			continue

		case string(b[i:i+2]) == "/*": //コメントアウト処理2
			for {
				i++
				if string(b[i:i+2]) == "*/" {
					break
				}
			}
			i += 2
			continue

		case string(b[i]) == " " || string(b[i]) == "\n": //改行または空白
			i++
			continue

		case b[i] == 13:
			i++
			continue

		case b[i] == 9: //タブコード
			i++
			continue

		case search(symbol, string(b[i])):
			//fmt.Println(string(b[i]))
			s := string(b[i])
			cur = tokenConnect(cur, s)
			i++
			fmt.Println(cur.word)

		case b[i] == 34: //ダブルクォートの文字列
			startIdx := i
			for {
				i++
				if b[i] == 34 {
					i++
					break
				}
			}
			//fmt.Println(string(b[startIdx:i]))
			s := string(b[startIdx:i])
			cur = tokenConnect(cur, s)
			//fmt.Println(cur.word)

		default:
			startIdx := i
			for {
				i++
				if string(b[i]) == " " || string(b[i]) == "\n" {
					break
				} else if search(symbol, string(b[i])) { //symbolで区切る
					break
				}
			}
			//fmt.Println(string(b[startIdx:i]))
			s := string(b[startIdx:i])
			cur = tokenConnect(cur, s)
			//fmt.Println(cur.word)
		}
	}
	return head.next.next
}

//tokenConnect　トークンを繋ぐ
func tokenConnect(cur *Token, s string) *Token {
	next := &Token{word: s}
	cur.next = next
	return next
}

//search: charで与えられた文字がarrayの要素かどうか調べる
func search(array []string, char string) bool {
	check := false
	for _, e := range array {
		if e == char {
			check = true
		}
	}
	return check
}

//ErrOutput エラー検出＆出力
func ErrOutput(e error) {
	if e != nil {
		panic(e)
	}
}
