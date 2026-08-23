package transaction

import (
	"context"

	"gorm.io/gorm"
)

type TransactionManager interface {
	ExecTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type TransactionManagerImp struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &TransactionManagerImp{db: db}
}

func (t *TransactionManagerImp) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "TxKey", tx)
		return fn(txCtx)
	})
}
