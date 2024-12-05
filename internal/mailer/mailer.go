package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/wneessen/go-mail"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	client *mail.Client
	sender string
}

func New(host string, port int, username, password, sender string) (Mailer, error) {
	// client, err := mail.NewClient(host, mail.WithPort(port), mail.WithUsername(username), mail.WithPassword(password), mail.WithTimeout(5*time.Second))
	client, err := mail.NewClient(host, mail.WithPort(port), mail.WithTimeout(5*time.Second), mail.WithTLSPortPolicy(mail.NoTLS)) // TEMPORARY
	if err != nil {
		return Mailer{}, err
	}

	return Mailer{
		client: client,
		sender: sender,
	}, nil
}

func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()

	// careful here cuz IgnoreInvalid versions of funcs is being used
	err = msg.To(recipient)
	if err != nil {
		return err
	}
	// msg.SetAddrHeaderIgnoreInvalid(mail.HeaderTo, recipient)
	err = msg.From(m.sender)
	if err != nil {
		return err
	}
	// msg.SetAddrHeaderIgnoreInvalid(mail.HeaderFrom, m.sender)

	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())

	for i := 1; i <= 3; i++ {
		err = m.client.DialAndSend(msg)

		// Attention here!
		if nil == err {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return err
}
