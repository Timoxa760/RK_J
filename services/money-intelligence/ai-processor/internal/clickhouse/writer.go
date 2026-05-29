package clickhouse

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"backend_project/services/money-intelligence/ai-processor/internal"
)

type Writer struct {
	conn driver.Conn
}

func NewWriter(ctx context.Context, host, user, password, database string) (*Writer, error) {
	addr := fmt.Sprintf("%s:9000", host)
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: database,
			Username: user,
			Password: password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:     10 * time.Second,
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("clickhouse: open: %w", err)
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("clickhouse: ping: %w", err)
	}
	log.Printf("clickhouse: connected to %s/%s", addr, database)
	return &Writer{conn: conn}, nil
}

func (w *Writer) InsertReceipt(ctx context.Context, receipt *internal.CategorizedReceipt) error {
	batch, err := w.conn.PrepareBatch(ctx, `
		INSERT INTO receipt_items (user_id, store, category, item_name, price, quantity, purchased_at, is_impulsive)
	`, driver.WithReleaseConnection())
	if err != nil {
		return fmt.Errorf("clickhouse: prepare batch: %w", err)
	}

	for _, item := range receipt.Items {
		imp := uint8(0)
		if item.IsImpulsive {
			imp = 1
		}
		if err := batch.Append(
			receipt.UserID,
			receipt.Store,
			item.Category,
			item.Name,
			item.Price,
			uint32(item.Quantity),
			receipt.Date,
			imp,
		); err != nil {
			return fmt.Errorf("clickhouse: append: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("clickhouse: send batch: %w", err)
	}
	return nil
}

func (w *Writer) Conn() driver.Conn {
	return w.conn
}

func (w *Writer) Close() error {
	return w.conn.Close()
}
