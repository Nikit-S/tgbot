package src

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var row = tgbotapi.NewInlineKeyboardRow
var urlB = tgbotapi.NewInlineKeyboardButtonURL
var switchB = tgbotapi.NewInlineKeyboardButtonSwitch
var dataB = tgbotapi.NewInlineKeyboardButtonData

var NumericInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	row(
		urlB("1.com", "http://1.com"),
		switchB("2sw", "open 2"),
		dataB("3", "3"),
	),
	row(
		dataB("4", "4"),
		dataB("5", "5"),
		dataB("6", "6"),
	),
)
