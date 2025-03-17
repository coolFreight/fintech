package equities

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"log/slog"
	"time"
)

type EquityService struct {
	client *marketdata.Client
	logger *slog.Logger
}

func NewEquityService(client *marketdata.Client, logger *slog.Logger) *EquityService {
	return &EquityService{client: client, logger: logger}
}

// Calculates the moving average for a given ticker starting with the current date
func (s *EquityService) MovingAverage(ticker string) (float64, error) {
	bars, err := s.client.GetBars(ticker, marketdata.GetBarsRequest{
		TimeFrame: marketdata.TimeFrame{1, marketdata.Day},
		Start:     time.Now().Add(-time.Duration(300) * 24 * time.Hour),
	})

	totalSum := 0.0
	sum := 0.0
	var volSum uint64
	window := [...]int{5, 10, 30, 60, 120, 200}
	days := 1
	idx := 0
	stop := len(bars) - 200
	for i := len(bars) - 1; i >= stop; i-- {
		bar := bars[i]
		sum += bar.Close
		volSum += bar.Volume
		if days == window[idx] {
			totalSum += sum
			avg := totalSum / float64(window[idx])
			s.logger.Info("Calculating moving average: ", "ticker", ticker, "avg", avg, "window", window[idx])
			sum = 0
			idx++
		}
		days++
	}

	//s.logger.Info("Calculating moving average: ", "ticker", ticker, "avg", avg, "volumeAggregate", volSum)

	if err != nil {
		s.logger.Error(err.Error())
		return 0, err
	}
	return 0, nil

}
