package main

import (
	"log"
	"os/exec"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

)



		


func main() {
	bot, err := tgbotapi.NewBotAPI("1329957296:AAHinSqFmmooAzaYa5PHK8wNOcrsKFehJYE")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)


	for update := range updates {
		if update.Message != nil && update.Message.Text == "/list" {
			cmd := exec.Command("docker", "ps")
			out, err := cmd.Output()
			//names, err := listContainers()
			if err != nil {
				log.Fatal(err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(out))
			bot.Send(msg)
		}

		/*if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}*/
	}
}