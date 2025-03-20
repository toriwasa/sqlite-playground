package usecase

import (
	"errors"
	"math"
	"sort"

	"github.com/toriwasa/sqlite-playground/internal/domain/models"
)

// 株価情報が空の場合のエラーメッセージ
const ErrEmptyStockPricesMessage = "stock prices cannot be empty"

// 異なる銘柄コードが含まれる場合のエラーメッセージ
const ErrDifferentStockIDsMessage = "all stock prices must have the same stock ID"

// CalculateStockPriceStatistics は、n日分の株価情報から統計情報を計算します。
// 入力された株価情報が空の場合はエラーを返します。
// 入力された株価情報の銘柄コードが一致しない場合はエラーを返します。
//
// 引数:
//   - dailyPrices: n日分の株価情報
//
// 戻り値:
//   - 株価統計情報
//   - エラー（処理中に問題が発生した場合）
func CalculateStockPriceStatistics(dailyPrices []models.DailyStockPrice) (models.DailyStockPriceStatistics, error) {
	// 入力バリデーション
	if len(dailyPrices) == 0 {
		return models.DailyStockPriceStatistics{}, errors.New(ErrEmptyStockPricesMessage)
	}

	// 最初の銘柄コードを取得
	firstStockID := dailyPrices[0].StockPrice.StockID

	// すべての株価情報の銘柄コードが一致することを確認
	for _, price := range dailyPrices {
		if price.StockPrice.StockID != firstStockID {
			return models.DailyStockPriceStatistics{}, errors.New(ErrDifferentStockIDsMessage)
		}
	}

	// 日付でソートするためのスライスをコピー
	sortedPrices := make([]models.DailyStockPrice, len(dailyPrices))
	copy(sortedPrices, dailyPrices)

	// 日付でソート
	sort.Slice(sortedPrices, func(i, j int) bool {
		return sortedPrices[i].PriceDate.Before(sortedPrices[j].PriceDate)
	})

	// 開始日と終了日を取得
	startDate := sortedPrices[0].PriceDate
	endDate := sortedPrices[len(sortedPrices)-1].PriceDate

	// 統計情報の計算
	sum := 0.0
	max := sortedPrices[0].StockPrice.Price
	min := sortedPrices[0].StockPrice.Price

	// 合計、最大値、最小値を計算
	for _, price := range sortedPrices {
		currentPrice := price.StockPrice.Price
		sum += currentPrice

		if currentPrice > max {
			max = currentPrice
		}

		if currentPrice < min {
			min = currentPrice
		}
	}

	// 平均値を計算
	count := float64(len(sortedPrices))
	average := sum / count

	// 標準偏差を計算
	sumSquaredDiff := 0.0
	for _, price := range sortedPrices {
		diff := price.StockPrice.Price - average
		sumSquaredDiff += diff * diff
	}
	standardDeviation := math.Sqrt(sumSquaredDiff / count)

	// 結果を構築
	return models.DailyStockPriceStatistics{
		StartDate: startDate,
		EndDate:   endDate,
		StockPriceStatistics: models.StockPriceStatistics{
			StockID:           firstStockID,
			Average:           average,
			Max:               max,
			Min:               min,
			StandardDeviation: standardDeviation,
		},
	}, nil
}
