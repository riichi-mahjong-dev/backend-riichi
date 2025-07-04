package utils

import (
	"log"
	"strconv"

	"github.com/go-gomail/gomail"
	"github.com/riichi-mahjong-dev/backend/configs"
)

type Emailer struct {
	Host        string
	Port        int
	FromAddress string
	FromName    string
	Username    string
	Password    string
}

func InitializeEmailer(emailConfig *configs.EmailConfig) (*Emailer, error) {
	portNumber, err := strconv.Atoi(emailConfig.EmailPort)

	if err != nil {
		return nil, err
	}

	return &Emailer{
		Host:        emailConfig.EmailHost,
		Port:        portNumber,
		FromAddress: emailConfig.FromAddress,
		FromName:    emailConfig.FromName,
		Username:    emailConfig.UserName,
		Password:    emailConfig.Password,
	}, nil
}

func (emailer *Emailer) SendEmail(to string, message *gomail.Message) error {
	message.SetHeader("From", emailer.FromAddress)
	message.SetHeader("To", to)
	dial := gomail.NewDialer(emailer.Host, emailer.Port, emailer.Username, emailer.Password)

	if err := dial.DialAndSend(message); err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	return nil
}
