package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	dbPath := filepath.Join("data", "app.db")

	// 0755：フォルダの権限設定
	// 自分が読み書きOK、他人は読み込みのみ許可
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, fmt.Errorf("dataディレクトリの作成に失敗しました： %w", err)
	}

	// sqliteドライバを使ってファイルを開く
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("dataフォルダのapp.dbファイルを開くことに失敗しました： %w", err)
	}
	// テーブル作成。
	// 2回目以降の起動でテーブルがすでに存在している場合は、作成しないようにする。
	query := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		body TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("notesテーブルの作成に失敗しました： %w", err)
	}
	return db, nil
}
