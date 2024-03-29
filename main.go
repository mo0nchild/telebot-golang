package main

import (
	"fmt"
	"log"
	"strings"

	htmlparser "github.com/mo0nchild/telebot-golang/gethtml"
	telebot "github.com/mo0nchild/telebot-golang/telegramBot"
)

var (
	//BotToken is using to access the HTTP API
	BotToken string = "your bot token"
	//AdminID is using to get admin privilege
	AdminID  int64 = "your admin id"
	messange string
)

//TeleBotMain is function for loop processing of the bot
func TeleBotMain(data telebot.UserData, key telebot.KeyData, Info *telebot.InfoVariable) {

	fmt.Printf("USERNAME: %s\tCHAT_ID: %d\tMESSANGE: %s\n",
		data.UserName, data.ChatID, data.UserMessange)

	if Info.State == "nothing" {
		switch strings.ToLower(data.UserMessange) {

		case "привет":
			messange = telebot.HelloAnswer[Info.CurrAnsCount]
			if Info.CurrAnsCount >= len(telebot.HelloAnswer)-1 {
				Info.CurrAnsCount = 0
			} else {
				Info.CurrAnsCount++
			}
		case "/find_info":
			messange = "Напишите, что хотите найти: "
			Info.State = "/find_info"

		case "/about":
			messange = "Я бот, написанный на языке GO{GoLang}.\nМой создатель [Mo0nChilld]\n"
			messange += "Он еще дорабатывает некоторые функции во мне :)"

		default:
			messange = "Не понимаю. Попробуй поздороваться со мной\n"
			messange += "Используй ключевое слово \"Привет\"\n"
			messange += "\nТакже ты можешь использовать команды из списка :)"
		}

		telebot.BotSendMsg(key.Body, messange, data.ChatID)

	} else {
		switch Info.State {
		case "/find_info":
			forSearch := data.UserMessange
			var info string = htmlparser.WikiParser(forSearch)
			log.Println(info)
			telebot.BotSendMsg(key.Body, info, data.ChatID)
			Info.State = "nothing"
		}
	}
	fmt.Printf("STATE: %s\tCURRENT_ANSWER_COUNT: %d\tINDEX_ID: %d\n",
		Info.State, Info.CurrAnsCount, Info.Index)
}

//TeleBotActivation is function to start telegram BOT
func TeleBotActivation() {
	bot := telebot.BotInit(BotToken, AdminID, false)
	if (bot.Body == nil) && (bot.Updates == nil) {
		log.Fatal("Can't Connect to Telegram Servers")
	}
	telebot.BotLoop(bot, func(data telebot.UserData, key telebot.KeyData,
		info *telebot.InfoVariable) {
		TeleBotMain(data, key, info)
	})
}

func main() {
	TeleBotActivation()
}
