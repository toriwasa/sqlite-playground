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

func TestGetDailyStockPricesByDateRange(t *testing.T) {
	// Arrange
	dbPath := "./test_stock_price_range.db"
	// テスト終了後にデータベースファイルを削除
	defer os.Remove(dbPath)

	// テスト用の日次株価情報を作成
	stockID := "7203"
	testPrices := []models.DailyStockPrice{
		{
			PriceDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: stockID,
				Price:   2800,
			},
		},
		{
			PriceDate: time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: stockID,
				Price:   2850,
			},
		},
		{
			PriceDate: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: stockID,
				Price:   2900,
			},
		},
		{
			PriceDate: time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: stockID,
				Price:   2950,
			},
		},
		{
			PriceDate: time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: stockID,
				Price:   3000,
			},
		},
		// 別の銘柄コードのデータも追加
		{
			PriceDate: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: "9984",
				Price:   5000,
			},
		},
	}

	// テーブルを初期化
	err := InitializeDailyStockPriceTable(dbPath, testPrices)
	if err != nil {
		t.Fatalf("Failed to initialize table: %v", err)
	}

	// テストケース
	testCases := []struct {
		name      string
		stockID   string
		startDate time.Time
		endDate   time.Time
		expected  int // 期待される結果の数
	}{
		{
			name:      "正常系: 範囲内のデータを取得",
			stockID:   stockID,
			startDate: time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC),
			expected:  3, // 2/2, 2/3, 2/4の3件
		},
		{
			name:      "境界値: 開始日と終了日に一致するデータを含む",
			stockID:   stockID,
			startDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			expected:  1, // 2/1の1件
		},
		{
			name:      "該当なし: 条件に一致するデータがない",
			stockID:   "0000", // 存在しない銘柄コード
			startDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC),
			expected:  0, // 該当なし
		},
		{
			name:      "別の銘柄コード: 特定の銘柄コードのみ取得",
			stockID:   "9984",
			startDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC),
			expected:  1, // 9984の2/3の1件
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			retrievedPrices, err := GetDailyStockPricesByDateRange(dbPath, tc.stockID, tc.startDate, tc.endDate)

			// Assert
			if err != nil {
				t.Fatalf("Failed to get prices by date range: %v", err)
			}

			// 取得したデータの数が期待値と一致することを確認
			if len(retrievedPrices) != tc.expected {
				t.Errorf("Expected %d prices, but got %d", tc.expected, len(retrievedPrices))
			}

			// 取得したデータが指定した条件に一致することを確認
			for _, price := range retrievedPrices {
				// 銘柄コードが一致することを確認
				if price.StockPrice.StockID != tc.stockID {
					t.Errorf("Expected stock ID %s, but got %s", tc.stockID, price.StockPrice.StockID)
				}

				// 日付が範囲内であることを確認
				priceDate := price.PriceDate
				if priceDate.Before(tc.startDate) || priceDate.After(tc.endDate) {
					t.Errorf("Price date %v is outside the range [%v, %v]", priceDate, tc.startDate, tc.endDate)
				}
			}
		})
	}
}
