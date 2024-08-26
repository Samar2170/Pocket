package main

import (
	"log"
	"os"
	"pocket/internal"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	args := os.Args[1:]
	switch args[0] {
	case "load":
		log.Println("loading accounts")
		err := internal.LoadAccounts()
		if err != nil {
			log.Println(err)
		}
	case "server":
		RunStorageServer()
	default:
		RunTelegramServer()
	}
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	RunStorageServer()
	// 	wg.Done()
	// }()

	// wg.Add(1)
	// go func() {
	RunTelegramServer()
	// 	wg.Done()
	// }()
	// wg.Wait()
}

func RunStorageServer() {
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
