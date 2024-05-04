package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/TWOG710/oamaru-live-info-get/api"
	"github.com/TWOG710/oamaru-live-info-get/util"
)

func main() {
	logFile, err := util.SetLogDir()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer logFile.Close()

	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
			os.Exit(1)
		}
	}()
}

func run() error {
	conf, err := util.LoadConfig()
	if err != nil {
		return err
	}

	ctx := context.Background()
	service, err := api.CreateService(ctx, conf.Yt_apikey)
	if err != nil {
		return err
	}

	channelId, err := api.GetChannelIDFromHandle(service, conf.Yt_handle)
	if err != nil {
		return err
	}

	isLive, err := api.IsLive(service, channelId)
	if err != nil {
		return err
	}

	if !isLive {
		log.Printf("[INFO] There are no live streaming on this channel (%s).", conf.Yt_handle)
		return nil
	}

	liveVideoId, err := api.GetLiveVideoID(service, channelId)
	if err != nil {
		return err
	}

	livechatId, err := api.GetLiveChatID(service, liveVideoId)
	if err != nil {
		return err
	}

	chat, err := api.GetChat(service, livechatId)
	if err != nil {
		return err
	}

	if chat != "" {
		log.Printf("[INFO] Send live chats to LINE. : chat :%s", chat)

		msg := conf.Message_foundChat + "\n" + chat
		err := api.SendMessage(conf.Line_channelSecret, conf.Line_channelToken, msg)
		if err != nil {
			return err
		}
	} else {
		log.Print("[INFO] There are not any live chats.")
	}

	viewers, err := api.GetConcurrentViewers(service, liveVideoId)
	if err != nil {
		return err
	}

	if viewers >= conf.Threshold_viewers {
		log.Printf("[INFO] Send the number of concurrent viewers to LINE. : concurrent viewers : %s", strconv.Itoa(viewers))

		msg := conf.Message_viewersIncreased + "\nViewers:" + strconv.Itoa(viewers)
		err := api.SendMessage(conf.Line_channelSecret, conf.Line_channelToken, msg)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[INFO] Conccurent viewers under 20. : viewers : %s", strconv.Itoa(viewers))
	}

	return nil
}
