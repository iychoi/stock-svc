package web_svc

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iychoi/stock-svc/finance_svc"
	log "github.com/sirupsen/logrus"
)

const (
	timeLayout      string = "Jan/02/2006 15:04:05"
	timeLayoutNoSec string = "Jan 02 03:04 PM"

	serviceAddress string = ":80"

	headerHTMLFile = "resources/header.html"
	footerHTMLFile = "resources/footer.html"
)

// Server ...
type WebSVC struct {
	Router       *mux.Router
	TimeService  *finance_svc.TimeSVC
	ChartService *finance_svc.ChartSVC
	PriceService *finance_svc.PriceSVC

	WebServer *http.Server
}

// InitWebSVC ...
func InitWebSVC(timeService *finance_svc.TimeSVC, chartService *finance_svc.ChartSVC, priceService *finance_svc.PriceSVC) (*WebSVC, error) {
	logger := log.WithFields(log.Fields{
		"package":  "web_svc",
		"function": "InitWebSVC",
	})

	webSVC := &WebSVC{
		Router:       mux.NewRouter(),
		TimeService:  timeService,
		ChartService: chartService,
		PriceService: priceService,
		WebServer:    nil,
	}

	webSVC.addHandlers()

	server := &http.Server{
		Addr:    serviceAddress,
		Handler: webSVC.Router,
	}

	webSVC.WebServer = server

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Error(err)
		}
	}()

	return webSVC, nil
}

// Close ...
func (svc *WebSVC) Close() error {
	logger := log.WithFields(log.Fields{
		"package":  "web_svc",
		"function": "Close",
	})

	err := svc.WebServer.Close()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// AddHandlers ...
func (svc *WebSVC) addHandlers() {
	svc.Router.HandleFunc("/", svc.getIndexHTMLHandler).Methods("GET")
	svc.Router.HandleFunc("/index", svc.getIndexHTMLHandler).Methods("GET")
	svc.Router.HandleFunc("/etf", svc.getEtfHTMLHandler).Methods("GET")
	svc.Router.HandleFunc("/faang", svc.getFaangHTMLHandler).Methods("GET")
	svc.Router.HandleFunc("/semiconductor", svc.getSemiconductorHTMLHandler).Methods("GET")
	svc.Router.HandleFunc("/crypto", svc.getCryptoHTMLHandler).Methods("GET")
	svc.Router.HandleFunc("/future", svc.getFutureHTMLHandler).Methods("GET")
	svc.Router.HandleFunc("/growth", svc.getGrowthHTMLHandler).Methods("GET")
	svc.Router.HandleFunc("/basic", svc.getBasicHTMLHandler).Methods("GET")

	// stock images
	svc.Router.HandleFunc("/chartimg/{symbol}/{period}/{interval}", svc.getChartImageHandler).Methods("GET")
}

// writeHTMLHeader ...
func (svc *WebSVC) writeHTMLHeader(w io.Writer) error {
	logger := log.WithFields(log.Fields{
		"package":  "web_svc",
		"function": "writeHTMLHeader",
	})

	data, err := ioutil.ReadFile(headerHTMLFile)
	if err != nil {
		logger.Error(err)
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// writeHTMLFooter ...
func (svc *WebSVC) writeHTMLFooter(w io.Writer) error {
	logger := log.WithFields(log.Fields{
		"package":  "web_svc",
		"function": "writeHTMLFooter",
	})

	data, err := ioutil.ReadFile(footerHTMLFile)
	if err != nil {
		logger.Error(err)
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
