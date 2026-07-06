package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"fit24/domain"
)

var (
	ErrInvalidName  = errors.New("имя должно содержать не менее 2-х символов")
	ErrInvalidPhone = errors.New("некорректный формат телефона")
	ErrInvalidEmail = errors.New("некорректный формат почты")
	ErrSpamDetected = errors.New("сообщение заблокировано спам-фильтром")
)

type leadService struct {
	repo domain.LeadRepository
}

func NewLeadService(repo domain.LeadRepository) domain.LeadService {
	return &leadService{repo: repo}
}

func (s *leadService) SubmitOrder(ctx context.Context, plan, name, phone, email string) (*domain.Order, error) {
	name = strings.TrimSpace(name)
	phone = strings.TrimSpace(phone)
	email = strings.TrimSpace(email)

	if len(name) < 2 {
		return nil, ErrInvalidName
	}

	cleanedPhone, err := s.normalizePhone(phone)
	if err != nil {
		return nil, ErrInvalidPhone
	}

	if !s.isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	order := &domain.Order{
		Plan:      plan,
		Name:      name,
		Phone:     cleanedPhone,
		Email:     email,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := s.repo.SaveOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("ошибка сохранения заказа: %w", err)
	}

	return order, nil
}

func (s *leadService) SubmitContact(ctx context.Context, name, phone, message string) (*domain.Contact, error) {
	name = strings.TrimSpace(name)
	phone = strings.TrimSpace(phone)
	message = strings.TrimSpace(message)

	if len(name) < 2 {
		return nil, ErrInvalidName
	}

	cleanedPhone, err := s.normalizePhone(phone)
	if err != nil {
		return nil, ErrInvalidPhone
	}

	if s.isSpam(message) {
		return nil, ErrSpamDetected
	}

	contact := &domain.Contact{
		Name:      name,
		Phone:     cleanedPhone,
		Message:   message,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := s.repo.SaveContact(ctx, contact); err != nil {
		return nil, fmt.Errorf("ошибка сохранения обращения: %w", err)
	}

	return contact, nil
}

func (s *leadService) normalizePhone(phone string) (string, error) {
	reg := regexp.MustCompile(`[\s\-\(\)]`)
	cleaned := reg.ReplaceAllString(phone, "")

	phoneRegex := regexp.MustCompile(`^(?:\+7|7|8)?[0-9]{10}$`)
	if !phoneRegex.MatchString(cleaned) {
		return "", errors.New("invalid phone structure")
	}

	if strings.HasPrefix(cleaned, "8") {
		cleaned = "+7" + cleaned[1:]
	} else if strings.HasPrefix(cleaned, "7") {
		cleaned = "+7" + cleaned[1:]
	} else if !strings.HasPrefix(cleaned, "+7") {
		cleaned = "+7" + cleaned
	}

	return cleaned, nil
}

func (s *leadService) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	return emailRegex.MatchString(email)
}

func (s *leadService) isSpam(message string) bool {
	lowerMsg := strings.ToLower(message)
	spamWords := []string{"casino", "crypto", "крипта", "казино", "заработок", "free money", "розыгрыш"}
	for _, word := range spamWords {
		if strings.Contains(lowerMsg, word) {
			return true
		}
	}
	return false
}
