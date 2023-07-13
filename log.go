package main

import (
	"log"
	"os"
)

func LogSettings(filename string) {
	// ログファイルを初期設定する
	logfile, err := os.OpenFile("./"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalln("ログファイルを作成できませんでした", err)
	}
	log.SetOutput(logfile)

}

func LogWrite(value string, err error) {
	// ログファイルに書き込む
	if err != nil {
		log.Printf("%s: %s", value, err)
	} else {
		log.Printf("%s", value)
	}
}

func LogFatal(value string, err error) {
	// ログファイルに書き込み終了する
	if err != nil {
		log.Fatalf("%s: %s", value, err)
	}
}
