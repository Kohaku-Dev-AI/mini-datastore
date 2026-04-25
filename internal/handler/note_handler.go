package handler

import (
	"database/sql"
	"net/http"

	"example.com/mini-datastore/internal/model"
	"github.com/labstack/echo/v4"
)

// NoteHandler はDB接続を保持する構造体です
type NoteHandler struct {
	DB *sql.DB
}

// CreateNote は新しいメモを保存します（POST/api/notes)
func (h *NoteHandler) CreateNote(c echo.Context) error {
	// ミドルウェアで保存した user_id を取り出す
	userID := c.Get("user_id").(string)

	// リクエストボディ（JSON)を構造体に読み込む
	var note model.Note
	if err := c.Bind(&note); err != nil {
		return err
	}

	// DBに保存する（SQL実行）
	query := `INSERT INTO notes (user_id, title, body) VALUES (?, ?, ?)`
	result, err := h.DB.Exec(query, userID, note.Title, note.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "保存に失敗しました"})
	}

	// 保存したデータのIDを取得してレスポンスに含める
	id, _ := result.LastInsertId()
	note.ID = int(id)
	note.UserID = userID

	return c.JSON(http.StatusCreated, note)
}

// GetNotes は自分のメモいちらんを取得します（GET /api/notes）
func (h *NoteHandler) GetNotes(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// 自分のデータだけを抽出するSQL
	query := `SELECT id, user_id, title, body, created_at FROM notes WHERE user_id = ?`
	rows, err := h.DB.Query(query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 結果をスライスに詰め込む
	var notes []model.Note
	for rows.Next() {
		var n model.Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.CreatedAt); err != nil {
			return err
		}
		notes = append(notes, n)
	}
	return c.JSON(http.StatusOK, notes)
}
