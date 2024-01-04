package pkg

import api "github.com/c1kzy/Telegram-API"

func IsResponseEmpty(body *api.WebHookReqBody) bool {
	if body.Message.Text == "" && body.Message.Location.Longitude == 0 && body.Message.Location.Latitude == 0 {
		return true
	}
	return false
}
