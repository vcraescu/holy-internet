package main

import (
	"github.com/spf13/viper"
	"github.com/vcraescu/holy-internet/internal/app/holyinternet"
	"log"
	"time"
	"github.com/vcraescu/holy-internet/pkg/utils"
	"fmt"
	"errors"
	"math/rand"
)

const (
	noInternetCheckInterval = time.Second * 5
)

var (
	startingTime time.Time
)

func init() {
	holyinternet.ReadConfig()
}

func main() {
	md := startMailerDaemon()
	app := holyinternet.NewApp(md)
	app.Run(onReady, onExit)
}

func onReady(app *holyinternet.App) {
	saints := viper.GetStringSlice("saints")
	down := false
	defaultSleepDuration := time.Second * time.Duration(viper.GetInt("pray.every"))
	sleepDuration := defaultSleepDuration

	log.Printf("Checking The Holy Internet connection every %s...", sleepDuration)

	for {
		time.Sleep(sleepDuration)
		if app.IsPaused() {
			continue
		}

		ok := holyinternet.IsInternetOK(saints)
		log.Printf("Internet: %v", ok)

		if !ok {
			if !down {
				startTimer()
				utils.NotifyCritical("Error", "Sorry! No internet for you")
			}
			down = true
			sleepDuration = noInternetCheckInterval
			app.Failed()
			continue
		}

		if down {
			app.DiscardFailure()
			down = false
			d := stopTimer()
			utils.NotifyNormal("Yay!", fmt.Sprintf("Internet was down for %s seconds", d))
			if err, _ := sendEmailToFollowers(app.MailerDaemon, d); err != nil {
				log.Println(err)
			}

			if _, err := sendCurse(app.MailerDaemon); err != nil {
				log.Println(err)
			}
		}

		sleepDuration = defaultSleepDuration
	}
}

func onExit(app *holyinternet.App) {
	log.Println("Exit")
}

func startMailerDaemon() (*holyinternet.MailerDaemon) {
	mailer := holyinternet.NewMailer(
		viper.GetString("email.host"),
		viper.GetString("email.username"),
		viper.GetString("email.password"),
		viper.GetString("email.port"),
		viper.GetBool("email.tls"),
	)

	md := holyinternet.NewMailerDaemon(mailer)
	md.Start()

	return md
}

func startTimer() {
	startingTime = time.Now()
}

func stopTimer() (time.Duration) {
	return time.Duration(time.Since(startingTime))
}

func sendEmailToFollowers(md *holyinternet.MailerDaemon, d time.Duration) (error, int) {
	emails := viper.GetStringSlice("followers.people")
	if len(emails) == 0 {
		return errors.New("no followers found"), 0
	}

	count := 0
	for _, email := range emails {
		msg := holyinternet.Message{
			Subject: "Internet was down",
			From:    viper.GetString("email.username"),
			Body:    fmt.Sprintf(viper.GetString("followers.message"), d),
			To:      email,
		}

		if err := md.Send(msg); err != nil {
			log.Println(err)
		}

		count++
	}

	return nil, count
}

func sendCurse(md *holyinternet.MailerDaemon) (bool, error) {
	messages := viper.GetStringSlice("curses.messages")
	if len(messages) == 0 {
		return false, errors.New("no curse messages found")
	}

	target := viper.GetString("curses.target")
	if target == "" {
		return false, errors.New("curse target not defined")
	}

	index := rand.Intn(len(messages))
	msg := holyinternet.Message{
		Subject: "Internet was down",
		From:    viper.GetString("email.username"),
		Body:    messages[index],
		To:      target,
	}

	if err := md.Send(msg); err != nil {
		return false, err
	}

	return true, nil
}
