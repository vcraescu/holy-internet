package main

import (
	"github.com/spf13/viper"
	"github.com/vcraescu/holy-internet/internal/app/holyinternet"
	"log"
	"time"
	"github.com/vcraescu/holy-internet/pkg/utils"
	"fmt"
)


var (
	startingTime time.Time
)

func init() {
	holyinternet.ReadConfig()
}

func main() {
	mailer := holyinternet.NewMailer(
		viper.GetString("email.host"),
		viper.GetString("email.username"),
		viper.GetString("email.password"),
		viper.GetString("email.port"),
	)

	md := holyinternet.NewMailerDaemon(mailer)
	md.Start()

	go forever(md)
	select{}
}

func startTimer() {
	startingTime = time.Now()
}

func stopTimer() (time.Duration) {
	return time.Duration(time.Since(startingTime))
}

func forever(md *holyinternet.MailerDaemon) {
	saints := viper.GetStringSlice("saints")

	down := false

	for {
		ok := holyinternet.IsInternetOK(saints)
		log.Printf("Internet: %v", ok)
		if !ok {
			if !down {
				startTimer()
				utils.NotifyCritical("Error", "Sorry! No internet for you")
			}
			down = true
		} else {
			if down {
				down = false
				d := stopTimer()
				utils.NotifyCritical("Yay!", fmt.Sprintf("Internet was down for %s seconds", d))
			}
		}

		time.Sleep(time.Second * 5)
		msg := holyinternet.Message{To: "viorel@wingravity.com", Body: "This is a test"}
		if err := md.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
