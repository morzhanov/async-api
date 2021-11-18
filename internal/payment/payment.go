package payment

import (
	"context"

	"github.com/jmoiron/sqlx"
	apipayment "github.com/morzhanov/async-api/api/payment"
	uuid "github.com/satori/go.uuid"
)

type pay struct {
	db *sqlx.DB
}

type Payment interface {
	ProcessPayment(ctx context.Context, in *apipayment.ProcessPaymentMessage) error
}

func (p *pay) ProcessPayment(ctx context.Context, in *apipayment.ProcessPaymentMessage) error {
	if _, err := p.db.QueryContext(
		ctx,
		`INSERT INTO payments (id, order_id, name, amount, status) VALUES ($id, $orderId, $name, $amount, $status)`,
		uuid.NewV4().String(), in.OrderID, in.Name, in.Amount, in.Status,
	); err != nil {
		return err
	}
	return nil
}

func NewPayment(db *sqlx.DB) Payment {
	return &pay{db}
}
