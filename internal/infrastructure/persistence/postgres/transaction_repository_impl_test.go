package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexduzi/challengepismo/internal/domain"
	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/alexduzi/challengepismo/internal/infrastructure/logger"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTransactionRepoTest(t *testing.T) (*TransactionRepositoryImpl, sqlmock.Sqlmock, *logger.MockLogger) {
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	mockLogger := logger.NewMockLogger(t)

	repo := NewTransactionRepository(sqlxDB, mockLogger)

	t.Cleanup(func() {
		db.Close()
	})

	return repo, mockDB, mockLogger
}

func TestTransactionRepository_GetAll_Success(t *testing.T) {
	repo, mockDB, _ := setupTransactionRepoTest(t)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"transaction_id", "account_id", "operation_type_id", "amount", "event_date", "created_at"}).
		AddRow(1, 1, 1, -50.0, now, now).
		AddRow(2, 1, 4, 100.0, now, now)

	mockDB.ExpectQuery(`SELECT \* FROM transactions`).WillReturnRows(rows)

	transactions, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	assert.Equal(t, int64(1), transactions[0].TransactionID)
	assert.Equal(t, -50.0, transactions[0].Amount)
	assert.Equal(t, int64(2), transactions[1].TransactionID)
	assert.Equal(t, 100.0, transactions[1].Amount)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTransactionRepository_GetAll_Empty(t *testing.T) {
	repo, mockDB, _ := setupTransactionRepoTest(t)

	rows := sqlmock.NewRows([]string{"transaction_id", "account_id", "operation_type_id", "amount", "event_date", "created_at"})

	mockDB.ExpectQuery(`SELECT \* FROM transactions`).WillReturnRows(rows)

	transactions, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, transactions)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTransactionRepository_GetAll_Error(t *testing.T) {
	repo, mockDB, _ := setupTransactionRepoTest(t)

	mockDB.ExpectQuery(`SELECT \* FROM transactions`).WillReturnError(errors.New("database error"))

	transactions, err := repo.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, transactions)
	assert.ErrorIs(t, err, exception.ErrDatabaseError)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTransactionRepository_GetByID_Success(t *testing.T) {
	repo, mockDB, _ := setupTransactionRepoTest(t)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"transaction_id", "account_id", "operation_type_id", "amount", "event_date", "created_at"}).
		AddRow(1, 1, 4, 100.0, now, now)

	mockDB.ExpectQuery(`SELECT \* FROM transactions WHERE transaction_id = \$1`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	transaction, err := repo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, int64(1), transaction.TransactionID)
	assert.Equal(t, int64(1), transaction.AccountID)
	assert.Equal(t, 4, transaction.OperationTypeID)
	assert.Equal(t, 100.0, transaction.Amount)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTransactionRepository_GetByID_NotFound(t *testing.T) {
	repo, mockDB, _ := setupTransactionRepoTest(t)

	mockDB.ExpectQuery(`SELECT \* FROM transactions WHERE transaction_id = \$1`).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	transaction, err := repo.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.ErrorIs(t, err, exception.ErrNotFound)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTransactionRepository_GetByID_DatabaseError(t *testing.T) {
	repo, mockDB, _ := setupTransactionRepoTest(t)

	mockDB.ExpectQuery(`SELECT \* FROM transactions WHERE transaction_id = \$1`).
		WithArgs(int64(1)).
		WillReturnError(errors.New("connection refused"))

	transaction, err := repo.GetByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.ErrorIs(t, err, exception.ErrDatabaseError)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTransactionRepository_Save_Success(t *testing.T) {
	repo, mockDB, mockLogger := setupTransactionRepoTest(t)

	mockLogger.On("Debug", "Executing SELECT query", mock.Anything).Return()

	now := time.Now()
	transaction := domain.Transaction{
		AccountID:       1,
		OperationTypeID: 4,
		Amount:          100.0,
	}

	rows := sqlmock.NewRows([]string{"transaction_id", "event_date", "created_at"}).
		AddRow(1, now, now)

	mockDB.ExpectQuery(`INSERT INTO transactions`).WillReturnRows(rows)

	savedTransaction, err := repo.Save(context.Background(), transaction)

	assert.NoError(t, err)
	assert.NotNil(t, savedTransaction)
	assert.Equal(t, int64(1), savedTransaction.TransactionID)
	assert.Equal(t, int64(1), savedTransaction.AccountID)
	assert.Equal(t, 100.0, savedTransaction.Amount)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTransactionRepository_Save_Error(t *testing.T) {
	repo, mockDB, mockLogger := setupTransactionRepoTest(t)

	mockLogger.On("Debug", "Executing SELECT query", mock.Anything).Return()
	mockLogger.On("Error", "INSERT query failed", mock.Anything).Return()

	transaction := domain.Transaction{
		AccountID:       1,
		OperationTypeID: 4,
		Amount:          100.0,
	}

	mockDB.ExpectQuery(`INSERT INTO transactions`).WillReturnError(errors.New("insert failed"))

	savedTransaction, err := repo.Save(context.Background(), transaction)

	assert.Error(t, err)
	assert.Nil(t, savedTransaction)
	assert.ErrorIs(t, err, exception.ErrDatabaseError)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTransactionRepository_Update_Success(t *testing.T) {
	repo, mockDB, _ := setupTransactionRepoTest(t)

	transaction := domain.Transaction{
		TransactionID:   1,
		AccountID:       1,
		OperationTypeID: 2,
		Amount:          -75.0,
	}

	rows := sqlmock.NewRows([]string{})

	mockDB.ExpectQuery(`UPDATE transactions`).WillReturnRows(rows)

	err := repo.Update(context.Background(), transaction)

	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTransactionRepository_Update_Error(t *testing.T) {
	repo, mockDB, _ := setupTransactionRepoTest(t)

	transaction := domain.Transaction{
		TransactionID:   1,
		AccountID:       1,
		OperationTypeID: 2,
		Amount:          -75.0,
	}

	mockDB.ExpectQuery(`UPDATE transactions`).WillReturnError(errors.New("update failed"))

	err := repo.Update(context.Background(), transaction)

	assert.Error(t, err)
	assert.ErrorIs(t, err, exception.ErrDatabaseError)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTransactionRepository_Delete_Success(t *testing.T) {
	repo, mockDB, _ := setupTransactionRepoTest(t)

	mockDB.ExpectExec(`DELETE FROM transactions WHERE transaction_id = \$1`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTransactionRepository_Delete_Error(t *testing.T) {
	repo, mockDB, _ := setupTransactionRepoTest(t)

	mockDB.ExpectExec(`DELETE FROM transactions WHERE transaction_id = \$1`).
		WithArgs(int64(1)).
		WillReturnError(errors.New("delete failed"))

	err := repo.Delete(context.Background(), 1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, exception.ErrDatabaseError)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}
