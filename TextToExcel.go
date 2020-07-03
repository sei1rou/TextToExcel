package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {
	flag.Parse()

	// ログファイル準備
	logfile, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	failOnError(err)
	defer logfile.Close()

	log.SetOutput(logfile)

	//log.SetOutput(os.Stdout)

	log.Print("Start\r\n")

	// ファイルを読み込んで二次元配列に入れる
	records := readfile(flag.Arg(0))

	// ファイルをエクセルに出力
	saveExcel(flag.Arg(0), records)

}

func readfile(filename string) [][]string {
	//入力ファイル準備
	infile, err := os.Open(filename)
	failOnError(err)
	defer infile.Close()

	reader := csv.NewReader(transform.NewReader(infile, japanese.ShiftJIS.NewDecoder()))
	reader.Comma = '\t'

	//CSVファイルを２次元配列に展開

	readrecords := make([][]string, 0)
	record, err := reader.Read() // 1行読み出す
	if err == io.EOF {
		return readrecords
	} else {
		failOnError(err)
	}
	colMax := len(record) - 1
	readrecords = append(readrecords, record[:colMax])

	for {
		record, err := reader.Read()[:colMax] // 1行読み出す
		if err == io.EOF {
			break
		} else {
			log.Print(record)
			log.Print(len(record))
			failOnError(err)
		}

		readrecords = append(readrecords, record[:colMax])
	}

	return readrecords
}

func saveExcel(filename string, recs [][]string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var err error

	//出力ファイル準備
	outDir, outfileName := filepath.Split(filename)
	pos := strings.LastIndex(outfileName, ".")
	outExcelName := outDir + outfileName[:pos] + ".xlsx"

	file = xlsx.NewFile()
	xlsx.SetDefaultFont(11, "ＭＳ Ｐゴシック") // デフォルトのフォントを指定
	sheet, err = file.AddSheet("データ")
	failOnError(err)

	for r, recRow := range recs {
		for c, recCell := range recRow {
			sheet.Cell(r, c).Value = recCell
		}
	}

	err = file.Save(outExcelName)
	failOnError(err)

	log.Print("Finesh !\r\n")

}
