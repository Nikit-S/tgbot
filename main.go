package main

import (
	"fmt"
	"log"
	"net/http"
	"testbot/cmp"
	"testbot/src"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var baseMsg = map[string]string{
	"start": `Это тестовый бот, который пока что ничего не умеет, но скоро всему научится`,
}

var commands = tgbotapi.NewSetMyCommands(
	tgbotapi.BotCommand{Command: "/hello", Description: "Пишет hello"},
	tgbotapi.BotCommand{Command: "/open", Description: "Открывает меню"},
)

func main() {

	bot, err := tgbotapi.NewBotAPI("...")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	if _, err := bot.Request(commands); err != nil {
		log.Fatal("Unable to set commands")
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	tgbotapi.NewWebhook("..." + bot.Token)

	//_, err = bot.SetWebhook(wh)
	//if err != nil {
	//	log.Fatal(err)
	//}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("127.0.0.1:8443", nil)

	cmp := cmp.Component{
		Msg:          &tgbotapi.MessageConfig{},
		Bot:          bot,
		HideKeyboard: true,
	}
	for update := range updates {
		cmp.Upd = &update
		if update.CallbackQuery != nil {
			fmt.Print(update)

			bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Получено"))
			cmp.Msg.ChatID = update.CallbackQuery.Message.Chat.ID
			cmp.Text(update.CallbackQuery.Data)
		}
		if update.Message != nil {
			cmp.Msg.ChatID = update.Message.Chat.ID
			//cmp.TextWithKeyboard("Hello", numericReplyKeyboard)
			cmp.TextWithButtons("Hello", src.NumericInlineKeyboard)
			//			msg.ChatID = update.Message.Chat.ID
			//			msg.Text = update.Message.Text
			//			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			//			switch update.Message.Command() {
			//			case "open":
			//				msg.ReplyMarkup = numericReplyKeyboard
			//				msg.Text = "Nav menu"
			//			case "start":
			//				msg.Text = baseMsg["start"]
			//
			//			}

		}
	}
}
