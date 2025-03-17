package main

import (
	"encoding/json"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/coolFreight/fintech/equities"
	"github.com/coolFreight/fintech/internal/client/pricing"
	"golang.org/x/net/websocket"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"
)

type TickerAggregate struct {
	Ticker     string
	Low        float64
	High       float64
	OpenPrice  float64
	ClosePrice float64
	Volume     float64 //volume for a given time period
	Timestamp  time.Time
}

const (
	TIMESPAN_SEC  = "second"
	TIMESPAN_MIN  = "minute"
	TIMESPAN_HOUR = "hour"
	TIMESPAN_DAY  = "day"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("calling alpaca %s", os.Getenv("APCA_PAPER_ENDPOINT"))
	client := alpaca.NewClient(alpaca.ClientOpts{
		// Alternatively you can set your Key and Secret using the
		// APCA_API_KEY_ID and APCA_API_SECRET_KEY environment variables
		APIKey:    os.Getenv("APCA_PAPER_API_KEY"),
		APISecret: os.Getenv("APCA_PAPER_API_SECRET"),
		BaseURL:   os.Getenv("APCA_PAPER_ENDPOINT"),
	})

	acct, err := client.GetAccount()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *acct)

	marketClient := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    os.Getenv("APCA_PAPER_API_KEY"),
		APISecret: os.Getenv("APCA_PAPER_API_SECRET"),
	})

	service := equities.NewEquityService(marketClient, logger)
	movingAverage, err := service.MovingAverage("YOU")

	if err != nil {
		fmt.Printf("%+v\n", err)
	} else {
		fmt.Printf("Moving average is %f\n", movingAverage)
	}

	type Connect struct {
		Action string   `json:"action"`
		Trades []string `json:"trades"`
	}

	type auth struct {
		Action string `json:"action"`
		Key    string `json:"key"`
		Secret string `json:"secret"`
	}

	origin := "https://paper-api.alpaca.markets/v2"
	url := "wss://stream.data.alpaca.markets/v2/test"
	ws, err := websocket.Dial(url, "", origin)

	if err != nil {
		log.Fatal(err)
	}
	read(ws, logger)

	fmt.Printf("Attempting to authenticate\n")

	authenticate := auth{
		Action: "auth",
		Key:    os.Getenv("APCA_PAPER_API_KEY"),
		Secret: os.Getenv("APCA_PAPER_API_SECRET"),
	}
	err = send(authenticate, ws, logger)
	if err != nil {
		log.Fatal("could not authenticate - ", err)
	}
	read(ws, logger)

	fmt.Printf("Requesting pricing\n")
	connect := Connect{Action: "subscribe", Trades: []string{"FAKEPACA"}}
	err = send(connect, ws, logger)
	read(ws, logger)
	read(ws, logger)

	var wg sync.WaitGroup
	prices := pricing.NewPricing(ws, logger)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for quote := range prices {
			fmt.Println(quote)
		}
	}()
	wg.Wait()

	//autho := auth{Action: "auth", Key: os.Getenv("APCA_PAPER_API_KEY"), Secret: os.Getenv("APCA_PAPER_API_SECRET")}
	//err = websocket.Message.Send(ws, autho)
	//if err != nil {
	//	log.Fatal("could not auth ", err)
	//}

	//if err != nil {
	//	fmt.Printf("%+v\n", err)
	//} else {
	//	fmt.Printf("Moving average is %f", movingAverage)
	//}
	//
	//q, err := marketClient.GetLatestQuote("SPXS", marketdata.GetLatestQuoteRequest{Currency: "USD", Feed: marketdata.IEX})

	//if err != nil {
	//	fmt.Printf("%+v\n", err)
	//}
	//fmt.Printf("%+v\n", q)

	//alpaca.SetBaseUrl("https://paper-api.alpaca.markets")

	// Get account information.
	//account, err := alpaca.GetAccount()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//// Calculate the difference between current balance and balance at the last market close.
	//balanceChange := account.Equity.Sub(account.LastEquity)

	//fmt.Println("Today's portfolio balance change:", balanceChange)

	//logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//
	//ctx := context.Background()
	//c := polygon.New(os.Getenv("POLYGON_API_KEY"))
	//
	//params := models.ListAggsParams{
	//	Ticker:     "ACHR",
	//	Multiplier: 1,
	//	Timespan:   TIMESPAN_DAY,
	//	From:       models.Millis(time.Date(2025, 3, 12, 14, 30, 0, 0, time.UTC)),
	//	To:         models.Millis(time.Date(2025, 3, 12, 21, 00, 0, 0, time.UTC)),
	//}
	//
	//resp := c.ListAggs(ctx, &params)
	//
	//if resp.Err() != nil {
	//	logger.Error("error fetching data", resp.Err())
	//	return
	//}
	//aggregates := make([]TickerAggregate, 0)
	//count := 0
	//var result float64
	//var day time.Time
	//var prev TickerAggregate
	//for resp.Next() {
	//	item := resp.Item()
	//	ta := TickerAggregate{
	//		Ticker:     item.Ticker,
	//		Low:        item.Low,
	//		High:       item.High,
	//		ClosePrice: item.Close,
	//		Volume:     item.Volume,
	//		OpenPrice:  item.Open,
	//		Timestamp:  time.Time(item.Timestamp),
	//	}
	//
	//	aggregates = append(aggregates, ta)
	//	// do something with the result
	//	//logger.Info("found :")
	//	count++
	//	//fmt.Println(ta)
	//	result = (math.Abs((item.Open - item.High)) / item.Open) * 100
	//	closedResult := (math.Abs((prev.ClosePrice - item.Close)) / prev.ClosePrice) * 100
	//	day = ta.Timestamp
	//	fmt.Println(fmt.Sprintf(" %s on %s had the potential to make  %f but closed at %f", params.Ticker, day, result, closedResult))
	//	prev = ta
	//}

}

func (t TickerAggregate) String() string {
	return fmt.Sprintf("TickerAggregate"+
		"{Ticker: %s, "+
		"High: %f, Low: %f, OpenPrice: %f, ClosePrice: %f, Time: %s}", t.Ticker, t.High, t.Low, t.OpenPrice, t.ClosePrice, t.Timestamp)
}

func send(data any, ws *websocket.Conn, logger *slog.Logger) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logger.Error(fmt.Sprintf("could not marshal data %s", data, err))
	}
	return websocket.Message.Send(ws, dataBytes)
}

func read(ws *websocket.Conn, logger *slog.Logger) {
	var msg = make([]byte, 512)
	var n int
	var err error
	if n, err = ws.Read(msg); err != nil {
		logger.Error("Could not read data", err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}
