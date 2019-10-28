package lexer

import (
	"bytes"

	"github.com/istsh/goimport-fmt/ast"
)

type importStatus uint

const (
	// NotYetReached means that the import block has not yet been reached
	NotYetReached importStatus = iota
	// UnderAnalysis means that the import block is being analyzed
	UnderAnalysis
	// Finished means that the import block has been analyzed
	Finished
)

// DividedSrc has the result of dividing Src into before import block, import block, and after import block
type DividedSrc struct {
	BeforeImportDivision []byte
	ImportDivision       [][]byte
	AfterImportDivision  []byte
}

// Lexer analyzes and divides Src
func Lexer(src []byte) *DividedSrc {
	// import status flag
	importStatus := NotYetReached

	// reading src...
	var line []byte
	ds := &DividedSrc{}
	for _, ch := range src {
		line = append(line, ch)

		if ch == '\n' {
			switch importStatus {
			case NotYetReached:
				ds.BeforeImportDivision = append(ds.BeforeImportDivision, line...)
				if bytes.Equal(line, []byte("import (\n")) {
					// start reading import block
					importStatus = UnderAnalysis
				}
			case UnderAnalysis:
				if bytes.Equal(line, []byte("\n")) {
					// only line feed is skipped.
					break
				}
				if bytes.HasPrefix(line, []byte("\t//")) {
					// comment is skipped.
					break
				}
				if bytes.Equal(line, []byte(")\n")) {
					ds.AfterImportDivision = append(ds.AfterImportDivision, line...)
					// finish reading import block
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

// GetImportDetails evaluates import paths and returns ImportDetails
func (ds *DividedSrc) GetImportDetails() (ast.ImportDetails, error) {
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
