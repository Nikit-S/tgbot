package cmp

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (cmp *Component) Text(text string) {
	cmp.Msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(cmp.HideKeyboard)
	cmp.Msg.Text = text
	cmp.Bot.Send(cmp.Msg)
}
