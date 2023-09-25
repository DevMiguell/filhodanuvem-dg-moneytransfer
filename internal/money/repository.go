package money

import (
	"context"
	"errors"
	"time"

	moneytransfer "github.com/filhodanuvem/dg-moneytransfer"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const transactionTimeout = 20 * time.Second
const defaultTimeout = 10 * time.Second

type Repository interface {
	OpenTransaction() (err error, cancelContext context.CancelFunc)
	Commit() error
	Rollback() error
	SelectBalanceByUserID(userID uuid.UUID) (moneytransfer.Balance, error)
	RemoveFromBalanceByUserID(amount int, userID uuid.UUID) error
	AddOnBalanceByUserID(amount int, userID uuid.UUID) error
}

type PostgresRepository struct {
	Conn *pgxpool.Pool
	tx   pgx.Tx
	ctx  context.Context
}

func (repo *PostgresRepository) OpenTransaction() (err error, cancelContext context.CancelFunc) {
	ctx, cancelContext := context.WithTimeout(context.Background(), transactionTimeout)

	repo.ctx = ctx
	options := pgx.TxOptions{}
	repo.tx, err = repo.Conn.BeginTx(ctx, options)

	return
}

func (repo *PostgresRepository) Commit() error {
	if repo.tx.Conn().IsClosed() {
		return errors.New("database transaction is closed already")
	}

	return repo.tx.Commit(repo.ctx)
}

func (repo *PostgresRepository) Rollback() error {
	if repo.tx.Conn().IsClosed() {
		return errors.New("database transaction is closed already")
	}

	return repo.tx.Rollback(repo.ctx)
}

func (repo *PostgresRepository) SelectBalanceByUserID(userID uuid.UUID) (moneytransfer.Balance, error) {
	ctx, cancel := context.WithTimeout(repo.ctx, defaultTimeout)
	defer cancel()

	sql := "SELECT id, amount, user_id FROM balances WHERE user_id = $1 FOR UPDATE "

	row := repo.tx.QueryRow(ctx, sql, userID)

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

func (repo *PostgresRepository) RemoveFromBalanceByUserID(amount int, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(repo.ctx, defaultTimeout)
	defer cancel()

	sql := "UPDATE balances SET amount = amount - $1 WHERE user_id = $2"

	_, err := repo.tx.Exec(ctx, sql, amount, userID)

	return err
}

func (repo *PostgresRepository) AddOnBalanceByUserID(amount int, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(repo.ctx, defaultTimeout)
	defer cancel()

	sql := "UPDATE balances SET amount = amount + $1 WHERE user_id = $2"

	_, err := repo.tx.Exec(ctx, sql, amount, userID)

	return err
}
