package data

import (
	"fmt"
	"math"
	"strconv"

	"github.com/alexxnosk/finproto/backend/trade_api/v1/marketdata/marketdata_service"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/type/decimal"
)

func StrToTimeFrame(tfStr string) (marketdata_service.TimeFrame, error) {
	switch tfStr {
	case "M1":
		return marketdata_service.TimeFrame_TIME_FRAME_M1, nil
	case "M5":
		return marketdata_service.TimeFrame_TIME_FRAME_M5, nil
	case "M15":
		return marketdata_service.TimeFrame_TIME_FRAME_M15, nil
	case "M30":
		return marketdata_service.TimeFrame_TIME_FRAME_M30, nil
	case "H1":
		return marketdata_service.TimeFrame_TIME_FRAME_H1, nil
	case "H2":
		return marketdata_service.TimeFrame_TIME_FRAME_H2, nil
	case "H4":
		return marketdata_service.TimeFrame_TIME_FRAME_H4, nil
	case "H8":
		return marketdata_service.TimeFrame_TIME_FRAME_H8, nil
	case "D":
		return marketdata_service.TimeFrame_TIME_FRAME_D, nil
	case "W":
		return marketdata_service.TimeFrame_TIME_FRAME_W, nil
	case "MN":
		return marketdata_service.TimeFrame_TIME_FRAME_MN, nil
	case "QR":
		return marketdata_service.TimeFrame_TIME_FRAME_QR, nil
	default:
		return marketdata_service.TimeFrame_TIME_FRAME_UNSPECIFIED, fmt.Errorf("unknown timeframe: %s", tfStr)
	}
}

func ConvertBarProtoToBar(proto *BarDecimal) (Bar, error) {
	var bar Bar

	if proto == nil {
		return Bar{}, fmt.Errorf("proto is nil")
	}

	if proto.Timestamp != nil {
		bar.Timestamp = proto.Timestamp.AsTime()
	}

	if proto.Open != nil {
		s := proto.Open.GetValue()
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return Bar{}, err
		}
		bar.Open = f
	}

	if proto.High != nil {
		s := proto.High.GetValue()
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return Bar{}, err
		}
		bar.High = f
	}

	if proto.Low != nil {
		s := proto.Low.GetValue()
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return Bar{}, err
		}
		bar.Low = f
	}

	if proto.Close != nil {
		s := proto.Close.GetValue()
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return Bar{}, err
		}
		bar.Close = f
	}

	if proto.Volume != nil {
		s := proto.Volume.GetValue()
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return Bar{}, err
		}
		bar.Volume = int64(math.Round(f))
	}

	return bar, nil
}

func ConvertBarToBarPG(bar Bar) BarPg {
	return BarPg{
		Timestamp: pgtype.Timestamptz{
			Time:  bar.Timestamp,
			Valid: true,
		},
		Open: pgtype.Float8{
			Float64: bar.Open,
			Valid:   true,
		},
		High: pgtype.Float8{
			Float64: bar.High,
			Valid:   true,
		},
		Low: pgtype.Float8{
			Float64: bar.Low,
			Valid:   true,
		},
		Close: pgtype.Float8{
			Float64: bar.Close,
			Valid:   true,
		},
		Volume: pgtype.Int8{
			Int64: bar.Volume,
			Valid: true,
		},
	}
}

func ConvertBarProtoToBarPG(proto *BarDecimal) (BarPg, error) {
	b, err := ConvertBarProtoToBar(proto)
	return ConvertBarToBarPG(b), err
}

func ConvertBarDecimalToBarPG(proto *BarDecimal) (BarPgUpdate, error) {
	barPG, err := ConvertBarProtoToBarPG(proto)
	BarSQLC := BarPgUpdate{
		Timestamp: barPG.Timestamp,
		Open:      barPG.Open,
		High:      barPG.High,
		Low:       barPG.Low,
		Close:     barPG.Close,
		Volume:    barPG.Volume,
	}
	return BarSQLC, err
}

func DecimalToFloat64E(d *decimal.Decimal) (float64, error) {
	if d == nil {
		return 0, nil
	}
	s := d.String()
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse decimal string %q: %w", s, err)
	}
	return f, nil
}

func DecimalToFloat64(d *decimal.Decimal) float64 {
	result, _ := DecimalToFloat64E(d)
	return result
}

func DecimalToInt64E(d *decimal.Decimal) (int64, error) {
	if d == nil {
		return 0, nil
	}
	val, err := DecimalToFloat64E(d)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

func DecimalToInt(d *decimal.Decimal) int64 {
	result, _ := DecimalToInt64E(d)
	return result
}
