package smtp

import (
	"crypto/tls"
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

// Send sends a plain text email via STARTTLS if supported.
func (c *Client) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", c.host, c.port)
	auth := smtp.PlainAuth("", c.user, c.pass, c.host)
	headers := map[string]string{
		"From":         c.from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=utf-8",
	}
	var sb strings.Builder
	for k, v := range headers {
		sb.WriteString(k + ": " + v + "\r\n")
	}
	sb.WriteString("\r\n" + body)

	// Try a TLS connection first
	tlsconfig := &tls.Config{ServerName: c.host, InsecureSkipVerify: false}
	// Use smtp.SendMail (will initiate STARTTLS internally if server supports?)
	// For explicit TLS we would implement custom dial, simplified here.
	return smtp.SendMail(addr, auth, c.from, []string{to}, []byte(sb.String()))
}

// SendOTP convenience helper
func (c *Client) SendOTP(to, purpose, code string) error {
	subject := fmt.Sprintf("Your %s OTP Code", strings.Title(strings.ReplaceAll(purpose, "_", " ")))
	body := fmt.Sprintf("Your %s verification code is: %s\nIt expires shortly. If you did not request this, please ignore.", purpose, code)
	return c.Send(to, subject, body)
}
