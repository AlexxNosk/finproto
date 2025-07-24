package data

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/alexxnosk/finproto/backend/trade_api/v1/assets/assets_service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) GetAssets(ctx context.Context) ([]Asset, error) {
	ctxWithToken, err := c.WithAuthToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth token: %w", err)
	}

	resp, err := c.AssetsService.Assets(ctxWithToken, &assets_service.AssetsRequest{})
	if err != nil {
		return nil, fmt.Errorf("AssetsService.Assets: %w", err)
	}

	var result []Asset
	for _, a := range resp.Assets {
		result = append(result, Asset{
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


func AssetsTable(operation string, token string) ([]Asset, int, error) {
    ctx := context.Background()
	client, err := NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return nil, 0, err
	}
	defer client.Close(ctx)

	tx, err := client.connPG.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	assetsFromGRPC, err := client.GetAssets(ctx)
	if err != nil {
		slog.Error("Client.GetAssets", "err", err.Error())
		return nil, 0, err
	}

	switch strings.ToLower(operation) {
	case "create":
		// Clean and insert all assets
		_, err := tx.Exec(ctx, "TRUNCATE assets RESTART IDENTITY CASCADE")
		if err != nil {
			return nil, 0, fmt.Errorf("truncate assets: %w", err)
		}

		count, err := insertAssets(ctx, tx, assetsFromGRPC)
		if err != nil {
			return nil, 0, err
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, 0, err
		}
		return nil, count, nil

	case "read":
		// Simple read logic
		rows, err := client.connPG.Query(ctx, "SELECT id, ticker, symbol, name, mic, type, external_id FROM assets")
		if err != nil {
			return nil, 0, fmt.Errorf("read assets: %w", err)
		}
		defer rows.Close()

		var assets []Asset
		for rows.Next() {
			var s Asset
			err = rows.Scan(&s.Id, &s.Ticker, &s.Symbol, &s.Name, &s.Mic, &s.Type, &s.Id)
			if err != nil {
				return nil, 0, err
			}
			assets = append(assets, s)
		}
		//PrintSecurities(assets)
		return assets, len(assets), nil

	case "update":
		// Fetch current external_ids
		rows, err := client.connPG.Query(ctx, "SELECT external_id FROM assets")
		if err != nil {
			return nil, 0, fmt.Errorf("fetch existing ids: %w", err)
		}
		defer rows.Close()

		existingIDs := make(map[string]struct{})
		for rows.Next() {
			var id string
			rows.Scan(&id)
			existingIDs[id] = struct{}{}
		}

		// Filter new assets
		var newAssets []Asset
		for _, s := range assetsFromGRPC {
			if _, exists := existingIDs[s.Id]; !exists {
				newAssets = append(newAssets, s)
			}
		}

		countNewAssets, err := insertAssets(ctx, tx, newAssets)
		if err != nil {
			return nil, 0, err
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, 0, err
		}
		return newAssets, countNewAssets, nil

	case "delete":
		_, err := tx.Exec(ctx, "DROP TABLE IF EXISTS assets")
		if err != nil {
			return nil, 0, fmt.Errorf("drop table: %w", err)
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, 0, err
		}
		return nil, 1, nil

	default:
		return nil, 0, fmt.Errorf("unsupported operation: %s", operation)
	}
}

func insertAssets(ctx context.Context, tx pgx.Tx, securities []Asset) (int, error) {
	insertQuery := `
        INSERT INTO assets (ticker, symbol, name, mic, type, external_id)
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
	slog.Info("Inserted assets", "count", count)
	return count, nil
}

// func PrintSecurities(sec []Asset) {
// 	for r, b := range sec {
// 		fmt.Println(r,
// 			"Ticker", b.Ticker,
// 			"Symbol", b.Symbol,
// 			"Name", b.Name,
// 			"Mic", b.Mic,
// 			"Type", b.Type,
// 			"Id", b.Id,
// 		)
// 	}
// }

// func ToSecurity(assets []*assets_service.Asset) []Asset {
// 	result := make([]Asset, len(assets))
// 	for i, asset := range assets {
// 		result[i] = Asset{
// 			Ticker: asset.Ticker,
// 			Symbol: asset.Symbol,
// 			Name:   asset.Name,
// 			Mic:    asset.Mic,
// 			Type:   asset.Type,
// 			Id:     asset.Id,
// 		}
// 	}
// 	return result
// }
