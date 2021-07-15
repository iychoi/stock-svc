package finance_svc

import (
	"fmt"
	"io/ioutil"
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

	stockMonitoringExpTime = 48 * time.Hour   // 2 days
	stockMonitoringTickMin = 15 * time.Minute // 15 min
	stockMonitoringTickDay = 30 * time.Minute // 30 min
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
	TimeService *TimeSVC
	// Charts to be monitored
	Charts              *cache.Cache
	MonitoringTickerMin *time.Ticker
	MonitoringTickerDay *time.Ticker
	MonitoringDone      chan bool
}

func InitChartSVC(timeService *TimeSVC) (*ChartSVC, error) {
	chartCache := cache.New(stockMonitoringExpTime, stockMonitoringExpTime)
	tickerMin := time.NewTicker(stockMonitoringTickMin)
	tickerDay := time.NewTicker(stockMonitoringTickDay)
	done := make(chan bool)

	chartSvc := &ChartSVC{
		TimeService:         timeService,
		Charts:              chartCache,
		MonitoringTickerMin: tickerMin,
		MonitoringTickerDay: tickerDay,
		MonitoringDone:      done,
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case <-tickerMin.C:
				// tick
				if timeService.GetMarketType(time.Now()) != Overnight {
					go chartSvc.renewChartsMinutes()
				}
			case <-tickerDay.C:
				// tick
				if timeService.GetMarketType(time.Now()) != Overnight {
					chartSvc.renewChartsDays()
				}
			}
		}
	}()

	return chartSvc, nil
}

// Close ...
func (svc *ChartSVC) Close() error {
	svc.Charts.Flush()

	svc.MonitoringTickerMin.Stop()
	svc.MonitoringTickerDay.Stop()
	svc.MonitoringDone <- true
	return nil
}

// GetChartData ...
func (svc *ChartSVC) GetChartData(symbol string, period ChartPeriod, interval ChartInterval) ([]byte, error) {
	logger := log.WithFields(log.Fields{
		"package":  "chart_svc",
		"function": "GetChartData",
	})

	err := svc.RequestChart(symbol, period, interval)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if chartCache, ok := svc.getChartCache(symbol, period, interval); ok {
		return ioutil.ReadFile(chartCache.LocalFilePath)
	} else {
		return nil, fmt.Errorf("could not get chart cache - %s, %s, %s", symbol, period, interval)
	}
}

// RequestChart ...
func (svc *ChartSVC) RequestChart(symbol string, period ChartPeriod, interval ChartInterval) error {
	logger := log.WithFields(log.Fields{
		"package":  "chart_svc",
		"function": "RequestChart",
	})

	logger.Infof("Request chart %s", symbol)

	if _, ok := svc.getChartCache(symbol, period, interval); !ok {
		return svc.makeChart(symbol, period, interval, true)
	} else {
		svc.renewChartCache(symbol, period, interval)
		return nil
	}
}

func (svc *ChartSVC) getChartCache(symbol string, period ChartPeriod, interval ChartInterval) (*StockChartData, bool) {
	filename := svc.makeChartFileName(symbol, period, interval)
	data, ok := svc.Charts.Get(filename)
	if ok {
		return data.(*StockChartData), ok
	} else {
		return nil, ok
	}
}

func (svc *ChartSVC) renewChartCache(symbol string, period ChartPeriod, interval ChartInterval) {
	filename := svc.makeChartFileName(symbol, period, interval)
	filepath := fmt.Sprintf("%s/%s", stockChartFileDir, filename)

	chartData := StockChartData{
		StockSymbol:   symbol,
		Period:        period,
		Interval:      interval,
		LocalFilePath: filepath,
	}

	// renew cache
	svc.Charts.SetDefault(filename, &chartData)
}

func (svc *ChartSVC) renewCharts() {
	logger := log.WithFields(log.Fields{
		"package":  "chart_svc",
		"function": "renewCharts",
	})

	for _, item := range svc.Charts.Items() {
		chartData := item.Object.(*StockChartData)
		err := svc.makeChart(chartData.StockSymbol, chartData.Period, chartData.Interval, false)
		if err != nil {
			logger.Error(err)
		}
	}
}

func (svc *ChartSVC) renewChartsMinutes() {
	logger := log.WithFields(log.Fields{
		"package":  "chart_svc",
		"function": "renewChartsMinutes",
	})

	for _, item := range svc.Charts.Items() {
		chartData := item.Object.(*StockChartData)
		if svc.isShortInterval(chartData.Interval) {
			err := svc.makeChart(chartData.StockSymbol, chartData.Period, chartData.Interval, false)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func (svc *ChartSVC) renewChartsDays() {
	logger := log.WithFields(log.Fields{
		"package":  "chart_svc",
		"function": "renewChartsDays",
	})

	for _, item := range svc.Charts.Items() {
		chartData := item.Object.(*StockChartData)
		if svc.isLongInterval(chartData.Interval) {
			err := svc.makeChart(chartData.StockSymbol, chartData.Period, chartData.Interval, false)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func (svc *ChartSVC) makeChart(symbol string, period ChartPeriod, interval ChartInterval, updateCache bool) error {
	logger := log.WithFields(log.Fields{
		"package":  "chart_svc",
		"function": "makeChart",
	})

	logger.Infof("Make chart %s", symbol)

	filename := svc.makeChartFileName(symbol, period, interval)
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

	_, err = svc.executeScript(stockChartBin, args)
	if err != nil {
		logger.Error(err)
		return err
	}

	if updateCache {
		// renew chart cache
		svc.renewChartCache(symbol, period, interval)
	}
	return nil
}

func (svc *ChartSVC) isShortInterval(interval ChartInterval) bool {
	if interval == ChartInteval1Min || interval == ChartInteval5Min {
		return true
	}
	return false
}

func (svc *ChartSVC) isLongInterval(interval ChartInterval) bool {
	if interval == ChartInteval1Min || interval == ChartInteval5Min {
		return false
	}
	return true
}

func (svc *ChartSVC) makeChartFileName(symbol string, period ChartPeriod, interval ChartInterval) string {
	safeSymbol := strings.TrimPrefix(symbol, "^")
	return fmt.Sprintf("%s_%s_%s.png", safeSymbol, period, interval)
}

func (svc *ChartSVC) executeScript(bin string, args []string) ([]byte, error) {
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
