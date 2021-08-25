package web_svc

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (svc *WebSVC) getCryptoHTMLHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "WebSVC",
		"function": "getCryptoHTMLHandler",
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
		"BTC-USD",
		"ETH-USD",
		"DOGE-USD",
		"XRP-USD",
		"ADA-USD",
		"BNB-USD",
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
