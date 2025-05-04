package mail

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	valueobjects "github.com/VanLavr/Diploma-fin/internal/domain/value_objects"
	"github.com/VanLavr/Diploma-fin/utils/config"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type mailer struct {
	SMTPHost     string
	SMTPPort     string
	AuthEmail    string
	AuthPassword string
	OAuthCode    string
}

func NewStudentMailer(cfg *config.Config) repositories.StudentMailer {
	return &mailer{
		SMTPHost:     cfg.SMTPHost,
		SMTPPort:     cfg.SMTPPort,
		AuthEmail:    cfg.AuthEmail,
		AuthPassword: cfg.AuthEmailPassword,
		OAuthCode:    cfg.SMTP2OAuthCode,
	}
}

// SendPassword implements repositories.StudentMailer.
func (this *mailer) SendPassword(ctx context.Context, studentEmail, password string) error {
	if studentEmail == "" || !strings.Contains(studentEmail, "@") {
		log.Logger.Error(errors.ErrInvalidData.Error(), errors.MethodKey, log.GetMethodName())
		return fmt.Errorf("correct student email is required")
	}
	fmt.Println("20", this.OAuthCode)

	subject := fmt.Sprintf("Subject: New user account for: %s\n", studentEmail)
	body := fmt.Sprintf(
		valueobjects.NewPasswordNotification,
		studentEmail,
		password,
	)

	msg := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth(
		"",
		this.AuthEmail,
		this.OAuthCode,
		this.SMTPHost,
	)

	if err := smtp.SendMail(
		this.SMTPHost+":"+this.SMTPPort,
		auth,
		this.AuthEmail,
		[]string{studentEmail},
		msg,
	); err != nil {
		fmt.Println("21")
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully to", studentEmail)
	return nil
}

func (this mailer) SendNotification(ctx context.Context, student models.Student, teacherEmail string, exam models.Exam) error {
	if teacherEmail == "" {
		return fmt.Errorf("teacher email is required")
	}
	fmt.Println("20", this.OAuthCode)

	subject := fmt.Sprintf("Subject: Exam Notification: %s\n", exam.Name)
	body := fmt.Sprintf(
		valueobjects.EmailText,
		student.FirstName,
		student.LastName,
		student.MiddleName,
		student.Group.Name,
		exam.Name,
	)

	msg := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth(
		"",
		this.AuthEmail,
		this.OAuthCode,
		this.SMTPHost,
	)

	if err := smtp.SendMail(
		this.SMTPHost+":"+this.SMTPPort,
		auth,
		this.AuthEmail,
		[]string{teacherEmail},
		msg,
	); err != nil {
		fmt.Println("21")
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully to", teacherEmail)
	return nil
}
