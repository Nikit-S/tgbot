package cmp

import (
	"database/sql"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	BotApi *tgbotapi.BotAPI
	Msg    *tgbotapi.MessageConfig
	Upd    *tgbotapi.Update
	Db     *sql.DB
}

func (s *Bot) Execute(c Component) {
	c.Execute(s)
}
func (b *Bot) CheckForRoot(update tgbotapi.Update) bool {
	var name string
	row, _ := b.Db.Query(`SELECT u_name from users where user_t = 'root'`)
	for row.Next() {
		row.Scan(&name)
		if name == update.Message.Text {
			return true
		}
	}
	return false
}

func (b *Bot) CheckForUser(update tgbotapi.Update) bool {
	var name string
	row, _ := b.Db.Query(`SELECT u_name from users where user_t = 'trainee'`)
	for row.Next() {
		row.Scan(&name)
		if name == update.Message.Text {
			return true
		}
	}
	return false
}

func (b *Bot) ShowAllUsersTo(id int64) {
	var name string
	row, _ := b.Db.Query(`SELECT u_name from public.users`)
	for row.Next() {
		row.Scan(&name)
		b.BotApi.Send(tgbotapi.NewMessage(id, name))
	}
}

func (b *Bot) SetIds(update tgbotapi.Update) {
	fmt.Println("Setting id for root: ", update.Message.Text)
	res, err := b.Db.Exec(`UPDATE public.users
	SET chat_id=$1,user_id=$2, username=$3
	WHERE u_name=$4;`,
		update.Message.Chat.ID,
		update.Message.From.ID,
		update.Message.Chat.UserName,
		update.Message.Text)
	fmt.Println(res, err)
}

func (b *Bot) Notify(query, msg string) {
	row, _ := b.Db.Query(query)
	var chatId int64
	for row.Next() {
		row.Scan(&chatId)
		b.BotApi.Send(tgbotapi.NewMessage(chatId, msg))
	}
}
