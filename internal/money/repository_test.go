package money

import (
	"context"
	"fmt"

	moneytransfer "github.com/filhodanuvem/dg-moneytransfer"
	"github.com/google/uuid"
)

type MockRepository struct {
	insufficientDebtor bool
	failToTopUp        bool

	expectedSelectQueryCounter int
	expectedRemoveQueryCounter int
	expectedInsertQueryCounter int
	expectedRollbackCounter    int
	expectedCommitCounter      int

	currentSelectQueryCounter int
	currentRemoveQueryCounter int
	currentInsertQueryCounter int
	currentRollbackCounter    int
	currentCommitCounter      int
}

func (repo *MockRepository) OpenTransaction() (error, context.CancelFunc) {
	cancel := func() {}
	return nil, cancel
}

func (repo *MockRepository) Commit() error {
	repo.currentCommitCounter++
	return nil
}

func (repo *MockRepository) Rollback() error {
	repo.currentRollbackCounter++
	return nil
}

func (repo *MockRepository) SelectBalanceByUserID(userID uuid.UUID) (moneytransfer.Balance, error) {
	repo.currentSelectQueryCounter++
	amount := 999999
	if repo.insufficientDebtor {
		amount = 0
	}
	return moneytransfer.Balance{Amount: amount}, nil
}

func (repo *MockRepository) RemoveFromBalanceByUserID(amount int, userID uuid.UUID) error {
	repo.currentRemoveQueryCounter++
	if repo.insufficientDebtor {
		return fmt.Errorf("insufficient balance on debtor account %s", userID)
	}
	return nil
}

func (repo *MockRepository) AddOnBalanceByUserID(amount int, userID uuid.UUID) error {
	repo.currentInsertQueryCounter++
	if repo.failToTopUp {
		return fmt.Errorf("internal database transaction insert on account %s", userID)
	}
	return nil
}

// mock methods
func (repo *MockRepository) allDatabaseOperationsWorked() {
	repo.expectedInsertQueryCounter = 1
	repo.expectedRemoveQueryCounter = 1
	repo.expectedSelectQueryCounter = 2
	repo.expectedCommitCounter = 1
}

func (repo *MockRepository) insufficientDebtorBalanceExpectsNoDeposit() {
	repo.insufficientDebtor = true
	repo.expectedSelectQueryCounter = 1
}

func (repo *MockRepository) failToTopUpRollbackTransaction() {
	repo.failToTopUp = true
	repo.expectedInsertQueryCounter = 1
	repo.expectedRemoveQueryCounter = 1
	repo.expectedSelectQueryCounter = 2
	repo.expectedRollbackCounter = 1
}

func (repo *MockRepository) check() error {
	if repo.expectedInsertQueryCounter != repo.currentInsertQueryCounter {
		return fmt.Errorf("expected %d insert queries, found %d", repo.expectedInsertQueryCounter, repo.currentInsertQueryCounter)
	}

	if repo.expectedRemoveQueryCounter != repo.currentRemoveQueryCounter {
		return fmt.Errorf("expected %d remove queries, found %d", repo.expectedRemoveQueryCounter, repo.currentRemoveQueryCounter)
	}

	if repo.expectedSelectQueryCounter != repo.currentSelectQueryCounter {
		return fmt.Errorf("expected %d select queries, found %d", repo.expectedSelectQueryCounter, repo.currentSelectQueryCounter)
	}

	if repo.expectedRollbackCounter != repo.currentRollbackCounter {
		return fmt.Errorf("expected %d rollback operation, found %d", repo.expectedRollbackCounter, repo.currentRollbackCounter)
	}

	if repo.expectedRollbackCounter != repo.currentRollbackCounter {
		return fmt.Errorf("expected %d rollback operation, found %d", repo.expectedRollbackCounter, repo.currentRollbackCounter)
	}

	if repo.expectedCommitCounter != repo.currentCommitCounter {
		return fmt.Errorf("expected %d commit operation, found %d", repo.expectedCommitCounter, repo.currentCommitCounter)
	}

	return nil
}
