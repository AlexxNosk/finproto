package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
)

// Formats dynamic table name, e.g., "SBER@MISX", "D" → "bars_sber_misx_d"
func assetTableName(symbol, tfStr, name string) string {
	s := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(symbol, "@", "_"), "-", "_"))
	t := strings.ToLower(tfStr)
	if n := strings.ToLower(name); n != "" {
		return fmt.Sprintf("bars_%s_%s_%s", s, t, n)
	}
	return fmt.Sprintf("bars_%s_%s", s, t)
}

// type AssetRequest struct {
// 	Request RequestParameters            `json:"request"`
// 	Tables  map[string]map[string]string `json:"tables"` // table name → column name → type
// }

// type RequestParameters struct {
// 	Symbol    string `json:"symbol"`
// 	Timeframe string `json:"timeframe"`
// 	StartDate string `json:"start_date"`
// 	EndDate   string `json:"end_date"`
// 	Operation string `json:"operation"` // "create", "update", "delete"
// }

type AssetRequest struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Operation string `json:"operation"` // "create", "update", "delete"
}

func AssetCreate(jsonData []byte, token string) error {
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
		if err := InsertDataInTables(ctx, client.connPG, barPgUpdate, req); err != nil {
			return fmt.Errorf("insert data failed: %w", err)
		}
	}

	// if err := InsertInfoInTables(ctx, client.connPG, req); err != nil {
	// 	return fmt.Errorf("insertion info in tables failed: %w", err)
	// }

	return nil
}

// func DynamicDataTables(ctx context.Context, conn *pgx.Conn, req AssetRequest) error {
// 	tx, err := conn.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("begin tx: %w", err)
// 	}
// 	defer tx.Rollback(ctx)

// 	var instrumentID int
// 	err = tx.QueryRow(ctx, `SELECT id FROM instruments WHERE symbol = $1`, req.Request.Symbol).Scan(&instrumentID)
// 	if err != nil {
// 		return fmt.Errorf("get asset_id: %w", err)
// 	}
// 	var timeframeID int
// 	err = tx.QueryRow(ctx, `SELECT id FROM timeframes WHERE code = $1`, req.Request.Timeframe).Scan(&timeframeID)
// 	if err != nil {
// 		return fmt.Errorf("get timeframe_id: %w", err)
// 	}

// 	for name, newTableParameters := range req.Tables {
// 		newTableName := assetTableName(req.Request.Symbol, req.Request.Timeframe, name)
// 		colDefs := []string{
// 			"id SERIAL PRIMARY KEY",
// 			"asset_id INT NOT NULL REFERENCES instruments(id) ON DELETE CASCADE",
// 			"timeframe_id INT NOT NULL REFERENCES timeframes(id)",
// 			"timestamp TIMESTAMPTZ",
// 		}

// 		for name, typeStr := range newTableParameters {
// 			colType := mapJSONTypeToPostgres(typeStr)
// 			if colType == "" {
// 				return fmt.Errorf("unsupported type: %s", typeStr)
// 			}
// 			colDefs = append(colDefs, fmt.Sprintf("%s %s", name, colType))
// 		}
// 		colDefs = append(colDefs, "created_at TIMESTAMPTZ DEFAULT now()", "UNIQUE (asset_id, timeframe_id, timestamp)")

// 		createSQL := fmt.Sprintf(`
// 			CREATE TABLE IF NOT EXISTS %s (
// 				%s
// 			);`, newTableName, strings.Join(colDefs, ",\n"))

// 		if _, err := tx.Exec(ctx, createSQL); err != nil {
// 			return fmt.Errorf("create table %s: %w", newTableName, err)
// 		}

// 		fmt.Println("Created new data table:", newTableName)
// 	}
// 	return tx.Commit(ctx)
// }

// func mapJSONTypeToPostgres(jsonType string) string {
// 	switch strings.ToLower(jsonType) {
// 	case "int", "int64":
// 		return "BIGINT"
// 	case "float", "float64":
// 		return "DOUBLE PRECISION"
// 	case "string":
// 		return "TEXT"
// 	case "bool", "boolean":
// 		return "BOOLEAN"
// 	default:
// 		return ""
// 	}
// }

// Creates asset_tables entry + dynamic table, wrapped in a transaction
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

	// Delete data_tables
	// rows, err := tx.Query(ctx, `
	// SELECT table_name FROM data_tables
	// WHERE asset_id = $1 AND timeframe_id = $2`,
	// 	assetID, timeframeID)
	// if err != nil {
	// 	return fmt.Errorf("query data_tables: %w", err)
	// }
	// defer rows.Close()

	// var tablesToDrop []string

	// for rows.Next() {
	// 	var dataTable string
	// 	if err := rows.Scan(&dataTable); err != nil {
	// 		return fmt.Errorf("scan data_table name: %w", err)
	// 	}
	// 	tablesToDrop = append(tablesToDrop, dataTable)
	// }

	// if err := rows.Err(); err != nil {
	// 	return fmt.Errorf("row iteration error: %w", err)
	// }

	// for _, dataTable := range tablesToDrop {
	// 	dropSQL := fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE`, dataTable)

	// 	if _, err := tx.Exec(ctx, dropSQL); err != nil {
	// 		return fmt.Errorf("drop data_table %s: %w", dataTable, err)
	// 	}
	// 	_, err = tx.Exec(ctx, `
	// 		DELETE FROM data_tables
	// 		WHERE asset_id = $1 AND timeframe_id = $2 AND table_name = $3`,
	// 		assetID, timeframeID, dataTable)
	// 	if err != nil {
	// 		return fmt.Errorf("delete metadata for %s: %w", dataTable, err)
	// 	}
	// 	fmt.Printf("Table %s was deleted!\n", dataTable)
	// }

	// Delete asset_tables
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

// Inserts one bar into the dynamic table for symbol + timeframe
func InsertDataInTables(ctx context.Context, conn *pgx.Conn, arg BarPgUpdate, req AssetRequest) error {
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

	// for name := range req.Tables {
	// 	newTableName := assetTableName(req.Request.Symbol, req.Request.Timeframe, name)
	// 	insertNewTableSQL := fmt.Sprintf(`
	// 		INSERT INTO %s (
	// 			asset_id, timeframe_id, timestamp
	// 		) VALUES ($1, $2, $3)
	// 		ON CONFLICT (asset_id, timeframe_id, timestamp) DO NOTHING;
	// 	`, newTableName)
	// 	_, err = tx.Exec(ctx, insertNewTableSQL, assetID, timeframeID, arg.Timestamp)
	// 	if err != nil {
	// 		return fmt.Errorf("insert data table: %w", err)
	// 	}
	// }

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
}

// func InsertInfoInTables(ctx context.Context, conn *pgx.Conn, req AssetRequest) error {
// 	tx, err := conn.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("begin tx: %w", err)
// 	}
// 	defer tx.Rollback(ctx)

// 	var assetID int
// 	err = tx.QueryRow(ctx, `SELECT id FROM instruments WHERE symbol = $1`, req.Request.Symbol).Scan(&assetID)
// 	if err != nil {
// 		return fmt.Errorf("get asset_id: %w", err)
// 	}

// 	var timeframeID int
// 	err = tx.QueryRow(ctx, `SELECT id FROM timeframes WHERE code = $1`, req.Request.Timeframe).Scan(&timeframeID)
// 	if err != nil {
// 		return fmt.Errorf("get timeframe_id: %w", err)
// 	}

// 	var tableID int
// 	err = tx.QueryRow(ctx, `
// 		SELECT id FROM asset_tables
// 		WHERE asset_id = $1 AND timeframe_id = $2`, assetID, timeframeID).Scan(&tableID)
// 	if err != nil {
// 		return fmt.Errorf("get timeframe_id: %w", err)
// 	}

// 	for name := range req.Tables {
// 		newTableName := assetTableName(req.Request.Symbol, req.Request.Timeframe, name)
// 		insertDataTableInfoSQL := `
// 			INSERT INTO data_tables (
// 				asset_id, timeframe_id, instrument_table_id, table_name, purpose
// 			) VALUES ($1, $2, $3, $4, $5)
// 			ON CONFLICT (asset_id, timeframe_id, table_name) DO NOTHING;
// 		`
// 		_, err = tx.Exec(ctx, insertDataTableInfoSQL, assetID, timeframeID, tableID, newTableName, name)
// 		if err != nil {
// 			return fmt.Errorf("insert data table: %w", err)
// 		}

// 		fmt.Println("inserted data table:", newTableName)
// 	}

// 	if err := tx.Commit(ctx); err != nil {
// 		return fmt.Errorf("commit failed: %w", err)
// 	}
// 	return nil
// }
