package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/iychoi/stock-svc/finance_svc"
	"github.com/iychoi/stock-svc/web_svc"
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
	log.Info("Starting Time Service...")
	timeSVC, err := finance_svc.InitTimeSVC()
	if err != nil {
		log.Fatal(err)
	}
	defer timeSVC.Close()
	log.Info("Time Service Started")

	log.Info("Starting Chart Service...")
	chartSVC, err := finance_svc.InitChartSVC(timeSVC)
	if err != nil {
		log.Fatal(err)
	}
	defer chartSVC.Close()
	log.Info("Chart Service Started")

	log.Info("Starting Price Service...")
	priceSVC, err := finance_svc.InitPriceSVC(timeSVC)
	if err != nil {
		log.Fatal(err)
	}
	defer priceSVC.Close()
	log.Info("Price Service Started")

	log.Info("Starting Web Service...")
	webSVC, err := web_svc.InitWebSVC(timeSVC, chartSVC, priceSVC)
	if err != nil {
		log.Fatal(err)
	}
	defer webSVC.Close()
	log.Info("Web Service Started")

	fmt.Println("Press Ctrl+C to stop server")
	waitForCtrlC()
}
