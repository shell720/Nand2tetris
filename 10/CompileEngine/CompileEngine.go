package compileEngine

import (
	"Compiler/typefile"
	"fmt"
)

func dfs(tree typefile.ParseVertex) int {
	if len(tree.ChildList) == 0 {
		if tree.Word != "" {
			fmt.Println(tree.Word)
			return 0
		}
	}
	for _, v := range tree.ChildList {
		dfs(v)
	}
	return 0
}

//CompilationEngine パーサ結果を返す
func CompilationEngine(t *typefile.Token) string { //(*t)でTokenへアクセス
	var ret typefile.ParseVertex
	var f string
	ret.Name = "HEAD"
	var childs []typefile.ParseVertex
	if t.Word == "class" {
		res := compileClass(&t, &f)
		//fmt.Println(res)
		childs = append(childs, res)
	} else {
		fmt.Println("Error: Not start with class")
	}

	if t.Next != nil {
		fmt.Println("Error: Not finish code")
	}
	ret.ChildList = childs
	//fmt.Println(ret.ChildList[0].ChildList[2].T.Word)
	dfs(ret)
	fmt.Println(ret)
	return f
}

func compileClass(t **typefile.Token, file *string) typefile.ParseVertex {
	*file += "<" + "class" + ">\n"
	var ret typefile.ParseVertex
	ret.Name = "class"
	var childs []typefile.ParseVertex
	for {
		var tmp typefile.ParseVertex
		tmp.Word = (*t).Word
		tmp.Tkind = (*t).Tkind
		ret.ChildNum++
		if (*t).Word == "static" || (*t).Word == "field" {
			tmp = compileClassVarDec(t, file) //classVarDecのラストで返す
			(*t) = (*t).Next
			childs = append(childs, tmp)
			continue
		} else if (*t).Word == "constructor" || (*t).Word == "function" || (*t).Word == "method" {
			tmp = compileSubroutine(t, file) //subroutineDecのラストで返す
			(*t) = (*t).Next
			childs = append(childs, tmp)
			continue
		} else if (*t).Word == "}" {
			childs = append(childs, tmp)
			markup(*t, file)
			*file += "</" + "class" + ">\n"
			break
		}
		markup(*t, file)
		(*t) = (*t).Next
		childs = append(childs, tmp)
	}
	ret.ChildList = childs
	return ret
}
func compileClassVarDec(t **typefile.Token, file *string) typefile.ParseVertex {
	*file += "<" + "classVarDec" + ">\n"
	var ret typefile.ParseVertex
	ret.Name = "classVarDec"
	var childs []typefile.ParseVertex
	for {
		ret.ChildNum++
		var tmp typefile.ParseVertex
		tmp.Word = (*t).Word
		tmp.Tkind = (*t).Tkind
		if (*t).Word == ";" {
			markup(*t, file)
			childs = append(childs, tmp)
			*file += "</" + "classVarDec" + ">\n"
			break
		}
		markup(*t, file)
		(*t) = (*t).Next
		childs = append(childs, tmp)
	}
	ret.ChildList = childs
	return ret
}
func compileSubroutine(t **typefile.Token, file *string) typefile.ParseVertex {
	*file += "<" + "subroutineDec" + ">\n"
	var ret typefile.ParseVertex
	ret.Name = "subroutineDec"
	var childs []typefile.ParseVertex
	for {
		var tmp typefile.ParseVertex
		if (*t).Word == "(" {
			markup(*t, file)
			tmp.Word = (*t).Word
			tmp.Tkind = (*t).Tkind
			childs = append(childs, tmp)
			(*t) = (*t).Next
			res := compileParameterList(t, file)
			childs = append(childs, res)
			continue
		} else if (*t).Word == "{" {
			*file += "<" + "subroutineBody" + ">\n"
			var res typefile.ParseVertex
			res.Name = "subroutineBody"
			var childs1 []typefile.ParseVertex

			var tmp1 typefile.ParseVertex
			tmp1.Word = (*t).Word
			tmp1.Tkind = (*t).Tkind
			childs1 = append(childs1, tmp1)
			markup(*t, file)
			(*t) = (*t).Next
			for { //varのぶん
				if (*t).Word == "var" {
					tmp2 := compileVarDec(t, file)
					childs1 = append(childs1, tmp2)
				} else {
					break
				}
			}
			if (*t).Word == "}" { // もしstatementが０個
				*file += "<" + "statements" + ">\n"
				*file += "</" + "statements" + ">\n"
				continue
			}
			tmp3 := compileStatements(t, file)
			childs1 = append(childs1, tmp3)
			res.ChildList = childs1
			childs = append(childs, res)
			continue
		} else if (*t).Word == "}" {
			markup(*t, file)
			tmp.Word = (*t).Word
			tmp.Tkind = (*t).Tkind
			childs = append(childs, tmp)
			*file += "</" + "subroutineBody" + ">\n"
			*file += "</" + "subroutineDec" + ">\n"
			break
		}
		markup(*t, file)
		tmp.Word = (*t).Word
		tmp.Tkind = (*t).Tkind
		(*t) = (*t).Next
		childs = append(childs, tmp)
	}
	ret.ChildList = childs
	return ret
}
func compileParameterList(t **typefile.Token, file *string) typefile.ParseVertex {
	//()中のみ処理
	*file += "<" + "parameterList" + ">\n"
	var ret typefile.ParseVertex
	ret.Name = "parameterList"
	var childs []typefile.ParseVertex
	for {
		var tmp typefile.ParseVertex
		if (*t).Word == ")" {
			*file += "</" + "parameterList" + ">\n"
			break
		}
		markup(*t, file)
		tmp.Word = (*t).Word
		tmp.Tkind = (*t).Tkind
		childs = append(childs, tmp)
		ret.ChildNum++
		(*t) = (*t).Next
	}
	ret.ChildList = childs
	return ret
}
func compileVarDec(t **typefile.Token, file *string) typefile.ParseVertex {
	*file += "<" + "varDec" + ">\n"
	var ret typefile.ParseVertex
	ret.Name = "varDec"
	var childs []typefile.ParseVertex
	for {
		var tmp typefile.ParseVertex
		ret.ChildNum++
		tmp.Word = (*t).Word
		tmp.Tkind = (*t).Tkind
		childs = append(childs, tmp)
		if (*t).Word == ";" {
			markup(*t, file)
			*file += "</" + "varDec" + ">\n"
			(*t) = (*t).Next
			break
		}
		markup(*t, file)
		(*t) = (*t).Next
	}
	ret.ChildList = childs
	return ret
}
func compileStatements(t **typefile.Token, file *string) typefile.ParseVertex {
	*file += "<" + "statements" + ">\n"
	var ret typefile.ParseVertex
	ret.Name = "statements"
	var childs []typefile.ParseVertex
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
	ret.ChildList = childs
	return ret
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
