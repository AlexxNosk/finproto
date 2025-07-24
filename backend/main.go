package main

import (
	"fmt"

	"github.com/alexxnosk/finproto/backend/config"
	"github.com/alexxnosk/finproto/backend/trade_api/data"
)

func main() {

	var jsonData = []byte(`{
  "request": {
    "symbol": "SBER@MISX",
    "timeframe": "D",
    "start_date": "01-01-2000",
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
	token := config.LoadConfig().TOKEN
	err := data.AssetCreate(jsonData, token)
	//_, err := data.AssetsUpload(token)
	if err != nil {
		fmt.Println("ERROR:", err)
	}

}
