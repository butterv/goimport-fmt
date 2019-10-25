package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"

	"github.com/istsh/goimport-fmt/ast"

	"github.com/istsh/goimport-fmt/config"
	"github.com/istsh/goimport-fmt/lexer"
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

type importStatus uint

const (
	NotYetReached importStatus = iota
	UnderAnalysis
	Finished
)

// 標準パッケージのgofmtの実装を参考にする。
// goimportsは関係なさそう。
// gofmtはimportエリアのソートはやっているが、標準パッケージとそれ以外で区別している模様。
// なのでサードパーティと自プロジェクトのパッケージをくっつけて記述すると空行が入らない。
// また、空行を1行挟んで記述すると、全く別のパッケージ群と判定するようで、その無駄な空行は削除してくれない。
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 対象となるファイルのパス
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

	config.Setup(ownProject)

	// ファイルをOpenする
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// importの各行
	var importLines []string
	// import開始フラグ
	importStatus := NotYetReached

	var beforeBlock []string
	var afterBlock []string
	for scanner.Scan() {
		// その行の内容
		lineStr := scanner.Text()

		if importStatus == NotYetReached && lineStr != "import (" {
			beforeBlock = append(beforeBlock, lineStr)
			continue
		}

		if importStatus == NotYetReached && lineStr == "import (" {
			// import部分の読み込み開始
			importStatus = UnderAnalysis
			continue
		}

		if importStatus == UnderAnalysis && lineStr == ")" {
			// import部分の読み込み終了
			importStatus = Finished
			continue
		}

		if importStatus == UnderAnalysis && lineStr != "" {
			// 対象行の内容を格納
			importLines = append(importLines, lineStr)
			continue
		}

		if importStatus == Finished {
			afterBlock = append(afterBlock, lineStr)
			continue
		}
	}

	if len(importLines) <= 1 {
		return
	}

	ids, err := lexer.Lexer(importLines)
	if err != nil {
		panic(err.Error())
	}

	nf, err := os.Create(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer nf.Close()
	for _, text := range beforeBlock {
		fmt.Fprintln(nf, text)
	}
	fmt.Fprintln(nf, "import (")

	sort.Sort(ast.ImportDetails(ids))
	beforePackageType := ast.Unknown
	for _, id := range ids {
		if beforePackageType != ast.Unknown && beforePackageType != id.PackageType {
			fmt.Fprintln(nf)
		}

		beforePackageType = id.PackageType
		if id.Alias == ast.NoAlias {
			fmt.Fprintf(nf, "\t\"%s\"\n", id.ImportStr)
		} else {
			fmt.Fprintf(nf, "\t%s \"%s\"\n", id.Alias, id.ImportStr)
		}
	}
	fmt.Fprintln(nf, ")")

	for _, text := range afterBlock {
		fmt.Fprintln(nf, text)
	}

	os.Exit(0)
}
