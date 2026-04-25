package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware は HTTPヘッダーに "X-User-Id" があるかチェックする門番
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// リクエストヘッダーから "X-User-Id" を取り出す
		userID := c.Request().Header.Get("X-User-Id")
		if userID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "X-User-Id header is required（なりすましIDが必要です）",
			})
		}
		// IDがあれば、後の処理が使いやすいように「コンテキスト」に保存しておく
		c.Set("user_id", userID)

		// 次の処理（実際のAPI処理）へ進む
		return next(c)
	}
}
