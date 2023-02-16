package middlewares

import (
	"movie_searcher/databases"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type DatabaseClient struct {
	DB *gorm.DB
}

func DatabaseService() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, _ := database.Connect()
			d := DatabaseClient{DB: session}

			defer d.DB.Close()

			d.DB.LogMode(true)
			c.Set("dbs", &d)

			if err := next(c); err != nil {
				return err
			}

			return nil
		}
	}
}
