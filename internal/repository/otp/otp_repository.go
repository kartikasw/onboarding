package otp

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"onboarding/pkg/config"
	"text/template"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	emailTemplate = `
					<!DOCTYPE html>
					<html>
					<body style="font-family: Helvetica, Arial; padding: 24px; color: #333;">
					  <h2>{{.Service}}</h2>
					  <p>To continue, use the verification code below:</p>
					  <div style="
					    margin: 24px 0;
					    padding: 16px 24px;
					    font-size: 32px;
					    font-weight: 700;
					    letter-spacing: 6px;
					    border: 1px solid #ccc;
					    border-radius: 8px;
					    background: #fafafa;
					  ">
					    {{.OTP}}
					  </div>
					  <p>This code expires in 5 minutes. Please do not share it.</p>
					</body>
					</html>`
)

type OtpRepository interface {
	SendOtp(ctx context.Context, email string, service ServiceType) error
	VerifyOtp(ctx context.Context, email string, otp string, service ServiceType) error
}

type IOtpRepository struct {
	redis   *redis.Client
	smtpCfg config.SMTP
}

func NewOtpRepository(redis *redis.Client, smtpCfg config.SMTP) OtpRepository {
	return &IOtpRepository{redis: redis, smtpCfg: smtpCfg}
}

func (i *IOtpRepository) SendOtp(ctx context.Context, email string, service ServiceType) error {
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	var body bytes.Buffer
	tmpl, err := template.New("otp").Parse(emailTemplate)
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}

	if err := tmpl.Execute(&body, map[string]string{
		"OTP":     otp,
		"Service": service.Name,
	}); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	if err := sendEmailSMTP(i.smtpCfg, email, service.Name, body.String()); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	key := fmt.Sprintf("otp:%s:%s", service.Code, email)
	err = i.redis.Set(ctx, key, otp, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("redis error: %w", err)
	}

	return nil
}

func (i *IOtpRepository) VerifyOtp(ctx context.Context, email string, otp string, service ServiceType) error {
	key := fmt.Sprintf("otp:%s:%s", service.Code, email)

	stored, err := i.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return errors.New("invalid or expired otp")
	}
	if err != nil {
		return fmt.Errorf("redis error: %w", err)
	}

	if stored != otp {
		return errors.New("invalid otp")
	}

	i.redis.Del(ctx, key)

	return nil
}

func sendEmailSMTP(cfg config.SMTP, to string, service string, htmlBody string) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         cfg.Host,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial error: %w", err)
	}

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return fmt.Errorf("smtp client error: %w", err)
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth error: %w", err)
	}

	if err := client.Mail(cfg.Username); err != nil {
		return err
	}

	if err := client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	// Build MIME email
	var msg bytes.Buffer
	msg.WriteString("From: " + cfg.FromName + "\r\n")
	msg.WriteString("To: " + to + "\r\n")
	msg.WriteString("Subject: " + service + " Verification Code\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	if _, err := w.Write(msg.Bytes()); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	client.Quit()
	return nil
}
