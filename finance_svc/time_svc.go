package finance_svc

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// MarketType ...
type MarketType string

const (
	PreMarket   MarketType = "Pre-market hours"
	DayMarket   MarketType = "Market hours"
	AfterMarket MarketType = "After-market hours"
	Overnight   MarketType = "Overnight hours"

	timeLayout     = "15:04:05"
	dateTimeLayout = "2006-01-02 15:04:05"

	PreMarketStartTime string = "07:00:00"
	MarketStartTime    string = "09:30:00"
	MarketEndTime      string = "16:00:00"
	AfterMarketEndTime string = "17:00:00"
)

// TimeSVC ...
type TimeSVC struct {
	NewYorkLocation *time.Location
	PhoenixLocation *time.Location

	PreMarketStartTime   time.Time
	MarketStartTime      time.Time
	AfterMarketStartTime time.Time
	AfterMarketEndTime   time.Time
}

func InitTimeSVC() (*TimeSVC, error) {
	logger := log.WithFields(log.Fields{
		"package":  "TimeSVC",
		"function": "InitTimeSVC",
	})

	newyorkLoc, err := time.LoadLocation("America/New_York")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	phoenixLoc, err := time.LoadLocation("America/Phoenix")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	premarketStartTime, err := time.ParseInLocation(timeLayout, PreMarketStartTime, newyorkLoc)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	marketStartTime, err := time.ParseInLocation(timeLayout, MarketStartTime, newyorkLoc)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	aftermarketStartTime, err := time.ParseInLocation(timeLayout, MarketEndTime, newyorkLoc)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	aftermarketEndTime, err := time.ParseInLocation(timeLayout, AfterMarketEndTime, newyorkLoc)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	timeSvc := &TimeSVC{
		NewYorkLocation: newyorkLoc,
		PhoenixLocation: phoenixLoc,

		PreMarketStartTime:   premarketStartTime,
		MarketStartTime:      marketStartTime,
		AfterMarketStartTime: aftermarketStartTime,
		AfterMarketEndTime:   aftermarketEndTime,
	}

	return timeSvc, nil
}

// Close ...
func (svc *TimeSVC) Close() error {
	return nil
}

// PhoenixToNewyork converts phoenix time to new york time
func (svc *TimeSVC) ToNewyork(t time.Time) time.Time {
	return t.In(svc.NewYorkLocation)
}

// NewyorkToPhoenix converts new york time to phoenix time
func (svc *TimeSVC) ToPhoenix(t time.Time) time.Time {
	return t.In(svc.PhoenixLocation)
}

// GetMarketType ...
func (svc *TimeSVC) GetMarketType(t time.Time) MarketType {
	logger := log.WithFields(log.Fields{
		"package":  "TimeSVC",
		"function": "GetMarketType",
	})

	nyTime := svc.ToNewyork(t)

	h, m, s := nyTime.Clock()
	nyTimeString := fmt.Sprintf("%02d:%02d:%02d", h, m, s)

	inputTime, err := time.ParseInLocation(timeLayout, nyTimeString, svc.NewYorkLocation)
	if err != nil {
		logger.Error(err)
		return Overnight
	}

	if inputTime.After(svc.AfterMarketEndTime) || inputTime.Before(svc.PreMarketStartTime) {
		return Overnight
	} else if inputTime.After(svc.PreMarketStartTime) && inputTime.Before(svc.MarketStartTime) {
		return PreMarket
	} else if inputTime.After(svc.MarketStartTime) && inputTime.Before(svc.AfterMarketStartTime) {
		return DayMarket
	} else if inputTime.After(svc.AfterMarketStartTime) && inputTime.Before(svc.AfterMarketEndTime) {
		return AfterMarket
	} else {
		return Overnight
	}
}
