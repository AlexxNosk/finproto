package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/alexxnosk/finproto/backend/trade_api/v1/assets/assets_service"
	"github.com/jackc/pgx/v5"
)

type AssetRequest struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Operation string `json:"operation"` // "create", "update", "delete"
}

// Formats dynamic table name, e.g., "SBER@MISX", "D" → "bars_sber_misx_d"
func assetTableName(symbol, tfStr, name string) string {
	s := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(symbol, "@", "_"), "-", "_"))
	t := strings.ToLower(tfStr)
	if n := strings.ToLower(name); n != "" {
		return fmt.Sprintf("bars_%s_%s_%s", s, t, n)
	}
	return fmt.Sprintf("bars_%s_%s", s, t)
}

func (c *Client) GetAssets(ctx context.Context) ([]AssetFinam, error) {
	ctxWithToken, err := c.WithAuthToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth token: %w", err)
	}

	resp, err := c.AssetsService.Assets(ctxWithToken, &assets_service.AssetsRequest{})
	if err != nil {
		return nil, fmt.Errorf("AssetsService.Assets: %w", err)
	}

	var result []AssetFinam
	for _, a := range resp.Assets {
		result = append(result, AssetFinam{
			Ticker: a.Ticker,
			Symbol: a.Symbol,
			Name:   a.Name,
			Mic:    a.Mic,
			Type:   a.Type,
			ID:     a.Id,
		})
	}

	return result, nil
}

func FinamAssetsTable(operation string, token string) ([]AssetFinamPG, int, error) {
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

	assetsFromFinam, err := client.GetAssets(ctx)
	seen := make(map[string]struct{})
	for _, s := range assetsFromFinam {
		if _, exists := seen[s.Symbol]; exists {
			slog.Warn("Duplicate symbol in input", "symbol", s.Symbol)
		}
		seen[s.Symbol] = struct{}{}
	}

	if err != nil {
		slog.Error("Client.GetAssets", "err", err.Error())
		return nil, 0, err
	}

	switch strings.ToLower(operation) {
	case "create":
		err := createFinamAssets(ctx, tx)
		if err != nil {
			return nil, 0, fmt.Errorf("create finam_asset table: %w", err)
		}

		assetsPG, count, err := insertAssets(ctx, tx, assetsFromFinam)
		if err != nil {
			return nil, 0, err
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, 0, err
		}
		return assetsPG, count, nil

	case "read":
		// Simple read logic
		rows, err := client.connPG.Query(ctx, "SELECT id, ticker, symbol, name, mic, type, finam_id FROM finam_assets")
		if err != nil {
			return nil, 0, fmt.Errorf("read assets: %w", err)
		}
		defer rows.Close()

		var assetsPG []AssetFinamPG
		for rows.Next() {
			var s AssetFinamPG
			err = rows.Scan(&s.ID, &s.Ticker, &s.Symbol, &s.Name, &s.Mic, &s.Type, &s.FinamID)
			if err != nil {
				return nil, 0, err
			}
			assetsPG = append(assetsPG, s)
		}
		//PrintSecurities(assets)
		return assetsPG, len(assetsPG), nil

	case "update":
		// Fetch current external_ids
		rows, err := client.connPG.Query(ctx, "SELECT finam_id FROM finam_assets")
		if err != nil {
			return nil, 0, fmt.Errorf("fetch existing ids: %w", err)
		}
		defer rows.Close()

		existingIDs := make(map[string]struct{})
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return nil, 0, fmt.Errorf("scan finam_id: %w", err)
			}
			existingIDs[id] = struct{}{}
		}

		// Filter new assets
		var newAssets []AssetFinam
		for _, s := range assetsFromFinam {
			if _, exists := existingIDs[s.ID]; !exists {
				newAssets = append(newAssets, s)
			}
		}

		assetsPG, countNewAssets, err := insertAssets(ctx, tx, newAssets)
		if err != nil {
			return nil, 0, err
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, 0, err
		}
		return assetsPG, countNewAssets, nil

	case "delete":
		res, err := tx.Exec(ctx, `
            DELETE FROM finam_assets
            WHERE id NOT IN (
                SELECT DISTINCT asset_id FROM asset_tables
            )
        `)
		if err != nil {
			return nil, 0, fmt.Errorf("delete unused assets: %w", err)
		}
		deletedCount := int(res.RowsAffected()) // <-- This gives the number of deleted rows

		if err := tx.Commit(ctx); err != nil {
			return nil, 0, err
		}
		slog.Info("Deleted unused assets", "count", deletedCount)
		return nil, deletedCount, nil

	case "delete_all":
		deletedCount, err := DeleteAllFinamAssetsAndDrop(ctx, tx)
		if err != nil {
			tx.Rollback(ctx)
			return nil, 0, err
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, 0, fmt.Errorf("commit failed: %w", err)
		}
		return nil, deletedCount, nil

	default:
		return nil, 0, fmt.Errorf("unsupported operation: %s", operation)
	}
}

func insertAssets(ctx context.Context, tx pgx.Tx, securities []AssetFinam) ([]AssetFinamPG, int, error) {
	insertQuery := `
        INSERT INTO finam_assets (ticker, symbol, name, mic, type, finam_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (symbol) DO NOTHING
        RETURNING id, ticker, symbol, name, mic, type, finam_id
    `

	var inserted []AssetFinamPG
	count := 0
	skipped := 0
	for _, s := range securities {
		var row AssetFinamPG
		err := tx.QueryRow(ctx, insertQuery, s.Ticker, s.Symbol, s.Name, s.Mic, s.Type, s.ID).
			Scan(&row.ID, &row.Ticker, &row.Symbol, &row.Name, &row.Mic, &row.Type, &row.FinamID)

		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				slog.Debug("Skipped insert due to conflict", "symbol", s.Symbol)
				skipped++
				continue
			}
			slog.Warn("Insert failed", "symbol", s.Symbol, "err", err)
			continue
		}

		inserted = append(inserted, row)
		count++
	}

	slog.Info("Insert summary", "inserted", count, "skipped", skipped)
	return inserted, count, nil
}

func createFinamAssets(ctx context.Context, tx pgx.Tx) error {
	exists, err := TableExists(ctx, tx, "finam_assets")
	if err != nil {
		return err
	}

	if exists {
		slog.Info("finam_assets exists: performing full cleanup and drop")

		if _, err := DeleteAllFinamAssetsAndDrop(ctx, tx); err != nil {
			return fmt.Errorf("delete_all failed: %w", err)
		}
	}

	slog.Info("Creating finam_assets and asset_tables")
	if err := CreateFinamAssetsAndAssetTables(ctx, tx); err != nil {
		return fmt.Errorf("create tables failed: %w", err)
	}

	return nil
}

func TableExists(ctx context.Context, tx pgx.Tx, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = $1
		)
	`
	var exists bool
	err := tx.QueryRow(ctx, query, tableName).Scan(&exists)
	return exists, err
}

func CreateFinamAssetsAndAssetTables(ctx context.Context, tx pgx.Tx) error {
	createAssets := `
	CREATE TABLE IF NOT EXISTS finam_assets (
		id           SERIAL PRIMARY KEY,
		ticker       VARCHAR NOT NULL,
		symbol       VARCHAR NOT NULL UNIQUE,
		name         VARCHAR,
		mic          VARCHAR NOT NULL,
		type         VARCHAR NOT NULL,
		finam_id     VARCHAR NOT NULL UNIQUE,
		updated_at   TIMESTAMPTZ DEFAULT now()
	)`
	if _, err := tx.Exec(ctx, createAssets); err != nil {
		return fmt.Errorf("create finam_assets: %w", err)
	}

	createAssetTables := `
	CREATE TABLE IF NOT EXISTS asset_tables (
		id             SERIAL PRIMARY KEY,
		asset_id       INT NOT NULL REFERENCES finam_assets(id) ON DELETE CASCADE,
		timeframe_id   INT NOT NULL,
		table_name     TEXT NOT NULL UNIQUE,
		created_at     TIMESTAMPTZ DEFAULT now(),
		updated_at     TIMESTAMPTZ DEFAULT now(),
		UNIQUE (asset_id, timeframe_id)
	)`
	if _, err := tx.Exec(ctx, createAssetTables); err != nil {
		return fmt.Errorf("create asset_tables: %w", err)
	}

	return nil
}

func DeleteAllFinamAssetsAndDrop(ctx context.Context, tx pgx.Tx) (int, error) {
	slog.Info("Starting full delete of finam_assets and related tables")

	// Step 1: Get asset IDs from finam_assets
	rows, err := tx.Query(ctx, `SELECT id FROM finam_assets`)
	if err != nil {
		return 0, fmt.Errorf("query asset ids: %w", err)
	}
	var assetIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("scan asset id: %w", err)
		}
		assetIDs = append(assetIDs, id)
	}
	rows.Close()

	deletedCount := len(assetIDs)

	// Step 2: Fetch table_names from asset_tables for these asset_ids
	var tableNames []string
	if deletedCount > 0 {
		rows, err := tx.Query(ctx, `SELECT table_name FROM asset_tables WHERE asset_id = ANY($1)`, assetIDs)
		if err != nil {
			return 0, fmt.Errorf("query asset_tables: %w", err)
		}
		for rows.Next() {
			var tbl string
			if err := rows.Scan(&tbl); err != nil {
				return 0, fmt.Errorf("scan table_name: %w", err)
			}
			tableNames = append(tableNames, tbl)
		}
		rows.Close()
	}

	// Step 3: Drop the actual physical tables listed in asset_tables
	seen := make(map[string]struct{})
	for _, tbl := range tableNames {
		if _, exists := seen[tbl]; exists {
			continue // skip duplicates
		}
		seen[tbl] = struct{}{}

		query := fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE`, tbl)
		if _, err := tx.Exec(ctx, query); err != nil {
			slog.Warn("Failed to drop table", "table", tbl, "err", err)
		} else {
			slog.Info("Dropped table", "table", tbl)
		}
	}

	// Step 4: Drop the finam_assets table itself (this cascades to asset_tables)
	_, err = tx.Exec(ctx, `DROP TABLE IF EXISTS finam_assets CASCADE`)
	if err != nil {
		return 0, fmt.Errorf("drop finam_assets: %w", err)
	}
	slog.Info("Dropped finam_assets and cascaded tables")

	return deletedCount, nil
}

func SingleAssetTable(jsonData []byte, token string) error {
	var req AssetRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	ctx := context.Background()
	client, err := NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return err
	}
	defer client.Close(ctx)

	switch strings.ToLower(req.Operation) {
	case "delete":
		return deleteFinamAsset(ctx, client.connPG, req)
	default:
		if strings.ToLower(req.Operation) == "update" {
			_, numb, err := FinamAssetsTable("update", token)
			if err != nil {
				return err
			}
			fmt.Printf("Update available finam_asssets: %v", numb)
		}
		if err := createAssetTable(ctx, client.connPG, req.Symbol, req.Timeframe); err != nil {
			return err
		}

	}

	bars, _, err := BarsFromFinam(ctx, client, req.Symbol, req.Timeframe, req.StartDate, req.EndDate)
	if err != nil {
		slog.Error("BarsFromFinam", "err", err)
		return err
	}

	for _, bar := range bars {
		barPgUpdate, err := ConvertBarDecimalToBarPG(bar)
		if err != nil {
			return err
		}
		if err := InsertInAssetTable(ctx, client.connPG, barPgUpdate, req); err != nil {
			return fmt.Errorf("insert data failed: %w", err)
		}
	}

	return nil
}

func createAssetTable(ctx context.Context, conn *pgx.Conn, symbol, tfStr string) error {
	table := assetTableName(symbol, tfStr, "")
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var assetID int
	err = tx.QueryRow(ctx, `SELECT id FROM finam_assets WHERE symbol = $1`, symbol).Scan(&assetID)
	if err != nil {
		return fmt.Errorf("get asset_id: %w", err)
	}

	var timeframeID int
	err = tx.QueryRow(ctx, `SELECT id FROM timeframes WHERE code = $1`, tfStr).Scan(&timeframeID)
	if err != nil {
		return fmt.Errorf("get timeframe_id: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO asset_tables (asset_id, timeframe_id, table_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (asset_id, timeframe_id)
		DO UPDATE SET updated_at = now();`, assetID, timeframeID, table)

	if err != nil {
		return fmt.Errorf("insert asset_tables: %w", err)
	}

	createSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			asset_id INT NOT NULL REFERENCES finam_assets(id) ON DELETE CASCADE,
			timeframe_id INT NOT NULL REFERENCES timeframes(id),
			timestamp TIMESTAMPTZ NOT NULL,
			open DOUBLE PRECISION,
			high DOUBLE PRECISION,
			low DOUBLE PRECISION,
			close DOUBLE PRECISION,
			volume BIGINT,
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now(),
			UNIQUE (asset_id, timeframe_id, timestamp)
		);`, table)

	if _, err := tx.Exec(ctx, createSQL); err != nil {
		return fmt.Errorf("create table %s: %w", table, err)
	}
	fmt.Printf("Created table: %s\n", table)
	return tx.Commit(ctx)
}

func deleteFinamAsset(ctx context.Context, connPG *pgx.Conn, req AssetRequest) error {
	tx, err := connPG.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var assetID int
	err = tx.QueryRow(ctx, `SELECT id FROM finam_assets WHERE symbol = $1`, req.Symbol).Scan(&assetID)
	if err != nil {
		return fmt.Errorf("get asset_id: %w", err)
	}

	var timeframeID int
	err = tx.QueryRow(ctx, `SELECT id FROM timeframes WHERE code = $1`, req.Timeframe).Scan(&timeframeID)
	if err != nil {
		return fmt.Errorf("get timeframe_id: %w", err)
	}
	var table string
	err = tx.QueryRow(ctx, `
		SELECT table_name FROM asset_tables
		WHERE asset_id = $1 AND timeframe_id = $2`,
		assetID, timeframeID).Scan(&table)
	if err != nil {
		return fmt.Errorf("get table_name: %w", err)
	}

	if _, err := tx.Exec(ctx, fmt.Sprintf(`DROP TABLE IF EXISTS %s CASCADE`, table)); err != nil {
		return fmt.Errorf("drop table %s: %w", table, err)
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM asset_tables
		WHERE asset_id = $1 AND timeframe_id = $2`,
		assetID, timeframeID); err != nil {
		return fmt.Errorf("delete from asset_tables: %w", err)
	}
	fmt.Printf("Table %s was deleted!\n", table)

	return tx.Commit(ctx)
}

func InsertInAssetTable(ctx context.Context, conn *pgx.Conn, arg BarPgUpdate, req AssetRequest) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var assetID int
	err = tx.QueryRow(ctx, `SELECT id FROM finam_assets WHERE symbol = $1`, req.Symbol).Scan(&assetID)
	if err != nil {
		return fmt.Errorf("get asset_id: %w", err)
	}

	var timeframeID int
	err = tx.QueryRow(ctx, `SELECT id FROM timeframes WHERE code = $1`, req.Timeframe).Scan(&timeframeID)
	if err != nil {
		return fmt.Errorf("get timeframe_id: %w", err)
	}

	table := assetTableName(req.Symbol, req.Timeframe, "")
	insertTableSQL := fmt.Sprintf(`
		INSERT INTO %s (
			asset_id, timeframe_id, timestamp,
			open, high, low, close, volume
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (asset_id, timeframe_id, timestamp) DO NOTHING;
	`, table)
	_, err = tx.Exec(ctx, insertTableSQL,
		assetID, timeframeID,
		arg.Timestamp, arg.Open, arg.High, arg.Low, arg.Close, arg.Volume,
	)
	if err != nil {
		return fmt.Errorf("insert table: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
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
