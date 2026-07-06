package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"fit24/domain"
)

var (
	ErrInvalidName  = errors.New("имя должно содержать не менее 2-х символов")
	ErrInvalidPhone = errors.New("некорректный формат телефона")
	ErrInvalidEmail = errors.New("некорректный формат Email")
	ErrSpamDetected = errors.New("сообщение заблокировано спам-фильтром")
)

type leadService struct {
	repo domain.LeadRepository
}

func NewLeadService(repo domain.LeadRepository) domain.LeadService {
	return &leadService{repo: repo}
}

func (s *leadService) SubmitOrder(ctx context.Context, plan, name, phone, email string) (*domain.Order, error) {
	slog.Info("начало обработки заявки",
		slog.String("plan", plan),
		slog.String("email", email),
	)

	name = strings.TrimSpace(name)
	phone = strings.TrimSpace(phone)
	email = strings.TrimSpace(email)

	if utf8.RuneCountInString(name) < 2 {
		slog.Warn("валидация имени не пройдена",
			slog.String("name", name),
		)
		return nil, ErrInvalidName
	}

	cleanedPhone, err := s.normalizePhone(phone)
	if err != nil {
		slog.Warn("валидация телефона не пройдена",
			slog.String("raw_phone", phone),
			slog.String("error", err.Error()),
		)
		return nil, ErrInvalidPhone
	}

	if !s.isValidEmail(email) {
		slog.Warn("валидация почты не пройдена",
			slog.String("email", email),
		)
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
		slog.Error("не удалось сохранить заявку в БД",
			slog.String("error", err.Error()),
			slog.String("phone", cleanedPhone),
		)
		return nil, fmt.Errorf("ошибка сохранения заявки: %w", err)
	}

	slog.Info("заявка успешно сохранена",
		slog.String("plan", plan),
		slog.String("phone", cleanedPhone),
	)

	return order, nil
}

func (s *leadService) SubmitContact(ctx context.Context, name, phone, message string) (*domain.Contact, error) {
	slog.Info("начало обработки обращения",
		slog.String("name", name),
	)

	name = strings.TrimSpace(name)
	phone = strings.TrimSpace(phone)
	message = strings.TrimSpace(message)

	if utf8.RuneCountInString(name) < 2 {
		slog.Warn("валидация имени в обращении не пройдена",
			slog.String("name", name),
		)
		return nil, ErrInvalidName
	}

	cleanedPhone, err := s.normalizePhone(phone)
	if err != nil {
		slog.Warn("валидация телефона в обращении не пройдена",
			slog.String("raw_phone", phone),
			slog.String("error", err.Error()),
		)
		return nil, ErrInvalidPhone
	}

	if s.isSpam(message) {
		slog.Warn("обращение заблокировано спам-фильтром",
			slog.String("name", name),
			slog.String("phone", cleanedPhone),
		)
		return nil, ErrSpamDetected
	}

	contact := &domain.Contact{
		Name:      name,
		Phone:     cleanedPhone,
		Message:   message,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := s.repo.SaveContact(ctx, contact); err != nil {
		slog.Error("не удалось сохранить обращение в БД",
			slog.String("error", err.Error()),
			slog.String("phone", cleanedPhone),
		)
		return nil, fmt.Errorf("ошибка сохранения обращения: %w", err)
	}

	slog.Info("обращение успешно сохранено",
		slog.String("phone", cleanedPhone),
	)

	return contact, nil
}

func (s *leadService) normalizePhone(phone string) (string, error) {
	reg := regexp.MustCompile(`[\s\-\(\)]`)
	cleaned := reg.ReplaceAllString(phone, "")
	slog.Debug("символы форматирования удалены", slog.String("cleaned", cleaned))

	phoneRegex := regexp.MustCompile(`^(?:\+7|7|8)?[0-9]{10}$`)
	if !phoneRegex.MatchString(cleaned) {
		slog.Debug("формат телефона не соответствует шаблону", slog.String("phone", cleaned))
		return "", errors.New("invalid phone structure")
	}

	if strings.HasPrefix(cleaned, "8") {
		cleaned = "+7" + cleaned[1:]
	} else if strings.HasPrefix(cleaned, "7") {
		cleaned = "+7" + cleaned[1:]
	} else if !strings.HasPrefix(cleaned, "+7") {
		cleaned = "+7" + cleaned
	}

	slog.Debug("нормализация телефона успешна", slog.String("result", cleaned))
	return cleaned, nil
}

func (s *leadService) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	valid := emailRegex.MatchString(email)
	slog.Debug("результат проверки email", slog.String("email", email), slog.Bool("valid", valid))
	return valid
}

func (s *leadService) isSpam(message string) bool {
	lowerMsg := strings.ToLower(message)
	spamWords := []string{"casino", "crypto", "крипта", "казино", "заработок", "free money", "розыгрыш"}
	for _, word := range spamWords {
		if strings.Contains(lowerMsg, word) {
			slog.Debug("обнаружено спам слово", slog.String("word", word), slog.String("message", message))
			return true
		}
	}
	return false
}
