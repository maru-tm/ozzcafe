package service

import (
	"fmt"
	"log"
	"ozzcafe/server/models"

	"gopkg.in/gomail.v2"
)

// EmailService отправляет email пользователю
type EmailService struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort int
}

// NewEmailService создает новый EmailService
func NewEmailService() *EmailService {
	// Используем данные из кода вашего друга для отправки email
	return &EmailService{
		From:     "elibrarysender@gmail.com", // Почта отправителя
		Password: "ocxwblzcockfwcud",         // Пароль или App password
		SMTPHost: "smtp.gmail.com",           // SMTP хост
		SMTPPort: 587,                        // SMTP порт
	}
}

// SendVerificationEmail отправляет email с токеном для подтверждения
func (s *EmailService) SendVerificationEmail(user *models.User) error {
	// Генерация токена для подтверждения email
	token, err := generateVerificationToken(user.Email)
	if err != nil {
		log.Println("Error generating token:", err)
		return err
	}

	// Ссылка для подтверждения email
	verificationURL := fmt.Sprintf("http://localhost:8080/verify?token=%s&email=%s", token, user.Email)

	// Тема и тело письма
	subject := "Please confirm your email address"
	body := fmt.Sprintf("Click the link to verify your email: %s", verificationURL)

	// Отправка письма
	err = s.sendEmail(user.Email, subject, body)
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	log.Printf("Verification email sent to: %s\n", user.Email)
	return nil
}

// sendEmail отправляет письмо с заданной темой и телом
func (s *EmailService) sendEmail(to, subject, body string) error {
	from := s.From // Используем email из структуры

	// Создание нового сообщения
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Создание соединения с сервером SMTP
	d := gomail.NewDialer(s.SMTPHost, s.SMTPPort, from, s.Password)

	// Отправка сообщения
	err := d.DialAndSend(m)
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}
