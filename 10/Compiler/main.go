package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

//Token そのトークンのワードと次へのポインタ、そのトークンの種類を所持する
type Token struct {
	word  string
	next  *Token
	tkind string
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
	head := strToToken(b) //headに開始トークン
	var xmloutput string
	tokenKind(head, &xmloutput)
	//fmt.Println(xmloutput)

	//パース部分
	t := head
	resultcompile := compilationEngine(t) // Tokenは参照渡し

	writeXML(fpath, true, resultcompile)
}

func compilationEngine(t *Token) string { //(*t)でTokenへアクセス
	var f string
	if t.word == "class" {
		compileClass(&t, &f)
	} else {
		fmt.Println("Error: Not start with class")
	}

	if t.next != nil {
		fmt.Println("Error: Not finish code")
	}

	//fmt.Println(f)
	return f
}

func compileClass(t **Token, file *string) {
	*file += "<" + "class" + ">\n"
	for {
		if (*t).word == "static" || (*t).word == "field" {
			compileClassVarDec(t, file) //classVarDecのラストで返す
			(*t) = (*t).next
			continue
		} else if (*t).word == "constructor" || (*t).word == "function" || (*t).word == "method" {
			compileSubroutine(t, file) //subroutineDecのラストで返す
			(*t) = (*t).next
			continue
		} else if (*t).word == "}" {
			markup(*t, file)
			*file += "</" + "class" + ">\n"
			break
		}
		markup(*t, file)
		(*t) = (*t).next
	}
}
func compileClassVarDec(t **Token, file *string) {
	*file += "<" + "classVarDec" + ">\n"
	for {
		if (*t).word == ";" {
			markup(*t, file)
			*file += "</" + "classVarDec" + ">\n"
			break
		}
		markup(*t, file)
		(*t) = (*t).next
	}
}
func compileSubroutine(t **Token, file *string) {
	*file += "<" + "subroutineDec" + ">\n"
	for {
		if (*t).word == "(" {
			markup(*t, file)
			(*t) = (*t).next
			compileParameterList(t, file)
		} else if (*t).word == "{" {
			*file += "<" + "subroutineBody" + ">\n"
			markup(*t, file)
			(*t) = (*t).next
			for { //varのぶん
				if (*t).word == "var" {
					compileVarDec(t, file)
				} else {
					break
				}
			}
			if (*t).word == "}" { // もしstatementが０個
				*file += "<" + "statements" + ">\n"
				*file += "</" + "statements" + ">\n"
				continue
			}
			compileStatements(t, file)
			continue
		} else if (*t).word == "}" {
			markup(*t, file)
			*file += "</" + "subroutineBody" + ">\n"
			*file += "</" + "subroutineDec" + ">\n"
			break
		}
		markup(*t, file)
		(*t) = (*t).next
	}
}
func compileParameterList(t **Token, file *string) {
	//()中のみ処理
	*file += "<" + "parameterList" + ">\n"
	for {
		if (*t).word == ")" {
			*file += "</" + "parameterList" + ">\n"
			break
		}
		markup(*t, file)
		(*t) = (*t).next
	}
}
func compileVarDec(t **Token, file *string) {
	*file += "<" + "varDec" + ">\n"
	for {
		if (*t).word == ";" {
			markup(*t, file)
			*file += "</" + "varDec" + ">\n"
			(*t) = (*t).next
			break
		}
		markup(*t, file)
		(*t) = (*t).next
	}
}
func compileStatements(t **Token, file *string) {
	*file += "<" + "statements" + ">\n"
	for { //再帰から戻ってくる時はそれぞれのステートメントの末尾の次のトークンで帰ってくること
		if (*t).word == "let" {
			compileLet(t, file)
		} else if (*t).word == "if" {
			compileIf(t, file)
		} else if (*t).word == "while" {
			compileWhile(t, file)
		} else if (*t).word == "do" {
			compileDo(t, file)
		} else if (*t).word == "return" {
			compileReturn(t, file)
		} else if (*t).word == "}" {
			*file += "</" + "statements" + ">\n"
			break
		}
	}
}
func compileDo(t **Token, file *string) {
	*file += "<" + "doStatement" + ">\n"
	markup(*t, file)
	*t = (*t).next
	compileSubCall(t, file)
	*t = (*t).next
	markup(*t, file)
	*file += "</" + "doStatement" + ">\n"
	*t = (*t).next

}
func compileLet(t **Token, file *string) {
	*file += "<" + "letStatement" + ">\n"
	for {
		if (*t).word == "[" {
			markup(*t, file)
			(*t) = (*t).next
			compileExpression(t, file)
		} else if (*t).word == "=" {
			markup(*t, file)
			(*t) = (*t).next
			compileExpression(t, file)
			markup(*t, file)
			break
		}

		markup(*t, file)
		(*t) = (*t).next
	}
	*file += "</" + "letStatement" + ">\n"
	(*t) = (*t).next

}
func compileWhile(t **Token, file *string) {
	*file += "<" + "whileStatement" + ">\n"
	markup(*t, file)
	(*t) = (*t).next
	markup(*t, file)
	(*t) = (*t).next
	compileExpression(t, file)
	markup(*t, file) //)
	(*t) = (*t).next
	markup(*t, file)
	(*t) = (*t).next
	compileStatements(t, file)
	markup(*t, file)
	*file += "</" + "whileStatement" + ">\n"
	(*t) = (*t).next
}
func compileReturn(t **Token, file *string) {
	*file += "<" + "returnStatement" + ">\n"
	markup(*t, file)
	(*t) = (*t).next
	if (*t).word != ";" {
		compileExpression(t, file)
	}
	markup(*t, file)
	*file += "</" + "returnStatement" + ">\n"
	(*t) = (*t).next
}
func compileIf(t **Token, file *string) {
	*file += "<" + "ifStatement" + ">\n"
	markup(*t, file)
	(*t) = (*t).next
	markup(*t, file)
	(*t) = (*t).next
	compileExpression(t, file)
	markup(*t, file) // )
	(*t) = (*t).next
	markup(*t, file)
	(*t) = (*t).next
	compileStatements(t, file)
	markup(*t, file)
	(*t) = (*t).next
	for {
		if (*t).word == "else" {
			markup(*t, file)
			(*t) = (*t).next
			markup(*t, file)
			(*t) = (*t).next
			compileStatements(t, file)
			markup(*t, file)
			(*t) = (*t).next
			continue
		} else {
			break
		}
	}
	*file += "</" + "ifStatement" + ">\n"

}
func compileExpression(t **Token, file *string) {
	op := []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}
	*file += "<" + "expression" + ">\n"
	compileTerm(t, file) //帰ってきた時に次のトークンがopなら続ける
	(*t) = (*t).next
	for {
		if search(op, (*t).word) {
			markup(*t, file)
			(*t) = (*t).next
			compileTerm(t, file)
			(*t) = (*t).next
		} else {
			break
		}
	}
	*file += "</" + "expression" + ">\n"
}
func compileTerm(t **Token, file *string) {
	*file += "<" + "term" + ">\n"
	if (*t).tkind == "identifier" {
		if ((*t).next).word == "(" {
			compileSubCall(t, file)
		} else if ((*t).next).word == "." {
			compileSubCall(t, file)
		} else if ((*t).next).word == "[" {
			markup(*t, file)
			(*t) = (*t).next
			markup(*t, file)
			(*t) = (*t).next
			compileExpression(t, file)
			markup(*t, file)
		} else {
			markup(*t, file)
		}
	} else if (*t).word == "(" {
		markup(*t, file)
		(*t) = (*t).next
		compileExpression(t, file)
		markup(*t, file)
	} else if (*t).word == "-" || (*t).word == "~" {
		markup(*t, file)
		(*t) = (*t).next
		compileTerm(t, file)
	} else {
		markup(*t, file)
	}
	*file += "<" + "/term" + ">\n"
}
func compileSubCall(t **Token, file *string) {
	//expressionListが含まれる, )まで処理してそのノードを呼び出し元に返す
	for {
		if (*t).word == "(" {
			markup(*t, file)
			*file += "<" + "expressionList" + ">\n"
			*t = (*t).next
			if (*t).word == ")" {
				continue
			}
			compileExpression(t, file)
			continue
		} else if (*t).word == ")" {
			*file += "</" + "expressionList" + ">\n"
			markup(*t, file)
			break
		} else if (*t).word == "," {
			markup(*t, file)
			*t = (*t).next
			compileExpression(t, file)
			continue
		} else {
			markup(*t, file)
		}
		*t = (*t).next
	}
}

//tokenKind: トークンの種類を決定する　(&トークン結果を書き込む)
func tokenKind(t *Token, file *string) {
	*file += "<tokens>\n"
	for { //終了はt.next　== nil
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
		t.tkind = kind

		markup(t, file)

		if t.next == nil {
			break
		} else {
			t = t.next
		}
	}
	*file += "</tokens>\n"
}

func markup(t *Token, file *string) {
	*file += "<" + t.tkind + "> "
	*file += outputTerminal(t.word)
	*file += " </" + t.tkind + ">\n"
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

func writeXML(fpath string, WritingIs bool, output string) {
	//ファイル出力のための下準備
	var pwd string
	var fname string
	pwd, fname = filepath.Split(fpath)
	fname = filename(fname)

	if WritingIs {
		//xmlファイルに出力
		file, _ := os.Create(pwd + fname + "j.xml")
		defer file.Close()
		file.Write(([]byte)(output))
	}
}

//ErrOutput エラー検出＆出力
func ErrOutput(e error) {
	if e != nil {
		panic(e)
	}
}
