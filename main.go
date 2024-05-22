package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	BASEDIR = "/Users/samararora/Desktop/PROJECTS/fxb/fxb/"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	hostname := os.Getenv("HOSTNAME")
	if err := os.Setenv("HOSTNAME", hostname); err != nil {
		panic(err)
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/:baseFolder/:subFolder/:fileName", func(c echo.Context) error {
		file, err := os.Open(BASEDIR + c.Param("baseFolder") + "/" + c.Param("subFolder") + "/" + c.Param("fileName"))
		if err != nil {
			return err
		}
		defer file.Close()
		return c.Stream(200, "image/jpeg", file)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
