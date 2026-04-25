package main

import (
	"log"
	"net/http"

	// 自分のプロジェクトのパッケージを読み込む
	"example.com/mini-datastore/internal/db"
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

	// echoインスタンスの作成
	e := echo.New()

	// ミドルウェアの設定（Logger: ログを出す、 Recover: 落ちても復活させる）
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// エンドポイントは全て X-User-Id のチェックを受ける
	e.Use(middleware.AuthMiddleware)

	// 最初のエンドポイント
	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// サーバー起動
	log.Println("サーバーをポート8080で起動します...")
	e.Logger.Fatal(e.Start(":8080"))
}
