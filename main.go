package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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
	l := 0
	for scanner.Scan() {
		l++
		// ここで一行ずつ処理
		fmt.Printf("L%d: %s\n", l, scanner.Text())
	}
}
