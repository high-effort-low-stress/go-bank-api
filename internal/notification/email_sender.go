package notification

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/resend/resend-go/v2"
)

type EmailRequest struct {
	From         string
	To           string
	Subject      string
	TemplatePath string
	TemplateData any
}

type EmailService interface {
	SendEmail(request *EmailRequest) error
}

type resendEmailService struct {
	client *resend.Client
}

func NewEmailService() (EmailService, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEY environment variable not set")
	}

	return &resendEmailService{
		client: resend.NewClient(apiKey),
	}, nil
}

func (s *resendEmailService) SendEmail(request *EmailRequest) error {
	body, err := s.buildEmailContent(request.TemplatePath, request.TemplateData)

	params := &resend.SendEmailRequest{
		From:    request.From,
		To:      []string{request.To},
		Subject: request.Subject,
		Html:    body,
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Verification email sent successfully to %s. Message ID: %s", params.To, sent.Id)
	return nil
}

func (s *resendEmailService) buildEmailContent(templatePath string, templateData any) (string, error) {
	mailTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	err = mailTemplate.Execute(&body, templateData)
	if err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return body.String(), nil
}
