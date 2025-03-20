package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/toriwasa/sqlite-playground/internal/infrastructures/db"
	"github.com/toriwasa/sqlite-playground/internal/infrastructures/file"
)

func main() {
	// コマンドライン引数を定義
	tsvPath := flag.String("tsv", "internal/data/sample_daily_stock_price.tsv", "Path to the TSV file")
	dbPath := flag.String("db", "stock_price.db", "Path to the SQLite database file")
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	// ログの設定
	log.SetPrefix("StockPriceImporter: ")
	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	// TSVファイルの存在確認
	if _, err := os.Stat(*tsvPath); os.IsNotExist(err) {
		log.Fatalf("TSV file not found: %s", *tsvPath)
	}

	// TSVファイルを読み込む
	log.Printf("Reading TSV file: %s", *tsvPath)
	dailyPrices, err := file.ReadDailyStockPriceFromTSV(*tsvPath)
	if err != nil {
		log.Fatalf("Failed to read TSV file: %v", err)
	}
	log.Printf("Read %d daily stock prices", len(dailyPrices))

	// データベースディレクトリを作成
	dbDir := filepath.Dir(*dbPath)
	if dbDir != "." {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}
	}

	// SQLiteデータベースを初期化
	log.Printf("Initializing SQLite database: %s", *dbPath)
	if err := db.InitializeDailyStockPriceTable(*dbPath, dailyPrices); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Printf("Database initialized successfully")

	// 確認のためにデータベースからデータを取得
	retrievedPrices, err := db.GetDailyStockPrices(*dbPath)
	if err != nil {
		log.Fatalf("Failed to retrieve data from database: %v", err)
	}
	log.Printf("Retrieved %d daily stock prices from database", len(retrievedPrices))

	// 成功メッセージを表示
	fmt.Printf("Successfully imported %d daily stock prices into %s\n", len(dailyPrices), *dbPath)
}
