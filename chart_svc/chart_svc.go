package chart_svc

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

const (
	stockChartBin     = "exec/stock_chart.py"
	stockChartFileDir = "charts"

	stockMonitoringExpTime = 7 * 24 * time.Hour // 1 week
	stockMonitoringTick    = 3 * time.Minute    // 3 min
)

type ChartPeriod string
type ChartInterval string

const (
	ChartPeriodMax    ChartPeriod = "max"
	ChartPeriod1Year  ChartPeriod = "1y"
	ChartPeriod2Year  ChartPeriod = "2y"
	ChartPeriod5Year  ChartPeriod = "5y"
	ChartPeriod10Year ChartPeriod = "10y"
	ChartPeriod6Month ChartPeriod = "6mo"
	ChartPeriod3Month ChartPeriod = "3mo"
	ChartPeriod1Month ChartPeriod = "1mo"
	ChartPeriod5Day   ChartPeriod = "5d"
	ChartPeriod1Day   ChartPeriod = "1d"

	ChartInteval1Min   ChartInterval = "1m"
	ChartInteval5Min   ChartInterval = "5m"
	ChartInteval30Min  ChartInterval = "30m"
	ChartInteval1Hour  ChartInterval = "1h"
	ChartInteval1Day   ChartInterval = "1d"
	ChartInteval5Day   ChartInterval = "5d"
	ChartInteval1Week  ChartInterval = "1wk"
	ChartInteval1Month ChartInterval = "1mo"
	ChartInteval3Month ChartInterval = "3mo"
)

type StockChartData struct {
	StockSymbol   string
	Period        ChartPeriod
	Interval      ChartInterval
	LocalFilePath string
}

// ChartSVC ...
type ChartSVC struct {
	// Stocks to be monitored
	Stocks *cache.Cache
	Ticker *time.Ticker
	Done   chan bool
}

func InitChartSVC() (*ChartSVC, error) {
	stockCache := cache.New(stockMonitoringExpTime, stockMonitoringExpTime)
	ticker := time.NewTicker(stockMonitoringTick)
	done := make(chan bool)

	chartSvc := &ChartSVC{
		Stocks: stockCache,
		Ticker: ticker,
		Done:   done,
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// tick
				chartSvc.renewChartFiles()
			}
		}
	}()

	return chartSvc, nil
}

// Close ...
func (chart *ChartSVC) Close() error {
	chart.Stocks.Flush()

	chart.Ticker.Stop()
	chart.Done <- true
	return nil
}

func (chart *ChartSVC) renewChartFiles() {
	logger := log.WithFields(log.Fields{
		"package":  "chart_svc",
		"function": "renewChartFiles",
	})

	for _, item := range chart.Stocks.Items() {
		chartData := item.Object.(*StockChartData)
		err := chart.makeChart(chartData.StockSymbol, chartData.Period, chartData.Interval, false)
		if err != nil {
			logger.Error(err)
		}
	}
}

// MakeChartFile ...
func (chart *ChartSVC) RequestChart(symbol string, period ChartPeriod, interval ChartInterval) error {
	logger := log.WithFields(log.Fields{
		"package":  "chart_svc",
		"function": "RequestChart",
	})

	logger.Infof("Request chart %s", symbol)
	return chart.makeChart(symbol, period, interval, true)
}

func (chart *ChartSVC) makeChart(symbol string, period ChartPeriod, interval ChartInterval, updateCache bool) error {
	logger := log.WithFields(log.Fields{
		"package":  "chart_svc",
		"function": "makeChart",
	})

	logger.Infof("Make chart %s", symbol)

	filename := makeChartFileName(symbol, period, interval)
	filepath := fmt.Sprintf("%s/%s", stockChartFileDir, filename)

	err := os.MkdirAll(stockChartFileDir, 0766)
	if err != nil {
		logger.Error(err)
		return err
	}

	args := []string{
		symbol,
		string(period),
		string(interval),
		filepath,
	}

	_, err = executeScript(stockChartBin, args)
	if err != nil {
		logger.Error(err)
		return err
	}

	chartData := StockChartData{
		StockSymbol:   symbol,
		Period:        period,
		Interval:      interval,
		LocalFilePath: filepath,
	}

	if updateCache {
		chart.Stocks.SetDefault(symbol, &chartData)
	}
	return nil
}

func makeChartFileName(symbol string, period ChartPeriod, interval ChartInterval) string {
	safeSymbol := strings.TrimPrefix(symbol, "^")
	return fmt.Sprintf("%s_%s_%s.png", safeSymbol, period, interval)
}

func executeScript(bin string, args []string) ([]byte, error) {
	logger := log.WithFields(log.Fields{
		"package":  "chart_svc",
		"function": "executeScript",
	})

	logger.Infof("Executing exec (%s) with arguments (%v)", bin, args)
	command := exec.Command(bin, args...)
	output, err := command.CombinedOutput()
	if err != nil {
		logger.Errorf("exec failed: %v\nCommand: %s\nArguments: %s\nOutput: %s\n", err, bin, args, string(output))
		return nil, fmt.Errorf("exec failed: %v\nCommand: %s\nArguments: %s\nOutput: %s", err, bin, args, string(output))
	}
	return output, err
}
