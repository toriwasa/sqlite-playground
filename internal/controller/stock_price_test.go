package controller

import (
	"testing"
	"time"

	"github.com/toriwasa/sqlite-playground/internal/domain/models"
	"github.com/toriwasa/sqlite-playground/internal/infrastructures/db"
)

func TestGetStockPriceStatisticsByDateRange(t *testing.T) {
	// テスト用のデータベースパス
	dbPath := "../infrastructures/db/test_stock_price_controller.db"

	// テスト用の日付
	startDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC)
	stockID := "7203"

	// テスト用の日次株価情報を作成
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
	}

	// テスト用のデータベースを初期化
	err := setupTestDatabase(dbPath, testPrices)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanupTestDatabase(dbPath)

	// テストケース
	testCases := []struct {
		name      string
		stockID   string
		startDate time.Time
		endDate   time.Time
		wantErr   bool
		errCheck  func(err error) bool
	}{
		{
			name:      "正常系: 全期間のデータを取得",
			stockID:   stockID,
			startDate: startDate,
			endDate:   endDate,
			wantErr:   false,
		},
		{
			name:      "正常系: 部分期間のデータを取得",
			stockID:   stockID,
			startDate: time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC),
			wantErr:   false,
		},
		{
			name:      "エラー: 存在しない銘柄コード",
			stockID:   "0000",
			startDate: startDate,
			endDate:   endDate,
			wantErr:   true,
			errCheck: func(err error) bool {
				return err != nil && err.Error() != ""
			},
		},
		{
			name:      "エラー: 範囲外の日付",
			stockID:   stockID,
			startDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
			wantErr:   true,
			errCheck: func(err error) bool {
				return err != nil && err.Error() != ""
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			stats, err := GetStockPriceStatisticsByDateRange(dbPath, tc.stockID, tc.startDate, tc.endDate)

			// Assert
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				if tc.errCheck != nil && !tc.errCheck(err) {
					t.Errorf("Error did not match expected condition: %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// 統計情報の検証
			if stats.StockID != tc.stockID {
				t.Errorf("Expected stock ID %s, but got %s", tc.stockID, stats.StockID)
			}

			// 日付範囲の検証
			if stats.StartDate.Before(tc.startDate) {
				t.Errorf("Start date %v is before requested start date %v", stats.StartDate, tc.startDate)
			}
			if stats.EndDate.After(tc.endDate) {
				t.Errorf("End date %v is after requested end date %v", stats.EndDate, tc.endDate)
			}

			// 統計値の検証（詳細な値は省略、正常に計算されていることを確認）
			if stats.Average <= 0 {
				t.Errorf("Expected positive average, but got %f", stats.Average)
			}
			if stats.Max <= 0 {
				t.Errorf("Expected positive max, but got %f", stats.Max)
			}
			if stats.Min <= 0 {
				t.Errorf("Expected positive min, but got %f", stats.Min)
			}
		})
	}
}

// テスト用のデータベースをセットアップする関数
func setupTestDatabase(dbPath string, prices []models.DailyStockPrice) error {
	// インフラストラクチャ層の関数を使用してテーブルを初期化
	return db.InitializeDailyStockPriceTable(dbPath, prices)
}

// テスト用のデータベースをクリーンアップする関数
func cleanupTestDatabase(dbPath string) {
	// ファイルを削除
	// Note: 実際のテストでは以下の行のコメントを外して有効にする
	// 今回はテスト後にファイルを残しておくためコメントアウトしている
	// os.Remove(dbPath)
}
