package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"testbot/cmp"
	"testbot/src"

	_ "github.com/lib/pq"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var baseMsg = map[string]string{
	"start": `Это тестовый бот, который пока что ничего не умеет, но скоро всему научится`,
}

var Chats map[int64]*cmp.Chat = make(map[int64]*cmp.Chat)

func index_handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		f, _ := os.ReadFile("./static/example.html")
		//
		//

		//fmt.Println("insides", e, f)
		fmt.Fprintf(w, string(f))
	case "POST":
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		q, err := url.ParseQuery(bodyString)
		if err != nil {
			//panic(err)
		}
		log.Println(q.Get("first_name"))
	}
}

//aaaaaaaaaaaaaaaaa

func main() {

	http.HandleFunc("/static", index_handler)

	DBconnStr := "host=" + "nbshtech.ru" + " user=" + "wb_order" + " dbname=" + "postgres" + " password=" + "mypass" + " sslmode=disable"

	Bot := cmp.Bot{}
	var err error
	Bot.Db, err = sql.Open("postgres", DBconnStr)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI("5392077863:AAHrIAh8uKMQqeFnD2fAFu7hh-50rjacZa8")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	//tgbotapi.NewWebhookWithCert("", tgbotapi.FilePath("cert.pem"))
	wh, _ := tgbotapi.NewWebhookWithCert("https://nbshtech.ru:8443/"+bot.Token, tgbotapi.FilePath("YOURPUBLIC.pem"))
	//wh.IPAddress = "185.20.226.116"
	//wh, _ := tgbotapi.NewWebhook("https://185.20.226.116:8443/" + bot.Token)
	//tgbotapi.NewWebhookWithCert("185.20.226.116:8443/", "cert.pem")

	bot.Request(wh)
	//_, err = bot.SetWebhook(wh)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println("===listenwebhook===")

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS(":8443", "YOURPUBLIC.pem", "YOURPRIVATE.key", nil)
	//go http.ListenAndServe("0.0.0.0:8443", nil)

	//fmt.Println("===webhookinfo===", wh.URL)
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	Bot.BotApi = bot
	Bot.Msg = &tgbotapi.MessageConfig{}

	for update := range updates {

		if update.CallbackQuery != nil {
			Bot.BotApi.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Отлично!"))
		}
		if update.Message != nil {
			if v, ok := Chats[update.Message.Chat.ID]; ok {
				v.Updates <- update
			} else {
				fmt.Println("starting new chat")
				v := &cmp.Chat{
					Updates: make(chan tgbotapi.Update),
					Id:      update.Message.Chat.ID}
				Chats[v.Id] = v
				go src.StartChatWithUser(&Bot, v)
				v.Updates <- update
			}

			//text := cmp.Text{
			//	Msg:    update.Message.Text,
			//	ChatId: update.Message.Chat.ID,
			//}
			//Bot.Execute(text)
		}
	}
}
