package finance_svc

import (
	"time"

	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	yahoo "github.com/tonymackay/go-yahoo-finance"
)

const (
	stockInfoCacheTimeout = 5 * time.Minute // 5 min
)

type StockInfo struct {
	Symbol             string
	StockName          string
	CurrentPrice       float64
	DayLow             float64
	DayHigh            float64
	Volume             int
	PriceChange        float64
	PriceChangePercent float64
}

// PriceSVC ...
type PriceSVC struct {
	TimeService         *TimeSVC
	StockCache          *cache.Cache
	OvernightStockCache *cache.Cache
}

func InitPriceSVC(timeService *TimeSVC) (*PriceSVC, error) {
	overnightStockCache := cache.New(cache.NoExpiration, cache.NoExpiration)
	stockCache := cache.New(stockInfoCacheTimeout, stockInfoCacheTimeout)

	priceSvc := &PriceSVC{
		TimeService:         timeService,
		OvernightStockCache: overnightStockCache,
		StockCache:          stockCache,
	}

	return priceSvc, nil
}

// Close ...
func (svc *PriceSVC) Close() error {
	svc.OvernightStockCache.Flush()
	svc.StockCache.Flush()
	return nil
}

// GetStockInfo ...
func (svc *PriceSVC) GetStockInfo(symbol string) (*StockInfo, error) {
	logger := log.WithFields(log.Fields{
		"package":  "PriceSVC",
		"function": "GetStockInfo",
	})

	if svc.TimeService.GetMarketType(time.Now()) == Overnight {
		// clear stock cache for daytime
		svc.StockCache.Flush()

		if cache, ok := svc.OvernightStockCache.Get(symbol); ok {
			return cache.(*StockInfo), nil
		} else {
			stockInfo, err := svc.getStockInfo(symbol)
			if err != nil {
				logger.Error(err)
				return nil, err
			}

			// no expiration
			svc.OvernightStockCache.Set(symbol, stockInfo, 12*time.Hour)
			return stockInfo, nil
		}
	} else {
		svc.OvernightStockCache.Flush()

		if cache, ok := svc.StockCache.Get(symbol); ok {
			return cache.(*StockInfo), nil
		} else {
			stockInfo, err := svc.getStockInfo(symbol)
			if err != nil {
				logger.Error(err)
				return nil, err
			}

			svc.StockCache.SetDefault(symbol, stockInfo)
			return stockInfo, nil
		}
	}
}

func (svc *PriceSVC) getStockInfo(symbol string) (*StockInfo, error) {
	logger := log.WithFields(log.Fields{
		"package":  "PriceSVC",
		"function": "getStockInfo",
	})

	quote, err := yahoo.Quote(symbol)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	stockInfo := &StockInfo{
		Symbol:             symbol,
		StockName:          symbol,
		CurrentPrice:       0,
		DayLow:             0,
		DayHigh:            0,
		Volume:             0,
		PriceChange:        0,
		PriceChangePercent: 0,
	}

	if len(quote.QuoteSummary.Result) > 0 {
		stockInfo.CurrentPrice = quote.QuoteSummary.Result[0].Price.RegularMarketPrice.Raw
		if len(quote.QuoteSummary.Result[0].Price.LongName) > 0 {
			stockInfo.StockName = quote.QuoteSummary.Result[0].Price.LongName
		} else if len(quote.QuoteSummary.Result[0].Price.ShortName) > 0 {
			stockInfo.StockName = quote.QuoteSummary.Result[0].Price.ShortName
		} else {
			stockInfo.StockName = symbol
		}

		stockInfo.DayLow = quote.QuoteSummary.Result[0].Price.RegularMarketDayLow.Raw
		stockInfo.DayHigh = quote.QuoteSummary.Result[0].Price.RegularMarketDayHigh.Raw
		stockInfo.Volume = quote.QuoteSummary.Result[0].Price.RegularMarketVolume.Raw
		stockInfo.PriceChange = quote.QuoteSummary.Result[0].Price.RegularMarketChange.Raw
		stockInfo.PriceChangePercent = quote.QuoteSummary.Result[0].Price.RegularMarketChangePercent.Raw
	}
	return stockInfo, nil
}
