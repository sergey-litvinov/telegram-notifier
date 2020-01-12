package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var telegramSendingChannel chan string = make(chan string, 10)

func startTelegramBot() {
	telegramToken, success := os.LookupEnv("TelegramToken")
	if !success {
		log.Panicf("Telegram token isn't specified")
	}
	telegramForwardToStr, success := os.LookupEnv("TelegramForwardTo")

	if !success {
		log.Panicf("Telegram forward to isn't specified")
	}

	telegramForwardTo, err := strconv.Atoi(telegramForwardToStr)
	if err != nil {
		log.Panicf("Failed to parse %s to int", telegramForwardToStr)
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panicf("Failed to start new bot instance", err)
	}

	telegramDebugStr, success := os.LookupEnv("TelegramDebug")
	if success {
		telegramDebug, err := strconv.ParseBool(telegramDebugStr)
		if err != nil {
			log.Panicf("Failed to parse %s : %v", telegramDebugStr, err)
		}

		bot.Debug = telegramDebug
	}

	go processSending(int64(telegramForwardTo), bot, telegramSendingChannel)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panicf("Failed to pull telegram updates: %v", err)
	}

	for update := range updates {
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
			log.Panicf("Failed to send telegram message: %v", err)
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
