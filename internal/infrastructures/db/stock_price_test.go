package db

import (
	"os"
	"testing"
	"time"

	"github.com/toriwasa/sqlite-playground/internal/domain/models"
)

func TestInitializeAndGetDailyStockPriceTable(t *testing.T) {
	// Arrange
	dbPath := "./test_stock_price.db"
	// テスト終了後にデータベースファイルを削除
	defer os.Remove(dbPath)

	// テスト用の日次株価情報を作成
	testPrices := []models.DailyStockPrice{
		{
			PriceDate: time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: "7203",
				Price:   2873,
			},
		},
		{
			PriceDate: time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: "7203",
				Price:   2963,
			},
		},
		{
			PriceDate: time.Date(2025, 2, 6, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: "7203",
				Price:   2903.5,
			},
		},
	}

	// Act - テーブルを初期化
	err := InitializeDailyStockPriceTable(dbPath, testPrices)

	// Assert
	if err != nil {
		t.Fatalf("Failed to initialize table: %v", err)
	}

	// Act - テーブルからデータを取得
	retrievedPrices, err := GetDailyStockPrices(dbPath)

	// Assert
	if err != nil {
		t.Fatalf("Failed to get prices: %v", err)
	}

	// 取得したデータの数が正しいことを確認
	if len(retrievedPrices) != len(testPrices) {
		t.Errorf("Expected %d prices, but got %d", len(testPrices), len(retrievedPrices))
	}

	// 各エントリの内容を確認
	for i, expected := range testPrices {
		// 同じ日付のエントリを探す
		var found bool
		for _, actual := range retrievedPrices {
			// 日付は時刻部分を無視して比較
			expectedDate := expected.PriceDate.Format("2006-01-02")
			actualDate := actual.PriceDate.Format("2006-01-02")

			if expectedDate == actualDate &&
				expected.StockPrice.StockID == actual.StockPrice.StockID &&
				expected.StockPrice.Price == actual.StockPrice.Price {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Entry %d not found in retrieved data: %+v", i, expected)
		}
	}

	// 再度初期化して上書きできることを確認
	newTestPrices := []models.DailyStockPrice{
		{
			PriceDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: "7203",
				Price:   3000,
			},
		},
	}

	// Act - テーブルを再初期化
	err = InitializeDailyStockPriceTable(dbPath, newTestPrices)

	// Assert
	if err != nil {
		t.Fatalf("Failed to re-initialize table: %v", err)
	}

	// Act - テーブルからデータを再取得
	retrievedPrices, err = GetDailyStockPrices(dbPath)

	// Assert
	if err != nil {
		t.Fatalf("Failed to get prices after re-initialization: %v", err)
	}

	// 取得したデータの数が新しいデータセットの数と一致することを確認
	if len(retrievedPrices) != len(newTestPrices) {
		t.Errorf("After re-initialization: Expected %d prices, but got %d", len(newTestPrices), len(retrievedPrices))
	}
}
