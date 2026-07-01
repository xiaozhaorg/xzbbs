package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/xiaozhaorg/xzbbs/internal/config"
)

// EmailService sends emails via SMTP
type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

// Send sends an email using SMTP
func (s *EmailService) Send(to, subject, htmlBody string) error {
	if !config.Global.Email.Enabled {
		return fmt.Errorf("email not configured")
	}

	cfg := config.Global.Email
	
	// Build MIME message
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s\r\n", cfg.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
	buf.WriteString(htmlBody)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return smtp.SendMail(addr, auth, cfg.From, []string{to}, buf.Bytes())
}

// SendVerificationEmail sends an email verification link
func (s *EmailService) SendVerificationEmail(to, token string) error {
	siteURL := "http://localhost:8080"
	
	link := fmt.Sprintf("%s/email/verify?token=%s", siteURL, token)
	
	tmpl := `<html><body>
		<p>请点击以下链接验证您的邮箱：</p>
		<p><a href="{{.Link}}">{{.Link}}</a></p>
		<p>如果不是您本人操作，请忽略此邮件。</p>
	</body></html>`
	
	t, _ := template.New("email").Parse(tmpl)
	var body bytes.Buffer
	t.Execute(&body, map[string]string{"Link": link})
	
	return s.Send(to, "请验证您的邮箱", body.String())
}

// SendPasswordResetEmail sends a password reset link
func (s *EmailService) SendPasswordResetEmail(to, token string) error {
	siteURL := "http://localhost:8080"
	
	link := fmt.Sprintf("%s/password/reset?token=%s", siteURL, token)
	
	tmpl := `<html><body>
		<p>请点击以下链接重置您的密码：</p>
		<p><a href="{{.Link}}">{{.Link}}</a></p>
		<p>该链接将在 30 分钟后失效。如果不是您本人操作，请忽略此邮件。</p>
	</body></html>`
	
	t, _ := template.New("email").Parse(tmpl)
	var body bytes.Buffer
	t.Execute(&body, map[string]string{"Link": link})
	
	return s.Send(to, "重置您的密码", body.String())
}
