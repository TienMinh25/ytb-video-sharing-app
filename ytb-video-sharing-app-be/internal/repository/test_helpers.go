package repository

import (
	"go.uber.org/mock/gomock"
	"testing"
	"ytb-video-sharing-app-be/mocks"
)

type testConfig struct {
	ctr  *gomock.Controller
	db   *mocks.MockDatabase
	tx   *mocks.MockTx
	row  *mocks.MockRow
	rows *mocks.MockRows
}

type MockSQLResult struct {
	LastInsertID int64
	RowAffected  int64
}

func (m *MockSQLResult) LastInsertId() (int64, error) {
	return m.LastInsertID, nil
}

func (m *MockSQLResult) RowsAffected() (int64, error) {
	return m.RowAffected, nil
}

func SetupTest(t *testing.T) *testConfig {
	ctr := gomock.NewController(t)
	db := mocks.NewMockDatabase(ctr)
	tx := mocks.NewMockTx(ctr)
	row := mocks.NewMockRow(ctr)
	rows := mocks.NewMockRows(ctr)

	return &testConfig{ctr: ctr, db: db, tx: tx, row: row, rows: rows}
}

func (cfg *testConfig) TearDownTest() {
	cfg.ctr.Finish()
}
