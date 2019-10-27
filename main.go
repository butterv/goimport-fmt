package main

import (
	"io/ioutil"
	"os"

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
// - コメントあり

type importStatus uint

const (
	NotYetReached importStatus = iota
	UnderAnalysis
	Finished
)

type codeDivision struct {
	tempSrc []byte

	beforeImportDivision []byte
	importDivision       []byte
	afterImportDivision  []byte
}

// 標準パッケージのgofmtの実装を参考にする。
// goimportsは関係なさそう。
// gofmtはimportエリアのソートはやっているが、標準パッケージとそれ以外で区別している模様。
// なのでサードパーティと自プロジェクトのパッケージをくっつけて記述すると空行が入らない。
// また、空行を1行挟んで記述すると、全く別のパッケージ群と判定するようで、その無駄な空行は削除してくれない。
func main() {
	// 標準パッケージのみで実装するからgo.modも不要。modulesをoffにしていいかも。

	// 処理を分割して考える
	// 1. コマンド引数や環境変数の読み込みは、initでまとめて行う。
	// 2. ファイルを開く
	// 3. `1行目からimportの直前まで`,`import部`,`importの直後から最後まで`の3つに分割する。
	// 4. import部を解析し、`標準パッケージ`,`サードパーティ`,`自プロジェクト`の情報を付与する。コメントも維持する。
	// 5. 空行を全て除去
	// 6. タイプ毎にソート
	// 7. タイプ間に空行を入れる。
	// 8. 分割して3つをファイルに書き込み、保存。

	// ファイルをOpenする
	filePath := config.GetEnv().GetFilePath()
	// TODO: 権限確認
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err.Error())
	}

	// import開始フラグ
	importStatus := NotYetReached

	// ファイルを読んでいく
	// TODO: 1文字ずつ解析していく必要はないので、stringで比較してもいいかも、
	codeDivision := codeDivision{}
	for _, lineByte := range src {
		codeDivision.tempSrc = append(codeDivision.tempSrc, lineByte)
		if lineByte == '\n' {
			lineStr := string(codeDivision.tempSrc)
			switch importStatus {
			case NotYetReached:
				if lineStr == "import (\n" {
					codeDivision.importDivision = append(codeDivision.importDivision, codeDivision.tempSrc...)
					// import部分の読み込み開始
					importStatus = UnderAnalysis
				} else {
					codeDivision.beforeImportDivision = append(codeDivision.beforeImportDivision, codeDivision.tempSrc...)
				}
			case UnderAnalysis:
				codeDivision.importDivision = append(codeDivision.importDivision, codeDivision.tempSrc...)
				if lineStr == ")\n" {
					// import部分の読み込み終了
					importStatus = Finished
				}
			case Finished:
				codeDivision.afterImportDivision = append(codeDivision.afterImportDivision, codeDivision.tempSrc...)
			}

			codeDivision.tempSrc = nil
		}
	}

	if len(codeDivision.importDivision) <= 1 {
		return
	}

	_, err = lexer.Lexer(codeDivision.importDivision)
	if err != nil {
		panic(err.Error())
	}
	//
	//nf, err := os.Create(filePath)
	//if err != nil {
	//	panic(err.Error())
	//}
	//defer nf.Close()
	//for _, text := range beforeBlock {
	//	fmt.Fprintln(nf, text)
	//}
	//fmt.Fprintln(nf, "import (")
	//
	//sort.Sort(ast.ImportDetails(ids))
	//beforePackageType := ast.Unknown
	//for _, id := range ids {
	//	if beforePackageType != ast.Unknown && beforePackageType != id.PackageType {
	//		fmt.Fprintln(nf)
	//	}
	//
	//	beforePackageType = id.PackageType
	//	if id.Alias == ast.NoAlias {
	//		fmt.Fprintf(nf, "\t\"%s\"\n", id.ImportStr)
	//	} else {
	//		fmt.Fprintf(nf, "\t%s \"%s\"\n", id.Alias, id.ImportStr)
	//	}
	//}
	//fmt.Fprintln(nf, ")")
	//
	//for _, text := range afterBlock {
	//	fmt.Fprintln(nf, text)
	//}

	os.Exit(0)
}
