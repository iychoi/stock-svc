package web_svc

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (svc *WebSVC) getIndexImageHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "WebSVC",
		"function": "getIndexImageHandler",
	})

	varMap := mux.Vars(r)
	index, ok := varMap["index"]
	if !ok {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "image/png")

	bytes, err := svc.FeerGreedIndexService.GetIndexData(index)
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
