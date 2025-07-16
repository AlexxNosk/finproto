package main1

import (
	"fmt"
	//"time"

	"github.com/alexxnosk/finproto/backend/trade_api/data"
	//"github.com/jackc/pgx/v5"
)

func main() {

var jsonData = []byte(`{
  "request": {
    "symbol": "VTBR@MISX",
    "timeframe": "D",
    "start_date": "01-01-2003",
    "end_date": "nil",
    "operation": "delete"
  },
  "tables": {
    "fundamentals": {
      "eps": "float64",
      "dividends": "float64",
      "currency": "string"
    },
    "news": {
      "headline": "string",
      "date": "string",
      "impact": "int64"
    }
  }
}`)

	err := data.InstrumentCreate(jsonData)
	if err != nil {
		fmt.Println("ERROR:", err)
	}

}
