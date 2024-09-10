package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"stockz-app/config"
	"stockz-app/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token not found")
			}
			return err
		}
		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("SECRETKEY")), nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return c.JSON(http.StatusUnauthorized, "Error : JWT Token is expired")
			}

			var user model.User
			config.DB.First(&user, claims["sub"])

			if user.ID == 0 {
				return c.JSON(http.StatusUnauthorized, "You are not authorized, Please login")
			}

			c.Set("user", user)

			return next(c)
		}
		return err
	}
}
