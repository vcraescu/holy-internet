package holyinternet

import (
	"net/smtp"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"time"
	"fmt"
)

type Mailer struct {
	host string
	username string
	password string
	port string
	client *smtp.Client
}

type Message struct {
	To string
	From string
	Subject string
	Body string
}

type MailerDaemon struct {
	mailer *Mailer
	ch chan Message
	started bool
}

func (msg *Message) GetTextBody() ([]byte) {
	body := fmt.Sprintf("From: %s\r\n", msg.From) +
		fmt.Sprintf("To: %s\r\n", msg.To) +
		fmt.Sprintf("Subject: %s\r\n\r\n", msg.Subject) +
		msg.Body + "\r\n"

	return []byte(body)
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
						panic(err)
					}

					opened = true
				}
				if err := md.mailer.Send(msg); err != nil {
					log.Println(err)
				}
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
		return errors.New("not started")
	}

	go func() {
		md.ch <- msg
	}()

	return nil
}

func NewMailer(host, username, password, port string) (*Mailer) {
	return &Mailer{
		host: host,
		username: username,
		password: password,
		port: port,
	}
}

func (m Mailer) Open() (error) {
	addr := net.JoinHostPort(m.host, m.port)

	tlsconfig := &tls.Config {
		InsecureSkipVerify: true,
		ServerName: m.host,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	log.Println(addr)
	if err != nil {
		return err
	}

	m.client, err = smtp.NewClient(conn, m.host)
	if err != nil {
		return nil
	}

	auth := smtp.PlainAuth(
		"",
		m.username,
		m.password,
		m.host,
	)

	return m.client.Auth(auth)
}

func (m Mailer) Close() (error) {
	if m.client == nil {
		return nil
	}

	return m.client.Quit()
}

func (m Mailer) Send(msg Message) (error) {
	log.Println("send")
    if err := m.client.Mail(m.username); err != nil {
    	return err
    }

    if err := m.client.Rcpt(msg.To); err != nil {
    	return err
    }

    w, err := m.client.Data()
    if err != nil {
    	return err
    }

    _, err = w.Write(msg.GetTextBody())
    if err != nil {
    	return err
    }

    return w.Close()
}
