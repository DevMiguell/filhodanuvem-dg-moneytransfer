package money

import (
	"testing"

	"github.com/filhodanuvem/dg-moneytransfer/internal/errors"
	"github.com/google/uuid"
)

func TestTransferSucceeded(t *testing.T) {
	repo := MockRepository{}
	service := TransferService{
		Repository: &repo,
	}

	user1 := uuid.New()
	user2 := uuid.New()
	amount := 10
	repo.allDatabaseOperationsWorked()

	err := service.Transfer(amount, user1, user2)
	if err != nil {
		t.Error(err)
	}

	if err := repo.check(); err != nil {
		t.Error(err)
	}
}

func TestDebtorWithInsufficientBalance(t *testing.T) {
	repo := MockRepository{}
	service := TransferService{
		Repository: &repo,
	}

	user1 := uuid.New()
	user2 := uuid.New()
	amount := 10
	repo.insufficientDebtorBalanceExpectsNoDeposit()

	err := service.Transfer(amount, user1, user2)
	if err == nil {
		t.Error("expected error of insufficient balance")
	}

	e, ok := err.(errors.Error)
	if !ok || e.Code != errors.CodeInsufficientBalance {
		t.Error("expected error of insufficient balance")
	}

	if err := repo.check(); err != nil {
		t.Error(err)
	}
}

func TestRollbackWhenFailsOnDeposit(t *testing.T) {
	repo := MockRepository{}
	service := TransferService{
		Repository: &repo,
	}

	user1 := uuid.New()
	user2 := uuid.New()
	amount := 10
	repo.failToTopUpRollbackTransaction()

	err := service.Transfer(amount, user1, user2)
	if err == nil {
		t.Error("expected error on sql transaction")
	}

	if err := repo.check(); err != nil {
		t.Error(err)
	}
}
