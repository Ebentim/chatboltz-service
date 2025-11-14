package smtp

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	textTemplate "text/template"
)

// EmailService handles email operations with template support
type EmailService struct {
	defaultClient *Client
	templates     map[string]*EmailTemplate
}

// EmailTemplate holds both HTML and text templates
type EmailTemplate struct {
	Subject  string
	HTMLTmpl *template.Template
	TextTmpl *textTemplate.Template
}

// EmailData represents data passed to email templates
type EmailData struct {
	Code string
}

// NewEmailService creates a new email service with default SMTP config
func NewEmailService(defaultConfig Config) (*EmailService, error) {
	service := &EmailService{
		defaultClient: NewClient(defaultConfig),
		templates:     make(map[string]*EmailTemplate),
	}

	if err := service.loadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load email templates: %w", err)
	}

	return service, nil
}

// SendOTP sends OTP email using templates
func (s *EmailService) SendOTP(to, purpose, code string) error {
	return s.SendOTPWithClient(s.defaultClient, to, purpose, code)
}

// SendOTPWithClient sends OTP email using custom SMTP client
func (s *EmailService) SendOTPWithClient(client interface{}, to, purpose, code string) error {
	smtpClient, ok := client.(*Client)
	if !ok {
		return fmt.Errorf("invalid client type")
	}
	tmpl, exists := s.templates[purpose]
	if !exists {
		return fmt.Errorf("template not found for purpose: %s", purpose)
	}

	data := EmailData{Code: code}

	// Render HTML content
	var htmlBuf bytes.Buffer
	if err := tmpl.HTMLTmpl.Execute(&htmlBuf, data); err != nil {
		return fmt.Errorf("failed to render HTML template: %w", err)
	}

	// Render text content
	var textBuf bytes.Buffer
	if err := tmpl.TextTmpl.Execute(&textBuf, data); err != nil {
		return fmt.Errorf("failed to render text template: %w", err)
	}

	return smtpClient.SendHTML(to, tmpl.Subject, htmlBuf.String(), textBuf.String())
}

// CreateClientFromConfig creates a new SMTP client from user config
func (s *EmailService) CreateClientFromConfig(config interface{}) interface{} {
	smtpConfig, ok := config.(Config)
	if !ok {
		return nil
	}
	return NewClient(smtpConfig)
}

// loadTemplates loads all OTP email templates
func (s *EmailService) loadTemplates() error {
	purposes := []string{"2fa", "password_reset", "login"}
	subjects := map[string]string{
		"2fa":            "Two-Factor Authentication Code",
		"password_reset": "Password Reset Code",
		"login":          "Login Verification Code",
	}

	for _, purpose := range purposes {
		htmlPath := filepath.Join("internal", "templates", "email", "otp", purpose+".html")
		textPath := filepath.Join("internal", "templates", "email", "otp", purpose+".txt")

		htmlTmpl, err := template.ParseFiles(htmlPath)
		if err != nil {
			return fmt.Errorf("failed to parse HTML template %s: %w", htmlPath, err)
		}

		textTmpl, err := textTemplate.ParseFiles(textPath)
		if err != nil {
			return fmt.Errorf("failed to parse text template %s: %w", textPath, err)
		}

		s.templates[purpose] = &EmailTemplate{
			Subject:  subjects[purpose],
			HTMLTmpl: htmlTmpl,
			TextTmpl: textTmpl,
		}
	}

	return nil
}
