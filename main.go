package main

import (
	"log"

	"github.com/toriwasa/sqlite-playground/internal/presenter/cui"
)

func main() {
	// DEBUGログのフォーマットを設定
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// DEBUGログのプレフィックスを設定
	log.SetPrefix("DEBUG: ")

	// CUIのメイン処理を呼び出す
	cui.Main()

}
