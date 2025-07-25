package finam

import (
	"fmt"
	"time"

	marketdata_service "github.com/alexxnosk/finproto/backend/trade_api/v1/marketdata/marketdata_service"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Информация о котировке
type Quote struct {
	Symbol    string // Символ инструмента
	Timestamp int64  // Метка времени
	Ask       float64
	Bid       float64
	Last      float64 // Цена последней сделки
}

func (q *Quote) Reset() {
	*q = Quote{} // копируем в q значение "пустого" объекта
}

func (q *Quote) Time() time.Time {
	return time.Unix(0, q.Timestamp).In(TzMoscow)
}

type QuoteStore struct {
	quoteState map[string]*Quote // Последнее состояние по символу
	//mu         sync.Mutex        // защита для конкурентного доступа
}

// processQuote обработаем сырые данные. Вернем срез Quote
func (qs *QuoteStore) processQuote(rq *marketdata_service.Quote) (Quote, error) {
	if rq == nil || rq.Symbol == "" {
		return Quote{}, fmt.Errorf("некорректная котировка: отсутствует символ")
	}
	// пока уберу. все делается в одном потоке (в listenQuoteStream)
	//qs.mu.Lock()
	//defer qs.mu.Unlock()

	// Получаем текущее состояние, если есть
	q, ok := qs.quoteState[rq.Symbol]
	if !ok {
		q = &Quote{
			Symbol: rq.Symbol,
		}
		qs.quoteState[rq.Symbol] = q
	}

	// Обновляем только непустые поля
	if rq.Timestamp != nil {
		q.Timestamp = rq.Timestamp.AsTime().UnixNano()
	}
	if rq.Ask != nil {
		q.Ask, _ = DecimalToFloat64E(rq.Ask)
	}
	if rq.Bid != nil {
		q.Bid, _ = DecimalToFloat64E(rq.Bid)
	}
	if rq.Last != nil {
		q.Last, _ = DecimalToFloat64E(rq.Last)
	}

	// Возвращаем копию
	return Quote{
		Symbol:    q.Symbol,
		Timestamp: q.Timestamp,
		Ask:       q.Ask,
		Bid:       q.Bid,
		Last:      q.Last,
	}, nil
}

// NewBarsRequest
func NewBarsRequest(symbol string, tf marketdata_service.TimeFrame, start, end time.Time) *marketdata_service.BarsRequest {
	i := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	
	return &marketdata_service.BarsRequest{
		Symbol:    symbol,
		Timeframe: tf,
		Interval:  i,
	}
}

func NewBarsRequestInterval(symbol string, tf marketdata_service.TimeFrame, start time.Time, days int32) *marketdata_service.BarsRequest{
	startTime := timestamppb.New(start)
	endTime := timestamppb.New(time.Now())
	end := startTime.AsTime().Add(time.Duration(days) * 24 * time.Hour)
	if end.Before(time.Now()){
		endTime = timestamppb.New(end)
	}
	if endTime.AsTime().Before(startTime.AsTime()){
		startTime = endTime
	}
	i := &interval.Interval{
		StartTime: startTime,
		EndTime:   endTime,
	}
	return &marketdata_service.BarsRequest{
		Symbol:    symbol,
		Timeframe: tf,
		Interval:  i,
	}
}
// NewQuoteRequest
func NewQuoteRequest(symbol string) *marketdata_service.QuoteRequest {
	return &marketdata_service.QuoteRequest{
		Symbol: symbol,
	}
}

// NewOrderBookRequest
func NewOrderBookRequest(symbol string) *marketdata_service.OrderBookRequest {
	return &marketdata_service.OrderBookRequest{
		Symbol: symbol,
	}
}

// NewLatestTradesRequest
func NewLatestTradesRequest(symbol string) *marketdata_service.LatestTradesRequest {
	return &marketdata_service.LatestTradesRequest{
		Symbol: symbol,
	}
}

// NewSubscribeQuoteRequest
func NewSubscribeQuoteRequest(symbols []string) *marketdata_service.SubscribeQuoteRequest {
	return &marketdata_service.SubscribeQuoteRequest{Symbols: symbols}
}

// NewSubscribeOrderBookRequest
func NewSubscribeOrderBookRequest(symbol string) *marketdata_service.SubscribeOrderBookRequest {
	return &marketdata_service.SubscribeOrderBookRequest{Symbol: symbol}
}

// NewSubscribeLatestTradesRequest
func NewSubscribeLatestTradesRequest(symbol string) *marketdata_service.SubscribeLatestTradesRequest {
	return &marketdata_service.SubscribeLatestTradesRequest{Symbol: symbol}
}
