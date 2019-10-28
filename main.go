package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
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
// - コメントあり

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
