package file

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/toriwasa/sqlite-playground/internal/domain/models"
)

// 日付フォーマット定数
const dateFormat = "2006/1/2"

// ReadDailyStockPriceFromTSV は指定されたTSVファイルから日次株価情報を読み込みます。
// TSVファイルは「銘柄コード\t日付\t株価」の形式である必要があります。
//
// 引数:
//   - filePath: 読み込むTSVファイルのパス
//
// 戻り値:
//   - 日次株価情報の配列
//   - エラー（ファイル読み込みや解析に失敗した場合）
func ReadDailyStockPriceFromTSV(filePath string) ([]models.DailyStockPrice, error) {
	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 結果を格納するスライス
	var dailyPrices []models.DailyStockPrice

	// スキャナーを作成
	scanner := bufio.NewScanner(file)

	// 各行を読み込む
	for scanner.Scan() {
		line := scanner.Text()

		// 空行をスキップ
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// タブで分割
		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			return nil, &InvalidTSVFormatError{Line: line}
		}

		// 銘柄コード
		stockID := strings.TrimSpace(fields[0])

		// 日付を解析
		dateStr := strings.TrimSpace(fields[1])
		priceDate, err := time.Parse(dateFormat, dateStr)
		if err != nil {
			return nil, &InvalidDateFormatError{DateStr: dateStr, Line: line}
		}

		// 株価を解析
		priceStr := strings.TrimSpace(fields[2])
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return nil, &InvalidPriceFormatError{PriceStr: priceStr, Line: line}
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

	// スキャナーのエラーをチェック
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dailyPrices, nil
}

// InvalidTSVFormatError はTSVファイルのフォーマットが不正な場合のエラー
type InvalidTSVFormatError struct {
	Line string
}

func (e *InvalidTSVFormatError) Error() string {
	return "invalid TSV format: expected 3 fields separated by tabs: " + e.Line
}

// InvalidDateFormatError は日付のフォーマットが不正な場合のエラー
type InvalidDateFormatError struct {
	DateStr string
	Line    string
}

func (e *InvalidDateFormatError) Error() string {
	return "invalid date format: " + e.DateStr + " in line: " + e.Line
}

// InvalidPriceFormatError は株価のフォーマットが不正な場合のエラー
type InvalidPriceFormatError struct {
	PriceStr string
	Line     string
}

func (e *InvalidPriceFormatError) Error() string {
	return "invalid price format: " + e.PriceStr + " in line: " + e.Line
}
