package symbolTable

import (
	"Compiler/typefile"
	"fmt"
)

func DFS(tree typefile.ParseVertex, file *string, p typefile.ParseVertex, idx int,
	ct map[string]typefile.TableValue, st map[string]typefile.TableValue) {
	if tree.Name != "" && tree.Name != "HEAD" && tree.Name != "subCall" {
		*file += "<" + tree.Name + ">\n"
	}
	if len(tree.ChildList) == 0 {
		if tree.Word != "" {
			markup(tree, file)
			symbol(tree, p, idx)
		}
	}
	for i, v := range tree.ChildList {
		DFS(v, file, tree, i, ct, st)
	}
	if tree.Name != "" && tree.Name != "HEAD" && tree.Name != "subCall" {
		*file += "</" + tree.Name + ">\n"
	}
}

func symbol(t typefile.ParseVertex, p typefile.ParseVertex, i int) {
	if t.Tkind == "identifier" {
		fmt.Print(t.Word, " ,")
		//fmt.Print(p.Name, " ,")
		if i != 1 {
			if p.ChildList[0].Word == "var" {
				fmt.Print("category: var")
				fmt.Print(", 型: ", p.ChildList[1].Word)
			} else if p.ChildList[0].Word == "static" {
				fmt.Print("category: static")
				fmt.Print(", 型: ", p.ChildList[1].Word)
			} else if p.ChildList[0].Word == "field" {
				fmt.Print("category: field")
				fmt.Print(", 型: ", p.ChildList[1].Word)
			}
		}
		if p.Name == "class" {
			fmt.Print("category: class")
		} else if p.Name == "parameterList" {
			fmt.Print("category: argument")
			fmt.Print(", 型", p.ChildList[i-1].Word)
		}
		fmt.Println("")
	}
}

func markup(t typefile.ParseVertex, file *string) {
	*file += "<" + t.Tkind + "> "
	*file += outputTerminal(t.Word)
	*file += " </" + t.Tkind + ">\n"
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
