package lexer

import (
	"fmt"
	"strings"

	"github.com/istsh/goimport-fmt/ast"
)

// TODO: []stringを引数にしたほうがいいかも。(1文字ずつ検証する必要がないから)
func Lexer(paths []byte) ([]*ast.ImportDetail, error) {
	var ids []*ast.ImportDetail

	var bs []byte
	for _, path := range paths {
		switch path {
		case '\n':
			fmt.Printf("%s\n", strings.TrimLeft(string(bs), "\t"))
			bs = nil
		case '(':
			bs = append(bs, path)
			fmt.Print("Undefined ")
		case ')':
			bs = append(bs, path)
			fmt.Print("Undefined ")
		default:
			bs = append(bs, path)
		}

		str := "github.com/istsh/imnoo"
		b := []byte(str)

		//
		//	//if path == "" || path == "\t" {
		//	//	continue
		//	//}
		//	//
		//	//trimStr := strings.Trim(path, "\t")
		//	//replaceStr := strings.Replace(trimStr, "\"", "", -1)
		//	//splitStrs := strings.Split(replaceStr, " ")
		//	//
		//	//if len(splitStrs) == 2 {
		//	//	id, err := ast.AnalyzeIncludeAlias(splitStrs[0], splitStrs[1])
		//	//	if err != nil {
		//	//		return nil, err
		//	//	}
		//	//	ids = append(ids, id)
		//	//} else {
		//	//	id, err := ast.Analyze(splitStrs[0])
		//	//	if err != nil {
		//	//		return nil, err
		//	//	}
		//	//	ids = append(ids, id)
		//	//}
	}

	return ids, nil
}
