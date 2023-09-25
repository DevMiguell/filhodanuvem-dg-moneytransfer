package money

func validateTransfer(t transferRequest) error {
	if t.Amount <= 0 {
		return errCodeInvalidAmountToTransfer
	}

	if t.BeneficiaryID == t.DebtorID {
		return errCodeSameDebtorAndBeneficiary
	}

	if t.BeneficiaryID == "" || t.DebtorID == "" {
		return errCodeMissingPart
	}

	return nil
}
