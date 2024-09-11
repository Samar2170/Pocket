package main

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunFXBStorageServer() {
	if err := os.Setenv("HOSTNAME", hostname); err != nil {
		panic(err)
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("fxb-storage/:baseFolder/:subFolder/:fileName", GetFXBStorage)
	e.Logger.Fatal(e.Start(":8080"))
}

func GetFXBStorage(c echo.Context) error {
	finalPath := basedir + c.Param("baseFolder") + "/" + c.Param("subFolder") + "/" + c.Param("fileName")
	file, err := os.Open(finalPath)
	imageFormat := strings.Split(file.Name(), ".")[len(strings.Split(file.Name(), "."))-1]
	if err != nil {
		return err
	}
	defer file.Close()
	return c.Stream(200, "image/"+imageFormat, file)
}
