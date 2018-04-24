package gotool


import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)


type (
	Email struct {
	From        string
	To          []string
	Cc          []string
	Subject     string
	ContentType string
	Content     string
	Attach      string
	}

	EmailClient struct {
	Host     string
	Port     int
	Username string
	Password string
	Message  *Email
	}

)

func NewEmail(from, subject, contentType, content, attach string, to, cc []string) *Email {
	return &Email{
		From:        from,
		Subject:     subject,
		ContentType: contentType,
		Content:     content,
		To:          to,
		Cc:          cc,
		Attach:      attach,
	}
}


func NewEmailClient(host, username, password string, port int, message *Email) *EmailClient {
	return &EmailClient{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Message:  message,
	}
}

func (c *EmailClient) SendMessage() (bool, error) {

	e := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	if 587 == c.Port {
		e.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	mail := gomail.NewMessage()
	mail.SetHeader("From", c.Message.From)
	mail.SetHeader("To", c.Message.To...)

	if len(c.Message.Cc) != 0 {
		mail.SetHeader("Cc", c.Message.Cc...)
	}

	mail.SetHeader("Subject", c.Message.Subject)
	mail.SetBody(c.Message.ContentType, c.Message.Content)

	if c.Message.Attach != "" {
		mail.Attach(c.Message.Attach)
	}

	if err := e.DialAndSend(mail); err != nil {
		return false, err
	}
	return true, nil
}