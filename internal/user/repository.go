package user

import (
	"context"
	"time"

	moneytransfer "github.com/filhodanuvem/dg-moneytransfer"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

const defaultTimeout = 10 * time.Second

type Repository struct {
	Conn *pgxpool.Pool
}

func (repo *Repository) SelectBalanceByUserID(userID uuid.UUID) (moneytransfer.Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	sql := "SELECT id, amount, user_id FROM balances WHERE user_id = $1"

	row := repo.Conn.QueryRow(ctx, sql, userID)

	var balance moneytransfer.Balance
	if err := row.Scan(
		&balance.ID,
		&balance.Amount,
		&balance.UserID,
	); err != nil {
		return moneytransfer.Balance{}, err
	}

	return balance, nil
}
