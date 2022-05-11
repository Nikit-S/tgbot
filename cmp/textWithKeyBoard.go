package cmp

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (cmp *Component) TextWithKeyboard(text string, keyboard tgbotapi.ReplyKeyboardMarkup) {
	cmp.Msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(cmp.HideKeyboard)
	cmp.Msg.Text = text
	cmp.Msg.ReplyMarkup = keyboard
	cmp.Bot.Send(cmp.Msg)
}
