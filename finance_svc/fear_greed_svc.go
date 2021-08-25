package finance_svc

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

const (
	chartCacheTimeout = 1 * time.Hour // 1 hour
)

type FeerGreedIndexInfo struct {
	ImageURL string
}

// FearGreedIndexSVC ...
type FearGreedIndexSVC struct {
	ChartCache *cache.Cache
}

func InitFearGreedIndexSVC() (*FearGreedIndexSVC, error) {
	chartCache := cache.New(chartCacheTimeout, chartCacheTimeout)

	indexSvc := &FearGreedIndexSVC{
		ChartCache: chartCache,
	}

	return indexSvc, nil
}

// Close ...
func (svc *FearGreedIndexSVC) Close() error {
	svc.ChartCache.Flush()
	return nil
}

// GetIndexData ...
func (svc *FearGreedIndexSVC) GetIndexData(name string) ([]byte, error) {
	switch name {
	case "stock":
		return svc.GetStockIndexData()
	case "crypto":
		return svc.GetCryptoIndexData()
	default:
		return nil, fmt.Errorf("unknown index name")
	}
}

// GetStockIndexData ...
func (svc *FearGreedIndexSVC) GetStockIndexData() ([]byte, error) {
	logger := log.WithFields(log.Fields{
		"package":  "FearGreedIndexSVC",
		"function": "GetStockIndexData",
	})

	info, err := svc.GetStockIndexInfo()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	response, err := http.Get(info.ImageURL)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

// GetStockIndexInfo ...
func (svc *FearGreedIndexSVC) GetStockIndexInfo() (*FeerGreedIndexInfo, error) {
	logger := log.WithFields(log.Fields{
		"package":  "FearGreedIndexSVC",
		"function": "GetStockIndexInfo",
	})

	if cache, ok := svc.ChartCache.Get("stock"); ok {
		return cache.(*FeerGreedIndexInfo), nil
	} else {
		indexInfo, err := svc.getStockIndexInfo()
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		svc.ChartCache.SetDefault("stock", indexInfo)
		return indexInfo, nil
	}
}

// GetCryptoIndexData ...
func (svc *FearGreedIndexSVC) GetCryptoIndexData() ([]byte, error) {
	logger := log.WithFields(log.Fields{
		"package":  "FearGreedIndexSVC",
		"function": "GetCryptoIndexData",
	})

	info, err := svc.GetCryptoIndexInfo()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	response, err := http.Get(info.ImageURL)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

// GetCryptoIndexInfo ...
func (svc *FearGreedIndexSVC) GetCryptoIndexInfo() (*FeerGreedIndexInfo, error) {
	logger := log.WithFields(log.Fields{
		"package":  "FearGreedIndexSVC",
		"function": "GetCryptoIndexInfo",
	})

	if cache, ok := svc.ChartCache.Get("crypto"); ok {
		return cache.(*FeerGreedIndexInfo), nil
	} else {
		indexInfo, err := svc.getCryptoIndexInfo()
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		svc.ChartCache.SetDefault("crypto", indexInfo)
		return indexInfo, nil
	}
}

func (svc *FearGreedIndexSVC) getStockIndexInfo() (*FeerGreedIndexInfo, error) {
	logger := log.WithFields(log.Fields{
		"package":  "FearGreedIndexSVC",
		"function": "getStockIndexInfo",
	})

	doc, err := htmlquery.LoadURL("https://money.cnn.com/data/fear-and-greed/")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	nodes, err := htmlquery.QueryAll(doc, "//div[@id='needleChart']")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	imageURL := ""
	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "style" {
				if strings.HasPrefix(attr.Val, "background-image:url") {
					vals := strings.Split(attr.Val, "'")
					if len(vals) == 3 {
						imageURL = vals[1]
					}
				}
			}
		}
	}

	return &FeerGreedIndexInfo{
		ImageURL: imageURL,
	}, nil
}

func (svc *FearGreedIndexSVC) getCryptoIndexInfo() (*FeerGreedIndexInfo, error) {
	return &FeerGreedIndexInfo{
		ImageURL: "https://alternative.me/crypto/fear-and-greed-index.png",
	}, nil
}
