package controller

import (
	"fmt"
	"time"

	"github.com/toriwasa/sqlite-playground/internal/domain/models"
	"github.com/toriwasa/sqlite-playground/internal/infrastructures/db"
	"github.com/toriwasa/sqlite-playground/internal/usecase"
)

// GetStockPriceStatisticsByDateRange は指定された銘柄コードと日付範囲に一致する
// 日次株価情報の統計を計算します。
//
// 引数:
//   - dbPath: SQLiteデータベースファイルのパス
//   - stockID: 取得する銘柄コード
//   - startDate: 取得する日付の始点（この日付を含む）
//   - endDate: 取得する日付の終点（この日付を含む）
//
// 戻り値:
//   - 日次株価統計情報
//   - エラー（データ取得や計算に失敗した場合）
func GetStockPriceStatisticsByDateRange(dbPath string, stockID string, startDate time.Time, endDate time.Time) (models.DailyStockPriceStatistics, error) {
	// インフラストラクチャ層から日次株価情報を取得
	dailyPrices, err := db.GetDailyStockPricesByDateRange(dbPath, stockID, startDate, endDate)
	if err != nil {
		return models.DailyStockPriceStatistics{}, fmt.Errorf("failed to get daily stock prices: %w", err)
	}

	// データが存在しない場合のエラー処理
	if len(dailyPrices) == 0 {
		return models.DailyStockPriceStatistics{}, fmt.Errorf("no stock prices found for stock ID %s between %s and %s",
			stockID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}

	// ユースケース層で統計情報を計算
	statistics, err := usecase.CalculateStockPriceStatistics(dailyPrices)
	if err != nil {
		return models.DailyStockPriceStatistics{}, fmt.Errorf("failed to calculate statistics: %w", err)
	}

	return statistics, nil
}
