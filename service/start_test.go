package service

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	api "github.com/c1kzy/Telegram-API"
	"github.com/stretchr/testify/assert"
)

func TestStartResponse(t *testing.T) {
	type args struct {
		body *api.WebHookReqBody
	}
	reqBody := &api.WebHookReqBody{
		Message: api.Message{
			Text: "/start",
			Chat: api.Chat{
				ID: 358383178,
			},
		},
	}

	replyMarkup := ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{KeyboardButton{Text: "Share my location", Location: true, OneTimeKeyboard: true, ResizeKeyboard: true}},
		},
	}
	jsonData, _ := json.Marshal(replyMarkup)

	tests := []struct {
		name    string
		args    args
		want    url.Values
		wantErr bool
	}{
		{
			name: "/start",
			args: args{body: reqBody},
			want: url.Values{
				"chat_id":      {strconv.Itoa(358383178)},
				"text":         {Start},
				"reply_markup": {string(jsonData)},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StartResponse(tt.args.body, 358383178)
			assert.Equal(t, got, tt.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
