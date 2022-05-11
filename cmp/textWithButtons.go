package cmp

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (cmp *Component) TextWithButtons(text string, buttons tgbotapi.InlineKeyboardMarkup) {
	cmp.Msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(cmp.HideKeyboard)
	cmp.Msg.Text = text
	cmp.Msg.ReplyMarkup = buttons
	cmp.Bot.Send(cmp.Msg)
}
