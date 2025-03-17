package pricing

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log/slog"
)

type Quote struct {
	MessageType  string   `json:"T"`
	Symbol       string   `json:"S"`
	TradeId      int      `json:"i"`
	ExchangeCode string   `json:"x"`
	Price        float64  `json:"p"`
	Size         int      `json:"s"`
	Condition    []string `json:"c"`
	Timestamp    string   `json:"t"`
	Tape         string   `json:"z"`
}

type Quotes struct {
	Q []Quote
}

func NewPricing(conn *websocket.Conn, logger *slog.Logger) <-chan []Quote {
	quoteChan := make(chan []Quote)
	go func() {
		for {
			var quotes []Quote
			if err := websocket.JSON.Receive(conn, &quotes); err != nil {
				logger.Error("Could not get pricing ", err)
				close(quoteChan)
				return
			}
			quoteChan <- quotes
		}
	}()
	return quoteChan
}

func (q Quote) String() string {
	return fmt.Sprintf("Quote"+
		"{Ticker: %s, "+
		"Price: %f, Size: %d,Time: %s}", q.Symbol, q.Price, q.Size, q.Timestamp)
}
