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
	// リクエストボディ（JSON)を構造体に読み込む
	var note model.Note
	if err := c.Bind(&note); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエスト形式が正しくありません"})
	}

	// タイトルチェックのバリデーション
	if note.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "タイトルは必須です"})
	}
	if len(note.Body) > 1000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "本文が長すぎます（1000文字以内）"})
	}
	// ミドルウェアで保存した user_id を取り出す
	userID := c.Request().Header.Get("X-User-Id")
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

// GetNotes は自分のメモ一覧を取得します（GET /api/notes）
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
	notes := make([]model.Note, 0)
	for rows.Next() {
		var n model.Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.CreatedAt); err != nil {
			return err
		}
		notes = append(notes, n)
	}
	return c.JSON(http.StatusOK, notes)
}

func (h *NoteHandler) GetNote(c echo.Context) error {
	id := c.Param("id")
	userID := c.Request().Header.Get("X-User-Id")

	var note model.Note
	err := h.DB.QueryRow(
		"SELECT id, user_id, title, body FROM notes WHERE id = ? AND user_id = ?",
		id, userID,
	).Scan(&note.ID, &note.UserID, &note.Title, &note.Body)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "指定されたメモは見つかりません"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "データの取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(c echo.Context) error {
	id := c.Param("id")
	userID := c.Request().Header.Get("X-User-Id")

	result, err := h.DB.Exec("DELETE FROM notes WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "データの削除に失敗しました"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "削除結果の確認に失敗しました"})
	}
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "指定されたメモは見つかりません"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *NoteHandler) UpdateNote(c echo.Context) error {
	// 入力値取得
	id := c.Param("id")
	userID := c.Request().Header.Get("X-User-Id")
	var note model.Note
	if err := c.Bind(&note); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエスト形式が正しくありません"})
	}

	// バリデーション
	if note.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "タイトルは必須です"})
	}
	if len(note.Body) > 1000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "本文が長すぎます（1000文字以内）"})
	}
	query := `UPDATE notes SET title = ?, body = ? WHERE id = ? AND user_id = ?`

	// UPDATE実行
	result, err := h.DB.Exec(query, note.Title, note.Body, id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新に失敗しました"})
	}

	// 影響桁数をチェック
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新結果の確認に失敗しました"})
	}
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "指定されたメモは見つかりません"})
	}

	// 更新後のノートを取得
	var updatedNote model.Note
	err = h.DB.QueryRow(
		"SELECT id, user_id, title, body, created_at FROM notes WHERE id = ? AND user_id = ?",
		id, userID,
	).Scan(&updatedNote.ID, &updatedNote.UserID, &updatedNote.Title, &updatedNote.Body, &updatedNote.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新後のデータ取得に失敗しました"})
	}

	// レスポンス
	return c.JSON(http.StatusOK, updatedNote)
}
