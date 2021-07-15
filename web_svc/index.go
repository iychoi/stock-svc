package web_svc

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (svc *WebSVC) getIndexHTMLHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "web_svc",
		"function": "getIndexHTMLHandler",
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
		"^TNX",
		"DX-Y.NYB",
		"^VIX",
		"^GSPC",
		"^DJI",
		"^IXIC",
		"^RUT",
		"^KS11",
		"BTC-USD",
		"ETH-USD",
		"FNGU",
		"SOXL",
		"TQQQ",
		"UPRO",
		"URTY",
		"TECL",
		"LABU",
		"ICLN",
		"UVXY",
	}

	err = svc.renderChartMapHTML(chartItems, w)
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
