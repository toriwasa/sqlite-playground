package usecase

import (
	"testing"
	"time"

	"github.com/toriwasa/sqlite-playground/internal/domain/models"
)

// TestCalculateStockPriceStatistics_ValidPrices は、有効な株価情報から正しい統計情報が計算されることをテストします。
func TestCalculateStockPriceStatistics_ValidPrices(t *testing.T) {
	// Arrange
	stockID := "1234"
	dailyPrices := []models.DailyStockPrice{
		{
			PriceDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: stockID,
				Price:   100.0,
			},
		},
		{
			PriceDate: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: stockID,
				Price:   110.0,
			},
		},
		{
			PriceDate: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: stockID,
				Price:   90.0,
			},
		},
	}

	expectedAverage := 100.0
	expectedMax := 110.0
	expectedMin := 90.0
	expectedStdDev := 8.16 // √((100-100)² + (110-100)² + (90-100)²) / 3 = √(0 + 100 + 100) / 3 = √200/3 ≈ 8.16

	// Act
	result, err := CalculateStockPriceStatistics(dailyPrices)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if result.StockID != stockID {
		t.Errorf("Expected StockID %s, but got %s", stockID, result.StockID)
	}

	if result.StartDate != dailyPrices[0].PriceDate {
		t.Errorf("Expected StartDate %v, but got %v", dailyPrices[0].PriceDate, result.StartDate)
	}

	if result.EndDate != dailyPrices[2].PriceDate {
		t.Errorf("Expected EndDate %v, but got %v", dailyPrices[2].PriceDate, result.EndDate)
	}

	if result.Average != expectedAverage {
		t.Errorf("Expected Average %f, but got %f", expectedAverage, result.Average)
	}

	if result.Max != expectedMax {
		t.Errorf("Expected Max %f, but got %f", expectedMax, result.Max)
	}

	if result.Min != expectedMin {
		t.Errorf("Expected Min %f, but got %f", expectedMin, result.Min)
	}

	// 浮動小数点の計算誤差を考慮して、標準偏差は近似値で比較
	stdDevDiff := result.StandardDeviation - expectedStdDev
	if stdDevDiff < -0.01 || stdDevDiff > 0.01 {
		t.Errorf("Expected StandardDeviation approximately %f, but got %f", expectedStdDev, result.StandardDeviation)
	}
}

// TestCalculateStockPriceStatistics_EmptyPrices は、空の株価情報リストでエラーが返されることをテストします。
func TestCalculateStockPriceStatistics_EmptyPrices(t *testing.T) {
	// Arrange
	var emptyPrices []models.DailyStockPrice

	// Act
	_, err := CalculateStockPriceStatistics(emptyPrices)

	// Assert
	if err == nil {
		t.Error("Expected an error for empty prices, but got nil")
	}

	if err.Error() != ErrEmptyStockPricesMessage {
		t.Errorf("Expected error message '%s', but got '%s'", ErrEmptyStockPricesMessage, err.Error())
	}
}

// TestCalculateStockPriceStatistics_DifferentStockIDs は、異なる銘柄コードの株価情報でエラーが返されることをテストします。
func TestCalculateStockPriceStatistics_DifferentStockIDs(t *testing.T) {
	// Arrange
	dailyPrices := []models.DailyStockPrice{
		{
			PriceDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: "1234",
				Price:   100.0,
			},
		},
		{
			PriceDate: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			StockPrice: models.StockPrice{
				StockID: "5678", // 異なる銘柄コード
				Price:   110.0,
			},
		},
	}

	// Act
	_, err := CalculateStockPriceStatistics(dailyPrices)

	// Assert
	if err == nil {
		t.Error("Expected an error for different stock IDs, but got nil")
	}

	if err.Error() != ErrDifferentStockIDsMessage {
		t.Errorf("Expected error message '%s', but got '%s'", ErrDifferentStockIDsMessage, err.Error())
	}
}

// TestCalculateStockPriceStatistics_SingleDay は、1日分の株価情報で正しく計算されることをテストします。
func TestCalculateStockPriceStatistics_SingleDay(t *testing.T) {
	// Arrange
	stockID := "1234"
	price := 100.0
	priceDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	dailyPrices := []models.DailyStockPrice{
		{
			PriceDate: priceDate,
			StockPrice: models.StockPrice{
				StockID: stockID,
				Price:   price,
			},
		},
	}

	// Act
	result, err := CalculateStockPriceStatistics(dailyPrices)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if result.StockID != stockID {
		t.Errorf("Expected StockID %s, but got %s", stockID, result.StockID)
	}

	if result.StartDate != priceDate {
		t.Errorf("Expected StartDate %v, but got %v", priceDate, result.StartDate)
	}

	if result.EndDate != priceDate {
		t.Errorf("Expected EndDate %v, but got %v", priceDate, result.EndDate)
	}

	if result.Average != price {
		t.Errorf("Expected Average %f, but got %f", price, result.Average)
	}

	if result.Max != price {
		t.Errorf("Expected Max %f, but got %f", price, result.Max)
	}

	if result.Min != price {
		t.Errorf("Expected Min %f, but got %f", price, result.Min)
	}

	if result.StandardDeviation != 0.0 {
		t.Errorf("Expected StandardDeviation 0.0, but got %f", result.StandardDeviation)
	}
}
