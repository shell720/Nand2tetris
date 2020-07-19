package main

import (
	"Compiler/compileEngine"
	"Compiler/jackTokenizer"
	"Compiler/symbolTable"
	"Compiler/typefile"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

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
				head := jackTokenizer.Tokenizer(argv[1] + f.Name())
				parse(head, argv[1]+f.Name())
			}
		}
	} else {
		head := jackTokenizer.Tokenizer(argv[1])
		parse(head, argv[1])
	}
}

func parse(head *typefile.Token, fpath string) {
	t := head
	parser := compileEngine.CompilationEngine(t)

	var result string
	var classSymTable map[string]typefile.TableValue
	var subrouSymTable map[string]typefile.TableValue
	symbolTable.DFS(parser, &result, parser, -1, classSymTable, subrouSymTable)
	//fmt.Println(result)
	writeXML(fpath, true, result)
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
