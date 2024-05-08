package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {

	// ログの設定
	//LogSettings("log.txt")
	LogWrite("Start", nil)

	// 入力ファイル名
	flag.Parse()
	textFilePath := flag.Arg(0)

	// 出力ファイル名
	excelDir, excelFilePath := filepath.Split(flag.Arg(0))
	pos := strings.LastIndex(excelFilePath, ".")
	excelFileName := excelDir + excelFilePath[:pos] + ".xlsx"

	//テキストファイルの読み込み
	textFile, err := os.Open(textFilePath)
	LogFatal("テキストファイルの読み込みができませんでした。", err)
	defer textFile.Close()

	reader := csv.NewReader(transform.NewReader(textFile, japanese.ShiftJIS.NewDecoder()))
	reader.Comma = '\t'

	// excelファイルの作成
	excelFile := excelize.NewFile()
	defer func() {
		err := excelFile.Close()
		LogFatal("エクセルファイルがクローズできない。", err)
	}()

	// ワークシートの設定
	excelFile.SetDefaultFont("游ゴシック")
	err = excelFile.SetSheetName("Sheet1", "データ") //デフォルトのシート名を"データ"に変更
	LogFatal("シート名が変更できませんでした。", err)

	// ストリームライターの設定(大容量ファイル対応の為)
	streamWriter, err := excelFile.NewStreamWriter("データ")
	LogFatal("ストリームライターの設定ができませんでした。", err)

	rowCount := 1
	for {
		textItems, err := reader.Read() //1行読みだす
		// _, err := reader.Read() //1行読みだす
		if err == io.EOF {
			break
		} else {
			LogFatal("テキストファイルの読み込みエラー", err)
		}

		// ストリームライターに１行書きだす
		// var rowItems []interface{}
		rowItems := make([]interface{}, len(textItems), len(textItems))
		for i, v := range textItems {
			//rowItems = append(rowItems, v)
			if v != "" {
				rowItems[i] = v
			}
		}
		err = streamWriter.SetRow(fmt.Sprintf("A%d", rowCount), rowItems)
		LogFatal("行に値を設定できませんでした。", err)

		// 1000件毎にログを表示する
		if rowCount % 1000 == 0 {
			LogWrite("処理中", fmt.Errorf("%d...", rowCount))
		}

		rowCount++

	}

	err = streamWriter.Flush()
	LogFatal("ストリームライターを終了できませんでした。", err)

	err = excelFile.SaveAs(excelFileName)
	LogFatal("エクセルファイルを保存できませんでした。", err)

	LogWrite("Finish!", nil)
}
