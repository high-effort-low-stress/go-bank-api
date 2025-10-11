package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/resend/resend-go/v2"
)

type EmailService interface {
	SendVerificationEmail(fullName, to, token string) error
}

type resendEmailService struct {
	client                  *resend.Client
	from                    string
	frontendVerificationURL string
	verificationTmpl        *template.Template
}

func NewEmailService() (EmailService, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEY environment variable not set")
	}

	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		return nil, fmt.Errorf("EMAIL_FROM environment variable not set")
	}

	frontendURL := os.Getenv("FRONTEND_VERIFICATION_URL")
	if frontendURL == "" {
		return nil, fmt.Errorf("FRONTEND_VERIFICATION_URL environment variable not set")
	}

	// Carrega o template de e-mail na inicialização
	tmpl, err := template.ParseFiles(os.Getenv("VERIFICATION_EMAIL_TEMPLATE"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse email template: %w", err)
	}

	return &resendEmailService{
		client:                  resend.NewClient(apiKey),
		from:                    from,
		frontendVerificationURL: frontendURL,
		verificationTmpl:        tmpl,
	}, nil
}

func (s *resendEmailService) SendVerificationEmail(fullName, to, token string) error {
	subject := "Bem-vindo ao GoBank! Confirme seu e-mail."
	verificationLink := fmt.Sprintf("%s?token=%s", s.frontendVerificationURL, token)

	// Prepara os dados para o template
	templateData := struct {
		FullName         string
		Subject          string
		VerificationLink string
	}{
		FullName:         fullName,
		Subject:          subject,
		VerificationLink: verificationLink,
	}

	var body bytes.Buffer
	if err := s.verificationTmpl.Execute(&body, templateData); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
		Html:    body.String(),
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Verification email sent successfully to %s. Message ID: %s", to, sent.Id)
	return nil
}
