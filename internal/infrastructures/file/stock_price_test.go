package file

import (
	"reflect"
	"testing"
	"time"

	"github.com/toriwasa/sqlite-playground/internal/domain/models"
)

func TestReadDailyStockPriceFromTSV(t *testing.T) {
	// Arrange
	filePath := "../../data/sample_daily_stock_price.tsv"
	expectedFirstPrice := models.DailyStockPrice{
		PriceDate: time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC),
		StockPrice: models.StockPrice{
			StockID: "7203",
			Price:   2873,
		},
	}

	// Act
	dailyPrices, err := ReadDailyStockPriceFromTSV(filePath)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if len(dailyPrices) == 0 {
		t.Fatal("Expected non-empty result, but got empty slice")
	}

	// 最初のエントリをチェック
	if !reflect.DeepEqual(dailyPrices[0], expectedFirstPrice) {
		t.Errorf("First entry mismatch.\nExpected: %+v\nGot: %+v", expectedFirstPrice, dailyPrices[0])
	}

	// 全てのエントリが同じ銘柄コードを持っていることを確認
	stockID := dailyPrices[0].StockPrice.StockID
	for i, price := range dailyPrices {
		if price.StockPrice.StockID != stockID {
			t.Errorf("Entry %d has different stock ID: expected %s, got %s", i, stockID, price.StockPrice.StockID)
		}
	}
}
