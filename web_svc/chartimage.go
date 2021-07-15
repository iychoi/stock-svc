package web_svc

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iychoi/stock-svc/finance_svc"
	log "github.com/sirupsen/logrus"
)

func (svc *WebSVC) getChartImageHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "web_svc",
		"function": "getChartImageHandler",
	})

	varMap := mux.Vars(r)
	symbol, ok := varMap["symbol"]
	if !ok {
		w.WriteHeader(500)
		return
	}

	period, ok := varMap["period"]
	if !ok {
		w.WriteHeader(500)
		return
	}

	interval, ok := varMap["interval"]
	if !ok {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "image/png")

	bytes, err := svc.ChartService.GetChartData(symbol, finance_svc.ChartPeriod(period), finance_svc.ChartInterval(interval))
	if err != nil {
		logger.Error(err)
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(500)
		return
	}
}
