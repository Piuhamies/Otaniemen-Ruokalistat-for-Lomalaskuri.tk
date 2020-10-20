package telegram

import (
	"fmt"
	"time"

	"otaniemenruokalistat.tk/ruokalista"

	"github.com/robfig/cron/v3"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Bot is the telegram bot that sends daily messages to a channel
var Bot *tb.Bot

const token = "your token here"

// Init starts the bot
func Init() {
	bot, err := tb.NewBot(tb.Settings{
		URL:    "https://api.telegram.org",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Token:  token,
	})
	if err != nil {
		panic(err)
	}

	bot.Handle(tb.OnText, func(m *tb.Message) {
		fmt.Println(m.Chat.ID)
	})

	bot.Handle("/ruoka", func(message *tb.Message) {
		weekday := int(time.Now().Weekday()) - 1
		if weekday > 4 {
			weekday = 0
		}
		ruokalista, _ := ruokalista.GetThisWeeksFood()
		paivanRuoka := ruokalista[weekday]
		ruoka := fmt.Sprintf("%s\nKotiruoka: %s\nKasvisruoka: %s",
			paivanRuoka.Viikonpäivä, paivanRuoka.Perus, paivanRuoka.Veg)
		bot.Send(message.Sender, ruoka)
	})
	bot.Handle("/lista", func(message *tb.Message) {
		ruokalista, _ := ruokalista.GetThisWeeksFood()
		ruoka := ""
		for _, d := range ruokalista {
			ruoka += fmt.Sprintf("%s\nKotiruoka: %s\nKasvisruoka: %s\n\n",
				d.Viikonpäivä, d.Perus, d.Veg)
		}
		bot.Send(message.Sender, ruoka)
	})
	chat, err := bot.ChatByID("-1001490826318")
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting cron job for bot to send updates...")
	c := cron.New()
	c.AddFunc("0 7 * * *", func() { sendUpdate(bot, chat) })
	c.Start()
	fmt.Println("Job scheduled")
	fmt.Println("Starting bot")
	bot.Start()
}

func sendUpdate(bot *tb.Bot, chat *tb.Chat) {
	weekday := int(time.Now().Weekday()) - 1
	if weekday > 4 {
		return
	}
	ruokalista, _ := ruokalista.GetThisWeeksFood()
	paivanRuoka := ruokalista[weekday]
	ruoka := fmt.Sprintf("%s\nKotiruoka: %s\nKasvisruoka: %s",
		paivanRuoka.Viikonpäivä, paivanRuoka.Perus, paivanRuoka.Veg)
	bot.Send(chat, ruoka)
}
