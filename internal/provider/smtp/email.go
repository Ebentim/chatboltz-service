package smtp

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Client struct {
	host string
	port string
	user string
	pass string
	from string
}

type Config struct {
	Host string
	Port string
	User string
	Pass string
	From string // optional override; if empty will use User
}

func NewClient(cfg Config) *Client {
	from := cfg.From
	if strings.TrimSpace(from) == "" {
		from = cfg.User
	}
	return &Client{host: cfg.Host, port: cfg.Port, user: cfg.User, pass: cfg.Pass, from: from}
}

// Send sends a plain text email
func (c *Client) Send(to, subject, body string) error {
	return c.sendEmail(to, subject, body, "text/plain")
}

// SendHTML sends an HTML email with text fallback
func (c *Client) SendHTML(to, subject, htmlBody, textBody string) error {
	body := c.buildMultipartBody(htmlBody, textBody)
	return c.sendEmail(to, subject, body, "multipart/alternative")
}

// sendEmail sends email with specified content type
func (c *Client) sendEmail(to, subject, body, contentType string) error {
	addr := fmt.Sprintf("%s:%s", c.host, c.port)
	auth := smtp.PlainAuth("", c.user, c.pass, c.host)
	headers := map[string]string{
		"From":         c.from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": contentType + "; charset=utf-8",
	}
	var sb strings.Builder
	for k, v := range headers {
		sb.WriteString(k + ": " + v + "\r\n")
	}
	sb.WriteString("\r\n" + body)

	return smtp.SendMail(addr, auth, c.from, []string{to}, []byte(sb.String()))
}

// buildMultipartBody creates multipart email body with HTML and text
func (c *Client) buildMultipartBody(htmlBody, textBody string) string {
	boundary := "boundary123456789"
	body := fmt.Sprintf(`--boundary123456789
Content-Type: text/plain; charset=utf-8

%s

--boundary123456789
Content-Type: text/html; charset=utf-8

%s

--boundary123456789--`, textBody, htmlBody)
	return strings.ReplaceAll(body, "boundary123456789", boundary)
}

// SendOTP sends OTP email (legacy method for backward compatibility)
func (c *Client) SendOTP(to, purpose, code string) error {
	subject, body := c.buildOTPEmail(purpose, code)
	return c.Send(to, subject, body)
}

// buildOTPEmail creates basic email content (legacy)
func (c *Client) buildOTPEmail(purpose, code string) (string, string) {
	switch purpose {
	case "2fa":
		return "Two-Factor Authentication Code",
			fmt.Sprintf(`Your two-factor authentication code is:

%s

This code will expire in 10 minutes.
If you didn't request this code, please ignore this email.`, code)
	case "password_reset":
		return "Password Reset Code",
			fmt.Sprintf(`Your password reset code is:

%s

This code will expire in 10 minutes.
If you didn't request a password reset, please ignore this email.`, code)
	case "login":
		return "Login Verification Code",
			fmt.Sprintf(`Your login verification code is:

%s

This code will expire in 10 minutes.
If you didn't attempt to log in, please secure your account immediately.`, code)
	default:
		return "Verification Code",
			fmt.Sprintf(`Your verification code is:

%s

This code will expire in 10 minutes.`, code)
	}
}
