package main

import (
	"log"
	"net/http"

	// 自分のプロジェクトのパッケージを読み込む
	"example.com/mini-datastore/internal/db"
	"example.com/mini-datastore/internal/handler"
	"example.com/mini-datastore/internal/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// DBの初期化
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("DBの初期化に失敗しました: %v", err)
	}
	defer database.Close()

	// ハンドラの初期化
	h := &handler.NoteHandler{DB: database}

	// echoインスタンスの作成
	e := echo.New()

	// ミドルウェアの設定（Logger: ログを出す、 Recover: 落ちても復活させる）
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// CORSミドルウェアを追加
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodDelete, http.MethodPut},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, "X-User-Id"},
	}))

	// エンドポイントは全て X-User-Id のチェックを受ける
	e.Use(middleware.AuthMiddleware)

	// 生存確認
	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// メモの保存（POST）
	e.POST("/api/notes", h.CreateNote)

	// メモの一覧取得（GET）
	e.GET("/api/notes", h.GetNotes)

	// メモを1件取得
	e.GET("/api/notes/:id", h.GetNote)

	// メモを1件削除
	e.DELETE("/api/notes/:id", h.DeleteNote)

	// メモを1件更新
	e.PUT("/api/notes/:id", h.UpdateNote)

	// サーバー起動
	log.Println("サーバーをポート8080で起動します...")
	e.Logger.Fatal(e.Start(":8080"))
}
