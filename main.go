package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/istsh/goimport-fmt/ast"
	"github.com/istsh/goimport-fmt/config"
	"github.com/istsh/goimport-fmt/lexer"
)

func setup() {
	filePathPtr := flag.String("filepath", "", "file path")
	ownProjectPtr := flag.String("ownproject", "", "own project")
	flag.Parse()

	filePath := *filePathPtr
	if filePath == "" {
		panic("file path not found")
	}
	ownProject := *ownProjectPtr
	if ownProject == "" {
		panic("own project not found")
	}

	goroot := os.Getenv("GOROOT")

	config.Set(goroot, filePath, ownProject)
}

// The import path can be divided into three types.
// 1. Standard package
// 2. Third-party package
// 3. Own project package
func main() {
	setup()
	// a target filepath
	filePath := config.GetFilePath()
	// permission: -rw-------(u=rw)
	var perm os.FileMode = 0600
	f, err := os.OpenFile(filePath, os.O_RDONLY, perm)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err.Error())
	}

	ds := lexer.Lexer(src)
	res, err := createResult(ds)
	if err != nil {
		panic(err.Error())
	}

	if bytes.Equal(src, res) {
		// No change
		return
	}

	err = ioutil.WriteFile(filePath, res, perm)
	if err != nil {
		panic(err.Error())
	}

	os.Exit(0)
}

func createResult(ds *lexer.DividedSrc) ([]byte, error) {
	ids, err := ds.GetImportDetails()
	if err != nil {
		return nil, err
	}

	sort.Sort(ids)

	var res []byte
	res = append(res, ds.BeforeImportDivision...)

	createBytes := func(id *ast.ImportDetail) []byte {
		var bs []byte
		bs = append(bs, '\t')
		if len(id.Alias) > 0 {
			bs = append(bs, id.Alias...)
			bs = append(bs, ' ')
		}
		bs = append(bs, '"')
		bs = append(bs, id.ImportPath...)
		bs = append(bs, '"')
		bs = append(bs, '\n')
		return bs
	}

	beforePackageType := ast.Unknown
	for _, id := range ids {
		switch id.PackageType {
		case ast.Standard:
			if beforePackageType == ast.Unknown {
				beforePackageType = ast.Standard
			}
			res = append(res, createBytes(id)...)
		case ast.ThirdParty:
			if beforePackageType == ast.Standard {
				beforePackageType = ast.ThirdParty
				res = append(res, '\n')
			}
			res = append(res, createBytes(id)...)
		case ast.OwnProject:
			if beforePackageType == ast.ThirdParty {
				beforePackageType = ast.OwnProject
				res = append(res, '\n')
			}
			res = append(res, createBytes(id)...)
		default:
			return nil, fmt.Errorf("unsupported package type: %d", id.PackageType)
		}
	}

	res = append(res, ds.AfterImportDivision...)
	return res, nil
}
