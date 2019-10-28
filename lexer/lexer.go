package lexer

import (
	"bytes"

	"github.com/istsh/goimport-fmt/ast"
)

type importStatus uint

const (
	NotYetReached importStatus = iota
	UnderAnalysis
	Finished
)

type DevidedSrc struct {
	BeforeImportDivision []byte
	ImportDivision       [][]byte
	AfterImportDivision  []byte
}

func Lexer(src []byte) *DevidedSrc {
	// import開始フラグ
	importStatus := NotYetReached

	// ファイルを読んでいく
	var line []byte
	ds := &DevidedSrc{}
	for _, ch := range src {
		line = append(line, ch)

		if ch == '\n' {
			switch importStatus {
			case NotYetReached:
				ds.BeforeImportDivision = append(ds.BeforeImportDivision, line...)
				if bytes.Equal(line, []byte("import (\n")) {
					// import部分の読み込み開始
					importStatus = UnderAnalysis
				}
			case UnderAnalysis:
				if bytes.Equal(line, []byte("\n")) {
					// 空行はスキップ
					break
				}
				if bytes.HasPrefix(line, []byte("\t//")) {
					// コメントもスキップ
					break
				}
				if bytes.Equal(line, []byte(")\n")) {
					ds.AfterImportDivision = append(ds.AfterImportDivision, line...)
					// import部分の読み込み終了
					importStatus = Finished
					break
				}

				ds.ImportDivision = append(ds.ImportDivision, line)
			case Finished:
				ds.AfterImportDivision = append(ds.AfterImportDivision, line...)
			}

			line = nil
		}
	}

	return ds
}

func (ds *DevidedSrc) GetImportDetails() (ast.ImportDetails, error) {
	var ids []*ast.ImportDetail
	for _, importPath := range ds.ImportDivision {
		trimBytes := bytes.TrimLeft(importPath, "\t")
		trimBytes = bytes.TrimRight(trimBytes, "\n")
		splitBytes := bytes.Split(bytes.ReplaceAll(trimBytes, []byte("\""), []byte("")), []byte(" "))

		var id *ast.ImportDetail
		var err error
		if len(splitBytes) <= 1 {
			id, err = ast.Analyze(splitBytes[0])
		} else {
			id, err = ast.AnalyzeIncludeAlias(splitBytes[0], splitBytes[1])
		}
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}
