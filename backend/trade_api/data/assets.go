package data

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/alexxnosk/finproto/backend/config"
	"github.com/alexxnosk/finproto/backend/trade_api/v1/assets/assets_service"
    "github.com/jackc/pgx/v5"

)
func (c *Client) GetAssets(ctx context.Context) ([]Security, error) {
	ctxWithToken, err := c.WithAuthToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth token: %w", err)
	}

	resp, err := c.AssetsService.Assets(ctxWithToken, &assets_service.AssetsRequest{})
	if err != nil {
		return nil, fmt.Errorf("AssetsService.Assets: %w", err)
	}

	var result []Security
	for _, a := range resp.Assets {
		result = append(result, Security{
			Ticker: a.Ticker,
			Symbol: a.Symbol,
			Name:   a.Name,
			Mic:    a.Mic,
			Type:   a.Type,
			Id:     a.Id,
		})
	}

	return result, nil
}




func AssetsUpload() (int, error) {
    cnf := config.LoadConfig()
    token := cnf.TOKEN

    ctx := context.Background()
    client, err := NewClient(ctx, token)
    if err != nil {
        slog.Error("NewClient", "err", err.Error())
        return 0, err
    }
    defer client.Close()

    securities, err := client.GetAssets(ctx)
    if err != nil {
        slog.Error("Client.GetAssets", "err", err.Error())
        return 0, err
    }

    // Connect using pgx
    connStr := "postgres://root:root@localhost:5434/finProto_db"
    conn, err := pgx.Connect(ctx, connStr)
    if err != nil {
        slog.Error("pgx.Connect", "err", err.Error())
        return 0, err
    }
    defer conn.Close(ctx)

    // Start a transaction
    tx, err := conn.Begin(ctx)
    if err != nil {
        return 0, err
    }
    defer tx.Rollback(ctx) // safe to defer, commit overrides this

    insertQuery := `
        INSERT INTO instruments (ticker, symbol, name, mic, type, external_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (symbol) DO NOTHING
    `
    count := 0

    for _, s := range securities {
        _, err := tx.Exec(ctx, insertQuery, s.Ticker, s.Symbol, s.Name, s.Mic, s.Type, s.Id)
        if err != nil {
            slog.Warn("Insert failed", "symbol", s.Symbol, "err", err)
            continue
        }
        count++
    }

    if err := tx.Commit(ctx); err != nil {
        return 0, err
    }

    slog.Info("Instruments uploaded", "count", count)
    return count, nil
}



func PrintSecurities(sec []Security){
	for r, b := range sec {
		fmt.Println( r,
			"Ticker", b.Ticker,
			"Symbol", b.Symbol,
			"Name", b.Name,
			"Mic", b.Mic,
			"Type", b.Type,
			"Id", b.Id,
		)
	}
}

func ToSecurity(assets []*assets_service.Asset) []Security {
	result := make([]Security, len(assets))
	for i, asset := range assets {
		result[i] = Security{
			Ticker: asset.Ticker,
			Symbol: asset.Symbol,
			Name:   asset.Name,
			Mic:    asset.Mic,
			Type:   asset.Type,
			Id:     asset.Id,
		}
	}
	return result
}