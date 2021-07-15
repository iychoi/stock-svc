package web_svc

import (
	"fmt"
	"io"
	"text/template"

	"github.com/leekchan/accounting"
	log "github.com/sirupsen/logrus"
)

const (
	chartMapHTMLFile       = "resources/chartmap.html"
	chartMapDetailHTMLFile = "resources/chartmap_detail.html"
)

type TemplateStockChartItem struct {
	Symbol              string
	StockName           string
	CurrentPrice        string
	PriceChange         string
	PriceChangePercent  string
	PriceChangePositive bool
}

type TemplateStockChartItems struct {
	Items []TemplateStockChartItem
}

// renderChartMapHTML ...
func (svc *WebSVC) renderChartMapHTML(chartItems []string, w io.Writer) error {
	logger := log.WithFields(log.Fields{
		"package":  "web_svc",
		"function": "renderChartMapHTML",
	})

	t, err := template.ParseFiles(chartMapHTMLFile)
	if err != nil {
		logger.Error(err)
		return err
	}

	// convert data
	dataItems := []TemplateStockChartItem{}

	for _, symbol := range chartItems {

		stockInfo, err := svc.PriceService.GetStockInfo(symbol)
		if err != nil {
			logger.Error(err)
			continue
		}

		changePositive := true
		if stockInfo.PriceChange < 0 {
			changePositive = false
		}

		dataItem := TemplateStockChartItem{
			Symbol:              symbol,
			StockName:           stockInfo.StockName,
			CurrentPrice:        "",
			PriceChange:         "",
			PriceChangePercent:  "",
			PriceChangePositive: changePositive,
		}

		ac := accounting.Accounting{
			Symbol:    "$",
			Precision: 2,
		}

		dataItem.CurrentPrice = ac.FormatMoney(stockInfo.CurrentPrice)
		if stockInfo.PriceChange > 0 {
			dataItem.PriceChange = fmt.Sprintf("+%.2f", stockInfo.PriceChange)
		} else {
			dataItem.PriceChange = fmt.Sprintf("%.2f", stockInfo.PriceChange)
		}

		if stockInfo.PriceChangePercent > 0 {
			dataItem.PriceChangePercent = fmt.Sprintf("+%.2f%%", stockInfo.PriceChangePercent*100)
		} else {
			dataItem.PriceChangePercent = fmt.Sprintf("%.2f%%", stockInfo.PriceChangePercent*100)
		}

		dataItems = append(dataItems, dataItem)
	}

	data := TemplateStockChartItems{
		Items: dataItems,
	}

	t.Execute(w, data)
	return nil
}

// renderChartMapDetailHTML ...
func (svc *WebSVC) renderChartMapDetailHTML(chartItems []string, w io.Writer) error {
	logger := log.WithFields(log.Fields{
		"package":  "web_svc",
		"function": "renderChartMapDetailHTML",
	})

	t, err := template.ParseFiles(chartMapDetailHTMLFile)
	if err != nil {
		logger.Error(err)
		return err
	}

	// convert data
	dataItems := []TemplateStockChartItem{}

	for _, symbol := range chartItems {

		stockInfo, err := svc.PriceService.GetStockInfo(symbol)
		if err != nil {
			logger.Error(err)
			continue
		}

		changePositive := true
		if stockInfo.PriceChange < 0 {
			changePositive = false
		}

		dataItem := TemplateStockChartItem{
			Symbol:              symbol,
			StockName:           stockInfo.StockName,
			CurrentPrice:        "",
			PriceChange:         "",
			PriceChangePercent:  "",
			PriceChangePositive: changePositive,
		}

		ac := accounting.Accounting{
			Symbol:    "$",
			Precision: 2,
		}

		dataItem.CurrentPrice = ac.FormatMoney(stockInfo.CurrentPrice)
		if stockInfo.PriceChange > 0 {
			dataItem.PriceChange = fmt.Sprintf("+%.2f", stockInfo.PriceChange)
		} else {
			dataItem.PriceChange = fmt.Sprintf("%.2f", stockInfo.PriceChange)
		}

		if stockInfo.PriceChangePercent > 0 {
			dataItem.PriceChangePercent = fmt.Sprintf("+%.2f%%", stockInfo.PriceChangePercent*100)
		} else {
			dataItem.PriceChangePercent = fmt.Sprintf("%.2f%%", stockInfo.PriceChangePercent*100)
		}

		dataItems = append(dataItems, dataItem)
	}

	data := TemplateStockChartItems{
		Items: dataItems,
	}

	t.Execute(w, data)
	return nil
}
