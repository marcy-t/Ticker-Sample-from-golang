package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Ticker Sample")

	// 1秒間に1回起動するtickerを作成する
	// 具体的にはTicker.Cというchan Timeなチャンネルに指定間隔で変数が
	// 投げられるように
	ticker := time.NewTicker(1 * time.Second)

	// 終了処理用にbool型のチャネルを用意する
	done := make(chan bool, 1)

	// シグナル受信用のチャネルを用意する
	sigsChan := make(chan os.Signal, 1)

	// プログラムがシグナルを受信した際に、チャンネルに通知するように設定する
	// 具体的には、os.Interrupt(SIGINTのこと)をプログラムが受信したら
	// sigs変数のチャネルに通知するよう設定しています
	// signal.Notify(sigsChan, os.Interrupt)
	signal.Notify(sigsChan, syscall.SIGTERM)

	cnt := 0

	// 以下の関数を、goroutineとして非同期に実行を開始する
	go func() {
		for {
			// ticker.Cチャネルに何か来るまでforループ
			// 1秒間に1回、Time型の変数が投げられます
			<-ticker.C
			cnt += 1
			fmt.Println(cnt)
		}
	}()

	// 以下の関数を、goroutineとして非同期に実行開始する
	go func() {
		// sigs チャネルに何かが来るまでgoroutineを止める
		<-sigsChan // チャネルにCrtl + Cが来るまでここで処理が止まる

		// 標準出力に"bye"と表示する
		fmt.Println("bye")

		// 後始末としてtickerを止める
		ticker.Stop()

		// doneチャネルにtrueを投げる
		done <- true
	}()
	// doneチャネルに何かが来るまでプログラムを止める
	<-done
}
