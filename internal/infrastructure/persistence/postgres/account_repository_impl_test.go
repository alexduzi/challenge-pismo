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

func setupAccountRepoTest(t *testing.T) (*AccountRepositoryImpl, sqlmock.Sqlmock, *logger.MockLogger) {
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	mockLogger := logger.NewMockLogger(t)

	repo := NewAccountRepository(sqlxDB, mockLogger)

	t.Cleanup(func() {
		db.Close()
	})

	return repo, mockDB, mockLogger
}

func TestAccountRepository_GetAll_Success(t *testing.T) {
	repo, mockDB, _ := setupAccountRepoTest(t)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"account_id", "document_number", "full_name", "email", "phone", "account_type", "balance", "created_at", "updated_at"}).
		AddRow(1, "12345678900", "John Doe", "john@email.com", "123456789", "savings", 100.0, now, now).
		AddRow(2, "98765432100", "Jane Doe", "jane@email.com", "987654321", "checking", 200.0, now, now)

	mockDB.ExpectQuery(`SELECT \* FROM accounts`).WillReturnRows(rows)

	accounts, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, accounts, 2)
	assert.Equal(t, int64(1), accounts[0].AccountID)
	assert.Equal(t, "John Doe", accounts[0].FullName)
	assert.Equal(t, int64(2), accounts[1].AccountID)
	assert.Equal(t, "Jane Doe", accounts[1].FullName)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAccountRepository_GetAll_Empty(t *testing.T) {
	repo, mockDB, _ := setupAccountRepoTest(t)

	rows := sqlmock.NewRows([]string{"account_id", "document_number", "full_name", "email", "phone", "account_type", "balance", "created_at", "updated_at"})

	mockDB.ExpectQuery(`SELECT \* FROM accounts`).WillReturnRows(rows)

	accounts, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, accounts)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAccountRepository_GetAll_Error(t *testing.T) {
	repo, mockDB, _ := setupAccountRepoTest(t)

	mockDB.ExpectQuery(`SELECT \* FROM accounts`).WillReturnError(errors.New("database error"))

	accounts, err := repo.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, accounts)
	assert.ErrorIs(t, err, exception.ErrDatabaseError)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAccountRepository_GetByID_Success(t *testing.T) {
	repo, mockDB, mockLogger := setupAccountRepoTest(t)

	mockLogger.On("Debug", "Executing SELECT query", mock.Anything).Return()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"account_id", "document_number", "full_name", "email", "phone", "account_type", "balance", "created_at", "updated_at"}).
		AddRow(1, "12345678900", "John Doe", "john@email.com", "123456789", "savings", 100.0, now, now)

	mockDB.ExpectQuery(`SELECT \* FROM accounts WHERE account_id = \$1`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	account, err := repo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, int64(1), account.AccountID)
	assert.Equal(t, "12345678900", account.DocumentNumber)
	assert.Equal(t, "John Doe", account.FullName)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAccountRepository_GetByID_NotFound(t *testing.T) {
	repo, mockDB, mockLogger := setupAccountRepoTest(t)

	mockLogger.On("Debug", "Executing SELECT query", mock.Anything).Return()
	mockLogger.On("Error", "SELECT query failed", mock.Anything).Return()

	mockDB.ExpectQuery(`SELECT \* FROM accounts WHERE account_id = \$1`).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	account, err := repo.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, account)
	assert.ErrorIs(t, err, exception.ErrNotFound)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAccountRepository_GetByID_DatabaseError(t *testing.T) {
	repo, mockDB, mockLogger := setupAccountRepoTest(t)

	mockLogger.On("Debug", "Executing SELECT query", mock.Anything).Return()
	mockLogger.On("Error", "SELECT query failed", mock.Anything).Return()

	mockDB.ExpectQuery(`SELECT \* FROM accounts WHERE account_id = \$1`).
		WithArgs(int64(1)).
		WillReturnError(errors.New("connection refused"))

	account, err := repo.GetByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, account)
	assert.ErrorIs(t, err, exception.ErrDatabaseError)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAccountRepository_Save_Success(t *testing.T) {
	repo, mockDB, mockLogger := setupAccountRepoTest(t)

	mockLogger.On("Debug", "Executing INSERT query", mock.Anything).Return()

	now := time.Now()
	account := domain.Account{
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@email.com",
		Phone:          "123456789",
		AccountType:    "savings",
	}

	rows := sqlmock.NewRows([]string{"account_id", "created_at", "updated_at"}).
		AddRow(1, now, now)

	mockDB.ExpectQuery(`INSERT INTO accounts`).
		WillReturnRows(rows)

	savedAccount, err := repo.Save(context.Background(), account)

	assert.NoError(t, err)
	assert.NotNil(t, savedAccount)
	assert.Equal(t, int64(1), savedAccount.AccountID)
	assert.Equal(t, "John Doe", savedAccount.FullName)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAccountRepository_Save_Error(t *testing.T) {
	repo, mockDB, mockLogger := setupAccountRepoTest(t)

	mockLogger.On("Debug", "Executing INSERT query", mock.Anything).Return()
	mockLogger.On("Error", "INSERT query failed", mock.Anything).Return()

	account := domain.Account{
		DocumentNumber: "12345678900",
		FullName:       "John Doe",
		Email:          "john@email.com",
		Phone:          "123456789",
		AccountType:    "savings",
	}

	mockDB.ExpectQuery(`INSERT INTO accounts`).
		WillReturnError(errors.New("insert failed"))

	savedAccount, err := repo.Save(context.Background(), account)

	assert.Error(t, err)
	assert.Nil(t, savedAccount)
	assert.ErrorIs(t, err, exception.ErrDatabaseError)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAccountRepository_Update_Success(t *testing.T) {
	repo, mockDB, _ := setupAccountRepoTest(t)

	now := time.Now()
	account := domain.Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
		FullName:       "John Updated",
		Email:          "john.updated@email.com",
		Phone:          "987654321",
		AccountType:    "checking",
	}

	rows := sqlmock.NewRows([]string{"updated_at"}).AddRow(now)

	mockDB.ExpectQuery(`UPDATE accounts`).WillReturnRows(rows)

	updatedAccount, err := repo.Update(context.Background(), account)

	assert.NoError(t, err)
	assert.NotNil(t, updatedAccount)
	assert.Equal(t, "John Updated", updatedAccount.FullName)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAccountRepository_Update_Error(t *testing.T) {
	repo, mockDB, _ := setupAccountRepoTest(t)

	account := domain.Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
		FullName:       "John Updated",
		Email:          "john.updated@email.com",
		Phone:          "987654321",
		AccountType:    "checking",
	}

	mockDB.ExpectQuery(`UPDATE accounts`).WillReturnError(errors.New("update failed"))

	updatedAccount, err := repo.Update(context.Background(), account)

	assert.Error(t, err)
	assert.Nil(t, updatedAccount)
	assert.ErrorIs(t, err, exception.ErrDatabaseError)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAccountRepository_Delete_Success(t *testing.T) {
	repo, mockDB, _ := setupAccountRepoTest(t)

	mockDB.ExpectExec(`DELETE FROM accounts WHERE account_id = \$1`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAccountRepository_Delete_Error(t *testing.T) {
	repo, mockDB, _ := setupAccountRepoTest(t)

	mockDB.ExpectExec(`DELETE FROM accounts WHERE account_id = \$1`).
		WithArgs(int64(1)).
		WillReturnError(errors.New("delete failed"))

	err := repo.Delete(context.Background(), 1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, exception.ErrDatabaseError)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}
