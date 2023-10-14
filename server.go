package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/addetz/go-weather-checker/apis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	TIMEOUT = 3 * time.Second
)

func main() {
	// Read port if one is set
	port := readPort()

	// Read API key
	weatherService := apis.NewWeatherService()

	// Initialise echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configure server
	s := http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           e,
		ReadTimeout:       TIMEOUT,
		ReadHeaderTimeout: TIMEOUT,
		WriteTimeout:      TIMEOUT,
		IdleTimeout:       TIMEOUT,
	}

	// Set up the root file
	e.Static("/", "layout")
	// Set up scripts
	e.File("/scripts.js", "scripts/scripts.js")
	e.File("/scripts.js.map", "scripts/scripts.js.map")

	e.GET("/weather/:city", func(c echo.Context) error {
		city := c.Param("city")
		response, err := weatherService.GetData(city)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, &apis.BackendResponse{
			Message:     fmt.Sprintf("Fetched data for %s", city),
			FeelsLike:   apis.ConvertCelsius(response.Main.FeelsLike),
			Temp:        apis.ConvertCelsius(response.Main.Temp),
			CityName:    city,
			Description: response.Weather[0].Description,
		})
	})

	log.Printf("Listening on :%s...\n", port)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

// readPort reads the SERVER_PORT environment variable if one is set
// or returns a default if none is found
func readPort() string {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return "8080"
	}
	return port
}
