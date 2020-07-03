package compileengine

import (
	"Compiler/typefile"
	"fmt"
)

//CompilationEngine パーサ結果を返す
func CompilationEngine(t *typefile.Token) string { //(*t)でTokenへアクセス
	var f string
	if t.Word == "class" {
		compileClass(&t, &f)
	} else {
		fmt.Println("Error: Not start with class")
	}

	if t.Next != nil {
		fmt.Println("Error: Not finish code")
	}

	fmt.Println(f)
	return f
}

func compileClass(t **typefile.Token, file *string) {
	*file += "<" + "class" + ">\n"
	for {
		if (*t).Word == "static" || (*t).Word == "field" {
			compileClassVarDec(t, file) //classVarDecのラストで返す
			(*t) = (*t).Next
			continue
		} else if (*t).Word == "constructor" || (*t).Word == "function" || (*t).Word == "method" {
			compileSubroutine(t, file) //subroutineDecのラストで返す
			(*t) = (*t).Next
			continue
		} else if (*t).Word == "}" {
			markup(*t, file)
			*file += "</" + "class" + ">\n"
			break
		}
		markup(*t, file)
		(*t) = (*t).Next
	}
}
func compileClassVarDec(t **typefile.Token, file *string) {
	*file += "<" + "classVarDec" + ">\n"
	for {
		if (*t).Word == ";" {
			markup(*t, file)
			*file += "</" + "classVarDec" + ">\n"
			break
		}
		markup(*t, file)
		(*t) = (*t).Next
	}
}
func compileSubroutine(t **typefile.Token, file *string) {
	*file += "<" + "subroutineDec" + ">\n"
	for {
		if (*t).Word == "(" {
			markup(*t, file)
			(*t) = (*t).Next
			compileParameterList(t, file)
		} else if (*t).Word == "{" {
			*file += "<" + "subroutineBody" + ">\n"
			markup(*t, file)
			(*t) = (*t).Next
			for { //varのぶん
				if (*t).Word == "var" {
					compileVarDec(t, file)
				} else {
					break
				}
			}
			if (*t).Word == "}" { // もしstatementが０個
				*file += "<" + "statements" + ">\n"
				*file += "</" + "statements" + ">\n"
				continue
			}
			compileStatements(t, file)
			continue
		} else if (*t).Word == "}" {
			markup(*t, file)
			*file += "</" + "subroutineBody" + ">\n"
			*file += "</" + "subroutineDec" + ">\n"
			break
		}
		markup(*t, file)
		(*t) = (*t).Next
	}
}
func compileParameterList(t **typefile.Token, file *string) {
	//()中のみ処理
	*file += "<" + "parameterList" + ">\n"
	for {
		if (*t).Word == ")" {
			*file += "</" + "parameterList" + ">\n"
			break
		}
		markup(*t, file)
		(*t) = (*t).Next
	}
}
func compileVarDec(t **typefile.Token, file *string) {
	*file += "<" + "varDec" + ">\n"
	for {
		if (*t).Word == ";" {
			markup(*t, file)
			*file += "</" + "varDec" + ">\n"
			(*t) = (*t).Next
			break
		}
		markup(*t, file)
		(*t) = (*t).Next
	}
}
func compileStatements(t **typefile.Token, file *string) {
	*file += "<" + "statements" + ">\n"
	for { //再帰から戻ってくる時はそれぞれのステートメントの末尾の次のトークンで帰ってくること
		if (*t).Word == "let" {
			compileLet(t, file)
		} else if (*t).Word == "if" {
			compileIf(t, file)
		} else if (*t).Word == "while" {
			compileWhile(t, file)
		} else if (*t).Word == "do" {
			compileDo(t, file)
		} else if (*t).Word == "return" {
			compileReturn(t, file)
		} else if (*t).Word == "}" {
			*file += "</" + "statements" + ">\n"
			break
		}
	}
}
func compileDo(t **typefile.Token, file *string) {
	*file += "<" + "doStatement" + ">\n"
	markup(*t, file)
	*t = (*t).Next
	compileSubCall(t, file)
	*t = (*t).Next
	markup(*t, file)
	*file += "</" + "doStatement" + ">\n"
	*t = (*t).Next

}
func compileLet(t **typefile.Token, file *string) {
	*file += "<" + "letStatement" + ">\n"
	for {
		if (*t).Word == "[" {
			markup(*t, file)
			(*t) = (*t).Next
			compileExpression(t, file)
		} else if (*t).Word == "=" {
			markup(*t, file)
			(*t) = (*t).Next
			compileExpression(t, file)
			markup(*t, file)
			break
		}

		markup(*t, file)
		(*t) = (*t).Next
	}
	*file += "</" + "letStatement" + ">\n"
	(*t) = (*t).Next

}
func compileWhile(t **typefile.Token, file *string) {
	*file += "<" + "whileStatement" + ">\n"
	markup(*t, file)
	(*t) = (*t).Next
	markup(*t, file)
	(*t) = (*t).Next
	compileExpression(t, file)
	markup(*t, file) //)
	(*t) = (*t).Next
	markup(*t, file)
	(*t) = (*t).Next
	compileStatements(t, file)
	markup(*t, file)
	*file += "</" + "whileStatement" + ">\n"
	(*t) = (*t).Next
}
func compileReturn(t **typefile.Token, file *string) {
	*file += "<" + "returnStatement" + ">\n"
	markup(*t, file)
	(*t) = (*t).Next
	if (*t).Word != ";" {
		compileExpression(t, file)
	}
	markup(*t, file)
	*file += "</" + "returnStatement" + ">\n"
	(*t) = (*t).Next
}
func compileIf(t **typefile.Token, file *string) {
	*file += "<" + "ifStatement" + ">\n"
	markup(*t, file)
	(*t) = (*t).Next
	markup(*t, file)
	(*t) = (*t).Next
	compileExpression(t, file)
	markup(*t, file) // )
	(*t) = (*t).Next
	markup(*t, file)
	(*t) = (*t).Next
	compileStatements(t, file)
	markup(*t, file)
	(*t) = (*t).Next
	for {
		if (*t).Word == "else" {
			markup(*t, file)
			(*t) = (*t).Next
			markup(*t, file)
			(*t) = (*t).Next
			compileStatements(t, file)
			markup(*t, file)
			(*t) = (*t).Next
			continue
		} else {
			break
		}
	}
	*file += "</" + "ifStatement" + ">\n"

}
func compileExpression(t **typefile.Token, file *string) {
	op := []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}
	*file += "<" + "expression" + ">\n"
	compileTerm(t, file) //帰ってきた時に次のトークンがopなら続ける
	(*t) = (*t).Next
	for {
		if search(op, (*t).Word) {
			markup(*t, file)
			(*t) = (*t).Next
			compileTerm(t, file)
			(*t) = (*t).Next
		} else {
			break
		}
	}
	*file += "</" + "expression" + ">\n"
}
func compileTerm(t **typefile.Token, file *string) {
	*file += "<" + "term" + ">\n"
	if (*t).Tkind == "identifier" {
		if ((*t).Next).Word == "(" {
			compileSubCall(t, file)
		} else if ((*t).Next).Word == "." {
			compileSubCall(t, file)
		} else if ((*t).Next).Word == "[" {
			markup(*t, file)
			(*t) = (*t).Next
			markup(*t, file)
			(*t) = (*t).Next
			compileExpression(t, file)
			markup(*t, file)
		} else {
			markup(*t, file)
		}
	} else if (*t).Word == "(" {
		markup(*t, file)
		(*t) = (*t).Next
		compileExpression(t, file)
		markup(*t, file)
	} else if (*t).Word == "-" || (*t).Word == "~" {
		markup(*t, file)
		(*t) = (*t).Next
		compileTerm(t, file)
	} else {
		markup(*t, file)
	}
	*file += "<" + "/term" + ">\n"
}
func compileSubCall(t **typefile.Token, file *string) {
	//expressionListが含まれる, )まで処理してそのノードを呼び出し元に返す
	for {
		if (*t).Word == "(" {
			markup(*t, file)
			*file += "<" + "expressionList" + ">\n"
			*t = (*t).Next
			if (*t).Word == ")" {
				continue
			}
			compileExpression(t, file)
			continue
		} else if (*t).Word == ")" {
			*file += "</" + "expressionList" + ">\n"
			markup(*t, file)
			break
		} else if (*t).Word == "," {
			markup(*t, file)
			*t = (*t).Next
			compileExpression(t, file)
			continue
		} else {
			markup(*t, file)
		}
		*t = (*t).Next
	}
}

func markup(t *typefile.Token, file *string) {
	*file += "<" + t.Tkind + "> "
	*file += outputTerminal(t.Word)
	*file += " </" + t.Tkind + ">\n"
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