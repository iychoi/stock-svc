package web_svc

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (svc *WebSVC) getFaangHTMLHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "WebSVC",
		"function": "getFaangHTMLHandler",
	})

	logger.Infof("Page access request from %s to %s", r.RemoteAddr, r.RequestURI)

	w.Header().Set("Content-Type", "text/html")

	// render header
	err := svc.writeHTMLHeader(w)
	if err != nil {
		logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}

	// Tickers?
	chartItems := []string{
		"GOOG",
		"FB",
		"AMZN",
		"AAPL",
		"NVDA",
		"TSLA",
		"NFLX",
		"BABA",
		"BIDU",
		"TWTR",
	}

	err = svc.renderChartMapDetailHTML(chartItems, w)
	if err != nil {
		logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}

	err = svc.writeHTMLFooter(w)
	if err != nil {
		logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}
}
