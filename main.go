package main

import (
	"fmt"
	"net/http"
	"time"
	"weatherbot/service"
	"weatherbot/weatherAPI"

	tgapi "github.com/c1kzy/Telegram-API"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/phuslu/log"
)

func main() {
	if envErr := godotenv.Load(); envErr != nil {
		log.Fatal().Err(fmt.Errorf("error loading .env file: %w", envErr))
	}

	cfg := &tgapi.Config{}
	wCfg := &weatherAPI.WeatherConfig{}

	log.DefaultLogger = log.Logger{
		Level:      log.InfoLevel,
		Caller:     cfg.Caller,
		TimeField:  cfg.TimeField,
		TimeFormat: time.RFC850,
		Writer:     &log.ConsoleWriter{},
	}

	if err := env.Parse(cfg); err != nil {
		log.Error().Err(err)
	}
	if wErr := env.Parse(wCfg); wErr != nil {
		log.Error().Err(wErr)
	}
	api := tgapi.GetAPI(cfg)
	weather := weatherAPI.GetWeatherAPI(wCfg)

	api.RegisterCommand("/start", service.StartResponse)
	api.RegisterInput(weather.WeatherRequest)

	http.HandleFunc("/telegram", api.TelegramHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("server start failed %v", err)
	}
	log.Info().Msg("Server started")

}
