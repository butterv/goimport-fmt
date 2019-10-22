package lexer

import (
	"strings"

	"github.com/istsh/goimport-fmt/ast"
)

func Lexer(strs []string) ([]*ast.ImportDetail, error) {
	var ids []*ast.ImportDetail

	if len(strs) == 0 {
		return ids, nil
	}

	if len(strs) == 1 {
		str := strs[0]
		trimStr := strings.Trim(str, "import ")
		replaceStr := strings.Replace(trimStr, "\"", "", -1)
		splitStrs := strings.Split(replaceStr, " ")
		if len(splitStrs) == 2 {
			id, err := ast.AnalyzeIncludeAlias(splitStrs[0], splitStrs[1])
			if err != nil {
				// TODO: error handling
			}
			ids = append(ids, id)
		} else {
			id, err := ast.Analyze(splitStrs[0])
			if err != nil {
				// TODO: error handling
			}
			ids = append(ids, id)
		}
		return ids, nil
	}

	for _, str := range strs {
		if str == "" {
			continue
		}

		if str == "import (" || str == ")" {
			ids = append(ids, &ast.ImportDetail{
				Alias:       ast.NoAlias,
				ImportStr:   str,
				PackageType: ast.Unknown,
			})
			continue
		}

		trimStr := strings.Trim(str, "\t")
		replaceStr := strings.Replace(trimStr, "\"", "", -1)
		splitStrs := strings.Split(replaceStr, " ")

		if len(splitStrs) == 2 {
			id, err := ast.AnalyzeIncludeAlias(splitStrs[0], splitStrs[1])
			if err != nil {
				// TODO: error handling
			}
			ids = append(ids, id)
		} else {
			id, err := ast.Analyze(splitStrs[0])
			if err != nil {
				// TODO: error handling
			}
			ids = append(ids, id)
		}
	}

	return ids, nil
}
