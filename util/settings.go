package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

const layout = "2006-01-02"

type ConfigFile struct {
	Yt_apikey                string `json:"yt_apikey"`
	Yt_handle                string `json:"yt_handle"`
	Line_channelSecret       string `json:"line_channelSecret"`
	Line_channelToken        string `json:"line_channelToken"`
	Message_foundChat        string `json:"message_foundChat"`
	Message_viewersIncreased string `json:"message_viewersIncreased"`
	Threshold_viewers        int    `json:"threshold_viewers"`
	Url                      struct {
		Colony_live_cam string `json:"colony_live_cam"`
		Nest_niwa       string `json:"nest_niwa"`
		Nest            string `json:"nest"`
	} `json:"url"`
}

func LoadConfig() (*ConfigFile, error) {
	configFileName, err := os.Open("./json/config.json")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("[ERROR] os.Open : %s %s : %w", file, strconv.Itoa(line), err)
	}
	defer configFileName.Close()

	raw, err := io.ReadAll(configFileName)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("[ERROR] io.ReadAll : %s %s : %w", file, strconv.Itoa(line), err)
	}

	var configFile ConfigFile
	if err := json.Unmarshal(raw, &configFile); err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("[ERROR] json.Unmarshal : %s %s : %w", file, strconv.Itoa(line), err)
	}

	return &configFile, nil
}

func SetLogDir() (*os.File, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("[ERROR] os.Getwd : %s %s : %w", file, strconv.Itoa(line), err)
	}

	logDir := currentDir + "/log"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.Mkdir(logDir, 0777)
		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("[ERROR] os.Mkdir : %s %s : %w", file, strconv.Itoa(line), err)
		}
	}

	now := time.Now()
	logfileName := logDir + "/log_" + now.Format(layout) + ".txt"

	logFile, err := os.OpenFile(logfileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("[ERROR] os.OpenFile : %s %s : %w", file, strconv.Itoa(line), err)
	}

	log.SetOutput(logFile)

	return logFile, nil
}
