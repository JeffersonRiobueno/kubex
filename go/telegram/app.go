package main

import (
	"log"
	"context"
	"strings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)
//https://api.telegram.org/bot1329957296:AAHinSqFmmooAzaYa5PHK8wNOcrsKFehJYE/setWebhook?remove

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
			containerList, err := getContainerList()
			if err != nil {
				log.Println("Error al obtener la lista de contenedores:", err)
				continue
			}

			// Construir una lista de los nombres de los contenedores
			var containerNames []string
			for _, container := range containerList {
				containerNames = append(containerNames, container.Names[0])
			}

			// Enviar una respuesta con los nombres de los contenedores
			responseText := "Los siguientes contenedores están corriendo en la máquina host:\n\n"
			responseText += strings.Join(containerNames, "\n")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
			_, err = bot.Send(msg)
			if err != nil {
				log.Println("Error al enviar la respuesta:", err)
			}
		}
	}
}

// Obtener una lista de los contenedores en la máquina host
func getContainerList() ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return containerList, nil
}