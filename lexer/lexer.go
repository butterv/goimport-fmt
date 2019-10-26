package lexer

import (
	"fmt"

	"github.com/istsh/goimport-fmt/ast"
)

func Lexer(paths []byte) ([]*ast.ImportDetail, error) {
	var ids []*ast.ImportDetail

	for _, path := range paths {
		if path == '\n' {
			fmt.Println()
		} else {
			fmt.Print(string(path))
		}

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
