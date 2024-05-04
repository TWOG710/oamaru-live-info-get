package api

import (
	"context"
	"fmt"
	"runtime"
	"strconv"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func CreateService(ctx context.Context, apiKey string) (*youtube.Service, error) {
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("[ERROR] youtube.NewService : %s : %w", file+strconv.Itoa(line), err)
	}

	return service, nil
}

func GetChannelIDFromHandle(service *youtube.Service, handle string) (string, error) {
	call := service.Search.List([]string{"id", "snippet"}).Q(handle)
	response, err := call.Do()
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return "", fmt.Errorf("[ERROR] service.Search.List : %s : %w", file+strconv.Itoa(line), err)
	}

	if len(response.Items) == 0 {
		_, file, line, _ := runtime.Caller(0)
		return "", fmt.Errorf("[ERROR] Couldn't get any response.Items. : %s", file+strconv.Itoa(line))
	}

	channelId := response.Items[0].Id.ChannelId
	if channelId == "" {
		_, file, line, _ := runtime.Caller(0)
		return "", fmt.Errorf("[ERROR] Couldn't get any channelId from the handle (%s). : %s", handle, file+strconv.Itoa(line))
	}

	return channelId, nil
}

func GetLiveVideoID(service *youtube.Service, channelId string) (string, error) {
	call := service.Search.List([]string{"id", "snippet"}).ChannelId(channelId).Type("video").EventType("Live").MaxResults(10)
	response, err := call.Do()
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return "", fmt.Errorf("[ERROR] service.Search.List : %s : %w", file+strconv.Itoa(line), err)
	}

	if len(response.Items) == 0 {
		_, file, line, _ := runtime.Caller(0)
		return "", fmt.Errorf("[ERROR] Couldn't get any response.Items. : %s", file+strconv.Itoa(line))
	}

	liveVideoId := response.Items[0].Id.VideoId
	return liveVideoId, nil
}

func GetLiveChatID(service *youtube.Service, liveVideoId string) (string, error) {
	call := service.Videos.List([]string{"liveStreamingDetails"}).Id(liveVideoId)
	response, err := call.Do()
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return "", fmt.Errorf("[ERROR] service.Search.List : %s : %w", file+strconv.Itoa(line), err)
	}

	if len(response.Items) == 0 {
		_, file, line, _ := runtime.Caller(0)
		return "", fmt.Errorf("[ERROR] Couldn't get any response.Items. : %s", file+strconv.Itoa(line))
	}

	liveChatId := response.Items[0].LiveStreamingDetails.ActiveLiveChatId
	return liveChatId, nil
}

func GetChat(service *youtube.Service, liveChatId string) (string, error) {
	// 1回のみ取得するためpagetokenは使用しない
	call := service.LiveChatMessages.List(liveChatId, []string{"snippet"})
	response, err := call.Do()
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return "", fmt.Errorf("[ERROR] service.LiveChatMessages.List : %s : %w", file+strconv.Itoa(line), err)
	}

	if response.PageInfo.TotalResults == 0 {
		return "", nil
	}

	if len(response.Items) == 0 {
		_, file, line, _ := runtime.Caller(0)
		return "", fmt.Errorf("[ERROR] Couldn't get any response.Items. : %s", file+strconv.Itoa(line))
	}

	chatText := response.Items[0].Snippet.DisplayMessage
	chatPubDate := response.Items[0].Snippet.PublishedAt
	return chatPubDate + chatText, nil
}

func GetConcurrentViewers(service *youtube.Service, liveVideoId string) (int, error) {
	call := service.Videos.List([]string{"liveStreamingDetails"}).Id(liveVideoId)
	response, err := call.Do()
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return 0, fmt.Errorf("[ERROR] service.Videos.List : %s : %w", file+strconv.Itoa(line), err)
	}

	if len(response.Items) == 0 {
		_, file, line, _ := runtime.Caller(0)
		return 0, fmt.Errorf("[ERROR] Couldn't get any response.Items. : %s", file+strconv.Itoa(line))
	}

	viewers := int(response.Items[0].LiveStreamingDetails.ConcurrentViewers)
	return viewers, nil
}

func IsLive(service *youtube.Service, channelId string) (bool, error) {
	call := service.Search.List([]string{"snippet"}).ChannelId(channelId).Type("channel").MaxResults(1)
	response, err := call.Do()
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return false, fmt.Errorf("[ERROR] service.Search.List : %s : %w", file+strconv.Itoa(line), err)
	}

	if len(response.Items) == 0 {
		_, file, line, _ := runtime.Caller(0)
		return false, fmt.Errorf("[ERROR] Couldn't get any response.Items. : %s", file+strconv.Itoa(line))
	}

	liveState := response.Items[0].Snippet.LiveBroadcastContent
	if liveState == "live" {
		return true, nil
	} else {
		return false, nil
	}
}
