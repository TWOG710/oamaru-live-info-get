package api

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

func SendMessage(channelSecret string, channelToken string, txt string) error {
	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("[ERROR] youtube.NewService : %s : %w", file+strconv.Itoa(line), err)
	}

	lineMessage := linebot.NewTextMessage(txt)
	if _, err := bot.BroadcastMessage(lineMessage).Do(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("[ERROR] youtube.NewService : %s : %w", file+strconv.Itoa(line), err)
	}

	return nil
}
