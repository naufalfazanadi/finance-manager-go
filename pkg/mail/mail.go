package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/go-mail/mail"

	"github.com/naufalfazanadi/finance-manager-go/pkg/config"
)

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// EmailData represents email content
type EmailData struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

// EmailTemplateData represents data for email templates
type EmailTemplateData struct {
	Name     string
	ResetURL string
}

// getEmailConfig creates email configuration from config first, then env as fallback
func getEmailConfig() *EmailConfig {
	cfg := config.GetConfig()
	smtpCfg := cfg.SMTP

	// If config is not set, fallback to env (shouldn't happen, but for safety)
	return &EmailConfig{
		SMTPHost:     smtpCfg.Host,
		SMTPPort:     smtpCfg.Port,
		SMTPUsername: smtpCfg.Username,
		SMTPPassword: smtpCfg.Password,
		FromEmail:    smtpCfg.FromEmail,
		FromName:     smtpCfg.FromName,
	}
}

// SendEmail sends an email using go-mail/mail
func SendEmail(emailData *EmailData) error {
	config := getEmailConfig()

	if config.SMTPUsername == "" || config.SMTPPassword == "" {
		return fmt.Errorf("SMTP credentials not configured")
	}

	if config.FromEmail == "" {
		return fmt.Errorf("from email not configured")
	}

	if len(emailData.To) == 0 {
		return fmt.Errorf("recipient email required")
	}

	// Create new message
	m := mail.NewMessage()

	// Set headers
	m.SetHeader("From", m.FormatAddress(config.FromEmail, config.FromName))
	m.SetHeader("To", emailData.To...)
	m.SetHeader("Subject", emailData.Subject)

	// Set body
	if emailData.IsHTML {
		m.SetBody("text/html", emailData.Body)
	} else {
		m.SetBody("text/plain", emailData.Body)
	}

	// Create dialer
	d := mail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendEmailWithTemplate sends an email with custom HTML template
func SendEmailWithTemplate(to, subject, htmlBody string) error {
	emailData := &EmailData{
		To:      []string{to},
		Subject: subject,
		Body:    htmlBody,
		IsHTML:  true,
	}

	return SendEmail(emailData)
}

// getTemplatesDir returns the templates directory path
func getTemplatesDir() string {
	// Get the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		// Fallback to relative path
		workingDir = "."
	}

	return filepath.Join(workingDir, "assets", "templates", "email")
}

// LoadTemplate loads and parses an email template with the provided data
func LoadTemplate(templateName string, data interface{}) (string, error) {
	templatesDir := getTemplatesDir()
	templatePath := filepath.Join(templatesDir, templateName)

	// Check if template file exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template file not found: %s", templatePath)
	}

	// Parse the template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// getEnv is no longer needed; config is used for env access
