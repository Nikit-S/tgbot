package src

import (
	"fmt"
	"log"
	"testbot/cmp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var commands = tgbotapi.NewSetMyCommandsWithScope(
	tgbotapi.BotCommandScope{Type: "chat", ChatID: -1},
	tgbotapi.BotCommand{Command: "/register", Description: "Попробовать войти снова"},
)

var rootCommands = tgbotapi.NewSetMyCommandsWithScope(
	tgbotapi.BotCommandScope{Type: "chat", ChatID: -1},
	tgbotapi.BotCommand{Command: "/show_all_users", Description: "Показывает всю таблицу"},
	tgbotapi.BotCommand{Command: "/notify_l0", Description: "Выслать ссылку на сдачу l0"},
	tgbotapi.BotCommand{Command: "/notify_l1", Description: "Выслать ссылку на сдачу l1"},
	tgbotapi.BotCommand{Command: "/notify_l2", Description: "Выслать ссылку на сдачу l2"},
)

func StartChatWithUser(b *cmp.Bot, ch *cmp.Chat) {
	fmt.Println("begin")
	update := <-ch.Updates
	fmt.Println("got update and executing StartChatWithUser")

	if update.Message.Command() == "start" {

		fmt.Println("START")

		if _, err := b.BotApi.Request(tgbotapi.NewDeleteMyCommandsWithScope(tgbotapi.BotCommandScope{ChatID: ch.Id, Type: "chat"})); err != nil {
			log.Println("Unable to UNset commands")
		} else {
			log.Println("UNset commands")
		}

		HelloScreen(ch.Id).Exec(b)
		update = <-ch.Updates
		if update.Message != nil {
			if b.CheckForRoot(update) {
				fmt.Println("Set id for root: ", update.Message.Text)
				b.SetIds(update)
				rootCommands.Scope.ChatID = ch.Id
				if _, err := b.BotApi.Request(rootCommands); err != nil {
					log.Println("Unable to set root commands")
				} else {
					log.Println("Def commands set")
				}
				b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "You were promoted to root"))
				go Admin(b, ch)
				return
			} else if b.CheckForUser(update) {
				b.SetIds(update)
				b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "You are trainee"))
				go Trainee(b, ch)
				return
			}

		}
	}
	commands.Scope.ChatID = ch.Id
	if _, err := b.BotApi.Request(commands); err != nil {
		log.Println("Unable to set def commands")
	} else {
		log.Println("Def commands set")
	}
	b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "You are not in trainee list check for typo or write to Alexander"))
	go NTrainee(b, ch)
}

func BasicLoop(b *cmp.Bot, ch *cmp.Chat) {

}

func ExecCommand(b *cmp.Bot, ch *cmp.Chat, upd tgbotapi.Update) error {
	switch upd.Message.Command() {
	case "superuser":
		b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "Enter your name again"))
		upd = <-ch.Updates
		if !b.CheckForRoot(upd) {
			return fmt.Errorf("You are not root")
		}
		rootCommands.Scope.ChatID = ch.Id
		if _, err := b.BotApi.Request(rootCommands); err != nil {
			log.Println("Unable to set root commands")
		} else {
			log.Println("Def commands set")
		}
		b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "You were promoted to root"))
		go Admin(b, ch)
		return nil
	case "start":
		b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "you have already started"))
	case "register":
		b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "Enter your name again"))
		upd = <-ch.Updates
		if b.CheckForUser(upd) {
			if _, err := b.BotApi.Request(tgbotapi.NewDeleteMyCommandsWithScope(tgbotapi.BotCommandScope{ChatID: ch.Id, Type: "chat"})); err != nil {
				log.Println("Unable to UNset commands")
			} else {
				log.Println("UNset commands")
			}
			b.SetIds(upd)
			go Trainee(b, ch)
			return nil
		}
	}
	return fmt.Errorf("no such command")
}

func ExecRootCommand(b *cmp.Bot, ch *cmp.Chat, upd tgbotapi.Update) error {
	//if !b.CheckForRoot(upd) {
	//	return fmt.Errorf("You are not root")
	//}
	switch upd.Message.Command() {
	case "show_all_users":
		b.ShowAllUsersTo(ch.Id)
	case "notify_l0":
		b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "Enter your message"))
		upd = <-ch.Updates
		msg := upd.Message.Text
		b.Notify(`SELECT chat_id from users where l0_status = 'делаю'`, msg)
	case "notify_l1":
		b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "Enter your message"))
		upd = <-ch.Updates
		msg := upd.Message.Text
		b.Notify(`SELECT chat_id from users where l1_status = 'делаю'`, msg)
	case "notify_l2":
		b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "Enter your message"))
		upd = <-ch.Updates
		msg := upd.Message.Text
		b.Notify(`SELECT chat_id from users where l2_status = 'делаю'`, msg)
	}
	return fmt.Errorf("no such command")
}

func Admin(b *cmp.Bot, ch *cmp.Chat) {

	rootCommands.Scope.ChatID = ch.Id
	b.BotApi.Request(rootCommands)
	for update := range ch.Updates {
		if update.Message.IsCommand() {
			ExecRootCommand(b, ch, update)
		}
	}
}

func Trainee(b *cmp.Bot, ch *cmp.Chat) {
	sendl0(b, ch)
	for update := range ch.Updates {
		if update.Message.IsCommand() {
			if ExecCommand(b, ch, update) == nil {
				return
			}
		}
		if update.CallbackQuery != nil && update.CallbackQuery.Data == "done_l0" {
			b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "Enter your git link with l0"))
			b.Db.Exec(`UPDATE public.users
			SET l0_link=$1, l0_status='сделано', l0_end_date=$2`,
				(<-ch.Updates).Message.Text,
				time.Now(),
			)
		}
	}

}

func sendl0(b *cmp.Bot, ch *cmp.Chat) {
	var name string
	b.Db.QueryRow(`SELECT link from tasks where l_num = 0`).Scan(&name)
	b.Db.Exec(`UPDATE public.users
	SET l0_status=$1
	WHERE chat_id=$2;`,
		"делаю",
		ch.Id)
	b.BotApi.Send(tgbotapi.NewMessage(ch.Id, "Вот ссылка на l0"))
	b.BotApi.Send(tgbotapi.NewMessage(ch.Id, name))
	msg := tgbotapi.NewMessage(ch.Id, "Как будешь готов тыкай Готово!")
	km := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Готово!", "done_l0"),
		),
	)
	msg.ReplyMarkup = km
	b.BotApi.Send(msg)
}

func NTrainee(b *cmp.Bot, ch *cmp.Chat) {
	for update := range ch.Updates {
		if update.Message.IsCommand() {
			if ExecCommand(b, ch, update) == nil {
				return
			}
		}
	}
}
