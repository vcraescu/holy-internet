package holyinternet

import (
	"errors"
	"log"
	"time"
)

type MailerDaemon struct {
	mailer  *Mailer
	ch      chan Message
	started bool
}

func NewMailerDaemon(mailer *Mailer) (*MailerDaemon) {
	return &MailerDaemon{
		mailer: mailer,
	}
}

func (md *MailerDaemon) Start() (error) {
	if md.started {
		return errors.New("already started")
	}

	md.started = true
	md.ch = make(chan Message)

	go func() {
		defer func() {
			md.mailer.Close()
		}()

		opened := false

		for {
			select {
			case msg := <-md.ch:
				if !opened {
					if err := md.mailer.Open(); err != nil {
						log.Printf("Mailer Daemon (Open): %v", err)
						continue
					}

					opened = true
				}
				if err := md.mailer.Send(msg); err != nil {
					log.Printf("Mailer Daemon (Send): %v", err)
				}

				log.Printf("Mailer Daemon (Send): Email sent to %v", msg.To)
			case <-time.After(time.Second * 30):
				if opened {
					md.mailer.Close()
					opened = false
				}
			}
		}
	}()

	return nil
}

func (md *MailerDaemon) Send(msg Message) (error) {
	if !md.started {
		return errors.New("mailer daemon not started")
	}

	go func() {
		md.ch <- msg
	}()

	return nil
}
