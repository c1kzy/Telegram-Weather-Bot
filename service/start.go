package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	api "github.com/c1kzy/Telegram-API"
)

// KeyboardButton struct for button text
type KeyboardButton struct {
	Text            string `json:"text"`
	Location        bool   `json:"request_location"`
	OneTimeKeyboard bool   `json:"one_time_keyboard"`
	ResizeKeyboard  bool   `json:"resize_keyboard"`
}

// ReplyKeyboardMarkup struct for keyboard layout
type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
}

func StartResponse(body *api.WebHookReqBody, chatID int) (url.Values, error) {
	replyMarkup := ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{KeyboardButton{Text: "Share my location", Location: true, OneTimeKeyboard: true, ResizeKeyboard: true}},
		},
	}
	jsonData, jsonErr := json.Marshal(replyMarkup)
	if jsonErr != nil {
		return url.Values{}, fmt.Errorf("error marshaling JSON: %w", jsonErr)
	}

	return url.Values{
		"chat_id":      {strconv.Itoa(chatID)},
		"text":         {Start},
		"reply_markup": {string(jsonData)},
	}, nil
}
