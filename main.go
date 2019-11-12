package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/corelof/umbrella_bot/userbase"
	"github.com/corelof/umbrella_bot/weather"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI
var helpMessage = "Umbrella bot can notify you about weather in Saint Petersburg. Data is provided by DarkSky API. Available commands:\n/start - Display this message\n/sub - Subscribe for notifications\n/unsub - Cancel subscription\n/help - Display this message"

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("bot_token"))
	if err != nil {
		log.Fatalln(err)
	}

	bot.Debug = false
	fmt.Printf("Authorized as %s\n", bot.Self.UserName)

	go worker()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		update := <-updates
		switch update.Message.Command() {
		case "sub":
			err := userbase.AddUser(update.Message.Chat.ID)
			if err != nil {
				fmt.Println(err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now you are subscribed")
			_, err = bot.Send(msg)
			if err != nil {
				fmt.Println(err)
			}
		case "unsub":
			err := userbase.RemoveUser(update.Message.Chat.ID)
			if err != nil {
				fmt.Println(err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now you are unsubscribed")
			_, err = bot.Send(msg)
			if err != nil {
				fmt.Println(err)
			}
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
			_, err = bot.Send(msg)
			if err != nil {
				fmt.Println(err)
			}
		case "help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
			_, err = bot.Send(msg)
			if err != nil {
				fmt.Println(err)
			}
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't understand you")
			_, err := bot.Send(msg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func worker() {
	spb, _ := time.LoadLocation("Europe/Moscow")
	for {
		var sleepDuration time.Duration
		if time.Now().In(spb).Hour() < 7 || time.Now().In(spb).Hour() == 7 && time.Now().In(spb).Minute() < 30 {
			sleepDuration = time.Minute * time.Duration(450-(time.Now().In(spb).Hour()*60+time.Now().In(spb).Minute()))
		} else {
			sleepDuration = time.Duration(450 + (24-time.Now().In(spb).Hour())*60 + (60 - time.Now().In(spb).Minute()))
		}
		time.Sleep(sleepDuration)
		go sendForecast()
	}
}

func generateForecastMessage() string {
	name, prob := weather.TodayForecast()
	var tip string
	if name == "rain" && prob >= 50 {
		tip = "Do not forget to take your umbrella"
	} else {
		tip = "Umbrella is not nessesary for today"
	}
	return fmt.Sprintf("Good morning! Probability of %s for today is %d%%. %s", name, prob, tip)
}

func sendForecast() {
	message := generateForecastMessage()
	users, err := userbase.AllUsers()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range users {
		msg := tgbotapi.NewMessage(v, message)
		_, err = bot.Send(msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
