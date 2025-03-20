package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/toriwasa/sqlite-playground/internal/infrastructures/db"
)

func main() {
	// コマンドライン引数を定義
	dbPath := flag.String("db", "stock_price.db", "Path to the SQLite database file")
	flag.Parse()

	// データベースからデータを取得
	dailyPrices, err := db.GetDailyStockPrices(*dbPath)
	if err != nil {
		log.Fatalf("Failed to retrieve data from database: %v", err)
	}

	// 結果を表示
	fmt.Printf("Found %d daily stock prices in database:\n\n", len(dailyPrices))
	fmt.Println("StockID\tDate\t\tPrice")
	fmt.Println("-------\t----------\t-------")

	for _, price := range dailyPrices {
		fmt.Printf("%s\t%s\t%.2f\n",
			price.StockPrice.StockID,
			price.PriceDate.Format("2006-01-02"),
			price.StockPrice.Price,
		)
	}
}
