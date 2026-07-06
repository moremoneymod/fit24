package domain

import "context"

type Order struct {
	CreatedAt string `json:"created_at"`
	Plan      string `json:"plan"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type Contact struct {
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Message   string `json:"message"`
}

type LeadRepository interface {
	SaveOrder(ctx context.Context, order *Order) error
	SaveContact(ctx context.Context, contact *Contact) error
}

type LeadService interface {
	SubmitOrder(ctx context.Context, plan, name, phone, email string) (*Order, error)
	SubmitContact(ctx context.Context, name, phone, message string) (*Contact, error)
}
