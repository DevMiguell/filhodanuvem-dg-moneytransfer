package money

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/filhodanuvem/dg-moneytransfer/internal/database"
	"github.com/filhodanuvem/dg-moneytransfer/internal/user"
	"github.com/google/uuid"
)

type transferRequest struct {
	DebtorID      string `json:"debtor_id"`
	BeneficiaryID string `json:"beneficiary_id"`
	Amount        int    `json:"amount"`
}

func UsersHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	conn, err := database.CreateConnection()
	if err != nil {
		responseFromError(err, w)
		return
	}
	defer conn.Close()

	id := strings.TrimPrefix(req.URL.Path, "/users/")
	userID, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	repo := user.Repository{
		Conn: conn,
	}

	balance, err := repo.SelectBalanceByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(balance)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(j)
}

func TransferHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	conn, err := database.CreateConnection()
	if err != nil {
		responseFromError(err, w)
		return
	}
	defer conn.Close()

	var tr transferRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&tr); err != nil {
		responseFromError(err, w)
		return
	}

	if err := validateTransfer(tr); err != nil {
		responseFromError(err, w)
		return
	}

	repo := PostgresRepository{Conn: conn}
	transferService := TransferService{Repository: &repo}
	debtorID, err := uuid.Parse(tr.DebtorID)
	if err != nil {
		responseFromError(err, w)
		return
	}
	beneficiaryID, err := uuid.Parse(tr.BeneficiaryID)
	if err != nil {
		responseFromError(err, w)
		return
	}
	err = transferService.Transfer(
		tr.Amount,
		debtorID,
		beneficiaryID,
	)

	if err != nil {
		responseFromError(err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
