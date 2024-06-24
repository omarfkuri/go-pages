package router

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const COMPILATION = "compilation.json"

type Compilation struct {
	Hash string
}

func Create[DB any](db *DB, sites... Site[DB]) (*echo.Echo, error) {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: (
				"${status}:${method} ${uri} in ${latency_human}\n" +
				"${bytes_out}B ${error}\n"),
		},
	))
	e.Use(middleware.Recover())


	e.Static("/dist", "static")

	var compilation Compilation

	err := getCompilation(&compilation)

	if err != nil {
		return e, err
	}

	globals := Globals{
		Hash: compilation.Hash,
	}

	e.Renderer = Templates()


	for _, site := range sites {
		if err = site.Setup(e, db, &globals); err != nil {
			return e, err
		}
	}

  	fmt.Println("Server initiated")
	e.Logger.Fatal(e.Start(":8080"))

	return e, nil
}

func getCompilation(comp *Compilation) error {
    
    jsonData, err := os.ReadFile(COMPILATION)
    
    if err != nil {
        return err
    }

    err = json.Unmarshal(jsonData, comp)
    if err != nil {
        return err
    }

    return nil
}