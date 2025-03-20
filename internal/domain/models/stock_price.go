package models

import (
	"time"
)

// 特定時点の株価情報を示す構造体
type StockPrice struct {
	// 銘柄コード文字列
	StockID string
	// 株価
	Price float64
}

// 日次の株価情報を示す構造体
type DailyStockPrice struct {
	// 株価が記録されている日付
	PriceDate time.Time
	// 株価情報
	StockPrice
}

// n個の株価情報の統計値を示す構造体
type StockPriceStatistics struct {
	// 株価情報のリスト
	StockPrices []StockPrice
	// 平均値
	Average float64
	// 最大値
	Max float64
	// 最小値
	Min float64
	// 標準偏差
	StandardDeviation float64
}
