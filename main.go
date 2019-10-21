package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// importは3種類
// 1. 標準パッケージ
// 2. サードパーティパッケージ
// 3. 自プロジェクトパッケージ
//
// importのルール
// - 種類毎にまとめて記述し、種類間は1行の空行を挟む
// - 種類内ではパスで昇順ソート
// - 1つしかない場合は丸括弧なし
// - パスはダブルクォートで挟む
//
// import文のパターン
// - パスのみの記述
// - エイリアスあり

type PackageType string

const (
	Unknown    PackageType = "unknown"
	Standard               = "standard"
	ThirdParty             = "third party"
	OwnProject             = "own project"
)

type Environment struct {
	GOROOT      string
	GOPATH      string
	GO111MODULE string
}

var Env *Environment

func init() {
	goroot := os.Getenv("GOROOT")
	fmt.Printf("goroot: %s\n", goroot)
	gopath := os.Getenv("GOPATH")
	fmt.Printf("gopath: %s\n", gopath)
	go111module := os.Getenv("GO111MODULE")
	fmt.Printf("go111module: %s\n", go111module)

	Env = &Environment{
		GOROOT:      goroot,
		GOPATH:      gopath,
		GO111MODULE: go111module,
	}
}

func main() {
	// 対象となるファイルのパス
	pathPtr := flag.String("filepath", "", "file path")
	flag.Parse()

	path := *pathPtr
	if path == "" {
		panic("file path not found")
	}

	fmt.Printf("target file: %s\n", path)

	// ファイルをOpenする
	f, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// 行数
	var line, importStart, importEnd uint
	// importの各行
	var importLines []string
	// import開始フラグ
	var importFlag bool
	for scanner.Scan() {
		line++
		// その行の内容
		lineStr := scanner.Text()
		// ここで一行ずつ処理
		//fmt.Printf("L%d: %s\n", line, lineStr)

		if !importFlag && lineStr == "import (" {
			// import部分の読み込み開始
			importFlag = true
			importStart = line
		}

		if importFlag {
			// 対象行の内容を格納
			importLines = append(importLines, lineStr)
		}

		if importFlag && lineStr == ")" {
			// import部分の読み込み終了
			importFlag = false
			importEnd = line
			break
		}
	}

	fmt.Printf("start: %d, end: %d\n", importStart, importEnd)

	fmt.Println("## Before")

	for _, path := range importLines {
		fmt.Printf("%s\n", path)
	}

	fmt.Println("## After")

	fset := token.NewFileSet()
	fp, err := parser.ParseFile(fset, path, nil, parser.Mode(0))
	if err != nil {
		panic(err.Error())
	}

	for _, d := range fp.Imports {
		//ast.Print(fset, d)
		p := strings.Replace(d.Path.Value, "\"", "", -1)

		var packageTyape = Unknown

		{
			isStandard, err := isStandardPackage(p)
			if err != nil {
				fmt.Printf("err: %s\n", err.Error())
			}
			if isStandard {
				packageTyape = Standard
			}
		}

		if packageTyape == Unknown {
			isThirdParty, err := isThirdPartyPackage(p)
			if err != nil {
				fmt.Printf("err: %s\n", err.Error())
			}
			if isThirdParty {
				packageTyape = ThirdParty
			}
		}

		if packageTyape == Unknown {
			isOwnProject, err := isOwnProjectPackage(p)
			if err != nil {
				fmt.Printf("err: %s\n", err.Error())
			}
			if isOwnProject {
				packageTyape = OwnProject
			}
		}

		fmt.Printf("%s(%s)\n", p, packageTyape)
	}
}

func isStandardPackage(path string) (bool, error) {
	p := fmt.Sprintf("%s/src/%s", Env.GOROOT, strings.Trim(path, "\""))
	//fmt.Printf("path: %s\n", p)

	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func isThirdPartyPackage(path string) (bool, error) {
	p := fmt.Sprintf("%s/pkg/mod/%s", Env.GOPATH, strings.Trim(path, "\""))
	//fmt.Printf("path: %s\n", p)

	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func isOwnProjectPackage(path string) (bool, error) {
	p := fmt.Sprintf("%s/src/%s", Env.GOPATH, strings.Trim(path, "\""))
	//fmt.Printf("path: %s\n", p)

	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
