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
