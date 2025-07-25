package data

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/alexxnosk/finproto/backend/trade_api/v1/marketdata/marketdata_service"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BarsFromFinam(ctx context.Context, client *Client, symbol, tfStr, startStr, endStr string) ([]*BarDecimal, *string, error) {
	var barDecimal []*BarDecimal

	if symbol == "" || tfStr == "" {
		return nil, nil, fmt.Errorf("symbol and timeframe are required")
	}

	ctx, err := client.WithAuthToken(ctx)
	if err != nil {
		slog.Error("BarsFromServer", "WithAuthToken", err.Error())
		return nil, nil, err
	}

	reqs, err := MakeBarRequests(symbol, tfStr, startStr, endStr)
	if err != nil {
		slog.Error("BarsFromServer", "MakeBarRequests", err.Error())
		return nil, nil, err
	}

	// for i, j := 0, len(reqs)-1; i < j; i, j = i+1, j-1 {
	// reqs[i], reqs[j] = reqs[j], reqs[i]
	// }

	var resolvedSymbol string
	for _, req := range reqs {
		bars, err := client.MarketDataService.Bars(ctx, req)
		if err != nil {
			slog.Error("BarsFromServer", "MarketDataService.Bars", err.Error())
			return nil, nil, err
		}
		barDecimalChank, symb := BarsResponseDecompose(bars)
		resolvedSymbol = symb
		barDecimal = append(barDecimal, barDecimalChank...)
	}
	return barDecimal, &resolvedSymbol, nil
}

func MakeBarRequests(symbol, tfStr, startStr, endStr string) ([]*marketdata_service.BarsRequest, error) {
	var requests []*marketdata_service.BarsRequest
	var startTime = time.Now()
	var endTime = time.Now()
	var err error
	tf, err := StrToTimeFrame(tfStr)
	if err != nil {
		slog.Error("unknown timeframe", "err", err)
		return nil, err
	}

	layout := "02-01-2006"

	if startStr != "nil" {
		startTime, err = time.Parse(layout, startStr)
		if err != nil {
			slog.Error("invalid start time format", "err", err)
			return nil, err
		}
	}

	if endStr != "nil" {
		endTime, err = time.Parse(layout, endStr)
		if err != nil {
			slog.Error("invalid end time format", "err", err)
			return nil, err
		}
	}

	// Define how many days each request chunk should cover
	var daysInInterval int
	switch tf {
	case 1:
		daysInInterval = 7 // 1-minute bars
	case 5, 9, 11, 12, 13, 15, 17:
		daysInInterval = 30 // 8-hour bars
	case 19:
		daysInInterval = 365 // daily bars
	case 20, 21, 22:
		daysInInterval = 365 * 5 // quarterly bars
	default:
		slog.Error("unsupported timeframe", "tf", tf)
		return nil, fmt.Errorf("unsupported timeframe: %v", tf)
	}

	// Create batched requests by time interval
	for current := startTime; current.Before(endTime); current = current.AddDate(0, 0, daysInInterval) {
		next := current.AddDate(0, 0, daysInInterval)
		if next.After(endTime) {
			next = endTime
		}

		requests = append(requests, &marketdata_service.BarsRequest{
			Symbol:    symbol,
			Timeframe: tf,
			Interval: &interval.Interval{
				StartTime: timestamppb.New(current),
				EndTime:   timestamppb.New(next),
			},
		})
	}
	// requests = append(requests, &marketdata_service.BarsRequest{
	// 	Symbol: symbol,
	// 	Timeframe: tf,
	// 	Interval: &interval.Interval{
	// 		StartTime: timestamppb.New(startTime),
	// 		EndTime:   timestamppb.New(endTime),
	// 	},
	// })
	return requests, nil
}

func BarsResponseDecompose(resp *marketdata_service.BarsResponse) ([]*BarDecimal, string) {
	var barDecimal []*BarDecimal
	for _, bar := range resp.Bars {
		barDecimal = append(barDecimal, &BarDecimal{
			Timestamp: bar.Timestamp,
			Open:      bar.Open,
			High:      bar.High,
			Low:       bar.Low,
			Close:     bar.Close,
			Volume:    bar.Volume,
		})
	}
	return barDecimal, resp.Symbol
}

func PrintBarsDecimal(barDecimal []*BarDecimal) {
	for r, b := range barDecimal {
		slog.Info("Bar", "row", r,
			"Timestamp", b.Timestamp.AsTime(),
			"Open", b.Open.Value,
			"High", b.High.Value,
			"Low", b.Low.Value,
			"Close", b.Close.Value,
			"Volume_int", b.Volume.Value,
		)
	}
}
