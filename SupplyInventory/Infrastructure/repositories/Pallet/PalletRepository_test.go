package infrastructure

import (
	"database/sql/driver"
	"regexp"
	"testing"

	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	entities "go_inventory/SupplyInventory/Domain/Entities"
)

// anyValue is used to match the auto-generated primary key value
type anyValue struct{}

func (a anyValue) Match(v driver.Value) bool { return true }

func setupMockRepo(t *testing.T) (sqlmock.Sqlmock, dbadapter.DBAdapter, func()) {
	sqlDB, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		sqlDB.Close()
		t.Fatalf("failed to open gorm with sqlmock: %v", err)
	}

	dbAdapter := dbadapter.NewGormAdapter(gdb)
	cleanup := func() { sqlDB.Close() }
	return sqlMock, dbAdapter, cleanup
}

// TestPalletRepository_Create tests the Create method of PalletRepository.
// It verifies that a new pallet entity is successfully inserted into the database
// using a mocked SQL database adapter. The test sets up expectations for a transaction
// (begin, insert, commit) and asserts that no errors occur during creation.
// The cleanup function returned by setupMockRepo ensures that mock resources,
// such as database connections or temporary states, are properly released after the test,
// preventing resource leaks and ensuring test isolation.
func TestPalletRepository_Create(t *testing.T) {
	// Set
	sqlMock, dbAdapter, cleanup := setupMockRepo(t)
	defer cleanup()
	repository := NewPalletRepository(dbAdapter)

	// Expectations
	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `pallets` (`name`,`pallet_rack_id`,`qr_code`,`qr_code_url`) VALUES (?,?,?,?)")).WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectCommit()

	// Actions
	palletEntity := &entities.PalletEntity{Name: "TestPallet", PalletRackID: 1}
	err := repository.Create(palletEntity)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestPalletRepository_FindByID(t *testing.T) {
	// Set
	sqlMock, dbAdapter, cleanup := setupMockRepo(t)
	defer cleanup()
	repository := NewPalletRepository(dbAdapter)

	// Expectations
	rows := sqlmock.NewRows([]string{"id", "name", "pallet_rack_id", "qr_code", "qr_code_url"}).AddRow(1, "TestPallet", 1, "", "")
	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `pallets` WHERE `pallets`.`id` = ? ORDER BY `pallets`.`id` LIMIT ?")).WillReturnRows(rows)

	// Actions
	foundPallet, err := repository.FindByID(1)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "TestPallet", foundPallet.Name)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestPalletRepository_List(t *testing.T) {
	// Set
	sqlMock, dbAdapter, cleanup := setupMockRepo(t)
	defer cleanup()
	repository := NewPalletRepository(dbAdapter)

	// Expectations
	rowsList := sqlmock.NewRows([]string{"id", "name", "pallet_rack_id", "qr_code", "qr_code_url"}).AddRow(1, "TestPallet", 1, "", "")
	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `pallets`")).WillReturnRows(rowsList)

	// Actions
	pallets, err := repository.List()

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, pallets, 1)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestPalletRepository_DeleteByID(t *testing.T) {
	// Set
	sqlMock, dbAdapter, cleanup := setupMockRepo(t)
	defer cleanup()
	repository := NewPalletRepository(dbAdapter)

	// Expectations
	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("DELETE FROM `pallets` WHERE `pallets`.`id` = ?")).WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectCommit()

	// Actions
	err := repository.DeleteByID(1)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
