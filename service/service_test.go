package service

import (
	"context"
	"errors"
	"testing"

	"fit24/domain"
)

type mockRepository struct {
	saveOrderFunc   func(ctx context.Context, order *domain.Order) error
	saveContactFunc func(ctx context.Context, contact *domain.Contact) error
}

func (m *mockRepository) SaveOrder(ctx context.Context, order *domain.Order) error {
	if m.saveOrderFunc != nil {
		return m.saveOrderFunc(ctx, order)
	}
	return nil
}

func (m *mockRepository) SaveContact(ctx context.Context, contact *domain.Contact) error {
	if m.saveContactFunc != nil {
		return m.saveContactFunc(ctx, contact)
	}
	return nil
}

func TestSubmitContact_SpamFilter(t *testing.T) {
	repo := &mockRepository{}
	svc := NewLeadService(repo)

	_, err := svc.SubmitContact(context.Background(), "Иван", "+79998887766", "заходи играть в казино")
	if !errors.Is(err, ErrSpamDetected) {
		t.Errorf("expected ErrSpamDetected, got: %v", err)
	}
}

func TestSubmitContact_InvalidName(t *testing.T) {
	repo := &mockRepository{}
	svc := NewLeadService(repo)

	_, err := svc.SubmitContact(context.Background(), "А", "+79998887766", "привет")
	if !errors.Is(err, ErrInvalidName) {
		t.Errorf("expected ErrInvalidName, got: %v", err)
	}
}
