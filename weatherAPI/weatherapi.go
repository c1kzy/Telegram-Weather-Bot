package weatherAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"weatherbot/pkg"

	api "github.com/c1kzy/Telegram-API"
	"github.com/phuslu/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type WeatherAPI struct {
	GeoAPI     string
	WeatherAPI string
}

var (
	lock             = sync.Mutex{}
	SingleWeatherAPI *WeatherAPI
)

// GetWeatherAPI is getting single instance for API
func GetWeatherAPI(cfg *WeatherConfig) *WeatherAPI {
	if SingleWeatherAPI == nil {
		lock.Lock()
		defer lock.Unlock()
		if SingleWeatherAPI == nil {
			SingleWeatherAPI = &WeatherAPI{
				GeoAPI:     fmt.Sprintf("http://api.openweathermap.org/geo/%s/direct?q=%%v&limit=%v&appid=%v", cfg.GeoAPI.Version, cfg.GeoAPI.Limit, cfg.API),
				WeatherAPI: fmt.Sprintf("https://api.openweathermap.org/data/%s/weather?lat=%%v&lon=%%v&appid=%v&units=%s", cfg.WeatherVersion, cfg.API, cfg.Units),
			}
			log.Info().Msg("Weather API created")
		}
	}
	return SingleWeatherAPI
}

func (w *WeatherAPI) WeatherRequest(body *api.WebHookReqBody, chatID int) (url.Values, error) {
	var (
		weather WeatherData
		lat     float64
		lon     float64
		cityErr error
	)
	//Checking if response is empty fixed the bug when it returns the weather for the Globe when user input was empty
	if pkg.IsResponseEmpty(body) {
		return url.Values{}, fmt.Errorf("response body is nil")
	}
	lat = body.Message.Location.Latitude
	lon = body.Message.Location.Longitude

	if body.Message.Text != "" && body.Message.Text[0] != '/' {
		lat, lon, cityErr = w.getWeatherByCityName(body.Message.Text)
		if cityErr != nil {
			return url.Values{}, cityErr
		}
	}
	resp, respErr := http.Get(w.getWeatherURL(lat, lon))
	if respErr != nil {
		return url.Values{}, respErr
	}

	respBody, respBodyErr := io.ReadAll(resp.Body)
	if respBodyErr != nil {
		return url.Values{}, respBodyErr
	}

	marshalErr := json.Unmarshal(respBody, &weather)
	if marshalErr != nil {
		log.Error().Err(fmt.Errorf("error unmarshalling JSON for holiday API request Error:%w", marshalErr))
	}

	text := fmt.Sprintf("Today is %v in %v\n🌡️Temperature %v°. Feels like %v°\n💨Wind speed %v", weather.Weather[0].Description, weather.Name, int(weather.Main.Temp), int(weather.Main.FeelsLike), float32(weather.Wind.Speed))

	return url.Values{
		"chat_id": {strconv.Itoa(chatID)},
		"text":    {text},
	}, nil

}

func (w *WeatherAPI) getWeatherByCityName(text string) (float64, float64, error) {
	var location []Location

	resp, err := http.Get(fmt.Sprintf(w.GeoAPI, cases.Title(language.Und, cases.NoLower).String(text)))
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("something went wrong during api request for city coordinates%w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("unable to read response body for city provided by user: %w", err)
	}

	marshalErr := json.Unmarshal(respBody, &location)
	if marshalErr != nil {
		return 0.0, 0.0, fmt.Errorf("json unmarshal error for city provided by user: %w", marshalErr)
	}

	if len(location) == 0 {
		return 0.0, 0.0, fmt.Errorf("geo response is empty. Invalid city")
	}

	return location[0].Lat, location[0].Lon, nil
}

func (w *WeatherAPI) getWeatherURL(lat, lon float64) string {
	return fmt.Sprintf(w.WeatherAPI, lat, lon)
}
