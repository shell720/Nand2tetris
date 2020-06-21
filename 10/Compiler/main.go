package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

//Token そのトークンのワードと次へのポインタを所持する
type Token struct {
	word string
	next *Token
}

//予約語
var symbol = []string{"{", "}", "(", ")", "[", "]", ".", ",",
	";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}
var keyword = []string{"class", "constructor", "function", "method",
	"field", "static", "var", "int", "char", "boolean", "void", "true",
	"false", "null", "this", "let", "do", "if", "else", "while", "return"}

func main() {
	argv := os.Args
	argc := len(argv)
	if argc != 2 {
		fmt.Println("Error: argument number")
	}
	//引数がファイルかディレクトリか場合わけ
	//ディレクトリとファイルの切り分けも
	finfo, _ := os.Stat(argv[1])
	if finfo.IsDir() {
		files, _ := ioutil.ReadDir(argv[1])
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".jack" {
				Tokenizer(argv[1] + f.Name())
			}
		}
	} else {
		Tokenizer(argv[1])
	}
}

//Tokenizer ファイルを開いて字句解析
func Tokenizer(fpath string) {
	//ファイルを開く
	f, err := os.Open(fpath)
	ErrOutput(err)
	defer f.Close()

	b, err := ioutil.ReadAll(f) // bをfor rangeでstring変換すると1文字ずつ取得できる
	ErrOutput(err)

	//文字列を字句解析
	var t *Token
	t = strToToken(b) //t.wordにはclassが入る、終了はt.next　== nil
	xmloutput := "<tokens>\n"
	for {
		tkind := tokenKind(t)
		xmloutput += "<" + tkind + "> "
		xmloutput += outputTerminal(t.word)
		xmloutput += " </" + tkind + ">\n"
		if t.next == nil {
			break
		} else {
			t = t.next
		}
	}
	xmloutput += "</tokens>\n"

	//ファイル出力のための下準備
	var pwd string
	var fname string
	pwd, fname = filepath.Split(fpath)
	fname = filename(fname)
	//xmlファイルに出力
	file, _ := os.Create(pwd + fname + "Tj.xml")
	defer file.Close()

	file.Write(([]byte)(xmloutput))
}

//tokenKind: トークンの種類を返す
func tokenKind(t *Token) string {
	var kind string
	_, err := strconv.Atoi(t.word)
	doubleQuote := "\""
	byteDQ := []byte(doubleQuote)
	switch {
	case search(keyword, t.word):
		kind = "keyword"
	case search(symbol, t.word):
		kind = "symbol"
	case t.word[0] == byteDQ[0]:
		kind = "stringConstant"
	case err == nil:
		kind = "integerConstant"
	default:
		kind = "identifier"
	}
	return kind
}

//strToToken 文字列を分割して終端文字に
//コメント削除、空白スペースを手がかりに分ける、symbolが現れても分ける
func strToToken(b []byte) *Token {
	l := len(b)
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
			//fmt.Println(cur.word)

		case b[i] == 34: //ダブルクォートで始まる文字列を処理
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

//outputTerminal: 終端文字の細かい変化
func outputTerminal(s string) string {
	var t string
	if s[0] == 34 {
		t = s[1 : len(s)-1]
	} else if s == "<" {
		t = "&lt;"
	} else if s == ">" {
		t = "&gt;"
	} else if s == "&" {
		t = "&amp;"
	} else {
		t = s
	}
	return t
}

//filename: 拡張子なしのファイル名を返す
func filename(f string) string {
	var end int
	dot := "."
	bytedot := []byte(dot)
	for i := 0; i < len(f); i++ {
		if f[i] == bytedot[0] {
			end = i
		}
	}
	return f[:end]
}

//ErrOutput エラー検出＆出力
func ErrOutput(e error) {
	if e != nil {
		panic(e)
	}
}
