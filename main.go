package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	basedir := os.Getenv("BASEDIR")
	hostname := os.Getenv("HOSTNAME")
	if err := os.Setenv("HOSTNAME", hostname); err != nil {
		panic(err)
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/:baseFolder/:subFolder/:fileName", func(c echo.Context) error {
		finalPath := basedir + c.Param("baseFolder") + "/" + c.Param("subFolder") + "/" + c.Param("fileName")
		file, err := os.Open(finalPath)
		imageFormat := strings.Split(file.Name(), ".")[len(strings.Split(file.Name(), "."))-1]
		if err != nil {
			return err
		}
		defer file.Close()
		return c.Stream(200, "image/"+imageFormat, file)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
