// プログラム外部とCUIでやりとりするためのパッケージ
package cui

import (
	"flag"
	"io"
	"log"
)

func initUsage() {
	// --help オプションをカスタマイズする
	flag.Usage = func() {
		println("Usage: 使い方")
		println("Example: 使用例")
		println("Description: 説明")
		println("Options:")
		// ここでデフォルトのオプション使い方を表示する
		flag.PrintDefaults()
	}
}

// メイン処理
func Main() {
	// -v でログを冗長に出力する
	// ポインタ型で返ってくるので値型に変換して取得する
	isVerbose := *flag.Bool("v", false, "output verbose log")

	// CLIのUsageを設定する
	initUsage()

	// コマンドライン引数を解析する
	flag.Parse()

	// verbose モードでない場合はログを出力しない
	if !isVerbose {
		log.SetOutput(io.Discard)
	}

	// パラメータ生成時にエラーが発生した場合はエラー内容を標準エラー出力してヘルプを表示する
	// if err != nil {
	// 	println("Error: " + err.Error())
	// 	flag.Usage()
	// 	return
	// }

}
