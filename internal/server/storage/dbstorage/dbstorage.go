package dbstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type DBStorage struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

const maxRetries = 3
const retryDelay = 1 * time.Second

func New(dbConn *pgxpool.Pool, log *zap.Logger) *DBStorage {
	return &DBStorage{db: dbConn, logger: log}
}

func GetConnect(connStr string, logger *zap.Logger) (*pgxpool.Pool, error) {
	var db *pgxpool.Pool
	var err error
	delay := retryDelay

	for i := 0; i < maxRetries; i++ {
		db, err = pgxpool.Connect(context.Background(), connStr)
		if err == nil {
			break
		}

		logger.Error("error connect to db. retry...", zap.Error(err))
		delay += 2 * time.Second
		time.Sleep(delay)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}

func (d *DBStorage) Ping() error {
	fmt.Println("ping")
	if d.db == nil {
		return fmt.Errorf("database is not connected")
	}
	return d.db.Ping(context.Background())
}

func (d *DBStorage) SetGauge(metric dto.PostMetricDto) error {
	query := `INSERT INTO metrics (name, metric_type, value) VALUES ($1, $2, $3)`
	_, err := d.db.Exec(context.Background(), query, metric.ID, metric.MType, metric.Value)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorage) SetCounter(metric dto.PostMetricDto) error {
	query := `INSERT INTO metrics (name, metric_type, delta, value) VALUES ($1, $2, $3, $4)`
	_, err := d.db.Exec(context.Background(), query, metric.ID, metric.MType, metric.Delta, metric.Value)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorage) Get(name string) (*models.Metric, error) {
	query := `SELECT * FROM metrics WHERE name = $1 ORDER BY created_at DESC LIMIT 1`
	row := d.db.QueryRow(context.Background(), query, name)

	var metric models.Metric
	err := row.Scan(&metric.ID, &metric.Name, &metric.Type, &metric.Value, &metric.Delta, &metric.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &metric, nil
}

func (d *DBStorage) GetAll() (*[]models.Metric, error) {
	query := `SELECT id, name, metric_type, value, delta, created_at FROM metrics ORDER BY created_at DESC LIMIT 27`
	rows, err := d.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.Metric
	for rows.Next() {
		var metric models.Metric
		err := rows.Scan(&metric.ID, &metric.Name, &metric.Type, &metric.Value, &metric.Delta, &metric.CreatedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return &metrics, nil
}

func (d *DBStorage) SetBatch(body []dto.PostMetricDto) error {
	tx, err := d.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, metric := range body {
		if metric.MType == models.CounterType {
			pollCount, _ := d.Get(metric.ID)

			if pollCount != nil {
				var newValue int64
				if pollCount.Delta != nil && metric.Delta != nil {
					newValue = *pollCount.Delta + *metric.Delta
				}
				_, err := tx.Exec(context.Background(), `UPDATE metrics SET delta = $1 WHERE name = $2`, newValue, metric.ID)
				if err != nil {
					return err
				}
				continue
			}

			_, err := tx.Exec(context.Background(), `INSERT INTO metrics (name, metric_type, delta) VALUES ($1, $2, $3)`, metric.ID, metric.MType, metric.Delta)
			if err != nil {
				return err
			}
			continue
		}

		_, err := tx.Exec(context.Background(), `INSERT INTO metrics (name, metric_type, value) VALUES ($1, $2, $3)`, metric.ID, metric.MType, metric.Value)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (d *DBStorage) Load(filePath string) error {
	return nil
}

func (d *DBStorage) Save(filePath string) error {
	return nil
}

func (d *DBStorage) Close() error {
	if d.db == nil {
		return nil
	}
	d.db.Close()
	return nil
}

func MigrateSQL(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS metrics (
		id SERIAL PRIMARY KEY,
		metric_type TEXT NOT NULL,
		name TEXT NOT NULL,
		value DOUBLE PRECISION,
		delta BIGINT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_metrics_name ON metrics (name);`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}
