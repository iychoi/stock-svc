package web_svc

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (svc *WebSVC) getGrowthHTMLHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "web_svc",
		"function": "getGrowthHTMLHandler",
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
		"U",
		"PYPL",
		"PLTR",
		"DOCU",
		"SNAP",
		"TDOC",
		"ADBE",
		"ROKU",
		"SPOT",
		"ETSY",
		"ZG",
		"EXPE",
		"ABNB",
		"UBER",
		"DIS",
		"SNOW",
		"COIN",
		"AGC",
		"CHPT",
		"PAYC",
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
