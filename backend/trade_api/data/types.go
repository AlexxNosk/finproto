package data

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/type/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Bar struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int64
}

type BarDecimal struct {
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Open      *decimal.Decimal       `protobuf:"bytes,2,opt,name=open,proto3" json:"open,omitempty"`
	High      *decimal.Decimal       `protobuf:"bytes,3,opt,name=high,proto3" json:"high,omitempty"`
	Low       *decimal.Decimal       `protobuf:"bytes,4,opt,name=low,proto3" json:"low,omitempty"`
	Close     *decimal.Decimal       `protobuf:"bytes,5,opt,name=close,proto3" json:"close,omitempty"`
	Volume    *decimal.Decimal       `protobuf:"bytes,6,opt,name=volume,proto3" json:"volume,omitempty"`
}

type BarPg struct {
	ID        int32
	InstrumentID int32
	TimeframeID	int32
	Timestamp pgtype.Timestamptz
	Open      pgtype.Float8
	High      pgtype.Float8
	Low       pgtype.Float8
	Close     pgtype.Float8
	Volume    pgtype.Int8
}

type BarPgUpdate struct {
	Timestamp pgtype.Timestamptz
	Open      pgtype.Float8
	High      pgtype.Float8
	Low       pgtype.Float8
	Close     pgtype.Float8
	Volume    pgtype.Int8
}

type Instrument struct {
	ID        int32
	Symbol    string
	Name      pgtype.Text
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type InstrumentTimeframe struct {
	InstrumentID int32
	TimeframeID  int32
}

type Timeframe struct {
	ID          int32
	Code        string
	Description pgtype.Text
}
