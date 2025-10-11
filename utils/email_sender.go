package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendMail(from, to, subject, body string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: RESEND_API_KEY environment variable not set.")
	}

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	sendMail := os.Getenv("SEND_MAIL")
	if sendMail != "true" {
		return nil
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("Error sending email: %s", err)
	}
	fmt.Println("Email sent successfully! Message ID:", sent.Id)

	return nil
}
