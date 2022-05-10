package notify

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"caas-monitoring-batch/util"
	"log"
	"os"
)

var config map[string]string

func init() {
	var err error
	config, err = util.ReadConfig("config.ini")
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

func SendChatBot(receivers []int64, message string, token string) {
	fmt.Println("token : " + token)
	fmt.Printf(" sns len : %v\n", len(receivers))
	if len(receivers) > 0 {
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			fmt.Println("Failed to get telegram client connection! :", err)
			return
		}
		bot.Debug = true

		for _, receiver := range receivers {
			msg := tgbotapi.NewMessage(receiver, message)
			botMsg, botErr := bot.Send(msg)
			fmt.Printf(">>>>> botMsg=[%v], botErr[%v]\n", botMsg, botErr)
		}
	}
}
