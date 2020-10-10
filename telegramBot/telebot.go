package telebot

import (
	"fmt"
	"log"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

//TelebotAPI need to easy getting *tgbotapi.BotAPI
type TelebotAPI *tgbotapi.BotAPI

//TelebotUpdates need to easy getting tgbotapi.UpdatesChannel
type TelebotUpdates tgbotapi.UpdatesChannel

//UserData is struct to sort all information about users
type UserData struct {
	UserName     string
	ChatID       int64
	UserMessange string
}

//Users is DataBase array..
var (
	UsersData []UserData
	UsersInfo []InfoVariable
	UsersID   []int64
)

//KeyData is struct to sort object....
type KeyData struct {
	Body    TelebotAPI
	Updates tgbotapi.UpdatesChannel
}

//InfoVariable is struct to sort object....
type InfoVariable struct {
	State        string
	CurrAnsCount int
	Index        int
}

//HelloAnswer is list for "Hello" answer
var HelloAnswer = [3]string{"Привет Тварь!", "Приветствую Дружище!",
	"Отвали от меня!\nОЙ, Привет тобишь :)"} //Specialy for VAGIF :)

//BotLoop is function to set loop logic telegram bot
func BotLoop(data KeyData, botLogic func(UserData, KeyData, *InfoVariable)) {
	fmt.Println("------------------NEW_UPDATE------------------")
	for update := range data.Updates {
		if update.Message == nil {
			continue
		}
		Сhecker := false
		IndexID := 0
		for i := 0; i < len(UsersID); i++ {
			//fmt.Println(update.Message.Chat.ID, UsersID[i])
			if UsersID[i] != update.Message.Chat.ID {
				Сhecker = true
			} else {
				Сhecker = false
				IndexID = i
				break
			}
		}
		if Сhecker {
			UsersID = append(UsersID, update.Message.Chat.ID)
			UsersInfo = append(UsersInfo, InfoVariable{
				State:        "nothing",
				CurrAnsCount: 0,
				Index:        len(UsersID) - 1,
			})
			IndexID = len(UsersID) - 1
		}
		botLogic(BotGetUserInfo(update), data, &UsersInfo[IndexID])
	}
}

//BotInit is function to init telegram bot
func BotInit(TOKEN string, admin int64, debugUnable bool) KeyData {

	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Println(err)
		return KeyData{nil, nil}
	}
	bot.Debug = debugUnable
	log.Printf("Authorized on account %s", bot.Self.UserName)

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60

	updates, err := bot.GetUpdatesChan(ucfg)
	UsersInfo = append(UsersInfo, InfoVariable{
		State:        "nothing",
		CurrAnsCount: 0,
		Index:        0,
	})
	UsersID = append(UsersID, admin)

	return KeyData{bot, updates}
}

//BotGetUserInfo is function to get messange from user
func BotGetUserInfo(upd tgbotapi.Update) UserData {
	var Name string = fmt.Sprintf("%s %s",
		upd.Message.Chat.FirstName, upd.Message.Chat.LastName)
	data := UserData{
		UserName:     Name,
		ChatID:       int64(upd.Message.Chat.ID),
		UserMessange: string(upd.Message.Text),
	}
	return data
}

//BotSendMsg is function to send messange to user
func BotSendMsg(bot *tgbotapi.BotAPI, messange string, chatID int64) {
	bot.Send(tgbotapi.NewMessage(chatID, messange))
}
