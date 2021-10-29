package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

var telegramSendingChannel chan string = make(chan string, 10)

func startTelegramBot(config Config, ctx context.Context) {

	bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.Panicf("Failed to start new bot instance: %s", err)
	}

	bot.Debug = config.Telegram.Debug
	go processSending(config.Telegram.ForwardTo, bot, telegramSendingChannel)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panicf("Failed to pull telegram updates: %s", err)
	}

	for update := range updates {

		if ctx.Err() != nil {
			log.Printf("Stopping telegram listener...")
			return
		}

		// ignore any non-Message\non-command updates
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "ping":
			msg.Text = "pong"
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panicf("Failed to send telegram message: %s", err)
		}
	}
}

func processSending(telegramForwardTo int64, bot *tgbotapi.BotAPI, messages <-chan string) {
	for message := range messages {

		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           telegramForwardTo,
				ReplyToMessageID: 0,
			},
			Text:                  message,
			DisableWebPagePreview: false,
			ParseMode:             "Markdown",
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panicf("Failed to send telegram message: %v", err)
		}
	}
}

func sendTelegramMessage(message string) {
	telegramSendingChannel <- message
}
