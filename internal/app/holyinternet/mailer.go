package holyinternet

import (
	"net/smtp"
	"fmt"
	"net"
	"crypto/tls"
	"log"
	"errors"
)

type Mailer struct {
	host     string
	username string
	password string
	port     string
	tls      bool
	client   *smtp.Client
}

type Message struct {
	To      string
	From    string
	Subject string
	Body    string
}

func (msg *Message) GetTextBody() ([]byte) {
	body := fmt.Sprintf("From: %s\r\n", msg.From) +
		fmt.Sprintf("To: %s\r\n", msg.To) +
		fmt.Sprintf("Subject: %s\r\n\r\n", msg.Subject) +
		msg.Body + "\r\n"

	return []byte(body)
}

func NewMailer(host, username, password, port string, tls bool) (*Mailer) {
	return &Mailer{
		host:     host,
		username: username,
		password: password,
		port:     port,
		tls:      tls,
	}
}

func (m *Mailer) Open() (error) {
	addr := net.JoinHostPort(m.host, m.port)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.host,
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return err
	}

	m.client, err = smtp.NewClient(conn, m.host)
	if err != nil {
		return nil
	}

	if m.tls {
		if err := m.client.StartTLS(tlsconfig); err != nil {
			return err
		}
	}

	auth := smtp.PlainAuth(
		"",
		m.username,
		m.password,
		m.host,
	)

	return m.client.Auth(auth)
}

func (m *Mailer) Close() (error) {
	if m.client == nil {
		return nil
	}

	return m.client.Quit()
}

func (m *Mailer) Send(msg Message) (error) {
	if m.client == nil {
		return errors.New("mail client not initialized")
	}

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
