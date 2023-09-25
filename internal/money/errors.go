package money

import (
	"log"
	"net/http"

	"github.com/filhodanuvem/dg-moneytransfer/internal/errors"
)

var errCodeInvalidAmountToTransfer = errors.New(
	errors.CodeInvalidAmountToTransfer,
	"Invalid amount to transfer, expected greater than zero",
	nil,
)

var errCodeSameDebtorAndBeneficiary = errors.New(
	errors.CodeSameDebtorAndBeneficiary,
	"Cannot transfer to its own account",
	nil,
)

var errCodeMissingPart = errors.New(
	errors.CodeMissingPart,
	"Missing beneficiary or debtor",
	nil,
)

func responseFromError(err error, w http.ResponseWriter) {
	e, ok := err.(errors.Error)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	switch e.Code {
	case errors.CodeInsufficientBalance:
		w.WriteHeader(http.StatusExpectationFailed)
	case errors.CodeSameDebtorAndBeneficiary:
		w.WriteHeader(http.StatusExpectationFailed)
	case errors.CodeInvalidAmountToTransfer:
		w.WriteHeader(http.StatusExpectationFailed)
	case errors.CodeMissingPart:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(e.Message))
}
