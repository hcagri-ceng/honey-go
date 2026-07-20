package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hcagri-ceng/honey-go/internal/models"
)

// Repository, veritabanı işlemlerini soyutlayan interface'imiz.
// Ağ katmanı sadece bu arayüzü bilecek.
type Repository interface {
	SaveEvent(ctx context.Context, event models.Event) error
}

// PostgresRepo, Repository interface'ini implemente eder.
type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// SaveEvent, olayı veritabanına asenkron olarak yazar.
func (r *PostgresRepo) SaveEvent(ctx context.Context, event models.Event) error {
	query := `
		INSERT INTO honeypot_logs (timestamp, source_ip, source_port, target_port, protocol, raw_payload)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		event.Timestamp,
		event.SourceIP,
		event.SourcePort,
		event.TargetPort,
		event.Protocol,
		event.Payload,
	)

	if err != nil {
		return fmt.Errorf("olay kaydedilemedi: %w", err)
	}

	return nil
}
