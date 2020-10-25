package main

import (
	tgbot "github.com/Syfaro/telegram-bot-api"
)

const (
	token = "1285311895:AAEW8x29YCw3Ux_5yx6e1wVW4VVbpVvKHg4"
)

func main() {
	bot, err := tgbot.NewBotAPI(token)
	err_handler(err)
	
	up := tgbot.NewUpdate(0)
	up.Timeout = 60
	updater, err := bot.GetUpdatesChan(up)
	err_handler(err)

	users := NewUserList()
	chanel := make(chan SendData)

	go RemindWorker(users, chanel)
	
	go sender(bot, chanel)
	
	
	M: for update := range updater{
		users.mu.Lock()
		
		for _, user := range users.users{
			if user.ID == update.Message.From.ID{
				user.MessageHandler(update)
				users.mu.Unlock()
				continue M
			} 
		}
		users.users = append(users.users, NewUser(update.Message.From.ID, update.Message.From.FirstName, chanel))
		msg := tgbot.NewMessage(update.Message.Chat.ID, "Добро пожаловать в мой органайзер)")
		bot.Send(msg)
		users.mu.Unlock()
		
	}

}

func err_handler(err error){
	if err != nil{
		panic(err)
	}
}

func sender(bot *tgbot.BotAPI, chanel chan SendData){
	for data := range chanel{
		msg := tgbot.NewMessage(int64(data.ID), data.Data)
		bot.Send(msg)
	}
}


