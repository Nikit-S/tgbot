package cmp

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Component struct {
	Bot          *tgbotapi.BotAPI
	Msg          *tgbotapi.MessageConfig
	Upd          *tgbotapi.Update
	HideKeyboard bool
	Execute      bool
}
