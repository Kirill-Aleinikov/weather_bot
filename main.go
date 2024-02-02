package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Gbot *tgbotapi.BotAPI
var Token string
var ChatId int64

func init() {
	_ = os.Setenv("TOKEN_NAME_IN_OS", "6838765270:AAEwcSFxYsfdA-QXI8j_bwazdrNfqC5vOi4")
	var err error
	Token = os.Getenv("TOKEN_NAME_IN_OS")

	if Gbot, err = tgbotapi.NewBotAPI(Token); err != nil {
		log.Panic(err)
	}

	Gbot.Debug = true

}
func delay(seconds uint8) {
	time.Sleep(time.Second * time.Duration(seconds))
}
func isStartMessage(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text == "/start"
}

func printIntro(update *tgbotapi.Update) {
	delay(1 / 2)
	Gbot.Send(tgbotapi.NewMessage(ChatId,
		"Привет! Я бот, который поможет тебе узнать погоду в твоем городе. Просто напиши мне название своего города, и я предоставлю тебе актуальную информацию о погоде. Например, 'Москва'."))
}

func isWeatherMessage(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text == "/weather"
}
func printWeather(update *tgbotapi.Update) {
	delay(1 / 2)
	Gbot.Send(tgbotapi.NewMessage(ChatId, "С пивком пойдет"))
}
func istemperature(update *tgbotapi.Update) bool {
	return update.Message != nil
}

type WeatherData struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func temperature(update *tgbotapi.Update) {
	var city string
	apiKey := "19c4c216d1a54908478fa1ade12e579d"
	city = update.Message.Text

	unit := "metric"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=%s&APPID=%s", city, unit, apiKey)

	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	if resp.StatusCode() != 200 {
		if resp.StatusCode() == 404 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Город: %s не найден, попробуйте указать город еще раз", city))
			_, err = Gbot.Send(msg)
			if err != nil {
				log.Println("Error sending message:", err)
			}
			return
		}
		log.Println("Error:", resp.String())
		return
	}

	var weather WeatherData
	err = json.Unmarshal(resp.Body(), &weather)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	roundedTemp := math.Round(weather.Main.Temp)
	if roundedTemp >= 24 && roundedTemp <= 30 {
		tempt := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("На улице жарко"))
		_, err = Gbot.Send(tempt)
	}
	if roundedTemp >= 17 && roundedTemp < 24 {
		tempt := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("На улице комфортно"))
		_, err = Gbot.Send(tempt)
	}
	if roundedTemp <= 10 {
		tempt := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("На улице холодно"))
		_, err = Gbot.Send(tempt)
	}
	if roundedTemp >= 32 {
		tempt := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Лучше не выходить на улицу, слишком жарко"))
		_, err = Gbot.Send(tempt)
	}
	delay(1)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Температура в %s: %.0f°C\n", strings.Title(city), roundedTemp))
	_, err = Gbot.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
	}

}

func main() {

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := Gbot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if isStartMessage(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			ChatId = update.Message.Chat.ID
			printIntro(&update)
			break
		}
		if isWeatherMessage(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			ChatId = update.Message.Chat.ID
			printWeather(&update)
			break

		}
		if istemperature(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			ChatId = update.Message.Chat.ID
			temperature(&update)

		}

	}

}
