package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hcagri-ceng/honey-go/internal/models"
	// CGO gerektirmeyen SQLite sürücüsü (sadece init için)
	_ "modernc.org/sqlite"
)

// SQLiteRepo, Repository interface'ini implemente eder.
type SQLiteRepo struct {
	db *sql.DB
}

// NewSQLiteRepo, veritabanı bağlantısını kurar ve tabloları oluşturur.
func NewSQLiteRepo(dbPath string) (*SQLiteRepo, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("veritabanı açılamadı: %w", err)
	}

	// Hackathon pratikliği: Tablo yoksa otomatik oluştur.
	query := `
	CREATE TABLE IF NOT EXISTS honeypot_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		source_ip TEXT NOT NULL,
		source_port INTEGER NOT NULL,
		target_port INTEGER NOT NULL,
		protocol TEXT NOT NULL,
		raw_payload BLOB
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("tablo oluşturulamadı: %w", err)
	}

	return &SQLiteRepo{db: db}, nil
}

// SaveEvent, olayı SQLite veritabanına yazar.
func (r *SQLiteRepo) SaveEvent(ctx context.Context, event models.Event) error {
	query := `
		INSERT INTO honeypot_logs (source_ip, source_port, target_port, protocol, raw_payload)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
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

func (r *SQLiteRepo) GetEvents() ([]models.Event, error) {
	query := `
		SELECT id, timestamp, source_ip, source_port, target_port, protocol, raw_payload
		FROM honeypot_logs
		ORDER BY timestamp DESC
		LIMIT 100
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("veriler okunamadı: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Timestamp, &e.SourceIP, &e.SourcePort, &e.TargetPort, &e.Protocol, &e.Payload); err != nil {
			continue // Bozuk satırı atla
		}
		events = append(events, e)
	}

	return events, nil
}
