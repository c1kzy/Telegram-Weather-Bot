package weatherAPI

import (
	"net/url"
	"strconv"
	"testing"
	"time"

	api "github.com/c1kzy/Telegram-API"
	"github.com/phuslu/log"
	"github.com/stretchr/testify/assert"
)

var weatherConfig WeatherConfig

func requestBody(t *testing.T, city string) *api.WebHookReqBody {
	reqBody := &api.WebHookReqBody{
		Message: api.Message{
			Text: city,
			Chat: api.Chat{
				ID: 358383178,
			},
		},
	}
	return reqBody
}

func TestWeatherAPI_WeatherRequest(t *testing.T) {
	log.DefaultLogger = log.Logger{
		Level:      log.DebugLevel,
		Caller:     1,
		TimeFormat: time.RFC850,
		Writer:     &log.ConsoleWriter{},
	}

	weatherAPI := GetWeatherAPI(&weatherConfig)

	tests := []struct {
		name     string
		have     url.Values
		city     string
		expected url.Values
	}{
		{
			name: "Test case",
			city: "Kharkiv",
			have: url.Values{
				"chat_id": {strconv.Itoa(358383178)},
				"text":    {"Today is light snow in Kharkiv"},
			},
			expected: url.Values{
				"chat_id": {strconv.Itoa(358383178)},
				"text":    {"Today is light snow in Kharkiv"},
			},
		},

		{
			name: "Empty test case",
			city: "",
			have: url.Values{
				"chat_id": {strconv.Itoa(358383178)},
				"text":    {""},
			},
			expected: url.Values{
				"chat_id": {strconv.Itoa(358383178)},
				"text":    {""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weatherAPI.WeatherRequest(requestBody(t, tt.city), 358383178)
			assert.Equal(t, tt.have, tt.expected)
		})
	}
}
