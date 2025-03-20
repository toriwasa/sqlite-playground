package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/glebarez/go-sqlite" // SQLiteドライバ
	"github.com/toriwasa/sqlite-playground/internal/domain/models"
)

// 日次株価テーブル名
const dailyStockPriceTableName = "daily_stock_price"

// 日次株価テーブル作成SQL
const createDailyStockPriceTableSQL = `
CREATE TABLE IF NOT EXISTS daily_stock_price (
    stock_id TEXT NOT NULL,
    price_date TEXT NOT NULL,
    price REAL NOT NULL,
    PRIMARY KEY (stock_id, price_date)
);
`

// InitializeDailyStockPriceTable はSQLiteのdaily_stock_priceテーブルを
// 引数で渡された日次株価情報配列で初期化します。
// テーブルが存在しない場合は作成し、存在する場合は全てのデータを削除してから
// 新しいデータを挿入します。
//
// 引数:
//   - dbPath: SQLiteデータベースファイルのパス
//   - dailyPrices: 挿入する日次株価情報の配列
//
// 戻り値:
//   - エラー（データベース操作に失敗した場合）
func InitializeDailyStockPriceTable(dbPath string, dailyPrices []models.DailyStockPrice) error {
	// データベース接続を開く
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// テーブルを作成
	_, err = db.Exec(createDailyStockPriceTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// トランザクションを開始
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 既存のデータを削除
	_, err = tx.Exec("DELETE FROM " + dailyStockPriceTableName)
	if err != nil {
		return fmt.Errorf("failed to delete existing data: %w", err)
	}

	// Prepared Statementを作成
	stmt, err := tx.Prepare("INSERT INTO " + dailyStockPriceTableName + " (stock_id, price_date, price) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// 各日次株価情報をテーブルに挿入
	for _, dailyPrice := range dailyPrices {
		// 日付をISO 8601形式の文字列に変換
		dateStr := dailyPrice.PriceDate.Format(time.RFC3339[:10]) // YYYY-MM-DD形式

		_, err = stmt.Exec(
			dailyPrice.StockPrice.StockID,
			dateStr,
			dailyPrice.StockPrice.Price,
		)
		if err != nil {
			return fmt.Errorf("failed to insert data: %w", err)
		}
	}

	// トランザクションをコミット
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetDailyStockPrices はSQLiteのdaily_stock_priceテーブルから
// 全ての日次株価情報を取得します。
//
// 引数:
//   - dbPath: SQLiteデータベースファイルのパス
//
// 戻り値:
//   - 日次株価情報の配列
//   - エラー（データベース操作に失敗した場合）
func GetDailyStockPrices(dbPath string) ([]models.DailyStockPrice, error) {
	// データベース接続を開く
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// クエリを実行
	rows, err := db.Query("SELECT stock_id, price_date, price FROM " + dailyStockPriceTableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	// 結果を格納するスライス
	var dailyPrices []models.DailyStockPrice

	// 各行を処理
	for rows.Next() {
		var stockID string
		var dateStr string
		var price float64

		// 行のデータを取得
		err := rows.Scan(&stockID, &dateStr, &price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// 日付文字列をtime.Time型に変換
		priceDate, err := time.Parse(time.RFC3339[:10], dateStr) // YYYY-MM-DD形式
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		// 日次株価情報を作成
		dailyPrice := models.DailyStockPrice{
			PriceDate: priceDate,
			StockPrice: models.StockPrice{
				StockID: stockID,
				Price:   price,
			},
		}

		// 結果に追加
		dailyPrices = append(dailyPrices, dailyPrice)
	}

	// エラーをチェック
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}

	return dailyPrices, nil
}
