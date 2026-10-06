package services

import (
	"car-dealership-system/internal/config"
	"fmt"

	"gopkg.in/gomail.v2"
)

// EmailService handles email sending operations
type EmailService struct {
	config config.EmailConfig
}

// NewEmailService creates a new email service
func NewEmailService(cfg config.EmailConfig) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// SendVerificationCode sends a verification code to the user's email
func (s *EmailService) SendVerificationCode(email, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.FromEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Код подтверждения - Car Dealership System")

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2c3e50; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f8f9fa; padding: 30px; border-radius: 5px; margin-top: 20px; }
        .code { font-size: 32px; font-weight: bold; color: #3498db; text-align: center; 
                letter-spacing: 5px; padding: 20px; background-color: white; 
                border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
        .warning { color: #e74c3c; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚗 Car Dealership System</h1>
        </div>
        <div class="content">
            <h2>Код подтверждения для входа</h2>
            <p>Здравствуйте!</p>
            <p>Вы запросили вход в систему Car Dealership. Используйте следующий код для завершения аутентификации:</p>
            
            <div class="code">%s</div>
            
            <p><span class="warning">⚠️ Важно:</span> Этот код действителен в течение 10 минут.</p>
            <p>Если вы не запрашивали этот код, просто проигнорируйте это письмо.</p>
        </div>
        <div class="footer">
            <p>Это автоматическое письмо, пожалуйста, не отвечайте на него.</p>
            <p>&copy; 2025 Car Dealership System. Все права защищены.</p>
        </div>
    </div>
</body>
</html>
	`, code)

	m.SetBody("text/html", body)

	// Alternative plain text version
	plainText := fmt.Sprintf(`
Car Dealership System - Код подтверждения

Здравствуйте!

Вы запросили вход в систему Car Dealership. 
Используйте следующий код для завершения аутентификации:

%s

⚠️ Важно: Этот код действителен в течение 10 минут.

Если вы не запрашивали этот код, просто проигнорируйте это письмо.

---
Это автоматическое письмо, пожалуйста, не отвечайте на него.
© 2025 Car Dealership System. Все права защищены.
	`, code)

	m.AddAlternative("text/plain", plainText)

	// Create dialer
	d := gomail.NewDialer(
		s.config.SMTPHost,
		s.config.SMTPPort,
		s.config.FromEmail,
		s.config.AppPassword,
	)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendWelcomeEmail sends a welcome email to new users
func (s *EmailService) SendWelcomeEmail(email, firstName string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.FromEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Добро пожаловать в Car Dealership System!")

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2c3e50; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f8f9fa; padding: 30px; border-radius: 5px; margin-top: 20px; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
        .button { display: inline-block; padding: 12px 24px; background-color: #3498db; 
                  color: white; text-decoration: none; border-radius: 5px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚗 Car Dealership System</h1>
        </div>
        <div class="content">
            <h2>Добро пожаловать, %s!</h2>
            <p>Спасибо за регистрацию в Car Dealership System.</p>
            <p>Теперь вы можете:</p>
            <ul>
                <li>Просматривать каталог автомобилей</li>
                <li>Добавлять автомобили в корзину</li>
                <li>Оформлять заказы</li>
                <li>Отслеживать историю покупок</li>
            </ul>
            <p>Мы рады видеть вас в нашей системе!</p>
        </div>
        <div class="footer">
            <p>&copy; 2025 Car Dealership System. Все права защищены.</p>
        </div>
    </div>
</body>
</html>
	`, firstName)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		s.config.SMTPHost,
		s.config.SMTPPort,
		s.config.FromEmail,
		s.config.AppPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
}

// SendOrderConfirmation sends order confirmation email
func (s *EmailService) SendOrderConfirmation(email, firstName string, orderID uint, totalAmount float64) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.FromEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", fmt.Sprintf("Подтверждение заказа #%d", orderID))

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2c3e50; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f8f9fa; padding: 30px; border-radius: 5px; margin-top: 20px; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
        .order-info { background-color: white; padding: 20px; border-radius: 5px; margin: 20px 0; }
        .total { font-size: 24px; font-weight: bold; color: #27ae60; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚗 Car Dealership System</h1>
        </div>
        <div class="content">
            <h2>Спасибо за ваш заказ, %s!</h2>
            <p>Ваш заказ успешно оформлен и принят в обработку.</p>
            
            <div class="order-info">
                <p><strong>Номер заказа:</strong> #%d</p>
                <p><strong>Общая сумма:</strong> <span class="total">$%.2f</span></p>
                <p><strong>Статус:</strong> В обработке</p>
            </div>
            
            <p>Мы свяжемся с вами в ближайшее время для уточнения деталей.</p>
            <p>Вы можете отслеживать статус заказа в личном кабинете.</p>
        </div>
        <div class="footer">
            <p>&copy; 2025 Car Dealership System. Все права защищены.</p>
        </div>
    </div>
</body>
</html>
	`, firstName, orderID, totalAmount)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		s.config.SMTPHost,
		s.config.SMTPPort,
		s.config.FromEmail,
		s.config.AppPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send order confirmation: %w", err)
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(to, firstName, token string) error {
	resetURL := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", token)

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.FromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Сброс пароля - Car Dealership System")

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2c3e50; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f8f9fa; padding: 30px; border-radius: 5px; margin-top: 20px; }
        .button { display: inline-block; padding: 15px 30px; background-color: #3498db; 
                  color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
        .warning { color: #e74c3c; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚗 Car Dealership System</h1>
        </div>
        <div class="content">
            <h2>Сброс пароля</h2>
            <p>Здравствуйте, %s!</p>
            <p>Вы запросили сброс пароля для вашего аккаунта в Car Dealership System.</p>
            <p>Нажмите на кнопку ниже, чтобы сбросить пароль:</p>
            
            <div style="text-align: center;">
                <a href="%s" class="button">Сбросить пароль</a>
            </div>
            
            <p>Или скопируйте эту ссылку в браузер:</p>
            <p style="word-break: break-all; color: #3498db;">%s</p>
            
            <p><span class="warning">⚠️ Важно:</span> Эта ссылка действительна в течение 1 часа.</p>
            <p>Если вы не запрашивали сброс пароля, просто проигнорируйте это письмо.</p>
        </div>
        <div class="footer">
            <p>Это автоматическое письмо, пожалуйста, не отвечайте на него.</p>
            <p>&copy; 2025 Car Dealership System. Все права защищены.</p>
        </div>
    </div>
</body>
</html>
	`, firstName, resetURL, resetURL)

	m.SetBody("text/html", body)

	// Alternative plain text version
	plainText := fmt.Sprintf(`
Car Dealership System - Сброс пароля

Здравствуйте, %s!

Вы запросили сброс пароля для вашего аккаунта.

Перейдите по ссылке для сброса пароля:
%s

Эта ссылка действительна в течение 1 часа.

Если вы не запрашивали сброс пароля, просто проигнорируйте это письмо.

С уважением,
Команда Car Dealership System
	`, firstName, resetURL)

	m.AddAlternative("text/plain", plainText)

	d := gomail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.FromEmail, s.config.AppPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}
