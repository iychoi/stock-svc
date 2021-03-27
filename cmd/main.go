package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/iychoi/stock-svc/chart_svc"
	log "github.com/sirupsen/logrus"
)

func waitForCtrlC() {
	var endWaiter sync.WaitGroup

	endWaiter.Add(1)
	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		<-signalChannel
		endWaiter.Done()
	}()

	endWaiter.Wait()
}

func main() {
	log.Info("Starting Chart Service...")
	chartSVC, err := chart_svc.InitChartSVC()
	if err != nil {
		log.Fatal(err)
	}
	defer chartSVC.Close()
	log.Info("Chart Service Started")

	chartSVC.RequestChart("TSLA", "1mo", "1d")
	chartSVC.RequestChart("FNGU", "1mo", "1d")

	fmt.Println("Press Ctrl+C to stop server")
	waitForCtrlC()
}
