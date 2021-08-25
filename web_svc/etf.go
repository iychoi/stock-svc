package web_svc

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (svc *WebSVC) getEtfHTMLHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "WebSVC",
		"function": "getEtfHTMLHandler",
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
		"FNGU",
		"SOXL",
		"TQQQ",
		"UDOW",
		"UPRO",
		"URTY",
		"TECL",
		"LABU",
		"BNKU",
		"ICLN",
		"CURE",
		"KRBN",
		"JETS",
		"NRGU",
		"RETL",
		"DFEN",
		"KORU",
		"NAIL",
		"TPOR",
		"VTV",
		"DRN",
		"XLB",
		"DBB",
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
