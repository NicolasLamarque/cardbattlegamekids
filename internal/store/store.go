package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open() (*Store, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("user config dir: %w", err)
	}
	appDir := filepath.Join(dir, "CardsBattle")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return nil, fmt.Errorf("create app dir: %w", err)
	}
	dbPath := filepath.Join(appDir, "cardsbattle.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	// Une seule connexion et un journal classique : on écrit pour de bon à
	// chaque commit plutôt que de laisser un pool multi-connexions ou un
	// journal WAL retarder ce qui touche vraiment le disque.
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA journal_mode=DELETE; PRAGMA synchronous=FULL;`); err != nil {
		return nil, fmt.Errorf("configure pragmas: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS characters (
			id INTEGER PRIMARY KEY,
			name_full TEXT NOT NULL,
			name_native TEXT,
			image_large TEXT,
			image_medium TEXT,
			image_local TEXT,
			favourites INTEGER NOT NULL,
			role TEXT,
			series_id INTEGER,
			series_title TEXT,
			series_color TEXT,
			force INTEGER NOT NULL,
			pv INTEGER NOT NULL,
			is_custom_stats INTEGER NOT NULL DEFAULT 0,
			custom_data TEXT NOT NULL DEFAULT '{}'
		);

		CREATE TABLE IF NOT EXISTS players (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			wallet_balance INTEGER NOT NULL DEFAULT 0,
			games_won INTEGER NOT NULL DEFAULT 0,
			games_played INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS draft_modifiers (
			id TEXT PRIMARY KEY,
			label TEXT NOT NULL,
			type TEXT NOT NULL,
			kind TEXT NOT NULL,
			amount INTEGER NOT NULL,
			rare INTEGER NOT NULL DEFAULT 0
		);
	`)
	return err
}
