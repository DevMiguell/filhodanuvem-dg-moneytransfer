package money

import (
	"github.com/filhodanuvem/dg-moneytransfer/internal/errors"
	"github.com/google/uuid"
)

type TransferService struct {
	Repository Repository
}

func (s *TransferService) Transfer(amount int, debtorID, beneficiaryID uuid.UUID) error {
	err, cancel := s.Repository.OpenTransaction()
	defer cancel()
	if err != nil {
		return err
	}

	err = s.removeFromBalance(amount, debtorID)
	if err != nil {
		return err
	}

	err = s.topUpBalance(amount, beneficiaryID)
	if err != nil {
		s.Repository.Rollback()
		return err
	}

	s.Repository.Commit()
	return nil
}

func (s *TransferService) removeFromBalance(amount int, userID uuid.UUID) error {
	debtorBalance, err := s.Repository.SelectBalanceByUserID(userID)
	if err != nil {
		return errors.New(errors.CodeInternalDatabaseError, "error on selecting balance", err)
	}

	if debtorBalance.Amount-amount < 0 {
		return errors.New(errors.CodeInsufficientBalance, "insufficient balance on debtor account", err)
	}

	err = s.Repository.RemoveFromBalanceByUserID(amount, userID)

	return err
}

func (s *TransferService) topUpBalance(amount int, userID uuid.UUID) error {
	_, err := s.Repository.SelectBalanceByUserID(userID)
	if err != nil {
		return errors.New(errors.CodeInternalDatabaseError, "error on selecting balance", err)
	}

	err = s.Repository.AddOnBalanceByUserID(amount, userID)

	return err
}
